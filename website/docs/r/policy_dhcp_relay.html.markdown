---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_dhcp_relay"
sidebar_current: "docs-nsxt-resource-policy-dhcp-relay"
description: A resource to configure a Dhcp Relay.
---

# nsxt_policy_dhcp_relay

This resource provides a method for the management of a Dhcp Relay.
 
## Example Usage

```hcl
resource "nsxt_policy_dhcp_relay" "test" {
    display_name      = "test"
    description       = "Terraform provisioned Dhcp Relay"
    server_addresses = ["10.0.0.2", "7001::2"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `server_addresses` - (Required) List of DHCP server addresses.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_dhcp_relay.test ID
```

The above command imports Dhcp Relay named `test` with the NSX Dhcp Relay ID `ID`.
