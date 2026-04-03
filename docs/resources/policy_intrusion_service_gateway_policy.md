---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_intrusion_service_gateway_policy"
description: A resource to configure an Intrusion Service Gateway Policy and its rules.
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
  stateful        = true

  rule {
    display_name       = "detect-inbound-threats"
    description        = "Detect threats in North-South inbound traffic"
    direction          = "IN"
    action             = "DETECT"
    scope              = [data.nsxt_policy_tier1_gateway.tier1_gw.path]
    source_groups      = [nsxt_policy_group.web_servers.path]
    destination_groups = [nsxt_policy_group.db_servers.path]
    services           = [nsxt_policy_service.http.path]
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.default.path]
    logged             = true
  }

  rule {
    display_name = "detect-prevent-outbound"
    description  = "Detect and prevent threats in outbound traffic"
    direction    = "OUT"
    action       = "DETECT_PREVENT"
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
* `domain` - (Optional) The domain to use for the resource. Defaults to `default`.
* `tag` - (Optional) A list of scope + tag pairs to associate with this policy.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `category` - (Required) Category of this policy. Must be one of: `Emergency`, `SystemRules`, `SharedPreRules`, `LocalGatewayRules`, `AutoServiceRules`, or `Default`. ForceNew.
* `comments` - (Optional) Comments for security policy lock/unlock.
* `locked` - (Optional) Indicates whether a security policy should be locked. Default is `false`.
* `sequence_number` - (Optional) This field is used to resolve conflicts between security policies across domains. Default is `0`.
* `stateful` - (Optional) When it is stateful, the state of the network connects are tracked and a stateful packet inspection is performed. Default is `true`.
* `tcp_strict` - (Optional) Ensures that a 3-way TCP handshake is done before the data packets are sent. Computed if not set.
* `rule` - (Optional) A list of rules for this policy. Each rule supports the following:
    * `display_name` - (Required) Display name of the rule.
    * `description` - (Optional) Description of the rule.
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
    * `sequence_number` - (Optional) Sequence number of this rule.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `path` - The NSX path of the policy resource.
* `rule`:
    * `revision` - Indicates current revision number of the rule.
    * `path` - The NSX path of the rule.
    * `sequence_number` - Sequence number of the rule.
    * `nsx_id` - NSX ID of the rule.

## Importing

An existing Intrusion Service Gateway Policy can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```shell
terraform import nsxt_policy_intrusion_service_gateway_policy.north_south_detect POLICY_PATH
```

The above command imports the policy named `north_south_detect` with the policy path `POLICY_PATH`.
Example: `/infra/domains/default/intrusion-service-gateway-policies/idps-gateway-policy`.
