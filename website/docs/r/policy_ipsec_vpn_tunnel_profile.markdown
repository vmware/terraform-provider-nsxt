---
subcategory: "Policy - Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ipsec_vpn_tunnel_profile"
description: A resource to configure a IPSec VPN Tunnel Profile.
---

# nsxt_policy_ipsec_vpn_tunnel_profile

This resource provides a method for the management of a IPSec VPN Ike Profile.

This resource is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_ipsec_vpn_tunnel_profile" "test" {
    display_name      = "test"
    description       = "Terraform-provisioned IPsec VPN Tunnel Profile-name"
    encryption_algorithms = ["AES_GCM_128"]
    digest_algorithms     = []
    dh_groups             = ["GROUP14"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `encryption_algorithms` - (Optional) Encryption algorithm to encrypt/decrypt the messages exchanged between IPSec VPN initiator and responder during tunnel negotiation. Default is AES_GCM_128.
* `digest_algorithms` - (Optional) Algorithm to be used for message digest. Default digest algorithm is implicitly covered by default encryption algorithm "AES_GCM_128".
* `dh_groups` - (Optional) Diffie-Hellman group to be used if PFS is enabled. Default is GROUP14.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_ipsec_vpn_tunnel_profile.test UUID
```

The above command imports IPSec VPN Ike Profile named `test` with the NSX IPSec VPN Ike Profile ID `UUID`.
