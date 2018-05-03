---
layout: "nsxt"
page_title: "NSXT: nsxt_logical_tier0 router"
sidebar_current: "docs-nsxt-resource-logical-tier0-router"
description: A resource to configure a logical Tier0 router in NSX.
---

# nsxt_logical_tier0_router

This resource provides a method for the management of a tier 0 logical router.

## Example Usage

```hcl
resource "nsxt_logical_tier0_router" "tier0_router" {
  display_name           = "RTR"
  description            = "ACTIVE-STANDBY Tier0 router provisioned by Terraform"
  high_availability_mode = "ACTIVE_STANDBY"
  edge_cluster_id        = "${data.nsxt_edge_cluster.edge_cluster.id}"

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description of the resource.
* `edge_cluster_id` - (Optional) Edge Cluster ID for the logical Tier0 router.
* `tag` - (Optional) A list of scope + tag pairs to associate with this logical Tier0 router.
* `high_availability_mode` - (Optional) High availability mode "ACTIVE_ACTIVE"/"ACTIVE_STANDBY".

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical Tier0 router.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `advertise_config_revision` - Indicates current revision number of the advertisement configuration object as seen by NSX-T API server. This attribute can be useful for debugging.
* `firewall_sections` - (Optional) The list of firewall sections for this router

## Importing

An existing logical tier0 router can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_logical_tier0_router.tier0_router UUID
```

The above command imports the logical tier 0 router named `tier0_router` with the NSX id `UUID`.
