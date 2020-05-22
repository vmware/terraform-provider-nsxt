/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
func isIPRange(v string) bool {
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

func isSingleIP(v string) bool {
	ip := net.ParseIP(v)
	return ip != nil
}

func isCidr(v string, allowMaxPrefix bool, isIP bool) bool {
	_, ipnet, err := net.ParseCIDR(v)
	if err != nil {
		return false
	}
	if ipnet == nil {
		return false
	}
	if isIP && v == ipnet.String() {
		return false
	}
	if !isIP && v != ipnet.String() {
		return false
	}

	_, bits := ipnet.Mask.Size()
	if !allowMaxPrefix && bits == 1 {
		return false
	}

	return true
}

func validatePortAddress() schema.SchemaValidateFunc {
	// Expects ip_address/prefix (prefix < 32)
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if !isCidr(v, false, true) {
			es = append(es, fmt.Errorf(
				"expected %s to contain a valid port address/prefix, got: %s", k, v))
		}
		return
	}
}

func validateCidrOrIPOrRange() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if !isCidr(v, true, false) && !isSingleIP(v) && !isIPRange(v) {
			es = append(es, fmt.Errorf(
				"expected %s to contain a valid CIDR or IP or Range, got: %s", k, v))
		}
		return
	}
}

func validateIPOrRange() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if !isSingleIP(v) && !isIPRange(v) {
			es = append(es, fmt.Errorf(
				"expected %s to contain a valid IP or Range, got: %s", k, v))
		}
		return
	}
}

func validateSingleIP() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if !isSingleIP(v) {
			es = append(es, fmt.Errorf(
				"expected %s to contain a valid IP, got: %s", k, v))
		}
		return
	}
}

func isValidStringUint(value string, bits int) bool {
	_, err := strconv.ParseUint(value, 10, bits)
	return (err == nil)
}

func isValidASN(value string) bool {
	tokens := strings.Split(value, ":")
	if len(tokens) != 2 {
		return false
	}
	// ASN is limited to 2 bytes
	if !isValidStringUint(tokens[0], 16) {
		return false
	}

	// Number is limited to 4 bytes
	if !isValidStringUint(tokens[1], 32) {
		return false
	}

	return true
}

func validateASNPair(i interface{}, k string) (s []string, es []error) {
	v, ok := i.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if !isValidASN(v) {
		err := fmt.Errorf("expected %s to contain ASN:<number> with ASN limited to 2 bytes, got: %s", k, v)
		es = append(es, err)
	}
	return
}

func validateIPorASNPair(i interface{}, k string) (s []string, es []error) {
	v, ok := i.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	err := fmt.Errorf("expected %s to contain <ASN>:<number> or <IP>:<number> with ASN limited to 2 bytes, got: %s", k, v)

	tokens := strings.Split(v, ":")
	if len(tokens) != 2 {
		es = append(es, err)
		return
	}

	if isSingleIP(tokens[0]) {
		if !isValidStringUint(tokens[1], 32) {
			es = append(es, err)
		}
		return
	}

	if !isValidASN(v) {
		es = append(es, err)
	}

	return
}

func validate4ByteASNPlain(i interface{}, k string) (s []string, es []error) {
	v, ok := i.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	if !isValidStringUint(v, 32) {
		es = append(es, fmt.Errorf("invalid ASN number %s", v))
		return
	}
	return
}

func validateIPRange() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if !isIPRange(v) {
			es = append(es, fmt.Errorf(
				"expected %s to contain a valid IP range, got: %s", k, v))
		}
		return
	}
}

func validateCidr() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if !isCidr(v, true, false) {
			es = append(es, fmt.Errorf(
				"expected %s to contain a valid CIDR, got: %s", k, v))
		}
		return
	}
}

func validateIPCidr() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if !isCidr(v, true, true) {
			es = append(es, fmt.Errorf(
				"expected %s to contain a valid CIDR, got: %s", k, v))
		}
		return
	}
}

func isPowerOfTwo(num int) bool {
	for num >= 2 {
		if num%2 != 0 {
			return false
		}
		num = num / 2
	}
	return num == 1
}

func validatePowerOf2(allowZero bool, maxVal int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return
		}

		if allowZero && v == 0 {
			return
		}

		if maxVal != 0 && v > maxVal {
			es = append(es, fmt.Errorf(
				"expected %s to be <= %d, got: %d", k, maxVal, v))
			return
		}

		if !isPowerOfTwo(v) {
			es = append(es, fmt.Errorf(
				"expected %s to be a power of 2, got: %d", k, v))
		}
		return
	}
}

var supportedSSLProtocols = []string{"SSL_V2", "SSL_V3", "TLS_V1", "TLS_V1_1", "TLS_V1_2"}

func validateSSLProtocols() schema.SchemaValidateFunc {
	return validation.StringInSlice(supportedSSLProtocols, false)
}

var supportedSSLCiphers = []string{
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
	"TLS_ECDH_ECDSA_WITH_AES_256_CBC_SHA",
	"TLS_ECDH_RSA_WITH_AES_256_CBC_SHA",
	"TLS_RSA_WITH_AES_256_CBC_SHA",
	"TLS_RSA_WITH_AES_128_CBC_SHA",
	"TLS_RSA_WITH_3DES_EDE_CBC_SHA",
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384",
	"TLS_RSA_WITH_AES_128_CBC_SHA256",
	"TLS_RSA_WITH_AES_128_GCM_SHA256",
	"TLS_RSA_WITH_AES_256_CBC_SHA256",
	"TLS_RSA_WITH_AES_256_GCM_SHA384",
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384",
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
	"TLS_ECDH_ECDSA_WITH_AES_128_CBC_SHA",
	"TLS_ECDH_ECDSA_WITH_AES_128_CBC_SHA256",
	"TLS_ECDH_ECDSA_WITH_AES_128_GCM_SHA256",
	"TLS_ECDH_ECDSA_WITH_AES_256_CBC_SHA384",
	"TLS_ECDH_ECDSA_WITH_AES_256_GCM_SHA384",
	"TLS_ECDH_RSA_WITH_AES_128_CBC_SHA",
	"TLS_ECDH_RSA_WITH_AES_128_CBC_SHA256",
	"TLS_ECDH_RSA_WITH_AES_128_GCM_SHA256",
	"TLS_ECDH_RSA_WITH_AES_256_CBC_SHA384",
	"TLS_ECDH_RSA_WITH_AES_256_GCM_SHA384",
}

func validateSSLCiphers() schema.SchemaValidateFunc {
	return validation.StringInSlice(supportedSSLCiphers, false)
}

func validatePolicyPath() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if !isPolicyPath(v) {
			es = append(es, fmt.Errorf("Invalid policy path: %s", v))
		}

		return
	}
}

func validateVLANId(i interface{}, k string) (s []string, es []error) {
	var vlan int
	vlan, ok := i.(int)
	if !ok {
		v, ok := i.(string)
		if ok {
			vlanVal, err := strconv.Atoi(v)
			if err != nil {
				es = append(es, fmt.Errorf("cannot convert %s to int", v))
				return
			}
			vlan = vlanVal
		} else {
			es = append(es, fmt.Errorf("Expected VLAN type to be int or string"))
			return
		}
	}
	if vlan < 0 || vlan > 4095 {
		es = append(es, fmt.Errorf("invalid VLAN ID %d", vlan))
		return
	}
	return
}

func validateStringIntBetween(start int, end int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		intVal, err := strconv.Atoi(v)
		if err != nil {
			es = append(es, fmt.Errorf("cannot convert %s to int", v))
			return
		}
		if intVal < start || intVal > end {
			es = append(es, fmt.Errorf("the value %s must be between %v and %v", v, start, end))
		}
		return
	}
}

func validateNsxtProviderHostFormat() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if strings.HasPrefix(v, "https://") || strings.HasPrefix(v, "http://") {
			es = append(es, fmt.Errorf("not expecting http:// or https:// in the host, but got %s", v))
		}

		return
	}
}
