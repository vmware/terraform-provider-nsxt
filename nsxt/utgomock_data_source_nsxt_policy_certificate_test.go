//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	gmModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestUnitNsxt_dataSourceNsxtPolicyCertificateRead(t *testing.T) {
	rt := "TlsCertificate"

	t.Run("success by id", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("cert-1"), DisplayName: str("cert-name"), Path: str("/infra/certificates/cert-1"), ResourceType: &rt,
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyCertificate()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "cert-1"})

		err := dataSourceNsxtPolicyCertificateRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "cert-1", d.Id())
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search-fail")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyCertificate()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "missing"})

		err := dataSourceNsxtPolicyCertificateRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
