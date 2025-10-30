---
subcategory: "Segments"
page_title: "NSXT: nsxt_policy_segment_port"
description: A resource to configure a network Segment Port.
---

# nsxt_policy_segment_port

This resource provides a method for the management of Segments Port.

## Example Usage

```hcl
resource "nsxt_policy_segment_port" "sample" {
  display_name = "segment-port1"
  description  = "NSX-t Segment port"
  segment_path = data.nsxt_policy_segment.segment1.path
  discovery_profile {
    ip_discovery_profile_path  = data.nsxt_policy_ip_discovery_profile.segprofile.path
    mac_discovery_profile_path = data.nsxt_policy_mac_discovery_profile.segprofile.path
  }
  security_profile {
    spoofguard_profile_path = data.nsxt_policy_spoofguard_profile.segprofile.path
    security_profile_path   = data.nsxt_policy_segment_security_profile.segprofile.path
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this policy.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `segment_path` - (Optional) Path of the segment port.
* `attachment` - (Optional) VIF attachment.
* `discovery_profile` - (Optional) IP and MAC discovery profiles for this segment port.
* `qos_profile` - (Optional) QoS profiles for this segment port.
* `security_profile` - (Optional) Security profiles for this segment port.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Segment Port.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing segment port can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_segment_port.segmentport1 GW_ID/SEGMENT_ID/ID
```

The above command imports the segment port named `segmentport1` with the NSX Gateway ID `GW_ID`, Segment ID `SEGMENT_ID` and Port ID `ID`.

```shell
terraform import nsxt_policy_segment_port.segmentport1 POLICY_PATH
```

The above command imports the segment port named `segmentport1` with the policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
