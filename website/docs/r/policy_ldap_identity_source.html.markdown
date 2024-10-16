---
subcategory: "User Management"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ldap_identity_source"
description: A resource to configure LDAP identity sources.
---

# nsxt_policy_ldap_identity_source

This resource provides a method for the management of LDAP identity sources.

## Example Usage

```hcl
resource "nsxt_policy_ldap_identity_source" "test" {
  nsx_id      = "Airius_LDAP"
  description = "Airius LDAP Identity Source"
  type        = "ActiveDirectory"
  domain_name = "airius.com"
  base_dn     = "DC=airius, DC=com"

  ldap_server {
    bind_identity = "nsxint@airius.com"
    password      = "forever big xray tempo"
    url           = "ldap://ldap-vip01.corp.airius.com"
    use_starttls  = true
    certificates = [
      trimspace(file("cert.pem")),
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `nsx_id` - (Required) The NSX ID of this resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `type` - (Required) Indicates the type of the LDAP identity source. Valid values are `ActiveDirectory`, `OpenLdap`.
* `domain_name` - (Required) Authentication domain name. This is the name of the authentication domain. When users log into NSX using an identity of the form "user@domain", NSX uses the domain portion to determine which LDAP identity source to use.
* `base_dn` - (Required) DN of subtree for user and group searches.
* `alternative_domain_names` - (Optional) Additional domains to be directed to this identity source. After parsing the "user@domain", the domain portion is used to select the LDAP identity source to use. Additional domains listed here will also be directed to this LDAP identity source. In Active Directory these are sometimes referred to as Alternative UPN Suffixes.
* `ldap_server` - (Required) List of LDAP servers that provide LDAP service for this identity source. Currently, only one LDAP server is supported.
    * `bind_identity` - (Optional) Username or DN for LDAP authentication.This user should have privileges to search the LDAP directory for groups and users. This user is also used in some cases (OpenLDAP) to look up an NSX user's distinguished name based on their NSX login name. If omitted, NSX will authenticate to the LDAP server using an LDAP anonymous bind operation. For Active Directory, provide a userPrincipalName (e.g. administrator@airius.com) or the full distinguished nane. For OpenLDAP, provide the distinguished name of the user (e.g. uid=admin, cn=airius, dc=com).
    * `certificates` - (Optional) TLS certificate(s) for LDAP server(s). If using LDAPS or STARTTLS, provide the X.509 certificate of the LDAP server in PEM format. This property is not required when connecting without TLS encryption and is ignored in that case.
    * `enabled` - (Optional) Allows the LDAP server to be enabled or disabled. When disabled, this LDAP server will not be used to authenticate users. Default value is `ture`.
    * `password` - (Optional) A password used when authenticating to the directory.
    * `url`- (Required) The URL for the LDAP server. Supported URL schemes are LDAP and LDAPS. Either a hostname or an IP address may be given, and the port number is optional and defaults to 389 for the LDAP scheme and 636 for the LDAPS scheme.
    * `use_starttls` - (Optional) If set to true, Use the StartTLS extended operation to upgrade the connection to TLS before sending any sensitive information. The LDAP server must support the StartTLS extended operation in order for this protocol to operate correctly. This option is ignored if the URL scheme is LDAPS.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ldap_identity_source.test ID
```
The above command imports LDAP identity source named `test` with the identifier `ID`.
