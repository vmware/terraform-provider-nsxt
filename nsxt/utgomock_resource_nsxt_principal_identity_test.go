//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run:
// mockgen -destination=mocks/nsx/trust_management/principal_identities/WithCertificateClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/trust_management/principal_identities/WithCertificateClient.go WithCertificateClient
// mockgen -destination=mocks/nsx/trust_management/PrincipalIdentitiesClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/trust_management/PrincipalIdentitiesClient.go PrincipalIdentitiesClient
// mockgen -destination=mocks/nsx/trust_management/CertificatesClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/trust_management/CertificatesClient.go CertificatesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	mpModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/trust_management"

	tmWrapper "github.com/vmware/terraform-provider-nsxt/api/nsx/trust_management"
	certPIWrapper "github.com/vmware/terraform-provider-nsxt/api/nsx/trust_management/principal_identities"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	certsmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/trust_management"
	pimocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/trust_management/principal_identities"
)

var (
	piID      = "pi-id"
	piName    = "principalname"
	piNodeID  = "nodeidentity1"
	piCertID  = "cert-id"
	piCertPem = "-----BEGIN CERTIFICATE-----\nMIIBxTCCAW+gAwIBAgIUTest\n-----END CERTIFICATE-----"
)

func piAPIResponse() mpModel.PrincipalIdentity {
	return mpModel.PrincipalIdentity{
		Id:            &piID,
		Name:          &piName,
		NodeId:        &piNodeID,
		CertificateId: &piCertID,
	}
}

func minimalPrincipalIdentityData() map[string]interface{} {
	return map[string]interface{}{
		"name":            piName,
		"node_id":         piNodeID,
		"certificate_pem": piCertPem,
		"is_protected":    true,
	}
}

func TestMockResourceNsxtPrincipalIdentityCreate(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Create success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockWithCert := pimocks.NewMockWithCertificateClient(ctrl)
		mockPIWrapper := &certPIWrapper.CertPrincipalIdentityClientContext{
			Client:     mockWithCert,
			ClientType: utl.Local,
		}
		originalWithCert := cliPrincipalIdentityWithCertificateClient
		cliPrincipalIdentityWithCertificateClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *certPIWrapper.CertPrincipalIdentityClientContext {
			return mockPIWrapper
		}
		defer func() { cliPrincipalIdentityWithCertificateClient = originalWithCert }()

		mockPIClient := certsmocks.NewMockPrincipalIdentitiesClient(ctrl)
		mockPIWrapperRead := &tmWrapper.PrincipalIdentityClientContext{
			Client:     mockPIClient,
			ClientType: utl.Local,
		}
		originalPIClient := cliPrincipalIdentitiesClient
		cliPrincipalIdentitiesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tmWrapper.PrincipalIdentityClientContext {
			return mockPIWrapperRead
		}
		defer func() { cliPrincipalIdentitiesClient = originalPIClient }()

		gomock.InOrder(
			mockWithCert.EXPECT().Create(gomock.Any()).Return(piAPIResponse(), nil),
			mockPIClient.EXPECT().Get(piID).Return(piAPIResponse(), nil),
		)

		res := resourceNsxtPrincipalIdentity()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrincipalIdentityData())

		err := resourceNsxtPrincipalIdentityCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, piID, d.Id())
		assert.Equal(t, piName, d.Get("name"))
	})

	t.Run("Create API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockWithCert := pimocks.NewMockWithCertificateClient(ctrl)
		mockPIWrapper := &certPIWrapper.CertPrincipalIdentityClientContext{
			Client:     mockWithCert,
			ClientType: utl.Local,
		}
		originalWithCert := cliPrincipalIdentityWithCertificateClient
		cliPrincipalIdentityWithCertificateClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *certPIWrapper.CertPrincipalIdentityClientContext {
			return mockPIWrapper
		}
		defer func() { cliPrincipalIdentityWithCertificateClient = originalWithCert }()

		mockWithCert.EXPECT().Create(gomock.Any()).Return(mpModel.PrincipalIdentity{}, errors.New("API error"))

		res := resourceNsxtPrincipalIdentity()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrincipalIdentityData())

		err := resourceNsxtPrincipalIdentityCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPrincipalIdentityRead(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPIClient := certsmocks.NewMockPrincipalIdentitiesClient(ctrl)
		mockPIWrapper := &tmWrapper.PrincipalIdentityClientContext{
			Client:     mockPIClient,
			ClientType: utl.Local,
		}
		original := cliPrincipalIdentitiesClient
		cliPrincipalIdentitiesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tmWrapper.PrincipalIdentityClientContext {
			return mockPIWrapper
		}
		defer func() { cliPrincipalIdentitiesClient = original }()

		mockPIClient.EXPECT().Get(piID).Return(piAPIResponse(), nil)

		res := resourceNsxtPrincipalIdentity()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrincipalIdentityData())
		d.SetId(piID)

		err := resourceNsxtPrincipalIdentityRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, piName, d.Get("name"))
		assert.Equal(t, piNodeID, d.Get("node_id"))
		assert.Equal(t, piCertID, d.Get("certificate_id"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPIClient := certsmocks.NewMockPrincipalIdentitiesClient(ctrl)
		mockPIWrapper := &tmWrapper.PrincipalIdentityClientContext{
			Client:     mockPIClient,
			ClientType: utl.Local,
		}
		original := cliPrincipalIdentitiesClient
		cliPrincipalIdentitiesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tmWrapper.PrincipalIdentityClientContext {
			return mockPIWrapper
		}
		defer func() { cliPrincipalIdentitiesClient = original }()

		notFoundErr := vapiErrors.NotFound{}
		mockPIClient.EXPECT().Get(piID).Return(mpModel.PrincipalIdentity{}, notFoundErr)

		res := resourceNsxtPrincipalIdentity()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrincipalIdentityData())
		d.SetId(piID)

		err := resourceNsxtPrincipalIdentityRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPrincipalIdentity()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrincipalIdentityData())

		err := resourceNsxtPrincipalIdentityRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPrincipalIdentityDelete(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Delete success without cert cleanup", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMpPI := certsmocks.NewMockPrincipalIdentitiesClient(ctrl)
		originalMpPI := cliMpPrincipalIdentitiesClient
		cliMpPrincipalIdentitiesClient = func(_ vapiProtocolClient.Connector) trust_management.PrincipalIdentitiesClient {
			return mockMpPI
		}
		defer func() { cliMpPrincipalIdentitiesClient = originalMpPI }()

		mockMpPI.EXPECT().Delete(piID).Return(nil)

		res := resourceNsxtPrincipalIdentity()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrincipalIdentityData())
		d.SetId(piID)

		err := resourceNsxtPrincipalIdentityDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete success with cert cleanup", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMpPI := certsmocks.NewMockPrincipalIdentitiesClient(ctrl)
		originalMpPI := cliMpPrincipalIdentitiesClient
		cliMpPrincipalIdentitiesClient = func(_ vapiProtocolClient.Connector) trust_management.PrincipalIdentitiesClient {
			return mockMpPI
		}
		defer func() { cliMpPrincipalIdentitiesClient = originalMpPI }()

		mockCerts := certsmocks.NewMockCertificatesClient(ctrl)
		mockCertsWrapper := &tmWrapper.CertificateClientContext{
			Client:     mockCerts,
			ClientType: utl.Local,
		}
		originalCerts := cliCertificatesClient
		cliCertificatesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tmWrapper.CertificateClientContext {
			return mockCertsWrapper
		}
		defer func() { cliCertificatesClient = originalCerts }()

		gomock.InOrder(
			mockMpPI.EXPECT().Delete(piID).Return(nil),
			mockCerts.EXPECT().Delete(piCertID).Return(nil),
		)

		res := resourceNsxtPrincipalIdentity()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrincipalIdentityData())
		d.SetId(piID)
		d.Set("certificate_id", piCertID)

		err := resourceNsxtPrincipalIdentityDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete PI fails but continues with cert cleanup", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMpPI := certsmocks.NewMockPrincipalIdentitiesClient(ctrl)
		originalMpPI := cliMpPrincipalIdentitiesClient
		cliMpPrincipalIdentitiesClient = func(_ vapiProtocolClient.Connector) trust_management.PrincipalIdentitiesClient {
			return mockMpPI
		}
		defer func() { cliMpPrincipalIdentitiesClient = originalMpPI }()

		mockCerts := certsmocks.NewMockCertificatesClient(ctrl)
		mockCertsWrapper := &tmWrapper.CertificateClientContext{
			Client:     mockCerts,
			ClientType: utl.Local,
		}
		originalCerts := cliCertificatesClient
		cliCertificatesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tmWrapper.CertificateClientContext {
			return mockCertsWrapper
		}
		defer func() { cliCertificatesClient = originalCerts }()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockMpPI.EXPECT().Delete(piID).Return(notFoundErr),
			mockCerts.EXPECT().Delete(piCertID).Return(nil),
		)

		res := resourceNsxtPrincipalIdentity()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrincipalIdentityData())
		d.SetId(piID)
		d.Set("certificate_id", piCertID)

		err := resourceNsxtPrincipalIdentityDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPrincipalIdentity()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrincipalIdentityData())

		err := resourceNsxtPrincipalIdentityDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
