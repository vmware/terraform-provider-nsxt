---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_intrusion_service_gateway_policy"
description: A resource to configure an Intrusion Service Gateway Policy and its rules for North-South traffic inspection.
---

# nsxt_policy_intrusion_service_gateway_policy

This resource provides a method for the management of an Intrusion Service Gateway Policy and its rules for North-South traffic (Gateway Firewall context).

This resource is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_tier1_gateway" "tier1_gw" {
  display_name = "tier1-gateway"
}

data "nsxt_policy_intrusion_service_profile" "default" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_intrusion_service_gateway_policy" "idps_gateway_policy" {
  display_name    = "intrusion-service-gateway-policy"
  description     = "North-South IDPS policy for gateway traffic"
  category        = "LocalGatewayRules"
  locked          = false
  sequence_number = 3

  rule {
    display_name       = "detect-inbound-threats"
    description        = "Detect threats in North-South inbound traffic"
    action             = "DETECT"
    direction          = "IN"
    ip_version         = "IPV4_IPV6"
    source_groups      = [nsxt_policy_group.web_servers.path]
    destination_groups = [nsxt_policy_group.db_servers.path]
    services           = [nsxt_policy_service.http.path]
    scope              = [data.nsxt_policy_tier1_gateway.tier1_gw.path]
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.default.path]
    logged             = true
  }

  rule {
    display_name = "detect-prevent-outbound"
    description  = "Detect and prevent threats in outbound traffic"
    action       = "DETECT_PREVENT"
    direction    = "OUT"
    ip_version   = "IPV4"
    scope        = [data.nsxt_policy_tier1_gateway.tier1_gw.path]
    ids_profiles = [data.nsxt_policy_intrusion_service_profile.default.path]
    logged       = true
  }

  lifecycle {
    create_before_destroy = true
  }
}
```

-> We recommend using `lifecycle` directive as in sample above, in order to avoid dependency issues when updating groups/services simultaneously with the rule.

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `domain` - (Optional) The domain to use for the resource. This domain must already exist. If not specified, this field is default to `default`.
* `tag` - (Optional) A list of scope + tag pairs to associate with this policy.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `category` - (Required) The category to use for priority of this Intrusion Service Gateway Policy. Must be one of: `SharedPreRules`, `LocalGatewayRules`, or `Default`.
* `comments` - (Optional) Comments for this Intrusion Service Gateway Policy including lock/unlock comments.
* `locked` - (Optional) A boolean value indicating if the policy is locked. If locked, no other users can update the resource. Default is `false`.
* `sequence_number` - (Optional) An int value used to resolve conflicts between intrusion service gateway policies across domains. Default is `0`.
* `stateful` - (Computed) A boolean value indicating if this Policy is stateful. Intrusion Service Gateway Policies are always stateful as they require connection state tracking for proper intrusion detection and prevention. This field is read-only and always returns `true`.
* `rule` - (Optional) A list of rules for this policy. Each rule supports the following:
    * `display_name` - (Required) Display name of the rule.
    * `description` - (Optional) Description of the rule.
    * `destination_groups` - (Optional) Set of group paths that serve as the destination for this rule. An empty set can be used to specify `ANY`. Default is `ANY`.
    * `destinations_excluded` - (Optional) A boolean value indicating negation of destination groups. Default is `false`.
    * `direction` - (Optional) The traffic direction for the rule. Must be one of: `IN`, `OUT` or `IN_OUT`. Default is `IN_OUT`.
    * `disabled` - (Optional) A boolean value to indicate the rule is disabled. Default is `false`.
    * `ip_version` - (Optional) The IP Protocol for the rule. Must be one of: `IPV4`, `IPV6` or `IPV4_IPV6`. Default is `IPV4_IPV6`.
    * `logged` - (Optional) A boolean flag to enable packet logging. Default is `false`.
    * `notes` - (Optional) Text for additional notes on changes for the rule.
    * `scope` - (Required) Set of policy paths where the rule is applied. These should be Tier-0 or Tier-1 gateway paths for North-South traffic inspection.
    * `services` - (Optional) Set of service paths to match for this rule. An empty set can be used to specify `ANY`. Default is `ANY`.
    * `source_groups` - (Optional) Set of group paths that serve as the source for this rule. An empty set can be used to specify `ANY`. Default is `ANY`.
    * `sources_excluded` - (Optional) A boolean value indicating negation of source groups. Default is `false`.
    * `log_label` - (Optional) Additional information which will be propagated to the rule syslog.
    * `tag` - (Optional) A list of scope + tag pairs to associate with this rule.
    * `action` - (Optional) Rule action for intrusion detection/prevention. One of `DETECT` or `DETECT_PREVENT`. Default is `DETECT`.
    * `ids_profiles` - (Required) Set of IDS profile paths for this rule. These profiles define the intrusion detection signatures to be applied.
    * `sequence_number` - (Optional) Sequence number to determine the order of rule processing within this policy.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `path` - The NSX path of the policy resource.
* `rule`:
    * `revision` - Indicates current revision number of the rule.
    * `path` - The NSX path of the rule.
    * `sequence_number` - Sequence number of the rule.
    * `rule_id` - Unique positive number that is assigned by the system and is useful for debugging.

## Importing

An existing Intrusion Service Gateway Policy can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_intrusion_service_gateway_policy.north_south_detect POLICY_PATH
```

The above command imports the policy named `north_south_detect` with the policy path `POLICY_PATH`.
Example: `/infra/domains/default/intrusion-service-gateway-policies/idps-gateway-policy`.
