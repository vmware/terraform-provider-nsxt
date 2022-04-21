---
layout: "nsxt"
page_title: "Frequently Asked Questions"
description: |-
  Frequently Asked Questions and Workarounds
---

<!-- TOC depthFrom:2 -->

- [Dependency Errors on Update or Destroy](#dependency-errors-lifecycle)
- [Dependency Error on Destroy](#dependency-error-nsx)
- [Authorization Error on VMC](#auth-error-vmc)
- [Error Unmarshalling Server Response on VMC](#proxy-vmc)
- [User is not authorized to perform this operation on the application on VMC](#domain-vmc)

<!-- /TOC -->

## Dependency Error on Update or Destroy <a name="dependency-errors-lifecycle"></a>

Consider the following error:
```
Error:  Failed to delete <object>: The object path=[..] cannot be deleted as either it
has children or it is being referenced by other objects..
```

Usually this error results from terraform engine assuming certain order of delete/update operation that is not consistent with NSX. In order to imply correct order on terraform and thus fix the issue, add the following clause to affected resources:

```hcl
resource "nsxt_policy_group" "example" {
  # ...

  lifecycle {
    create_before_destroy = true
  }
}
```

However, sometimes the error above is symptom of misconfiguration, i.e. there are legitimate dependencies on the platform that prevent given delete/update.


## Dependency Error on Destroy <a name="dependency-error-nsx"></a>

Consider same error as above:
```
Error:  Failed to delete <object>: The object path=[..] cannot be deleted as either it has
children or it is being referenced by other objects..
```

Sometimes this error is due to the fact that certain resource cleanup on NSX needs more time. For now, the workaround would be to re-run the destroy command after few seconds. In future versions of provider, this issue will be solved with automatic retry.


## Authorization Error on VMC <a name="auth-error-vmc"></a>

Consider the following error:
```
User is not authorized to perform this operation on the application.
Please contact the system administrator to get access. (code 401)
```

Assuming user permissions are sufficient, this error is usually due to missing `domain` configuration on the resource. For example, configuring group resource on VMC needs to contain domain configuration:

```hcl
resource "nsxt_policy_group" "group1" {
  display_name = "tf-group1"
  domain       = "cgw"
}
```

Be sure to also specify `enforcement_point` as `vmc-enforcementpoint` in provider section.


## Error Unmarshalling Server Response on VMC <a name="proxy-vmc"></a>

Consider the following error:
```
Failed to read <object type> (Error unmarshalling server response)
```

This issue is usually caused by proxy timeout on VMC side. Re-apply can help. If you run into this a lot, please consider direct connection to NSX in your VMC environment to avoid the proxy overhead. If this is not an option for you, another way to reduce the load on the proxy could be to avoid heavy usage of data sources, and instead use policy path directly in your configuration.


## User is not authorized to perform this operation on the application on VMC <a name="domain-vmc"></a>

Consider the following error:
```
User is not authorized to perform this operation on the application.
Please contact the system administrator to get access. (code 401)
```

Provided your VMC token is accurate and you have sufficient permissions, it is likely that `domain` is missing from resource configuration. Domains of VMC are different from default value that is set in the provider.For group resource, the fix would be specifying relevant domain like in the example below:

```
resource "nsxt_policy_group" "test" {
  display_name = "test"
  domain = "cgw"

  ..
}
```

