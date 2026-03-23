// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services/custom_signature_versions"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

const (
	idpsCustomSignatureCreateActionAdd = "ADD_CUSTOM_SIGNATURES"
	idpsCustomSignatureActionValidate  = "VALIDATE"
	idpsCustomSignatureActionPublish   = "PUBLISH"
	idpsCustomSignatureActionCancel    = "CANCEL"
	// Full policy path prefix for custom signature versions (UI uses this in deleted_signatures for VALIDATE).
	idpsCustomSignatureVersionsPathPrefix = "/infra/settings/firewall/security/intrusion-services/custom-signature-versions"
)

func resourceNsxtPolicyIdpsCustomSignature() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIdpsCustomSignatureCreate,
		Read:   resourceNsxtPolicyIdpsCustomSignatureRead,
		Update: resourceNsxtPolicyIdpsCustomSignatureUpdate,
		Delete: resourceNsxtPolicyIdpsCustomSignatureDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyIdpsCustomSignatureImportState,
		},

		Schema: map[string]*schema.Schema{
			"signature_version_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the custom signature version (e.g. \"default\") to which this signature belongs.",
				ForceNew:    true,
			},
			"signature": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Raw custom signature rule in Snort syntax.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the custom signature.",
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the signature is enabled.",
			},
			"publish": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, validate and publish the signature so it becomes active in IDS/IPS. Default is false (signature remains in preview).",
			},
			// Computed
			"path": {
				Type:        schema.TypeString,
				Description: "Absolute path of the custom signature.",
				Computed:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "Display name (derived from signature).",
				Computed:    true,
			},
			"revision": {
				Type:        schema.TypeInt,
				Description: "Revision of the custom signature.",
				Computed:    true,
			},
			"signature_id": {
				Type:        schema.TypeString,
				Description: "System-assigned signature ID (e.g. SID in Snort).",
				Computed:    true,
			},
			"original_signature_id": {
				Type:        schema.TypeString,
				Description: "Original signature ID from the rule (sid in Snort).",
				Computed:    true,
			},
			"original_signature": {
				Type:        schema.TypeString,
				Description: "Original raw signature as stored.",
				Computed:    true,
			},
			"validation_status": {
				Type:        schema.TypeString,
				Description: "Validation status: VALID, INVALID, PENDING, WARNING.",
				Computed:    true,
			},
			"validation_message": {
				Type:        schema.TypeString,
				Description: "Validation message or error if invalid.",
				Computed:    true,
			},
			"severity": {
				Type:        schema.TypeString,
				Description: "Severity of the signature.",
				Computed:    true,
			},
			"signature_severity": {
				Type:        schema.TypeString,
				Description: "Vendor-assigned signature severity.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Signature name (from msg in rule).",
				Computed:    true,
			},
			"signature_revision": {
				Type:        schema.TypeString,
				Description: "Signature revision from the rule.",
				Computed:    true,
			},
		},
	}
}

func resourceNsxtPolicyIdpsCustomSignatureImportState(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("import ID must be of the form <signature_version_id>/<signature_id>")
	}
	d.Set("signature_version_id", parts[0])
	d.SetId(d.Id())
	return []*schema.ResourceData{d}, nil
}

func parseCustomSignatureCompositeID(id string) (versionID, sigID string, err error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("custom signature ID must be <version_id>/<signature_id>, got %q", id)
	}
	return parts[0], parts[1], nil
}

// parseCustomSignatureCompositeIDOrLegacy parses id as "version_id/signature_id". If id contains no "/", treats as legacy state (bare signature_id) and returns ("default", id, true). Caller should SetId("default/"+id) when legacy is true to migrate state.
func parseCustomSignatureCompositeIDOrLegacy(id string) (versionID, sigID string, legacy bool, err error) {
	if id == "" {
		return "", "", false, fmt.Errorf("custom signature ID must be <version_id>/<signature_id>, got empty")
	}
	if !strings.Contains(id, "/") {
		return "default", id, true, nil
	}
	versionID, sigID, err = parseCustomSignatureCompositeID(id)
	return versionID, sigID, false, err
}

// resourceNsxtPolicyIdpsCustomSignatureExists checks whether a custom signature exists (Get or list fallback for preview).
func resourceNsxtPolicyIdpsCustomSignatureExists(connector client.Connector, id string) (bool, error) {
	versionID, sigID, _, err := parseCustomSignatureCompositeIDOrLegacy(id)
	if err != nil {
		return false, err
	}
	customSigsClient := custom_signature_versions.NewCustomSignaturesClient(connector)
	_, err = customSigsClient.Get(versionID, sigID)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		found, _ := idpsCustomSignatureFindByID(connector, versionID, sigID)
		return found != nil, nil
	}
	return false, err
}

func resourceNsxtPolicyIdpsCustomSignatureCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	if isPolicyGlobalManager(m) {
		return globalManagerOnlyError()
	}

	versionID := d.Get("signature_version_id").(string)
	signature := d.Get("signature").(string)
	if signature == "" {
		return fmt.Errorf("signature is required")
	}

	// Add raw signature via CustomSignatureVersions API
	versionsClient := intrusion_services.NewCustomSignatureVersionsClient(connector)
	signatures := []string{signature}
	rawSigs := model.IdsRawSignatures{
		Signatures: signatures,
	}
	err := versionsClient.Create(versionID, rawSigs, idpsCustomSignatureCreateActionAdd)
	if err != nil {
		return handleCreateError("IdsCustomSignature", versionID, err)
	}

	// List custom signatures (include preview to see newly added) to find the one we just added
	customSigsClient := custom_signature_versions.NewCustomSignaturesClient(connector)
	includePreview := custom_signature_versions.CustomSignatures_LIST_INCLUDE_PREVIEW_CUSTOM_SIGNATURES
	pageSize := int64(1000)
	listResult, err := customSigsClient.List(versionID, nil, &includePreview, nil, nil, &pageSize, nil, nil)
	if err != nil {
		return handleListError("IdsCustomSignature", err)
	}

	// When multiple signatures match the same rule text (e.g. leftover from a previous run), take the last match so we get the one just added.
	var found *model.IdsCustomSignature
	for i := range listResult.Results {
		sig := &listResult.Results[i]
		if sig.OriginalSignature != nil && strings.TrimSpace(*sig.OriginalSignature) == strings.TrimSpace(signature) {
			found = sig
		}
	}
	if found == nil || found.Id == nil {
		return fmt.Errorf("failed to find newly created custom signature in list (version %s)", versionID)
	}

	// Store composite ID: versionID + last path segment (or Id). Data source resolves via Get or list fallback.
	sigID := *found.Id
	if found.Path != nil && *found.Path != "" {
		parts := strings.Split(strings.Trim(*found.Path, "/"), "/")
		if len(parts) > 0 {
			sigID = parts[len(parts)-1]
		}
	}
	compositeID := versionID + "/" + sigID
	d.SetId(compositeID)

	// Optionally validate and publish
	if d.Get("publish").(bool) {
		if err := idpsCustomSignatureValidateAndPublish(connector, versionID); err != nil {
			// Some backends reject VALIDATE (empty) right after ADD; try Publish only
			versionsClient := intrusion_services.NewCustomSignatureVersionsClient(connector)
			versionObj, verr := versionsClient.Get(versionID)
			if verr == nil {
				payload := model.CustomSignatureValidationPayload{}
				if versionObj.Revision != nil {
					payload.Revision = versionObj.Revision
				}
				if pubErr := custom_signature_versions.NewCustomSignaturesClient(connector).Create(versionID, payload, idpsCustomSignatureActionPublish); pubErr == nil {
					// Publish-only succeeded
				} else {
					return fmt.Errorf("failed to validate/publish after create: %w", err)
				}
			} else {
				return fmt.Errorf("failed to validate/publish after create: %w", err)
			}
		}
	} else {
		// Validate only so status is updated
		if err := idpsCustomSignatureValidate(connector, versionID); err != nil {
			log.Printf("[WARN] validate after create failed (signature still created): %v", err)
		}
	}

	return resourceNsxtPolicyIdpsCustomSignatureRead(d, m)
}

func idpsCustomSignatureValidate(connector client.Connector, versionID string) error {
	versionObj, err := intrusion_services.NewCustomSignatureVersionsClient(connector).Get(versionID)
	if err != nil {
		return err
	}
	payload := model.CustomSignatureValidationPayload{}
	if versionObj.Revision != nil {
		payload.Revision = versionObj.Revision
	}
	customSigsClient := custom_signature_versions.NewCustomSignaturesClient(connector)
	return customSigsClient.Create(versionID, payload, idpsCustomSignatureActionValidate)
}

func idpsCustomSignatureValidateAndPublish(connector client.Connector, versionID string) error {
	if err := idpsCustomSignatureValidate(connector, versionID); err != nil {
		return err
	}
	versionObj, err := intrusion_services.NewCustomSignatureVersionsClient(connector).Get(versionID)
	if err != nil {
		return err
	}
	payload := model.CustomSignatureValidationPayload{}
	if versionObj.Revision != nil {
		payload.Revision = versionObj.Revision
	}
	customSigsClient := custom_signature_versions.NewCustomSignaturesClient(connector)
	return customSigsClient.Create(versionID, payload, idpsCustomSignatureActionPublish)
}

func resourceNsxtPolicyIdpsCustomSignatureRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	if isPolicyGlobalManager(m) {
		return globalManagerOnlyError()
	}

	versionID, sigID, legacy, err := parseCustomSignatureCompositeIDOrLegacy(d.Id())
	if err != nil {
		return err
	}
	if legacy {
		d.SetId("default/" + d.Id())
	}

	customSigsClient := custom_signature_versions.NewCustomSignaturesClient(connector)
	obj, err := customSigsClient.Get(versionID, sigID)
	if err != nil {
		if isNotFoundError(err) {
			// Preview signatures may not be gettable by Get(); fall back to listing and find by id/path
			foundObj, findErr := idpsCustomSignatureFindByID(connector, versionID, sigID)
			if findErr != nil {
				return handleReadError(d, "IdsCustomSignature", d.Id(), findErr)
			}
			if foundObj == nil {
				d.SetId("")
				return nil
			}
			obj = *foundObj
		} else {
			return handleReadError(d, "IdsCustomSignature", d.Id(), err)
		}
	}

	// When the backend has multiple signatures with the same path segment (e.g. preview and published),
	// Get/FindByID may return the wrong one. Prefer the signature that matches our state to avoid drift.
	if stateSig, ok := d.Get("signature").(string); ok && stateSig != "" && obj.OriginalSignature != nil {
		if strings.TrimSpace(*obj.OriginalSignature) != strings.TrimSpace(stateSig) {
			if match := idpsCustomSignatureFindByIDAndContent(connector, versionID, sigID, stateSig); match != nil {
				obj = *match
			}
		}
	}

	d.Set("signature_version_id", versionID)
	d.Set("path", obj.Path)
	d.Set("display_name", obj.DisplayName)
	d.Set("revision", obj.Revision)
	d.Set("signature_id", obj.SignatureId)
	d.Set("original_signature_id", obj.OriginalSignatureId)
	d.Set("original_signature", obj.OriginalSignature)
	d.Set("validation_status", obj.ValidationStatus)
	d.Set("validation_message", obj.ValidationMessage)
	d.Set("severity", obj.Severity)
	d.Set("signature_severity", obj.SignatureSeverity)
	d.Set("name", obj.Name)
	// Preview signatures often return enable=false until published; preserve config when false to avoid drift.
	// When not in state (e.g. import), use schema default true so import does not produce a spurious diff.
	if obj.Enable != nil && *obj.Enable {
		d.Set("enable", true)
	} else if v, ok := d.GetOk("enable"); ok {
		d.Set("enable", v)
	} else {
		d.Set("enable", true)
	}
	d.Set("signature_revision", obj.SignatureRevision)
	if obj.OriginalSignature != nil {
		d.Set("signature", *obj.OriginalSignature)
	}
	return nil
}

// idpsCustomSignatureFindByID lists custom signatures (preview and/or published) and returns the first matching sigID (id or path segment).
func idpsCustomSignatureFindByID(connector client.Connector, versionID, sigID string) (*model.IdsCustomSignature, error) {
	found := idpsCustomSignatureFindByIDAndContent(connector, versionID, sigID, "")
	if found == nil {
		return nil, nil
	}
	return found, nil
}

func idpsCustomSignaturePathSegment(sig *model.IdsCustomSignature) string {
	if sig.Id != nil && *sig.Id != "" {
		return *sig.Id
	}
	if sig.Path != nil && *sig.Path != "" {
		parts := strings.Split(strings.Trim(*sig.Path, "/"), "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}
	return ""
}

// idpsCustomSignatureFindByIDAndContent lists custom signatures and returns the one matching sigID and optionally expectedSignature.
// When expectedSignature is non-empty, returns the last match (in list order) with that OriginalSignature to prefer the correct one when
// multiple signatures share the same path segment (e.g. preview vs published). When expectedSignature is empty, returns the first match by sigID.
func idpsCustomSignatureFindByIDAndContent(connector client.Connector, versionID, sigID, expectedSignature string) *model.IdsCustomSignature {
	customSigsClient := custom_signature_versions.NewCustomSignaturesClient(connector)
	pageSize := int64(1000)
	matchByContent := expectedSignature != ""
	expectedTrim := strings.TrimSpace(expectedSignature)
	var found *model.IdsCustomSignature
	for _, include := range []*string{
		ptr(custom_signature_versions.CustomSignatures_LIST_INCLUDE_PREVIEW_CUSTOM_SIGNATURES),
		ptr(custom_signature_versions.CustomSignatures_LIST_INCLUDE_CUSTOM_SIGNATURES),
		nil,
	} {
		listResult, err := customSigsClient.List(versionID, nil, include, nil, nil, &pageSize, nil, nil)
		if err != nil {
			continue
		}
		for i := range listResult.Results {
			sig := &listResult.Results[i]
			if idpsCustomSignaturePathSegment(sig) != sigID {
				continue
			}
			if !matchByContent {
				return sig
			}
			if sig.OriginalSignature != nil && strings.TrimSpace(*sig.OriginalSignature) == expectedTrim {
				found = sig
			}
		}
	}
	return found
}

func ptr(s string) *string { return &s }

func resourceNsxtPolicyIdpsCustomSignatureUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	if isPolicyGlobalManager(m) {
		return globalManagerOnlyError()
	}

	versionID, sigID, legacy, err := parseCustomSignatureCompositeIDOrLegacy(d.Id())
	if err != nil {
		return err
	}
	if legacy {
		d.SetId("default/" + d.Id())
	}

	versionsClient := intrusion_services.NewCustomSignatureVersionsClient(connector)
	customSigsClient := custom_signature_versions.NewCustomSignaturesClient(connector)
	versionObj, err := versionsClient.Get(versionID)
	if err != nil {
		return handleReadError(d, "IdsCustomSignatureVersion", versionID, err)
	}

	// If raw signature changed, we need to send a modification
	if d.HasChange("signature") {
		newSig := d.Get("signature").(string)
		mod := model.CustomSignatureModification{
			SignatureId:  &sigID,
			RawSignature: &newSig,
		}
		payload := model.CustomSignatureValidationPayload{
			ModifiedSignatures: []model.CustomSignatureModification{mod},
		}
		if versionObj.Revision != nil {
			payload.Revision = versionObj.Revision
		}
		if err := customSigsClient.Create(versionID, payload, idpsCustomSignatureActionValidate); err != nil {
			return handleUpdateError("IdsCustomSignature", d.Id(), err)
		}
	}

	if d.Get("publish").(bool) {
		if err := idpsCustomSignatureValidateAndPublish(connector, versionID); err != nil {
			return handleUpdateError("IdsCustomSignature", d.Id(), err)
		}
	}

	return resourceNsxtPolicyIdpsCustomSignatureRead(d, m)
}

func resourceNsxtPolicyIdpsCustomSignatureDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	if isPolicyGlobalManager(m) {
		return globalManagerOnlyError()
	}

	versionID, sigID, legacy, err := parseCustomSignatureCompositeIDOrLegacy(d.Id())
	if err != nil {
		return err
	}
	if legacy {
		d.SetId("default/" + d.Id())
	}

	versionsClient := intrusion_services.NewCustomSignatureVersionsClient(connector)
	customSigsClient := custom_signature_versions.NewCustomSignaturesClient(connector)
	versionObj, err := versionsClient.Get(versionID)
	if err != nil {
		return handleReadError(d, "IdsCustomSignatureVersion", versionID, err)
	}

	// Delete is a two-step process matching the UI: (1) VALIDATE with DeletedSignatures, (2) PUBLISH (empty payload).
	// UI sends deleted_signatures: ["/infra/.../custom-signature-versions/{versionID}/signatures-preview/{sigID}"].
	// Build that full path first so VALIDATE stages the correct deletion.
	fullPreviewPath := fmt.Sprintf("%s/%s/signatures-preview/%s", idpsCustomSignatureVersionsPathPrefix, versionID, sigID)
	deletedRef := fullPreviewPath
	tryPaths := []string{
		fullPreviewPath,
		sigID,
		"/" + sigID,
	}
	if obj, getErr := customSigsClient.Get(versionID, sigID); getErr == nil && obj.Path != nil && *obj.Path != "" {
		tryPaths = append(tryPaths, *obj.Path, strings.TrimPrefix(*obj.Path, "/"))
	} else if foundObj, _ := idpsCustomSignatureFindByID(connector, versionID, sigID); foundObj != nil && foundObj.Path != nil && *foundObj.Path != "" {
		tryPaths = append(tryPaths, *foundObj.Path, strings.TrimPrefix(*foundObj.Path, "/"))
	}
	tryPaths = append(tryPaths, "/signatures-preview/"+sigID, "signatures-preview/"+sigID, "/signatures/"+sigID, "signatures/"+sigID)

	var validateErr error
	for _, path := range tryPaths {
		if path == "" {
			continue
		}
		validatePayload := model.CustomSignatureValidationPayload{
			DeletedSignatures: []string{path},
		}
		if versionObj.Revision != nil {
			validatePayload.Revision = versionObj.Revision
		}
		validateErr = customSigsClient.Create(versionID, validatePayload, idpsCustomSignatureActionValidate)
		if validateErr == nil {
			log.Printf("[DEBUG] IdpsCustomSignature delete: VALIDATE accepted path %q for version %s", path, versionID)
			deletedRef = path
			break
		}
		log.Printf("[DEBUG] IdpsCustomSignature delete: VALIDATE with path %q failed: %v", path, validateErr)
		if !strings.Contains(validateErr.Error(), "invalid") && !strings.Contains(validateErr.Error(), "500012") {
			break
		}
	}
	if validateErr != nil {
		previewCount := idpsCustomSignatureCountPreview(connector, versionID)
		if previewCount > 0 {
			log.Printf("[INFO] Delete via VALIDATE failed; using CANCEL to remove unpublished signatures (version %s)", versionID)
			cancelPayload := model.CustomSignatureValidationPayload{}
			if versionObj.Revision != nil {
				cancelPayload.Revision = versionObj.Revision
			}
			if err = customSigsClient.Create(versionID, cancelPayload, idpsCustomSignatureActionCancel); err != nil {
				return handleDeleteError("IdsCustomSignature", d.Id(), err)
			}
			d.SetId("")
			return nil
		}
		return handleDeleteError("IdsCustomSignature", d.Id(), validateErr)
	}

	// Step 2: Re-fetch version (revision may have changed after VALIDATE) then PUBLISH to apply the staged deletion.
	versionAfterValidate, getErr := versionsClient.Get(versionID)
	if getErr != nil {
		return handleDeleteError("IdsCustomSignature", d.Id(), getErr)
	}
	publishPayload := model.CustomSignatureValidationPayload{}
	if versionAfterValidate.Revision != nil {
		publishPayload.Revision = versionAfterValidate.Revision
	}
	if err := customSigsClient.Create(versionID, publishPayload, idpsCustomSignatureActionPublish); err != nil {
		return handleDeleteError("IdsCustomSignature", d.Id(), err)
	}

	// If still present (e.g. published signature or eventual consistency), retry VALIDATE+PUBLISH once with fresh revision.
	if stillFound, _ := idpsCustomSignatureFindByID(connector, versionID, sigID); stillFound != nil {
		previewCount := idpsCustomSignatureCountPreview(connector, versionID)
		if previewCount > 0 {
			// Unpublished path: CANCEL removes all preview so the signature is gone.
			log.Printf("[INFO] Signature still in preview after PUBLISH; using CANCEL to remove unpublished (version %s)", versionID)
			versionObj2, _ := versionsClient.Get(versionID)
			cancelPayload := model.CustomSignatureValidationPayload{}
			if versionObj2.Revision != nil {
				cancelPayload.Revision = versionObj2.Revision
			}
			_ = customSigsClient.Create(versionID, cancelPayload, idpsCustomSignatureActionCancel)
		} else {
			// Retry VALIDATE with full UI path first, then re-fetch and PUBLISH.
			fullPath := fmt.Sprintf("%s/%s/signatures-preview/%s", idpsCustomSignatureVersionsPathPrefix, versionID, sigID)
			pathsToTry := []string{fullPath, deletedRef, "/signatures-preview/" + sigID, "signatures-preview/" + sigID}
			for _, p := range pathsToTry {
				v, getErr := versionsClient.Get(versionID)
				if getErr != nil {
					continue
				}
				valPayload := model.CustomSignatureValidationPayload{DeletedSignatures: []string{p}}
				if v.Revision != nil {
					valPayload.Revision = v.Revision
				}
				if customSigsClient.Create(versionID, valPayload, idpsCustomSignatureActionValidate) != nil {
					continue
				}
				v2, _ := versionsClient.Get(versionID)
				pubPayload := model.CustomSignatureValidationPayload{}
				if v2.Revision != nil {
					pubPayload.Revision = v2.Revision
				}
				if customSigsClient.Create(versionID, pubPayload, idpsCustomSignatureActionPublish) == nil {
					break
				}
			}
		}
	}

	d.SetId("")
	return nil
}

// idpsCustomSignatureCountPreview returns the number of custom signatures in preview (unpublished) for the version.
func idpsCustomSignatureCountPreview(connector client.Connector, versionID string) int {
	client := custom_signature_versions.NewCustomSignaturesClient(connector)
	include := custom_signature_versions.CustomSignatures_LIST_INCLUDE_PREVIEW_CUSTOM_SIGNATURES
	pageSize := int64(1000)
	listResult, err := client.List(versionID, nil, &include, nil, nil, &pageSize, nil, nil)
	if err != nil {
		return 0
	}
	if listResult.Results == nil {
		return 0
	}
	return len(listResult.Results)
}

// ResourceNsxtPolicyIdpsCustomSignatureExistsByContent returns true if any custom signature in the version
// has OriginalSignature containing contentSubstring. Used by acceptance tests to verify delete (no signature with that content).
func ResourceNsxtPolicyIdpsCustomSignatureExistsByContent(connector client.Connector, versionID, contentSubstring string) (bool, error) {
	customSigsClient := custom_signature_versions.NewCustomSignaturesClient(connector)
	pageSize := int64(1000)
	for _, include := range []*string{
		ptr(custom_signature_versions.CustomSignatures_LIST_INCLUDE_PREVIEW_CUSTOM_SIGNATURES),
		ptr(custom_signature_versions.CustomSignatures_LIST_INCLUDE_CUSTOM_SIGNATURES),
		nil,
	} {
		listResult, err := customSigsClient.List(versionID, nil, include, nil, nil, &pageSize, nil, nil)
		if err != nil {
			continue
		}
		for i := range listResult.Results {
			sig := &listResult.Results[i]
			if sig.OriginalSignature != nil && strings.Contains(*sig.OriginalSignature, contentSubstring) {
				return true, nil
			}
		}
	}
	return false, nil
}
