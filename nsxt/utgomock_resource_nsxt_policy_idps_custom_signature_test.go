//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run:
// mockgen -destination=mocks/infra/settings/firewall/security/intrusion_services/custom_signature_versions/CustomSignaturesClient.go -package=mocks -source=<sdk>/services/nsxt/infra/settings/firewall/security/intrusion_services/custom_signature_versions/CustomSignaturesClient.go CustomSignaturesClient
// mockgen -destination=mocks/infra/settings/firewall/security/intrusion_services/CustomSignatureVersionsClient.go -package=mocks -source=<sdk>/services/nsxt/infra/settings/firewall/security/intrusion_services/CustomSignatureVersionsClient.go CustomSignatureVersionsClient

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services/custom_signature_versions"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	versionmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/settings/firewall/security/intrusion_services"
	custsigmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/settings/firewall/security/intrusion_services/custom_signature_versions"
)

// resourceNsxtPolicyIdpsCustomSignature talks to the SDK-generated clients via the
// cliIdsCustomSignaturesClient / cliIdsCustomSignatureVersionsClient package-level vars
// (matching the swap-for-testing pattern used throughout this package), so the helper
// functions below are mockable via setupIdpsCustomSignatureMocks. Guard conditions,
// composite ID parsing, and schema validation are covered separately below.

func setupIdpsCustomSignatureMocks(ctrl *gomock.Controller) (*custsigmocks.MockCustomSignaturesClient, *versionmocks.MockCustomSignatureVersionsClient, func()) {
	mockSigs := custsigmocks.NewMockCustomSignaturesClient(ctrl)
	mockVersions := versionmocks.NewMockCustomSignatureVersionsClient(ctrl)

	origSigs := cliIdsCustomSignaturesClient
	cliIdsCustomSignaturesClient = func(_ client.Connector) custom_signature_versions.CustomSignaturesClient {
		return mockSigs
	}
	origVersions := cliIdsCustomSignatureVersionsClient
	cliIdsCustomSignatureVersionsClient = func(_ client.Connector) intrusion_services.CustomSignatureVersionsClient {
		return mockVersions
	}
	return mockSigs, mockVersions, func() {
		cliIdsCustomSignaturesClient = origSigs
		cliIdsCustomSignatureVersionsClient = origVersions
	}
}

func minimalIdpsCustomSigData() map[string]interface{} {
	return map[string]interface{}{
		"signature_version_id": "default",
		"signature":            `alert tcp any any -> any 80 (msg:"Test"; sid:9000001; rev:1;)`,
	}
}

func TestMockResourceNsxtPolicyIdpsCustomSignatureGuard(t *testing.T) {
	t.Run("Create fails for global manager", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsCustomSignature()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsCustomSigData())

		err := resourceNsxtPolicyIdpsCustomSignatureCreate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})

	t.Run("Read fails for global manager", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsCustomSignature()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsCustomSigData())
		d.SetId("default/sig-1")

		err := resourceNsxtPolicyIdpsCustomSignatureRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})

	t.Run("Update fails for global manager", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsCustomSignature()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsCustomSigData())
		d.SetId("default/sig-1")

		err := resourceNsxtPolicyIdpsCustomSignatureUpdate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})

	t.Run("Delete fails for global manager", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsCustomSignature()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsCustomSigData())
		d.SetId("default/sig-1")

		err := resourceNsxtPolicyIdpsCustomSignatureDelete(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})
}

func TestMockResourceNsxtPolicyIdpsCustomSignatureIDParsing(t *testing.T) {
	t.Run("Composite ID parsed correctly", func(t *testing.T) {
		versionID, sigID, err := parseCustomSignatureCompositeID("myversion/sig-123")
		require.NoError(t, err)
		assert.Equal(t, "myversion", versionID)
		assert.Equal(t, "sig-123", sigID)
	})

	t.Run("Invalid composite ID returns error", func(t *testing.T) {
		_, _, err := parseCustomSignatureCompositeID("noseparator")
		require.Error(t, err)
	})

	t.Run("Legacy ID (no slash) treated as default version", func(t *testing.T) {
		versionID, sigID, legacy, err := parseCustomSignatureCompositeIDOrLegacy("legacy-sig-id")
		require.NoError(t, err)
		assert.True(t, legacy)
		assert.Equal(t, "default", versionID)
		assert.Equal(t, "legacy-sig-id", sigID)
	})

	t.Run("Empty ID returns error", func(t *testing.T) {
		_, _, _, err := parseCustomSignatureCompositeIDOrLegacy("")
		require.Error(t, err)
	})
}

func TestUnitNsxt_resourceNsxtPolicyIdpsCustomSignatureImportState(t *testing.T) {
	res := resourceNsxtPolicyIdpsCustomSignature()

	t.Run("valid composite ID sets signature_version_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("myversion/sig-1")
		out, err := resourceNsxtPolicyIdpsCustomSignatureImportState(d, nil)
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, "myversion", d.Get("signature_version_id"))
	})

	t.Run("invalid ID errors", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("no-slash")
		_, err := resourceNsxtPolicyIdpsCustomSignatureImportState(d, nil)
		require.Error(t, err)
	})
}

func TestUnitNsxt_idpsCustomSignaturePathSegment(t *testing.T) {
	t.Run("prefers Id", func(t *testing.T) {
		id := "sig-1"
		sig := &model.IdsCustomSignature{Id: &id}
		assert.Equal(t, "sig-1", idpsCustomSignaturePathSegment(sig))
	})

	t.Run("falls back to last path segment", func(t *testing.T) {
		path := "/infra/settings/firewall/security/intrusion-services/custom-signature-versions/default/signatures-preview/sig-2"
		sig := &model.IdsCustomSignature{Path: &path}
		assert.Equal(t, "sig-2", idpsCustomSignaturePathSegment(sig))
	})

	t.Run("empty when neither is set", func(t *testing.T) {
		assert.Equal(t, "", idpsCustomSignaturePathSegment(&model.IdsCustomSignature{}))
	})
}

func TestMockNsxtResourceNsxtPolicyIdpsCustomSignatureExists(t *testing.T) {
	t.Run("Get succeeds means it exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		mockSigs.EXPECT().Get("default", "sig-1").Return(model.IdsCustomSignature{}, nil)

		exists, err := resourceNsxtPolicyIdpsCustomSignatureExists(nil, "default/sig-1")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("NotFound falls back to list search", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		mockSigs.EXPECT().Get("default", "sig-1").Return(model.IdsCustomSignature{}, vapiErrors.NotFound{})
		sigID := "sig-1"
		mockSigs.EXPECT().List("default", gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.IdsCustomSignatureListResult{Results: []model.IdsCustomSignature{{Id: &sigID}}}, nil,
		).AnyTimes()

		exists, err := resourceNsxtPolicyIdpsCustomSignatureExists(nil, "default/sig-1")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("other errors propagate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		mockSigs.EXPECT().Get("default", "sig-1").Return(model.IdsCustomSignature{}, vapiErrors.InternalServerError{})

		_, err := resourceNsxtPolicyIdpsCustomSignatureExists(nil, "default/sig-1")
		require.Error(t, err)
	})

	t.Run("invalid composite id errors", func(t *testing.T) {
		_, err := resourceNsxtPolicyIdpsCustomSignatureExists(nil, "")
		require.Error(t, err)
	})
}

func TestMockNsxtIdpsCustomSignatureValidate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, mockVersions, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		rev := int64(3)
		mockVersions.EXPECT().Get("default").Return(model.IdsCustomSignatureVersion{Revision: &rev}, nil)
		mockSigs.EXPECT().Create("default", gomock.Any(), idpsCustomSignatureActionValidate).Return(nil)

		err := idpsCustomSignatureValidate(nil, "default")
		require.NoError(t, err)
	})

	t.Run("version Get error propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, mockVersions, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		mockVersions.EXPECT().Get("default").Return(model.IdsCustomSignatureVersion{}, vapiErrors.InternalServerError{})

		err := idpsCustomSignatureValidate(nil, "default")
		require.Error(t, err)
	})
}

func TestMockNsxtIdpsCustomSignatureValidateAndPublish(t *testing.T) {
	t.Run("success validates then publishes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, mockVersions, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		rev := int64(3)
		mockVersions.EXPECT().Get("default").Return(model.IdsCustomSignatureVersion{Revision: &rev}, nil).Times(2)
		mockSigs.EXPECT().Create("default", gomock.Any(), idpsCustomSignatureActionValidate).Return(nil)
		mockSigs.EXPECT().Create("default", gomock.Any(), idpsCustomSignatureActionPublish).Return(nil)

		err := idpsCustomSignatureValidateAndPublish(nil, "default")
		require.NoError(t, err)
	})

	t.Run("validate failure short-circuits before publish", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, mockVersions, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		rev := int64(3)
		mockVersions.EXPECT().Get("default").Return(model.IdsCustomSignatureVersion{Revision: &rev}, nil)
		mockSigs.EXPECT().Create("default", gomock.Any(), idpsCustomSignatureActionValidate).Return(vapiErrors.InternalServerError{})

		err := idpsCustomSignatureValidateAndPublish(nil, "default")
		require.Error(t, err)
	})
}

func TestMockNsxtIdpsCustomSignatureFindByIDAndContent(t *testing.T) {
	t.Run("returns first match by id when no content filter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		sigID := "sig-1"
		mockSigs.EXPECT().List("default", gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.IdsCustomSignatureListResult{Results: []model.IdsCustomSignature{{Id: &sigID}}}, nil,
		).AnyTimes()

		found := idpsCustomSignatureFindByIDAndContent(nil, "default", "sig-1", "")
		require.NotNil(t, found)
		assert.Equal(t, "sig-1", *found.Id)
	})

	t.Run("returns match by content when a content filter is given", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		sigID := "sig-1"
		content := "alert tcp"
		mockSigs.EXPECT().List("default", gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.IdsCustomSignatureListResult{Results: []model.IdsCustomSignature{{Id: &sigID, OriginalSignature: &content}}}, nil,
		).AnyTimes()

		found := idpsCustomSignatureFindByIDAndContent(nil, "default", "sig-1", "alert tcp")
		require.NotNil(t, found)
	})

	t.Run("List errors are skipped across include values; no match returns nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		mockSigs.EXPECT().List("default", gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.IdsCustomSignatureListResult{}, vapiErrors.InternalServerError{},
		).AnyTimes()

		found := idpsCustomSignatureFindByIDAndContent(nil, "default", "sig-1", "")
		assert.Nil(t, found)
	})
}

func TestMockNsxtIdpsCustomSignatureFindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
	defer restore()
	sigID := "sig-1"
	mockSigs.EXPECT().List("default", gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
		model.IdsCustomSignatureListResult{Results: []model.IdsCustomSignature{{Id: &sigID}}}, nil,
	).AnyTimes()

	found, err := idpsCustomSignatureFindByID(nil, "default", "sig-1")
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, "sig-1", *found.Id)
}

func TestMockNsxtIdpsCustomSignatureCountPreview(t *testing.T) {
	t.Run("counts results", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		sigID1, sigID2 := "sig-1", "sig-2"
		mockSigs.EXPECT().List("default", gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.IdsCustomSignatureListResult{Results: []model.IdsCustomSignature{{Id: &sigID1}, {Id: &sigID2}}}, nil,
		)

		count := idpsCustomSignatureCountPreview(nil, "default")
		assert.Equal(t, 2, count)
	})

	t.Run("error returns zero", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		mockSigs.EXPECT().List("default", gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.IdsCustomSignatureListResult{}, vapiErrors.InternalServerError{},
		)

		count := idpsCustomSignatureCountPreview(nil, "default")
		assert.Equal(t, 0, count)
	})
}

func TestMockNsxtResourceNsxtPolicyIdpsCustomSignatureExistsByContent(t *testing.T) {
	t.Run("finds matching content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		content := "alert tcp any any"
		mockSigs.EXPECT().List("default", gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.IdsCustomSignatureListResult{Results: []model.IdsCustomSignature{{OriginalSignature: &content}}}, nil,
		).AnyTimes()

		found, err := ResourceNsxtPolicyIdpsCustomSignatureExistsByContent(nil, "default", "alert tcp")
		require.NoError(t, err)
		assert.True(t, found)
	})

	t.Run("no match returns false", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		mockSigs.EXPECT().List("default", gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.IdsCustomSignatureListResult{Results: []model.IdsCustomSignature{}}, nil,
		).AnyTimes()

		found, err := ResourceNsxtPolicyIdpsCustomSignatureExistsByContent(nil, "default", "nomatch")
		require.NoError(t, err)
		assert.False(t, found)
	})
}
