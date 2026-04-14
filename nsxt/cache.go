package nsxt

import (
	"fmt"
	"log"
	"os"
	"reflect"
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

// compositeCacheEntry describes a parent search type that needs a second child search and merge
// before results are stored (e.g. GatewayPolicy + Rule where search omits embedded rules).
type compositeCacheEntry struct {
	childSearchType string
	merge           func(parents, children []*data.StructValue) ([]*data.StructValue, error)
}

var compositeCacheRegistry = map[string]compositeCacheEntry{
	"gatewaypolicy": {
		childSearchType: "rule",
		merge:           mergeGatewayPolicyCacheSearchResults,
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
	case "config-scope":
		return cacheConfigScoped
	case "global":
		return cacheGlobal
	default:
		if raw != "" {
			if _, loaded := invalidNSXTCacheModeWarn.LoadOrStore(raw, struct{}{}); !loaded {
				log.Printf("[WARNING] Invalid %s=%q; expected disabled, off, config-scope, or global. Caching disabled.", envNSXTCacheMode, raw)
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

func isRefreshPhase(d *schema.ResourceData) bool {
	return d.Id() != "" && !d.HasChangesExcept()
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
			// Keep DisplayName as an alias key for data source lookups.
			ret[*resource.DisplayName] = obj
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
	if val, ok := c.data[query]; ok {
		return val[resourceID], nil
	}
	return nil, fmt.Errorf("element is not found")
}

func (c *resourceTypeCache) writeCache(query string, resourceType string, d *schema.ResourceData, m interface{}, connector client.Connector) error {
	c.cacheMis += 1
	runID := m.(nsxtClients).CommonConfig.contextID
	log.Printf("[DEBUG] Cache miss: populating cache for resourceType=%s query=%q", resourceType, query)
	err := c.getListOfPolicyResources(query, d, connector, getSessionContext(d, m), resourceType, runID)
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

func getCacheQueryKey(resourceType string, d *schema.ResourceData, m interface{}) string {
	context := getSessionContext(d, m)
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
	if val, _ := tc.getQueryResult(query, resourceID); val != nil {
		tc.cacheHit += 1
		log.Printf("[DEBUG] Cache hit: resourceType=%s id=%s query=%q (hit=%d miss=%d)", resourceType, resourceID, query, tc.cacheHit, tc.cacheMis)
		return val, nil
	}
	log.Printf("[DEBUG] Cache lookup miss: resourceType=%s id=%s query=%q (hit=%d miss=%d)", resourceType, resourceID, query, tc.cacheHit, tc.cacheMis)
	err := tc.writeCache(query, resourceType, d, m, connector)
	if err != nil {
		return nil, err
	}
	return tc.getQueryResult(query, resourceID)
}

func (c *resourceTypeCache) getListOfPolicyResources(query string, d *schema.ResourceData, connector client.Connector, context utl.SessionContext, resourceType string, runID string) error {
	// Now we have access to the proper runID passed from writeCache
	additionalQuery := buildTagQuery(d, runID)
	resultList, err := listPolicyResources(connector, context, resourceType, &additionalQuery)
	if err != nil {
		return fmt.Errorf("error listing resource %s %w", resourceType, err)
	}

	entry, composite := compositeCacheRegistry[resourceType]
	if !composite {
		tmp := converListToMap(resultList)
		c.data[query] = tmp
		return nil
	}

	// Child objects (e.g. rules) must not use the full parent tag query (policy user tags
	// like color:orange vs rule tags color:blue). In tag cache mode (!global search cache),
	// still scope children by provider-managed tags when runID is set; if that returns
	// nothing (rules may not carry those tags in NSX), fall back to path-scoped list.
	var childAdditional *string
	if !isGlobalSearchCacheMode() {
		if q := providerManagedTagsSearchQuery(runID); q != "" {
			childAdditional = &q
		}
	}
	childList, err := listPolicyResources(connector, context, entry.childSearchType, childAdditional)
	if err != nil {
		return fmt.Errorf("error listing composite child resource %s for parent %s: %w", entry.childSearchType, resourceType, err)
	}
	if childAdditional != nil && len(resultList) > 0 && len(childList) == 0 {
		childList, err = listPolicyResources(connector, context, entry.childSearchType, nil)
		if err != nil {
			return fmt.Errorf("error listing composite child resource %s for parent %s: %w", entry.childSearchType, resourceType, err)
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

// groupRulesByValidParentPath buckets rules by trimmed ParentPath when that path is in validParentPaths.
// Rules with nil/empty ParentPath or unknown parent paths are skipped. Order within each bucket matches
// the order rules appear in the input slice (typically child search result order).
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

// attachRulesByParentPath groups rules onto parents where rule.ParentPath matches getPath(parent) (trimmed).
// Rules without ParentPath or with unknown parents are skipped. Parents are returned in input order.
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
					return &typedVal, cacheUsed, cacheAttempted, nil
				}
			}
		}
	}

	obj, err := backendRead()
	if err != nil {
		return nil, cacheUsed, cacheAttempted, err
	}

	// Handle tag-based cache mode: validate and patch provider-managed tags if missing
	if !isGlobalSearchCacheMode() {
		_, patchErr := ensureProviderManagedTagsWithPatchFunc(obj, m, patchFunc)
		if patchErr != nil {
			// Log the error but don't fail the read operation
			log.Printf("[WARNING] Failed to patch provider-managed tags for %s %s: %v", resourceType, resourceID, patchErr)
		}
	}

	return obj, cacheUsed, cacheAttempted, nil
}

// providerManagedTagsSearchFragments returns NSX search tag clauses for provider-managed
// tags (managed-by=terraform, tf-run-id). Empty when runID is unset.
func providerManagedTagsSearchFragments(runID string) []string {
	if runID == "" {
		return nil
	}
	return []string{
		fmt.Sprintf("tags.scope:%s", escapeSpecialCharacters("nsx-tf/managed-by")),
		fmt.Sprintf("tags.tag:%s", escapeSpecialCharacters("terraform")),
		fmt.Sprintf("tags.scope:%s", escapeSpecialCharacters("tf-run-id")),
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

// buildTagQuery extracts tags from resource data and builds NSX-T search query string
// In tag search mode, automatically appends provider-managed tags if not present
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

	// Tags are stored as *schema.Set
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

		// Extract scope and tag values
		scope, hasScope := tagMap["scope"]
		tag, hasTag := tagMap["tag"]
		// Build query parts for scope and tag
		if hasScope && scope != nil && scope.(string) != "" {
			scopeStr := escapeSpecialCharacters(scope.(string))
			if isManagedDefaultTagScope(&scopeStr) {
				managedTagPresent = true
			}
			tagQueries = append(tagQueries, fmt.Sprintf("tags.scope:%s", escapeSpecialCharacters(scope.(string))))
		}
		if hasTag && tag != nil && tag.(string) != "" {
			tagQueries = append(tagQueries, fmt.Sprintf("tags.tag:%s", escapeSpecialCharacters(tag.(string))))
		}
	}
	// In tag search mode (not global search cache mode), automatically append provider-managed tags if not present
	if !isGlobalSearchCacheMode() && !managedTagPresent && runID != "" {
		tagQueries = append(tagQueries, providerManagedTagsSearchFragments(runID)...)
	}

	if len(tagQueries) == 0 {
		return ""
	}

	// Join all tag queries with AND
	return strings.Join(tagQueries, " AND ")
}

// InvalidateCacheForResourceType clears all cached entries for the given resource type.
// Must be called after write operations (create/update) to prevent stale reads in the
// subsequent perpetual-diff plan, where isRefreshPhase returns true and cache is used.
// This is added cause the testcases runs create/update/data-source read in a single run.
func InvalidateCacheForResourceType(resourceType string) {
	if !IsCacheEnabled() {
		return
	}
	tc := gcache.getTypeCache(resourceType)
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.data = make(map[string]map[string]*data.StructValue)
	log.Printf("[DEBUG] Cache invalidated for resourceType=%s", resourceType)
}

// ensureProviderManagedTagsWithPatchFunc checks and patches tags using the provided patch function
func ensureProviderManagedTagsWithPatchFunc[T any](obj T, m interface{}, patchFunc func(obj T) error) (interface{}, error) {
	// Use reflection to check if the object has a Tags field
	objValue := reflect.ValueOf(obj)
	for objValue.IsValid() && objValue.Kind() == reflect.Ptr {
		if objValue.IsNil() {
			// No object to patch
			return nil, nil
		}
		objValue = objValue.Elem()
	}
	if !objValue.IsValid() || objValue.Kind() != reflect.Struct {
		return nil, nil
	}

	tagsField := objValue.FieldByName("Tags")
	if !tagsField.IsValid() {
		// Resource doesn't have tags, skip validation
		return nil, nil
	}

	// Extract current tags from the resource
	var currentTags []model.Tag
	if tagsField.Kind() == reflect.Slice && !tagsField.IsNil() {
		for i := 0; i < tagsField.Len(); i++ {
			tag := tagsField.Index(i).Interface().(model.Tag)
			currentTags = append(currentTags, tag)
		}
	}

	// Get expected provider-managed tags
	runID := m.(nsxtClients).CommonConfig.contextID
	expectedManagedTags := getProviderManagedDefaultTags(runID)

	// Check if provider-managed tags are missing
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

	// Merge current user tags with expected provider-managed tags
	userTags := make([]model.Tag, 0)
	for _, tag := range currentTags {
		if !isManagedDefaultTag(tag) {
			userTags = append(userTags, tag)
		}
	}

	// Create merged tags
	mergedTags := mergeManagedDefaultAndUserTags(expectedManagedTags, userTags)

	// Update the resource object's Tags field using reflection
	if tagsField.CanSet() {
		tagsField.Set(reflect.ValueOf(mergedTags))

		// Patch the updated object to NSX API using the provided patch function
		if patchFunc != nil {
			err := patchFunc(obj)
			if err != nil {
				return nil, fmt.Errorf("failed to patch tags to NSX API: %w", err)
			}
			log.Printf("[DEBUG] Patched provider-managed tags successfully")
		}

		return obj, nil
	}

	return nil, fmt.Errorf("Provider managed tags field is not settable")
}
