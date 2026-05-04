---
subcategory: "DNS"
page_title: "NSXT: nsxt_policy_dns_zone"
description: A resource to configure a Policy DNS Zone.
---

# nsxt_policy_dns_zone

This resource provides a method for the management of a DNS Zone attached to a Policy DNS Service.

A zone defines the domain name, optional VPC scope for split-horizon DNS, default TTL, and SOA parameters. The `dns_domain_name` field is immutable after creation.

This resource is applicable to NSX Policy Manager and requires NSX 9.2.0 or higher.

## Example Usage

```hcl
resource "nsxt_policy_dns_zone" "example" {
  parent_path     = nsxt_policy_dns_service.svc.path
  display_name    = "example-zone"
  description     = "Terraform provisioned DNS Zone"
  dns_domain_name = "example.com"
  ttl             = 300

  soa {
    primary_nameserver = "ns1.example.com."
    responsible_party  = "admin.example.com."
    refresh_interval   = 3600
    retry_interval     = 900
    expire_time        = 604800
    negative_cache_ttl = 300
  }

  tag {
    scope = "env"
    tag   = "prod"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `parent_path` - (Required, Force New) Policy path of the parent `nsxt_policy_dns_service`.
* `dns_domain_name` - (Required, Force New) The domain name for this zone (e.g. `example.com`) or reverse-notation (e.g. `12.168.192.in-addr.arpa`). Immutable after creation.
* `scope` - (Optional) Policy path to a single VPC within this project. When set, only workloads in the specified VPC can resolve this zone (split-horizon DNS). When unset, all VPCs can resolve the zone.
* `ttl` - (Optional) Default Time-To-Live in seconds for DNS records in this zone (30-86400). Default: `300`.
* `soa` - (Optional) SOA (Start of Authority) record parameters. All sub-fields are optional; omitted fields use system defaults.
    * `primary_nameserver` - (Optional) FQDN of the primary nameserver. Must end with a trailing dot.
    * `responsible_party` - (Optional) Zone administrator email in DNS format. Must end with a trailing dot.
    * `serial_number` - (Optional, Computed) Zone serial number.
    * `refresh_interval` - (Optional) Refresh interval in seconds.
    * `retry_interval` - (Optional) Retry interval in seconds.
    * `expire_time` - (Optional) Expire time in seconds.
    * `negative_cache_ttl` - (Optional) Negative cache TTL in seconds (0-86400).

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `path` - The NSX policy path of the resource.

## Importing

An existing object can be [imported][docs-import] into this resource using its full policy path:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_dns_zone.example PATH
```

The expected import ID is the full policy path: `/orgs/[org]/projects/[project]/dns-services/[dns-service]/zones/[zone]`.
The `parent_path` attribute is automatically derived from the import path.
