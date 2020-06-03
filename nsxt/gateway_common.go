package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt"
	global_policy "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"strings"
)

func getFailoverModeSchema(defaultValue string) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "Failover mode",
		Default:      defaultValue,
		Optional:     true,
		ValidateFunc: validation.StringInSlice(policyFailOverModeValues, false),
	}
}

func getIPv6NDRAPathSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "The path of an IPv6 NDRA profile",
		Optional:    true,
		Computed:    true,
	}
}

func getIPv6DadPathSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "The path of an IPv6 DAD profile",
		Optional:    true,
		Computed:    true,
	}
}

func getPolicyEdgeClusterPathSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "The path of the edge cluster connected to this gateway",
		Optional:     true,
		ValidateFunc: validatePolicyPath(),
	}
}

func getPolicyLocaleServiceSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		Description:   "Locale Service for the gateway",
		ConflictsWith: []string{"edge_cluster_path"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"edge_cluster_path": {
					Type:         schema.TypeString,
					Description:  "The path of the edge cluster connected to this gateway",
					Required:     true,
					ValidateFunc: validatePolicyPath(),
				},
				"preferred_edge_paths": {
					Type:        schema.TypeSet,
					Description: "Paths of specific edge nodes",
					Optional:    true,
					Elem:        getElemPolicyPathSchema(),
				},
				"path": getPathSchema(),
			},
		},
	}
}

func getIpv6ProfilePathsFromSchema(d *schema.ResourceData) []string {
	var profiles []string
	if d.Get("ipv6_ndra_profile_path") != "" {
		profiles = append(profiles, d.Get("ipv6_ndra_profile_path").(string))
	}
	if d.Get("ipv6_dad_profile_path") != "" {
		profiles = append(profiles, d.Get("ipv6_dad_profile_path").(string))
	}
	return profiles
}

func setIpv6ProfilePathsInSchema(d *schema.ResourceData, paths []string) error {
	for _, path := range paths {
		if strings.HasPrefix(path, "/infra/ipv6-ndra-profiles") {
			d.Set("ipv6_ndra_profile_path", path)
		}
		if strings.HasPrefix(path, "/infra/ipv6-dad-profiles") {
			d.Set("ipv6_dad_profile_path", path)
		}
	}
	return nil
}

func getGatewayInterfaceSubnetsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of IP addresses and network prefixes for this interface",
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateIPCidr(),
		},
		Required: true,
	}
}

func getMtuSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Description:  "Maximum transmission unit specifies the size of the largest packet that a network protocol can transmit",
		ValidateFunc: validation.IntAtLeast(64),
	}
}

var gatewayInterfaceUrpfModeValues = []string{
	model.Tier0Interface_URPF_MODE_NONE,
	model.Tier0Interface_URPF_MODE_STRICT,
}

func getGatewayInterfaceUrpfModeSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "Unicast Reverse Path Forwarding mode",
		ValidateFunc: validation.StringInSlice(gatewayInterfaceUrpfModeValues, false),
		Default:      model.Tier0Interface_URPF_MODE_STRICT,
	}
}

func listPolicyGatewayLocaleServices(connector *client.RestConnector, gwID string, listLocaleServicesFunc func(*client.RestConnector, string, *string) (model.LocaleServicesListResult, error)) ([]model.LocaleServices, error) {
	var results []model.LocaleServices
	var cursor *string
	var count int64
	total := int64(0)

	for {
		listResponse, err := listLocaleServicesFunc(connector, gwID, cursor)
		if err != nil {
			return results, err
		}
		cursor = listResponse.Cursor
		count = *listResponse.ResultCount
		results = append(results, listResponse.Results...)
		if total == 0 {
			// first response
			total = count
		}
		if int64(len(results)) >= total {
			return results, nil
		}
	}
}

func initChildLocaleService(serviceStruct *model.LocaleServices, markForDelete bool) (*data.StructValue, error) {
	childService := model.ChildLocaleServices{
		ResourceType:    "ChildLocaleServices",
		LocaleServices:  serviceStruct,
		MarkedForDelete: &markForDelete,
	}

	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)
	dataValue, err := converter.ConvertToVapi(childService, model.ChildLocaleServicesBindingType())

	if err != nil {
		return nil, err[0]
	}

	return dataValue.(*data.StructValue), nil
}

func initGatewayLocaleServices(d *schema.ResourceData, connector *client.RestConnector, listLocaleServicesFunc func(*client.RestConnector, string, bool) ([]model.LocaleServices, error)) ([]*data.StructValue, error) {
	var localeServices []*data.StructValue

	services := d.Get("locale_service").([]interface{})
	if len(services) == 0 {
		return localeServices, nil
	}

	existingServices := make(map[string]bool)
	if len(d.Id()) > 0 {
		// This is an update - we might need to delete locale services
		existingServiceObjects, errList := listLocaleServicesFunc(connector, d.Id(), true)
		if errList != nil {
			return nil, errList
		}

		for _, obj := range existingServiceObjects {
			existingServices[*obj.Id] = true
		}
	}
	lsType := "LocaleServices"
	for _, service := range services {
		cfg := service.(map[string]interface{})
		edgeClusterPath := cfg["edge_cluster_path"].(string)
		edgeNodes := interface2StringList(cfg["preferred_edge_paths"].(*schema.Set).List())
		path := cfg["path"].(string)

		var serviceID string
		if path != "" {
			serviceID = getPolicyIDFromPath(path)
		} else {
			serviceID = newUUID()
		}
		serviceStruct := model.LocaleServices{
			Id:                 &serviceID,
			ResourceType:       &lsType,
			EdgeClusterPath:    &edgeClusterPath,
			PreferredEdgePaths: edgeNodes,
		}

		dataValue, err := initChildLocaleService(&serviceStruct, false)
		if err != nil {
			return localeServices, err
		}

		localeServices = append(localeServices, dataValue)
		existingServices[serviceID] = false
	}
	// Add instruction to delete services that are no longer present in intent
	for id, shouldDelete := range existingServices {
		if shouldDelete {
			serviceStruct := model.LocaleServices{
				Id:           &id,
				ResourceType: &lsType,
			}
			dataValue, err := initChildLocaleService(&serviceStruct, true)
			if err != nil {
				return localeServices, err
			}
			localeServices = append(localeServices, dataValue)
		}
	}

	return localeServices, nil
}

func policyInfraPatch(obj model.Infra, isGlobalManager bool, connector *client.RestConnector, enforceRevision bool) error {
	if isGlobalManager {
		infraClient := global_policy.NewDefaultGlobalInfraClient(connector)
		gmObj, err := convertModelBindingType(obj, model.InfraBindingType(), gm_model.InfraBindingType())
		if err != nil {
			return err
		}

		return infraClient.Patch(gmObj.(gm_model.Infra), &enforceRevision)
	}

	infraClient := nsx_policy.NewDefaultInfraClient(connector)
	return infraClient.Patch(obj, &enforceRevision)
}
