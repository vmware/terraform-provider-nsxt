---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_lb_client_ssl_profile"
description: A resource to configure a LB Client SSL Profile.
---

# nsxt_policy_lb_client_ssl_profile

This resource provides a method for the management of a LBClientSslProfile.

This resource is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_lb_client_ssl_profile" "test" {
  display_name          = "test"
  description           = "Terraform provisioned LBClientSslProfile"
  cipher_group_label    = "HIGH_COMPATIBILITY"
  ciphers               = ["TLS_ECDH_ECDSA_WITH_AES_128_CBC_SHA"]
  is_fips               = true
  is_secure             = true
  prefer_server_ciphers = true
  protocols             = ["TLS_V1_2"]
  session_cache_enabled = true
  session_cache_timeout = 2
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `cipher_group_label` - (Optional) A label of cipher group which is mostly consumed by GUI. Possible values are: `BALANCED`, `HIGH_SECURITY`, `HIGH_COMPATIBILITY` and `CUSTOM`.
* `ciphers` - (Optional) Supported SSL cipher list to client side. Possible values are: `TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`, `TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384`, `TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA`,`TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA`, `TLS_ECDH_ECDSA_WITH_AES_256_CBC_SHA`, `TLS_ECDH_RSA_WITH_AES_256_CBC_SHA `, `TLS_RSA_WITH_AES_256_CBC_SHA`, `TLS_RSA_WITH_AES_128_CBC_SHA`, `TLS_RSA_WITH_3DES_EDE_CBC_SHA`,`TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA`, `TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256`, `TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384`, `TLS_RSA_WITH_AES_128_CBC_SHA256`,`TLS_RSA_WITH_AES_128_GCM_SHA256`, `TLS_RSA_WITH_AES_256_CBC_SHA256`,`TLS_RSA_WITH_AES_256_GCM_SHA384`, `TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA`, `TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256`, `TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256`,  `TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384`, `TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384`, `TLS_ECDH_ECDSA_WITH_AES_128_CBC_SHA`, `TLS_ECDH_ECDSA_WITH_AES_128_CBC_SHA256`, `TLS_ECDH_ECDSA_WITH_AES_128_GCM_SHA256`,`TLS_ECDH_ECDSA_WITH_AES_256_CBC_SHA384`, `TLS_ECDH_ECDSA_WITH_AES_256_GCM_SHA384`, `TLS_ECDH_RSA_WITH_AES_128_CBC_SHA`, `TLS_ECDH_RSA_WITH_AES_128_CBC_SHA256`, `TLS_ECDH_RSA_WITH_AES_128_GCM_SHA256`,`TLS_ECDH_RSA_WITH_AES_256_CBC_SHA384`, `TLS_ECDH_RSA_WITH_AES_256_GCM_SHA384`
* `is_fips` - (Optional) This flag is set to true when all the ciphers and protocols are FIPS compliant. It is set to false when one of the ciphers or protocols are not FIPS compliant. Read-only property, its value will be decided automatically based on the result of applying this configuration.
* `is_secure` - (Optional) This flag is set to true when all the ciphers and protocols are secure. It is set to false when one of the ciphers or protocols is insecure.  Read-only property, its value will be decided automatically based on the result of applying this configuration.
* `prefer_server_ciphers` - (Optional) During SSL handshake as part of the SSL client Hello client sends an ordered list of ciphers that it can support (or prefers) and typically server selects the first one from the top of that list it can also support. For Perfect Forward Secrecy(PFS), server could override the client's preference. Default is `true`.
* `protocols` - (Optional) Protocols used by the LB Client SSL profile. SSL versions TLS1.1 and TLS1.2 are supported and enabled by default. SSLv2, SSLv3, and TLS1.0 are supported, but disabled by default. Possible values are:`SSL_V2`, `SSL_V3`, `TLS_V1`, `TLS_V1_1`, `TLS_V1_2`, SSL versions TLS1.1 and TLS1.2 are supported and enabled by default. SSLv2, SSLv3, and TLS1.0 are supported, but disabled by default.
* `session_cache_enabled` - (Optional) SSL session caching allows SSL client and server to reuse previously negotiated security parameters avoiding the expensive public key operation during handshake. Default is `true`.
* `session_cache_timeout` - (Optional) Session cache timeout specifies how long the SSL session parameters are held on to and can be reused. format: int64, default is `300`.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_lb_client_ssl_profile.test UUID
```

The above command imports LBClientSslProfile named `test` with the NSX ID `UUID`.
