// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/QosProfilesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/QosProfilesClient.go QosProfilesClient

package nsxt

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	qosmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	qosDisplayName = "fooname"
	qosDescription = "this is a mock"
	qosPath        = "/infra/qos-profiles/my-nsxt_policy_qos_profile"
	qosRevision    = int64(1)
	qosID          = "my-nsxt_policy_qos_profile"
)

func newGoMockProviderClient() nsxtClients {
	mockProviderClient := constructMockProviderClient()
	mockProviderClient.NsxtClientConfig.HTTPClient = &http.Client{}
	mockProviderClient.PolicyHTTPClient = &http.Client{}
	return mockProviderClient
}

func TestMockResourceNsxtPolicyQosProfileCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQosProfilesSDK := qosmocks.NewMockQosProfilesClient(ctrl)

	mockWrapper := &cliinfra.QosProfileClientContext{
		Client:     mockQosProfilesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliQosProfilesClient
	defer func() { cliQosProfilesClient = originalCli }()
	cliQosProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.QosProfileClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		dscpMode := "UNTRUSTED"
		dscpPriority := int64(0)
		classOfService := int64(0)

		mockQosProfilesSDK.EXPECT().
			Patch(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)
		mockQosProfilesSDK.EXPECT().
			Get(gomock.Any()).
			Return(model.QosProfile{
				DisplayName:    &qosDisplayName,
				Description:    &qosDescription,
				Path:           &qosPath,
				Revision:       &qosRevision,
				ClassOfService: &classOfService,
				Dscp: &model.QosDscp{
					Mode:     &dscpMode,
					Priority: &dscpPriority,
				},
			}, nil)

		res := resourceNsxtPolicyQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":                  qosDisplayName,
			"description":                   qosDescription,
			"class_of_service":              0,
			"dscp_trusted":                  false,
			"dscp_priority":                 0,
			"ingress_rate_shaper":           []interface{}{},
			"ingress_broadcast_rate_shaper": []interface{}{},
			"egress_rate_shaper":            []interface{}{},
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyQosProfileCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockQosProfilesSDK.EXPECT().
			Get("existing-id").
			Return(model.QosProfile{Id: &qosID}, nil)

		res := resourceNsxtPolicyQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"nsx_id":                        "existing-id",
			"display_name":                  qosDisplayName,
			"description":                   qosDescription,
			"class_of_service":              0,
			"dscp_trusted":                  false,
			"dscp_priority":                 0,
			"ingress_rate_shaper":           []interface{}{},
			"ingress_broadcast_rate_shaper": []interface{}{},
			"egress_rate_shaper":            []interface{}{},
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyQosProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyQosProfileRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQosProfilesSDK := qosmocks.NewMockQosProfilesClient(ctrl)

	mockWrapper := &cliinfra.QosProfileClientContext{
		Client:     mockQosProfilesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliQosProfilesClient
	defer func() { cliQosProfilesClient = originalCli }()
	cliQosProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.QosProfileClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		dscpMode := "TRUSTED"
		dscpPriority := int64(53)
		classOfService := int64(5)

		mockQosProfilesSDK.EXPECT().
			Get(qosID).
			Return(model.QosProfile{
				DisplayName:    &qosDisplayName,
				Description:    &qosDescription,
				Path:           &qosPath,
				Revision:       &qosRevision,
				ClassOfService: &classOfService,
				Dscp: &model.QosDscp{
					Mode:     &dscpMode,
					Priority: &dscpPriority,
				},
			}, nil)

		res := resourceNsxtPolicyQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(qosID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyQosProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, qosDisplayName, d.Get("display_name"))
		assert.Equal(t, qosDescription, d.Get("description"))
		assert.Equal(t, qosPath, d.Get("path"))
		assert.Equal(t, int(qosRevision), d.Get("revision"))
		assert.Equal(t, 5, d.Get("class_of_service"))
		assert.True(t, d.Get("dscp_trusted").(bool))
		assert.Equal(t, 53, d.Get("dscp_priority"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		// d.Id() is empty

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyQosProfileRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining QosProfile ID")
	})
}

func TestMockResourceNsxtPolicyQosProfileUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQosProfilesSDK := qosmocks.NewMockQosProfilesClient(ctrl)

	mockWrapper := &cliinfra.QosProfileClientContext{
		Client:     mockQosProfilesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliQosProfilesClient
	defer func() { cliQosProfilesClient = originalCli }()
	cliQosProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.QosProfileClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		updatedName := "updated-name"
		dscpMode := "TRUSTED"
		dscpPriority := int64(10)
		classOfService := int64(3)

		mockQosProfilesSDK.EXPECT().
			Patch(qosID, gomock.Any(), gomock.Any()).
			Return(nil)
		mockQosProfilesSDK.EXPECT().
			Get(qosID).
			Return(model.QosProfile{
				DisplayName:    &updatedName,
				Description:    &qosDescription,
				Path:           &qosPath,
				Revision:       &qosRevision,
				ClassOfService: &classOfService,
				Dscp: &model.QosDscp{
					Mode:     &dscpMode,
					Priority: &dscpPriority,
				},
			}, nil)

		res := resourceNsxtPolicyQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":                  updatedName,
			"description":                   qosDescription,
			"class_of_service":              3,
			"dscp_trusted":                  true,
			"dscp_priority":                 10,
			"ingress_rate_shaper":           []interface{}{},
			"ingress_broadcast_rate_shaper": []interface{}{},
			"egress_rate_shaper":            []interface{}{},
		})
		d.SetId(qosID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyQosProfileUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, updatedName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":                  qosDisplayName,
			"description":                   qosDescription,
			"class_of_service":              0,
			"dscp_trusted":                  false,
			"dscp_priority":                 0,
			"ingress_rate_shaper":           []interface{}{},
			"ingress_broadcast_rate_shaper": []interface{}{},
			"egress_rate_shaper":            []interface{}{},
		})
		// d.Id() is empty

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyQosProfileUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining QosProfile ID")
	})
}

func TestMockResourceNsxtPolicyQosProfileDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQosProfilesSDK := qosmocks.NewMockQosProfilesClient(ctrl)

	mockWrapper := &cliinfra.QosProfileClientContext{
		Client:     mockQosProfilesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliQosProfilesClient
	defer func() { cliQosProfilesClient = originalCli }()
	cliQosProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.QosProfileClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockQosProfilesSDK.EXPECT().
			Delete(qosID, gomock.Any()).
			Return(nil)

		res := resourceNsxtPolicyQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(qosID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyQosProfileDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		// d.Id() is empty

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyQosProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining QosProfile ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockQosProfilesSDK.EXPECT().
			Delete(qosID, gomock.Any()).
			Return(errors.New("API error"))

		res := resourceNsxtPolicyQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(qosID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyQosProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}
