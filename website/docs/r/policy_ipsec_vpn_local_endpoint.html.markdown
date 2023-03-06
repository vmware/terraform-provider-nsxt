---
subcategory: "VPN"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ipsec_vpn_local_endpoint"
description: A resource to configure a IPSec VPN Local Endpoint.
---
# nsxt_policy_ipsec_vpn_local_endpoint

This resource provides a method for the management of a IPSec VPN Local Endpoint.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_ipsec_vpn_local_endpoint" "test" {
  display_name     = "test"
  service_path     = nsxt_policy_ipsec_vpn_service.test.path
  description      = "Terraform provisioned IPSec VPN Local Endpoint"
  local_address    = "20.20.0.10"
  local_id         = "test"
  certificate_path = data.nsxt_policy_certificate.cert.path
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `local_address` - (Required) Local IPv4 IP address.
* `local_id` - (Optional) Local id for the local endpoint.
* `certificate_path` - (Optional) Policy path referencing site certificate.
* `trust_ca_paths` - (Optional) List of trust ca certificate paths.
* `trust_crl_paths` - (Optional) List of trust CRL paths.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.


## Importing

An existing segment can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ipsec_vpn_local_endpoint.enpoint1 ID
```

The above command imports the endpoint named `endpoint1` with the NSX ID `ID`.
