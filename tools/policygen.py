#!/usr/bin/python
# -*- coding: latin-1 -*-

#  Copyright (c) 2017 VMware, Inc. All Rights Reserved.
#  SPDX-License-Identifier: MPL-2.0

# Generating Data Sources:
# 1. Run the following command: python policygen.py <resources struct name> data
# 2. The generated files end up in your current directory; move them to the proper directory.
# 3. Add the data source to provider.go in the DataSourcesMap.
# 4. If needed, format the generated go files with: go fmt
# 5. Make sure go is happy with the format by running golint and go vet on the source files.
# 6. Run the generated test(s) to ensure they pass.
# 7. Add the doc file to nsxt.erb

import sys
import re
import os

PACKAGE_NAME = "nsxt"

# TODO: clone when opensourced
# TODO: support non-infra resources/data sources
try:
    GO_PATH = os.environ['GOPATH']
except Exception:
    print("Please set GOPATH in your environment")
    exit(2)

BASE_SDK_PATH = "%s/src/github.com/vmware/go-nsx-t-policy" % GO_PATH
STRUCTS_FILE = "%s/bindings/nsx_policy/model/ModelPackageTypes.go" % BASE_SDK_PATH
TEMPLATE_RESOURCE_FILE = "resource_nsxt_policy_template"
TEMPLATE_RESOURCE_TEST_FILE = "resource_nsxt_policy_test_template"
TEMPLATE_RESOURCE_DOC_FILE = "resource_nsxt_policy_doc_template"
TEMPLATE_DATA_SOURCE_FILE = "data_source_nsxt_policy_template"
TEMPLATE_DATA_SOURCE_TEST_FILE = "data_source_nsxt_policy_test_template"
TEMPLATE_DATA_SOURCE_DOC_FILE = "data_source_nsxt_policy_doc_template"

DONT_SPLIT_US = ["IP", "LB", "SSL", "TCP", "UDP"]

LAST_COMMON_ATTR = "MarkedForDelete"

TYPE_MAP = {"string": "schema.TypeString",
            "int32": "schema.TypeInt",
            "int64": "schema.TypeInt",
            "bool": "schema.TypeBool"}

DEFAULT_VALUES_MAP = {"string": '"test"',
                      "int32": 2,
                      "int64": 2,
                      "bool": 'true'}

TEST_CREATE_VALUES_MAP = {"string": '"test-create"',
                          "int32": '"2"',
                          "int64": '"2"',
                          "bool": '"true"'}

TEST_UPDATE_VALUES_MAP = {"string": '"test-update"',
                          "int32": '"5"',
                          "int64": '"5"',
                          "bool": '"false"'}

TYPECAST_MAP = {"int64": "int", "int32": "int"}

metadata = {}

def to_lower_separated(name):
    tmp = re.sub(r'([A-Z])', r'_\1', name).lower()
    for abbr in DONT_SPLIT_US:
        split_abbr = re.sub(r'([A-Z])', r'_\1', abbr).lower()
        tmp = tmp.replace(split_abbr[1:], abbr.lower())
    return tmp[1:]

def lowercase_first(name):
    return name[:1].lower() + name[1:]

def name_to_upper(name):
    return name.title().replace('_', '')

def to_upper_separated(name):
    tmp = re.sub(r'([A-Z])', r'_\1', name).upper()[1:]
    for abbr in DONT_SPLIT_US:
        split_abbr = re.sub(r'([A-Z])', r'_\1', abbr)
        tmp = tmp.replace(split_abbr[1:], abbr)

    return tmp

def is_list_complex_attr(attr):
    if attr['type'].startswith('[]'):
        # this is a list.
        if attr['type'][2:] not in TYPE_MAP:
            # complex type: needs to be in a single form
            return True
    return False

def get_const_values_name(attr_name):
    return "%s%sValues" % (lowercase_first(metadata['resource']), attr_name)

def build_enums():

    result = ""
    for attr in metadata['attrs']:
        if not attr['const_needed']:
            continue

        result += "\nvar %s = []string{\n" % get_const_values_name(attr['name'])
        const_candidate = "%s_%s" % (metadata['resource'], to_upper_separated(attr['name']))
        for const in metadata['constants'].keys():
            if const.startswith(const_candidate):
                result += "model.%s,\n" % const

        result += "}\n"

    return result

def get_value_for_const(attr, num=1):
    const_candidate = "%s_%s" % (metadata['resource'], to_upper_separated(attr['name']))
    for name, value in metadata['constants'].items():
        if name.startswith(const_candidate):
            num -= 1
            if num == 0:
                return value

def build_schema_attr(attr):

    result = ""
    is_array = False
    attr['const_needed'] = False
    attr_type = attr['type']
    # Computed and required attributes are not ref. At the moment we can't
    # distinguish between them and assume required, since this will necessarily
    # error out in test.
    # TODO: parse yaml for this purpose

    optional = not attr['required']
    if attr['list'] and attr_type in TYPE_MAP:
        # Handle arrays. By default, arrays are translated to sets
        # assuming in most cases order is not significant. When order is
        # significant, these should be changed to lists
        # TODO: add both with choice comment
        attr['helper'] = "%sListFromSchemaSet" % name_to_upper(attr_type)

    if attr_type not in TYPE_MAP:
        # TODO: support sub-structs
        print("Skipping attribute %s due to mysterious type %s" % (attr['name'], attr_type))
        return result

    result += "\"%s\": {\n" % attr['schema_name']

    if is_array:
        result += "Type:        schema.TypeSet,\n"
    else:
        result += "Type:        %s,\n" % TYPE_MAP[attr_type]

    #TODO: concatenate multi-string comment, and remove enum speces
    #if attr['description']:
    #    result += "Description: \"%s\"," % attr['description']

    const_candidate = "%s_%s" % (metadata['resource'], to_upper_separated(attr['name']))
    if attr['type'] == "string":
        for const in metadata['constants'].keys():
            if const.startswith(const_candidate):
                attr['const_needed'] = True

    validation = None
    if attr['const_needed']:
        validation = "ValidateFunc: validation.StringInSlice(%s, false),\n" % get_const_values_name(attr['name'])

    if is_array:
        result += "Elem: &schema.Schema{\n"
        result += "Type:        %s,\n" % TYPE_MAP[attr_type]
        if validation:
            result += validation
        result +=  "},\n"
    elif validation:
        result += validation

    if optional:
        result += "Optional:    true,\n"
    else:
        result += "Required:    true,\n"

    result += "},\n"

    return result

def build_get_attr_from_schema(attr):
    if 'helper' in attr:
        # helper function name already computed - this is the case for arrays
        return '%s := get%s(d, "%s")\n' % (
            lowercase_first(attr['name']),
            attr['helper'],
            attr['schema_name'])

    if attr['type'] in TYPECAST_MAP:
        # type casting is needed
        return '%s := %s(d.Get("%s").(%s))\n' % (
                 lowercase_first(attr['name']),
                 attr['type'],
                 attr['schema_name'],
                 TYPECAST_MAP[attr['type']])
    return '%s := d.Get("%s").(%s)\n' % (
                 lowercase_first(attr['name']),
                 attr['schema_name'],
                 attr['type'])




def load_resource_metadata():

    with open(STRUCTS_FILE, 'r') as f:
        lines = f.readlines()

    metadata["constants"] = {}
    metadata["attrs"] = []
    stage = ""
    description = ""
    resource = metadata['resource']
    for line in lines:
        # load constants for this resource
        if line.startswith("const %s" % resource):
            tokens = line.split()
            const = tokens[1]
            value = tokens[3]
            metadata["constants"][const] = value
            continue

        if line.startswith("type %s struct" % resource):
            stage = "struct"
            continue

        if stage == "struct":
            if LAST_COMMON_ATTR in line:
                stage = "attrs"
                continue

        if stage == "attrs":
            if line.startswith("}"):
                # end of type struct
                stage = ""
                continue

            if line.strip().startswith('//'):
                description += line
                continue

            tokens = line.split()
            attr = tokens[0]
            full_type = tokens[1]
            is_ref = False
            is_list = False
            if full_type.startswith("*"):
                is_ref = True
                full_type = full_type[1:]

            if full_type.startswith("[]"):
                is_list = True
                full_type = full_type[2:]

            schema_name = to_lower_separated(attr)
            if is_list and full_type not in TYPE_MAP:
                if schema_name.endswith('es'):
                    schema_name = schema_name[:2]
                if schema_name.endswith('s'):
                    schema_name = schema_name[:1]

            metadata["attrs"].append({'name': attr,
                                      'description': description,
                                      'type': full_type,
                                      'list': is_list,
                                      'ref': is_ref,
                                      'schema_name': schema_name,
                                      'required': not is_ref and not is_list})
            description = ""


def build_schema_attrs():
    result = ""
    for attr in metadata['attrs']:
        result += build_schema_attr(attr)

    return result


def build_get_attrs_from_schema():
    result = ""
    for attr in metadata['attrs']:
        result += build_get_attr_from_schema(attr)

    return result


def build_set_attrs_in_obj():
    result = ""
    for attr in metadata['attrs']:
        result += "%s: &%s,\n" % (attr['name'], lowercase_first(attr['name']))

    return result


def build_set_obj_attrs_in_schema():
    result = ""
    for attr in metadata['attrs']:
        result += 'd.Set("%s", obj.%s)\n'% (attr['schema_name'], attr['name'])

    return result


def build_test_attrs(required_only=False):

    result = ""
    for attr in metadata['attrs']:
        if required_only and not attr['required']:
            continue
        result += "%s = " % attr['schema_name']
        if attr['type'] == 'string':
            result += '"%s"'
        else:
            result += '%s'

        result += "\n"

    return result


def build_test_attrs_map(is_create=True):

    result = ""
    values_map = TEST_CREATE_VALUES_MAP if is_create else TEST_UPDATE_VALUES_MAP
    for attr in metadata['attrs']:
        value = ""
        if attr['const_needed']:
            num = 1 if is_create else 2
            value = get_value_for_const(attr, num)
        elif attr['type'] in values_map:
            value = values_map[attr['type']]
        if value:
            result += '"%s": %s,\n' % (attr["schema_name"], value)

    return result


def build_check_test_attrs(is_create=True):

    result = ""
    for attr in metadata['attrs']:
        result += 'resource.TestCheckResourceAttr(testResourceName, "%s", ' % attr['schema_name']
        result += 'accTestPolicy%s%sAttributes["%s"]),\n' % (metadata['resource'], "Create" if is_create else "Update", attr['schema_name'])

    return result


def build_doc_attrs():
    doc_attrs = ""
    for attr in metadata['attrs']:
        optional = attr['ref'] or attr['list']
        if attr['const_needed']:
            value = get_value_for_const(attr)
        else:
            if attr['type'] in DEFAULT_VALUES_MAP:
                value = DEFAULT_VALUES_MAP[attr['type']]
            else:
                value = "FILL VALUE FOR %s" % attr['type']

        if not value:
            value = "FILL ENUM VALUE"
        doc_attrs += "%s = %s\n" % (attr['schema_name'], value)

    return doc_attrs


def build_doc_attrs_reference():
    result = ""

    for attr in metadata['attrs']:
        result += "* `%s` - (%s)\n" % (attr['schema_name'],
                                   "Required" if attr['required'] else "Optional")
    return result


def build_test_attrs_sprintf(required_only=False):
    result = ""

    for attr in metadata['attrs']:
        if required_only and not attr['required']:
            continue
        result += ', attrMap["%s"]' % attr['schema_name']
    return result


def generate_replace_targets():
    metadata['schema_attrs'] = build_schema_attrs()
    metadata['enums'] = build_enums()
    metadata['set_attrs_in_obj'] = build_set_attrs_in_obj()
    metadata['get_attrs_from_schema'] = build_get_attrs_from_schema()
    metadata['set_obj_attrs_in_schema'] = build_set_obj_attrs_in_schema()
    metadata['doc_attrs'] = build_doc_attrs()
    metadata['doc_attrs_reference'] = build_doc_attrs_reference()
    metadata['test_attrs'] = build_test_attrs()
    metadata['test_required_attrs'] = build_test_attrs(True)
    metadata['test_attrs_create'] = build_test_attrs_map()
    metadata['test_attrs_update'] = build_test_attrs_map(False)
    metadata['check_attrs_create'] = build_check_test_attrs()
    metadata['check_attrs_update'] = build_check_test_attrs(False)
    metadata['test_attrs_sprintf'] = build_test_attrs_sprintf()
    metadata['test_required_attrs_sprintf'] = build_test_attrs_sprintf(True)


def replace_templates(line):
    result = line.replace("<!RESOURCE!>", metadata['resource'])
    result = result.replace("<!RESOURCES!>", "%ss" % metadata['resource'])
    result = result.replace("<!MODULE!>", metadata['module'])
    result = result.replace("<!SCHEMA_ATTRS!>", metadata['schema_attrs'])
    result = result.replace("<!GET_ATTRS_FROM_SCHEMA!>", metadata['get_attrs_from_schema'])
    result = result.replace("<!SET_ATTRS_IN_OBJ!>", metadata['set_attrs_in_obj'])
    result = result.replace("<!SET_OBJ_ATTRS_IN_SCHEMA!>", metadata['set_obj_attrs_in_schema'])
    result = result.replace("<!ENUMS!>", metadata['enums'])
    result = result.replace("<!resource_lower!>", metadata['resource_lower'])
    result = result.replace("<!resource-lower!>", metadata['resource-lower'])
    result = result.replace("<!TEST_ATTRS!>", metadata['test_attrs'])
    result = result.replace("<!TEST_REQUIRED_ATTRS!>", metadata['test_required_attrs'])
    result = result.replace("<!TEST_ATTRS_CREATE!>", metadata['test_attrs_create'])
    result = result.replace("<!TEST_ATTRS_UPDATE!>", metadata['test_attrs_update'])
    result = result.replace("<!TEST_ATTRS_SPRINTF!>", metadata['test_attrs_sprintf'])
    result = result.replace("<!TEST_REQUIRED_ATTRS_SPRINTF!>", metadata['test_required_attrs_sprintf'])
    result = result.replace("<!CHECK_ATTRS_CREATE!>", metadata['check_attrs_create'])
    result = result.replace("<!CHECK_ATTRS_UPDATE!>", metadata['check_attrs_update'])
    result = result.replace("<!DOC_ATTRS!>", metadata['doc_attrs'])
    result = result.replace("<!DOC_ATTRS_REFERENCE!>", metadata['doc_attrs_reference'])
    return result


def main():

    if len(sys.argv) < 2:
        print("Usage: %s <resource name> [data]" % sys.argv[0])
        sys.exit()

    generate_data_source = False
    if len(sys.argv) == 3 and sys.argv[2] == 'data':
        generate_data_source = True

    resource = sys.argv[1]
    metadata['resource'] = resource
    metadata['module'] = resource[0].lower() + resource[1:] + 's'
    resource_lower = to_lower_separated(resource)
    metadata['resource_lower'] = resource_lower
    metadata['resource-lower'] = resource_lower.replace('_','-')
    print("Building resource from %s" % resource)

    print("Loading metadata..")
    load_resource_metadata()
    generate_replace_targets()

    spec = {}
    if generate_data_source:
        spec = {TEMPLATE_DATA_SOURCE_FILE: "data_source_nsxt_policy_%s.go" % resource_lower,
                TEMPLATE_DATA_SOURCE_TEST_FILE: "data_source_nsxt_policy_%s_test.go" % resource_lower,
                TEMPLATE_DATA_SOURCE_DOC_FILE: "policy_%s.html.markdown" % resource_lower}

    else:
        spec = {TEMPLATE_RESOURCE_FILE: "resource_nsxt_policy_%s.go" % resource_lower,
                TEMPLATE_RESOURCE_TEST_FILE: "resource_nsxt_policy_%s_test.go" % resource_lower,
                TEMPLATE_RESOURCE_DOC_FILE: "policy_%s.html.markdown" % resource_lower}

    for source, target in spec.items():
        print("Generating %s.." % target)
        with open(source, 'r') as f:
            lines = f.readlines()

        with open(target, 'w') as f:
            for line in lines:
                f.write(replace_templates(line))

    print("Done.")

main()
