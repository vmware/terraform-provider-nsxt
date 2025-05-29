---
subcategory: "Firewall"
page_title: "NSXT: nsxt_policy_l7_access_profile"
description: A resource to configure a L7AccessProfile.
---

# nsxt_policy_l7_access_profile

This resource provides a method for the management of a L7AccessProfile.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_l7_access_profile" "test" {
  display_name          = "test_l7_access_profile"
  description           = "Terraform provisioned L7AccessProfile"
  default_action_logged = true
  default_action        = "ALLOW"

  l7_access_entry {
    display_name = "entry1"
    action       = "ALLOW"

    attribute {
      key              = "APP_ID"
      values           = ["SSL"]
      attribute_source = "SYSTEM"

      sub_attribute {
        key    = "TLS_CIPHER_SUITE"
        values = ["TLS_RSA_EXPORT_WITH_RC4_40_MD5"]
      }
    }

    disabled        = true
    logged          = "true"
    sequence_number = "100"
  }

  l7_access_entry {
    display_name = "entry2"
    action       = "REJECT_WITH_RESPONSE"

    attribute {
      key              = "URL_CATEGORY"
      values           = ["Abused Drugs"]
      attribute_source = "SYSTEM"
    }
    sequence_number = "200"
  }

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
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
* `default_action` - (Optional) The action to be applied to all the services. Valid values are `ALLOW`, `REJECT`, `REJECT_WITH_RESPONSE`.
* `default_action_logged` - (Optional) Flag to activate packet logging. Default is deactivated.
* `l7_access_entry` - (Optional) Property containing L7 access entries for Policy L7 Access Profile.
  * `nsx_id` - (Optional) The NSX ID of this L7 access entry. If set, this ID will be used to create the resource.
  * `display_name` - (Required) Display name of the L7 access entry.
  * `description` - (Optional) Description of the L7 access entry.
  * `action` - (Optional) The action to be applied to all the services.  Valid values are `ALLOW`, `REJECT`, `REJECT_WITH_RESPONSE`.
  * `attribute` - (Optional) Property containing attributes/sub-attributes for Policy L7 Access Profile.
    * `attribute_source` - (Optional) Source of attribute value i.e whether system defined or custom value. Valid values are `SYSTEM`, `CUSTOM`.
    * `custom_url_partial_match` - (Optional) True value for this flag will be treated as a partial match for custom url.
    * `description` - (Optional) Description for attribute value.
    * `key` - (Optional) Key for attribute. Supported Attribute Keys are APP_ID, URL_CATEGORY, CUSTOM_URL.
    * `isALGType` - (Optional) Describes whether the APP_ID value is ALG type or not.
    * `values` - (Optional) Multiple attribute values can be specified as elements of array.
    * `metadata` - (Optional) This is optional part that can hold additional data about the attribute key/values. Example - For URL CATEGORY key , it specified super category for url category value. This is generic array and can hold multiple meta information about key/values in future
      * `key` - (Optional) Key for metadata.
      * `value` - (Optional) Value for metadata key.
    * `sub_attribute` - (Optional) Reference to sub attributes for the attribute
      * `key` - (Optional) Key for sub attribute
      * `values` - (Optional) Multiple sub attribute values can be specified as elements of array.
  * `disabled` - (Optional) Flag to deactivate the entry. Default is activated.
  * `logged` - (Optional) Flag to activate packet logging. Default is deactivated.
  * `sequence_number` - (Optional) Determines the order of the entry in this profile. If no sequence number is specified in the payload, a value of 0 is assigned by default. If there are multiple rules with the same sequence number then their order is not deterministic.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `l7_access_entry` - (Optional) Property containing L7 access entries for Policy L7 Access Profile.
  * `path` - The NSX path of the L7 access entry.
  * `revision` - Indicates current revision number of the L7 access entry as seen by NSX API server. This attribute can be useful for debugging.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_l7_access_profile.test PATH
```

The above command imports L7AccessProfile named `test` with the NSX path `PATH`.
