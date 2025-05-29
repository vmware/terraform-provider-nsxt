---
subcategory: "Load Balancer"
page_title: "NSXT: nsxt_policy_lb_server_ssl_profile"
description: A resource to configure a Load Balancer Server SSL Profile.
---

# nsxt_policy_lb_server_ssl_profile

This resource provides a method for the management of a Load Balancer Server SSL Profile.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_lb_server_ssl_profile" "test" {
  display_name          = "test"
  description           = "Terraform provisioned profile"
  cipher_group_label    = "BALANCED"
  session_cache_enabled = true
}
```

## Custom Cipher Profile Example Usage

```hcl
resource "nsxt_policy_lb_server_ssl_profile" "test" {
  display_name          = "test"
  description           = "Terraform provisioned profile"
  cipher_group_label    = "CUSTOM"
  ciphers               = ["TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"]
  session_cache_enabled = false
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `cipher_group_label` - (Required) One of `BALANCED`, `HIGH_SECURITY`, `HIGH_COMPATIBILITY` or `CUSTOM`.
* `ciphers` - (Optional) Supported SSL cipher list. Can only be specified if `cipher_group_label` is set to `CUSTOM`, otherwise the list is assigned by NSX.
* `protocols` - (Optional) Since NSX 4.2, only `TLS_V1_2` is supported.
* `session_cache_enabled` - (Optional) SSL session caching allows SSL client and server to reuse previously
negotiated security parameters avoiding the expensive public key operation during handshake.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `is_secure` - (Optional) This flag is set to true when all the ciphers and protocols are secure.
It is set to false when one of the ciphers or protocols is insecure.
* `is_fips` - (Optional) This flag is set to true when all the ciphers and protocols are FIPS
compliant.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_lb_server_ssl_profile.test PATH
```

The above command imports Server SSL profile named `test` with the NSX path `PATH`.
