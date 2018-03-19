---
layout: "nsxt"
page_title: "Provider: NSXT"
sidebar_current: "docs-nsxt-index"
description: |-
  The VMware NSX-T Terraform Provider
---

# The NSX Terraform Provider

The NSX Terraform provider gives the NSX administrator a way to automate NSX to provide virtualized networking and security services using both ESXi and KVM based hypervisor hosts as well as container networking and security.

More information on NSX can be found on the [NSX Product Page](https://www.vmware.com/products/nsx.html)

Documentation on the NSX platform can be found on the [NSX Documentation Page](https://docs.vmware.com/en/VMware-NSX-T/index.html)

Please use the navigation to the left to read about available data sources and resources.

## Basic Configuration of the NSX Terraform Provider

In order to use the NSX Terraform provider you must first configure the provider to communicate with the VMmare NSX manager. The NSX manager is the system which serves the NSX REST API and provides a way to configure the desired state of the NSX system. The configuration of the NSX provider requires the IP address, hostname, or FQDN of the NSX manager.

The NSX provider offers several ways to authenticate to the NSX manager. Credentials can be provided statically or provided as environment variables. In addition, client certificates can be used for authentication. For authentication with certificates Terraform will require a certificate file and private key file in PEM format. To use client certificates the client certificate needs to be registered with NSX-T manager prior to invoking Terraform.

The provider also can accept both signed and self-signed server certificates. It is recommended that in production environments you only use certificates signed by a certificate authority. NSX ships by default with a self-signed server certificates as the hostname of the NSX manager is not known until the NSX administrator determines what name or IP to use.

Setting the `allow_unverified_ssl` parameter to `true` will direct the Terraform client to skip server certificate verification. This is not recommended in production deployments as it is recommended that you use trusted connection using certificates signed by a certificate authority.

With the `ca_file` parameter you can also specify a file that contains your certificate authority certificate in PEM format to verify certificates with a certificate authority.

There are also a number of other parameters that can be set to tune how the provider connects to the NSX REST API. It is recommended you leave these to the defaults unless you experience issues in which case they can be tuned to optimize the system in your environment.

Note that in all of the examples you will need to update the `host`, `username`, and `password` settings to match those configured in your NSX deployment.

### Example of Configuration with Credentials

```hcl
provider "nsxt" {
  host                 = "192.168.110.41"
  username             = "admin"
  password             = "default"
  allow_unverified_ssl = true
  max_retries          = 10
  retry_min_delay      = 500
  retry_max_delay      = 5000
  retry_on_statuses    = [429]
}

```

### Example of Setting Environment Variables

```hcl
export NSX_MANAGER_HOST="192.168.110.41"
export NSX_USERNAME="admin"
export NSX_PASSWORD="default"
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

### Example Usage

The following abridged example demonstrates a current basic usage of the provider to create Logical Switch, Tier1 router, connect the Logical Switch to the T1 router as well as connect the T1 router to the provider T0 router.

```hcl
# Configure the VMware NSX-T Provider
provider "nsxt" {
  host                 = "192.168.110.41"
  username             = "admin"
  password             = "default"
  allow_unverified_ssl = true
}

# Define NSX-T Tag in order to be able easily to search for the created objects in NSX
variable "nsx_tag_scope" {
    default = "project"
}
variable "nsx_tag" {
    default = "terraform-demo"
}


# Create the data sources we will need to refer to later
data "nsxt_transport_zone" "overlay_tz" {
    display_name = "tz1"
}
data "nsxt_logical_tier0_router" "tier0_router" {
  display_name = "DefaultT0Router"
}
data "nsxt_edge_cluster" "edge_cluster1" {
    display_name = "EdgeCluster1"
}

# Create NSX-T Logical Switch
resource "nsxt_logical_switch" "switch1" {
    admin_state = "UP"
    description = "LS created by Terraform"
    display_name = "TfLogicalSwitch"
    transport_zone_id = "${data.nsxt_transport_zone.overlay_tz.id}"
    replication_mode = "MTEP"
    tag {
	scope = "${var.nsx_tag_scope}"
	tag = "${var.nsx_tag}"
    }
    tag {
	scope = "tenant"
	tag = "second_example_tag"
    }
}

# Create T1 router
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
	tag = "${var.nsx_tag}"
    }
}

# Create a port on the T0 router. We will connect the T1 router to this port
resource "nsxt_logical_router_link_port_on_tier0" "link_port_tier0" {
    description       = "TIER0_PORT1 provisioned by Terraform"
    display_name      = "TIER0_PORT1"
    logical_router_id = "${data.nsxt_logical_tier0_router.tier0_router.id}"
    tag {
	scope = "${var.nsx_tag_scope}"
	tag = "${var.nsx_tag}"
    }
}

# Create a T1 uplink port and connect it to T0 router
resource "nsxt_logical_router_link_port_on_tier1" "link_port_tier1" {
    description                   = "TIER1_PORT1 provisioned by Terraform"
    display_name                  = "TIER1_PORT1"
    logical_router_id             = "${nsxt_logical_tier1_router.tier1_router.id}"
    linked_logical_router_port_id = "${nsxt_logical_router_link_port_on_tier0.link_port_tier0.id}"
    tag {
	scope = "${var.nsx_tag_scope}"
	tag = "${var.nsx_tag}"
    }
}

# Create a switchport on our logical switch
resource "nsxt_logical_port" "logical_port1" {
    admin_state       = "UP"
    description       = "LP1 provisioned by Terraform"
    display_name      = "LP1"
    logical_switch_id = "${nsxt_logical_switch.switch1.id}"
    tag {
	scope = "${var.nsx_tag_scope}"
	tag = "${var.nsx_tag}"
    }
}

# Create downlink port on the T1 router and connect it to the switchport we created above
# The ip_address will be default gateway for VMs connected to this logical switch
resource "nsxt_logical_router_downlink_port" "downlink_port" {
    description                   = "DP1 provisioned by Terraform"
    display_name                  = "DP1"
    logical_router_id             = "${nsxt_logical_tier1_router.tier1_router.id}"
    linked_logical_switch_port_id = "${nsxt_logical_port.logical_port1.id}"
    ip_address                    = "192.168.245.1/24"
    tag {
	scope = "${var.nsx_tag_scope}"
	tag = "${var.nsx_tag}"
    }
}

```

In order to be able to connect VMs to the newly created logical switch a new vpshere_network datasource need to be defined.
```hcl
data "vsphere_network" "terraform_switch1" {
    name = "${nsxt_logical_switch.switch1.display_name}"
    datacenter_id = "${data.vsphere_datacenter.dc.id}"
    depends_on = ["nsxt_logical_switch.switch1"]
}

```
The datasource above should be refered in network_id inside network_interface section for vsphere_virtual_machine resource.

## Feature Requests, Bug Reports, and Contributing

For more information how how to submit feature requests, bug reports, or details on how to make your own contributions to the provider, see the NSX-T provider project page.
