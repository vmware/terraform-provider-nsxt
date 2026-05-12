// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services/signature_versions"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyIdpsSystemSignatures() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIdpsSystemSignaturesRead,

		Schema: map[string]*schema.Schema{
			"id":   getDataSourceIDSchema(),
			"path": getPathSchema(),
			"version_id": {
				Type:        schema.TypeString,
				Description: "Signature version ID to query. If not specified, uses the active version",
				Optional:    true,
				Computed:    true,
			},
			"severity": {
				Type:         schema.TypeString,
				Description:  "Filter signatures by severity level",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"LOW", "MEDIUM", "HIGH", "CRITICAL"}, false),
			},
			"product_affected": {
				Type:        schema.TypeString,
				Description: "Filter signatures by affected product",
				Optional:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "Filter signatures by display name pattern (case-insensitive substring match)",
				Optional:    true,
			},
			"class_type": {
				Type:        schema.TypeString,
				Description: "Filter signatures by class type",
				Optional:    true,
			},
			"signatures": {
				Type:        schema.TypeList,
				Description: "List of system signatures matching the filter criteria",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Signature ID",
							Computed:    true,
						},
						"signature_id": {
							Type:        schema.TypeString,
							Description: "Unique signature identifier",
							Computed:    true,
						},
						"display_name": {
							Type:        schema.TypeString,
							Description: "Display name of the signature",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the signature",
							Computed:    true,
						},
						"severity": {
							Type:        schema.TypeString,
							Description: "Severity level of the signature",
							Computed:    true,
						},
						"class_type": {
							Type:        schema.TypeString,
							Description: "Class type of the signature",
							Computed:    true,
						},
						"categories": {
							Type:        schema.TypeList,
							Description: "Categories of the signature",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"product_affected": {
							Type:        schema.TypeString,
							Description: "Product affected by this signature",
							Computed:    true,
						},
						"attack_target": {
							Type:        schema.TypeString,
							Description: "Attack target of the signature",
							Computed:    true,
						},
						"cvss": {
							Type:        schema.TypeString,
							Description: "CVSS severity rating",
							Computed:    true,
						},
						"cvssv2": {
							Type:        schema.TypeString,
							Description: "CVSS v2 score",
							Computed:    true,
						},
						"cvssv3": {
							Type:        schema.TypeString,
							Description: "CVSS v3 score",
							Computed:    true,
						},
						"cves": {
							Type:        schema.TypeList,
							Description: "List of CVE IDs associated with this signature",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"urls": {
							Type:        schema.TypeList,
							Description: "List of reference URLs for this signature",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"signature_revision": {
							Type:        schema.TypeString,
							Description: "Signature revision number",
							Computed:    true,
						},
						"path": {
							Type:        schema.TypeString,
							Description: "Policy path of the signature",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNsxtPolicyIdpsSystemSignaturesRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	if isPolicyGlobalManager(m) {
		return localManagerOnlyError()
	}

	// Get filter parameters
	versionID := d.Get("version_id").(string)
	severityFilter := d.Get("severity").(string)
	productFilter := d.Get("product_affected").(string)
	displayNameFilter := d.Get("display_name").(string)
	classTypeFilter := d.Get("class_type").(string)

	// If version_id not specified, get the active version
	if versionID == "" {
		versionsClient := intrusion_services.NewSignatureVersionsClient(connector)
		if versionsClient == nil {
			return policyResourceNotSupportedError()
		}

		versionsList, err := versionsClient.List(nil, nil, nil, nil, nil, nil)
		if err != nil {
			return handleListError("IdsSignatureVersion", err)
		}

		// Find the active version
		var activeVersion *model.IdsSignatureVersion
		for _, version := range versionsList.Results {
			if version.State != nil && *version.State == "ACTIVE" {
				activeVersion = &version
				break
			}
		}

		// If no active version found, try to use the first version or DEFAULT
		if activeVersion == nil {
			if len(versionsList.Results) > 0 {
				activeVersion = &versionsList.Results[0]
			} else {
				// Default to "DEFAULT" version as fallback
				versionID = "DEFAULT"
			}
		}

		if activeVersion != nil && activeVersion.Id != nil {
			versionID = *activeVersion.Id
		}

		if versionID == "" {
			return fmt.Errorf("No signature version found")
		}
	}

	// Set version_id in state
	d.Set("version_id", versionID)

	// Get signatures client for the specific version
	signaturesClient := signature_versions.NewSignaturesClient(connector)
	if signaturesClient == nil {
		return policyResourceNotSupportedError()
	}

	// List signatures for the version with pagination
	// Fetch all pages automatically without exposing pagination to users
	var allSignatures []model.IdsSignature
	var cursor *string
	total := 0

	for {
		signaturesList, err := signaturesClient.List(versionID, cursor, nil, nil, nil, nil, nil)
		if err != nil {
			return handleListError("IdsSignature", err)
		}

		allSignatures = append(allSignatures, signaturesList.Results...)

		if total == 0 && signaturesList.ResultCount != nil {
			// First response
			total = int(*signaturesList.ResultCount)
		}

		cursor = signaturesList.Cursor
		if len(allSignatures) >= total {
			break
		}
	}

	// Apply filters
	var filteredSignatures []model.IdsSignature
	for _, sig := range allSignatures {
		// Apply severity filter
		if severityFilter != "" && sig.Severity != nil && *sig.Severity != severityFilter {
			continue
		}

		// Apply product_affected filter
		if productFilter != "" && sig.ProductAffected != nil && *sig.ProductAffected != productFilter {
			continue
		}

		// Apply display_name filter (case-insensitive substring match)
		if displayNameFilter != "" {
			if sig.DisplayName == nil {
				continue
			}
			if !strings.Contains(strings.ToLower(*sig.DisplayName), strings.ToLower(displayNameFilter)) {
				continue
			}
		}

		// Apply class_type filter
		if classTypeFilter != "" && sig.ClassType != nil && *sig.ClassType != classTypeFilter {
			continue
		}

		filteredSignatures = append(filteredSignatures, sig)
	}

	// Convert to schema format
	signatures := make([]map[string]interface{}, 0, len(filteredSignatures))
	for _, sig := range filteredSignatures {
		sigMap := make(map[string]interface{})

		sigMap["id"] = sig.Id
		sigMap["signature_id"] = sig.SignatureId
		sigMap["display_name"] = sig.DisplayName
		sigMap["name"] = sig.Name
		sigMap["severity"] = sig.Severity
		sigMap["class_type"] = sig.ClassType
		sigMap["product_affected"] = sig.ProductAffected
		sigMap["attack_target"] = sig.AttackTarget
		sigMap["cvss"] = sig.Cvss
		sigMap["cvssv2"] = sig.Cvssv2
		sigMap["cvssv3"] = sig.Cvssv3
		sigMap["signature_revision"] = sig.SignatureRevision
		sigMap["path"] = sig.Path
		sigMap["categories"] = sig.Categories
		sigMap["cves"] = sig.Cves
		sigMap["urls"] = sig.Urls

		signatures = append(signatures, sigMap)
	}

	// Set a unique ID based on version and filters
	idParts := []string{versionID}
	if severityFilter != "" {
		idParts = append(idParts, "severity", severityFilter)
	}
	if productFilter != "" {
		idParts = append(idParts, "product", productFilter)
	}
	if displayNameFilter != "" {
		idParts = append(idParts, "name", displayNameFilter)
	}
	if classTypeFilter != "" {
		idParts = append(idParts, "class", classTypeFilter)
	}
	d.SetId(strings.Join(idParts, "-"))

	// Set the path
	d.Set("path", fmt.Sprintf("/infra/settings/firewall/security/intrusion-services/signature-versions/%s/signatures", versionID))

	// Set signatures in state
	return d.Set("signatures", signatures)
}
