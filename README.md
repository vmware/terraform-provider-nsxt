# Terraform NSX-T Provider

This is the repository for the Terraform NSX Provider, which one can use with
Terraform to work with [VMware NSX-T][vmware-nsxt].

[vmware-nsxt]: https://www.vmware.com/products/nsx.html

For general information about Terraform, visit the [official
website][tf-website] and the [GitHub project page][tf-github].

[tf-website]: https://terraform.io/
[tf-github]: https://github.com/hashicorp/terraform

Documentation on the NSX platform can be found at the [NSX-T Documentation page](https://docs.vmware.com/en/VMware-NSX-T/index.html)

# Using the Provider

The latest version of this provider requires Terraform v0.12 or higher to run.

The VMware supported version of the provider requires NSX version 3.0.0 onwards and Terraform 0.12 onwards.
Version 2.0.0 of the provider offers NSX consumption via policy APIs, which is the recommended way.
Most policy resources are supported with NSX version 3.0.0 onwards, however some resources or attributes require later releases. Please refer to documentation for more details.
The recommended vSphere provider to be used in conjunction with the NSX-T Terraform Provider is 1.3.3 or above.

Note that you need to run `terraform init` to fetch the provider before
deploying.

## Full Provider Documentation

The provider is documented in full on the Terraform website and can be found
[here](https://registry.terraform.io/providers/vmware/nsxt/latest). Check the provider documentation for details on entering your connection information and how to get started with writing configuration for vSphere resources.

### Controlling the provider version

Note that you can also control the provider version. This requires the use of a `provider` block in your Terraform configuration if you have not added one already.

The syntax is as follows:

```hcl
provider "nsxt" {
  version = "~> 3.2"
  ...
}
```

Version locking uses a pessimistic operator, so this version lock would mean
anything within the 3.x namespace, including or after 3.0.0. [Read more][provider-vc] on provider version control.

[provider-vc]: https://www.terraform.io/docs/configuration/providers.html#provider-versions

# Automated Installation (Recommended)

Download and initialization of Terraform providers is with the “terraform init” command. This applies to the NSX-T provider as well. Once the provider block for the NSX-T provider is specified in your .tf file, “terraform init” will detect a need for the provider and download it to your environment.
You can list versions of providers installed in your environment by running “terraform version” command:

```hcl
$ ./terraform version
Terraform v1.0.0
on linux_amd64
+ provider registry.terraform.io/vmware/nsxt v3.2.4
```

# Manual Installation

**NOTE:** Unless you are [developing](#developing-the-provider) or require a
pre-release bugfix or feature, you will want to use the officially released
version of the provider (see [the section above](#using-the-provider)).

**NOTE:** Recommended way to compile the provider is using [Go Modules](https://blog.golang.org/using-go-modules).

**NOTE:** For terraform 0.13, please refer to [provider installation configuration][install-013] in order to use custom provider.

[install-013]: https://www.terraform.io/docs/commands/cli-config.html#provider-installation

## Cloning the Project

First, you will want to clone the repository to
`$GOPATH/src/github.com/vmware/terraform-provider-nsxt`:

```sh
mkdir -p $GOPATH/src/github.com/vmware
cd $GOPATH/src/github.com/vmware
git clone https://github.com/vmware/terraform-provider-nsxt.git
```

## Building and Installing the Provider

Recommended golang version is go1.17 onwards.
After the clone has been completed, you can enter the provider directory and build the provider.

```sh
cd $GOPATH/src/github.com/vmware/terraform-provider-nsxt
make
```

After the build is complete, copy the provider executable `terraform-provider-nsxt` into location specified in your provider installation configuration. Make sure to delete provider lock files that might exist in your working directory due to prior provider usage. Run `terraform init`.
For developing, consider using [dev overrides configuration][dev-overrides]. Please note that `terraform init` should not be used with dev overrides.

[dev-overrides]: https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers

# Developing the Provider

**NOTE:** Before you start work on a feature, please make sure to check the
[issue tracker][gh-issues] and existing [pull requests][gh-prs] to ensure that
work is not being duplicated. For further clarification, you can also ask in a
new issue.

[gh-issues]: https://github.com/vmware/terraform-provider-nsxt/issues
[gh-prs]: https://github.com/vmware/terraform-provider-nsxt/pulls

If you wish to work on the provider, you'll first need [Go][go-website]
installed on your machine (version 1.14+ is recommended). You'll also need to
correctly setup a [GOPATH][gopath], as well as adding `$GOPATH/bin` to your
`$PATH`.

[go-website]: https://golang.org/
[gopath]: http://golang.org/doc/code.html#GOPATH

See [Manual Installation](#manual-installation) for details on building the
provider.

# Testing the Provider

**NOTE:** Testing the NSX-T provider is currently a complex operation as it
requires having a NSX-T manager endpoint to test against, which should be
hosting a standard configuration for a NSX-T cluster. To cover Global Manager
test cases, NSX-T Global Manager suite needs to be preconfigured.

## Configuring Environment Variables

Most of the tests in this provider require a comprehensive list of environment
variables to run. See the individual `*_test.go` files in the [`nsxt/`](nsxt/)
directory for more details, in addition to
[`tests_utils.go`](nsxt/tests_utils.go) for details on some tunables that can be
used to specify the locations of certain pre-created resources that some tests
require.

Minimum environment variable :
```sh
$ export NSXT_MANAGER_HOST="192.168.110.41"
$ export NSXT_USERNAME="admin"
$ export NSXT_PASSWORD="MyPassword123!"
$ export NSXT_ALLOW_UNVERIFIED_SSL=true
```

## Running the Acceptance Tests

After this is done, you can run the acceptance tests by running:

```sh
$ make testacc
```

If you want to run against a specific set of tests, run `make testacc` with the
`TESTARGS` parameter containing the run mask as per below:

```sh
make testacc TESTARGS="-run=TestAccResourceNsxtPolicyTier0Gateway"
```

This following example would run all of the acceptance tests matching
`TestAccResourceNsxtPolicyTier0Gateway`. Change this for the specific tests you want
to run.

# Interoperability

The following versions of NSX are supported:

 * NSX-T 3.2.*
 * NSX-T 3.1.*
 * NSX-T 3.0.*
 * NSX-T 2.5.* support is limited with provider version 3.2.x and above
 
Some specific resources and attributed may require recent versions of NSX-T. Please refer to documentation for more details.

# Support

The NSX Terraform provider is now VMware supported as well as community supported. For bugs and feature requests please open a Github Issue and label it appropriately or contact VMware support.

# License

Copyright © 2015-2021 VMware, Inc. All Rights Reserved.

The NSX Terraform provider is available under [MPL2.0 license](https://github.com/vmware/terraform-provider-nsxt/blob/master/LICENSE.txt).
