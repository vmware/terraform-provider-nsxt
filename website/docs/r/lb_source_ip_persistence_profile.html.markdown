---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_source_ip_persistence_profile"
description: |-
  Provides a resource to configure lb source ip persistence profile on NSX-T manager
---

# nsxt_lb_source_ip_persistence_profile

Provides a resource to configure lb source ip persistence profile on NSX-T manager

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usage

```hcl
resource "nsxt_lb_source_ip_persistence_profile" "lb_source_ip_persistence_profile" {
  description              = "lb_source_ip_persistence_profile provisioned by Terraform"
  display_name             = "lb_source_ip_persistence_profile"
  persistence_shared       = "true"
  ha_persistence_mirroring = "true"
  purge_when_full          = "true"
  timeout                  = "100"

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
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb source ip persistence profile.
* `persistence_shared` - (Optional) A boolean flag which reflects whether the cookie persistence is private or shared.
* `ha_persistence_mirroring` - (Optional) A boolean flag which reflects whether persistence entries will be synchronized to the HA peer.
* `timeout` - (Optional) Persistence expiration time in seconds, counted from the time all the connections are completed. Defaults to 300 seconds.
* `purge_when_full` - (Optional) A boolean flag which reflects whether entries will be purged when the persistence table is full. Defaults to true.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb source ip persistence profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb source ip persistence profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_source_ip_persistence_profile.lb_source_ip_persistence_profile UUID
```

The above would import the lb source ip persistence profile named `lb_source_ip_persistence_profile` with the nsx id `UUID`
