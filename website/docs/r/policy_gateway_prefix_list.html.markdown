---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_gateway_prefix_list"
description: A resource to configure a Tier 0 Gateway Prefix Liston NSX Policy manager.
---

# nsxt_policy_gateway_prefix_list

This resource provides a method for the management of a Tier-0 Gateway Prefix List.

# Example Usage

```hcl
data "nsxt_policy_tier0_gateway" "gw1" {
  display_name = "gw1"
}

resource "nsxt_policy_gateway_prefix_list" "pf1" {
  gateway_path = data.nsxt_policy_tier0_gateway.gw1.path
  display_name = "t0_prefix_list"
  description  = "Prefix list for tier0 GW gw1"

  prefix {
    action  = "PERMIT"
    network = "4.4.0.0/20"
    le      = 23
    ge      = 20
  }

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `gateway_path` - (Required) Gateway where the prefix is defined.
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `prefix` - (Required) A list of network prefixes.
  * `action` - (Optional) PERMIT or DENY Action for the prefix list. The default value is PERMIT.
  * `le` - (Optional) Prefix length less than or equal to, between 0-128. (0 means no value)
  * `ge` - (Optional) Prefix length greater than or equal to, between 0-128. (0 means no value).
  * `network` - (Optional) Network prefix in CIDR format. If not set it will match ANY network.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing policy Tier-0 Gateway prefix list can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_gateway_prefix_list.pf1 GW-ID/ID
```

The above command imports the policy Tier-0 gateway prefix list named `pf1` with the NSX Policy ID `ID` on Tier0 Gateway `GW-ID`.
