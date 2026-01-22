---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_idps_cluster_config"
description: A resource to configure IDPS cluster activation/deactivation.
---

# nsxt_policy_idps_cluster_config

This resource provides a method for the management of IDPS (Intrusion Detection and Prevention System) cluster configuration. It allows you to activate or deactivate IDPS at the cluster level.

This resource is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
resource "nsxt_policy_idps_cluster_config" "cluster1" {
  display_name = "cluster1-idps-config"
  description  = "IDPS configuration for cluster1"
  ids_enabled  = true

  cluster {
    target_id   = "domain-c123"
    target_type = "VC_Cluster"
  }

  tag {
    scope = "environment"
    tag   = "production"
  }

  tag {
    scope = "team"
    tag   = "security"
  }

  tag {
    scope = "cost-center"
    tag   = "cc-1001"
  }
}
```

## Example Usage - Multi-Tenancy

```hcl

resource "nsxt_policy_idps_cluster_config" "cluster1" {
  display_name = "cluster1-idps-config"
  description  = "IDPS configuration for cluster1"
  ids_enabled  = true

  cluster {
    target_id   = "domain-c123"
    target_type = "VC_Cluster"
  }

  tag {
    scope = "environment"
    tag   = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `ids_enabled` - (Required) Enable or disable IDS for the cluster (TF-IDPS-001).
* `cluster` - (Required) Cluster reference configuration. Exactly one cluster block is required.
    * `target_id` - (Required) Cluster target ID (e.g., domain-c123). This cannot be changed after creation.
    * `target_type` - (Required) Target type (e.g., VC_Cluster). This cannot be changed after creation.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the IDPS cluster configuration. This will be the same as the cluster `target_id`.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX policy path of the resource.

## Importing

An existing IDPS cluster configuration can be [imported][docs-import] into this resource, via the following command:

```shell
terraform import nsxt_policy_idps_cluster_config.cluster1 3e13d462-88ff-4313-8e7b-56ac0aafd917:domain-c46
```

The above command imports the IDPS cluster configuration named `cluster1` with the cluster target ID `3e13d462-88ff-4313-8e7b-56ac0aafd917:domain-c46`.

[docs-import]: https://www.terraform.io/docs/import/

## Notes

* This resource manages IDPS activation at the cluster level
* The `target_id` and `target_type` in the cluster block cannot be changed after resource creation
* Disabling IDPS on a cluster (`ids_enabled = false`) will stop all intrusion detection and prevention activities for that cluster
* The `target_id` should be the vCenter cluster identifier (e.g., domain-c123)
* The `target_type` should typically be "VC_Cluster" for vCenter clusters
* The resource ID will automatically be set to the cluster `target_id` as required by NSX-T
* Only one IDPS cluster configuration can exist per cluster target ID
