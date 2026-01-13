package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicySegmentPort() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicySegmentPortRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"context":      getContextSchemaWithSpec(utl.SessionContextSpec{IsRequired: false, IsComputed: false, IsVpc: false, AllowDefaultProject: false, FromGlobal: true}),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"vif_id": {
				Type:        schema.TypeString,
				Description: "Segment Port attachment id",
				Optional:    true,
			},
			"segment_path": {
				Type:        schema.TypeString,
				Description: "Segment path",
				Optional:    true,
			},
		},
	}
}

func dataSourceNsxtPolicySegmentPortRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()
	var segmentPortList []*data.StructValue
	var err error
	if vifId, exists := d.GetOkExists("vif_id"); exists {
		addQuery := fmt.Sprintf("attachment.id:%s", escapeSpecialCharacters(vifId.(string)))
		segmentPortList, err = listPolicyResources(connector, getSessionContext(d, m), "SegmentPort", &addQuery)
		if err != nil {
			return fmt.Errorf("Error getting the Segment Port list : %v", err)
		}
		if len(segmentPortList) == 0 {
			return fmt.Errorf("Segment Port with VIF attachment id %s not found", vifId)
		}
	} else if portName, exists := d.GetOkExists("id"); exists {
		addQuery := fmt.Sprintf("id:%s", escapeSpecialCharacters(portName.(string)))
		segmentPortList, err = listPolicyResources(connector, getSessionContext(d, m), "SegmentPort", &addQuery)
		if err != nil {
			return fmt.Errorf("Error getting the Segment Port list : %v", err)
		}
		if len(segmentPortList) == 0 {
			return fmt.Errorf("Segment Port with id %s not found", portName)
		}
	} else if portName, exists := d.GetOkExists("display_name"); exists {
		portNameStr := portName.(string)
		addQuery := fmt.Sprintf("display_name:%s*", escapeSpecialCharacters(portNameStr))
		segmentPortList, err = listPolicyResources(connector, getSessionContext(d, m), "SegmentPort", &addQuery)
		if err != nil {
			return fmt.Errorf("Error getting the Segment Port list : %v", err)
		}
		if len(segmentPortList) == 0 {
			return fmt.Errorf("Segment Port with display_name %s not found", portName)
		}
		// Filter results to prefer exact matches, then prefix matches
		var perfectMatch, prefixMatch []*data.StructValue
		for _, result := range segmentPortList {
			dataValue, errors := converter.ConvertToGolang(result, model.SegmentPortBindingType())
			if len(errors) > 0 {
				continue
			}
			segmentPort := dataValue.(model.SegmentPort)
			if segmentPort.DisplayName != nil {
				if *segmentPort.DisplayName == portNameStr {
					perfectMatch = append(perfectMatch, result)
				} else if strings.HasPrefix(*segmentPort.DisplayName, portNameStr) {
					prefixMatch = append(prefixMatch, result)
				}
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple Segment Ports with display_name '%s'", portNameStr)
			}
			segmentPortList = perfectMatch
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple Segment Ports with display_name starting with '%s'", portNameStr)
			}
			segmentPortList = prefixMatch
		} else {
			return fmt.Errorf("Segment Port with display_name %s not found", portName)
		}
	} else {
		return fmt.Errorf("Atleast one of vif_id, display_name or id should be set")
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
	if segmentPort.Attachment != nil && segmentPort.Attachment.Id != nil {
		d.Set("vif_id", segmentPort.Attachment.Id)
	}
	d.Set("segment_path", getSegmentPathFromPortPath(*segmentPort.Path))

	return nil
}

func getSegmentPathFromPortPath(portPath string) string {
	parts := strings.Split(portPath, "/")
	parts = parts[:len(parts)-2]

	return strings.Join(parts, "/")
}
