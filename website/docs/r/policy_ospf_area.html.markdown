---
subcategory: "OSPF"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ospf_area"
description: A resource to configure a OSPF Area.
---

# nsxt_policy_ospf_area

This resource provides a method for the management of OSPF Area. Only a single area is supported so far per Tier-0 Gateway OSPF Config.

This resource is applicable to NSX Policy Manager only.
This resource is supported with NSX 3.1.1 onwards.

## Example Usage

```hcl
resource "nsxt_policy_ospf_area" "test" {
  display_name = "test"
  description  = "Terraform provisioned OSPF Area"
  ospf_path    = nsxt_policy_ospf_config.test.path
  area_id      = "15"
  area_type    = "NORMAL"

  auth_mode  = "PASSWORD"
  secret_key = "af12a1ed"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `ospf_path` - (Required) The policy path to the OSPF configuration on particular Tier-0 Gateway.
* `area_id   - (Required) OSPF Area ID in either decimal or dotted format.
* `area_type` - (Optional) OSPF Area type, one of `NORMAL` or `NSSA`. Default is `NSSA`.
* `auth_mode` - (Optional) OSPF Authentication mode, one of `NONE`, `PASSWORD` or `MD5`. By default, OSPF authentication is disabled with mode `NONE`.
* `key_id` - (Optional) Authentication secret key id, required for authenication mode `MD5`. This attribute is sensitive.
* `secret_key` - (Optional) Authentication secret key, required for authentication mode other than `NONE`. This attribute is sensitive. Length should not exceed 8 characters.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing OSPF Area can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ospf_area.test GW-ID/LOCALE-SERVICE-ID/ID
```

The above command imports OSPF Area named `test` with NSX ID `ID` on Tier-0 Gateway `GW-ID` and Locale Service `LOCALE-SERVICE-ID`.
