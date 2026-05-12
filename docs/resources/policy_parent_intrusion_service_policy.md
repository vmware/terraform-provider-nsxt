---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_parent_intrusion_service_policy"
description: A resource to configure an Intrusion Service (IDS) Policy without embedded rules.
---

# nsxt_policy_parent_intrusion_service_policy

This resource provides a method for the management of an Intrusion Service (IDS) Policy without embedded rules.

Note: to avoid unexpected behavior, don't use this resource and resource `nsxt_policy_intrusion_service_policy` to manage the same Intrusion Service Policy at the same time.
To configure rules under this resource, please use resource `nsxt_policy_intrusion_service_policy_rule` to manage rules separately.

This resource is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_intrusion_service_profile" "default" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_parent_intrusion_service_policy" "parent_policy" {
  display_name    = "intrusion-svc-parent-policy"
  description     = "Parent policy for standalone IDS rules"
  category        = "ThreatRules"
  locked          = false
  sequence_number = 3

  tag {
    scope = "env"
    tag   = "production"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "nsxt_policy_intrusion_service_policy_rule" "detect_rule" {
  display_name     = "detect-threats"
  description      = "Detect threats in East-West traffic"
  policy_path      = nsxt_policy_parent_intrusion_service_policy.parent_policy.path
  action           = "DETECT_PREVENT"
  direction        = "IN"
  ip_version       = "IPV4"
  oversubscription = "BYPASSED"
  sequence_number  = 1
  ids_profiles     = [data.nsxt_policy_intrusion_service_profile.default.path]
  logged           = true
}
```

-> We recommend using `lifecycle` directive as in the sample above, in order to avoid dependency issues when updating groups/services simultaneously with the rule.

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_parent_intrusion_service_policy" "parent_policy" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name    = "intrusion-svc-parent-policy"
  description     = "Parent policy for standalone IDS rules"
  locked          = false
  sequence_number = 3
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `domain` - (Optional) The domain to use for the resource. This domain must already exist. If not specified, this field is default to `default`.
* `tag` - (Optional) A list of scope + tag pairs to associate with this policy.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `category` - (Optional) Category for this policy. Must be one of: `ThreatRules` or `EmergencyThreatRules`. Default is `ThreatRules`.
* `comments` - (Optional) Comments for IDS policy lock/unlock.
* `locked` - (Optional) Indicates whether the policy should be locked. If locked by a user, no other user would be able to modify this policy. Default is `false`.
* `sequence_number` - (Optional) This field is used to resolve conflicts between IDS policies across domains. Default is `0`.
* `stateful` - (Computed) A boolean value indicating if this Policy is stateful. Intrusion Service Policies are always stateful as they require connection state tracking for proper intrusion detection and prevention. This field is read-only and always returns `true`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `path` - The NSX path of the policy resource.

## Importing

An existing Intrusion Service Policy can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_parent_intrusion_service_policy.parent_policy domain/ID
```

The above command imports the policy named `parent_policy` under NSX domain `domain` with the NSX Policy ID `ID`.

```shell
terraform import nsxt_policy_parent_intrusion_service_policy.parent_policy POLICY_PATH
```

The above command imports the policy named `parent_policy` with the NSX policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
