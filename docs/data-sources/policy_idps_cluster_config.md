---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_idps_cluster_config"
description: A data source to read IDPS cluster configuration information.
---

# nsxt_policy_idps_cluster_config

This data source provides information about an IDPS cluster configuration in NSX Policy manager.

This data source is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_idps_cluster_config" "cluster1" {
  display_name = "cluster1-idps-config"
}
```

## Argument Reference

* `id` - (Optional) The ID of the IDPS cluster configuration to retrieve.
* `display_name` - (Optional) The Display Name prefix of the IDPS cluster configuration to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `ids_enabled` - Whether IDPS is enabled on this cluster.
* `cluster` - Cluster reference configuration.
    * `target_id` - Cluster target ID (e.g., domain-c123).
    * `target_type` - Target type (e.g., VC_Cluster).
