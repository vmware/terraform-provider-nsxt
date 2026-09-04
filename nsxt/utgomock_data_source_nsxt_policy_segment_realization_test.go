//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// The mocks used here (MockStateClient, MockSegmentsClient) already exist under mocks/infra
// and mocks/infra/segments, reused from utgomock_resource_nsxt_policy_segment_test.go.

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	segmentsapi "github.com/vmware/terraform-provider-nsxt/api/infra/segments"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	segmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	segstatemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/segments"
)

func setupSegmentRealizationMocks(t *testing.T, ctrl *gomock.Controller) (*segstatemocks.MockStateClient, *segmocks.MockSegmentsClient, func()) {
	t.Helper()
	mockStateSDK := segstatemocks.NewMockStateClient(ctrl)
	stateWrapper := &segmentsapi.SegmentConfigurationStateClientContext{Client: mockStateSDK, ClientType: utl.Local}

	mockSegmentsSDK := segmocks.NewMockSegmentsClient(ctrl)
	segmentsWrapper := &cliinfra.SegmentClientContext{Client: mockSegmentsSDK, ClientType: utl.Local}

	originalState := cliSegmentStateClient
	originalSegments := cliSegmentsClient
	cliSegmentStateClient = func(sessionContext utl.SessionContext, connector client.Connector) *segmentsapi.SegmentConfigurationStateClientContext {
		return stateWrapper
	}
	cliSegmentsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SegmentClientContext {
		return segmentsWrapper
	}
	return mockStateSDK, mockSegmentsSDK, func() {
		cliSegmentStateClient = originalState
		cliSegmentsClient = originalSegments
	}
}

func segmentRealizationTestData() map[string]interface{} {
	return map[string]interface{}{
		"path": "/infra/segments/seg-rlz-1",
	}
}

func TestMockDataSourceNsxtPolicySegmentRealizationRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStateSDK, mockSegmentsSDK, restore := setupSegmentRealizationMocks(t, ctrl)
	defer restore()

	t.Run("success", func(t *testing.T) {
		successState := model.SegmentConfigurationState_STATE_SUCCESS
		segName := "seg-rlz-1-name"

		mockStateSDK.EXPECT().Get(
			"seg-rlz-1", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(model.SegmentConfigurationState{State: &successState}, nil)

		mockSegmentsSDK.EXPECT().Get("seg-rlz-1").Return(model.Segment{DisplayName: &segName}, nil)

		ds := dataSourceNsxtPolicySegmentRealization()
		d := schema.TestResourceDataRaw(t, ds.Schema, segmentRealizationTestData())

		err := dataSourceNsxtPolicySegmentRealizationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, successState, d.Get("state"))
		assert.Equal(t, segName, d.Get("network_name"))
	})

	t.Run("state API error", func(t *testing.T) {
		mockStateSDK.EXPECT().Get(
			"seg-rlz-1", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(model.SegmentConfigurationState{}, errors.New("state failed"))

		ds := dataSourceNsxtPolicySegmentRealization()
		d := schema.TestResourceDataRaw(t, ds.Schema, segmentRealizationTestData())

		err := dataSourceNsxtPolicySegmentRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to get realization information")
	})

	t.Run("segment lookup error after realization success", func(t *testing.T) {
		successState := model.SegmentConfigurationState_STATE_SUCCESS

		mockStateSDK.EXPECT().Get(
			"seg-rlz-1", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(model.SegmentConfigurationState{State: &successState}, nil)

		mockSegmentsSDK.EXPECT().Get("seg-rlz-1").Return(model.Segment{}, errors.New("get segment failed"))

		ds := dataSourceNsxtPolicySegmentRealization()
		d := schema.TestResourceDataRaw(t, ds.Schema, segmentRealizationTestData())

		err := dataSourceNsxtPolicySegmentRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "get segment failed")
	})
}
