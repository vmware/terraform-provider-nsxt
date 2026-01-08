---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_cluster_security_config"
description: A resource to configure cluster security features like DFW.
---

# nsxt_policy_cluster_security_config

This resource provides a method to configure security features on compute clusters in NSX-T. It is primarily used to enable or disable Distributed Firewall (DFW) on a cluster, which is a prerequisite for features like IDPS (Intrusion Detection and Prevention System).

This resource manages existing cluster security configuration that is automatically created by NSX when a cluster is added. It cannot create new cluster configurations or delete them - it can only update the security feature settings.

~> **Note:** This resource is available for NSX-T 9.1.0 and above.

## Example Usage

### Enable DFW on a Cluster

```hcl
data "nsxt_compute_collection" "cluster1" {
  display_name = "Compute-Cluster-01"
}

resource "nsxt_policy_cluster_security_config" "cluster1_dfw" {
  cluster_id  = data.nsxt_compute_collection.cluster1.id
  dfw_enabled = true
}
```

### Use with IDPS Cluster Config

```hcl
data "nsxt_compute_collection" "cluster1" {
  display_name = "Compute-Cluster-01"
}

# Enable DFW on the cluster (required for IDPS)
resource "nsxt_policy_cluster_security_config" "cluster1_security" {
  cluster_id  = data.nsxt_compute_collection.cluster1.id
  dfw_enabled = true
}

# Create IDPS cluster configuration (depends on DFW being enabled)
resource "nsxt_policy_idps_cluster_config" "cluster1_idps" {
  depends_on = [nsxt_policy_cluster_security_config.cluster1_security]

  display_name = "cluster1-idps"
  ids_enabled  = true

  cluster {
    target_id   = data.nsxt_compute_collection.cluster1.id
    target_type = "VC_Cluster"
  }
}
```

### Disable DFW on a Cluster

```hcl
data "nsxt_compute_collection" "cluster1" {
  display_name = "Compute-Cluster-01"
}

resource "nsxt_policy_cluster_security_config" "cluster1_dfw" {
  cluster_id  = data.nsxt_compute_collection.cluster1.id
  dfw_enabled = false
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The cluster external ID (e.g., "uuid:domain-c20"). This is typically obtained from the `nsxt_compute_collection` data source's `id` attribute. Changing this forces a new resource to be created.
* `dfw_enabled` - (Optional) Enable or disable Distributed Firewall (DFW) on the cluster. DFW must be enabled before IDPS can be configured. Defaults to the current state if not specified.

## Attribute Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - The cluster ID.
* `display_name` - The display name of the cluster security configuration.
* `description` - The description of the cluster security configuration.
* `path` - The NSX path of the cluster security configuration.
* `revision` - The revision number of the cluster security configuration.

## Importing

An existing cluster security configuration can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```shell
terraform import nsxt_policy_cluster_security_config.cluster1_dfw <cluster-id>
```

The above command imports the cluster security configuration for the cluster with ID `<cluster-id>`. The cluster ID is the external ID of the compute collection (e.g., "uuid:domain-c20").

Example:

```shell
terraform import nsxt_policy_cluster_security_config.cluster1_dfw "5c9e6be8-ef56-4987-acff-c8e29e524e5c:domain-c20"
```

## Important Notes

### Cluster Configuration Lifecycle

* **Auto-created**: Cluster security configuration is automatically created by NSX when a cluster is added to NSX-T. This resource does not create new configurations.
* **Cannot delete**: Cluster security configurations cannot be deleted. When this resource is destroyed, it will disable all configured features but the configuration itself will remain in NSX.
* **Update only**: This resource manages the security feature settings of an existing cluster configuration.

### DFW Prerequisites

* **NSX Version**: DFW cluster configuration management requires NSX-T 9.1.0 or higher.
* **Cluster Registration**: The cluster must be registered with NSX-T before you can configure its security features.
* **IDPS Dependency**: If you plan to use IDPS features, DFW must be enabled on the cluster first.

### Best Practices

1. **Use with Data Sources**: Always obtain the cluster ID from the `nsxt_compute_collection` data source's `id` attribute to ensure you're using the correct external ID format.

2. **Dependencies**: When using with IDPS or other features that require DFW, use `depends_on` to ensure proper ordering:

   ```hcl
   resource "nsxt_policy_idps_cluster_config" "example" {
     depends_on = [nsxt_policy_cluster_security_config.example]
     # ...
   }
   ```

3. **State Management**: If DFW is already enabled on a cluster outside of Terraform, you can import the existing configuration to bring it under Terraform management.

### Troubleshooting

#### Error: "Cluster not found"

* Verify the cluster is registered with NSX-T
* Check that you're using the correct `external_id` from the compute collection data source

#### Error: "Feature not supported"

* Verify you're running NSX-T 9.1.0 or higher
* Check that the cluster is properly prepared for the feature

#### IDPS fails to enable

* Ensure DFW is enabled first
* Use `depends_on` to establish the dependency relationship
* Verify the cluster has sufficient resources

## See Also

* [nsxt_compute_collection](../data-sources/compute_collection.md) data source
* [nsxt_policy_idps_cluster_config](./policy_idps_cluster_config.md) resource
* [nsxt_policy_cluster_security_config](../data-sources/policy_cluster_security_config.md) data source
