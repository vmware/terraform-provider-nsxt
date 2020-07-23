---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_context_profile"
sidebar_current: "docs-nsxt-resource-policy-context-profile"
description: A resource to configure a Context Profile.
---

# nsxt_policy_context_profile

This resource provides a method for the management of a Context Profile.
 
## Example Usage

```hcl
resource "nsxt_policy_context_profile" "test" {
    display_name      = "test"
    description       = "Terraform provisioned ContextProfile"
    attribute {
        data_type   = "STRING"
        description = "test-domain-name-attribute"
        key         = "DOMAIN_NAME"
        value       = ["*-myfiles.sharepoint.com"]
    }
    attribute {
        data_type   = "STRING"
        description = "test-app-id-attribute"
        key         = "APP_ID"
        value       = ["SSL"]
        sub_attribute {
            data_type = "STRING"
            key       = "TLS_VERSION"
            value     = ["SSL_V3"]
        }
    }
}

```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `attribute` - (Required) A repeatable block to specify attributes for the context profile. At least one block is required.
  * `datatype` - (Required) Datatype for `attribute`, must be `STRING`.
  * `description` - (Optional) Description of the attribute.
  * `key` - (Required) A string value indicating key for the `attribute`. Must be one of `APP_ID`, `DOMAIN_NAME`, or `URL_CATEGORY`.
  * `value` - (Required) A list of string indicating values for the `attribute`. Must be a subset of the preset list of valid values for the `key` on NSX.
  * `sub_attribute` - (Optional) A repeatable block to specify sub attributes for the attribute. This configuration is only valid when `value` has only one element, and `sub_attribute` is supported for that value on NSX.
    * `datatype` - (Required for `sub_attribute`) Datatype for `sub_attribute`, must be `STRING`.
    * `key` - (Required for `sub_attribute`) A string value indicating key for the `sub_attribute`. Must be one of `TLS_CIPHER_SUITE`, `TLS_VERSION`, or `CIFS_SMB_VERSION`.
    * `value` - (Required for `sub_attribute`) A list of string indicating values for the `sub_attribute`. Must be a subset of the preset list of valid values for the `key` on NSX.
 
## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `attribute`:
  * `is_alg_type` - Describes whether the APP_ID value is ALG type or not.

## Importing

An existing Context Profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_context_profile.test UUID
```

The above command imports Context Profile named `test` with the NSX Context Profile ID `UUID`.
