// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"

	sdkinfra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	projectinfra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/infra"
)

func testAccTlsCsrWithDaysValidSpec(displayName string) model.TlsCsrWithDaysValid {
	algo := model.TlsCsrWithDaysValid_ALGORITHM_RSA
	keySize := int64(2048)
	days := int64(365)
	isCa := false
	cn := displayName
	return model.TlsCsrWithDaysValid{
		DisplayName: &displayName,
		Algorithm:   &algo,
		KeySize:     &keySize,
		DaysValid:   &days,
		IsCa:        &isCa,
		Subject: &model.Principal{
			Attributes: []model.KeyValue{
				{Key: ptr("CN"), Value: &cn},
			},
		},
	}
}

func testAccTlsCertificateSelfSignedCreateGlobal(displayName string) (string, error) {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return "", err
	}
	cert, err := sdkinfra.NewCsrsClient(connector).Selfsign0(testAccTlsCsrWithDaysValidSpec(displayName))
	if err != nil {
		return "", err
	}
	if cert.Id == nil || *cert.Id == "" {
		return "", fmt.Errorf("tls certificate Selfsign0 returned empty id")
	}
	return *cert.Id, nil
}

func testAccTlsCertificateSelfSignedCreateProject(displayName string) (string, error) {
	projectID := testAccMultitenancyProjectID()
	if projectID == "" {
		return "", fmt.Errorf("project id is not set (NSXT_PROJECT_ID or NSXT_VPC_PROJECT_ID)")
	}
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return "", err
	}
	cert, err := projectinfra.NewCsrsClient(connector).Selfsign0(defaultOrgID, projectID, testAccTlsCsrWithDaysValidSpec(displayName))
	if err != nil {
		return "", err
	}
	if cert.Id == nil || *cert.Id == "" {
		return "", fmt.Errorf("project tls certificate Selfsign0 returned empty id")
	}
	return *cert.Id, nil
}

func testAccMultitenancyProjectID() string {
	if testAccIsVPC() {
		return os.Getenv("NSXT_VPC_PROJECT_ID")
	}
	return os.Getenv("NSXT_PROJECT_ID")
}

func testAccTlsCertificateDeleteGlobal(certID string) error {
	if certID == "" {
		return nil
	}
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return err
	}
	return sdkinfra.NewCertificatesClient(connector).Delete(certID)
}

func testAccTlsCertificateDeleteProject(certID string) error {
	if certID == "" {
		return nil
	}
	projectID := testAccMultitenancyProjectID()
	if projectID == "" {
		return nil
	}
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return err
	}
	return projectinfra.NewCertificatesClient(connector).Delete(defaultOrgID, projectID, certID)
}
