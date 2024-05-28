/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type policyPathTest struct {
	path string
	parents []string
}

func TestParseStandardPolicyPath(t *testing.T) {

	testData := []policyPathTest{
		policyPathTest{
			path: "/infra/tier-1s/mygw1",
			parents: []string{"mygw1"},
		},
		policyPathTest{
			path: "/global-infra/tier-0s/mygw1",
			parents: []string{"mygw1"},
		},
		policyPathTest{
			path: "/orgs/myorg/projects/myproj/infra/tier-1s/mygw1",
			parents:[]string{"myorg", "myproj", "mygw1"},
		},
		policyPathTest{
			path: "/orgs/myorg/projects/myproj/vpcs/myvpc/tier-1s/mygw1",
			parents:[]string{"myorg", "myproj", "myvpc", "mygw1"},
		},
		policyPathTest{
			path: "/orgs/myorg/projects/myproj/infra/domains/d1/groups/g1",
			parents:[]string{"myorg", "myproj", "d1", "g1"},
		},
		policyPathTest{
			path: "/global-infra/tier-1s/t15/ipsec-vpn-services/default/sessions/xxx-yyy-xxx",
			parents: []string{"t15", "default", "xxx-yyy-xxx"},
		},
		policyPathTest{
			path: "/infra/evpn-tenant-configs/{config-id}",
			parents: []string{"{config-id}"},
		},
		policyPathTest{
			path: "/infra/tier-0s/{tier-0-id}/locale-services/{locale-service-id}/service-interfaces/{interface-id}",
			parents: []string{"{tier-0-id}", "{locale-service-id}", "{interface-id}"},
		},
	}

	for _, test := range testData {
		parents, err := parseStandardPolicyPath(test.path)
		assert.Nil(t, err)
		assert.Equal(t, test.parents, parents)
	}
}

func TestIsPolicyPath(t *testing.T) {

	testData := []string{
		"/infra/tier-1s/mygw1",
		"/global-infra/tier-1s/mygw1",
		"/orgs/infra/tier-1s/mygw1",
		"/orgs/myorg/projects/myproj/domains/d",
	}

	for _, test := range testData {
		pp := isPolicyPath(test)
		assert.True(t, pp)
	}
}


func TestNegativeParseStandardPolicyPath(t *testing.T) {

	testData := []string{
		"/some-infra/tier-1s/mygw1",
		"orgs/infra/tier-1s/mygw1-1",
		"orgs/infra  /tier-1s/mygw1-1",
	}

	for _, test := range testData {
		_, err := parseStandardPolicyPath(test)
		assert.NotNil(t, err)
	}
}
