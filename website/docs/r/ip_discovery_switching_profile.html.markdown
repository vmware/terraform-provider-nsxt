---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_ip_discovery_switching_profile"
description: |-
  Provides a resource to configure IP discovery switching profile on NSX-T manager
---

# nsxt_ip_discovery_switching_profile

Provides a resource to configure IP discovery switching profile on NSX-T manager

## Example Usage

```hcl
resource "nsxt_ip_discovery_switching_profile" "ip_discovery_switching_profile" {
  description           = "ip_discovery_switching_profile provisioned by Terraform"
  display_name          = "ip_discovery_switching_profile"
  vm_tools_enabled      = "false"
  arp_snooping_enabled  = "true"
  dhcp_snooping_enabled = "false"
  arp_bindings_limit    = "1"

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
* `tag` - (Optional) A list of scope + tag pairs to associate with this IP discovery switching profile.
* `arp_snooping_enabled` - (Optional) A boolean flag iIndicates whether ARP snooping is enabled.
* `vm_tools_enabled` - (Optional) A boolean flag iIndicates whether VM tools will be enabled. This option is only supported on ESX where vm-tools is installed.
* `dhcp_snooping_enabled` - (Optional) A boolean flag iIndicates whether DHCP snooping is enabled.
* `arp_bindings_limit` - (Optional) Limit for the amount of ARP bindings.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the IP discovery switching profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing IP discovery switching profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_ip_discovery_switching_profile.ip_discovery_switching_profile UUID
```

The above would import the IP discovery switching profile named `ip_discovery_switching_profile` with the nsx id `UUID`
