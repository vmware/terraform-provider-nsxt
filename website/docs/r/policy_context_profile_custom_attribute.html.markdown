---
subcategory: "Firewall"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_context_profile_custom_attribute"
description: A resource to configure a Context Profile FQDN or URL Custom attribute.
---

# nsxt_policy_context_profile_custom_attribute

This resource provides a method for the management of a Context Profile FQDN or URL Custom attributes.
This resource is supported with NSX 4.1.0 onwards.

## Example Usage

```hcl
resource "nsxt_policy_context_profile_custom_attribute" "test" {
  key       = "DOMAIN_NAME"
  attribute = "test.somesite.com"
}

```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_context_profile_custom_attribute" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  key       = "DOMAIN_NAME"
  attribute = "test.somesite.com"
}
```

## Argument Reference

The following arguments are supported:
Note: `key`, `attribute` must be present.

* `key` - (Required) Policy Custom Attribute Key. Valid values are "DOMAIN_NAME" and "CUSTOM_URL"
* `attribute` - (Required) FQDN or URL to be used as custom attribute.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Importing

An existing Context Profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_context_profile_custom_attribute.test DOMAIN_NAME~test.somesite.com
```

The above command imports Context Profile FQDN attribute named `test` with FQDN `test.somesite.com`.
