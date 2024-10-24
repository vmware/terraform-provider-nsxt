---
layout: "nsxt"
page_title: "VMware NSX Terraform Provider NSX upgrade support"
description: |-
  VMware NSX Terraform Provider NSX upgrade support
---

# NSX Provider for NSX Upgrade

NSX Terraform Provider offers support for [upgrading NSX components](https://docs.vmware.com/en/VMware-NSX/4.1/upgrade/GUID-E04242D7-EF09-4601-8906-3FA77FBB06BD.html).

The provider support consists of several resources and data sources which can perform a complete upgrade cycle.
The NSX Terraform provider support upgrades for NSX on ESX hosts, NSX Edge appliances, and NSX manager appliances.

**NOTE:** The upgrade process using the Terraform NSXT provider assumes that the Terraform NSX provider uses NSX upgrade 
components exclusively. Therefore, concurrent changess to upgrade components via other means, e.g. UI, might fail the process.
**NOTE:** The NSX upgrade process should start each upgrade with a clean state. So whenever an upgrade is executed, e.g. 
from v3.2.3 to v4.1.1 the Terraform state file should be retained until the upgrade is completed. Then, if the same plan 
is used later on to upgrade from v4.1.1 to v4.2.0, the Terraform state file should be deleted prior to applying the plan.

## Upgrade process steps

### Preparation for upgrade

Upgrade preparation consists of several actions, and is performed using [nsxt_upgrade_prepare](../resources/upgrade_prepare.html.markdown):
* Uploading the upgrade and precheck bundles.
* Acceptance of the user agreement.
* Reporting of failed pre-checks.

**NOTE:** Only one `nsxt_upgrade_prepare` resource should be used in the upgrade plan.

```hcl
resource "nsxt_upgrade_prepare" "prepare_res" {
  upgrade_bundle_url    = "http://url-to-upgrade-bundle.mub"
  precheck_bundle_url   = "http://url-to-precheck-bundle.pub"
  accept_user_agreement = true
}
```

**NOTE:** Acceptance of the license agreement is required to initiate the upgrade process.
**NOTE:** Precheck bundle is usable only while upgrading a NSX manager of version v4.1.1 and above.

### Checking the readiness for upgrade

When preparation is complete, it is required to validate the readiness for upgrade using the 
[nsxt_upgrade_prepare_ready](../data-sources/upgrade_prepare_ready.html.markdown) data source.
When NSX preparation results with failures or warnings, the data source will fail the upgrade.
Warnings should be acknowledged as described below.

For example:
```
Error: 
There are unacknowledged warnings in prechecks:
  Component: EDGE, code: pUBCheck, description: Checks if new version of pub available on VMware Download site.
  Component: MP, code: backupOperationCheck, description: Warns if backup is not taken recently and blocks management plane upgrade if backup is in progress

Please address these errors from NSX or using nsxt_upgrade_precheck_acknowledge resource
```

Prechecks which are in `failure' or in warning state can also be obtained from the `nsxt_upgrade_prepare` resource, 
using its `failed_prechecks` attribute.

### Acknowledgement of upgrade issues and warnings

The execution of `nsxt_upgrade_prepare` and `nsxt_upgrade_prepare_ready` could result with various warnings which would 
halt the upgrade process. However, warnings can be suppressed with the [nsxt_upgrade_precheck_acknowledge](../resources/upgrade_precheck_acknowledge.html.markdown) resource.

The resource uses the `precheck_ids` to conclude which prechecks are safe to ignore for this upgrade cycle.
The `nsxt_upgrade_precheck_acknowledge` resource requires the `target_version` from the `nsxt_upgrade_prepare` resource 
to make sure that the acknowledged warnings are from the correct execution, and to enforce execution of the 
`nsxt_upgrade_precheck_acknowledge` resource after preparation is complete.

The `nsxt_upgrade_prepare_ready` data source verifies that the preparation phase of the upgrade was concluded successfully 
and enable the execution of the upgrade. In case that there are still open issues, these will raise as described above.

```hcl
resource "nsxt_upgrade_precheck_acknowledge" "test" {
  precheck_ids   = ["backupOperationCheck", "pUBCheck"]
  target_version = nsxt_upgrade_prepare.prepare_res.target_version
}

data "nsxt_upgrade_prepare_ready" "ready" {
  upgrade_prepare_id = nsxt_upgrade_prepare.prepare_res.id
  depends_on         = [nsxt_upgrade_precheck_acknowledge.test]
}
```

### Running the upgrade

In order to configure and execute upgrade of NSXT edges, hosts, and managers, use [nsxt_upgrade_run](../resources/upgrade_run.html.markdown) resource. 

This resource will start with configuration of upgrade unit groups and upgrade plan settings for `EDGE`and `HOST` 
components. Host component groups may be created within this resource with custom host lists and attributes, or it can 
use the predefined groups which NSX provides.

Following the configuration of the edge, host groups, it will execute the upgrade for each component. If there are 
either `EDGE` or `HOST` groups excluded from the upgrade, the corresponding component will be paused after enabled 
groups get upgraded. At the same time, `MP` upgrade won't be processed, because that requires `EDGE` and `HOST` to be 
fully upgraded. Refer to the exposed attribute for current upgrade state details. For example, check `upgrade_plan`
for current upgrade plan, which also includes the plan not specified in this resource and `state` for upgrade status of 
each component and UpgradeUnitGroups in the component. For more details, please check NSX admin guide.

The order of the groups within the `nsxt_upgrade_run` resource defined the order of the upgrade execution.

In the example below, host_group `HostWithDifferentSettings` will use different settings than the defaulted ones. The 
group configuration is initiated with the `nsxt_upgrade_run` while the other groups in this example are predefined and 
consumed using the data sources `nsxt_edge_upgrade_group` and `nsxt_host_upgrade_group`.

**NOTE:** Only one `nsxt_upgrade_run` resource should be used in the upgrade plan.

```hcl
data "nsxt_policy_host_transport_node" "host_transport_node1" {
  display_name = "HostTransportNode1"
}

data "nsxt_policy_host_transport_node" "host_transport_node2" {
  display_name = "HostTransportNode2"
}

data "nsxt_policy_host_transport_node" "host_transport_node3" {
  display_name = "HostTransportNode3"
}

data "nsxt_host_upgrade_group" "hg1" {
  upgrade_prepare_id = nsxt_upgrade_prepare.test.id
  display_name       = "Group for Compute-Cluster1"
}

data "nsxt_edge_upgrade_group" "eg1" {
  upgrade_prepare_id = nsxt_upgrade_prepare.test.id
  display_name       = "edgegroup-EDGECLUSTER1"
}

data "nsxt_edge_upgrade_group" "eg2" {
  upgrade_prepare_id = nsxt_upgrade_prepare.test.id
  display_name       = "edgegroup-EDGECLUSTER2"
}

resource "nsxt_upgrade_run" "run" {
  upgrade_prepare_ready_id = data.nsxt_upgrade_prepare_ready.ready.id

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
    display_name = "HostWithDifferentSettings"
    parallel     = false
    hosts        = [data.nsxt_policy_host_transport_node.host_transport_node1.id]
  }

  host_group {
    display_name       = "OtherHostsWithDifferentSettings"
    parallel           = true
    rebootless_upgrade = false
    hosts = [
      data.nsxt_policy_host_transport_node.host_transport_node2.id,
    data.nsxt_policy_host_transport_node.host_transport_node3.id]
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

### Upgrade run stage with NSX v9.0 and above

NSX v9.0 introduces a few changes to the flow of execution of the run resource: the upgrade begins with updating the NSX 
management plane, then Edge appliances are upgraded, followed by ESXi hosts. Finally, the upgrade process is finalized, 
when the entire process is complete, including the upgrade of the ESXi software.

When the ESXi software is not upgraded, the finalization stage will fail. There is an additional component in the 
`nsxt_upgrade_run` resource which allows execution without the finalization stage. This allows successful execution of 
the NSX management plane, the Edge appliances and the NSX bits on the ESXi hosts, and postpone the ESXi OS upgrade to 
a later time.

### Post upgrade checks

Upgrade post check data sources can be used to examine the results of the edge and host upgrades, to conclude if the 
upgrade has completed successfully.

The example below uses Terraform `check` to test the upgrade results.

```hcl
data "nsxt_upgrade_postcheck" "pc1" {
  upgrade_run_id = nsxt_upgrade_run.run.id
  type           = "EDGE"
}

data "nsxt_upgrade_postcheck" "pc2" {
  upgrade_run_id = nsxt_upgrade_run.run.id
  type           = "HOST"
}

check "edge_post_check" {
  assert {
    condition     = data.nsxt_upgrade_postcheck.pc1.failed_group == null
    error_message = "Edge group upgrade check failed"
  }
}

check "host_post_check" {
  assert {
    condition     = data.nsxt_upgrade_postcheck.pc2.failed_group == null
    error_message = "Host group upgrade check failed"
  }
}
```
