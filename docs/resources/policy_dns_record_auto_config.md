---
subcategory: "DNS"
page_title: "NSXT: nsxt_policy_dns_record_auto_config"
description: A resource to configure a Policy DNS Record Auto Config.
---

# nsxt_policy_dns_record_auto_config

This resource provides a method for the management of a DNS Record Auto Config (`ProjectDnsAutoRecordConfig`) within an NSX project.

When an IP is allocated from the referenced IP block, an A record is automatically created in the referenced DNS zone using the configured `naming_pattern`. The `ip_block_path` field is immutable after creation.

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

  display_name   = "vm-auto-records"
  description    = "Terraform provisioned DNS Auto Record Config"
  ip_block_path  = data.nsxt_policy_ip_block.workloads.path
  zone_path      = nsxt_policy_dns_zone.example.path
  naming_pattern = "{vm_name}-{index}"
  ttl            = 300

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
* `zone_path` - (Required) Policy path to a locally-owned `ProjectDnsZone` where auto-created A records are placed.
* `naming_pattern` - (Optional) Template for generating DNS record names from VM attributes. Supports `{vm_name}`, `{ip}`, `{index}`, `{vpc_name}`. Default: `{vm_name}`.
* `ttl` - (Optional) Time-To-Live in seconds for auto-created A records (30-86400). Default: `300`.

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
