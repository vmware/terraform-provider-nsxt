---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_upgrade_run"
description: A resource to configure and execute upgrade of NSXT cluster.
---

# nsxt_upgrade_run

This resource provides a method to configure and execute upgrade of NSXT edges, hosts, and managers.

It will first configure upgrade unit groups and upgrade plan settings for `EDGE`
and `HOST` components. Then it will execute the upgrade for each component. If 
there are either `EDGE` or `HOST` groups disabled from the upgrade, the corresponding 
component will be paused after enabled groups get upgraded. At the same time, `MP` upgrade
won't be processed, because that requires `EDGE` and `HOST` to be fully upgraded. Refer to
the exposed attribute for current upgrade state details. For example, check `upgrade_plan`
for current upgrade plan, which also includes the plan not specified in this resource and
`state` for upgrade status of each component and UpgradeUnitGroups in the component. For more
details, please check NSX admin guide.

If upgrade post-checks are configured to be run, it will trigger the upgrade post-check.
Please use data source `nsxt_upgrade_postcheck` to retrieve results of upgrade post-checks.

## Example Usage

```hcl
resource "nsxt_upgrade_run" "run1" {
  upgrade_prepare_ready_id = nsxt_upgrade_prepare_ready.test.id

  edge_group {
    id      = data.nsxt_edge_upgrade_group.eg1.id
    enabled = false
  }

  edge_group {
    id      = data.nsxt_edge_upgrade_group.eg2.id
    enabled = true
  }

  host_group {
    id       = data.nsxt_host_upgrade_group.hg1.id
    parallel = true
  }
  host_group {
    display_name = "TEST123"
    parallel     = false
    hosts        = ["2fa96cdc-6b82-4284-a69a-18a21a6b6d0c"]
  }

  edge_upgrade_setting {
    parallel           = true
    post_upgrade_check = true
  }

  host_upgrade_setting {
    parallel           = true
    post_upgrade_check = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `upgrade_prepare_ready_id` - (Required) ID of corresponding `nsxt_upgrade_prepare_ready` resource. Updating this field will trigger replacement (destroy and create) of this resource.
* `edge_group` - (Optional) EDGE component upgrade unit group configurations. Groups will be reordered following the order they present in this field.
    * `id` - (Required) ID of the upgrade unit group.
    * `enabled` - (Optional) Flag to indicate whether upgrade of this group is enabled or not. Default: True.
    * `parallel` - (Optional) Upgrade method to specify whether upgrades of UpgradeUnits in this group are performed in parallel or serially. Default: True.
    * `pause_after_each_upgrade_unit` - (Optional) Flag to indicate whether upgrade should be paused after upgrade of each upgrade-unit. Default: False.
* `host_group` - (Optional) HOST component upgrade unit group configurations. Groups will be reordered following the order they present in this field.
    * `id` - (Optional) ID of the upgrade unit group. Should exist only for predefined groups. When creating a custom host group, the value is assigned by NSX.
    * `display_name` - (Optional) The display name of the host group. Should be assigned only for custom host groups and must be unique.
    * `enabled` - (Optional) Flag to indicate whether upgrade of this group is enabled or not. Default: True.
    * `parallel` - (Optional) Upgrade method to specify whether upgrades of UpgradeUnits in this group are performed in parallel or serially. Default: True.
    * `pause_after_each_upgrade_unit` - (Optional) Flag to indicate whether upgrade should be paused after upgrade of each upgrade-unit. Default: False.
    * `upgrade_mode` - (Optional) Upgrade mode. Supported values: `maintenance_mode`, `in_place`, `stage_in_vlcm`.
    * `maintenance_mode_config_vsan_mode` - (Optional) Maintenance mode config of vsan mode. Supported values: `evacuate_all_data`, `ensure_object_accessibility`, `no_action`.
    * `maintenance_mode_config_evacuate_powered_off_vms` - (Optional) Maintenance mode config of whether evacuate powered off vms.
    * `rebootless_upgrade` - (Optional) Flag to indicate whether to use rebootless upgrade. Default: True.
    * `hosts` - (Optional) The list of hosts to be associated with a custom group.
* `edge_upgrade_setting` - (Optional) EDGE component upgrade plan setting.
    * `parallel` - (Optional) Upgrade Method to specify whether upgrades of UpgradeUnitGroups in this component are performed serially or in parallel. Default: True.
    * `post_upgrade_check` - (Optional) Flag to indicate whether run post upgrade check after upgrade. Default: True.
* `host_upgrade_setting` - (Optional) HOST component upgrade plan setting.
    * `parallel` - (Optional) Upgrade Method to specify whether upgrades of UpgradeUnitGroups in this component are performed serially or in parallel. Default: True.
    * `post_upgrade_check` - (Optional) Flag to indicate whether run post upgrade check after upgrade. Default: True.
    * `stop_on_error` - (Optional) Flag to indicate whether to pause the upgrade plan execution when an error occurs. Default: False.
* `finalize_upgrade_setting` - (Optional) FINALIZE_UPGRADE component upgrade plan setting.
    * `enabled` - (Optional) Finalize upgrade after completion of all the components' upgrade is complete. Default: True.   

## Argument Reference

In addition to arguments listed above, the following attributes are exported:

* `upgrade_plan` - (Computed) Upgrade plan for current upgrade. Upgrade unit groups that are not defined in `edge_group` or `host_group` will also be included here.
    * `type` - Component type.
    * `id` - ID of the upgrade unit group.
    * `enabled` - Flag to indicate whether upgrade of this group is enabled or not.
    * `parallel` - Upgrade method to specify whether the upgrade is to be performed in parallel or serially.
    * `pause_after_each_upgrade_unit` - Flag to indicate whether upgrade should be paused after upgrade of each upgrade-unit.
    * `extended_config` - Extended configuration for the group.
* `state` - (Computed) Upgrade states of each component
    * `type` - Component type.
    * `status` - Upgrade status of component.
    * `details` - Details about the upgrade status.
    * `target_version` - Target component version
    * `group_state` - State of upgrade group
       * `group_id` - Upgrade group ID
       * `group_name` - Upgrade group name
       * `status` - Upgrade status of the upgrade group

## Importing

Importing is not supported for this resource.
