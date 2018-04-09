---
layout: "nsxt"
page_title: "Provider: NSXT"
sidebar_current: "docs-nsxt-index"
description: |-
  The VMware NSX-T Terraform Provider
---

# The NSX Terraform Provider

The NSX Terraform provider gives the NSX administrator a way to automate NSX to
provide virtualized networking and security services using both ESXi and KVM
based hypervisor hosts as well as container networking and security.

More information on NSX can be found on the [NSX Product
Page](https://www.vmware.com/products/nsx.html)

Documentation on the NSX platform can be found on the [NSX Documentation
Page](https://docs.vmware.com/en/VMware-NSX-T/index.html)

Please use the navigation to the left to read about available data sources and
resources.

## Basic Configuration of the NSX Terraform Provider

In order to use the NSX Terraform provider you must first configure the
provider to communicate with the VMware NSX manager. The NSX manager is the
system which serves the NSX REST API and provides a way to configure the
desired state of the NSX system. The configuration of the NSX provider requires
the IP address, hostname, or FQDN of the NSX manager.

The NSX provider offers several ways to authenticate to the NSX manager.
Credentials can be provided statically or provided as environment variables. In
addition, client certificates can be used for authentication. For
authentication with certificates Terraform will require a certificate file and
private key file in PEM format. To use client certificates the client
certificate needs to be registered with NSX-T manager prior to invoking
Terraform.

The provider also can accept both signed and self-signed server certificates.
It is recommended that in production environments you only use certificates
signed by a certificate authority. NSX ships by default with a self-signed
server certificates as the hostname of the NSX manager is not known until the
NSX administrator determines what name or IP to use.

Setting the `allow_unverified_ssl` parameter to `true` will direct the
Terraform client to skip server certificate verification. This is not
recommended in production deployments as it is recommended that you use trusted
connection using certificates signed by a certificate authority.

With the `ca_file` parameter you can also specify a file that contains your
certificate authority certificate in PEM format to verify certificates with a
certificate authority.

There are also a number of other parameters that can be set to tune how the
provider connects to the NSX REST API. It is recommended you leave these to the
defaults unless you experience issues in which case they can be tuned to
optimize the system in your environment.

Note that in all of the examples you will need to update the `host`,
`username`, and `password` settings to match those configured in your NSX
deployment.

### Example of Configuration with Credentials

```hcl
provider "nsxt" {
  host                     = "192.168.110.41"
  username                 = "admin"
  password                 = "default"
  allow_unverified_ssl     = true
  max_retries              = 10
  retry_min_delay          = 500
  retry_max_delay          = 5000
  retry_on_status_codes    = [429]
}

```

### Example of Setting Environment Variables

```hcl
export NSXT_MANAGER_HOST="192.168.110.41"
export NSXT_USERNAME="admin"
export NSXT_PASSWORD="default"
```

### Example using a Client Certificate

```hcl
provider "nsxt" {
  host                  = "192.168.110.41"
  client_auth_cert_file = "mycert.pem"
  client_auth_key_file  = "mykey.pem"
  allow_unverified_ssl  = true
}

```

### Example with Certificate Authority Certificate

```hcl
provider "nsxt" {
  host     = "10.160.94.11"
  username = "admin"
  password = "qwerty"
  ca_file  = "myca.pem"
}

```

## Argument Reference

The following arguments are used to configure the VMware NSX-T Provider:

* `host` - (Required) The host name or IP address of the NSX-T manager. Can also
  be specified with the `NSXT_MANAGER_HOST` environment variable.
* `username` - (Required) The user name to connect to the NSX-T manager as. Can
  also be specified with the `NSXT_USERNAME` environment variable.
* `password` - (Required) The password for the NSX-T manager user. Can also be
  specified with the `NSXT_PASSWORD` environment variable.
* `client_auth_cert_file` - (Optional) The path to a certificate file for
  certificate authorization. Can also be specified with the
  `NSXT_CLIENT_AUTH_CERT_FILE` environment variable.
* `client_auth_key_file` - (Optional) The path to a private key file for the
  certificate supplied to `client_auth_cert_file`. Can also be specified with
  the `NSXT_CLIENT_AUTH_KEY_FILE` environment variable.
* `allow_unverified_ssl` - (Optional) Boolean that can be set to true to disable
  SSL certificate verification. This should be used with care as it could allow
  an attacker to intercept your auth token. If omitted, default value is
  `false`. Can also be specified with the `NSXT_ALLOW_UNVERIFIED_SSL`
  environment variable.
* `ca_file` - (Optional) The path to an optional CA certificate file for SSL
  validation. Can also be specified with the `NSXT_CA_FILE` environment
  variable.
* `max_retries` - (Optional) The maximum number of retires before failing an API
  request. Default: `10` Can also be specified with the `NSXT_MAX_RETRIES`
  environment variable.
* `retry_min_delay` - (Optional) The minimum delay, in milliseconds, between
  retires made to the API. Default:`500`. Can also be specified with the
  `NSXT_RETRY_MIN_DELAY` environment variable.
* `retry_max_delay` - (Optional) The maximum delay, in milliseconds, between
  retires made to the API. Default:`5000`. Can also be specified with the
  `NSXT_RETRY_MAX_DELAY` environment variable.
* `retry_on_status_codes` - (Optional) A list of HTTP status codes to retry on.
  By default, the provider will retry on HTTP error 429 (too many requests),
  essentially retrying on throttled connections. Can also be specified with the
  `NSXT_RETRY_ON_STATUS_CODES` environment variable.

## NSX Logical Networking

The NSX Terraform provider can be used to manage logical networking and
security constructs in NSX. This includes logical switching, routing and
firewall.

### Logical Networking and Security Example Usage

The following example demonstrates using the NSX Terraform provider to create a
logical switch and tier1 logical router. It then connects the logical switch to
the tier1 logical router and uplinks the T1 router to a pre-created T0 router.

#### Example variables.tf File

This file allows you to define some variables that can be reused in multiple
.tf files.

```hcl
variable "nsx_manager" {}
variable "nsx_username" {}
variable "nsx_password" {}
```
#### Example terraform.tfvars File

This file allows you to set some variables that can be reused in multiple .tf
files.

```hcl
nsx_manager = "192.168.110.41"

nsx_username = "admin"

nsx_password = "default"
```

#### Example nsx.tf file

This file will define the logical networking topology that Terraform will
create in NSX.

```hcl
#
# The first step is to configure the VMware NSX provider to connect to the NSX
# REST API running on the NSX manager.
#
provider "nsxt" {
  host                  = "${var.nsx_manager}"
  username              = "${var.nsx_username}"
  password              = "${var.nsx_password}"
  allow_unverified_ssl  = true
  max_retries           = 10
  retry_min_delay       = 500
  retry_max_delay       = 5000
  retry_on_status_codes = [429]
}

#
# Here we show that you define a NSX tag which can be used later to easily to
# search for the created objects in NSX.
#
variable "nsx_tag_scope" {
  default = "project"
}

variable "nsx_tag" {
  default = "terraform-demo"
}

#
# This part of the example shows some data sources we will need to refer to
# later in the .tf file. They include the transport zone, tier 0 router and
# edge cluster.
#
data "nsxt_transport_zone" "overlay_tz" {
  display_name = "tz1"
}

#
# The tier 0 router (T0) is considered a "provider" router that is pre-created
# by the NSX admin. A T0 router is used for north/south connectivity between
# the logical networking space and the physical networking space. Many tier 1
# routers will be connected to a tier 0 router.
#
data "nsxt_logical_tier0_router" "tier0_router" {
  display_name = "DefaultT0Router"
}

data "nsxt_edge_cluster" "edge_cluster1" {
  display_name = "EdgeCluster1"
}

#
# This shows the settings required to create a NSX logical switch to which you
# can attach virtual machines.
#
resource "nsxt_logical_switch" "switch1" {
  admin_state       = "UP"
  description       = "LS created by Terraform"
  display_name      = "TfLogicalSwitch"
  transport_zone_id = "${data.nsxt_transport_zone.overlay_tz.id}"
  replication_mode  = "MTEP"

  tag {
    scope = "${var.nsx_tag_scope}"
    tag   = "${var.nsx_tag}"
  }

  tag {
    scope = "tenant"
    tag   = "second_example_tag"
  }
}

#
# In this part of the example the settings are defined that are required to
# create a T1 router. In NSX a T1 router is often used on a per user, tenant,
# or application basis. Each application may have it's own T1 router. The T1
# router provides the default gateway for machines on logical switches
# connected to the T1 router.
#
resource "nsxt_logical_tier1_router" "tier1_router" {
  description                 = "Tier1 router provisioned by Terraform"
  display_name                = "TfTier1"
  failover_mode               = "PREEMPTIVE"
  high_availability_mode      = "ACTIVE_STANDBY"
  edge_cluster_id             = "${data.nsxt_edge_cluster.edge_cluster1.id}"
  enable_router_advertisement = true
  advertise_connected_routes  = true
  advertise_static_routes     = false
  advertise_nat_routes        = true

  tag {
    scope = "${var.nsx_tag_scope}"
    tag   = "${var.nsx_tag}"
  }
}

#
# This resource creates a logical port on the T0 router. We will connect the T1
# router to this port to enable connectivity from the tenant / application
# networks to the networks to the cloud.
#
resource "nsxt_logical_router_link_port_on_tier0" "link_port_tier0" {
  description       = "TIER0_PORT1 provisioned by Terraform"
  display_name      = "TIER0_PORT1"
  logical_router_id = "${data.nsxt_logical_tier0_router.tier0_router.id}"

  tag {
    scope = "${var.nsx_tag_scope}"
    tag   = "${var.nsx_tag}"
  }
}

#
# Here we create a tier 1 router uplink port and connect it to T0 router port
# created in the previous example.
#
resource "nsxt_logical_router_link_port_on_tier1" "link_port_tier1" {
  description                   = "TIER1_PORT1 provisioned by Terraform"
  display_name                  = "TIER1_PORT1"
  logical_router_id             = "${nsxt_logical_tier1_router.tier1_router.id}"
  linked_logical_router_port_id = "${nsxt_logical_router_link_port_on_tier0.link_port_tier0.id}"

  tag {
    scope = "${var.nsx_tag_scope}"
    tag   = "${var.nsx_tag}"
  }
}

#
# Like their physical counterpart a logical switch can have switch ports. In
# this example Terraform will create a logical switch port on a logical switch.
#
resource "nsxt_logical_port" "logical_port1" {
  admin_state       = "UP"
  description       = "LP1 provisioned by Terraform"
  display_name      = "LP1"
  logical_switch_id = "${nsxt_logical_switch.switch1.id}"

  tag {
    scope = "${var.nsx_tag_scope}"
    tag   = "${var.nsx_tag}"
  }
}

#
# In order to connect a logical switch to a tier 1 logical router we will need
# a downlink port on the tier 1 router and will need to  connect it to the
# switch port we created above.
#
# The IP address provided in the `ip_address` property will be default gateway
# for virtual machines connected to this logical switch.
#
resource "nsxt_logical_router_downlink_port" "downlink_port" {
  description                   = "DP1 provisioned by Terraform"
  display_name                  = "DP1"
  logical_router_id             = "${nsxt_logical_tier1_router.tier1_router.id}"
  linked_logical_switch_port_id = "${nsxt_logical_port.logical_port1.id}"
  ip_address                    = "192.168.245.1/24"

  tag {
    scope = "${var.nsx_tag_scope}"
    tag   = "${var.nsx_tag}"
  }
}
```

In order to be able to connect VMs to the newly created logical switch a new
`vpshere_network` datasource need to be defined.

```hcl
data "vsphere_network" "terraform_switch1" {
    name = "${nsxt_logical_switch.switch1.display_name}"
    datacenter_id = "${data.vsphere_datacenter.dc.id}"
    depends_on = ["nsxt_logical_switch.switch1"]
}

```

The datasource in the above example should be referred in `network_id` inside
`network_interface` section for `vsphere_virtual_machine` resource.

## Feature Requests, Bug Reports, and Contributing

For more information how how to submit feature requests, bug reports, or
details on how to make your own contributions to the provider, see the [NSX-T
provider project page][nsxt-provider-project-page].

[nsxt-provider-project-page]: https://github.com/terraform-providers/terraform-provider-nsxt
