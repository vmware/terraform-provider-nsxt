package nsxt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type resourceTypeCache struct {
	mu       sync.RWMutex
	data     map[string]map[string]*data.StructValue // key: query, key: resource ID
	cacheHit int
	cacheMis int
}

type typeScopedCache struct {
	mu    sync.RWMutex
	byTyp map[string]*resourceTypeCache
}

var gcache = &typeScopedCache{byTyp: make(map[string]*resourceTypeCache)}

// postWriteByKey tracks recently written resources to bypass cache once on the next read.
var postWriteByKey sync.Map // string -> struct{}


// errCacheUseBackendDirect indicates the cache bucket exists but the ID is missing, so use direct GET.
var errCacheUseBackendDirect = errors.New("nsxt cache: use direct API read")

// compositeCacheEntry defines a parent type with an extra child search and merge step.
type compositeCacheEntry struct {
	childSearchType string
	merge           func(parents, children []*data.StructValue) ([]*data.StructValue, error)
}

var compositeCacheRegistry = map[string]compositeCacheEntry{
	"gatewaypolicy": {
		childSearchType: "rule",
		merge:           mergeGatewayPolicyCacheSearchResults,
	},
	"securitypolicy": {
		childSearchType: "rule",
		merge:           mergeSecurityPolicyCacheSearchResults,
	},
}

const envNSXTCacheMode = "NSXT_CACHE_MODE"

type cacheMode int

const (
	cacheDisabled cacheMode = iota
	cacheConfigScoped
	cacheGlobal
)

var invalidNSXTCacheModeWarn sync.Map // raw env value -> struct{}; log once per distinct invalid value

func postWriteKey(resourceType, resourceID string) string {
	return resourceType + "::" + resourceID
}

// MarkPostWriteForResourceTypeKey marks a single resource instance as just written.
// The next CacheAwareResourceRead/TryCacheRead for that exact resourceID will bypass cache once.
func MarkPostWriteForResourceTypeKey(resourceType, resourceID string) {
	if !IsCacheEnabled() {
		return
	}
	if resourceType == "" || resourceID == "" {
		return
	}
	postWriteByKey.Store(postWriteKey(resourceType, resourceID), struct{}{})
}

// MarkPostWriteAndInvalidateCacheForResourceType marks resourceID as just written and invalidates
// the type's cache bucket. resourceID must be the same key used in the corresponding
// CacheAwareResourceRead/TryCacheRead call (typically d.Id(), or the parent's ID for child resources).
func MarkPostWriteAndInvalidateCacheForResourceType(resourceType, resourceID string) {
	MarkPostWriteForResourceTypeKey(resourceType, resourceID)
	InvalidateCacheForResourceType(resourceType)
}

func currentCacheMode() cacheMode {
	raw := strings.TrimSpace(os.Getenv(envNSXTCacheMode))
	lower := strings.ToLower(raw)
	switch lower {
	case "", "disabled", "off":
		return cacheDisabled
	case "config_scope":
		return cacheConfigScoped
	case "global":
		return cacheGlobal
	default:
		if raw != "" {
			if _, loaded := invalidNSXTCacheModeWarn.LoadOrStore(raw, struct{}{}); !loaded {
				// strconv.Quote neutralizes newlines/control chars in env (gosec G706).
				log.Printf("[ERROR]: Invalid %s value %s; expected disabled, off, config_scope, or global (caching disabled)", envNSXTCacheMode, strconv.Quote(raw)) //nolint:gosec
			}
		}
		return cacheDisabled
	}
}

func IsCacheEnabled() bool {
	return currentCacheMode() != cacheDisabled
}

func isGlobalSearchCacheMode() bool {
	return currentCacheMode() == cacheGlobal
}

// isConfigScopedCacheMode reports whether tag-filtered caching (with provider-managed tags) is active.
func isConfigScopedCacheMode() bool {
	return currentCacheMode() == cacheConfigScoped
}

func isRefreshPhase(d *schema.ResourceData) bool {
	return d.Id() != ""
}

func isCacheEnabledForRead(d *schema.ResourceData) bool {
	return IsCacheEnabled() && isRefreshPhase(d)
}

// CacheKeyForResourceID returns the key that CacheAwareResourceRead/TryCacheRead will use
// for this resource. Path-indexed types key by path (when set); all others key by id.
// Use this when calling MarkPostWriteAndInvalidateCacheForResourceType to ensure the keys match.
func CacheKeyForResourceID(resourceType string, d *schema.ResourceData) string {
	if shouldIndexByPath(resourceType) {
		if p, ok := d.GetOk("path"); ok {
			if path, ok2 := p.(string); ok2 && path != "" {
				return path
			}
		}
	}
	return d.Id()
}

func shouldIndexByPath(resourceType string) bool {
	switch resourceType {
	case "PolicyNatRule", "SegmentPort", "IpAddressPoolStaticSubnet", "IpAddressPoolBlockSubnet", "Service":
		return true
	default:
		return false
	}
}

// indexCacheMapKey inserts key→obj if key is not already present (first-seen wins).
// Used for both primary keys (id) and secondary keys (display_name, path).
func indexCacheMapKey(ret map[string]*data.StructValue, key string, obj *data.StructValue) {
	if _, seen := ret[key]; seen {
		return
	}
	ret[key] = obj
}

func getStructValueStringField(obj *data.StructValue, field string) (string, bool) {
	if obj == nil {
		return "", false
	}
	fields := obj.Fields()
	v, ok := fields[field]
	if !ok || v == nil {
		return "", false
	}
	if sv, ok := v.(*data.StringValue); ok {
		return sv.Value(), true
	}
	if ov, ok := v.(*data.OptionalValue); ok {
		if !ov.IsSet() {
			return "", false
		}
		inner := ov.Value()
		if sv, ok := inner.(*data.StringValue); ok {
			return sv.Value(), true
		}
	}
	return "", false
}

func converListToMapByType(list []*data.StructValue, resourceType string) map[string]*data.StructValue {
	converter := bindings.NewTypeConverter()
	ret := make(map[string]*data.StructValue)
	for _, obj := range list {
		dataValue, errors := converter.ConvertToGolang(obj, model.PolicyConfigResourceBindingType())
		if len(errors) > 0 {
			return nil
		}
		resource := dataValue.(model.PolicyConfigResource)

		// Primary indexing via PolicyConfigResource conversion.
		if resource.Id != nil {
			id := *resource.Id
			indexCacheMapKey(ret, id, obj)
			// Index both raw and extracted policy IDs when Search returns a full path.
			if strings.Contains(id, "/") {
				indexCacheMapKey(ret, getPolicyIDFromPath(id), obj)
			}
		}
		if resource.DisplayName != nil {
			indexCacheMapKey(ret, *resource.DisplayName, obj)
		}
		if shouldIndexByPath(resourceType) && resource.Path != nil {
			indexCacheMapKey(ret, *resource.Path, obj)
		}

		// Also index from raw StructValue fields to cover cases where SDK conversion loses fields.
		if id, ok := getStructValueStringField(obj, "id"); ok {
			indexCacheMapKey(ret, id, obj)
			if strings.Contains(id, "/") {
				indexCacheMapKey(ret, getPolicyIDFromPath(id), obj)
			}
		}
		if dn, ok := getStructValueStringField(obj, "display_name"); ok {
			indexCacheMapKey(ret, dn, obj)
		}
		if shouldIndexByPath(resourceType) {
			if p, ok := getStructValueStringField(obj, "path"); ok {
				indexCacheMapKey(ret, p, obj)
			}
		}
	}
	return ret
}

func getQueryString(resourceType string, context utl.SessionContext) string {
	switch context.ClientType {
	case utl.Global:
		return fmt.Sprintf("resource_type:%s AND marked_for_delete:false AND context:Global", resourceType)
	case utl.Local:
		return fmt.Sprintf("resource_type:%s AND marked_for_delete:false AND context:Local", resourceType)
	case utl.VPC, utl.Multitenancy:
		return fmt.Sprintf("resource_type:%s AND marked_for_delete:false AND context:%s-%s", resourceType, context.ProjectID, context.VPCID)
	default:
		return fmt.Sprintf("resource_type:%s AND marked_for_delete:false", resourceType)
	}
}

func (c *resourceTypeCache) getQueryResult(query string, resourceID string) (*data.StructValue, error) {
	if inner, ok := c.data[query]; ok {
		if v := inner[resourceID]; v != nil {
			return v, nil
		}
		// As a safety net, scan the bucket for a matching raw "id" field.
		for _, v := range inner {
			if v == nil {
				continue
			}
			if id, ok := getStructValueStringField(v, "id"); ok && id == resourceID {
				return v, nil
			}
		}
		return nil, errCacheUseBackendDirect
	}
	return nil, fmt.Errorf("element is not found")
}

func (c *resourceTypeCache) writeCache(query string, resourceType string, d *schema.ResourceData, m interface{}, connector client.Connector) error {
	if _, ok := c.data[query]; ok {
		log.Printf("[DEBUG] Cache skip bulk refill for resourceType=%s query=%q (bucket already present)", resourceType, query) //nolint:gosec
		return nil
	}
	runID := m.(nsxtClients).CommonConfig.contextID
	log.Printf("[DEBUG] Cache miss: populating cache for resourceType=%s query=%q", resourceType, query) //nolint:gosec
	err := c.getListOfPolicyResources(query, d, connector, getEffectiveCacheContext(d, m), resourceType, runID)
	if err != nil {
		return err
	}
	return nil
}

func (c *typeScopedCache) getTypeCache(resourceType string) *resourceTypeCache {
	c.mu.Lock()
	defer c.mu.Unlock()
	if tc, ok := c.byTyp[resourceType]; ok {
		return tc
	}
	tc := &resourceTypeCache{data: make(map[string]map[string]*data.StructValue)}
	c.byTyp[resourceType] = tc
	return tc
}

// getEffectiveCacheContext derives the cache context, preferring parent_path when context{} is absent.
func getEffectiveCacheContext(d *schema.ResourceData, m interface{}) utl.SessionContext {
	if pp, ok := d.GetOk("parent_path"); ok {
		if parentPath := pp.(string); parentPath != "" {
			return getParentContext(d, m, parentPath)
		}
	}
	return getSessionContext(d, m)
}

func getCacheQueryKey(resourceType string, d *schema.ResourceData, m interface{}) string {
	clients := m.(nsxtClients)
	context := getEffectiveCacheContext(d, m)
	query := getQueryString(resourceType, context)
	runID := clients.CommonConfig.contextID
	additionalQuery := buildTagQuery(d, runID)
	// Prefix with manager host so aliased provider blocks pointing at different managers
	// never share a cache bucket even when context_id and resource types are identical.
	host := clients.Host
	base := fmt.Sprintf("[%s]%s", host, query)
	if additionalQuery == "" {
		return base
	}
	return fmt.Sprintf("%s AND %s ", base, additionalQuery)
}

func (c *typeScopedCache) readCache(resourceID string, resourceType string, d *schema.ResourceData, m interface{}, connector client.Connector) (interface{}, error) {
	tc := c.getTypeCache(resourceType)
	query := getCacheQueryKey(resourceType, d, m)

	// Fast path: read lock allows concurrent hits without blocking each other.
	tc.mu.RLock()
	val, qerr := tc.getQueryResult(query, resourceID)
	if val != nil {
		tc.cacheHit++
		log.Printf("[DEBUG] Cache hit: resourceType=%s id=%s query=%q (hit=%d miss=%d)", resourceType, resourceID, query, tc.cacheHit, tc.cacheMis) //nolint:gosec
		tc.mu.RUnlock()
		return val, nil
	}
	if errors.Is(qerr, errCacheUseBackendDirect) {
		tc.cacheMis++
		log.Printf("[DEBUG] Cache lookup miss: resourceType=%s id=%s query=%q (hit=%d miss=%d)", resourceType, resourceID, query, tc.cacheHit, tc.cacheMis) //nolint:gosec
		tc.mu.RUnlock()
		return nil, errCacheUseBackendDirect
	}
	tc.mu.RUnlock()

	// Slow path: take exclusive lock to populate the bucket, then double-check.
	tc.mu.Lock()
	defer tc.mu.Unlock()
	val, qerr = tc.getQueryResult(query, resourceID)
	if val != nil {
		tc.cacheHit++
		log.Printf("[DEBUG] Cache hit: resourceType=%s id=%s query=%q (hit=%d miss=%d)", resourceType, resourceID, query, tc.cacheHit, tc.cacheMis) //nolint:gosec
		return val, nil
	}
	if errors.Is(qerr, errCacheUseBackendDirect) {
		tc.cacheMis++
		log.Printf("[DEBUG] Cache lookup miss: resourceType=%s id=%s query=%q (hit=%d miss=%d)", resourceType, resourceID, query, tc.cacheHit, tc.cacheMis) //nolint:gosec
		return nil, errCacheUseBackendDirect
	}
	tc.cacheMis++
	log.Printf("[DEBUG] Cache lookup miss: resourceType=%s id=%s query=%q (hit=%d miss=%d)", resourceType, resourceID, query, tc.cacheHit, tc.cacheMis) //nolint:gosec
	err := tc.writeCache(query, resourceType, d, m, connector)
	if err != nil {
		return nil, err
	}
	val, _ = tc.getQueryResult(query, resourceID)
	if val != nil {
		return val, nil
	}
	return nil, errCacheUseBackendDirect
}

func (c *resourceTypeCache) getListOfPolicyResources(query string, d *schema.ResourceData, connector client.Connector, context utl.SessionContext, resourceType string, runID string) error {
	additionalQuery := buildTagQuery(d, runID)
	log.Printf("[DEBUG] Cache search query: resourceType=%s query=%q additionalQuery=%q", resourceType, query, additionalQuery) //nolint:gosec
	resultList, err := listPolicyResources(connector, context, resourceType, &additionalQuery)
	log.Printf("[DEBUG] Cache search results: resourceType=%s query=%q additionalQuery=%q results=%d", resourceType, query, additionalQuery, len(resultList)) //nolint:gosec
	if err != nil && len(resultList) == 0 {
		return fmt.Errorf("error listing resource %s %w", resourceType, err)
	}
	if err != nil {
		log.Printf("[WARNING] Partial search results for resourceType=%s query=%q: %d parent results before error: %v", resourceType, query, len(resultList), err) //nolint:gosec
	}

	entry, composite := compositeCacheRegistry[resourceType]
	if !composite {
		tmp := converListToMapByType(resultList, resourceType)
		if tmp == nil {
			return fmt.Errorf("error converting resources to cache map for resource type %s", resourceType)
		}
		c.data[query] = tmp
		return nil
	}

	// Child rules may be untagged, so try provider-managed tags first and then fall back to no tag filter.
	var childAdditional *string
	if !isGlobalSearchCacheMode() {
		if q := providerManagedTagsSearchQuery(runID); q != "" {
			childAdditional = &q
		}
	}
	childList, childErr := listPolicyResources(connector, context, entry.childSearchType, childAdditional)
	if childErr != nil && len(childList) == 0 {
		return fmt.Errorf("error listing composite child resource %s for parent %s: %w", entry.childSearchType, resourceType, childErr)
	}
	if childErr != nil {
		log.Printf("[WARNING] Partial child search for resourceType=%s childType=%s query=%q: %d results before error: %v", resourceType, entry.childSearchType, query, len(childList), childErr) //nolint:gosec
	}
	if childAdditional != nil && len(resultList) > 0 && len(childList) == 0 {
		childList, childErr = listPolicyResources(connector, context, entry.childSearchType, nil)
		if childErr != nil && len(childList) == 0 {
			return fmt.Errorf("error listing composite child resource %s for parent %s: %w", entry.childSearchType, resourceType, childErr)
		}
		if childErr != nil {
			log.Printf("[WARNING] Partial child search (fallback) for resourceType=%s childType=%s query=%q: %d results before error: %v", resourceType, entry.childSearchType, query, len(childList), childErr) //nolint:gosec
		}
	}

	mergedSVs, err := entry.merge(resultList, childList)
	if err != nil {
		return err
	}
	tmp := converListToMapByType(mergedSVs, resourceType)
	if tmp == nil {
		return fmt.Errorf("error converting merged resources to cache map for resource type %s", resourceType)
	}
	c.data[query] = tmp
	return nil
}

// convertCachedValue converts a raw cache value to the typed model, strips provider-managed tags,
// and returns (typedPtr, true) on success or (nil, false) on conversion/type-assertion failure.
func convertCachedValue[T any](val interface{}, resourceType, resourceID string, bindingType bindings.BindingType) (*T, bool) {
	sv, ok := val.(*data.StructValue)
	if !ok {
		return nil, false
	}
	converter := bindings.NewTypeConverter()
	goVal, convErrs := converter.ConvertToGolang(sv, bindingType)
	if len(convErrs) > 0 {
		log.Printf("[WARNING] Cache: conversion failed for resourceType=%s id=%s (%v); discarding cached value", resourceType, resourceID, convErrs[0]) //nolint:gosec
		return nil, false
	}
	typedVal, ok := goVal.(T)
	if !ok {
		log.Printf("[WARNING] Cache: type assertion failed for resourceType=%s id=%s; discarding cached value", resourceType, resourceID) //nolint:gosec
		return nil, false
	}
	if !isGlobalSearchCacheMode() {
		stripProviderManagedTagsFromAny(&typedVal)
	}
	return &typedVal, true
}

// TryCacheRead reads from cache only and never falls back to a backend GET.
func TryCacheRead[T any](d *schema.ResourceData, m interface{}, connector client.Connector, resourceID string, resourceType string, bindingType bindings.BindingType) (*T, bool, bool, error) {
	cacheUsed := false
	cacheAttempted := false
	if isRefreshPhase(d) && IsCacheEnabled() {
		if _, ok := postWriteByKey.LoadAndDelete(postWriteKey(resourceType, resourceID)); ok {
			// Bypass cache after a write for this resource instance (one-shot).
			return nil, cacheUsed, true, nil
		}
		cacheAttempted = true
		val, err := gcache.readCache(resourceID, resourceType, d, m, connector)
		if err == nil {
			if typedVal, ok := convertCachedValue[T](val, resourceType, resourceID, bindingType); ok {
				return typedVal, true, cacheAttempted, nil
			}
		}
	}
	return nil, cacheUsed, cacheAttempted, nil
}

func structValuesToModels[T any](list []*data.StructValue, bt bindings.BindingType) ([]T, error) {
	converter := bindings.NewTypeConverter()
	out := make([]T, 0, len(list))
	for _, obj := range list {
		dataValue, errors := converter.ConvertToGolang(obj, bt)
		if len(errors) > 0 {
			var zero T
			return nil, fmt.Errorf("converting %T for cache: %w", zero, errors[0])
		}
		v, ok := dataValue.(T)
		if !ok {
			return nil, fmt.Errorf("converting for cache: unexpected type %T", dataValue)
		}
		out = append(out, v)
	}
	return out, nil
}

func modelsToStructValues[T any](models []T, bt bindings.BindingType) ([]*data.StructValue, error) {
	converter := bindings.NewTypeConverter()
	out := make([]*data.StructValue, 0, len(models))
	for i := range models {
		dataValue, errors := converter.ConvertToVapi(models[i], bt)
		if len(errors) > 0 {
			var zero T
			return nil, fmt.Errorf("converting %T to struct value: %w", zero, errors[0])
		}
		sv, ok := dataValue.(*data.StructValue)
		if !ok {
			return nil, fmt.Errorf("converting to struct value: expected *data.StructValue, got %T", dataValue)
		}
		out = append(out, sv)
	}
	return out, nil
}

func structValuesToRules(list []*data.StructValue) []model.Rule {
	converter := bindings.NewTypeConverter()
	out := make([]model.Rule, 0, len(list))
	for _, obj := range list {
		dataValue, errors := converter.ConvertToGolang(obj, model.RuleBindingType())
		if len(errors) > 0 {
			continue
		}
		rule, ok := dataValue.(model.Rule)
		if !ok {
			continue
		}
		out = append(out, rule)
	}
	return out
}

// groupRulesByValidParentPath indexes rules by ParentPath, restricted to known parent paths.
// Rules with a missing, empty, or unrecognized ParentPath are dropped.
func groupRulesByValidParentPath(validParentPaths map[string]struct{}, rules []model.Rule) map[string][]model.Rule {
	byParent := make(map[string][]model.Rule)
	for _, r := range rules {
		if r.ParentPath == nil {
			continue
		}
		pp := strings.TrimSpace(*r.ParentPath)
		if pp == "" {
			continue
		}
		if _, ok := validParentPaths[pp]; !ok {
			continue
		}
		byParent[pp] = append(byParent[pp], r)
	}
	return byParent
}

// attachRulesByParentPath assigns rules to the parent whose Path matches rule.ParentPath.
// Unmatched rules are dropped; parent order is preserved.
func attachRulesByParentPath[P any](parents []P, rules []model.Rule, getPath func(P) *string, setRules func(*P, []model.Rule)) []P {
	validPaths := make(map[string]struct{})
	for _, p := range parents {
		path := getPath(p)
		if path != nil {
			k := strings.TrimSpace(*path)
			if k != "" {
				validPaths[k] = struct{}{}
			}
		}
	}
	byParent := groupRulesByValidParentPath(validPaths, rules)

	out := make([]P, len(parents))
	for i := range parents {
		p := parents[i]
		key := ""
		if path := getPath(p); path != nil {
			key = strings.TrimSpace(*path)
		}
		bucket := byParent[key]
		setRules(&p, bucket)
		out[i] = p
	}
	return out
}

func mergeGatewayPolicyCacheSearchResults(parents, children []*data.StructValue) ([]*data.StructValue, error) {
	gp, err := structValuesToModels[model.GatewayPolicy](parents, model.GatewayPolicyBindingType())
	if err != nil {
		return nil, err
	}
	rules := structValuesToRules(children)
	merged := attachRulesByParentPath(gp, rules,
		func(p model.GatewayPolicy) *string { return p.Path },
		func(p *model.GatewayPolicy, r []model.Rule) { p.Rules = r },
	)
	return modelsToStructValues(merged, model.GatewayPolicyBindingType())
}

func mergeSecurityPolicyCacheSearchResults(parents, children []*data.StructValue) ([]*data.StructValue, error) {
	sp, err := structValuesToModels[model.SecurityPolicy](parents, model.SecurityPolicyBindingType())
	if err != nil {
		return nil, err
	}
	rules := structValuesToRules(children)
	merged := attachRulesByParentPath(sp, rules,
		func(p model.SecurityPolicy) *string { return p.Path },
		func(p *model.SecurityPolicy, r []model.Rule) { p.Rules = r },
	)
	return modelsToStructValues(merged, model.SecurityPolicyBindingType())
}

func CacheAwareResourceRead[T any](d *schema.ResourceData, m interface{}, connector client.Connector, resourceID string, resourceType string, bindingType bindings.BindingType, backendRead func() (*T, error), patchFunc func(obj *T) error) (*T, bool, bool, error) {
	cacheAttempted := false
	postWriteBypass := false

	if isRefreshPhase(d) && IsCacheEnabled() {
		if _, ok := postWriteByKey.LoadAndDelete(postWriteKey(resourceType, resourceID)); ok {
			// Ensure read-your-writes semantics: bypass cache once right after a write.
			// Search-backed cache may be briefly stale immediately after a write.
			postWriteBypass = true
			cacheAttempted = true
		}
		if !postWriteBypass {
			cacheAttempted = true
			val, err := gcache.readCache(resourceID, resourceType, d, m, connector)
			if err == nil {
				if typedVal, ok := convertCachedValue[T](val, resourceType, resourceID, bindingType); ok {
					return typedVal, true, cacheAttempted, nil
				}
			}
		}
	}

	if postWriteBypass {
		log.Printf("[DEBUG] Cache post-write bypass: direct GET for resourceType=%s id=%s", resourceType, resourceID) //nolint:gosec
	} else {
		log.Printf("[DEBUG] Cache fallback: direct GET for resourceType=%s id=%s", resourceType, resourceID) //nolint:gosec
	}
	obj, err := backendRead()
	if err != nil {
		return nil, false, cacheAttempted, err
	}

	// Only stamp provider-managed tags in config_scope mode.
	if isConfigScopedCacheMode() {
		_, patchErr := ensureProviderManagedTagsWithPatchFunc(obj, m, patchFunc)
		if patchErr != nil {
			log.Printf("[WARNING] Failed to patch provider-managed tags for %s %s: %v", resourceType, resourceID, patchErr) //nolint:gosec
		}
	}
	// Strip provider-managed tags from state in non-global mode.
	if !isGlobalSearchCacheMode() {
		stripProviderManagedTagsFromAny(obj)
	}

	return obj, false, cacheAttempted, nil
}

// providerManagedTagsSearchFragments returns NSX search fragments for provider-managed tags.
func providerManagedTagsSearchFragments(runID string) []string {
	if runID == "" {
		return nil
	}
	return []string{
		fmt.Sprintf("tags.scope:%s", escapeSpecialCharacters(managedDefaultTagScope)),
		fmt.Sprintf("tags.tag:%s", escapeSpecialCharacters(runID)),
	}
}

// providerManagedTagsSearchQuery is the AND-joined search fragment for provider-managed tags.
func providerManagedTagsSearchQuery(runID string) string {
	fr := providerManagedTagsSearchFragments(runID)
	if len(fr) == 0 {
		return ""
	}
	return strings.Join(fr, " AND ")
}

// buildTagQuery builds the NSX search tag filter, adding provider-managed tags when needed.
func buildTagQuery(d *schema.ResourceData, runID string) string {
	managedTagPresent := false
	if d == nil {
		return ""
	}

	shouldAddProviderTags := !isGlobalSearchCacheMode() && runID != ""

	tags, exists := d.GetOk("tag")
	if !exists {
		if shouldAddProviderTags {
			return providerManagedTagsSearchQuery(runID)
		}
		return ""
	}

	tagSet, ok := tags.(*schema.Set)
	if !ok {
		if shouldAddProviderTags {
			return providerManagedTagsSearchQuery(runID)
		}
		return ""
	}
	if tagSet.Len() == 0 {
		if shouldAddProviderTags {
			return providerManagedTagsSearchQuery(runID)
		}
		return ""
	}

	var tagQueries []string
	for _, tagInterface := range tagSet.List() {
		tagMap := tagInterface.(map[string]interface{})

		scope, hasScope := tagMap["scope"]
		tag, hasTag := tagMap["tag"]
		if hasScope && scope != nil && scope.(string) != "" {
			rawScope := scope.(string)
			if isManagedDefaultTagScope(&rawScope) {
				managedTagPresent = true
			}
			tagQueries = append(tagQueries, fmt.Sprintf("tags.scope:%s", escapeSpecialCharacters(rawScope)))
		}
		if hasTag && tag != nil && tag.(string) != "" {
			tagQueries = append(tagQueries, fmt.Sprintf("tags.tag:%s", escapeSpecialCharacters(tag.(string))))
		}
	}
	if !isGlobalSearchCacheMode() && !managedTagPresent && runID != "" {
		tagQueries = append(tagQueries, providerManagedTagsSearchFragments(runID)...)
	}

	if len(tagQueries) == 0 {
		return ""
	}

	return strings.Join(tagQueries, " AND ")
}

// InvalidateCacheForResourceType clears the cache bucket for a resource type.
func InvalidateCacheForResourceType(resourceType string) {
	if !IsCacheEnabled() {
		return
	}
	tc := gcache.getTypeCache(resourceType)
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.data = make(map[string]map[string]*data.StructValue)
	log.Printf("[DEBUG] Cache invalidated for resourceType=%s", resourceType) //nolint:gosec
}

// stripProviderManagedTagsFromAny removes nsx-tf/tf-run-id tags from an NSX model object via reflection.
func stripProviderManagedTagsFromAny(obj interface{}) {
	if obj == nil {
		return
	}
	objValue := reflect.ValueOf(obj)
	for objValue.IsValid() && objValue.Kind() == reflect.Pointer {
		if objValue.IsNil() {
			return
		}
		objValue = objValue.Elem()
	}
	if !objValue.IsValid() || objValue.Kind() != reflect.Struct {
		return
	}
	tagsField := objValue.FieldByName("Tags")
	if !tagsField.IsValid() || !tagsField.CanSet() || tagsField.Kind() != reflect.Slice {
		return
	}
	current, ok := tagsField.Interface().([]model.Tag)
	if !ok {
		return
	}
	userTags := make([]model.Tag, 0, len(current))
	for _, tag := range current {
		if !isManagedDefaultTag(tag) {
			userTags = append(userTags, tag)
		}
	}
	tagsField.Set(reflect.ValueOf(userTags))
}

// ensureProviderManagedTagsWithPatchFunc adds provider-managed tags to obj via patchFunc when missing.
func ensureProviderManagedTagsWithPatchFunc[T any](obj T, m interface{}, patchFunc func(obj T) error) (interface{}, error) {
	objValue := reflect.ValueOf(obj)
	for objValue.IsValid() && objValue.Kind() == reflect.Pointer {
		if objValue.IsNil() {
			return nil, nil
		}
		objValue = objValue.Elem()
	}
	if !objValue.IsValid() || objValue.Kind() != reflect.Struct {
		return nil, nil
	}

	tagsField := objValue.FieldByName("Tags")
	if !tagsField.IsValid() {
		return nil, nil
	}

	var currentTags []model.Tag
	if tagsField.Kind() == reflect.Slice && !tagsField.IsNil() {
		for i := 0; i < tagsField.Len(); i++ {
			tag := tagsField.Index(i).Interface().(model.Tag)
			currentTags = append(currentTags, tag)
		}
	}

	runID := m.(nsxtClients).CommonConfig.contextID
	expectedManagedTags := getProviderManagedDefaultTags(runID)

	needsPatch := false
	for _, expectedTag := range expectedManagedTags {
		found := false
		for _, currentTag := range currentTags {
			if currentTag.Scope != nil && expectedTag.Scope != nil &&
				currentTag.Tag != nil && expectedTag.Tag != nil &&
				*currentTag.Scope == *expectedTag.Scope &&
				*currentTag.Tag == *expectedTag.Tag {
				found = true
				break
			}
		}
		if !found {
			needsPatch = true
			break
		}
	}

	if !needsPatch {
		return nil, nil
	}
	log.Printf("[DEBUG] Provider-managed tag missing; patching object to add scope=%s", managedDefaultTagScope) //nolint:gosec

	userTags := make([]model.Tag, 0)
	for _, tag := range currentTags {
		if !isManagedDefaultTag(tag) {
			userTags = append(userTags, tag)
		}
	}

	mergedTags := mergeManagedDefaultAndUserTags(expectedManagedTags, userTags)

	if tagsField.CanSet() {
		tagsField.Set(reflect.ValueOf(mergedTags))

		if patchFunc != nil {
			err := patchFunc(obj)
			if err != nil {
				return nil, fmt.Errorf("failed to patch tags to NSX API: %w", err)
			}
			log.Printf("[DEBUG] Patched provider-managed tags successfully")
			// Bump Revision to match the server-side PATCH increment and avoid stale-revision errors.
			revField := objValue.FieldByName("Revision")
			if revField.IsValid() && revField.Kind() == reflect.Pointer && !revField.IsNil() {
				revField.Elem().SetInt(revField.Elem().Int() + 1)
			}
		}

		return obj, nil
	}

	return nil, fmt.Errorf("Provider managed tags field is not settable")
}
