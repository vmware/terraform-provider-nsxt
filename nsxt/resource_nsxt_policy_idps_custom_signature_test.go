// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	// Custom signature version ID used in tests (e.g. "default" on NSX).
	testAccIdpsCustomSignatureVersionID = "default"
	// Base SID for test Snort rules; each test uses a different offset to avoid conflicts.
	testAccIdpsCustomSignatureSIDBase = 5999991
	// CheckDestroy: retry the "is it gone?" check (condition-based) for tests that use it (basic, update, signature_diff).
	testAccIdpsCustomSignatureDestroyRetries = 16
	testAccIdpsCustomSignatureDestroyDelay   = 2 * time.Second
)

func TestAccResourceNsxtPolicyIdpsCustomSignature_basic(t *testing.T) {
	testResourceName := "nsxt_policy_idps_custom_signature.test"
	// Unique content per test so we never match a leftover signature from another run.
	signature := fmt.Sprintf(`alert tcp any any -> any any (msg:"TF Acc Test Basic"; content:"tf_acc_basic_%d"; nocase; metadata:signature_severity Medium; sid:%d; rev:1;)`, testAccIdpsCustomSignatureSIDBase, testAccIdpsCustomSignatureSIDBase)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIdpsCustomSignatureCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsCustomSignatureBasic(signature),
				Check: resource.ComposeTestCheckFunc(
					// Only verify outcomes that are 100% guaranteed when create/read succeed.
					testAccNsxtPolicyIdpsCustomSignatureExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, "signature_version_id", testAccIdpsCustomSignatureVersionID),
					resource.TestCheckResourceAttr(testResourceName, "publish", "false"),
				),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable", "publish"},
			},
		},
	})
}

func TestAccResourceNsxtPolicyIdpsCustomSignature_updateSignature(t *testing.T) {
	testResourceName := "nsxt_policy_idps_custom_signature.test"
	sid := testAccIdpsCustomSignatureSIDBase + 1
	sig1 := fmt.Sprintf(`alert tcp any any -> any any (msg:"TF Acc Test Update 1"; content:"tf_acc_update1_%d"; nocase; metadata:signature_severity Low; sid:%d; rev:1;)`, sid, sid)
	sig2 := fmt.Sprintf(`alert tcp any any -> any any (msg:"TF Acc Test Update 2"; content:"tf_acc_update2_%d"; nocase; metadata:signature_severity High; sid:%d; rev:2;)`, sid, sid)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIdpsCustomSignatureCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsCustomSignatureBasic(sig1),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsCustomSignatureExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, "signature_version_id", testAccIdpsCustomSignatureVersionID),
				),
			},
			{
				Config: testAccNsxtPolicyIdpsCustomSignatureBasic(sig2),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsCustomSignatureExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, "signature_version_id", testAccIdpsCustomSignatureVersionID),
				),
			},
		},
	})
}

// TestAccResourceNsxtPolicyIdpsCustomSignature_withPublish mirrors the UI delete flow:
// create a published custom signature, then remove it (terraform destroy = VALIDATE then PUBLISH), then verify it is gone.
func TestAccResourceNsxtPolicyIdpsCustomSignature_withPublish(t *testing.T) {
	testResourceName := "nsxt_policy_idps_custom_signature.test"
	sid := testAccIdpsCustomSignatureSIDBase + 2
	contentSubstring := fmt.Sprintf("tf_acc_publish_%d", sid)
	signature := fmt.Sprintf(`alert tcp any any -> any any (msg:"TF Acc Test Publish"; content:"%s"; nocase; metadata:signature_severity Medium; sid:%d; rev:1;)`, contentSubstring, sid)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIdpsCustomSignatureCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			// Step 1: Create custom signature and publish (like UI: add, validate, publish).
			{
				Config: testAccNsxtPolicyIdpsCustomSignatureWithPublish(signature, true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsCustomSignatureExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, "signature_version_id", testAccIdpsCustomSignatureVersionID),
					resource.TestCheckResourceAttr(testResourceName, "publish", "true"),
				),
			},
			// Step 2: Remove the resource (no custom signature in config). Terraform runs destroy = Delete (VALIDATE then PUBLISH, like UI).
			// Then verify the signature is gone from the backend (by content).
			{
				Config:             testAccNsxtPolicyIdpsCustomSignatureNoResource(),
				Destroy:            false,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsCustomSignatureGoneByContent(testAccIdpsCustomSignatureVersionID, contentSubstring),
				),
			},
		},
	})
}

func testAccNsxtPolicyIdpsCustomSignatureExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource %s has no ID", resourceName)
		}
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		exists, err := resourceNsxtPolicyIdpsCustomSignatureExists(connector, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("custom signature %s does not exist", rs.Primary.ID)
		}
		return nil
	}
}

// testAccNsxtPolicyIdpsCustomSignatureCheckDestroy verifies the signature is gone from the backend.
// We retry the check (condition-based) to allow for eventual consistency; no fixed sleep before a single check.
func testAccNsxtPolicyIdpsCustomSignatureCheckDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_idps_custom_signature" {
			continue
		}
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		var lastErr error
		for attempt := 0; attempt < testAccIdpsCustomSignatureDestroyRetries; attempt++ {
			exists, err := resourceNsxtPolicyIdpsCustomSignatureExists(connector, rs.Primary.ID)
			if err != nil {
				return err
			}
			if !exists {
				return nil
			}
			lastErr = fmt.Errorf("custom signature %s still exists after destroy (attempt %d/%d)", rs.Primary.ID, attempt+1, testAccIdpsCustomSignatureDestroyRetries)
			if attempt < testAccIdpsCustomSignatureDestroyRetries-1 {
				time.Sleep(testAccIdpsCustomSignatureDestroyDelay)
			}
		}
		return lastErr
	}
	return nil
}

func testAccNsxtPolicyIdpsCustomSignatureBasic(signature string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_idps_custom_signature" "test" {
  signature_version_id = %q
  signature            = %q
  enable               = true
  publish              = false
}`, testAccIdpsCustomSignatureVersionID, signature)
}

func testAccNsxtPolicyIdpsCustomSignatureWithPublish(signature string, publish bool) string {
	return fmt.Sprintf(`
resource "nsxt_policy_idps_custom_signature" "test" {
  signature_version_id = %q
  signature            = %q
  enable               = true
  publish              = %v
}`, testAccIdpsCustomSignatureVersionID, signature, publish)
}

// testAccNsxtPolicyIdpsCustomSignatureNoResource returns a minimal valid config with no custom signature resource,
// so Terraform will destroy the resource from the previous step (triggers Delete: VALIDATE then PUBLISH).
func testAccNsxtPolicyIdpsCustomSignatureNoResource() string {
	return `# no resources - step 2 destroys custom signature from step 1`
}

// testAccNsxtPolicyIdpsCustomSignatureGoneByContent verifies no custom signature in the version contains contentSubstring.
// Used after destroy to verify delete (VALIDATE then PUBLISH) removed the signature. Retries for eventual consistency.
func testAccNsxtPolicyIdpsCustomSignatureGoneByContent(versionID, contentSubstring string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		var lastErr error
		for attempt := 0; attempt < testAccIdpsCustomSignatureDestroyRetries; attempt++ {
			exists, err := ResourceNsxtPolicyIdpsCustomSignatureExistsByContent(connector, versionID, contentSubstring)
			if err != nil {
				return err
			}
			if !exists {
				return nil
			}
			lastErr = fmt.Errorf("custom signature with content %q still exists after destroy (attempt %d/%d)", contentSubstring, attempt+1, testAccIdpsCustomSignatureDestroyRetries)
			if attempt < testAccIdpsCustomSignatureDestroyRetries-1 {
				time.Sleep(testAccIdpsCustomSignatureDestroyDelay)
			}
		}
		return lastErr
	}
}
