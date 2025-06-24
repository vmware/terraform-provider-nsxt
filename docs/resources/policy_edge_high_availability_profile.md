---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_edge_high_availability_profile"
description: A resource to configure a PolicyEdgeHighAvailabilityProfile.
---

# nsxt_policy_edge_high_availability_profile

This resource provides a method for the management of a PolicyEdgeHighAvailabilityProfile.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_edge_high_availability_profile" "test" {
  display_name              = "test"
  description               = "Terraform provisioned PolicyEdgeHighAvailabilityProfile"
  bfd_probe_interval        = 750
  bfd_allowed_hops          = 140
  bfd_declare_dead_multiple = 7

}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `site_path` - (Optional) The path of the site which the Edge Transport Node belongs to. `path` field of the existing `nsxt_policy_site` can be used here. Defaults to default site path.
* `enforcement_point` - (Optional) The ID of enforcement point under given `site_path` to manage the Edge Transport Node. Defaults to default enforcement point.
* `bfd_probe_interval` - (Optional) The time interval (in millisecond) between probe packets for heartbeat purpose. Default: 500.
* `bfd_allowed_hops` - (Optional) Value of BFD allowed hops. Default: 255.
* `bfd_declare_dead_multiple` - (Optional) Number of times a packet is missed before BFD declares the neighbor down. Default: 3.
* `standby_relocation_config` - (Optional) Stand by relocation flag.
    * `standby_relocation_threshold` - (Required) Standby service context relocation wait time.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_edge_high_availability_profile.test PATH
```

The above command imports PolicyEdgeHighAvailabilityProfile named `test` with the NSX path `PATH`.
