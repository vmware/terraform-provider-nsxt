---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: switching_profile"
description: A switching profile data source.
---

# nsxt_switching_profile

The switching profile data source provides information about switching profiles configured in NSX. A switching profile is a template that defines the settings of one or more logical switches. There can be both factory default and user defined switching profiles. One example of a switching profile is a quality of service (QoS) profile which defines the QoS settings of all switches that use the defined switch profile.

## Example Usage

```hcl
data "nsxt_switching_profile" "qos_profile" {
  display_name = "qos-profile"
}
```

## Argument Reference

* `id` - (Optional) The ID of Switching Profile to retrieve.
* `display_name` - (Optional) The Display Name of the Switching Profile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `resource_type` - The resource type representing the specific type of this switching profile.
* `description` - The description of the switching profile.
