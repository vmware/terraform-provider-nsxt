---
subcategory: "VPN"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ipsec_vpn_tunnel_profile"
description: A resource to configure a IPSec VPN Tunnel Profile.
---

# nsxt_policy_ipsec_vpn_tunnel_profile

This resource provides a method for the management of a IPSec VPN Tunnel Profile.

This resource is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_ipsec_vpn_tunnel_profile" "test" {
  display_name          = "test"
  description           = "Terraform provisioned IPSec VPN Ike Profile"
  df_policy             = "COPY"
  encryption_algorithms = ["AES_128"]
  digest_algorithms     = ["SHA2_256"]
  dh_groups             = ["GROUP14"]
  sa_life_time          = 7200
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `digest_algorithms` - (Required) Set of algorithms to be used for message digest during IKE negotiation. Default is `SHA2_256`.
* `dh_groups` - (Required) Diffie-Hellman group to be used if PFS is enabled. Default is GROUP14.
* `df_policy` - (Optional) Defragmentation policy, one of `COPY` or `CLEAR`. `COPY` copies the defragmentation bit from the inner IP packet into the outer packet. `CLEAR` ignores the defragmentation bit present in the inner packet. Default is `COPY`.
* `encryption_algorithms` - (Optional) Set of encryption algorithms to be used during IKE negotiation.
* `sa_life_time` - (Optional) SA lifetime specifies the expiry time of security association. Default is 3600.
* `enable_perfect_forward_secrecy` - (Optional) Enable perfect forward secrecy. Default is True.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ipsec_vpn_tunnel_profile.test UUID
```

The above command imports IPSec VPN IKE Profile named `test` with the NSX ID `UUID`.
