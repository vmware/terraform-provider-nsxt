---
subcategory: "Firewall"
page_title: "NSXT: nsxt_policy_parent_gateway_policy"
description: A resource to Gateway security policies without rules.
---

# nsxt_policy_parent_gateway_policy

This resource provides a method for the management of a Gateway Policy without rules. Users can use nsxt_policy_gateway_policy_rule resource to add rules to this gateway policy if needed.

This resource is applicable to NSX Global Manager and NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_parent_gateway_policy" "test" {
  display_name    = "tf-gw-policy"
  description     = "Terraform provisioned Gateway Policy"
  category        = "LocalGatewayRules"
  locked          = false
  sequence_number = 3
  stateful        = true
  tcp_strict      = false

  tag {
    scope = "color"
    tag   = "orange"
  }

  lifecycle {
    create_before_destroy = true
  }
}
```

## Example Usage with Global Manager

```hcl
resource "nsxt_policy_domain" "france" {
  display_name = "France"
  sites        = ["Paris"]
}

resource "nsxt_policy_parent_gateway_policy" "test" {
  display_name    = "tf-gw-policy"
  description     = "Terraform provisioned Gateway Policy"
  domain          = nsxt_policy_domain.france.id
  category        = "LocalGatewayRules"
  locked          = false
  sequence_number = 3
  stateful        = true
  tcp_strict      = false

  tag {
    scope = "color"
    tag   = "orange"
  }

  lifecycle {
    create_before_destroy = true
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_parent_gateway_policy" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name    = "tf-gw-policy"
  description     = "Terraform provisioned Gateway Policy"
  category        = "LocalGatewayRules"
  locked          = false
  sequence_number = 3
  stateful        = true
  tcp_strict      = false

  tag {
    scope = "color"
    tag   = "orange"
  }

  lifecycle {
    create_before_destroy = true
  }
}
```

-> We recommend using `lifecycle` directive as in samples above, in order to avoid dependency issues when updating groups/services simultaneously with the rule.

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `category` - (Required) The category to use for priority of this Gateway Policy. For local manager must be one of: `Emergency`, `SystemRules`, `SharedPreRules`, `LocalGatewayRules`, `AutoServiceRules` and `Default`. For global manager must be `SharedPreRules` or `LocalGatewayRules`.
* `description` - (Optional) Description of the resource.
* `domain` - (Optional) The domain to use for the Gateway Policy. This domain must already exist. For VMware Cloud on AWS use `cgw`.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Gateway Policy.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the Gateway Policy resource.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `comments` - (Optional) Comments for this Gateway Policy including lock/unlock comments.
* `locked` - (Optional) A boolean value indicating if the policy is locked. If locked, no other users can update the resource.
* `sequence_number` - (Optional) An int value used to resolve conflicts between security policies across domains
* `stateful` - (Optional) A boolean value to indicate if this Policy is stateful. When it is stateful, the state of the network connects are tracked and a stateful packet inspection is performed.
* `tcp_strict` - (Optional) A boolean value to enable/disable a 3 way TCP handshake is done before the data packets are sent.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Security Policy.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

~> **NOTE:** `display_name` argument for service entries is not supported for NSX 3.2.x and below.

## Importing

An existing Gateway Policy can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_parent_gateway_policy.gwpolicy1 ID
```

The above command imports the policy Gateway Policy named `gwpolicy1` with the NSX Policy id `ID`.

If the Policy to import isn't in the `default` domain, the domain name can be added to the `ID` before a slash.

For example to import a Group with `ID` in the `MyDomain` domain:

```shell
terraform import nsxt_policy_parent_gateway_policy.gwpolicy1 MyDomain/ID
```
