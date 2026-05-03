---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_intrusion_service_policy"
description: A resource to configure Intrusion Service Policy and its rules for East-West traffic inspection.
---

# nsxt_policy_intrusion_service_policy

This resource provides a method for the management of Intrusion Service (IDS) Policy and rules under it.

This resource is applicable to NSX Policy Manager (NSX version 3.1.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_intrusion_service_profile" "default" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_intrusion_service_policy" "policy1" {
  display_name = "policy1"
  description  = "Terraform provisioned Policy"
  category     = "ThreatRules"
  locked       = false

  rule {
    display_name = "rule1"
    action       = "DETECT"
    services     = [nsxt_policy_service.http.path]
    logged       = true
    ids_profiles = [data.nsxt_policy_intrusion_service_profile.default.path]
  }

  rule {
    display_name          = "rule2"
    description           = ""
    action                = "DETECT_PREVENT"
    direction             = "IN"
    ip_version            = "IPV4"
    oversubscription      = "BYPASSED"
    source_groups         = [nsxt_policy_group.external_clients.path]
    destination_groups    = [nsxt_policy_group.web_servers.path]
    services              = [nsxt_policy_service.https.path]
    notes                 = "Disabled till Sunday"
    ids_profiles          = [data.nsxt_policy_intrusion_service_profile.default.path]
    logged                = true
    disabled              = true
    sources_excluded      = true
    destinations_excluded = true
  }

  rule {
    display_name       = "rule3"
    source_groups      = [nsxt_policy_group.trusted_admins.path]
    destination_groups = [nsxt_policy_group.external_clients.path, nsxt_policy_group.db_servers.path]
    action             = "EXEMPT"
    oversubscription   = "DROPPED"
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.default.path]
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_intrusion_service_profile" "default" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_intrusion_service_policy" "policy1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "policy1"
  description  = "Terraform provisioned Policy"
  locked       = false

  rule {
    display_name       = "rule1"
    source_groups      = [nsxt_policy_group.cats.path]
    destination_groups = [nsxt_policy_group.cats.path, nsxt_policy_group.dogs.path]
    action             = "DETECT"
    services           = [nsxt_policy_service.icmp.path]
    logged             = true
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.default.path]
  }

  rule {
    display_name       = "rule2"
    source_groups      = [nsxt_policy_group.fish.path]
    destination_groups = [nsxt_policy_group.cats.path]
    sources_excluded   = true
    scope              = [nsxt_policy_group.aquarium.path]
    action             = "DETECT_PREVENT"
    services           = [nsxt_policy_service.udp.path]
    logged             = true
    disabled           = true
    notes              = "Disabled till Sunday"
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.default.path]
  }
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
* `category` - (Optional) Category of this policy. Must be one of: `ThreatRules` or `EmergencyThreatRules`. Default is `ThreatRules`.
* `comments` - (Optional) Comments for this Intrusion Service Policy including lock/unlock comments.
* `locked` - (Optional) A boolean value indicating if the policy is locked. If locked, no other users can update the resource. Default is `false`.
* `sequence_number` - (Optional) An int value used to resolve conflicts between intrusion service policies across domains. Default is `0`.
* `stateful` - (Computed) A boolean value indicating if this Policy is stateful. Intrusion Service Policies are always stateful as they require connection state tracking for proper intrusion detection and prevention. This field is read-only and always returns `true`.
* `rule` - (Optional) A repeatable block to specify rules for the Policy. Each rule includes the following fields:
    * `display_name` - (Required) Display name of the resource.
    * `description` - (Optional) Description of the resource.
    * `action` - (Optional) Rule action, one of `DETECT`, `DETECT_PREVENT`, `EXEMPT`. Default is `DETECT`. Note: `EXEMPT` is only supported from NSX version 9.1.0 onwards.
    * `destination_groups` - (Optional) Set of group paths that serve as the destination for this rule. An empty set can be used to specify `ANY`. Default is `ANY`.
    * `source_groups` - (Optional) Set of group paths that serve as the source for this rule. An empty set can be used to specify `ANY`. Default is `ANY`.
    * `destinations_excluded` - (Optional) A boolean value indicating negation of destination groups. Default is `false`.
    * `sources_excluded` - (Optional) A boolean value indicating negation of source groups. Default is `false`.
    * `scope` - (Optional) Set of policy object paths where the rule is applied for East-West traffic inspection.
    * `direction` - (Optional) The traffic direction for the rule. Must be one of: `IN`, `OUT` or `IN_OUT`. Default is `IN_OUT`.
    * `disabled` - (Optional) A boolean value to indicate the rule is disabled. Default is `false`.
    * `ip_version` - (Optional) The IP Protocol for the rule. Must be one of: `IPV4`, `IPV6` or `IPV4_IPV6`. Default is `IPV4_IPV6`.
    * `logged` - (Optional) A boolean flag to enable packet logging. Default is `false`.
    * `notes` - (Optional) Text for additional notes on changes for this rule.
    * `ids_profiles` - (Required) Set of IDS profile paths for this rule. These profiles define the intrusion detection signatures to be applied.
    * `oversubscription` - (Optional) Action to take when IDPS engine is oversubscribed. One of `BYPASSED`, `DROPPED` or `INHERIT_GLOBAL`. Default is `INHERIT_GLOBAL`. `BYPASSED`: Traffic bypasses IDPS when oversubscribed. `DROPPED`: Traffic is dropped when oversubscribed. `INHERIT_GLOBAL`: Inherit the behavior from the global IDPS settings.
    * `services` - (Optional) Set of service paths to match for this rule. An empty set can be used to specify `ANY`. Default is `ANY`.
    * `log_label` - (Optional) Additional information (string) which will be propagated to the rule syslog for this rule.
    * `tag` - (Optional) A list of scope + tag pairs to associate with this Rule.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the IDS Policy.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `rule`:
    * `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
    * `path` - The NSX policy path for this rule.
    * `sequence_number` - Sequence number for this rule, as defined by order of rules in the list.
    * `rule_id` - Unique positive number that is assigned by the system and is useful for debugging.

## Importing

An existing policy can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_intrusion_service_policy.policy1 domain/ID
```

The above command imports the policy named `policy1` under NSX domain `domain` with the NSX Policy ID `ID`.

```shell
terraform import nsxt_policy_intrusion_service_policy.policy1 POLICY_PATH
```

The above command imports the policy named `policy1` under NSX domain `domain` with the NSX policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
