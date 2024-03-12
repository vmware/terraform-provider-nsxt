---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_metadata_proxy"
description: A resource to configure a metadata proxy in NSX Policy.
---

# nsxt_policy_metadata_proxy

This resource provides a method for the management of a metadata proxy which can be used within NSX Policy.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_metadata_proxy" "test" {
  display_name      = "metadata_proxy"
  description       = "TF configure MD proxy"
  edge_cluster_path = data.nsxt_policy_edge_cluster.test.path
  secret            = "Confidential!"
  server_address    = "http://192.168.1.1:6400"
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `crypto_protocols` - (Optional) Metadata proxy supported cryptographic protocols.
* `edge_cluster_path` - (Required) Policy path to Edge Cluster.
* `enable_standby_relocation` - (Optional) Flag to enable standby relocation. Default is false.
* `preferred_edge_paths` - (Optional) Preferred Edge Paths.
* `secret` - (Required) Secret.
* `server_address` - (Required) Server Address. This field is a URL. Example formats - http://1.2.3.4:3888/path, http://text-md-proxy:5001/. Port number should be between 3000-9000.
* `server_certificates` - (Optional) Policy paths to Certificate Authority (CA) certificates.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful
  for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_metadata_proxy.proxy UUID
```

The above command imports metadata proxy named `proxy` with the NSX ID `UUID`.

```
terraform import nsxt_policy_metadata_proxy.proxy POLICY_PATH
```

The above command imports the metadata proxy named `proxy` with policy path `POLICY_PATH`.
