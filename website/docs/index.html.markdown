---
layout: "nsxt"
page_title: "Provider: NSXT"
sidebar_current: "docs-nsxt-index"
description: |-
  The VMWare NSXT provider is used to configure VMWare NSX-T Manager.
---

# NSXT Provider

This provider is used to configure VMWare NSX-T Manager. The provider needs to be configured with NSX-T Manager host and credentials.

Use the navigation to the left to read about available data sources and resources.

## Example Usage

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

## Client Authentication

NSXT Provider offers few ways of authentication. Credentials can be provided statically as shown in example above or provided as environment variables:

```hcl
export NSX_MANAGER_HOST="10.160.94.11"
export NSX_USERNAME="admin"
export NSX_PASSWORD="qwerty"
```

In addition, self-signed client certificate can be used for authentication. Terraform will require certificate file and private key file in PEM format. In this case, the client certificate needs to be registered with NSX-T Manager prior to invoking terraform.


```hcl
provider "nsxt" {
  host                  = "10.160.94.11"
  client_auth_cert_file = "mycert.pem"
  client_auth_key_file  = "mykey.pem"
  allow_unverified_ssl  = true
}

```

## Server Authentication

The "insecure" provider parameter set to true (shown in examples above) will direct terraform client to skip server certificate verification. However, it is recommended to use trusted connection. The example below shows provides CA file (in PEM format) to verify server certificate.

```hcl
provider "nsxt" {
  host     = "10.160.94.11"
  username = "admin"
  password = "qwerty"
  ca_file  = "myca.pem"
}

```
