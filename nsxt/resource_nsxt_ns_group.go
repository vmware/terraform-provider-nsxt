package nsxt

import(
    "github.com/hashicorp/terraform/helper/schema"
    api "github.com/vmware/go-vmware-nsxt"
    "github.com/vmware/go-vmware-nsxt/manager"
    "net/http"
    "fmt"
)

func resourceNsGroup() *schema.Resource {
    return &schema.Resource{
        Create: resourceNsGroupCreate,
        Read: resourceNsGroupRead,
        Update: resourceNsGroupUpdate,
        Delete: resourceNsGroupDelete,

        Schema: map[string]*schema.Schema{
            "revision": getRevisionSchema(),
            "system_owned": getSystemOwnedSchema(),
            "description": &schema.Schema{
                Type:        schema.TypeString,
                Description: "Description of this resource",
                Optional:    true,
            },
            "display_name": &schema.Schema{
                Type:        schema.TypeString,
                Description: "Defaults to ID if not set",
                Optional:    true,
            },
            "tags": getTagsSchema(),
            "member_count": &schema.Schema{
                Type:        schema.TypeInt,
                Description: "Count of the members added to this NSGroup",
                Computed:    true,
            },
            "members": &schema.Schema{
                Type:     schema.TypeList,
                Description: "Reference to the direct/static members of the NSGroup.",
                Optional: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "target_type": &schema.Schema{
                            Type:     schema.TypeString,
                            Required: true,
                        },
                        "value": &schema.Schema{
                            Type:     schema.TypeString,
                            Required: true,
                        },
                    },
                },
            },
            // TODO(asarfaty): currently only NsGroupTagExpression?
            // "membership_criteria": &schema.Schema{
            //     Type:     schema.TypeList,
            //     Description: "List of tag or ID expressions which define the membership criteria for this NSGroup.",
            //     Optional: true,
            //     Elem: &schema.Resource{
            //         Schema: map[string]*schema.Schema{
            //             "resource_type": &schema.Schema{
            //                 Type:     schema.TypeString,
            //                 Default:  "NSGroupTagExpression",
            //                 Optional: true,
            //             },
            //             "target_type": &schema.Schema{
            //                 Type:     schema.TypeString,
            //                 Required: true,
            //             },
            //             "scope": &schema.Schema{
            //                 Type:     schema.TypeString,
            //                 Required: true,
            //             },
            //             "tag": &schema.Schema{
            //                 Type:     schema.TypeString,
            //                 Required: true,
            //             },
            //             "scope_op": &schema.Schema{
            //                 Type:     schema.TypeString,
            //                 Default:  "EQUALS",
            //                 Optional: true,
            //             },
            //             "tag_op": &schema.Schema{
            //                 Type:     schema.TypeString,
            //                 Default:  "EQUALS",
            //                 Optional: true,
            //             },
            //         },
            //     },
            // },
        },
    }
}

func getMembershipCriteriaFromSchema(d *schema.ResourceData) []manager.NsGroupTagExpression {
    criterias := d.Get("membership_criteria").([]interface{})
    var expresionList []manager.NsGroupTagExpression
    for _, criteria := range criterias {
        data := criteria.(map[string]interface{})
        elem := manager.NsGroupTagExpression{
            ResourceType: data["resource_type"].(string),
            Scope: data["scope"].(string),
            ScopeOp: data["scope_op"].(string),
            Tag: data["tag"].(string),
            TagOp: data["tag_op"].(string),
            TargetType: data["target_type"].(string),
        }
        expresionList = append(expresionList, elem)
    }
    return expresionList
}

func setMembershipCriteriaInSchema(d *schema.ResourceData, membershipCriterias []manager.NsGroupTagExpression) {
    var expresionList []map[string]interface{}
    for _, criteria := range membershipCriterias {
        elem := make(map[string]interface{})
        elem["resource_type"] = criteria.ResourceType
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
    members := d.Get("members").([]interface{})
    var expresionList []manager.NsGroupSimpleExpression
    for _, member := range members {
        data := member.(map[string]interface{})
        elem := manager.NsGroupSimpleExpression{
            ResourceType: "NSGroupSimpleExpression",
            Op: "EQUALS",
            TargetProperty: "id",
            TargetType: data["target_type"].(string),
            Value: data["value"].(string),
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
    d.Set("members", expresionList)
}

func resourceNsGroupCreate(d *schema.ResourceData, m interface{}) error {

    nsxClient := m.(*api.APIClient)

    description := d.Get("description").(string)
    display_name := d.Get("display_name").(string)
    tags := getTagsFromSchema(d)
    members := getMembersFromSchema(d)
    //membership_criteria := getMembershipCriteriaFromSchema(d)
    ns_group := manager.NsGroup {
        Description: description,
        DisplayName: display_name,
        Tags: tags,
        Members: members,
        //MembershipCriteria: membership_criteria,
    }

    ns_group, resp, err := nsxClient.GroupingObjectsApi.CreateNSGroup(nsxClient.Context, ns_group)

    if err != nil {
        return fmt.Errorf("Error during NsGroup create: %v", err)
    }

    if resp.StatusCode != http.StatusCreated {
        fmt.Printf("Unexpected status returned")
        return nil
    }
    d.SetId(ns_group.Id)
    
    return resourceNsGroupRead(d, m)
}

func resourceNsGroupRead(d *schema.ResourceData, m interface{}) error {

    nsxClient := m.(*api.APIClient)

    id := d.Id()
    if id == "" {
        return fmt.Errorf("Error obtaining logical object id")
    }

    localVarOptionals := make(map[string]interface{})
    localVarOptionals["populateReferences"] = true
    ns_group, resp, err := nsxClient.GroupingObjectsApi.ReadNSGroup(nsxClient.Context, id, localVarOptionals)
    if resp.StatusCode == http.StatusNotFound {
        fmt.Printf("NsGroup not found")
        d.SetId("")
        return nil
    }
    if err != nil {
        return fmt.Errorf("Error during NsGroup read: %v", err)
    }

    d.Set("revision", ns_group.Revision)
    d.Set("system_owned", ns_group.SystemOwned)
    d.Set("description", ns_group.Description)
    d.Set("display_name", ns_group.DisplayName)
    setTagsInSchema(d, ns_group.Tags)
    d.Set("member_count", ns_group.MemberCount)
    setMembersInSchema(d, ns_group.Members)
    //setMembershipCriteriaInSchema(d, ns_group.MembershipCriteria)

    return nil
}

func resourceNsGroupUpdate(d *schema.ResourceData, m interface{}) error {

    nsxClient := m.(*api.APIClient)

    id := d.Id()
    if id == "" {
        return fmt.Errorf("Error obtaining logical object id")
    }

    revision := int64(d.Get("revision").(int))
    description := d.Get("description").(string)
    display_name := d.Get("display_name").(string)
    tags := getTagsFromSchema(d)
    members := getMembersFromSchema(d)
    //membership_criteria := getMembershipCriteriaFromSchema(d)
    ns_group := manager.NsGroup {
        Revision: revision,
        Description: description,
        DisplayName: display_name,
        Tags: tags,
        Members: members,
        //MembershipCriteria: membership_criteria,
    }

    ns_group, resp, err := nsxClient.GroupingObjectsApi.UpdateNSGroup(nsxClient.Context, id, ns_group)

    if err != nil || resp.StatusCode == http.StatusNotFound {
        return fmt.Errorf("Error during NsGroup update: %v", err)
    }

    return resourceNsGroupRead(d, m)
}

func resourceNsGroupDelete(d *schema.ResourceData, m interface{}) error {

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
        fmt.Printf("NsGroup not found")
        d.SetId("")
    }
return nil
}
