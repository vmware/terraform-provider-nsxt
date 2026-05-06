---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_parent_intrusion_service_policy"
description: A data source to retrieve Parent Intrusion Service (IDS) Policy metadata.
---

# nsxt_policy_parent_intrusion_service_policy

This data source provides metadata information about an existing Parent Intrusion Service (IDS) Policy configured on NSX for DFW (Distributed Firewall) context.
It can be useful for fetching policy path to use in `nsxt_policy_intrusion_service_policy_rule` resource.

~> **NOTE:** This data source retrieves only the parent policy metadata (id, display_name, description, path, sequence_number, etc.) **without embedded rules**, allowing you to refer a policy's path for creating standalone rules. For different use cases, consider:

* `nsxt_policy_intrusion_service_policy` - For **IDPS DFW policy with embedded rules**
* `nsxt_policy_intrusion_service_policy_rule` - For **individual standalone IDPS DFW rules**

This data source is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
# Get parent policy metadata for rule creation
data "nsxt_policy_parent_intrusion_service_policy" "parent_policy" {
  display_name = "production-ids-policy"
}

# Create standalone rule using parent policy path
resource "nsxt_policy_intrusion_service_policy_rule" "new_rule" {
  display_name = "new-detection-rule"
  policy_path  = data.nsxt_policy_parent_intrusion_service_policy.parent_policy.path
  action       = "DETECT"
  # ... other rule configuration
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_parent_intrusion_service_policy" "parent_policy" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "production-ids-policy"
}
```

## Argument Reference

* `id` - (Optional) The ID of the policy to retrieve.
* `display_name` - (Optional) The display name of the policy to retrieve.
* `domain` - (Optional) The domain of the policy. Defaults to `default`.
* `category` - (Optional) Category of the policy.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `stateful` - When it is stateful, the state of the network connects are tracked and a stateful packet inspection is performed. Note: Intrusion Service Policies are always stateful.
* `locked` - Indicates whether a security policy should be locked. If the security policy is locked by a user, then no other user would be able to modify this security policy.
* `sequence_number` - This field is used to resolve conflicts between multiple policies that have rules that match the same packet.
* `comments` - Comments for security policy lock/unlock.
* `tag` - A list of scope + tag pairs to associate with this policy.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
