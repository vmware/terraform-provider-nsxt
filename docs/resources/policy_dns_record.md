---
subcategory: "DNS"
page_title: "NSXT: nsxt_policy_dns_record"
description: A resource to configure a Policy DNS Record.
---

# nsxt_policy_dns_record

This resource provides a method for the management of DNS records within an NSX project.

A DNS record (type `ProjectDnsRecord`) resides in a DNS zone referenced by `zone_path`. Supports A, AAAA, CNAME, PTR, and NS record types. The `record_type` and `zone_path` fields are immutable after creation.

This resource is applicable to NSX Policy Manager and requires NSX 9.2.0 or higher.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_dns_record" "example_a" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }

  display_name  = "www-record"
  description   = "Terraform provisioned A record"
  record_name   = "www"
  record_type   = "A"
  record_values = ["192.168.1.10"]
  zone_path     = nsxt_policy_dns_zone.example.path
  ttl           = 300
}

resource "nsxt_policy_dns_record" "example_cname" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }

  display_name  = "alias-record"
  record_name   = "alias"
  record_type   = "CNAME"
  record_values = ["canonical.example.com."]
  zone_path     = nsxt_policy_dns_zone.example.path
  ttl           = 300
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Required) The context which the object belongs to.
    * `project_id` - (Required) The ID of the project which the object belongs to.
* `record_name` - (Required) DNS record name or, for PTR records, the host-octet label within the reverse zone.
* `record_type` - (Required, Force New) DNS record type. One of `A`, `AAAA`, `CNAME`, `PTR`, `NS`. Immutable after creation.
* `record_values` - (Required) DNS record data values consistent with `record_type`. At least one value is required.
* `zone_path` - (Required, Force New) Policy path to the `ProjectDnsZone` for this record. Immutable after creation.
* `ttl` - (Optional) Time-To-Live in seconds (30-86400). Overrides the zone's default TTL. Default: `300`.
* `ip_address` - (Optional) IPv4 address mapped by a PTR record. Only applicable when `record_type` is `PTR`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `path` - The NSX policy path of the resource.
* `fqdn` - System-computed fully qualified domain name (read-only).

## Importing

An existing object can be [imported][docs-import] into this resource using its full policy path:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_dns_record.example PATH
```

The expected path format is `/orgs/[org]/projects/[project]/dns-records/[dns-record]`.
