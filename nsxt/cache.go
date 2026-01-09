package nsxt

import (
	"fmt"
	"os"
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

// SetCacheEnabled allows enabling/disabling cache for benchmarking
var cacheEnabled = os.Getenv("NSXT_ENABLE_CACHE") == "true"

func IsCacheEnabled() bool {
	return cacheEnabled
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
		} else if resource.DisplayName != nil {
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
	err := c.getListOfPolicyResources(query, d, connector, getSessionContext(d, m), resourceType)
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
	additionalQuery := buildTagQuery(d)
	if additionalQuery == "" {
		return query
	}
	return fmt.Sprintf("%s AND %s", query, additionalQuery)
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

func (c *resourceTypeCache) getListOfPolicyResources(query string, d *schema.ResourceData, connector client.Connector, context api.SessionContext, resourceType string) error {
	additionalQuery := buildTagQuery(d)
	resultList, err := listPolicyResources(connector, context, resourceType, &additionalQuery)
	if err != nil {
		return fmt.Errorf("error listing resource %s %w", resourceType, err)
	}
	tmp := converListToMap(resultList)
	//convert list to map
	c.data[query] = tmp
	return nil
}

func CacheAwareResourceRead[T any](d *schema.ResourceData, m interface{}, connector client.Connector, resourceID string, resourceType string, bindingType bindings.BindingType, backendRead func() (T, error)) (T, bool, bool, error) {
	var zero T
	cacheUsed := false
	cacheAttempted := false

	if isRefreshPhase(d) && IsCacheEnabled() {
		cacheAttempted = true
		val, err := gcache.readCache(resourceID, resourceType, d, m, connector)
		if err == nil {
			converter := bindings.NewTypeConverter()
			goVal, convErrs := converter.ConvertToGolang(val.(*data.StructValue), bindingType)
			if len(convErrs) == 0 {
				obj, ok := goVal.(T)
				if ok {
					cacheUsed = true
					return obj, cacheUsed, cacheAttempted, nil
				}
			}
		}
	}

	obj, err := backendRead()
	if err != nil {
		return zero, cacheUsed, cacheAttempted, err
	}
	return obj, cacheUsed, cacheAttempted, nil
}

// buildTagQuery extracts tags from resource data and builds NSX-T search query string
func buildTagQuery(d *schema.ResourceData) string {
	if d == nil {
		return ""
	}
	tags, exists := d.GetOk("tag")
	if !exists {
		return ""
	}

	// Tags are stored as *schema.Set
	tagSet, ok := tags.(*schema.Set)
	if !ok {
		return ""
	}
	if tagSet.Len() == 0 {
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
			tagQueries = append(tagQueries, fmt.Sprintf("tags.scope:%s", escapeSpecialCharacters(scope.(string))))
		}
		if hasTag && tag != nil && tag.(string) != "" {
			tagQueries = append(tagQueries, fmt.Sprintf("tags.tag:%s", escapeSpecialCharacters(tag.(string))))
		}
	}

	if len(tagQueries) == 0 {
		return ""
	}

	// Join all tag queries with AND
	return strings.Join(tagQueries, " AND ")
}
