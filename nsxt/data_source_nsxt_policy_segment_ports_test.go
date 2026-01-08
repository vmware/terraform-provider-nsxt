package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceNsxtPolicySegmentPorts_bySegmentPath(t *testing.T) {
	testAccDataSourceNsxtPolicySegmentPorts_bySegmentPath(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccDataSourceNsxtPolicySegmentPorts_bySegmentPathMultitenancy(t *testing.T) {
	testAccDataSourceNsxtPolicySegmentPorts_bySegmentPath(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccDataSourceNsxtPolicySegmentPorts_bySegmentPath(t *testing.T, withContext bool, preCheck func()) {
	segmentName := getAccTestResourceName()
	port1Name := getAccTestResourceName()
	port2Name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_segment_ports.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySegmentCheckDestroy(state, segmentName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicySegmentPortsReadBySegmentPathTemplate(segmentName, port1Name, port2Name, withContext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "items.#"),
					resource.TestCheckResourceAttr(testResourceName, "items.#", "2"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.id"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.path"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.1.id"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.1.path"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.1.segment_path"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicySegmentPorts_bySegmentPathAndDisplayName(t *testing.T) {
	testAccDataSourceNsxtPolicySegmentPorts_bySegmentPathAndDisplayName(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccDataSourceNsxtPolicySegmentPorts_bySegmentPathAndDisplayNameMultitenancy(t *testing.T) {
	testAccDataSourceNsxtPolicySegmentPorts_bySegmentPathAndDisplayName(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccDataSourceNsxtPolicySegmentPorts_bySegmentPathAndDisplayName(t *testing.T, withContext bool, preCheck func()) {
	segmentName := getAccTestResourceName()
	port1Name := getAccTestResourceName()
	port2Name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_segment_ports.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySegmentCheckDestroy(state, segmentName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicySegmentPortsReadBySegmentPathAndDisplayNameTemplate(segmentName, port1Name, port2Name, withContext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "items.#"),
					resource.TestCheckResourceAttr(testResourceName, "items.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "items.0.display_name", port1Name),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.id"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.path"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.segment_path"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicySegmentPorts_byDisplayName(t *testing.T) {
	testAccDataSourceNsxtPolicySegmentPorts_byDisplayName(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccDataSourceNsxtPolicySegmentPorts_byDisplayNameMultitenancy(t *testing.T) {
	testAccDataSourceNsxtPolicySegmentPorts_byDisplayName(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccDataSourceNsxtPolicySegmentPorts_byDisplayName(t *testing.T, withContext bool, preCheck func()) {
	segmentName := getAccTestResourceName()
	portName := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_segment_ports.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySegmentCheckDestroy(state, segmentName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicySegmentPortsReadByDisplayNameTemplate(segmentName, portName, withContext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "items.#"),
					resource.TestCheckResourceAttr(testResourceName, "items.0.display_name", portName),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.id"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.path"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.segment_path"),
				),
			},
		},
	})
}

func testAccNsxtPolicySegmentPortsReadBySegmentPathTemplate(segmentName, port1Name, port2Name string, withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}

	tzName := getOverlayTransportZoneName()
	if withContext {
		return testAccNsxtPolicySegmentNoTransportZoneTemplate(segmentName, "12.12.2.1/24", withContext) + fmt.Sprintf(`

resource "nsxt_policy_segment_port" "port1" {
  %s
  display_name = "%s"
  description  = "Acceptance Test Port 1"
  segment_path = nsxt_policy_segment.test.path
}

resource "nsxt_policy_segment_port" "port2" {
  %s
  display_name = "%s"
  description  = "Acceptance Test Port 2"
  segment_path = nsxt_policy_segment.test.path
}

data "nsxt_policy_segment_ports" "test" {
  %s
  segment_path = nsxt_policy_segment.test.path
  depends_on   = [nsxt_policy_segment_port.port1, nsxt_policy_segment_port.port2]
}
`, context, port1Name, context, port2Name, context)
	}

	return testAccNsxtPolicySegmentImportTemplate(tzName, segmentName, withContext) + fmt.Sprintf(`

resource "nsxt_policy_segment_port" "port1" {
  %s
  display_name = "%s"
  description  = "Acceptance Test Port 1"
  segment_path = nsxt_policy_segment.test.path
}

resource "nsxt_policy_segment_port" "port2" {
  %s
  display_name = "%s"
  description  = "Acceptance Test Port 2"
  segment_path = nsxt_policy_segment.test.path
}

data "nsxt_policy_segment_ports" "test" {
  %s
  segment_path = nsxt_policy_segment.test.path
  depends_on   = [nsxt_policy_segment_port.port1, nsxt_policy_segment_port.port2]
}
`, context, port1Name, context, port2Name, context)
}

func testAccNsxtPolicySegmentPortsReadBySegmentPathAndDisplayNameTemplate(segmentName, port1Name, port2Name string, withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}

	tzName := getOverlayTransportZoneName()
	if withContext {
		return testAccNsxtPolicySegmentNoTransportZoneTemplate(segmentName, "12.12.2.1/24", withContext) + fmt.Sprintf(`

resource "nsxt_policy_segment_port" "port1" {
  %s
  display_name = "%s"
  description  = "Acceptance Test Port 1"
  segment_path = nsxt_policy_segment.test.path
}

resource "nsxt_policy_segment_port" "port2" {
  %s
  display_name = "%s"
  description  = "Acceptance Test Port 2"
  segment_path = nsxt_policy_segment.test.path
}

data "nsxt_policy_segment_ports" "test" {
  %s
  segment_path = nsxt_policy_segment.test.path
  display_name = "%s"
  depends_on   = [nsxt_policy_segment_port.port1, nsxt_policy_segment_port.port2]
}
`, context, port1Name, context, port2Name, context, port1Name)
	}

	return testAccNsxtPolicySegmentImportTemplate(tzName, segmentName, withContext) + fmt.Sprintf(`

resource "nsxt_policy_segment_port" "port1" {
  %s
  display_name = "%s"
  description  = "Acceptance Test Port 1"
  segment_path = nsxt_policy_segment.test.path
}

resource "nsxt_policy_segment_port" "port2" {
  %s
  display_name = "%s"
  description  = "Acceptance Test Port 2"
  segment_path = nsxt_policy_segment.test.path
}

data "nsxt_policy_segment_ports" "test" {
  %s
  segment_path = nsxt_policy_segment.test.path
  display_name = "%s"
  depends_on   = [nsxt_policy_segment_port.port1, nsxt_policy_segment_port.port2]
}
`, context, port1Name, context, port2Name, context, port1Name)
}

func testAccNsxtPolicySegmentPortsReadByDisplayNameTemplate(segmentName, portName string, withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}

	tzName := getOverlayTransportZoneName()
	if withContext {
		return testAccNsxtPolicySegmentNoTransportZoneTemplate(segmentName, "12.12.2.1/24", withContext) + fmt.Sprintf(`

resource "nsxt_policy_segment_port" "test" {
  %s
  display_name = "%s"
  description  = "Acceptance Test Port"
  segment_path = nsxt_policy_segment.test.path
}

data "nsxt_policy_segment_ports" "test" {
  %s
  display_name = "%s"
  depends_on   = [nsxt_policy_segment_port.test]
}
`, context, portName, context, portName)
	}

	return testAccNsxtPolicySegmentImportTemplate(tzName, segmentName, withContext) + fmt.Sprintf(`

resource "nsxt_policy_segment_port" "test" {
  %s
  display_name = "%s"
  description  = "Acceptance Test Port"
  segment_path = nsxt_policy_segment.test.path
}

data "nsxt_policy_segment_ports" "test" {
  %s
  display_name = "%s"
  depends_on   = [nsxt_policy_segment_port.test]
}
`, context, portName, context, portName)
}
