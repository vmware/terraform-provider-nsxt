---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_lb_tcp_monitor_profile"
description: A resource to configure a LBTcpMonitorProfile.
---

# nsxt_policy_lb_tcp_monitor_profile

This resource provides a method for the management of a LBTcpMonitorProfile.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_lb_tcp_monitor_profile" "test" {
  display_name = "test"
  description  = "Terraform provisioned LBTcpMonitorProfile"
  receive      = "test"
  send         = "test"
  fall_count   = 2
  interval     = 2
  monitor_port = 8080
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
* `receive` - (Optional) The expected data string to be received from the response, can be anywhere in the response.
* `send` - (Optional) The data to be sent to the monitored server.
* `fall_count` - (Optional) Mark member status DOWN if the healtcheck fails consecutively for fall_count times.
* `interval` - (Optional) Active healthchecks are initiated periodically, at a configurable interval (in seconds), to each member of the Group.
* `monitor_port` - (Optional) Typically, monitors perform healthchecks to Group members using the member IP address and pool_port. However, in some cases, customers prefer to run healthchecks against a different port than the pool member port which handles actual application traffic. In such cases, the port to run healthchecks against can be specified in the monitor_port value.
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
terraform import nsxt_policy_lb_tcp_monitor_profile.test UUID
```

The above command imports LBTcpMonitorProfile named `test` with the NSX ID `UUID`.
