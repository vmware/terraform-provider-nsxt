package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"log"
	"net/http"
)

var nsGroupTargetTypeValues = []string{"NSGroup", "IPSet", "LogicalPort", "LogicalSwitch", "MACSet"}
var nsGroupMembershipCriteriaTargetTypeValues = []string{"LogicalPort", "LogicalSwitch", "VirtualMachine"}
var nsGroupTagOperationValues = []string{"EQUALS"}

func resourceNsxtNsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtNsGroupCreate,
		Read:   resourceNsxtNsGroupRead,
		Update: resourceNsxtNsGroupUpdate,
		Delete: resourceNsxtNsGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"tag": getTagsSchema(),
			"member": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "Reference to the direct/static members of the NSGroup.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_type": &schema.Schema{
							Type:         schema.TypeString,
							Description:  "Type of the resource on which this expression is evaluated",
							Required:     true,
							ValidateFunc: validation.StringInSlice(nsGroupTargetTypeValues, false),
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Description: "Value that satisfies this expression",
							Required:    true,
						},
					},
				},
			},
			"membership_criteria": &schema.Schema{
				Type:        schema.TypeList,
				Description: "List of tag expressions which define the membership criteria for this NSGroup.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_type": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(nsGroupMembershipCriteriaTargetTypeValues, false),
						},
						"scope": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"tag": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"scope_op": &schema.Schema{
							Type:         schema.TypeString,
							Default:      "EQUALS",
							Optional:     true,
							ValidateFunc: validation.StringInSlice(nsGroupTagOperationValues, false),
						},
						"tag_op": &schema.Schema{
							Type:         schema.TypeString,
							Default:      "EQUALS",
							Optional:     true,
							ValidateFunc: validation.StringInSlice(nsGroupTagOperationValues, false),
						},
					},
				},
			},
		},
	}
}

func getMembershipCriteriaFromSchema(d *schema.ResourceData) []manager.NsGroupTagExpression {
	criterias := d.Get("membership_criteria").([]interface{})
	var expresionList []manager.NsGroupTagExpression
	for _, criteria := range criterias {
		data := criteria.(map[string]interface{})
		elem := manager.NsGroupTagExpression{
			ResourceType: "NSGroupTagExpression",
			Scope:        data["scope"].(string),
			ScopeOp:      data["scope_op"].(string),
			Tag:          data["tag"].(string),
			TagOp:        data["tag_op"].(string),
			TargetType:   data["target_type"].(string),
		}
		expresionList = append(expresionList, elem)
	}
	return expresionList
}

func setMembershipCriteriaInSchema(d *schema.ResourceData, membershipCriterias []manager.NsGroupTagExpression) {
	var expresionList []map[string]interface{}
	for _, criteria := range membershipCriterias {
		elem := make(map[string]interface{})
		elem["scope"] = criteria.Scope
		elem["scope_op"] = criteria.ScopeOp
		elem["tag"] = criteria.Tag
		elem["tag_op"] = criteria.TagOp
		elem["target_type"] = criteria.TargetType
		expresionList = append(expresionList, elem)
	}
	d.Set("membership_criteria", expresionList)
}

func getMembersFromSchema(d *schema.ResourceData) []manager.NsGroupSimpleExpression {
	members := d.Get("member").(*schema.Set).List()
	var expresionList []manager.NsGroupSimpleExpression
	for _, member := range members {
		data := member.(map[string]interface{})
		elem := manager.NsGroupSimpleExpression{
			ResourceType:   "NSGroupSimpleExpression",
			Op:             "EQUALS",
			TargetProperty: "id",
			TargetType:     data["target_type"].(string),
			Value:          data["value"].(string),
		}
		expresionList = append(expresionList, elem)
	}
	return expresionList
}

func setMembersInSchema(d *schema.ResourceData, members []manager.NsGroupSimpleExpression) {
	var expresionList []map[string]interface{}
	for _, member := range members {
		elem := make(map[string]interface{})
		elem["target_type"] = member.TargetType
		elem["value"] = member.Value
		expresionList = append(expresionList, elem)
	}
	d.Set("member", expresionList)
}

func resourceNsxtNsGroupCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	members := getMembersFromSchema(d)
	membershipCriteria := getMembershipCriteriaFromSchema(d)
	nsGroup := manager.NsGroup{
		Description:        description,
		DisplayName:        displayName,
		Tags:               tags,
		Members:            members,
		MembershipCriteria: membershipCriteria,
	}

	nsGroup, resp, err := nsxClient.GroupingObjectsApi.CreateNSGroup(nsxClient.Context, nsGroup)

	if err != nil {
		return fmt.Errorf("Error during NsGroup create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during NsGroup create: %v", resp.StatusCode)
	}
	d.SetId(nsGroup.Id)

	return resourceNsxtNsGroupRead(d, m)
}

func resourceNsxtNsGroupRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	localVarOptionals := make(map[string]interface{})
	localVarOptionals["populateReferences"] = true
	nsGroup, resp, err := nsxClient.GroupingObjectsApi.ReadNSGroup(nsxClient.Context, id, localVarOptionals)
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] NsGroup %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during NsGroup read: %v", err)
	}

	d.Set("revision", nsGroup.Revision)
	d.Set("description", nsGroup.Description)
	d.Set("display_name", nsGroup.DisplayName)
	setTagsInSchema(d, nsGroup.Tags)
	setMembersInSchema(d, nsGroup.Members)
	setMembershipCriteriaInSchema(d, nsGroup.MembershipCriteria)

	return nil
}

func resourceNsxtNsGroupUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	members := getMembersFromSchema(d)
	membershipCriteria := getMembershipCriteriaFromSchema(d)
	nsGroup := manager.NsGroup{
		Revision:           revision,
		Description:        description,
		DisplayName:        displayName,
		Tags:               tags,
		Members:            members,
		MembershipCriteria: membershipCriteria,
	}

	nsGroup, resp, err := nsxClient.GroupingObjectsApi.UpdateNSGroup(nsxClient.Context, id, nsGroup)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during NsGroup update: %v", err)
	}

	return resourceNsxtNsGroupRead(d, m)
}

func resourceNsxtNsGroupDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.GroupingObjectsApi.DeleteNSGroup(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during NsGroup delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] NsGroup %s not found", id)
		d.SetId("")
	}
	return nil
}
