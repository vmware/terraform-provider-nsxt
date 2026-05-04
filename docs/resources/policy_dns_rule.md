---
subcategory: "DNS"
page_title: "NSXT: nsxt_policy_dns_rule"
description: A resource to configure a Policy DNS Rule.
---

# nsxt_policy_dns_rule

This resource provides a method for the management of DNS rules attached to a Policy DNS Service.

A DNS rule defines an action applied to DNS queries whose domain matches the rule's `domain_patterns` using longest-prefix match. Supported actions are `UPDATE_MEMBERSHIP` (dynamically populate an FQDN Group with resolved IPs) and `FORWARD` (proxy to upstream DNS servers or to a shared DNS zone).

This resource is applicable to NSX Policy Manager and requires NSX 9.2.0 or higher.

## Example Usage

```hcl
resource "nsxt_policy_dns_rule" "forward_rule" {
  parent_path      = nsxt_policy_dns_service.svc.path
  display_name     = "forward-rule"
  description      = "Terraform provisioned DNS forward rule"
  action_type      = "FORWARD"
  domain_patterns  = ["*.external.com", "updates.example.com"]
  upstream_servers = ["8.8.8.8", "8.8.4.4"]

  tag {
    scope = "env"
    tag   = "prod"
  }
}

resource "nsxt_policy_dns_rule" "shared_zone_rule" {
  parent_path      = nsxt_policy_dns_service.svc.path
  display_name     = "shared-zone-rule"
  action_type      = "FORWARD"
  shared_zone_path = nsxt_policy_dns_zone.shared.path
}

resource "nsxt_policy_dns_rule" "membership_rule" {
  parent_path     = nsxt_policy_dns_service.svc.path
  display_name    = "membership-rule"
  action_type     = "UPDATE_MEMBERSHIP"
  domain_patterns = ["*.internal.example.com"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `parent_path` - (Required, Force New) Policy path of the parent `nsxt_policy_dns_service`.
* `action_type` - (Required) Action to apply to DNS queries matching `domain_patterns`. One of `FORWARD` or `UPDATE_MEMBERSHIP`.
* `domain_patterns` - (Optional) Domain name patterns matched via longest-prefix match. Supports wildcards (e.g. `*.example.com`). Required unless `action_type` is `FORWARD` and `shared_zone_path` is set (in which case domain patterns are derived from the shared zone).
* `upstream_servers` - (Optional) Upstream DNS server IP addresses to forward matching queries to. Only valid when `action_type` is `FORWARD`. Mutually exclusive with `shared_zone_path`. Maximum 3 entries.
* `shared_zone_path` - (Optional) Policy path to a `DnsZone` shared with this project. Only valid when `action_type` is `FORWARD`. Mutually exclusive with `upstream_servers`. When set, `domain_patterns` may be omitted.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `path` - The NSX policy path of the resource.

## Importing

An existing object can be [imported][docs-import] into this resource using its full policy path:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_dns_rule.example PATH
```

The expected import ID is the full policy path: `/orgs/[org]/projects/[project]/dns-services/[dns-service]/rules/[rule]`.
The `parent_path` attribute is automatically derived from the import path.
