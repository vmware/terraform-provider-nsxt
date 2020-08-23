# Terraform NSX-T Provider

This is the repository for the Terraform NSX Provider, which one can use with
Terraform to work with [VMware NSX-T][vmware-nsxt].

[vmware-nsxt]: https://www.vmware.com/products/nsx.html

For general information about Terraform, visit the [official
website][tf-website] and the [GitHub project page][tf-github].

[tf-website]: https://terraform.io/
[tf-github]: https://github.com/hashicorp/terraform

This provider plugin is maintained by a collaboration between
[VMware](https://www.vmware.com/) and the Terraform team at
[HashiCorp](https://www.hashicorp.com/).

Documentation on the NSX platform can be found at the [NSX-T Documentation page](https://docs.vmware.com/en/VMware-NSX-T/index.html)

# Using the Provider

The latest version of this provider requires Terraform v0.12 or higher to run.

The VMware supported version of the provider requires NSX version 2.2 onwards and Terraform 0.12 onwards.
Version 2.0.0 of the provider offers NSX consumption via policy APIs, which is the recommended way.
Most policy resources are supported with NSX version 2.5 onwards, however some resources or attributes require NSX 3.0 onwards. Please refer to documentation for more details.
The recommended vSphere provider to be used in conjunction with the NSX-T Terraform Provider is 1.3.3 or above.

Note that you need to run `terraform init` to fetch the provider before
deploying. Read about the provider split and other changes to TF v0.10.0 in the
official release announcement found [here][tf-0.10-announce].

[tf-0.10-announce]: https://www.hashicorp.com/blog/hashicorp-terraform-0-10/

## Full Provider Documentation

The provider is documented in full on the Terraform website and can be found
[here][tf-nsxt-docs]. Check the provider documentation for details on entering
your connection information and how to get started with writing configuration
for vSphere resources.

[tf-nsxt-docs]: https://www.terraform.io/docs/providers/nsxt/index.html

### Controlling the provider version

Note that you can also control the provider version. This requires the use of a
`provider` block in your Terraform configuration if you have not added one
already.

The syntax is as follows:

```hcl
provider "nsxt" {
  version = "~> 1.0"
  ...
}
```

Version locking uses a pessimistic operator, so this version lock would mean
anything within the 1.x namespace, including or after 1.0.0. [Read
more][provider-vc] on provider version control.

[provider-vc]: https://www.terraform.io/docs/configuration/providers.html#provider-versions

# Automated Installation (Recommended)

Download and initialization of Terraform providers is with the “terraform init” command. This applies to the NSX-T provider as well. Once the provider block for the NSX-T provider is specified in your .tf file, “terraform init” will detect a need for the provider and download it to your environment.
You can list versions of providers installed in your environment by running “terraform version” command:

```hcl
$ ./terraform version
Terraform v0.12.7
+ provider.nsxt v2.0.0
+ provider.vsphere v1.5.0
```

# Manual Installation

**NOTE:** Unless you are [developing](#developing-the-provider) or require a
pre-release bugfix or feature, you will want to use the officially released
version of the provider (see [the section above](#using-the-provider)).

**NOTE:** Recommended way to compile the provider is using [Go Modules](https://blog.golang.org/using-go-modules), however vendored dependencies are still supported.

**NOTE:** For terraform 0.13, please refer to [provider installation configuration][install-013] in order to use custom provider.
[install-013]: https://www.terraform.io/docs/commands/cli-config.html#provider-installation

**NOTE:** Note that if the provider is manually copied to your running folder (rather than fetched with the “terraform init” based on provider block), Terraform is not aware of the version of the provider you’re running. It will appear as “unversioned”:
```hcl
$ ./terraform version
Terraform v0.12.7
+ provider.nsxt (unversioned)
+ provider.vsphere v1.5.0
```
Since Terraform has no indication of version, it cannot upgrade in a native way, based on the “version” attribute in provider block.
In addition, this may cause difficulties in housekeeping and issue reporting.

## Cloning the Project

First, you will want to clone the repository to
`$GOPATH/src/github.com/vmware/terraform-provider-nsxt`:

```sh
mkdir -p $GOPATH/src/github.com/vmware
cd $GOPATH/src/github.com/vmware
git clone https://github.com/vmware/terraform-provider-nsxt.git
```

## Building and Installing the Provider

After the clone has been completed, you can enter the provider directory and build the provider.

```sh
cd $GOPATH/src/github.com/vmware/terraform-provider-nsxt
make
```

After the build is complete, if your terraform running folder does not match your GOPATH environment, you need to copy the `terraform-provider-nsxt` executable to your running folder and re-run `terraform init` to make terraform aware of your local provider executable.

After this, your project-local `.terraform/plugins/ARCH/lock.json` (where `ARCH`
matches the architecture of your machine) file should contain a SHA256 sum that
matches the local plugin. Run `shasum -a 256` on the binary to verify the values
match.

# Developing the Provider

**NOTE:** Before you start work on a feature, please make sure to check the
[issue tracker][gh-issues] and existing [pull requests][gh-prs] to ensure that
work is not being duplicated. For further clarification, you can also ask in a
new issue.

[gh-issues]: https://github.com/vmware/terraform-provider-nsxt/issues
[gh-prs]: https://github.com/vmware/terraform-provider-nsxt/pulls

If you wish to work on the provider, you'll first need [Go][go-website]
installed on your machine (version 1.11+ is **required**). You'll also need to
correctly setup a [GOPATH][gopath], as well as adding `$GOPATH/bin` to your
`$PATH`.

[go-website]: https://golang.org/
[gopath]: http://golang.org/doc/code.html#GOPATH

See [Manual Installation](#manual-installation) for details on building the
provider.

# Testing the Provider

**NOTE:** Testing the NSX-T provider is currently a complex operation as it
requires having a NSX-T manager endpoint to test against, which should be
hosting a standard configuration for a NSX-T cluster.

## Configuring Environment Variables

Most of the tests in this provider require a comprehensive list of environment
variables to run. See the individual `*_test.go` files in the [`nsxt/`](nsxt/)
directory for more details, in addition to
[`tests_utils.go`](nsxt/tests_utils.go) for details on some tunables that can be
used to specify the locations of certain pre-created resources that some tests
require.

## Running the Acceptance Tests

After this is done, you can run the acceptance tests by running:

```sh
$ make testacc
```

If you want to run against a specific set of tests, run `make testacc` with the
`TESTARGS` parameter containing the run mask as per below:

```sh
make testacc TESTARGS="-run=TestAccResourceNsxtLogicalSwitch"
```

This following example would run all of the acceptance tests matching
`TestAccResourceNsxtLogicalSwitch`. Change this for the specific tests you want
to run.

# Interoperability

The following versions of NSX are supported:

 * NSX-T 3.0
 * NSX-T 2.5.*
 * NSX-T 2.4.*
 * NSX-T 2.3.*
 * NSX-T 2.2.*
 

Some specific resources may require later versions of NSX-T.

# Support

The NSX Terraform provider is now VMware supported as well as community supported. For bugs and feature requests please open a Github Issue and label it appropriately or contact VMware support.

# License

Copyright © 2015-2019 VMware, Inc. All Rights Reserved.

The NSX Terraform provider is available under [MPL2.0 license](https://github.com/vmware/terraform-provider-nsxt/blob/master/LICENSE.txt).
