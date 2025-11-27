---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_security_policy_container_cluster"
description: A resource to configure a Security Policy Container Cluster.
---

# nsxt_policy_security_policy_container_cluster

This resource provides a method for the management of Container Clusters associated with security policies.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_container_cluster" "cluster" {
  display_name = "containercluster1"
}

resource "nsxt_policy_parent_security_policy" "policy1" {
  display_name = "policy1"
  category     = "Application"
}

resource "nsxt_policy_security_policy_container_cluster" "antreacluster" {
  display_name           = "cluster1"
  description            = "Terraform provisioned SecurityPolicyContainerCluster"
  policy_path            = nsxt_policy_parent_security_policy.policy1.path
  container_cluster_path = data.nsxt_policy_container_cluster.cluster.path
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `policy_path` - (Required) The path of the Security Policy which the object belongs to
* `container_cluster_path` - (Required) Path to the container cluster entity in NSX

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```shell
terraform import nsxt_policy_security_policy_container_cluster.antreacluster PATH
```

The above command imports Security Policy Container Cluster named `antreacluster` with the NSX path `PATH`.
