---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_intrusion_service_gateway_policy_rule"
description: A resource to configure a single rule in an Intrusion Service Gateway Policy.
---

# nsxt_policy_intrusion_service_gateway_policy_rule

This resource provides a method for the management of a single rule in an Intrusion Service Gateway Policy.

This resource is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_tier1_gateway" "tier1_gw" {
  display_name = "tier1-gateway"
}

data "nsxt_policy_intrusion_service_profile" "default" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_parent_intrusion_service_gateway_policy" "north_south_detect" {
  display_name    = "tf-intrusion-svc-gw-policy"
  category        = "LocalGatewayRules"
  locked          = false
  sequence_number = 3
  stateful        = true

  lifecycle {
    create_before_destroy = true
  }
}

resource "nsxt_policy_intrusion_service_gateway_policy_rule" "detect_inbound" {
  display_name    = "detect-inbound-threats"
  description     = "Detect threats in North-South inbound traffic"
  policy_path     = nsxt_policy_parent_intrusion_service_gateway_policy.north_south_detect.path
  action          = "DETECT"
  direction       = "IN"
  sequence_number = 1
  scope           = [data.nsxt_policy_tier1_gateway.tier1_gw.path]
  ids_profiles    = [data.nsxt_policy_intrusion_service_profile.default.path]
  logged          = true
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `policy_path` - (Required) Path of the Intrusion Service Gateway Policy this rule belongs to. ForceNew.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `destination_groups` - (Optional) Set of group paths that serve as the destination for this rule.
* `destinations_excluded` - (Optional) Negation of destination groups. Default is `false`.
* `direction` - (Optional) Traffic direction. One of `IN`, `OUT`, or `IN_OUT`. Default is `IN_OUT`.
* `disabled` - (Optional) Flag to disable the rule. Default is `false`.
* `ip_version` - (Optional) IP version. One of `IPV4`, `IPV6`, or `IPV4_IPV6`. Default is `IPV4_IPV6`.
* `logged` - (Optional) Flag to enable packet logging. Default is `false`.
* `notes` - (Optional) Text for additional notes on changes.
* `scope` - (Required) Set of policy paths where the rule is applied (e.g., Tier-0/Tier-1 gateway paths).
* `services` - (Optional) Set of service paths to match.
* `source_groups` - (Optional) Set of group paths that serve as the source for this rule.
* `sources_excluded` - (Optional) Negation of source groups. Default is `false`.
* `log_label` - (Optional) Additional information which will be propagated to the rule syslog.
* `tag` - (Optional) A list of scope + tag pairs to associate with this rule.
* `action` - (Optional) Rule action. One of `DETECT` or `DETECT_PREVENT`. Default is `DETECT`.
* `ids_profiles` - (Required) Set of IDS profile paths for this rule.
* `sequence_number` - (Required) Sequence number of this rule.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `path` - The NSX path of the policy resource.

## Importing

An existing Intrusion Service Gateway Policy Rule can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```shell
terraform import nsxt_policy_intrusion_service_gateway_policy_rule.detect_inbound RULE_PATH
```

Example: `/infra/domains/default/intrusion-service-gateway-policies/gateway-idps-policy/rules/detect-inbound-threats`.
