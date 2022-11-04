---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_mac_management_switching_profile"
description: |-
  Provides a resource to configure MAC management switching profile on NSX-T manager
---

# nsxt_mac_management_switching_profile

Provides a resource to configure MAC management switching profile on NSX-T manager

## Example Usage

```hcl
resource "nsxt_mac_management_switching_profile" "mac_management_switching_profile" {
  description        = "mac_management_switching_profile provisioned by Terraform"
  display_name       = "mac_management_switching_profile"
  mac_change_allowed = "true"

  mac_learning {
    enabled                  = "true"
    limit                    = "4096"
    limit_policy             = "ALLOW"
    unicast_flooding_allowed = "false"
  }

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this MAC management switching profile.
* `mac_change_allowed` - (Optional) A boolean flag indicating allowing source MAC address change.
* `mac_learning` - (Optional) Mac learning configuration:
  * `enabled` - (Optional) A boolean flag indicating allowing source MAC address learning.
  * `unicast_flooding_allowed` - (Optional) A boolean flag indicating allowing flooding for unlearned MAC for ingress traffic. Can be True only if mac_learning is enabled.
  * `limit` - (Optional) The maximum number of MAC addresses that can be learned on this port.
  * `limit_policy` - (Optional) The policy after MAC Limit is exceeded: ALLOW/DROP.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the MAC management switching profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing MAC management switching profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_mac_management_switching_profile.mac_management_switching_profile UUID
```

The above would import the MAC management switching profile named `mac_management_switching_profile` with the nsx id `UUID`
