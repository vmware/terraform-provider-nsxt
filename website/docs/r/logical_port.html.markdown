---
layout: "nsxt"
page_title: "NSXT: nsxt_logical_port"
sidebar_current: "docs-nsxt-resource-logical-port"
description: |-
  Provides a resource to configure Logical Port (LP) on NSX-T Manager.
---

# nsxt_logical_port

Provides a resource to configure Logical Port (LP) on NSX-T Manager.

## Example Usage

```hcl
resource "nsxt_logical_port" "logical_port" {
    admin_state = "UP"
    description = "LP1 provisioned by Terraform"
    display_name = "LP1"
    logical_switch_id = "${nsxt_logical_switch.switch1.id}"
    tag {
        scope = "color"
        tag = "blue"
    }
    switching_profile_id {
        key = "${data.nsxt_switching_profile.qos_profile.resource_type}",
        value = "${data.nsxt_switching_profile.qos_profile.id}"
    }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description of this resource.
* `logical_switch_id` - (Required) Logical switch ID for the logical port.
* `admin_state` - (Required) Admin state for the logical port. Accepted values - 'UP' or 'DOWN'.
* `switching_profile_id` - (Optional) List of IDs of switching profiles (of various types) to be associated with this switch. Default switching profiles will be used if not specified.
* `tag` - (Optional) A list of scope + tag pairs to associate with this logical switch.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical switch.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing Logical Port can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import nsxt_logical_port.x id
```

The above would import the Logical Port named `x` with the nsx id `id`
