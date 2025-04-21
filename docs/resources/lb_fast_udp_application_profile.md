---
subcategory: "Deprecated"
page_title: "NSXT: nsxt_lb_fast_udp_application_profile"
description: |-
  Provides a resource to configure LB fast UDP application profile on NSX manager
---

# nsxt_lb_fast_udp_application_profile

Provides a resource to configure LB fast UDP application profile on NSX manager

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usage

```hcl
resource "nsxt_lb_fast_udp_application_profile" "lb_fast_udp_profile" {
  description       = "lb_fast_udp_application_profile provisioned by Terraform"
  display_name      = "lb_fast_udp_application_profile"
  idle_timeout      = "1800"
  ha_flow_mirroring = "false"

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
* `idle_timeout` - (Optional) Timeout in seconds to specify how long an idle UDP connection in ESTABLISHED state should be kept for this application before cleaning up. The default value will be 300 seconds
* `ha_flow_mirroring` - (Optional) A boolean flag which reflects whether flow mirroring is enabled, and all the flows to the bounded virtual server are mirrored to the standby node. By default this is disabled.
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb fast udp profile.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb fast udp profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing lb fast udp profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_lb_fast_udp_application_profile.lb_fast_udp_profile UUID
```

The above would import the LB fast UDP application profile named `lb_fast_udp_profile` with the nsx id `UUID`
