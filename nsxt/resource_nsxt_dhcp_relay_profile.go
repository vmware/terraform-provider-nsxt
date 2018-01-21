package nsxt

import(
    "github.com/hashicorp/terraform/helper/schema"
    api "github.com/vmware/go-vmware-nsxt"
    "github.com/vmware/go-vmware-nsxt/manager"
    "net/http"
    "fmt"
)

func resourceDhcpRelayProfile() *schema.Resource {
    return &schema.Resource{
        Create: resourceDhcpRelayProfileCreate,
        Read: resourceDhcpRelayProfileRead,
        Update: resourceDhcpRelayProfileUpdate,
        Delete: resourceDhcpRelayProfileDelete,

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
            "server_addresses": &schema.Schema{
                Type:        schema.TypeSet,
                Description: "Set of dhcp relay server addresses",
                Elem:        &schema.Schema{Type: schema.TypeString},
                Required:    true,
            },
        },
    }
}

func resourceDhcpRelayProfileCreate(d *schema.ResourceData, m interface{}) error {

    nsxClient := m.(*api.APIClient)

    description := d.Get("description").(string)
    display_name := d.Get("display_name").(string)
    tags := getTagsFromSchema(d)
    server_addresses := getStringListFromSchemaSet(d, "server_addresses")
    dhcp_relay_profile := manager.DhcpRelayProfile {
        Description: description,
        DisplayName: display_name,
        Tags: tags,
        ServerAddresses: server_addresses,
    }

    dhcp_relay_profile, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateDhcpRelayProfile(nsxClient.Context, dhcp_relay_profile)

    if err != nil {
        return fmt.Errorf("Error during DhcpRelayProfile create: %v", err)
    }

    if resp.StatusCode != http.StatusCreated {
        fmt.Printf("Unexpected status returned")
        return nil
    }
    d.SetId(dhcp_relay_profile.Id)

    return resourceDhcpRelayProfileRead(d, m)
}

func resourceDhcpRelayProfileRead(d *schema.ResourceData, m interface{}) error {

    nsxClient := m.(*api.APIClient)

    id := d.Id()
    if id == "" {
        return fmt.Errorf("Error obtaining dhcp relay profile id")
    }

    dhcp_relay_profile, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadDhcpRelayProfile(nsxClient.Context, id)
    if resp.StatusCode == http.StatusNotFound {
        fmt.Printf("DhcpRelayProfile not found")
        d.SetId("")
        return nil
    }
    if err != nil {
        return fmt.Errorf("Error during DhcpRelayProfile read: %v", err)
    }

    d.Set("Revision", dhcp_relay_profile.Revision)
    d.Set("SystemOwned", dhcp_relay_profile.SystemOwned)
    d.Set("Description", dhcp_relay_profile.Description)
    d.Set("DisplayName", dhcp_relay_profile.DisplayName)
    setTagsInSchema(d, dhcp_relay_profile.Tags)
    d.Set("ServerAddresses", dhcp_relay_profile.ServerAddresses)

    return nil
}

func resourceDhcpRelayProfileUpdate(d *schema.ResourceData, m interface{}) error {

    nsxClient := m.(*api.APIClient)

    id := d.Id()
    if id == "" {
        return fmt.Errorf("Error obtaining dhcp relay profile id")
    }

    revision := int64(d.Get("revision").(int))
    description := d.Get("description").(string)
    display_name := d.Get("display_name").(string)
    tags := getTagsFromSchema(d)
    server_addresses := interface2StringList(d.Get("server_addresses").(*schema.Set).List())
    dhcp_relay_profile := manager.DhcpRelayProfile {
        Revision: revision,
        Description: description,
        DisplayName: display_name,
        Tags: tags,
        ServerAddresses: server_addresses,
    }

    dhcp_relay_profile, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateDhcpRelayProfile(nsxClient.Context, id, dhcp_relay_profile)

    if err != nil || resp.StatusCode == http.StatusNotFound {
        return fmt.Errorf("Error during DhcpRelayProfile update: %v", err)
    }

    return resourceDhcpRelayProfileRead(d, m)
}

func resourceDhcpRelayProfileDelete(d *schema.ResourceData, m interface{}) error {

    nsxClient := m.(*api.APIClient)

    id := d.Id()
    if id == "" {
        return fmt.Errorf("Error obtaining dhcp relay profile id")
    }

    resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteDhcpRelayProfile(nsxClient.Context, id)
    if err != nil {
        return fmt.Errorf("Error during DhcpRelayProfile delete: %v", err)
    }

    if resp.StatusCode == http.StatusNotFound {
        fmt.Printf("DhcpRelayProfile not found")
        d.SetId("")
    }
    return nil
}
