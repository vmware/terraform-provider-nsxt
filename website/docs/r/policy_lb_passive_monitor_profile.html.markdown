---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_lb_passive_monitor_profile"
description: A resource to configure a LBPassiveMonitorProfile.
---

# nsxt_policy_lb_passive_monitor_profile

This resource provides a method for the management of a LBPassiveMonitorProfile.

This resource is applicable NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_lb_passive_monitor_profile" "test" {
  display_name = "test"
  description  = "Terraform provisioned LBPassiveMonitorProfile"
  max_fails    = 2
  timeout      = 2

}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `max_fails` - (Optional) Number of consecutive failures before a member is considered temporarily unavailable.
* `timeout` - (Optional) After this timeout period, the member is tried again for a new connection to see if it is available.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_lb_passive_monitor_profile.test UUID
```

The above command imports LBPassiveMonitorProfile named `test` with the NSX ID `UUID`.
