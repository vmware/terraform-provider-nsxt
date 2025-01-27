---
subcategory: "Firewall"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_l7_access_profile"
description: A resource to configure a L7AccessProfile.
---

# nsxt_policy_l7_access_profile

This resource provides a method for the management of a L7AccessProfile.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_l7_access_profile" "test" {
    display_name      = "test"
    description       = "Terraform provisioned L7AccessProfile"
    parent_path = FILL PARENT RESOURCE.path
default_action_logged = FILL VALUE FOR schema.TypeBool
default_action = "ALLOW"
entry_count = FILL VALUE FOR schema.TypeInt

}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `parent_path` - (Required) Path of parent object
* `default_action_logged` - (Optional) Flag to activate packet logging. Default is deactivated.
* `default_action` - (Optional) The action to be applied to all the services.
* `entry_count` - (Optional) The count of entries in the L7 profile.

* `l7_access_entry` - (Optional) Property containing L7 access entries for Policy L7 Access Profile.

  * `action` - (Optional) The action to be applied to all the services.
  * `attribute` - (Optional) Property containing attributes/sub-attributes for Policy L7 Access Profile.
    * `sub_attribute` - (Optional) Reference to sub attributes for the attribute
      * `datatype` - (Optional) Datatype for sub attribute
      * `value` - (Optional) Multiple sub attribute values can be specified as elements of array.

      * `key` - (Optional) Key for sub attribute
    * `attribute_source` - (Optional) Source of attribute value i.e whether system defined or custom value
    * `custom_url_partial_match` - (Optional) True value for this flag will be treated as a partial match for custom url
    * `description` - (Optional) Description for attribute value
    * `key` - (Optional) URL_Reputation is currently not available. Please do not use it in Attribute Key while creating context profile
    * `datatype` - (Optional) Datatype for attribute
    * `isALGType` - (Optional) Describes whether the APP_ID value is ALG type or not.
    * `value` - (Optional) Multiple attribute values can be specified as elements of array.

    * `metadata` - (Optional) This is optional part that can hold additional data about the attribute key/values.
Example - For URL CATEGORY key , it specified super category for url category value.
This is generic array and can hold multiple meta information about key/values in future

      * `value` - (Optional) Value for metadata key
      * `key` - (Optional) Key for metadata
  * `disabled` - (Optional) Flag to deactivate the entry. Default is activated.
  * `logged` - (Optional) Flag to activate packet logging. Default is deactivated.
  * `sequence_number` - (Optional) Determines the order of the entry in this profile. If no sequence number is
specified in the payload, a value of 0 is assigned by default. If there are
multiple rules with the same sequence number then their order is not deterministic.



## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_l7_access_profile.test PATH
```

The above command imports L7AccessProfile named `test` with the NSX path `PATH`.
