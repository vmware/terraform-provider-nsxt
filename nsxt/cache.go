package nsxt

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "github.com/vmware/terraform-provider-nsxt/api/utl"
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

func IsCacheEnabled() bool {
	return os.Getenv("NSXT_ENABLE_CACHE") == "true"
}

func isGlobalSearchCacheMode() bool {
	return os.Getenv("USE_GLOBAL_SEARCH_CACHE") == "true"
}

func isRefreshPhase(d *schema.ResourceData) bool {
	return d.Id() != "" && !d.HasChangesExcept()
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
		return val, nil
	}
	err := tc.writeCache(query, resourceType, d, m, connector)
	if err != nil {
		return nil, err
	}
	return tc.getQueryResult(query, resourceID)
}

func (c *resourceTypeCache) getListOfPolicyResources(query string, d *schema.ResourceData, connector client.Connector, context api.SessionContext, resourceType string, runID string) error {
	// Now we have access to the proper runID passed from writeCache
	additionalQuery := buildTagQuery(d, runID)
	resultList, err := listPolicyResources(connector, context, resourceType, &additionalQuery)
	if err != nil {
		return fmt.Errorf("error listing resource %s %w", resourceType, err)
	}
	tmp := converListToMap(resultList)
	//convert list to map
	c.data[query] = tmp
	return nil
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
		_, patchErr := ensureProviderManagedTagsWithPatchFunc(obj, d, m, patchFunc)
		if patchErr != nil {
			// Log the error but don't fail the read operation
			fmt.Printf("[WARN] Failed to patch provider-managed tags for %s %s: %v\n", resourceType, resourceID, patchErr)
		}
	}

	return obj, cacheUsed, cacheAttempted, nil
}

// buildTagQuery extracts tags from resource data and builds NSX-T search query string
// In tag search mode, automatically appends provider-managed tags if not present
func buildTagQuery(d *schema.ResourceData, runID string) string {
	managedTagPresent := false
	if d == nil {
		return ""
	}

	shouldAddProviderTags := !isGlobalSearchCacheMode() && runID != ""
	providerManagedTagQueries := func() []string {
		return []string{
			fmt.Sprintf("tags.scope:%s", escapeSpecialCharacters("nsx-tf/managed-by")),
			fmt.Sprintf("tags.tag:%s", escapeSpecialCharacters("terraform")),
			fmt.Sprintf("tags.scope:%s", escapeSpecialCharacters("tf-run-id")),
			fmt.Sprintf("tags.tag:%s", escapeSpecialCharacters(runID)),
		}
	}
	providerManagedTagQuery := func() string {
		return strings.Join(providerManagedTagQueries(), " AND ")
	}

	tags, exists := d.GetOk("tag")
	if !exists {
		if shouldAddProviderTags {
			return providerManagedTagQuery()
		}
		return ""
	}

	// Tags are stored as *schema.Set
	tagSet, ok := tags.(*schema.Set)
	if !ok {
		if shouldAddProviderTags {
			return providerManagedTagQuery()
		}
		return ""
	}
	if tagSet.Len() == 0 {
		if shouldAddProviderTags {
			return providerManagedTagQuery()
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
		// Append provider-managed tags: managed-by=terraform and tf-run-id=<runID>
		tagQueries = append(tagQueries, providerManagedTagQueries()...)
	}

	if len(tagQueries) == 0 {
		return ""
	}

	// Join all tag queries with AND
	return strings.Join(tagQueries, " AND ")
}

// ensureProviderManagedTagsWithPatchFunc checks and patches tags using the provided patch function
func ensureProviderManagedTagsWithPatchFunc[T any](obj T, d *schema.ResourceData, m interface{}, patchFunc func(obj T) error) (interface{}, error) {
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
			fmt.Println("pathc func success")
		}

		return obj, nil
	}

	return nil, fmt.Errorf("Tags field is not settable")
}
