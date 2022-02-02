---
subcategory: "Policy - Segments"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_mac_discovery_profile"
description: A resource to configure a MacDiscoveryProfile.
---

# nsxt_policy_mac_discovery_profile

This resource provides a method for the management of a MacDiscoveryProfile.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_mac_discovery_profile" "test" {
  display_name                     = "test"
  description                      = "Terraform provisioned MacDiscoveryProfile"
  mac_change_enabled               = true
  mac_learning_enabled             = true
  mac_limit                        = 4096
  mac_limit_policy                 = "ALLOW"
  remote_overlay_mac_limit         = 2048
  unknown_unicast_flooding_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `mac_change_enabled` - (Optional) Allowing source MAC address learning
* `mac_learning_enabled` - (Optional) Allowing source MAC address learning
* `mac_limit` - (Optional) The maximum number of MAC addresses that can be learned on this port
* `mac_limit_policy` - (Optional) The policy after MAC Limit is exceeded
* `remote_overlay_mac_limit` - (Optional) The maximum number of MAC addresses learned on an overlay Logical Switch
* `unknown_unicast_flooding_enabled` - (Optional) Allowing flooding for unlearned MAC for ingress traffic


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_mac_discovery_profile.test UUID
```

The above command imports MacDiscoveryProfile named `test` with the NSX MacDiscoveryProfile ID `UUID`.
