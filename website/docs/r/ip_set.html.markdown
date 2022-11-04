---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_ip_set"
description: A resource that can be used to configure an IP set in NSX.
---

# nsxt_ip_set

This resources provides a way to configure an IP set in NSX. An IP set is a collection of IP addresses. It is often used in the configuration of the NSX firewall.

## Example Usage

```hcl
resource "nsxt_ip_set" "ip_set1" {
  description  = "IS provisioned by Terraform"
  display_name = "IS"

  tag {
    scope = "color"
    tag   = "blue"
  }

  ip_addresses = ["1.1.1.1", "2.2.2.2"]
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this IP set.
* `ip_addresses` - (Optional) IP addresses.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the IP set.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing IP set can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_ip_set.ip_set1 UUID
```

The above command imports the IP set named `ip_set1` with the NSX id `UUID`.
