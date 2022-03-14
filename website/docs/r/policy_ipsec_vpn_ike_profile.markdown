---
subcategory: "Policy - Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ipsec_vpn_ike_profile"
description: A resource to configure a IPSec VPN Ike Profile.
---

# nsxt_policy_ipsec_vpn_ike_profile

This resource provides a method for the management of a IPSec VPN Ike Profile.

This resource is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_ipsec_vpn_ike_profile" "test" {
    display_name      = "test"
    description       = "Terraform provisioned IPSec VPN Ike Profile"
    encryption_algorithms = ["AES_128"]
    digest_algorithms     = ["SHA2_256"]
    dh_groups             = ["GROUP14"]
    ike_version           = "IKE_V2"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `ike_version` - (Optional) IKE protocol version to be used. IKE-Flex will initiate IKE-V2 and responds to both IKE-V1 and IKE-V2.
* `encryption_algorithms` - (Optional) Encryption algorithm is used during Internet Key Exchange(IKE) negotiation. Default is AES_128.
* `digest_algorithms` - (Optional) Algorithm to be used for message digest during Internet Key Exchange(IKE) negotiation. Default is SHA2_256.
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
terraform import nsxt_policy_ipsec_vpn_ike_profile.test UUID
```

The above command imports IPSec VPN Ike Profile named `test` with the NSX IPSec VPN Ike Profile ID `UUID`.
