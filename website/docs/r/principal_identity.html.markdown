---
subcategory: "User Management"
layout: "nsxt"
page_title: "NSXT: nsxt_principal_identity"
description: A resource to configure principal identities.
---

# nsxt_principal_identity

This resource provides a method for the management of Principal Identities.

~> **NOTE:** This resource requires NSX version 4.0.0 or higher.

## Example Usage

```hcl
resource "nsxt_principal_identity" "test" {
  name            = "open-stack"
  node_id         = "node-2"
  certificate_pem = trimspace(file("cert.pem"))
  roles_for_path {
    path  = "/orgs/default"
    roles = ["org_admin"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `is_protected` - (optional) Indicates whether the entities created by this principal should be protected.
* `name` - (Required) Name of the principal.
* `node_id` - (Required) Unique node-id of a principal. This is used primarily in the case where a cluster of nodes is used to make calls to the NSX Manager and the same `name` is used so that the nodes can access and modify the same data while still accessing NSX through their individual secret (certificate or JWT). In all other cases this can be any string.
* `certificate_pem` - (Required) PEM encoding of the certificate to be associated with this principal identity.
* `roles_for_path` - (Required) A list of The roles that are associated with the user, limiting them to a path. In case the path is '/', the roles apply everywhere.
    * `path` - (Required) Path of the entity in parent hierarchy.
    * `roles` - (Required) A list of identifiers for the roles to associate with the given user limited to a path.

Once a Principal Identity is created, it can't be modified. Modification of above arguments will cause the current PI on NSX to be deleted and recreated. Certificate updates is also handled in the same way. 

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `certificate_id` - NSX certificate ID of the imported `certificate_pem`.

# Build-in NSX roles

`roles_for_path.roles` accepts user created roles as well as roles native to NSX. For reference, the following is a list of native roles as of NSX 4.1.2
- `network_engineer`: Network Admin
- `support_bundle_collector`: Support Bundle Collector
- `security_op`: Security Operator
- `lb_auditor`: LB Operator
- `netx_partner_admin`: NETX Partner Admin
- `project_admin`: Project Admin
- `auditor`: Auditor
- `network_op`: Network Operator
- `enterprise_admin`: Enterprise Admin
- `lb_admin`: LB Admin
- `gi_partner_admin`: GI Partner Admin
- `vpn_admin`: VPN Admin
- `vpc_admin`: VPC Admin
- `security_engineer`: Security Admin

The permission matrix for above roles is available on [NSX documentation](https://docs.vmware.com/en/VMware-NSX/4.1/administration/GUID-26C44DE8-1854-4B06-B6DA-A2FD426CDF44.html)

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_principal_identity.test PRINCIPAL_IDENTITY_ID
```
The above command imports Principal Identity named `test` with the identifier `PRINCIPAL_IDENTITY_ID`.
