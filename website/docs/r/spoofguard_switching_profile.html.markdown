---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_spoofguard_switching_profile"
description: |-
  Provides a resource to configure spoofguard switching profile on NSX-T manager
---

# nsxt_spoofguard_switching_profile

Provides a resource to configure spoofguard switching profile on NSX-T manager

## Example Usage

```hcl
resource "nsxt_spoofguard_switching_profile" "spoofguard_switching_profile" {
  description                       = "spoofguard_switching_profile provisioned by Terraform"
  display_name                      = "spoofguard_switching_profile"
  address_binding_whitelist_enabled = "true"

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
* `tag` - (Optional) A list of scope + tag pairs to associate with this spoofguard switching profile.
* `address_binding_whitelist_enabled` - (Optional) A boolean flag indicating whether this profile overrides the default system wide settings for Spoof Guard when assigned to ports.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the spoofguard switching profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing spoofguard switching profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_spoofguard_switching_profile.spoofguard_switching_profile UUID
```

The above would import the spoofguard switching profile named `spoofguard_switching_profile` with the nsx id `UUID`
