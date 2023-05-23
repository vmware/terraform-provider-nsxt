---
subcategory: "Load Balancer"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_lb_service"
description: A resource to configure a Load Balancer Service.
---

# nsxt_policy_lb_service

This resource provides a method for the management of a Load Balancer Service.

This resource is applicable to NSX Policy Manager.

In order to enforce correct order of create/delete, it is recommended to add
depends_on clause to lb service.

## Example Usage

```hcl
data "nsxt_policy_tier1_gateway" "test" {
  display_name = "test"
}


resource "nsxt_policy_lb_service" "test" {
  display_name      = "test"
  description       = "Terraform provisioned Service"
  connectivity_path = data.nsxt_policy_tier1_gateway.test.path
  size              = "SMALL"
  enabled           = true
  error_log_level   = "ERROR"
  depends_on        = [nsxt_policy_tier1_gateway_interface.tier1_gateway_interface]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `size` - (Optional) Load Balancer Service size, one of `SMALL`, `MEDIUM`, `LARGE`, `XLARGE`, `DLB`. Default is `SMALL`. Please note that XLARGE is only supported since NSX 3.0.0
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `connectivity_path` - (Optional) Tier1 Gateway where this service is instantiated. In future, other objects will be supported.
* `enabled` - (Optional) Flag to enable the service.
* `error_log_level` - (Optional) Log level for the service, one of `DEBUG`, `INFO`, `WARNING`, `ERROR`, `CRITICAL`, `ALERT`, `EMERGENCY`. Default is `INFO`.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing service can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_lb_service.test ID
```

The above command imports LBService named `test` with the NSX Load Balancer Service ID `ID`.
