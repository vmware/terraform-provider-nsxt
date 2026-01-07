package nsxt

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "github.com/vmware/terraform-provider-nsxt/api/utl"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type cache struct {
	mu       sync.RWMutex
	data     map[string]map[string]*data.StructValue //key: query , key DisplayName and objects
	cacheHit int
	cacheMis int
}

var gcache = &cache{data: make(map[string]map[string]*data.StructValue), cacheHit: 0, cacheMis: 0}

func converListToMap(list []*data.StructValue) map[string]*data.StructValue {
	converter := bindings.NewTypeConverter()
	ret := make(map[string]*data.StructValue)
	for _, obj := range list {
		dataValue, errors := converter.ConvertToGolang(obj, model.PolicyConfigResourceBindingType())
		if len(errors) > 0 {
			return nil
		}
		resource := dataValue.(model.PolicyConfigResource)
		if resource.DisplayName != nil { //TBD: DisplayName will be changed to ID
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

func (c *cache) getQueryResult(query string, displayName string) (*data.StructValue, error) {
	if val, ok := c.data[query]; ok {
		return val[displayName], nil
	}
	return nil, fmt.Errorf("element is not found")
}

func trackTime(start time.Time, name string) {
	_ = time.Since(start).Seconds()
	_ = name
}

func (c *cache) writeCache(query string, resourceType string, d *schema.ResourceData, m interface{}, connector client.Connector) error {
	start := time.Now()
	defer trackTime(start, fmt.Sprint("One Time taken to populate the cache "))
	c.cacheMis += 1
	err := c.getListOfPolicyResources(query, d, connector, getSessionContext(d, m), resourceType)
	if err != nil {
		return err
	}
	return nil
}

func (c *cache) readCache(displayName string, resourceType string, d *schema.ResourceData, m interface{}, connector client.Connector) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	query := getQueryString(resourceType, getSessionContext(d, m))
	if val, _ := c.getQueryResult(query, displayName); val != nil {
		c.cacheHit += 1
		return val, nil
	}
	err := c.writeCache(query, resourceType, d, m, connector)
	if err != nil {
		return nil, err
	}
	return c.getQueryResult(query, displayName)
}

func (c *cache) getListOfPolicyResources(query string, d *schema.ResourceData, connector client.Connector, context api.SessionContext, resourceType string) error {
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

func isRefreshPhase(d *schema.ResourceData) bool {
	return d.Id() != "" && !d.HasChangesExcept()
}

// SetCacheEnabled allows enabling/disabling cache for benchmarking
var cacheEnabled = os.Getenv("NSXT_ENABLE_CACHE") == "true"

func IsCacheEnabled() bool {
	return cacheEnabled
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
