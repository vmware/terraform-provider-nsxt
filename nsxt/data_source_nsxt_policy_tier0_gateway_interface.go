package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	// t0interface "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s"
)

func dataSourceNsxtPolicyTier0GatewayInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyTier0GatewayInterfaceRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"t0_gateway_name": {
				Type:        schema.TypeString,
				Description: "The name of the Tier0 gateway where the interface is linked",
				Required:    true,
			},
			"path": getPathSchema(),
			"edge_cluster_path": {
				Type:        schema.TypeString,
				Description: "The path of the edge cluster connected to the Tier0 gateway linked to this interface",
				Optional:    true,
				Computed:    true,
			},
			"segment_path": {
				Type:        schema.TypeString,
				Description: "The path of the edge cluster connected to the Tier0 gateway linked to this interface",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyTier0GatewayInterfaceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	t0Gw := d.Get("t0_gateway_name").(string)
	obj, err := policyDataSourceResourceRead(d, connector, getSessionContext(d, m), "Tier0Interface", nil)
	if err != nil {
		return err
	}

	// t0interface.NewLocaleServicesClient()
	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToGolang(obj, model.Tier0InterfaceBindingType())
	if len(errors) > 0 {
		return errors[0]
	}
	gwInterface := dataValue.(model.Tier0Interface)
	isT0, gwID, _, _ := parseGatewayInterfacePolicyPath(*gwInterface.Path)
	if isT0 && t0Gw == gwID {
		err := d.Set("path", *gwInterface.Path)
		if err != nil {
			return fmt.Errorf("Error while setting interface path : %v", err)
		}
		err = d.Set("edge_cluster_path", *gwInterface.EdgePath)
		if err != nil {
			return fmt.Errorf("Error while setting edge path : %v", err)
		}
		err = d.Set("description", *gwInterface.Description)
		if err != nil {
			return fmt.Errorf("Error while setting the interface description : %v", err)
		}
		err = d.Set("segment_path", *gwInterface.SegmentPath)
		if err != nil {
			return fmt.Errorf("Error while setting the segment connected to the interface : %v", err)
		}
	}
	// Single edge cluster is not informative for global manager
	// if isPolicyGlobalManager(m) {
	// 	err = d.Set("edge_cluster_path", "")
	// 	if err != nil {
	// 		return fmt.Errorf("Error while setting edge path : %v", err)
	// 	}
	// }
	return nil

}
