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

The current version of this provider requires Terraform v0.10.2 or higher to
run.

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

# Building The Provider

**NOTE:** Unless you are [developing](#developing-the-provider) or require a
pre-release bugfix or feature, you will want to use the officially released
version of the provider (see [the section above](#using-the-provider)).

## Cloning the Project

First, you will want to clone the repository to
`$GOPATH/src/github.com/terraform-providers/terraform-provider-nsxt`:

```sh
mkdir -p $GOPATH/src/github.com/terraform-providers
cd $GOPATH/src/github.com/terraform-providers
git clone git@github.com:terraform-providers/terraform-provider-nsxt
```

## Running the Build

After the clone has been completed, you can enter the provider directory and
build the provider.

```sh
cd $GOPATH/src/github.com/terraform-providers/terraform-provider-nsxt
make build
```

## Installing the Local Plugin

After the build is complete, copy the `terraform-provider-nsxt` binary into
the same path as your `terraform` binary, and re-run `terraform init`.

After this, your project-local `.terraform/plugins/ARCH/lock.json` (where `ARCH`
matches the architecture of your machine) file should contain a SHA256 sum that
matches the local plugin. Run `shasum -a 256` on the binary to verify the values
match.

# Developing the Provider

**NOTE:** Before you start work on a feature, please make sure to check the
[issue tracker][gh-issues] and existing [pull requests][gh-prs] to ensure that
work is not being duplicated. For further clarification, you can also ask in a
new issue.

[gh-issues]: https://github.com/terraform-providers/terraform-provider-nsxt/issues
[gh-prs]: https://github.com/terraform-providers/terraform-provider-nsxt/pulls

If you wish to work on the provider, you'll first need [Go][go-website]
installed on your machine (version 1.9+ is **required**). You'll also need to
correctly setup a [GOPATH][gopath], as well as adding `$GOPATH/bin` to your
`$PATH`.

[go-website]: https://golang.org/
[gopath]: http://golang.org/doc/code.html#GOPATH

See [Building the Provider](#building-the-provider) for details on building the
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

 * NSX-T 2.1.*

# Support

The NSX Terraform provider is community supported. For bugs and feature requests please open a Github Issue and label it appropriately. As this is a community supported solution there is no SLA for resolutions.

# License

Copyright © 2015-2018 VMware, Inc. All Rights Reserved.

The NSX Terraform provider is available under [MPL2.0 license](https://github.com/terraform-providers/terraform-provider-nsxt/blob/master/LICENSE.txt).
