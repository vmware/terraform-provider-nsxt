---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_lb_fast_tcp_application_profile"
description: A resource to configure a Load Balancer Fast TCP Application Profile.
---

# nsxt_policy_lb_fast_tcp_application_profile

This resource provides a method for the management of a Load Balancer Fast TCP Application Profile.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_lb_fast_tcp_application_profile" "test" {
  display_name              = "test"
  description               = "Terraform provisioned profile"
  idle_timeout              = 120
  close_timeout             = 10
  ha_flow_mirroring_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `idle_timeout` - (Optional) Timeout in seconds to specify how long TCP connection can remain in ESTABLISHED state. Default is 1800.
* `close_timeout` - (Optional) Timeout in seconds to specify how long TCP connection can remain in closed state (both FINs received or a RST is received). Default is 8.
* `ha_flow_mirroring_enabled` - (Optional) If enabled, all the flows to the bounded virtual server are mirrored to the standby node. Default is False.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: fast_tcps://www.terraform.io/cli/import

```
terraform import nsxt_policy_lb_fast_tcp_application_profile.test PATH
```

The above command imports LB Application Profile named `test` with the NSX path `PATH`.
