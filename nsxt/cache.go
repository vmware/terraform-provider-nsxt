package nsxt

import (
	"fmt"
	"sync"

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
		if resource.DisplayName != nil {
			ret[*resource.DisplayName] = obj
		}
	}
	fmt.Println("conver to map ", ret)
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

func (c *cache) getQueryResult(query string, displayName string) *data.StructValue {
	if val, ok := c.data[query]; ok {
		return val[displayName]
	}
	return nil
}

func (c *cache) writeCache(query string, resourceType string, d *schema.ResourceData, m interface{}, connector client.Connector) error {
	c.cacheMis += 1
	fmt.Println("cacheMis ", c.cacheMis)
	err := c.getListOfPolicyResources(query, d, connector, getSessionContext(d, m), resourceType)
	if err != nil {
		fmt.Println("failed to read the cache", err)
		return err
	}
	return nil
}

func (c *cache) readCache(displayName string, resourceType string, d *schema.ResourceData, m interface{}, connector client.Connector) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	query := getQueryString(resourceType, getSessionContext(d, m))
	if val := c.getQueryResult(query, displayName); val != nil {
		c.cacheHit += 1
		fmt.Println("cacheHit ", c.cacheHit)
		return val, nil
	}
	err := c.writeCache(query, resourceType, d, m, connector)
	if err != nil {
		return nil, err
	}
	return c.getQueryResult(query, displayName), nil
}

func (c *cache) getListOfPolicyResources(query string, d *schema.ResourceData, connector client.Connector, context api.SessionContext, resourceType string) error {
	resultList, err := listPolicyResources(connector, context, resourceType, nil)
	if err != nil {
		return fmt.Errorf("error listing resource %s %w", resourceType, err)
	}
	fmt.Println("listPolicyResources ", resultList)
	tmp := converListToMap(resultList)
	//convert list to map
	c.data[query] = tmp
	return nil
}

func isRefreshPhase(d *schema.ResourceData) bool {
	return d.Id() != "" && !d.HasChangesExcept()
}
