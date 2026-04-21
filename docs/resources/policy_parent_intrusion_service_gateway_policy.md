---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_parent_intrusion_service_gateway_policy"
description: A resource to configure an Intrusion Service Gateway Policy without embedded rules for North-South traffic inspection.
---

# nsxt_policy_parent_intrusion_service_gateway_policy

This resource provides a method for the management of an Intrusion Service Gateway Policy without embedded rules. Users can use `nsxt_policy_intrusion_service_gateway_policy_rule` resource to add rules to this gateway policy if needed.

This resource is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_tier1_gateway" "tier1_gw" {
  display_name = "tier1-gateway"
}

data "nsxt_policy_intrusion_service_profile" "default_ids" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_parent_intrusion_service_gateway_policy" "parent_policy" {
  display_name    = "intrusion-svc-gw-parent-policy"
  description     = "Parent policy for standalone IDPS gateway rules"
  category        = "LocalGatewayRules"
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

resource "nsxt_policy_intrusion_service_gateway_policy_rule" "detect_inbound" {
  display_name    = "detect-inbound-threats"
  description     = "Detect threats in inbound traffic"
  policy_path     = nsxt_policy_parent_intrusion_service_gateway_policy.parent_policy.path
  action          = "DETECT_PREVENT"
  direction       = "IN"
  ip_version      = "IPV4"
  sequence_number = 1
  scope           = [data.nsxt_policy_tier1_gateway.tier1_gw.path]
  ids_profiles    = [data.nsxt_policy_intrusion_service_profile.default_ids.path]
  logged          = true
}
```

-> We recommend using `lifecycle` directive as in the sample above, in order to avoid dependency issues when updating groups/services simultaneously with the rule.

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `domain` - (Optional) The domain to use for the resource. This domain must already exist. If not specified, this field defaults to `default`.
* `tag` - (Optional) A list of scope + tag pairs to associate with this policy.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `category` - (Required) The category to use for priority of this Intrusion Service Gateway Policy. Must be one of: `SharedPreRules`, `LocalGatewayRules`, or `Default`.
* `comments` - (Optional) Comments for this Intrusion Service Gateway Policy including lock/unlock comments.
* `locked` - (Optional) A boolean value indicating if the policy is locked. If locked, no other users can update the resource. Default is `false`.
* `sequence_number` - (Optional) An int value used to resolve conflicts between intrusion service gateway policies across domains. Default is `0`.
* `stateful` - (Computed) A boolean value indicating if this Policy is stateful. Intrusion Service Gateway Policies are always stateful as they require connection state tracking for proper intrusion detection and prevention. This field is read-only and always returns `true`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `path` - The NSX path of the policy resource.

## Importing

An existing Intrusion Service Gateway Policy can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_parent_intrusion_service_gateway_policy.parent_policy POLICY_PATH
```

The above command imports the policy named `parent_policy` with the policy path `POLICY_PATH`.

If the Policy to import isn't in the `default` domain, the domain name can be added to the `ID` before a slash.

```shell
terraform import nsxt_policy_parent_intrusion_service_gateway_policy.parent_policy MyDomain/ID
```
