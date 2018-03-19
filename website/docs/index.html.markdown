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

## Feature Requests, Bug Reports, and Contributing

For more information how how to submit feature requests, bug reports, or details on how to make your own contributions to the provider, see the NSX-T provider project page.
