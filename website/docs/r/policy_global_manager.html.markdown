---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_global_manager"
description: A resource to configure Policy Global Manager.
---

# nsxt_policy_global_manager

This resource provides a method for the management of Policy Global Manager.

This resource is applicable to NSX Global Manager.

## Example Usage

Activate Global Manager as ACTIVE
```hcl
resource "nsxt_policy_global_manager" "test" {
  display_name = "test"
  mode         = "ACTIVE"
}
```

Attach a backup Global Manager as STANDBY
```hcl
resource "nsxt_policy_global_manager" "test" {
  display_name = "test"
  description  = "Terraform provisioned Global Manager"
  connection_info {
    fqdn       = "192.168.230.230"
    username   = "admin"
    password   = "somepasswd"
    thumbprint = "207d65dcb6f17aa5a1ef2365ee6ae0b396867baa92464e5f8a46f6853708b9ef"
  }
  mode = "STANDBY"
}
```
## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `fail_if_rtt_exceeded` - (Optional) Fail onboarding if maximum RTT exceeded. Default is true.
* `maximum_rtt` - (Optional) Maximum acceptable packet round trip time (RTT). Default is 250.
* `connection_info` - (Optional) Connection information.
    * `fqdn` - (Optional) Fully Qualified Domain Name of the Management Node.
    * `password` - (Optional) Password.
    * `site_uuid` - (Optional) ID of Global Manager.
    * `thumbprint` - (Optional) Thumbprint of Enforcement Point.
    * `username` - (Optional) Username.
* `mode` - (Required) Mode of the global manager. Allowed values are `ACTIVE`, `STANDBY`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_global_manager.test POLICY_PATH
```
The above command imports Global Manager named `test` with policy path `POLICY_PATH`.
