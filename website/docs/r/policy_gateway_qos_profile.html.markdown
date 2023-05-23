---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_gateway_qos_profile"
description: A resource to configure a Gateway QoS Profile.
---

# nsxt_policy_gateway_qos_profile

This resource provides a method for the management of a Gateway QoS Profile.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_gateway_qos_profile" "test" {
  display_name        = "test"
  description         = "Terraform provisioned profile"
  burst_size          = 10
  committed_bandwidth = 20
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `burst_size` - (Optional) Maximum amount of traffic that can be transmitted at peak bandwidth rate (bytes)
* `committed_bandwidth` - (Optional) Bandwidth is limited to line rate when the value configured is greater than line rate (mbps)
* `excess_action` - (Optional) Action on traffic exceeding bandwidth. Currently only `DROP` is supported.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_gateway_qos_profile.test UUID
```

The above command imports profile named `test` with the NSX ID `UUID`.
