---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_cluster_security_config"
description: A data source to read cluster security configuration.
---

# nsxt_policy_cluster_security_config

This data source provides information about cluster security configuration in NSX-T. It can be used to check the status of security features like Distributed Firewall (DFW) on a compute cluster.

~> **Note:** This data source is available for NSX-T 9.1.0 and above.

## Example Usage

### Read Cluster Security Configuration

```hcl
data "nsxt_compute_collection" "cluster1" {
  display_name = "Compute-Cluster-01"
}

data "nsxt_policy_cluster_security_config" "cluster1_security" {
  cluster_id = data.nsxt_compute_collection.cluster1.id
}

output "dfw_enabled" {
  value = data.nsxt_policy_cluster_security_config.cluster1_security.dfw_enabled
}
```

### Check DFW Status Before Creating IDPS Config

```hcl
data "nsxt_compute_collection" "cluster1" {
  display_name = "Compute-Cluster-01"
}

data "nsxt_policy_cluster_security_config" "cluster1_security" {
  cluster_id = data.nsxt_compute_collection.cluster1.id
}

# Only create IDPS config if DFW is enabled
resource "nsxt_policy_idps_cluster_config" "cluster1_idps" {
  count = data.nsxt_policy_cluster_security_config.cluster1_security.dfw_enabled ? 1 : 0

  display_name = "cluster1-idps"
  ids_enabled  = true

  cluster {
    target_id   = data.nsxt_compute_collection.cluster1.id
    target_type = "VC_Cluster"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The cluster external ID (e.g., "uuid:domain-c20"). This is typically obtained from the `nsxt_compute_collection` data source's `id` attribute.

## Attribute Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - The cluster ID.
* `display_name` - The display name of the cluster security configuration.
* `description` - The description of the cluster security configuration.
* `path` - The NSX path of the cluster security configuration.
* `dfw_enabled` - Whether Distributed Firewall (DFW) is enabled on the cluster.

## Example Outputs

### Check if DFW is Enabled

```hcl
output "is_dfw_enabled" {
  description = "Whether DFW is enabled on the cluster"
  value       = data.nsxt_policy_cluster_security_config.cluster1_security.dfw_enabled
}
```

### Conditional Resource Creation

```hcl
resource "null_resource" "dfw_check" {
  count = data.nsxt_policy_cluster_security_config.cluster1_security.dfw_enabled ? 0 : 1

  provisioner "local-exec" {
    command = "echo 'Warning: DFW is not enabled on the cluster'"
  }
}
```

## Important Notes

### Prerequisites

* **NSX Version**: Requires NSX-T 9.1.0 or higher
* **Cluster Registration**: The cluster must be registered with NSX-T
* **Permissions**: Requires read access to cluster security configuration

### Use Cases

This data source is useful for:

1. **Validation**: Checking if DFW is enabled before creating IDPS resources
2. **Conditional Logic**: Making decisions based on DFW status
3. **Monitoring**: Tracking DFW enablement status across clusters
4. **Documentation**: Generating reports of cluster security configurations

### Troubleshooting

#### Error: "Cluster not found"

* Verify the cluster is registered with NSX-T
* Check that you're using the correct `id` from the compute collection data source
* Ensure you have read permissions for the cluster

#### Error: "API not supported"

* Verify you're running NSX-T 9.1.0 or higher
* Check that the cluster security configuration API is available

## See Also

* [nsxt_compute_collection](./compute_collection.md) data source
* [nsxt_policy_cluster_security_config](../resources/policy_cluster_security_config.md) resource
* [nsxt_policy_idps_cluster_config](../resources/policy_idps_cluster_config.md) resource
