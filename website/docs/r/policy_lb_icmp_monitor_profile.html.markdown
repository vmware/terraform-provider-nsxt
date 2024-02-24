---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_lb_icmp_monitor_profile"
description: A resource to configure a LBIcmpMonitorProfile.
---

# nsxt_policy_lb_icmp_monitor_profile

This resource provides a method for the management of a LBIcmpMonitorProfile.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_lb_icmp_monitor_profile" "test" {
  display_name = "test"
  description  = "Terraform provisioned LBIcmpMonitorProfile"
  data_length  = 2
  fall_count   = 2
  interval     = 2
  rise_count   = 2
  timeout      = 2

}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `data_length` - (Optional) The data size (in bytes) of the ICMP healthcheck packet.
* `fall_count` - (Optional) Mark member status DOWN if the healtcheck fails consecutively for fall_count times.
* `interval` - (Optional) Active healthchecks are initiated periodically, at a configurable interval (in seconds), to each member of the Group.
* `rise_count` - (Optional) Bring a DOWN member UP if rise_count successive healthchecks succeed.
* `timeout` - (Optional) Timeout specified in seconds. After a healthcheck is initiated, if it does not complete within a certain period, then also the healthcheck is considered to be unsuccessful. Completing a healthcheck within timeout means establishing a connection (TCP or SSL), if applicable, sending the request and receiving the response, all within the configured timeout.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_lb_icmp_monitor_profile.test UUID
```

The above command imports LBIcmpMonitorProfile named `test` with the NSX ID `UUID`.
