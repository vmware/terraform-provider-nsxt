/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security


import (
	cjson "github.com/gibson042/canonicaljson-go"
)

// Used for transforming JSON messages into their
// canonical representation.
//
// The canonical form is defined by the following rules:
//     1. Non-significant(1) whitespace characters MUST NOT be used
//     2. Non-significant(1) line endings MUST NOT be used
//     3. Entries (set of name/value pairs) in JSON objects MUST be sorted
//        lexicographically(2) by their names based on UCS codepoint values
//     4. Arrays MUST preserve their initial ordering
//
// Link to the IEFT proposal:
// https://datatracker.ietf.org/doc/draft-staykov-hu-json-canonical-form/

type JSONCanonicalEncoder struct {
}

func NewJSONCanonicalEncoder() *JSONCanonicalEncoder {
	return &JSONCanonicalEncoder{}
}

func (j *JSONCanonicalEncoder) Marshal(v interface{}) ([]byte, error) {
	return cjson.Marshal(v)
}
