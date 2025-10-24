---
page_title: "Support for NSX Upgrades"
description: |-
  Support for NSX Upgrades
---

# Support for NSX Upgrades

The provider includes support for upgrading [NSX components](https://techdocs.broadcom.com/us/en/vmware-cis/nsx/vmware-nsx/4-2/upgrade-guide.html).

The provider support consists of several resources and data sources which can perform a complete upgrade cycle.

The provider supports upgrades for NSX on ESX hosts, NSX Edge appliances, and NSX Manager appliances.

~> **NOTE:** The upgrade process using the provider assumes that only the provider is used to upgrade NSX components.
Therefore, concurrent changes to upgrade components via other means, e.g. UI, might fail the process.

~> **NOTE:** The NSX upgrade process should start each upgrade with a clean state. So whenever an upgrade is executed, e.g.
from v3.2.3 to v4.1.1 the Terraform state file should be retained until the upgrade is completed. Then, if the same plan
is used later on to upgrade from v4.1.1 to v4.2.0, the Terraform state file should be deleted prior to applying the plan.

**NOTE:** For details concerning specific NSX version, and limitations concerning supported scenarios and limitations, please review the [NSX upgrade guide](https://techdocs.broadcom.com/us/en/vmware-cis/nsx/vmware-nsx/4-2/upgrade-guide/nsx-t-upgrade-guide.html).

## Upgrade process steps

### Preparation for upgrade

Upgrade preparation consists of several actions, and is performed using [nsxt_upgrade_prepare](../resources/upgrade_prepare.md):

* Uploading the upgrade and precheck bundles.
* Acceptance of the user agreement.
* Reporting of failed pre-checks.

~> **NOTE:** Only one `nsxt_upgrade_prepare` resource should be used in the upgrade plan.

```hcl
resource "nsxt_upgrade_prepare" "prepare_res" {
  upgrade_bundle_url    = "http://url-to-upgrade-bundle.mub"
  precheck_bundle_url   = "http://url-to-precheck-bundle.pub"
  accept_user_agreement = true
}
```

~> **NOTE:** Acceptance of the license agreement is required to initiate the upgrade process.

~> **NOTE:** Precheck bundle is usable only while upgrading a NSX manager of version v4.1.1 and above.

### Checking the readiness for upgrade

When preparation is complete, it is required to validate the readiness for upgrade using the
[nsxt_upgrade_prepare_ready](../data-sources/upgrade_prepare_ready.md) data source.
When NSX preparation results with failures or warnings, the data source will fail the upgrade.
Warnings should be acknowledged as described below.

For example:

```shell
Error:
There are unacknowledged warnings in prechecks:
  Component: EDGE, code: pUBCheck, description: Checks if new version of pub available on VMware Download site.
  Component: MP, code: backupOperationCheck, description: Warns if backup is not taken recently and blocks management plane upgrade if backup is in progress

Please address these errors from NSX or using nsxt_upgrade_precheck_acknowledge resource
```

Prechecks which are in `failure' or in warning state can also be obtained from the`nsxt_upgrade_prepare` resource,
using its `failed_prechecks` attribute.

### Acknowledgement of upgrade issues and warnings

The execution of `nsxt_upgrade_prepare` and `nsxt_upgrade_prepare_ready` could result with various warnings which would
halt the upgrade process. However, warnings can be suppressed with the [nsxt_upgrade_precheck_acknowledge](../resources/upgrade_precheck_acknowledge.md) resource.

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

In order to configure and execute upgrade of NSXT edges, hosts, and managers, use [nsxt_upgrade_run](../resources/upgrade_run.md) resource.

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

~> **NOTE:** Only one `nsxt_upgrade_run` resource should be used in the upgrade plan.

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

### Upgrading Global Manager

NSX Global Manager upgrade is similar to Local Manager. Yet, Global Manager has no Edge or Host components and therefore
Doesn't allow the configuration of these components in upgrade groups in the `nsxt_upgrade_run` resource.

The example below uses Terraform to upgrade a Global Manager.

```hcl
provider "nsxt" {
  host                 = "global.manager.somedomain.org"
  username             = "admin"
  password             = "AdminPassword"
  allow_unverified_ssl = true
  global_manager       = true // This is required to indicate that we upgrade a Global Manager
}

resource "nsxt_upgrade_prepare" "gm_prepare_res" {
  upgrade_bundle_url    = var.upgrade_bundle_url
  accept_user_agreement = true
}

resource "nsxt_upgrade_precheck_acknowledge" "gm_precheck_ack" {
  provider = nsxt.gm_nsxt

  precheck_ids   = var.gm_precheck_warns
  target_version = nsxt_upgrade_prepare.gm_prepare_res.target_version
}

data "nsxt_upgrade_prepare_ready" "gm_ready" {
  provider = nsxt.gm_nsxt

  upgrade_prepare_id = nsxt_upgrade_prepare.gm_prepare_res.id
  depends_on         = [nsxt_upgrade_precheck_acknowledge.gm_precheck_ack]
}

resource "nsxt_upgrade_run" "gm_run" {
  provider                 = nsxt.gm_nsxt
  upgrade_prepare_ready_id = data.nsxt_upgrade_prepare_ready.gm_ready.id
}
```

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
