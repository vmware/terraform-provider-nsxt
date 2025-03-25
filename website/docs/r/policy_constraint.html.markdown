---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_constraint"
description: A resource to configure a Constraint (Quota).
---

# nsxt_policy_constraint

This resource provides a method for the management of a Constraint.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_constraint" "test" {
    display_name = "demo-quota"
    description  = "Terraform provisioned Constraint"
    message      = "too many objects mate"

    target {
        path_prefix = "/orgs/default/projects/demo/"
    }

    instance_count {
        count                = 4
        target_resource_type = "StaticRoutes"
    }
    
    instance_count {
        count                = 1
        target_resource_type = "Infra.Tier1.PolicyDnsForwarder"
    }
    
    instance_count {
        count                = 20
        target_resource_type = "Infra.Domain.Group"
    }
}
```

## Example Usage - Multi-Tenancy

```hcl
resource "nsxt_policy_constraint" "test" {
  context {
      project_id = "demo"
  }

  display_name = "demo1-quota"

  target {
      path_prefix = "/orgs/default/projects/demo/vpcs/demo1/"
  }

  instance_count {
      count                = 4
      target_resource_type = "Org.Project.Vpc.PolicyNat.PolicyVpcNatRule"
  }
}
```

## Argument Reference

The following arguments are supported:

* `context` - (Optional) The context which the object belongs to
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `message` - (Optional) User friendly message to be shown to users upon violation.
* `target` - (Optional) Targets for the constraints to be enforced
  * `path_prefix` - (Optional) Prefix match to the path, needs to end with `\`
* `instance_count` - (Optional) Constraint details
  * `target_resource_type` - (Required) Type of the resource that should be limited to certain instance count
  * `operator` - (Optional) Either `<=` or `<`. Default is `<=`
  * `count` - (Required) Limit of instances
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
terraform import nsxt_policy_constraint.test PATH
```

The above command imports Constraint named `test` with the NSX path `PATH`.
