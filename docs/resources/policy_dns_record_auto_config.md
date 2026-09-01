---
subcategory: "DNS"
page_title: "NSXT: nsxt_policy_dns_record_auto_config"
description: A resource to configure a Policy DNS Record Auto Config.
---

# nsxt_policy_dns_record_auto_config

This resource provides a method for the management of a DNS Record Auto Config (`DnsAutoRecordConfig`) within an NSX project.

When an IP is allocated from the referenced IP block, an A record is automatically created in the `a_record_zone_path` DNS zone. An optional PTR record is also auto-created in the `ptr_record_zone_path` zone if set. The `ip_block_path` field is immutable after creation.

This resource is applicable to NSX Policy Manager and requires NSX 9.2.0 or higher.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_dns_record_auto_config" "example" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }

  display_name         = "vm-auto-records"
  description          = "Terraform provisioned DNS Auto Record Config"
  ip_block_path        = data.nsxt_policy_ip_block.workloads.path
  a_record_zone_path   = nsxt_policy_dns_zone.forward_zone.path
  ptr_record_zone_path = nsxt_policy_dns_zone.reverse_zone.path
  ttl                  = 300

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
* `context` - (Required) The context which the object belongs to.
    * `project_id` - (Required) The ID of the project which the object belongs to.
* `ip_block_path` - (Required, Force New) Policy path to the IP block from which workload IPs are allocated. Immutable after creation. Only one auto record config per IP block is allowed within a project.
* `a_record_zone_path` - (Required) Policy path to a `DnsZone` (local or shared) in which auto-created A records are placed. The zone must exist and be visible to this project.
* `ptr_record_zone_path` - (Optional) Policy path to a `DnsZone` (local or shared) in which auto-created PTR (reverse DNS) records are placed. When absent, no PTR record is auto-created.
* `ttl` - (Optional) Time-To-Live in seconds for auto-created A and PTR records (30-86400). Default: `300`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `path` - The NSX policy path of the resource.

## Importing

An existing object can be [imported][docs-import] into this resource using its full policy path:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_dns_record_auto_config.example PATH
```

The expected path format is `/orgs/[org]/projects/[project]/dns-auto-record-configs/[config]`.
