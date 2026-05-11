---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_parent_intrusion_service_gateway_policy"
description: A data source to retrieve Parent Intrusion Service Gateway Policy metadata.
---

# nsxt_policy_parent_intrusion_service_gateway_policy

This data source provides metadata information about an existing Parent Intrusion Service Gateway Policy configured on NSX.
It can be useful for fetching policy path to use in `nsxt_policy_intrusion_service_gateway_policy_rule` resource.

~> **NOTE:** This data source retrieves only the parent policy metadata (id, display_name, description, path, sequence_number, etc.) **without embedded rules**, allowing you to refer a policy's path for creating standalone rules. For different use cases, consider:

* `nsxt_policy_intrusion_service_gateway_policy` - For **IDPS Gateway policy with embedded rules**
* `nsxt_policy_intrusion_service_gateway_policy_rule` - For **individual standalone IDPS Gateway rules**

This data source is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
# Get parent gateway policy metadata for rule creation
data "nsxt_policy_parent_intrusion_service_gateway_policy" "parent_gw_policy" {
  display_name = "production-gateway-ids-policy"
}

# Create standalone gateway rule using parent policy path
resource "nsxt_policy_intrusion_service_gateway_policy_rule" "new_gw_rule" {
  display_name = "new-north-south-detection"
  policy_path  = data.nsxt_policy_parent_intrusion_service_gateway_policy.parent_gw_policy.path
  action       = "DETECT"
  # ... other rule configuration
}
```

## Argument Reference

* `id` - (Optional) The ID of the policy to retrieve.
* `display_name` - (Optional) The display name of the policy to retrieve.
* `domain` - (Optional) The domain of the policy. Defaults to `default`.
* `category` - (Optional) Category of the policy.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `stateful` - When it is stateful, the state of the network connects are tracked and a stateful packet inspection is performed. Note: Intrusion Service Gateway Policies are always stateful.
* `locked` - Indicates whether a security policy should be locked. If the security policy is locked by a user, then no other user would be able to modify this security policy.
* `sequence_number` - This field is used to resolve conflicts between multiple policies that have rules that match the same packet.
* `comments` - Comments for security policy lock/unlock.
* `tag` - A list of scope + tag pairs to associate with this policy.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
