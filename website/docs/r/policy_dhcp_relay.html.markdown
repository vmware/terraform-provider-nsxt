---
subcategory: "DHCP"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_dhcp_relay"
description: A resource to configure a Dhcp Relay.
---

# nsxt_policy_dhcp_relay

This resource provides a method for the management of a Dhcp Relay.

This resource is applicable to NSX Policy Manager and NSX Global Manager.
 
## Example Usage

```hcl
resource "nsxt_policy_dhcp_relay" "test" {
  display_name     = "test"
  description      = "Terraform provisioned Dhcp Relay"
  server_addresses = ["10.0.0.2", "7001::2"]
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_dhcp_relay" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name     = "test"
  description      = "Terraform provisioned Dhcp Relay"
  server_addresses = ["10.0.0.2", "7001::2"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `server_addresses` - (Required) List of DHCP server addresses.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_dhcp_relay.test ID
```
The above command imports Dhcp Relay named `test` with the NSX Dhcp Relay ID `ID`.

```
terraform import nsxt_policy_dhcp_relay.test POLICY_PATH
```
The above command imports Dhcp Relay named `test` with the NSX Dhcp Relay policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
