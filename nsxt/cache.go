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
	"sync/atomic"

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

// postWritePhase disables cache reads after a write (Create/Update/Delete).
// Post-write Reads always use a direct GET to pick up the just-written state.
var postWritePhase atomic.Bool

// errCacheUseBackendDirect signals that the bucket exists but the requested ID is absent;
// the caller must fall back to a direct API GET.
var errCacheUseBackendDirect = errors.New("nsxt cache: use direct API read")

// compositeCacheEntry pairs a parent type with a child search + merge step.
// Used when NSX search omits embedded children (e.g. GatewayPolicy rules).
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

// isConfigScopedCacheMode reports whether tag-filtered caching is active.
// Provider-managed tags (nsx-tf/tf-run-id) are only injected in this mode.
func isConfigScopedCacheMode() bool {
	return currentCacheMode() == cacheConfigScoped
}

func isRefreshPhase(d *schema.ResourceData) bool {
	return d.Id() != "" && !postWritePhase.Load()
}

func isCacheEnabledForRead(d *schema.ResourceData) bool {
	return IsCacheEnabled() && isRefreshPhase(d)
}

func converListToMap(list []*data.StructValue) map[string]*data.StructValue {
	converter := bindings.NewTypeConverter()
	ret := make(map[string]*data.StructValue)
	for _, obj := range list {
		dataValue, errors := converter.ConvertToGolang(obj, model.PolicyConfigResourceBindingType())
		if len(errors) > 0 {
			return nil
		}
		resource := dataValue.(model.PolicyConfigResource)
		if resource.Id != nil {
			ret[*resource.Id] = obj
		}
		if resource.DisplayName != nil {
			name := *resource.DisplayName
			if existing, seen := ret[name]; seen {
				// Duplicate display_name: mark ambiguous with nil so the caller
				// falls back to a direct GET (which returns a proper multi-match error).
				if existing != nil {
					ret[name] = nil
				}
				// Already nil (3rd+ duplicate) — leave as nil.
			} else {
				ret[name] = obj
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

// getEffectiveCacheContext returns the correct SessionContext for cache key derivation.
// Resources without a context{} block (e.g. TransitGatewayAttachment) encode their
// scope in parent_path; getSessionContext alone would return Local for those.
// Checking parent_path first mirrors what the resource's own CRUD code does.
func getEffectiveCacheContext(d *schema.ResourceData, m interface{}) utl.SessionContext {
	if pp, ok := d.GetOk("parent_path"); ok {
		if parentPath := pp.(string); parentPath != "" {
			return getParentContext(d, m, parentPath)
		}
	}
	return getSessionContext(d, m)
}

func getCacheQueryKey(resourceType string, d *schema.ResourceData, m interface{}) string {
	context := getEffectiveCacheContext(d, m)
	query := getQueryString(resourceType, context)
	runID := m.(nsxtClients).CommonConfig.contextID
	additionalQuery := buildTagQuery(d, runID)
	if additionalQuery == "" {
		return query
	}
	return fmt.Sprintf("%s AND %s ", query, additionalQuery)
}

func (c *typeScopedCache) readCache(resourceID string, resourceType string, d *schema.ResourceData, m interface{}, connector client.Connector) (interface{}, error) {
	tc := c.getTypeCache(resourceType)

	tc.mu.Lock()
	defer tc.mu.Unlock()

	query := getCacheQueryKey(resourceType, d, m)
	val, qerr := tc.getQueryResult(query, resourceID)
	if val != nil {
		tc.cacheHit += 1
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
	if err != nil && len(resultList) == 0 {
		return fmt.Errorf("error listing resource %s %w", resourceType, err)
	}
	if err != nil {
		log.Printf("[WARNING] Partial search results for resourceType=%s query=%q: %d parent results before error: %v", resourceType, query, len(resultList), err) //nolint:gosec
	}

	entry, composite := compositeCacheRegistry[resourceType]
	if !composite {
		tmp := converListToMap(resultList)
		if tmp == nil {
			return fmt.Errorf("error converting resources to cache map for resource type %s", resourceType)
		}
		c.data[query] = tmp
		return nil
	}

	// Rules may not carry user tags, so scope child search by provider-managed tags only.
	// If that yields nothing (rules untagged in NSX), retry without any tag filter.
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
	tmp := converListToMap(mergedSVs)
	if tmp == nil {
		return fmt.Errorf("error converting merged resources to cache map for resource type %s", resourceType)
	}
	c.data[query] = tmp
	return nil
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
	cacheUsed := false
	cacheAttempted := false
	if isRefreshPhase(d) && IsCacheEnabled() {
		cacheAttempted = true
		val, err := gcache.readCache(resourceID, resourceType, d, m, connector)
		if err == nil {
			converter := bindings.NewTypeConverter()
			goVal, convErrs := converter.ConvertToGolang(val.(*data.StructValue), bindingType)
			if len(convErrs) == 0 {
				typedVal, ok := goVal.(T)
				if ok {
					cacheUsed = true
					// Remove provider-managed tags before returning; they must not appear in TF state.
					if !isGlobalSearchCacheMode() {
						stripProviderManagedTagsFromAny(&typedVal)
					}
					return &typedVal, cacheUsed, cacheAttempted, nil
				}
				log.Printf("[WARNING] Cache: type assertion failed for resourceType=%s id=%s; discarding cached value, falling back to direct GET", resourceType, resourceID) //nolint:gosec
			} else {
				log.Printf("[WARNING] Cache: conversion failed for resourceType=%s id=%s (%v); discarding cached value, falling back to direct GET", resourceType, resourceID, convErrs[0]) //nolint:gosec
			}
		}
	}

	log.Printf("[DEBUG] Cache fallback: direct GET for resourceType=%s id=%s", resourceType, resourceID) //nolint:gosec
	obj, err := backendRead()
	if err != nil {
		return nil, cacheUsed, cacheAttempted, err
	}

	// Only stamp provider-managed tags in config_scope mode; other modes have no tag ownership.
	if isConfigScopedCacheMode() {
		_, patchErr := ensureProviderManagedTagsWithPatchFunc(obj, m, patchFunc)
		if patchErr != nil {
			log.Printf("[WARNING] Failed to patch provider-managed tags for %s %s: %v", resourceType, resourceID, patchErr) //nolint:gosec
		}
	}
	// Strip provider-managed tags from state; global mode returns tags as-is.
	if !isGlobalSearchCacheMode() {
		stripProviderManagedTagsFromAny(obj)
	}

	return obj, cacheUsed, cacheAttempted, nil
}

// providerManagedTagsSearchFragments returns NSX search tag clauses for provider-managed
// tag nsx-tf/tf-run-id (value = context_id). Empty when runID is unset.
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

// buildTagQuery builds the NSX search tag filter from the resource's tag block.
// In config_scope mode, provider-managed tags are appended automatically if absent.
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

// InvalidateCacheForResourceType clears the cache bucket for a resource type and
// sets postWritePhase so that the immediately following Read uses a direct GET.
// Call after every Create, Update, and Delete.
func InvalidateCacheForResourceType(resourceType string) {
	if !IsCacheEnabled() {
		return
	}
	tc := gcache.getTypeCache(resourceType)
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.data = make(map[string]map[string]*data.StructValue)
	postWritePhase.Store(true)
	log.Printf("[DEBUG] Cache invalidated for resourceType=%s; post-write phase set, cache bypassed for remaining reads", resourceType) //nolint:gosec
}

// stripProviderManagedTagsFromAny removes nsx-tf/tf-run-id tags from any NSX model object
// via reflection. Provider-managed tags must not appear in TF state to prevent spurious diffs.
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

// ensureProviderManagedTagsWithPatchFunc stamps provider-managed tags onto obj via patchFunc if missing.
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
			// Bump Revision to match the server-side increment from the PATCH;
			// without this the next apply sends a stale revision and NSX rejects it (error 500071).
			revField := objValue.FieldByName("Revision")
			if revField.IsValid() && revField.Kind() == reflect.Pointer && !revField.IsNil() {
				revField.Elem().SetInt(revField.Elem().Int() + 1)
			}
		}

		return obj, nil
	}

	return nil, fmt.Errorf("Provider managed tags field is not settable")
}
