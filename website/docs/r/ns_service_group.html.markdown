---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_ns_service_group"
description: |-
  Provides a resource to configure NS service group on NSX-T manager
---

# nsxt_ns_service_group

Provides a resource to configure NS service group on NSX-T manager

## Example Usage

```hcl

data "nsxt_ns_service" "dns" {
  display_name = "DNS"
}

resource "nsxt_ip_protocol_ns_service" "prot17" {
  display_name = "ip_prot"
  protocol     = "17"
}

resource "nsxt_ns_service_group" "ns_service_group" {
  description  = "ns_service_group provisioned by Terraform"
  display_name = "ns_service_group"
  members      = [nsxt_ip_protocol_ns_service.prot17.id, data.nsxt_ns_service.dns.id]

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
* `tag` - (Optional) A list of scope + tag pairs to associate with this NS service group.
* `members` - (Required) List of NSServices IDs that can be added as members to an NSServiceGroup. All members should be of the same L2 type: Ethernet, or Non Ethernet.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the NS service group.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing ns service group can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_ns_service_group.ns_service_group UUID
```

The above would import the NS service group named `ns_service_group` with the nsx id `UUID`
