---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_network_span"
description: A resource to configure a NetworkSpan.
---

# nsxt_policy_network_span

This resource provides a method for the management of a NetworkSpan.

This resource is applicable to NSX Global Manager, NSX Policy Manager with version 9.1.0 onwards and VMC.

## Example Usage

```hcl
resource "nsxt_policy_network_span" "netspan" {
  display_name = "test"
  description  = "Terraform provisioned NetworkSpan"
  exclusive    = true

}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `exclusive` - (Optional) When exclusive span enabled, VC clusters in  the span are not shared across other Span including system default span.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```shell
terraform import nsxt_policy_network_span.netspan PATH
```

The above command imports NetworkSpan named `netspan` with the NSX path `PATH`.
