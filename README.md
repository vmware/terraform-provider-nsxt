

# terraform-provider-nsxt
This is terraform provider for vmware NSX-T

## Overview

Supported data sources:

* TransportZone
* SwitchingProfile

Supported resources:

* L4PortSetNsService
* LogicalPort
* LogicalSwitch

## Try it out

### Prerequisites

* Go 1.9.x onwards
* Terraform 0.10.x
* This repo makes use of go-vmware-nsxt library

### Build & Run

1. go get github.com/vmware/terraform-provider-nsxt
2. go build -o terraform-provider-nsxt
3. copy terraform-provider-nsxt to terraform running folder

## Contributing

The terraform-provider-nsxt project team welcomes contributions from the community. Before you start working with terraform-provider-nsxt, please read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on as an open-source patch. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License

This terraform provider is available under MPL2.0 license.
