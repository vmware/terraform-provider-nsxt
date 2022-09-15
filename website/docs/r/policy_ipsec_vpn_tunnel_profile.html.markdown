---
subcategory: "Policy - VPN"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ipsec_vpn_local_endpoint"
description: A resource to configure a IPSec VPN Local Endpoint.
---

# nsxt_policy_ipsec_vpn_local_endpoint

This resource provides a method for the management of IPSec VPN Local Endpoint.

This resource is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_ipsec_vpn_local_endpoint" "test" {
    display_name  = "test"
    service_path  = "Terraform provisioned IpsecVpnTunnelProfile"
    local_address = "10.203.1.14"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `local_address` - (Required) IPV4 Address of local endpoint.
* `local_id` - (Optional) Local identifier.
* `certificate_path` - (optional) Policy path referencing site certificate.
* `trust_ca_paths` - (Optional) List of policy paths referencing certificate authority to verify peer certificates.
* `trust_crl_paths` - (Optional)  List of policy paths referencing certificate revocation list to peer certificates.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ipsec_vpn_local_endpoint.test path
```

The above command imports endpoint named `test` with NSX path `path`.
