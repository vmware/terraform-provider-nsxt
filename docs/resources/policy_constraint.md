---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_constraint"
description: A resource to configure a Constraint (Quota).
---

# nsxt_policy_constraint

This resource provides a method for the management of a Constraint.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.

## Example Usage

```hcl
resource "nsxt_policy_constraint" "test" {
  display_name = "demo-quota"
  description  = "Terraform provisioned Constraint"
  message      = "too many objects mate"

  target {
    path_prefix = "/orgs/default/projects/demo"
  }

  instance_count {
    count                = 4
    target_resource_type = "StaticRoutes"
  }

  instance_count {
    count                = 1
    target_resource_type = "Infra.Tier1.PolicyDnsForwarder"
  }

  instance_count {
    count                = 20
    target_resource_type = "Infra.Domain.Group"
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
resource "nsxt_policy_constraint" "test" {
  context {
    project_id = "demo"
  }

  display_name = "demo1-quota"

  target {
    path_prefix = "/orgs/default/projects/demo/vpcs/demo1"
  }

  instance_count {
    count                = 4
    target_resource_type = "Org.Project.Vpc.PolicyNat.PolicyVpcNatRule"
  }
}
```

## Argument Reference

The following arguments are supported:

* `context` - (Optional) The context which the object belongs to
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `message` - (Optional) User friendly message to be shown to users upon violation.
* `target` - (Optional) Targets for the constraints to be enforced
    * `path_prefix` - (Optional) Prefix match to the path
* `instance_count` - (Optional) Constraint details
    * `target_resource_type` - (Required) Type of the resource that should be limited in instance count (refer to the table below)
    * `operator` - (Optional) Either `<=` or `<`. Default is `<=`
    * `count` - (Required) Limit of instances
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Target resource types

|Object|project + VPC|project only|VPC only|
|------|-------------|------------|--------|
|Group|Group|Infra.Domain.Group|Org.Project.Vpc.Group|
|Service||Infra.Service||
|Service Entry||Infra.Service.ServiceEntry||
|TLS Certificate||Infra.TlsCertificate||
|TLS CRL||Infra.TlsCrl||
|All Firewall Rules|Rule|||
|Security Policy|SecurityPolicy|Infra.Domain.SecurityPolicy|Org.Project.Vpc.SecurityPolicy|
|Security Policy Rule|SecurityPolicy.Rule|Infra.Domain.SecurityPolicy.Rule|Org.Project.Vpc.SecurityPolicy.Rule|
|Gateway Policy|SecurityPolicy|Infra.Domain.SecurityPolicy|Org.Project.Vpc.SecurityPolicy|
|Gateway Policy Rule|GatewayPolicy.Rule|Infra.Domain.GatewayPolicy.Rule|Org.Project.Vpc.GatewayPolicy.Rule|
|IDS Security Policy||Infra.Domain.IdsPolicy||
|IDS Security Policy Rule||Infra.Domain.IdsPolicy.Rule||
|Session Timer Profile||Infra.PolicyFirewallSessionTimerProfile||
|Flood Protection Profile||Infra.FloodProtectionProfile||
|DNS Security Profile||Infra.DnsSecurityProfile||
|Context Profile||Infra.PolicyContextProfile||
|l7 Access Profile||Infra.L7AccessProfile||
|Tier1 Gateway||Infra.Tier1||
|Segment||Infra.Segment||
|Segment Port||Infra.Segment.SegmentPort||
|Subnet|||Org.Project.Vpc.Subnet|
|Subnet Port|||Org.Project.Vpc.Subnet.SubnetPort|
|Segment Security Profile||Infra.SegmentSecurityProfile||
|Segment QoS Profile||Infra.QosProfile||
|Segment IP Discovery Profile||Infra.IpDiscoveryProfile||
|Segment MAC Discovery Profile||Infra.MacDiscoveryProfile||
|Segment Spoof Guard Profile||Infra.SpoofGuardProfile||
|IPv6 NDRA Profile||Infra.Ipv6NdraProfile||
|IPv6 DAD Profile||Infra.Ipv6DadProfile||
|Gateway QoS Profile||Infra.GatewayQosProfile||
|Static Routes|StaticRoutes|Infra.Tier1.StaticRoutes|Org.Project.Vpc.StaticRoutes|
|NAT Rule|NatRule|Infra.Tier1.PolicyNat.PolicyNatRule|Org.Project.Vpc.PolicyNat.PolicyNatRule|
|DNS Forwarder Zone||Infra.PolicyDnsForwarderZone||
|DNS Forwarder||Infra.Tier1.PolicyDnsForwarder||
|IP Address Block||Infra.IpAddressBlock||
|IP Address Pool||Infra.IpAddressPool||
|IP Address Pool Subnet||Infra.IpAddressPool.IpAddressPoolSubnet||
|IP Address Allocation||Infra.IpAddressPool.IpAddressAllocation||
|DHCP Server Config||Infra.DhcpServerConfig||
|IPSec VPN Service||Infra.Tier1.IPSecVpnService||
|IPSec VPN Session||Infra.Tier1.IPSecVpnService.IPSecVpnSession||
|IPSec VPN Local Endpoint||Infra.Tier1.IPSecVpnService.IPSecVpnLocalEndpoint||
|IPSec VPN Tunnel Profile||Infra.IPSecVpnTunnelProfile||
|IPSec VPN IKE Profile||Infra.IPSecVpnIkeProfile||
|IPSec VPN DPD Profile||Infra.IPSecVpnDpdProfile||
|L2 VPN Service||Infra.Tier1.L2VpnService||
|L2 VPN Session||Infra.Tier1.L2VpnService.L2VpnSession||
|VPC||Org.Project.Vpc||

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_constraint.test PATH
```

The above command imports Constraint named `test` with the NSX path `PATH`.
