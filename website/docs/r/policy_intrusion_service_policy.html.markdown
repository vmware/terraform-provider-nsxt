---
subcategory: "Firewall"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_intrusion_service_policy"
description: A resource to configure Intrusion Service Policy and its rules.
---

# nsxt_policy_intrusion_service_policy

This resource provides a method for the management of Intrusion Service (IDS) Policy and rules under it.

This resource is applicable to NSX Policy Manager and VMC (NSX version 3.1.0 onwards).

## Example Usage

```hcl
resource "nsxt_policy_intrusion_service_policy" "policy1" {
  display_name = "policy1"
  description  = "Terraform provisioned Policy"
  locked       = false
  stateful     = true

  rule {
    display_name       = "rule1"
    destination_groups = [nsxt_policy_group.cats.path, nsxt_policy_group.dogs.path]
    action             = "DETECT"
    services           = [nsxt_policy_service.icmp.path]
    logged             = true
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.default.path]
  }

  rule {
    display_name     = "rule2"
    source_groups    = [nsxt_policy_group.fish.path]
    sources_excluded = true
    action           = "DETECT_PREVENT"
    services         = [nsxt_policy_service.udp.path]
    logged           = true
    disabled         = true
    notes            = "Disabled till Sunday"
    ids_profiles     = [data.nsxt_policy_intrusion_service_profile.default.path]
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_intrusion_service_policy" "policy1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "policy1"
  description  = "Terraform provisioned Policy"
  locked       = false
  stateful     = true

  rule {
    display_name       = "rule1"
    destination_groups = [nsxt_policy_group.cats.path, nsxt_policy_group.dogs.path]
    action             = "DETECT"
    services           = [nsxt_policy_service.icmp.path]
    logged             = true
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.default.path]
  }

  rule {
    display_name     = "rule2"
    source_groups    = [nsxt_policy_group.fish.path]
    sources_excluded = true
    scope            = [nsxt_policy_group.aquarium.path]
    action           = "DETECT_PREVENT"
    services         = [nsxt_policy_service.udp.path]
    logged           = true
    disabled         = true
    notes            = "Disabled till Sunday"
    ids_profiles     = [data.nsxt_policy_intrusion_service_profile.someprofile.path]
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `domain` - (Optional) The domain to use for the resource. This domain must already exist. For VMware Cloud on AWS use `cgw`. If not specified, this field is default to `default`.
* `tag` - (Optional) A list of scope + tag pairs to associate with this policy.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
* `comments` - (Optional) Comments for IDS policy lock/unlock.
* `locked` - (Optional) Indicates whether the policy should be locked. If locked by a user, no other user would be able to modify this policy.
* `sequence_number` - (Optional) This field is used to resolve conflicts between IDS policies across domains.
* `stateful` - (Optional) If true, state of the network connects are tracked and a stateful packet inspection is performed. Default is true.
* `rule` - (Optional) A repeatable block to specify rules for the Policy. Each rule includes the following fields:
  * `display_name` - (Required) Display name of the resource.
  * `description` - (Optional) Description of the resource.
  * `action` - (Optional) Rule action, one of `DETECT`, `DETECT_PREVENT`. Default is `DETECT`.
  * `destination_groups` - (Optional) Set of group paths that serve as destination for this rule.
  * `source_groups` - (Optional) Set of group paths that serve as source for this rule.
  * `destinations_excluded` - (Optional) A boolean value indicating negation of destination groups.
  * `sources_excluded` - (Optional) A boolean value indicating negation of source groups.
  * `scope` - (Optional) Set of policy object paths where the rule is applied.
  * `direction` - (Optional) Traffic direction, one of `IN`, `OUT` or `IN_OUT`. Default is `IN_OUT`.
  * `disabled` - (Optional) Flag to disable this rule. Default is false.
  * `ip_version` - (Optional) Version of IP protocol, one of `IPV4`, `IPV6`, `IPV4_IPV6`. Default is `IPV4_IPV6`.
  * `logged` - (Optional) Flag to enable packet logging. Default is false.
  * `notes` - (Optional) Additional notes on changes.
  * `ids_profiles` - (Required) Set of IDS profile paths relevant for this rule.
  * `services` - (Optional) Set of service paths to match.
  * `log_label` - (Optional) Additional information (string) which will be propagated to the rule syslog.
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

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_intrusion_service_policy.policy1 domain/ID
```
The above command imports the policy named `policy1` under NSX domain `domain` with the NSX Policy ID `ID`.

```
terraform import nsxt_policy_intrusion_service_policy.policy1 POLICY_PATH
```
The above command imports the policy named `policy1` under NSX domain `domain` with the NSX policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
