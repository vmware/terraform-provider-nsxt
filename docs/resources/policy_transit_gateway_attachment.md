---
subcategory: "VPC"
page_title: "NSXT: nsxt_policy_transit_gateway_attachment"
description: A resource to configure a Transit Gateway Attachment.
---

# nsxt_policy_transit_gateway_attachment

This resource provides a method for the management of a Transit Gateway Attachment.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_project" "test_proj" {
  display_name = "test_project"
}

data "nsxt_policy_transit_gateway" "test_tgw" {
  context {
    project_id = nsxt_policy_project.test_proj.id
  }
  id = "default"
}

resource "nsxt_policy_gateway_connection" "test_gw_conn" {
  display_name = "test_gw_conn"
  tier0_path   = "/infra/tier-0s/test-t0"
}

resource "nsxt_policy_transit_gateway_attachment" "test_tgw_att" {
  display_name    = "test"
  parent_path     = data.nsxt_policy_transit_gateway.test_tgw.path
  connection_path = nsxt_policy_gateway_connection.test_gw_conn.path

  admin_state = "UP"
  urpf_mode   = "STRICT"

  route_advertisement_types = ["PUBLIC", "TGW_PRIVATE"]
}
```

### Using CNA path instead of connection path

```hcl
resource "nsxt_policy_transit_gateway_attachment" "cna_att" {
  display_name = "test-cna"
  parent_path  = data.nsxt_policy_transit_gateway.test_tgw.path
  cna_path     = "/orgs/default/projects/proj1/vpcs/vpc1/subnets/sub1"
}
```

## Argument Reference

The following arguments are supported:

* `parent_path` - (Required) Policy path of the parent transit gateway.
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `connection_path` - (Optional) Policy path of the transit gateway external connection (e.g. a `nsxt_policy_gateway_connection`). Exactly one of `connection_path` or `cna_path` must be provided.
* `cna_path` - (Optional) Policy path to a centralized network attachment. Exactly one of `connection_path` or `cna_path` must be provided.
* `admin_state` - (Optional/Computed) Administrative state of the attachment interface. Accepted values: `UP`, `DOWN`.
* `urpf_mode` - (Optional/Computed) Unicast Reverse Path Forwarding mode. Accepted values: `NONE`, `STRICT`.
* `route_advertisement_types` - (Optional/Computed) List of route types to advertise from this attachment. Accepted values: `PUBLIC` (subnets and NAT IPs from external networks), `TGW_PRIVATE` (transit gateway private subnets, only advertised when `allow_private` is set on the BGP neighbor). Defaults to `["PUBLIC"]` when not specified.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_transit_gateway_attachment.test PATH
```

The above command imports Transit Gateway Attachment named `test` with the policy path `PATH`.
