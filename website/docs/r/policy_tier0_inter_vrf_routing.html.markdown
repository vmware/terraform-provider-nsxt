---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_tier0_inter_vrf_routing"
description: A resource to configure tier0 inter VRF routing

---

# nsxt_policy_tier0_inter_vrf_routing

This resource provides a method for the management of tier0 inter VRF routing.
This resource is supported with NSX 4.1.0 onwards.

## Example Usage

```hcl
resource "nsxt_policy_tier0_inter_vrf_routing" "test" {
  display_name = "test-vrf-route"
  gateway_path = nsxt_policy_tier0_gateway.t0.path
  target_path  = nsxt_policy_tier0_gateway.t0-vrf.path

  bgp_route_leaking {
    address_family = "IPV4"
  }
  static_route_advertisement {
    advertisement_rule {
      name                      = "test"
      action                    = "PERMIT"
      prefix_operator           = "GE"
      route_advertisement_types = ["TIER0_CONNECTED", "TIER0_NAT"]
      subnets                   = ["192.168.240.0/24", "192.168.241.0/24"]
    }
    in_filter_prefix_list = [
      join("/", [nsxt_policy_tier0_gateway.parent.path, "prefix-lists/IPSEC_LOCAL_IP"]),
      join("/", [nsxt_policy_tier0_gateway.parent.path, "prefix-lists/DNS_FORWARDER_IP"])
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `gateway_path` - (Required) The NSX Policy path to the Tier0 Gateway for this resource.
* `bgp_route_leaking` - (Optional) Import / export BGP routes
  * `address_family` - (Optional) Address family type. Valid values are "IPV4", "IPV6".
  * `in_filter` - (Optional) List of route map paths for IN direction.
  * `out_filter` - (Optional) List of route map paths for OUT direction.
* `static_route_advertisement` - (Optional) Advertise subnet to target peers as static routes.
  * `advertisement_rule` - (Optional) Route advertisement rules.
    * `action` - (Optional) Action to advertise routes. Valid values are "PERMIT", "DENY". Default is "PERMIT".
    * `name` - (Optional) Display name for rule.
    * `prefix_operator` - (Optional) Prefix operator to match subnets. Valid values are "GE", "EQ". Default is "GE".
    * `route_advertisement_types` - (Optional) Enable different types of route advertisements. Valid values are "TIER0_STATIC", "TIER0_CONNECTED", "TIER0_NAT", "TIER0_DNS_FORWARDER_IP", "TIER0_IPSEC_LOCAL_ENDPOINT", "TIER1_STATIC", "TIER1_CONNECTED", "TIER1_LB_SNAT", "TIER1_LB_VIP", "TIER1_NAT", "TIER1_DNS_FORWARDER_IP", "TIER1_IPSEC_LOCAL_ENDPOINT".
    * `subnets` - (Optional) Network CIDRs.
  * `in_filter_prefix_list` - (Optional) Paths of ordered Prefix list.
* `target_path` - (Required) Policy path to tier0/vrf belongs to the same parent tier0.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing tier0 inter VRF routing configuration can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_tier0_inter_vrf_routing.vrfroute POLICY_PATH
```
The above command imports the policy tier0 inter VRF routing configuration named `vrfroute` for policy path `POLICY_PATH`.
