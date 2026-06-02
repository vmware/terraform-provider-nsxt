---
subcategory: "Beta"
page_title: "NSXT: policy_baremetal_server_tags"
description: A resource for applying tags to Bare Metal Servers.
---

# nsxt_policy_baremetal_server_tags

This resource provides a way to configure tags on Bare Metal Servers discovered and managed by NSX-T.
Tags applied to bare metal servers can be used in dynamic group membership criteria and security policies.

## Example Usage

```hcl
# Discover bare metal servers
data "nsxt_policy_baremetal_servers" "production_servers" {
  # Filter for production servers if needed
}

# Apply tags to a specific bare metal server
resource "nsxt_policy_baremetal_server_tags" "web_server_tags" {
  external_id = data.nsxt_policy_baremetal_servers.production_servers.results[0].external_id

  tag {
    scope = "environment"
    tag   = "production"
  }

  tag {
    scope = "application"
    tag   = "web-server"
  }

  tag {
    scope = "owner"
    tag   = "platform-team"
  }
}

# Use tagged servers in a dynamic group
resource "nsxt_policy_group" "production_web_servers" {
  display_name = "Production-Web-Servers"
  description  = "Dynamic group of production web servers"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "environment|production"
    }
  }

  conjunction {
    operator = "AND"
  }

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "application|web-server"
    }
  }
}

# Apply security policy to tagged servers
resource "nsxt_policy_security_policy" "web_server_security" {
  display_name = "Web-Server-Security"
  description  = "Security policy for web servers"
  category     = "Application"

  rule {
    display_name       = "Allow-HTTP-HTTPS"
    source_groups      = ["ANY"]
    destination_groups = [nsxt_policy_group.production_web_servers.path]
    action             = "ALLOW"
    services           = ["HTTP", "HTTPS"]
    logged             = true
  }
}
```

## Example with Multiple Servers

```hcl
# Tag multiple servers using for_each
resource "nsxt_policy_baremetal_server_tags" "production_servers" {
  for_each = {
    for idx, server in data.nsxt_policy_baremetal_servers.all.results :
    idx => server.external_id
    if length(regexall("prod", lower(server.display_name))) > 0
  }

  external_id = each.value

  tag {
    scope = "environment"
    tag   = "production"
  }

  tag {
    scope = "managed-by"
    tag   = "terraform"
  }
}
```

## Argument Reference

* `external_id` - (Required) External ID of the bare metal server to tag. This is typically obtained from the `nsxt_policy_baremetal_servers` data source.
* `tag` - (Optional) A list of scope + tag pairs to associate with this bare metal server. Default is empty (no tags applied).

The `tag` block supports:

* `scope` - (Required) Tag scope.
* `tag` - (Required) Tag value.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the bare metal server being tagged.

## Import

An NSX Policy Bare Metal Server Tags resource can be imported using the bare metal server external ID:

```shell
terraform import nsxt_policy_baremetal_server_tags.web_server_tags 71be0142-2ed1-1d53-9c60-5564cf4b7e2e
```

The above command imports the bare metal server tags for the server with external ID `71be0142-2ed1-1d53-9c60-5564cf4b7e2e`.

## Notes

* This resource requires NSX-T version 9.0.0 or higher (Bare Metal Server support)
* The bare metal server must be discovered and registered with NSX-T through a compute manager
* Tags applied through this resource will be visible in the NSX-T UI and available for use in policy groups
* Only tag operations are supported; other bare metal server properties cannot be modified through this resource
* This resource uses NSX Policy APIs and requires local manager access (not supported on Global Manager)
* When the resource is destroyed, all tags managed by this resource are removed from the bare metal server
