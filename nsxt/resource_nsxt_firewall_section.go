package nsxt

import(
    "github.com/hashicorp/terraform/helper/schema"
    api "github.com/vmware/go-vmware-nsxt"
    "github.com/vmware/go-vmware-nsxt/manager"
    "net/http"
    "fmt"
)

func resourceFirewallSection() *schema.Resource {
    return &schema.Resource{
        Create: resourceFirewallSectionCreate,
        Read: resourceFirewallSectionRead,
        Update: resourceFirewallSectionUpdate,
        Delete: resourceFirewallSectionDelete,

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
            "is_default": &schema.Schema{
                Type:        schema.TypeBool,
                Description: "It is a boolean flag which reflects whether a firewall section is default section or not. Each Layer 3 and Layer 2 section will have at least and at most one default section",
                Computed:    true,
            },
            "rule_count": &schema.Schema{
                Type:        schema.TypeInt,
                Description: "Number of rules in this section",
                Computed:    true,
            },
            "section_type": &schema.Schema{
                Type:        schema.TypeString,
                Description: "Type of the rules which a section can contain. Only homogeneous sections are supported",
                Required:    true,
            },
            "stateful": &schema.Schema{
                Type:        schema.TypeBool,
                Description: "Stateful or Stateless nature of firewall section is enforced on all rules inside the section. Layer3 sections can be stateful or stateless. Layer2 sections can only be stateless",
                Required:    true,
                ForceNew:    true,
            },
            "applied_tos": getResourceReferencesSchema(false, false),
            // TODO(asarfaty): add rules
        },
    }
}

func resourceFirewallSectionCreate(d *schema.ResourceData, m interface{}) error {

    nsxClient := m.(*api.APIClient)

    description := d.Get("description").(string)
    display_name := d.Get("display_name").(string)
    tags := getTagsFromSchema(d)
    applied_tos := getResourceReferencesFromSchema(d, "applied_tos")
    is_default := d.Get("is_default").(bool)
    rule_count := int64(d.Get("rule_count").(int))
    section_type := d.Get("section_type").(string)
    stateful := d.Get("stateful").(bool)
    firewall_section := manager.FirewallSection {
        Description: description,
        DisplayName: display_name,
        Tags: tags,
        AppliedTos: applied_tos,
        IsDefault: is_default,
        RuleCount: rule_count,
        SectionType: section_type,
        Stateful: stateful,
    }

    localVarOptionals := make(map[string]interface{})
    firewall_section, resp, err := nsxClient.ServicesApi.AddSection(nsxClient.Context, firewall_section, localVarOptionals)

    if err != nil {
        return fmt.Errorf("Error during FirewallSection create: %v", err)
    }

    if resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("Unexpected status returned during FirewallSection create: %v", resp.StatusCode)
    }
    d.SetId(firewall_section.Id)

    return resourceFirewallSectionRead(d, m)
}

func resourceFirewallSectionRead(d *schema.ResourceData, m interface{}) error {

    nsxClient := m.(*api.APIClient)

    id := d.Id()
    if id == "" {
        return fmt.Errorf("Error obtaining logical object id")
    }

    firewall_section, resp, err := nsxClient.ServicesApi.GetSectionWithRulesListWithRules(nsxClient.Context, id)
    if resp.StatusCode == http.StatusNotFound {
        fmt.Printf("FirewallSection %s not found", id)
        d.SetId("")
        return nil
    }
    if err != nil {
        return fmt.Errorf("Error during FirewallSection %s read: %v", id, err)
    }

    d.Set("revision", firewall_section.Revision)
    d.Set("system_owned", firewall_section.SystemOwned)
    d.Set("description", firewall_section.Description)
    d.Set("display_name", firewall_section.DisplayName)
    setTagsInSchema(d, firewall_section.Tags)
    setResourceReferencesInSchema(d, firewall_section.AppliedTos, "applied_tos")
    d.Set("is_default", firewall_section.IsDefault)
    d.Set("rule_count", firewall_section.RuleCount)
    d.Set("section_type", firewall_section.SectionType)
    d.Set("stateful", firewall_section.Stateful)

    return nil
}

func resourceFirewallSectionUpdate(d *schema.ResourceData, m interface{}) error {

    nsxClient := m.(*api.APIClient)

    id := d.Id()
    if id == "" {
        return fmt.Errorf("Error obtaining logical object id")
    }

    revision := int64(d.Get("revision").(int))
    description := d.Get("description").(string)
    display_name := d.Get("display_name").(string)
    tags := getTagsFromSchema(d)
    applied_tos := getResourceReferencesFromSchema(d, "applied_tos")
    is_default := d.Get("is_default").(bool)
    rule_count := int64(d.Get("rule_count").(int))
    section_type := d.Get("section_type").(string)
    stateful := d.Get("stateful").(bool)
    firewall_section := manager.FirewallSection {
        Revision: revision,
        Description: description,
        DisplayName: display_name,
        Tags: tags,
        AppliedTos: applied_tos,
        IsDefault: is_default,
        RuleCount: rule_count,
        SectionType: section_type,
        Stateful: stateful,
    }

    firewall_section, resp, err := nsxClient.ServicesApi.UpdateSection(nsxClient.Context, id, firewall_section)

    if err != nil || resp.StatusCode == http.StatusNotFound {
        return fmt.Errorf("Error during FirewallSection %s update: %v", id, err)
    }

    return resourceFirewallSectionRead(d, m)
}

func resourceFirewallSectionDelete(d *schema.ResourceData, m interface{}) error {

    nsxClient := m.(*api.APIClient)

    id := d.Id()
    if id == "" {
        return fmt.Errorf("Error obtaining logical object id")
    }

    localVarOptionals := make(map[string]interface{})
    localVarOptionals["cascade"] = true
    resp, err := nsxClient.ServicesApi.DeleteSection(nsxClient.Context, id, localVarOptionals)
    if err != nil {
        return fmt.Errorf("Error during FirewallSection %s delete: %v", id, err)
    }

    if resp.StatusCode == http.StatusNotFound {
        fmt.Printf("FirewallSection %s not found", id)
        d.SetId("")
    }
return nil
}
