---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_lb_server_ssl_profile"
description: A resource to configure a VPC Load Balancer Server SSL Profile.
---

# nsxt_vpc_lb_server_ssl_profile

This resource provides a method for the management of a Load Balancer Server SSL Profile.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.0.2 onwards.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_vpc" "demovpc" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "vpc1"
}

resource "nsxt_vpc_lb_server_ssl_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }
  display_name          = "test"
  description           = "Terraform provisioned profile"
  cipher_group_label    = "BALANCED"
  session_cache_enabled = true
}
```

## Custom Cipher Profile Example Usage

```hcl
resource "nsxt_vpc_lb_server_ssl_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }
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
* `context` - (Required) The context which the object belongs to. Requires a VPC context.
    * `project_id` - (Required) The ID of the project which the object belongs to.
    * `vpc_id` - (Required) The ID of the VPC which the object belongs to.
* `cipher_group_label` - (Required) One of `BALANCED`, `HIGH_SECURITY`, `HIGH_COMPATIBILITY` or `CUSTOM`.
* `ciphers` - (Optional) Supported SSL cipher list. Can only be specified if `cipher_group_label` is set to `CUSTOM`, otherwise the list is assigned by NSX.
* `protocols` - (Optional) Since NSX 4.2, only `TLS_V1_2` is supported.
* `session_cache_enabled` - (Optional) SSL session caching allows SSL client and server to reuse previously
negotiated security parameters avoiding the expensive public key operation during handshake.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `is_secure` - (Optional) This flag is set to true when all the ciphers and protocols are secure.
It is set to false when one of the ciphers or protocols is insecure.
* `is_fips` - (Optional) This flag is set to true when all the ciphers and protocols are FIPS
compliant.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc_lb_server_ssl_profile.test PATH
```

The above command imports a VPC LB Server SSL Profile named `test` with the policy path `PATH`.
The path must be a VPC path of the form `/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lb-server-ssl-profiles/[id]`.
Importing by ID only is not supported, since the project and VPC are derived from the path.
