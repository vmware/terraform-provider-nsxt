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

func dataSourceNsxtPolicySegmentPorts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicySegmentPortsRead,

		Schema: map[string]*schema.Schema{
			"context": getContextSchemaWithSpec(utl.SessionContextSpec{IsRequired: false, IsComputed: false, IsVpc: false, AllowDefaultProject: false, FromGlobal: true}),
			"display_name": {
				Type:        schema.TypeString,
				Description: "Display name to filter segment ports",
				Optional:    true,
			},
			"vif_id": {
				Type:        schema.TypeString,
				Description: "Segment Port attachment id to filter segment ports",
				Optional:    true,
			},
			"segment_path": {
				Type:        schema.TypeString,
				Description: "Segment path to filter segment ports",
				Optional:    true,
			},
			"items": {
				Type:        schema.TypeList,
				Description: "List of segment ports matching the criteria",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Unique identifier of the segment port",
							Computed:    true,
						},
						"display_name": {
							Type:        schema.TypeString,
							Description: "Display name of the segment port",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of the segment port",
							Computed:    true,
						},
						"path": {
							Type:        schema.TypeString,
							Description: "Policy path of the segment port",
							Computed:    true,
						},
						"vif_id": {
							Type:        schema.TypeString,
							Description: "VIF attachment id of the segment port",
							Computed:    true,
						},
						"segment_path": {
							Type:        schema.TypeString,
							Description: "Parent segment path",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNsxtPolicySegmentPortsRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()
	var segmentPortList []*data.StructValue
	var err error
	var query string

	vifID, vifIDExists := d.GetOkExists("vif_id")
	displayName, displayNameExists := d.GetOkExists("display_name")
	segmentPath, segmentPathExists := d.GetOkExists("segment_path")

	if vifIDExists {
		query = fmt.Sprintf("attachment.id:%s", escapeSpecialCharacters(vifID.(string)))

		if segmentPathExists {
			query = fmt.Sprintf("%s AND parent_path:%s", query, escapeSpecialCharacters(segmentPath.(string)))
		}

		segmentPortList, err = listPolicyResources(connector, getSessionContext(d, m), "SegmentPort", &query)
		if err != nil {
			return fmt.Errorf("Error getting the Segment Port list: %v", err)
		}

		if displayNameExists {
			segmentPortList = filterSegmentPortsByDisplayName(segmentPortList, displayName.(string), converter)
		}
	} else if segmentPathExists {
		query = fmt.Sprintf("parent_path:%s", escapeSpecialCharacters(segmentPath.(string)))

		if displayNameExists {
			query = fmt.Sprintf("%s AND display_name:%s", query, escapeSpecialCharacters(displayName.(string)))
		}

		segmentPortList, err = listPolicyResources(connector, getSessionContext(d, m), "SegmentPort", &query)
		if err != nil {
			return fmt.Errorf("Error getting the Segment Port list: %v", err)
		}
	} else if displayNameExists {
		query = fmt.Sprintf("display_name:%s", escapeSpecialCharacters(displayName.(string)))
		segmentPortList, err = listPolicyResources(connector, getSessionContext(d, m), "SegmentPort", &query)
		if err != nil {
			return fmt.Errorf("Error getting the Segment Port list: %v", err)
		}
	} else {
		return fmt.Errorf("At least one of vif_id, segment_path, or display_name should be set")
	}

	items := make([]map[string]interface{}, 0, len(segmentPortList))
	for _, portData := range segmentPortList {
		dataValue, errors := converter.ConvertToGolang(portData, model.SegmentPortBindingType())
		if len(errors) > 0 {
			return fmt.Errorf("Error converting segment port: %v", errors[0])
		}

		segmentPort := dataValue.(model.SegmentPort)
		item := make(map[string]interface{})

		if segmentPort.Id != nil {
			item["id"] = *segmentPort.Id
		}
		if segmentPort.DisplayName != nil {
			item["display_name"] = *segmentPort.DisplayName
		}
		if segmentPort.Description != nil {
			item["description"] = *segmentPort.Description
		}
		if segmentPort.Path != nil {
			item["path"] = *segmentPort.Path
			item["segment_path"] = getSegmentPathFromPortPath(*segmentPort.Path)
		}
		if segmentPort.Attachment != nil && segmentPort.Attachment.Id != nil {
			item["vif_id"] = *segmentPort.Attachment.Id
		}

		items = append(items, item)
	}

	id := generateSegmentPortsDataSourceID(vifID, segmentPath, displayName, vifIDExists, segmentPathExists, displayNameExists)
	d.SetId(id)
	d.Set("items", items)

	return nil
}

func filterSegmentPortsByDisplayName(segmentPortList []*data.StructValue, displayName string, converter *bindings.TypeConverter) []*data.StructValue {
	filtered := make([]*data.StructValue, 0)

	for _, portData := range segmentPortList {
		dataValue, errors := converter.ConvertToGolang(portData, model.SegmentPortBindingType())
		if len(errors) > 0 {
			continue
		}

		segmentPort := dataValue.(model.SegmentPort)
		if segmentPort.DisplayName != nil && *segmentPort.DisplayName == displayName {
			filtered = append(filtered, portData)
		}
	}

	return filtered
}

func generateSegmentPortsDataSourceID(vifID, segmentPath, displayName interface{}, vifIDExists, segmentPathExists, displayNameExists bool) string {
	parts := make([]string, 0, 3)

	if vifIDExists {
		parts = append(parts, fmt.Sprintf("vif_id=%s", vifID))
	}
	if segmentPathExists {
		parts = append(parts, fmt.Sprintf("segment_path=%s", segmentPath))
	}
	if displayNameExists {
		parts = append(parts, fmt.Sprintf("display_name=%s", displayName))
	}

	return strings.Join(parts, ";")
}
