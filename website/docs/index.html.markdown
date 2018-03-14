---
layout: "nsxt"
page_title: "Provider: NSXT"
sidebar_current: "docs-nsxt-index"
description: |-
  The VMware NSX-T Terraform Provider 
---

# The NSX-T Terraform Provider

The NSX-T Terraform Provider provides a way to automate NSX to provide virtualized networking and security services using both vSphere and KVM.

More information on NSX can be found on the [NSX Product Page](https://www.vmware.com/products/nsx.html)

Documenation on the NSX plaform can be found on the [NSX Documenatation page](https://docs.vmware.com/en/VMware-NSX-T/index.html)

## NSXTProvider
In order to use the NSX-T Terraform Provider you must first configure the proveder to communicate with the VMmare NSX-T manager. This required the IP address, hostname or FQDN of the NSX manager along with the username and password of the user to which the provider will authenticate. 

The provider also supports using both signed and self-signed certificates. It is recommended that in production environments you only use certificates signed by a certificate authority.

Ther are also a number of other parameters that can be set to tune how the provider connects to the NSX REST API.

Use the navigation to the left to read about available data sources and resources.

### Example Usage

```hcl
provider "nsxt" {
  host                 = "10.160.94.11"
  username             = "admin"
  password             = "qwerty"
  allow_unverified_ssl = true
  max_retries          = 10
  retry_min_delay      = 500
  retry_max_delay      = 5000
  retry_on_statuses    = [429]
}

```

### Client Authentication

NSXT Provider offers few ways of authentication. Credentials can be provided statically as shown in example above or provided as environment variables:

```hcl
export NSX_MANAGER_HOST="10.160.94.11"
export NSX_USERNAME="admin"
export NSX_PASSWORD="qwerty"
```

In addition, client certificates can be used for authentication. Terraform will require certificate file and private key file in PEM format. In this case, the client certificate needs to be registered with NSX-T manager prior to invoking Terraform.


```hcl
provider "nsxt" {
  host                  = "10.160.94.11"
  client_auth_cert_file = "mycert.pem"
  client_auth_key_file  = "mykey.pem"
  allow_unverified_ssl  = true
}

```

### Server Authentication

The "insecure" provider parameter set to true (shown in examples above) will direct Terraform client to skip server certificate verification. This is not recommended in production deployments and it is recommended to use trusted connection using certificates signed by a certificate authority.. The example below shows provides CA file (in PEM format) to verify server certificate.

```hcl
provider "nsxt" {
  host     = "10.160.94.11"
  username = "admin"
  password = "qwerty"
  ca_file  = "myca.pem"
}

```
