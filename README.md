

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

## Documentation

## Releases & Major Branches

## Contributing

The terraform-provider-nsxt project team welcomes contributions from the community. If you wish to contribute code and you have not
signed our contributor license agreement (CLA), our bot will update the issue when you open a Pull Request. For any
questions about the CLA process, please refer to our [FAQ](https://cla.vmware.com/faq). For more detailed information,
refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License

This terraform provider is available under MPL2.0 license.
