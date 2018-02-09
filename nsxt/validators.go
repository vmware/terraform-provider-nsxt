/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net"
	"strconv"
	"strings"
)

// Validations for Port objects

func isSinglePort(v string) bool {
	i, err := strconv.ParseUint(v, 10, 32)
	if err != nil {
		return false
	}
	if i > 65536 {
		return false
	}
	return true
}

func isPortRange(v string) bool {
	s := strings.Split(v, "-")
	if len(s) != 2 {
		return false
	}
	if !isSinglePort(s[0]) || !isSinglePort(s[1]) {
		return false
	}
	return true
}

func validatePortRange() schema.SchemaValidateFunc {
	// A single port num or a range of ports
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if !isPortRange(value) && !isSinglePort(value) {
			errors = append(errors, fmt.Errorf(
				"expected %q to be a port range or a single port. Got %s", k, value))
		}
		return
	}
}

func validateSinglePort() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if !isSinglePort(value) {
			errors = append(errors, fmt.Errorf(
				"expected %q to be a single port number. Got %s", k, value))
		}
		return
	}
}

// Validations for IP objects
func IsIpRange(v string) bool {
	s := strings.Split(v, "-")
	if len(s) != 2 {
		return false
	}
	ip1 := net.ParseIP(s[0])
	ip2 := net.ParseIP(s[1])
	if ip1 == nil || ip2 == nil {
		return false
	}
	return true
}

func IsSingleIp(v string) bool {
	ip := net.ParseIP(v)
	if ip == nil {
		return false
	}
	return true
}

func IsCidr(v string) bool {
	_, ipnet, err := net.ParseCIDR(v)
	if err != nil {
		return false
	}
	if ipnet == nil || v != ipnet.String() {
		return false
	}
	return true
}

func ValidateCidrOrIPOrRange() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if !IsCidr(v) && !IsSingleIp(v) && !IsIpRange(v) {
			es = append(es, fmt.Errorf(
				"expected %s to contain a valid CIDR or IP or Range, got: %s", k, v))
		}
		return
	}
}

func ValidateSingleIP() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if !IsSingleIp(v) {
			es = append(es, fmt.Errorf(
				"expected %s to contain a valid IP, got: %s", k, v))
		}
		return
	}
}
