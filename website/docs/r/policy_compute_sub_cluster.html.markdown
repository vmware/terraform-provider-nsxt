---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_compute_sub_cluster"
description: A resource to configure a Compute Sub Cluster.
---

# nsxt_policy_compute_sub_cluster

This resource provides a method for the management of Compute Sub Cluster.
This resource is supported with NSX 3.2.2 onwards.

## Example Usage

```hcl
resource "nsxt_policy_compute_sub_cluster" "test" {
  display_name          = "subcluster1"
  compute_collection_id = data.nsxt_compute_collection.cc1.id
  discovered_node_ids   = [data.nsxt_discovered_node.dn.id]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `site_path` - (Optional) The path of the site which this Sub Cluster belongs to. `path` field of the existing `nsxt_policy_site` can be used here. Defaults to default site path.
* `enforcement_point` - (Optional) The ID of enforcement point under given `site_path` to manage the Host Transport Node. Defaults to default enforcement point.
* `compute_collection_id` - (Required) ID of compute collection under which sub-cluster is created
* `discovered_node_ids` - (Optional)  Discovered node IDs under this subcluster

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `discovered_node_ids` - (Optional)  Discovered node IDs under this subcluster

## Importing

An existing Sub Cluster can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_compute_sub_cluster.test POLICY_PATH
```
The above command imports Sub Cluster named `test` with the policy path `POLICY_PATH`.
