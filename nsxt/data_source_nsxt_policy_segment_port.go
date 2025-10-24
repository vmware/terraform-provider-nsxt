package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicySegmentPort() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicySegmentPortRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchema(false, false, false),
			"vif_id": {
				Type:        schema.TypeString,
				Description: "Segment Port attachment id",
				Required:    true,
			},
		},
	}
}

func dataSourceNsxtPolicySegmentPortRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()
	vifId := d.Get("vif_id").(string)
	addQuery := fmt.Sprintf("attachment.id:%s", escapeSpecialCharacters(vifId))
	segmentPortList, err := listPolicyResources(connector, getSessionContext(d, m), "SegmentPort", &addQuery)
	if err != nil {
		return fmt.Errorf("Error getting the Segment Port list : %v", err)
	}
	if len(segmentPortList) == 0 {
		return fmt.Errorf("Segment Port with VIF attachment id %s not found", vifId)
	}
	dataValue, errors := converter.ConvertToGolang(segmentPortList[0], model.SegmentPortBindingType())
	if len(errors) > 0 {
		return errors[0]
	}
	segmentPort := dataValue.(model.SegmentPort)
	d.SetId(*segmentPort.Id)
	d.Set("display_name", segmentPort.DisplayName)
	d.Set("description", segmentPort.Description)
	d.Set("path", segmentPort.Path)

	return nil
}
