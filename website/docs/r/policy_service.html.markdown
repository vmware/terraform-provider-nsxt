---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_service"
sidebar_current: "docs-nsxt-resource-policy-service"
description: A resource that can be used to configure a networking and security service in NSX Policy.
---

# nsxt_policy_service

This resource provides a way to configure a networking and security service which can be used within NSX Policy.

## Example Usage

```hcl
resource "nsxt_policy_service" "service_icmp" {
  description  = "ICMP service provisioned by Terraform"
  display_name = "S1"

  icmp_entry {
    display_name = "ICMP-entry"
    description  = "ICMP v4 entry"
    protocol     = "ICMPv4"
    icmp_code    = "1"
    icmp_type    = "3"
  }

  tag {
    scope = "color"
    tag   = "blue"
  }
}

resource "nsxt_policy_service" "service_l4port" {
  description  = "L4 ports service provisioned by Terraform"
  display_name = "S1"

  l4_port_set_entry {
    display_name      = "TCP80"
    description       = "TCP port 80 entry"
    protocol          = "TCP"
    destination_ports = [ "80" ]
  }

  tag {
    scope = "color"
    tag   = "pink"
  }
}


```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
The service must contain at least 1 entry (of at least one of the types), and possibly more.
* `icmp_entry` - (Optional) Set of ICMP type service entries. Each with the following attributes:
    * `display_name` - (Optional) Display name of the service entry.
    * `description` - (Optional) Description of the service entry.
    * `protocol` - (Required) Version of ICMP protocol ICMPv4 or ICMPv6.
    * `icmp_code` - (Optional) ICMP message code.
    * `icmp_type` - (Optional) ICMP message type.
* `l4_port_set_entry` - (Optional) Set of L4 ports set service entries. Each with the following attributes:
    * `display_name` - (Optional) Display name of the service entry.
    * `description` - (Optional) Description of the service entry.
    * `protocol` - (Required) L4 protocol. Accepted values - 'TCP' or 'UDP'.
    * `destination_ports` - (Optional) Set of destination ports.
    * `source_ports` - (Optional) Set of source ports.
* `igmp_entry` - (Optional) Set of IGMP type service entries. Each with the following attributes:
    * `display_name` - (Optional) Display name of the service entry.
    * `description` - (Optional) Description of the service entry.
* `ether_type_entry` - (Optional) Set of Ether type service entries. Each with the following attributes:
    * `display_name` - (Optional) Display name of the service entry.
    * `description` - (Optional) Description of the service entry.
    * `ether_type` - (Required) Type of the encapsulated protocol.
* `ip_protocol_entry` - (Optional) Set of IP Protocol type service entries. Each with the following attributes:
    * `display_name` - (Optional) Display name of the service entry.
    * `description` - (Optional) Description of the service entry.
    * `protocol` - (Required) IP protocol number.
* `algorithm_entry` - (Optional) Set of Algorithm type service entries. Each with the following attributes:
    * `display_name` - (Optional) Display name of the service entry.
    * `description` - (Optional) Description of the service entry.
    * `destination_port` - (Required) a single destination port.
    * `source_ports` - (Optional) Set of source ports/ranges.
    * `algorithm` - (Required) Algorithm one of "ORACLE_TNS", "FTP", "SUN_RPC_TCP", "SUN_RPC_UDP", "MS_RPC_TCP", "MS_RPC_UDP", "NBNS_BROADCAST", "NBDG_BROADCAST", "TFTP"

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the service.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing service can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_service.service_icmp ID
```

The above service imports the service named `service_icmp` with the NSX ID `ID`.
