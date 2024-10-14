<!--
© Broadcom. All Rights Reserved.
The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
SPDX-License-Identifier: MPL-2.0
-->

<!-- markdownlint-disable first-line-h1 no-inline-html -->

<img src="images/icon-color.svg" alt="VMware NSX" width="150">

# Testing the Terraform Provider for VMware NSX

Testing the Terraform Provider for VMware NSX is currently a complex operation
as it requires having an NSX Local Manager endpoint to test against, which
should be hosting a standard configuration for a NSX Manager cluster. To cover
NSX Global Manager test cases, NSX Global Manager suite needs to be
pre-configured.

## Configuring Environment Variables

Most of the tests for the provider require a comprehensive list of environment
variables to run. Please refer to the individual `*_test.go` files in the
[`nsxt/`](nsxt/) directory for more details. In addition, refer to the
`tests_utils.go` for details on some tunables that can be used to specify the
locations of certain pre-created resources that some tests require.

Minimum Environment Variable:

```sh
$ export NSXT_MANAGER_HOST="nsx-01.example.com"
$ export NSXT_USERNAME="admin"
$ export NSXT_PASSWORD="VMw@re123!VMw@re123!"
$ export NSXT_ALLOW_UNVERIFIED_SSL=true
```

## Running the Acceptance Tests

You can run the acceptance tests by running:

```sh
$ make testacc
```

If you want to run against a specific set of tests, run `make testacc` with the
`TESTARGS` parameter containing the run mask. For example:

```sh
make testacc TESTARGS="-run=TestAccResourceNsxtPolicyTier0Gateway"
```

This following example would run all of the acceptance tests matching
`TestAccResourceNsxtPolicyTier0Gateway`. Change this for the specific tests you
want to run.
