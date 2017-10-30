package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt"
	"net/http"
	"regexp"
	"testing"
)

func TestNSXLogicalSwitchBasic(t *testing.T) {

	switchName := fmt.Sprintf("test-nsx-logical-switch-overlay")
	updateSwitchName := fmt.Sprintf("%s-update", switchName)
	testResourceName := "nsxt_logical_switch.test"
	novlan := "0"
	replicationMode := "MTEP"
	fmt.Printf("Test logical switch display_name is %s\n", switchName)
	// overlay TZ id
	transportZoneId := "dbb54c3c-41ff-4a98-8dbf-1c4b84a1a94d"
	//transportZoneId := getTzIdByType("Overlay")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLogicalSwitchCheckDestroy(state, switchName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccNSXLogicalSwitchNoTZIDTemplate(switchName),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config: testAccNSXLogicalSwitchCreateTemplate(switchName, transportZoneId, novlan, replicationMode),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(switchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", switchName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "transport_zone_id", transportZoneId),
					resource.TestCheckResourceAttr(testResourceName, "admin_state", "UP"),
					resource.TestCheckResourceAttr(testResourceName, "replication_mode", replicationMode),
					resource.TestCheckResourceAttr(testResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "vlan", novlan),
				),
			},
			{
				Config: testAccNSXLogicalSwitchUpdateTemplate(updateSwitchName, transportZoneId, novlan, replicationMode),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(updateSwitchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateSwitchName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "transport_zone_id", transportZoneId),
					resource.TestCheckResourceAttr(testResourceName, "admin_state", "UP"),
					resource.TestCheckResourceAttr(testResourceName, "replication_mode", replicationMode),
					resource.TestCheckResourceAttr(testResourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "vlan", novlan),
				),
			},
		},
	})
}

func TestNSXLogicalSwitchvlan(t *testing.T) {

	switchName := fmt.Sprintf("test-nsx-logical-switch-vlan")
	updateSwitchName := fmt.Sprintf("%s-update", switchName)
	// vlan TZ Id
	transportZoneId := "3507a11d-4ce3-44c8-a495-2a17b43855f4"
	//getTzIdByType("VLAN")
	testResourceName := "nsxt_logical_switch.test"

	origvlan := "1"
	updatedvlan := "2"
	replicationMode := ""
	fmt.Printf("Test logical switch display_name is %s\n", switchName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLogicalSwitchCheckDestroy(state, switchName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccNSXLogicalSwitchNovlanTemplate(switchName, transportZoneId),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config: testAccNSXLogicalSwitchCreateTemplate(switchName, transportZoneId, origvlan, replicationMode),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(switchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", switchName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "transport_zone_id", transportZoneId),
					resource.TestCheckResourceAttr(testResourceName, "admin_state", "UP"),
					resource.TestCheckResourceAttr(testResourceName, "vlan", origvlan),
				),
			},
			{
				Config: testAccNSXLogicalSwitchUpdateTemplate(updateSwitchName, transportZoneId, updatedvlan, replicationMode),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(updateSwitchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateSwitchName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "transport_zone_id", transportZoneId),
					resource.TestCheckResourceAttr(testResourceName, "admin_state", "UP"),
					resource.TestCheckResourceAttr(testResourceName, "vlan", updatedvlan),
				),
			},
		},
	})

}

func getTzIdByType(TransportZoneType string) string {
	//fmt.Printf("DEBUG ADIT getTzIdByType %s\n", TransportZoneType)
	//fmt.Printf("DEBUG ADIT getTzIdByType meta %v\n", testAccProvider.Meta())
	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

	obj_list, _, err := nsxClient.NetworkTransportApi.ListTransportZones(nsxClient.Context, nil)
	if err != nil {
		fmt.Printf("Error while reading transport zones: %v\n", err)
		return ""
	}
	// go over the list and return the first of this type
	for _, obj_in_list := range obj_list.Results {
		if obj_in_list.TransportType == TransportZoneType {
			return obj_in_list.Id
		}
	}
	return ""
}

func testAccNSXLogicalSwitchExists(display_name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX logical switch resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX logical switch resource ID not set in resources ")
		}

		logicalSwitch, responseCode, err := nsxClient.LogicalSwitchingApi.GetLogicalSwitch(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving logical switch ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if logical switch %s exists. HTTP return code was %d", resourceID, responseCode)
		}

		if display_name == logicalSwitch.DisplayName {
			return nil
		}
		return fmt.Errorf("NSX logical switch %s wasn't found", display_name)
	}
}

func testAccNSXLogicalSwitchCheckDestroy(state *terraform.State, display_name string) error {

	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_logical_switch" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		logicalSwitch, responseCode, err := nsxClient.LogicalSwitchingApi.GetLogicalSwitch(nsxClient.Context, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving logical switch ID %s. Error: %v", resourceID, err)
		}

		if display_name == logicalSwitch.DisplayName {
			return fmt.Errorf("NSX logical switch %s still exists", display_name)
		}
	}
	return nil
}

func testAccNSXLogicalSwitchNoTZIDTemplate(switchName string) string {
	return fmt.Sprintf(`
resource "nsxt_logical_switch" "test" {
display_name = "%s"
admin_state = "UP"
description = "Acceptance Test"
replication_mode = "MTEP"
}`, switchName)
}

func testAccNSXLogicalSwitchNovlanTemplate(switchName string, transportZoneId string) string {
	return fmt.Sprintf(`
resource "nsxt_logical_switch" "test" {
display_name = "%s"
admin_state = "UP"
description = "Acceptance Test"
transport_zone_id = "%s"
}`, switchName)
}

func testAccNSXLogicalSwitchCreateTemplate(switchName string, transportZoneId string, vlan string, replicationMode string) string {
	return fmt.Sprintf(`
resource "nsxt_logical_switch" "test" {
display_name = "%s"
admin_state = "UP"
description = "Acceptance Test"
transport_zone_id = "%s"
replication_mode = "%s"
vlan = "%s"
tags = [
    {
	    Scope = "scope1"
        tag = "tag1"
    }
]
}`, switchName, transportZoneId, replicationMode, vlan)
}

func testAccNSXLogicalSwitchUpdateTemplate(switchUpdateName, transportZoneId string, vlan string, replicationMode string) string {
	return fmt.Sprintf(`
resource "nsxt_logical_switch" "test" {
display_name = "%s"
admin_state = "UP"
description = "Acceptance Test Update"
transport_zone_id = "%s"
replication_mode = "%s"
vlan = "%s"
tags = [
	{
		Scope = "scope1"
    	tag = "tag1"
    },
	{
		Scope = "scope2"
    	tag = "tag2"
    },
]
}`, switchUpdateName, transportZoneId, replicationMode, vlan)
}
