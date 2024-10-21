---
subcategory: "Realization"
layout: "nsxt"
page_title: "NSXT: compute_manager_realization"
description: Compute manager resource realization information.
---

# nsxt_compute_manager_realization

This data source provides information about the realization of a compute manager resource on NSX manager. This data source will wait until realization is determined as either success or error. It is recommended to use this data source if further configuration depends on compute manager realization.

## Example Usage

```hcl
resource "nsxt_compute_manager" "test" {
  description  = "Terraform provisioned Compute Manager"
  display_name = "test"
  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  server = "192.168.244.144"

  credential {
    username_password_login {
      username = "user"
      password = "pass"
    }
  }
  origin_type = "vCenter"
}

data "nsxt_compute_manager_realization" "test" {
  id      = nsxt_compute_manager.test.id
  timeout = 60
}
```

## Argument Reference

* `id` - (Required) ID of the resource.
* `delay` - (Optional) Delay (in seconds) before realization polling is started. Default is set to 1.
* `timeout` - (Optional) Timeout (in seconds) for realization polling. Default is set to 1200.
* `check_registration` - (Optional) Check if registration of compute manager is complete.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `state` - The realization state of the resource. Transitional states are: "pending", "in_progress", "in_sync", "unknown". Target states are: "success", "failed", "partial_success", "orphaned", "error".
* `registration_status` - Overall registration status of desired configuration. Transitional statuses are "CONNECTING", "REGISTERING". Target statuses are: "REGISTERED", "UNREGISTERED", "REGISTERED_WITH_ERRORS"
