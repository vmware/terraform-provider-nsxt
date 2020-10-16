---
subcategory: "Policy - Firewall"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_context_profile"
description: A resource to configure a Context Profile.
---

# nsxt_policy_context_profile

This resource provides a method for the management of a Context Profile.
This resource is supported with NSX 3.0.0 onwards.

## Example Usage

```hcl
resource "nsxt_policy_context_profile" "test" {
  display_name = "test"
  description  = "Terraform provisioned ContextProfile"
  domain_name {
    description = "test-domain-name-attribute"
    value       = ["*-myfiles.sharepoint.com"]
  }
  app_id {
    description = "test-app-id-attribute"
    value       = ["SSL"]
    sub_attribute {
      tls_version = ["SSL_V3"]
    }
  }
}

```

## Argument Reference

The following arguments are supported:
Note: At least one of `app_id`, `domain_name`, or `url_category` must present.

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `app_id` - (Optional) A block to specify app id attributes for the context profile. Only one block is allowed.
  * `description` - (Optional) Description of the attribute.
  * `value` - (Required) A list of string indicating values for the `app_id`. Must be a subset of valid values for `app_id` on NSX.
  * `sub_attribute` - (Optional) A block to specify sub attribute for the `app_id`. Only one block is allowed.
    * `tls_cipher_suite` - (Optional) A list of string indicating values for `tls_cipher_suite`, only applicable to `SSL`.
    * `tls_version` - (Optional) A list of string indicating values for `tls_version`, only applicable to `SSL`.
    * `cifs_smb_version` - (Optional) A list of string indicating values for `cifs_smb_version`, only applicable to `CIFS`.
* `domain_name` - (Optional) A block to specify domain name (FQDN) attributes for the context profile. Only one block is allowed.
  * `description` - (Optional) Description of the attribute.
  * `value` - (Required) A list of string indicating values for the `domain_name`. Must be a subset of valid values for `domain_name` on NSX.
* `url_category` - (Optional) A block to specify url category attributes for the context profile. Only one block is allowed.
  * `description` - (Optional) Description of the attribute.
  * `value` - (Required) A list of string indicating values for the `url_category`. Must be a subset of valid values for `url_category` on NSX.
 
## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `app_id`:
  * `is_alg_type` - Describes whether the APP_ID value is ALG type or not.

## Importing

An existing Context Profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_context_profile.test UUID
```

The above command imports Context Profile named `test` with the NSX Context Profile ID `UUID`.
