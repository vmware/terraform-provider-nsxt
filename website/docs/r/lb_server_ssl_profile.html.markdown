---
layout: "nsxt"
page_title: "NSXT: nsxt_lb_server_ssl_profile"
sidebar_current: "docs-nsxt-resource-lb-server-ssl-profile"
description: |-
  Provides a resource to configure lb server ssl profile on NSX-T manager
---

# nsxt_lb_server_ssl_profile

Provides a resource to configure lb server ssl profile on NSX-T manager

## Example Usage

```hcl
resource "nsxt_lb_server_ssl_profile" "lb_server_ssl_profile" {
  description           = "lb_server_ssl_profile provisioned by Terraform"
  display_name          = "lb_server_ssl_profile"
  protocols             = ["TLS_V1_2"]
  ciphers               = ["TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256", "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384"]
  session_cache_enabled = true

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb server ssl profile.
* `ciphers` - (Optional) supported SSL cipher list to server side.
* `protocols` - (Optional) SSL versions TLS1.1 and TLS1.2 are supported and enabled by default. SSLv2, SSLv3, and TLS1.0 are supported, but disabled by default.
* `session_cache_enabled` - (Optional) SSL session caching allows SSL server and server to reuse previously negotiated security parameters avoiding the expensive public key operation during handshake. Defaults to true.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb server ssl profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `is_secure` - This flag is set to true when all the ciphers and protocols are secure. It is set to false when one of the ciphers or protocols is insecure.


## Importing

An existing lb server ssl profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_lb_server_ssl_profile.lb_server_ssl_profile UUID
```

The above would import the lb server ssl profile named `lb_server_ssl_profile` with the nsx id `UUID`
