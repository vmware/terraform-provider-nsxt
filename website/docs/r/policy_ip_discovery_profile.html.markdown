---
subcategory: "Segments"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ip_discovery_profile"
description: A resource to configure an IP discovery profile.
---

# nsxt_policy_ip_discovery_profile

This resource provides a method for the management of IP discovery profiles.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_ip_discovery_profile" "ip_discovery_profile" {
  description  = "ip discovery profile provisioned by Terraform"
  display_name = "ip_discovery_profile1"

  arp_nd_binding_timeout         = 20
  duplicate_ip_detection_enabled = false

  arp_binding_limit     = 140
  arp_snooping_enabled  = false
  dhcp_snooping_enabled = false
  vmtools_enabled       = false

  dhcp_snooping_v6_enabled = false
  nd_snooping_enabled      = false
  nd_snooping_limit        = 12
  vmtools_v6_enabled       = false
  tofu_enabled             = false

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_ip_discovery_profile" "ip_discovery_profile" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  description  = "ip discovery profile provisioned by Terraform"
  display_name = "ip_discovery_profile1"

  arp_nd_binding_timeout         = 20
  duplicate_ip_detection_enabled = false

  arp_binding_limit     = 140
  arp_snooping_enabled  = false
  dhcp_snooping_enabled = false
  vmtools_enabled       = false

  dhcp_snooping_v6_enabled = false
  nd_snooping_enabled      = false
  nd_snooping_limit        = 12
  vmtools_v6_enabled       = false
  tofu_enabled             = false

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this policy.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `arp_nd_binding_timeout` - (Optional) ARP and ND cache timeout (in minutes)
* `duplicate_ip_detection_enabled` - (Optional) Duplicate IP detection
* `arp_binding_limit` - (Optional) Maximum number of ARP bindings
* `arp_snooping_enabled` - (Optional) Is ARP snooping enabled or not
* `dhcp_snooping_enabled` - (Optional) Is DHCP snooping enabled or not
* `vmtools_enabled` - (Optional) Is VM tools enabled or not
* `dhcp_snooping_v6_enabled` - (Optional)  Is DHCP snooping v6 enabled or not
* `nd_snooping_enabled` - (Optional) Is ND snooping enabled or not
* `nd_snooping_limit` - (Optional) Maximum number of ND (Neighbor Discovery Protocol) bindings
* `vmtools_v6_enabled` - (Optional) Use VMTools to learn the IPv6 addresses which are configured on interfaces of a VM
* `tofu_enabled` - (Optional) Indicates whether "Trust on First Use(TOFU)" paradigm is enabled

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ip_discovery_profile.ip_discovery_profile ID
```
The above command imports the IP discovery profile named `ip_discovery_profile` with the NSX ID `ID`.

```
terraform import nsxt_policy_ip_discovery_profile.ip_discovery_profile POLICY_PATH
```
The above command imports the IP discovery profile named `ip_discovery_profile` with policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
