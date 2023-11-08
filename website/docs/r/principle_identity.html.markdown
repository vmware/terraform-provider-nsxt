---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_principle_identity"
description: A resource to configure principle identities.
---

# nsxt_principle_identity

This resource provides a method for the management of Principle Identities.

## Example Usage

```hcl
resource "nsxt_principle_identity" "test" {
  name            = "open-stack"
  node_id         = "node-2"
  certificate_pem = trimspace(file("cert.pem"))
  roles_for_path {
    path = "/orgs/default"

    role {
      role = "org_admin"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `is_protected` - (optional) Indicates whether the entities created by this principal should be protected.
* `name` - (Required) Name of the principal.
* `node_id` - (Required) Unique node-id of a principal. This is used primarily in the case where a cluster of nodes is used to make calls to the NSX Manager and the same `name` is used so that the nodes can access and modify the same data while still accessing NSX through their individual secret (certificate or JWT). In all other cases this can be any string.
* `certificate_pem` - (Required) PEM encoding of the certificate to be associated with this principle identity.
* `roles_for_path` - (Required) A list of The roles that are associated with the user, limiting them to a path. In case the path is '/', the roles apply everywhere.
    * `path` - (Required) Path of the entity in parent hierarchy.
    * `role` - (Required) A list of identifiers for the roles to associate with the given user limited to a path.

Once a Principle Identity is created, it can't be modified. Modification of above arguments will cause the current PI on NSX to be deleted and recreated. Certificate updates is also handled in the same way. 

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `certificate_id` - NSX certificate ID of the imported `certificate_pem`.


## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_principle_identity.test PRINCIPLE_IDENTITY_ID
```
The above command imports Principle Identity named `test` with the identifier `PRINCIPLE_IDENTITY_ID`.
