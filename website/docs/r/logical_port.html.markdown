---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_logical_port"
description: A resource that can be used to configure a Logical Port in NSX.
---

# nsxt_logical_port

This resource provides a resource to configure a logical port on a logical switch in the NSX system. Like physical switches a logical switch can have one or more ports which can be connected to virtual machines or logical routers.

## Example Usage

```hcl
resource "nsxt_logical_port" "logical_port" {
  admin_state       = "UP"
  description       = "LP1 provisioned by Terraform"
  display_name      = "LP1"
  logical_switch_id = nsxt_logical_switch.switch1.id

  tag {
    scope = "color"
    tag   = "blue"
  }

  switching_profile_id {
    key   = data.nsxt_switching_profile.qos_profile.resource_type
    value = data.nsxt_switching_profile.qos_profile.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description of this resource.
* `logical_switch_id` - (Required) Logical switch ID for the logical port.
* `admin_state` - (Optional) Admin state for the logical port. Accepted values - 'UP' or 'DOWN'. The default value is 'UP'.
* `switching_profile_id` - (Optional) List of IDs of switching profiles (of various types) to be associated with this switch. Default switching profiles will be used if not specified.
* `tag` - (Optional) A list of scope + tag pairs to associate with this logical port.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical port.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing Logical Port can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_logical_port.logical_port UUID
```

The above command imports the logical port named `logical_port` with the NSX id `UUID`.
