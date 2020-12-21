---
layout: "nsxt"
page_title: "Frequently Asked Questions"
description: |-
  Frequently Asked Questions and Workarounds
---

<!-- TOC -->

- [Dependency Errors on Update or Destroy](#dependency-errors-lifecycle)
- [Dependency Error on Destroy](#dependency-error-nsx)
- [Authorization Error on VMC ](#auth-error-vmc)

<!-- /TOC -->

## Dependency Error on Update or Destroy

Consider the following error:
```
Error:  Failed to delete <object>: The object path=[..] cannot be deleted as either it has children or it is being referenced by other objects..
```

Usually this error results from terraform engine assuming certain order of delete/update operation that is not consistent with NSX. In order to imply correct order on terraform and thus fix the issue, add the following clause to affected resources:

```
resource "nsxt_policy_group" "example" {
  # ...

  lifecycle {
    create_before_destroy = true
  }
}
```

However, sometimes the error above is symptom of misconfiguration, i.e. there are legitimate dependencies on the platform that prevent given delete/update.


## Dependency Error on Destroy

Consider same error as above:
```
Error:  Failed to delete <object>: The object path=[..] cannot be deleted as either it has children or it is being referenced by other objects..
```

Sometimes this error is due to the fact that certain resource cleanup on NSX needs more time. For now, the workaround would be to re-run the destroy command after few seconds. In future versions of provider, this issue will be solved with automatic retry.


## Authorization Error on VMC

Consider the following error:
```
User is not authorized to perform this operation on the application. Please contact the system administrator to get access. (code 401)
```

Assuming user permissions are sufficient, this error is usually due to missing `domain` configuration on the resource. For example, configuring group resource on VMC needs to contain domain configuration:

```hcl
resource "nsxt_policy_group" "group1" {
  display_name = "tf-group1"
  domain       = "cgw"
}
```

Be sure to also specify `enforcement_point` as `vmc-enforcementpoint` in provider section.
