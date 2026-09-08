//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/LbPoolsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/LbPoolsClient.go LbPoolsClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	lbPoolMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	lbPoolID          = "lb-pool-1"
	lbPoolDisplayName = "lb-pool-fooname"
	lbPoolDescription = "lb pool mock"
	lbPoolPath        = "/infra/lb-pools/lb-pool-1"
	lbPoolRevision    = int64(1)
	lbPoolAlgorithm   = model.LBPool_ALGORITHM_ROUND_ROBIN
	lbPoolMinActive   = int64(1)
)

func TestMockResourceNsxtPolicyLBPoolCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPoolSDK := lbPoolMocks.NewMockLbPoolsClient(ctrl)
	mockWrapper := &cliinfra.LBPoolClientContext{
		Client:     mockPoolSDK,
		ClientType: utl.Local,
	}

	originalCli := cliLbPoolsClient
	defer func() { cliLbPoolsClient = originalCli }()
	cliLbPoolsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.LBPoolClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockPoolSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockPoolSDK.EXPECT().Get(gomock.Any()).Return(minimalLBPoolModel(), nil)

		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPoolData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockPoolSDK.EXPECT().Get("existing-id").Return(model.LBPool{Id: &lbPoolID}, nil)

		res := resourceNsxtPolicyLBPool()
		data := minimalLBPoolData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyLBPoolRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPoolSDK := lbPoolMocks.NewMockLbPoolsClient(ctrl)
	mockWrapper := &cliinfra.LBPoolClientContext{
		Client:     mockPoolSDK,
		ClientType: utl.Local,
	}

	originalCli := cliLbPoolsClient
	defer func() { cliLbPoolsClient = originalCli }()
	cliLbPoolsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.LBPoolClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockPoolSDK.EXPECT().Get(lbPoolID).Return(minimalLBPoolModel(), nil)

		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(lbPoolID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, lbPoolDisplayName, d.Get("display_name"))
		assert.Equal(t, lbPoolDescription, d.Get("description"))
		assert.Equal(t, lbPoolPath, d.Get("path"))
		assert.Equal(t, int(lbPoolRevision), d.Get("revision"))
		assert.Equal(t, lbPoolID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining LBPool ID")
	})
}

func TestMockResourceNsxtPolicyLBPoolUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPoolSDK := lbPoolMocks.NewMockLbPoolsClient(ctrl)
	mockWrapper := &cliinfra.LBPoolClientContext{
		Client:     mockPoolSDK,
		ClientType: utl.Local,
	}

	originalCli := cliLbPoolsClient
	defer func() { cliLbPoolsClient = originalCli }()
	cliLbPoolsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.LBPoolClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockPoolSDK.EXPECT().Update(lbPoolID, gomock.Any()).Return(minimalLBPoolModel(), nil)
		mockPoolSDK.EXPECT().Get(lbPoolID).Return(minimalLBPoolModel(), nil)

		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPoolData())
		d.Set("revision", 1)
		d.SetId(lbPoolID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPoolData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining LBPool ID")
	})
}

func TestMockResourceNsxtPolicyLBPoolDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPoolSDK := lbPoolMocks.NewMockLbPoolsClient(ctrl)
	mockWrapper := &cliinfra.LBPoolClientContext{
		Client:     mockPoolSDK,
		ClientType: utl.Local,
	}

	originalCli := cliLbPoolsClient
	defer func() { cliLbPoolsClient = originalCli }()
	cliLbPoolsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.LBPoolClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockPoolSDK.EXPECT().Delete(lbPoolID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(lbPoolID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining LBPool ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockPoolSDK.EXPECT().Delete(lbPoolID, gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(lbPoolID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestMockResourceNsxtPolicyLBPoolExists(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPoolSDK := lbPoolMocks.NewMockLbPoolsClient(ctrl)
	mockWrapper := &cliinfra.LBPoolClientContext{
		Client:     mockPoolSDK,
		ClientType: utl.Local,
	}

	originalCli := cliLbPoolsClient
	defer func() { cliLbPoolsClient = originalCli }()
	cliLbPoolsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.LBPoolClientContext {
		return mockWrapper
	}

	t.Run("Exists returns true when Get succeeds", func(t *testing.T) {
		mockPoolSDK.EXPECT().Get(lbPoolID).Return(minimalLBPoolModel(), nil)

		exists, err := resourceNsxtPolicyLBPoolExists(lbPoolID, getPolicyConnector(newGoMockProviderClient()), false)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Exists returns false when not found", func(t *testing.T) {
		mockPoolSDK.EXPECT().Get(lbPoolID).Return(model.LBPool{}, vapiErrors.NotFound{})

		exists, err := resourceNsxtPolicyLBPoolExists(lbPoolID, getPolicyConnector(newGoMockProviderClient()), false)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Exists propagates other errors", func(t *testing.T) {
		mockPoolSDK.EXPECT().Get(lbPoolID).Return(model.LBPool{}, vapiErrors.InternalServerError{})

		_, err := resourceNsxtPolicyLBPoolExists(lbPoolID, getPolicyConnector(newGoMockProviderClient()), false)
		require.Error(t, err)
	})
}

func TestUnitNsxt_getPolicyPoolMembersFromSchema(t *testing.T) {
	res := resourceNsxtPolicyLBPool()
	data := minimalLBPoolData()
	data["member"] = []interface{}{
		map[string]interface{}{
			"display_name":               "member-1",
			"admin_state":                "ENABLED",
			"backup_member":              false,
			"ip_address":                 "1.1.1.1",
			"max_concurrent_connections": 5,
			"port":                       "80",
			"weight":                     2,
		},
	}
	d := schema.TestResourceDataRaw(t, res.Schema, data)

	members := getPolicyPoolMembersFromSchema(d)
	require.Len(t, members, 1)
	assert.Equal(t, "member-1", *members[0].DisplayName)
	assert.Equal(t, "1.1.1.1", *members[0].IpAddress)
	assert.Equal(t, "80", *members[0].Port)
	assert.EqualValues(t, 5, *members[0].MaxConcurrentConnections)
}

func TestUnitNsxt_setPolicyPoolMembersInSchema(t *testing.T) {
	res := resourceNsxtPolicyLBPool()
	d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPoolData())

	name, port := "member-1", "80"
	weight := int64(3)
	maxConn := int64(10)
	adminState := "ENABLED"
	backup := false
	ip := "2.2.2.2"
	err := setPolicyPoolMembersInSchema(d, []model.LBPoolMember{
		{
			DisplayName:              &name,
			AdminState:               &adminState,
			BackupMember:             &backup,
			IpAddress:                &ip,
			Port:                     &port,
			Weight:                   &weight,
			MaxConcurrentConnections: &maxConn,
		},
	})
	require.NoError(t, err)

	members := d.Get("member").([]interface{})
	require.Len(t, members, 1)
	elem := members[0].(map[string]interface{})
	assert.Equal(t, "member-1", elem["display_name"])
	assert.Equal(t, "80", elem["port"])
}

func TestUnitNsxt_getPolicyPoolMemberGroupFromSchema(t *testing.T) {
	res := resourceNsxtPolicyLBPool()

	t.Run("nil when member_group is not set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPoolData())
		assert.Nil(t, getPolicyPoolMemberGroupFromSchema(d))
	})

	t.Run("builds group with ipv4-only filter", func(t *testing.T) {
		data := minimalLBPoolData()
		data["member_group"] = []interface{}{
			map[string]interface{}{
				"group_path":       "/infra/domains/default/groups/g1",
				"allow_ipv4":       true,
				"allow_ipv6":       false,
				"max_ip_list_size": 10,
				"port":             "8080",
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		group := getPolicyPoolMemberGroupFromSchema(d)
		require.NotNil(t, group)
		assert.Equal(t, model.LBPoolMemberGroup_IP_REVISION_FILTER_IPV4, *group.IpRevisionFilter)
		assert.EqualValues(t, 8080, *group.Port)
		assert.EqualValues(t, 10, *group.MaxIpListSize)
	})

	t.Run("builds group with dual-stack filter", func(t *testing.T) {
		data := minimalLBPoolData()
		data["member_group"] = []interface{}{
			map[string]interface{}{
				"group_path": "/infra/domains/default/groups/g1",
				"allow_ipv4": true,
				"allow_ipv6": true,
				"port":       "",
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		group := getPolicyPoolMemberGroupFromSchema(d)
		require.NotNil(t, group)
		assert.Equal(t, model.LBPoolMemberGroup_IP_REVISION_FILTER_IPV4_IPV6, *group.IpRevisionFilter)
	})

	t.Run("builds group with ipv6-only filter", func(t *testing.T) {
		data := minimalLBPoolData()
		data["member_group"] = []interface{}{
			map[string]interface{}{
				"group_path": "/infra/domains/default/groups/g1",
				"allow_ipv4": false,
				"allow_ipv6": false,
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		group := getPolicyPoolMemberGroupFromSchema(d)
		require.NotNil(t, group)
		assert.Equal(t, model.LBPoolMemberGroup_IP_REVISION_FILTER_IPV6, *group.IpRevisionFilter)
	})
}

func TestUnitNsxt_setPolicyPoolMemberGroupInSchema(t *testing.T) {
	res := resourceNsxtPolicyLBPool()

	t.Run("clears schema when group is nil", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPoolData())
		require.NoError(t, setPolicyPoolMemberGroupInSchema(d, nil))
	})

	t.Run("sets fields from group", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPoolData())
		path := "/infra/domains/default/groups/g1"
		filter := model.LBPoolMemberGroup_IP_REVISION_FILTER_IPV4_IPV6
		port := int64(80)
		maxSize := int64(5)
		err := setPolicyPoolMemberGroupInSchema(d, &model.LBPoolMemberGroup{
			GroupPath:        &path,
			IpRevisionFilter: &filter,
			Port:             &port,
			MaxIpListSize:    &maxSize,
		})
		require.NoError(t, err)

		groups := d.Get("member_group").([]interface{})
		require.Len(t, groups, 1)
		elem := groups[0].(map[string]interface{})
		assert.Equal(t, true, elem["allow_ipv4"])
		assert.Equal(t, true, elem["allow_ipv6"])
		assert.Equal(t, "80", elem["port"])
	})
}

func TestUnitNsxt_getPolicyPoolSnatFromSchema(t *testing.T) {
	res := resourceNsxtPolicyLBPool()

	t.Run("DISABLED type", func(t *testing.T) {
		data := minimalLBPoolData()
		data["snat"] = []interface{}{map[string]interface{}{"type": "DISABLED"}}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		val, err := getPolicyPoolSnatFromSchema(d)
		require.NoError(t, err)
		require.NotNil(t, val)
	})

	t.Run("AUTOMAP type", func(t *testing.T) {
		data := minimalLBPoolData()
		data["snat"] = []interface{}{map[string]interface{}{"type": "AUTOMAP"}}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		val, err := getPolicyPoolSnatFromSchema(d)
		require.NoError(t, err)
		require.NotNil(t, val)
	})

	t.Run("IPPOOL type with CIDR and single IP", func(t *testing.T) {
		data := minimalLBPoolData()
		data["snat"] = []interface{}{map[string]interface{}{
			"type":              "IPPOOL",
			"ip_pool_addresses": []interface{}{"10.0.0.0/24", "10.0.1.5"},
		}}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		val, err := getPolicyPoolSnatFromSchema(d)
		require.NoError(t, err)
		require.NotNil(t, val)
	})

	t.Run("no snat configured returns nil", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPoolData())
		val, err := getPolicyPoolSnatFromSchema(d)
		require.NoError(t, err)
		assert.Nil(t, val)
	})
}

func TestUnitNsxt_setPolicyPoolSnatInSchema(t *testing.T) {
	res := resourceNsxtPolicyLBPool()

	t.Run("nil snat is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPoolData())
		require.NoError(t, setPolicyPoolSnatInSchema(d, nil))
	})

	t.Run("round-trips DISABLED, AUTOMAP, and IPPOOL", func(t *testing.T) {
		for _, snatData := range []map[string]interface{}{
			{"type": "DISABLED"},
			{"type": "AUTOMAP"},
			{"type": "IPPOOL", "ip_pool_addresses": []interface{}{"10.0.0.0/24"}},
		} {
			data := minimalLBPoolData()
			data["snat"] = []interface{}{snatData}
			d := schema.TestResourceDataRaw(t, res.Schema, data)
			val, err := getPolicyPoolSnatFromSchema(d)
			require.NoError(t, err)
			require.NotNil(t, val)

			d2 := schema.TestResourceDataRaw(t, res.Schema, minimalLBPoolData())
			require.NoError(t, setPolicyPoolSnatInSchema(d2, val))
		}
	})
}

func minimalLBPoolModel() model.LBPool {
	return model.LBPool{
		Id:               &lbPoolID,
		DisplayName:      &lbPoolDisplayName,
		Description:      &lbPoolDescription,
		Path:             &lbPoolPath,
		Revision:         &lbPoolRevision,
		Algorithm:        &lbPoolAlgorithm,
		MinActiveMembers: &lbPoolMinActive,
	}
}

func minimalLBPoolData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbPoolDisplayName,
		"description":  lbPoolDescription,
	}
}
