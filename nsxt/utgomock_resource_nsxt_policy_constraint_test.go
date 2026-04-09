//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	constraintID       = "constraint-1"
	constraintName     = "constraint-fooname"
	constraintRevision = int64(1)
	constraintPath     = "/infra/constraints/constraint-1"
)

func TestMockResourceNsxtPolicyConstraintCreate(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConstraintsSDK := inframocks.NewMockConstraintsClient(ctrl)
	constraintWrapper := &cliinfra.ConstraintClientContext{
		Client:     mockConstraintsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliConstraintsClient
	defer func() { cliConstraintsClient = originalCli }()
	cliConstraintsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ConstraintClientContext {
		return constraintWrapper
	}

	res := resourceNsxtPolicyConstraint()

	t.Run("Create_success", func(t *testing.T) {
		mockConstraintsSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockConstraintsSDK.EXPECT().Get(gomock.Any()).Return(model.Constraint{
			DisplayName: &constraintName,
			Path:        &constraintPath,
			Revision:    &constraintRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": constraintName,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConstraintCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, constraintName, d.Get("display_name"))
	})

	t.Run("Create_fails_below_version_9_0_0", func(t *testing.T) {
		util.NsxVersion = "3.2.0"
		defer func() { util.NsxVersion = "9.0.0" }()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": constraintName,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConstraintCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.0.0")
	})

	t.Run("Create_fails_when_already_exists", func(t *testing.T) {
		mockConstraintsSDK.EXPECT().Get(constraintID).Return(model.Constraint{}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"nsx_id":       constraintID,
			"display_name": constraintName,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConstraintCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockConstraintsSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": constraintName,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConstraintCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyConstraintRead(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConstraintsSDK := inframocks.NewMockConstraintsClient(ctrl)
	constraintWrapper := &cliinfra.ConstraintClientContext{
		Client:     mockConstraintsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliConstraintsClient
	defer func() { cliConstraintsClient = originalCli }()
	cliConstraintsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ConstraintClientContext {
		return constraintWrapper
	}

	res := resourceNsxtPolicyConstraint()

	t.Run("Read_success", func(t *testing.T) {
		mockConstraintsSDK.EXPECT().Get(constraintID).Return(model.Constraint{
			DisplayName: &constraintName,
			Path:        &constraintPath,
			Revision:    &constraintRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(constraintID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConstraintRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, constraintName, d.Get("display_name"))
	})

	t.Run("Read_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConstraintRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Constraint ID")
	})

	t.Run("Read_fails_when_Get_returns_not_found", func(t *testing.T) {
		mockConstraintsSDK.EXPECT().Get(constraintID).Return(model.Constraint{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(constraintID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConstraintRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyConstraintUpdate(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConstraintsSDK := inframocks.NewMockConstraintsClient(ctrl)
	constraintWrapper := &cliinfra.ConstraintClientContext{
		Client:     mockConstraintsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliConstraintsClient
	defer func() { cliConstraintsClient = originalCli }()
	cliConstraintsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ConstraintClientContext {
		return constraintWrapper
	}

	res := resourceNsxtPolicyConstraint()

	t.Run("Update_success", func(t *testing.T) {
		mockConstraintsSDK.EXPECT().Update(constraintID, gomock.Any()).Return(model.Constraint{}, nil)
		mockConstraintsSDK.EXPECT().Get(constraintID).Return(model.Constraint{
			DisplayName: &constraintName,
			Path:        &constraintPath,
			Revision:    &constraintRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": constraintName,
		})
		d.SetId(constraintID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConstraintUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConstraintUpdate(d, m)
		require.Error(t, err)
	})

	t.Run("Update_fails_when_Update_returns_error", func(t *testing.T) {
		mockConstraintsSDK.EXPECT().Update(constraintID, gomock.Any()).Return(model.Constraint{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": constraintName,
		})
		d.SetId(constraintID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConstraintUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyConstraintDelete(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConstraintsSDK := inframocks.NewMockConstraintsClient(ctrl)
	constraintWrapper := &cliinfra.ConstraintClientContext{
		Client:     mockConstraintsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliConstraintsClient
	defer func() { cliConstraintsClient = originalCli }()
	cliConstraintsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ConstraintClientContext {
		return constraintWrapper
	}

	res := resourceNsxtPolicyConstraint()

	t.Run("Delete_success", func(t *testing.T) {
		mockConstraintsSDK.EXPECT().Delete(constraintID).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(constraintID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConstraintDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConstraintDelete(d, m)
		require.Error(t, err)
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockConstraintsSDK.EXPECT().Delete(constraintID).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(constraintID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConstraintDelete(d, m)
		require.Error(t, err)
	})
}
