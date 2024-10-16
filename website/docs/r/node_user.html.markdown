---
subcategory: "User Management"
layout: "nsxt"
page_title: "NSXT: nsxt_node_user"
description: A resource to configure Node Users.
---

# nsxt_node_user

This resource provides a method for the management of NSX node users. These node user accounts can log in to the NSX web-based user interface or access API.

~> **NOTE:** This resource requires NSX version 4.1.0 or higher.

## Example Usage

```hcl
resource "nsxt_node_user" "test_user" {
  active                    = true
  full_name                 = "John Doe"
  password                  = "Str0ng_Pwd!Wins$"
  username                  = "johndoe123"
  password_change_frequency = 180
  password_change_warning   = 30
}
```

## Argument Reference

* `active` - (Optional) If this account should be activated or deactivated. Default value is `true`.
* `full_name` - (Required) The full name of this user.
* `password` - (Optional) Password of this user. Password must be specified for creating `ACTIVE` accounts. Password updates will only take effect after deactivating the account, then reactivating it.
* `username` - (Required) User login name.
* `password_change_frequency` - (Optional) Number of days password is valid before it must be changed. This can be set to 0 to indicate no password change is required or a positive integer up to 9999. By default local user passwords must be changed every 90 days.
* `password_change_warning` - (Optional) Number of days before user receives warning message of password expiration.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `last_password_change` - Number of days since password was last changed.
* `password_reset_required` - Boolean value that states if a password reset is required.
* `status` - Status of the user. This value can be `ACTIVE` indicating authentication attempts will be successful if the correct credentials are specified. The value can also be `PASSWORD_EXPIRED` indicating authentication attempts will fail because the user's password has expired and must be changed. Or, this value can be `NOT_ACTIVATED` indicating the user's password has not yet been set and must be set before the user can authenticate.
* `user_id` - Numeric id for the user.
## Importing

An existing Node User can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import
```
terraform import nsxt_node_user.user1 USER_ID
```
The above command imports the User `user1` with the Numeric id (`USER_ID`) of the user.
