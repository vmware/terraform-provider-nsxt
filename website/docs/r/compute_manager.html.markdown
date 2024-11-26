---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_compute_manager"
description: A resource to configure a Compute Manager.
---

# nsxt_compute_manager

This resource provides a method for the management of a Compute Manager.
This resource is supported with NSX 4.1.0 onwards.

## Example Usage

```hcl
resource "nsxt_compute_manager" "test" {
  description  = "Terraform provisioned Compute Manager"
  display_name = "test"
  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  server = "192.168.244.144"

  credential {
    username_password_login {
      username = "user"
      password = "pass"
    }
  }
  origin_type = "vCenter"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `access_level_for_oidc` - (Optional) Specifies access level to NSX from the compute manager. Accepted values - 'FULL' or 'LIMITED'. The default value is 'FULL'.
* `create_service_account` - (Optional) Specifies whether service account is created or not on compute manager. Note that `false` value for this attribute is no longer supported with NSX 9.0 and above.
* `credential` - (Required) Login credentials for the compute manager. Should contain exactly one credential enlisted below: 
  * `saml_login` - (Optional) A login credential specifying saml token.
    * `thumbprint` - (Required) Thumbprint of the server.
    * `token` - (Required) The saml token to login to server.
  * `session_login` - (Optional) A login credential specifying session_id.
    * `session_id` - (Required) The session_id to login to server.
    * `thumbprint` - (Required) Thumbprint of the login server.
  * `username_password_login` - (Optional) A login credential specifying a username and password.
    * `password` - (Required) The authentication password for login.
    * `thumbprint` - (Required) Thumbprint of the login server.
    * `username` - (Required) The username for login.
  * `verifiable_asymmetric_login` - (Optional) A verifiable asymmetric login credential.
    * `asymmetric_credential` - (Required) Asymmetric login credential.
    * `credential_key` - (Required) Credential key.
    * `credential_verifier` - (Required) Credential verifier.
  * `extension_certificate` - (Optional) Specifies certificate for compute manager extension.
    * `pem_encoded` - (Required) PEM encoded certificate data.
    * `private_key` - (Required) Private key of certificate.
  * `multi_nsx` - (Optional) Specifies whether multi nsx feature is enabled for compute manager.
  * `origin_type` - (Optional) Compute manager type like vCenter. Defaults to vCenter if not set.
  * `reverse_proxy_https_port` - (Optional) Proxy https port of compute manager.
  * `server` - (Required) IP address or hostname of compute manager.
  * `set_as_oidc_provider` - (Optional) Specifies whether compute manager has been set as OIDC provider.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `origin_properties` - Key-Value map of additional specific properties of compute manager.
  * `key` - Key.
  * `value` - Value.

## Importing

An existing Compute Manager can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_compute_manager.test UUID
```
The above command imports Compute Manager named `test` with the NSX Compute Manager ID `UUID`.
