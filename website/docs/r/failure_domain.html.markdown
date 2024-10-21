---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_failure_domain"
description: A resource to configure a Failure Domain.
---

# nsxt_failure_domain

This resource provides a method for the management of a Failure Domain.

This resource is supported with NSX 3.0.0 onwards.

## Example Usage

```hcl
resource "nsxt_failure_domain" "failure_domain" {
  display_name            = "failuredomain1"
  description             = "failuredomain"
  preferred_edge_services = "active"
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this failure domain.
* `preferred_edge_services` - Set preference for failure domain. Set preference for edge transport node failure domain which will be considered while doing auto placement of logical router, DHCP and MDProxy on edge node. `active`: For preemptive failover mode, active edge cluster member allocation prefers this failure domain. `standby`: For preemptive failover mode, standby edge cluster member allocation prefers this failure domain. Default will be `no_preference`. It means no explicit preference.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing Failure Domain can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_failure_domain.failure_domain UUID
```
The above command imports Failure Domain named `test` with the NSX Failure Domain ID `UUID`.
