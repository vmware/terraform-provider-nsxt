#!/usr/bin/python
# -*- coding: latin-1 -*-

#  Copyright (c) 2017 VMware, Inc. All Rights Reserved.
#  SPDX-License-Identifier: MPL-2.0


import sys
import re

PACKAGE_NAME = "nsxt"
SDK_PACKAGE_NAME = "api"
MANAGER_PACKAGE_NAME = "manager"

IGNORE_ATTRS = ["Links", "Schema", "Self", "Id", "ResourceType", "CreateTime", "CreateUser", "LastModifiedTime", "LastModifiedUser", "SystemOwned"]
COMPUTED_ATTRS = ["CreateTime", "CreateUser", "LastModifiedTime", "LastModifiedUser", "SystemOwned"]
COMPUTED_AND_OPTIONAL_ATTRS = ["DisplayName"]
FORCENEW_ATTRS = ["TransportZoneId"]
# TODO: ServiceBindings
VIP_SCHEMA_ATTRS = ["Tags", "SwitchingProfileIds", "Revision", "AddressBindings"]
VIP_GETTER_ATTRS = ["Tags", "SwitchingProfileIds", "AddressBindings"]
VIP_SETTER_ATTRS = VIP_GETTER_ATTRS

TYPE_MAP = {"string": "schema.TypeString",
            "int32": "schema.TypeInt",
            "int64": "schema.TypeInt",
            "bool": "schema.TypeBool"}

TYPECAST_MAP = {"int64": "int", "int32": "int"}

indent = 0


def convert_name(name):
    tmp = re.sub(r'([A-Z])', r'_\1', name).lower()
    return tmp[1:]


def is_list_complex_attr(attr):
    if attr['type'].startswith('[]'):
        # this is a list.
        if attr['type'][2:] not in TYPE_MAP:
            # complex type: needs to be in a single form
            return True
    return False


def get_attr_fixed_name(attr):
    fixed_name = attr['name']
    if is_list_complex_attr(attr) and fixed_name.endswith('s'):
        # remove last s
        fixed_name = fixed_name[:-1]
    fixed_name = convert_name(fixed_name)
    return fixed_name


def shift():
    global indent
    indent += 1


def unshift():
    global indent
    indent -= 1


def pretty_writeln(f, line):
    for i in range(indent):
        f.write("\t")
    f.write(line)
    f.write("\n")


def write_header(f):
    pretty_writeln(f, "/* Copyright © 2018 VMware, Inc. All Rights Reserved.")
    pretty_writeln(f, "   SPDX-License-Identifier: MPL-2.0 */\n")

    pretty_writeln(f, "package %s\n" % PACKAGE_NAME)
    pretty_writeln(f, "import (")
    shift()
    pretty_writeln(f, "\"github.com/hashicorp/terraform/helper/schema\"")
    pretty_writeln(f, "%s \"github.com/vmware/go-vmware-nsxt\"" % SDK_PACKAGE_NAME)
    pretty_writeln(f, "\"github.com/vmware/go-vmware-nsxt/%s\"" % MANAGER_PACKAGE_NAME)
    pretty_writeln(f, "\"net/http\"")
    pretty_writeln(f, "\"fmt\"")
    unshift()
    pretty_writeln(f, ")\n")


def write_attr(f, attr):
    fixed_name = get_attr_fixed_name(attr)
    if attr['name'] in VIP_SCHEMA_ATTRS:
        pretty_writeln(f, "\"%s\": get%sSchema()," % (fixed_name, attr['name']))
        return

    if attr['type'] not in TYPE_MAP:
        print("Skipping attribute %s due to mysterious type %s" % (attr['name'], attr['type']))
        return

    pretty_writeln(f, "\"%s\": &schema.Schema{" % fixed_name)
    shift()
    pretty_writeln(f, "Type:        %s," % TYPE_MAP[attr['type']])

    comment = ' '
    if attr['comment']:
        comment = attr['comment']
    if attr['name'] == 'DisplayName' and comment == 'Defaults to ID if not set':
        comment = "The display name of this resource. " + comment
    pretty_writeln(f, "Description: \"%s\"," % comment)

    if attr['optional']:
        pretty_writeln(f, "Optional:    true,")
    else:
        pretty_writeln(f, "Required:    true,")
    if attr['name'] in FORCENEW_ATTRS:
        pretty_writeln(f, "ForceNew:    true,")
    if attr['name'] in COMPUTED_ATTRS or attr['name'] in COMPUTED_AND_OPTIONAL_ATTRS:
        pretty_writeln(f, "Computed:    true,")

    unshift()
    pretty_writeln(f, "},")

def write_func_header(f, resource, operation):
    f.write("\n")
    pretty_writeln(f, "func resourceNsxt%s%s(d *schema.ResourceData, m interface{}) error {" %
            (resource, operation))
    f.write("\n")
    shift()

def write_nsxclient(f):
    pretty_writeln(f, "nsxClient := m.(*%s.APIClient)\n" % SDK_PACKAGE_NAME)


def write_get_id(f):
    pretty_writeln(f, "id := d.Id()")
    pretty_writeln(f, "if id == \"\" {")
    shift()
    pretty_writeln(f, "return fmt.Errorf(\"Error obtaining logical object id\")")
    unshift()
    pretty_writeln(f, "}\n")


def write_error_check(f, resource, operation):
    if operation == "update":
        pretty_writeln(f, "if err != nil || resp.StatusCode == http.StatusNotFound {")
    else:
        pretty_writeln(f, "if err != nil {")

    shift()
    pretty_writeln(f, "return fmt.Errorf(\"Error during %s %s: " % (resource, operation) + '%v", err)')
    unshift()
    pretty_writeln(f, "}\n")

def write_object(f, resource, attrs, is_create=True):
    used_attrs = []
    for attr in attrs:
        if (is_create and attr['name'] == 'Revision') or attr['name'] in COMPUTED_ATTRS:
            # Revision is irrelevant in create
            continue

        used_attrs.append(attr['name'])
        fixed_name = get_attr_fixed_name(attr)
        if attr['name'] in VIP_GETTER_ATTRS:
            pretty_writeln(f, "%s := get%sFromSchema(d)" % (
                convert_name(attr['name']), attr['name']))
            continue

        if attr['type'] in TYPECAST_MAP:
            # type casting is needed
            pretty_writeln(f, "%s := %s(d.Get(\"%s\").(%s))" %
                    (convert_name(attr['name']),
                     attr['type'],
                     fixed_name,
                     TYPECAST_MAP[attr['type']]))
        else:
            pretty_writeln(f, "%s := d.Get(\"%s\").(%s)" %
                        (convert_name(attr['name']),
                         fixed_name,
                         attr['type']))

    pretty_writeln(f, "%s := %s.%s {" % (convert_name(resource), MANAGER_PACKAGE_NAME, resource))
    shift()
    for attr in used_attrs:
        pretty_writeln(f, "%s: %s," % (attr, convert_name(attr)))

    unshift()

    pretty_writeln(f, "}\n")

def write_create_func(f, resource, attrs, api_section):

    lower_resource = convert_name(resource)
    write_func_header(f, resource, "Create")

    write_nsxclient(f)

    write_object(f, resource, attrs)

    pretty_writeln(f, "%s, resp, err := nsxClient.%s.Create%s(nsxClient.Context, %s)" % (
        lower_resource, api_section, resource, lower_resource))

    f.write("\n")
    write_error_check(f, resource, "create")

    pretty_writeln(f, "if resp.StatusCode != http.StatusCreated {")
    shift()
    pretty_writeln(f, "return fmt.Errorf(\"Unexpected status returned during %s create: " % resource + '%v", resp.StatusCode)')
    unshift()
    pretty_writeln(f, "}")


    pretty_writeln(f, "d.SetId(%s.Id)\n" % lower_resource)

    pretty_writeln(f, "return resourceNsxt%sRead(d, m)" % resource)
    unshift()
    pretty_writeln(f, "}")


def write_read_func(f, resource, attrs, api_section):

    lower_resource = convert_name(resource)
    write_func_header(f, resource, "Read")

    write_nsxclient(f)
    write_get_id(f)

    # For some resources this is GET and for other it is read
    pretty_writeln(f, "//TerraGen TODO - select the right command for this resource, and delete this comment")
    pretty_writeln(f, "%s, resp, err := nsxClient.%s.Get%s(nsxClient.Context, id)" %
            (lower_resource, api_section, resource))
    pretty_writeln(f, "%s, resp, err := nsxClient.%s.Read%s(nsxClient.Context, id)" %
            (lower_resource, api_section, resource))

    pretty_writeln(f, "if resp.StatusCode == http.StatusNotFound {")
    shift()
    pretty_writeln(f, "fmt.Printf(\"%s " % resource + '%s not found\", id)')
    pretty_writeln(f, 'd.SetId("")')
    pretty_writeln(f, "return nil")
    unshift()
    pretty_writeln(f, "}")

    write_error_check(f, resource, "read")

    for attr in attrs:
        if attr['name'] in IGNORE_ATTRS:
            continue

        if attr['name'] in VIP_SETTER_ATTRS:
            pretty_writeln(f, "set%sInSchema(d, %s.%s)" % (
                attr['name'], lower_resource, attr['name']))
            continue

        fixed_name = get_attr_fixed_name(attr)
        pretty_writeln(f, "d.Set(\"%s\", %s.%s)" %
                (fixed_name, lower_resource, attr['name']))

    f.write("\n")
    pretty_writeln(f, "return nil")
    unshift()
    pretty_writeln(f, "}")


def write_update_func(f, resource, attrs, api_section):

    lower_resource = convert_name(resource)
    write_func_header(f, resource, "Update")

    write_nsxclient(f)
    write_get_id(f)

    write_object(f, resource, attrs, is_create=False)
    pretty_writeln(f, "%s, resp, err := nsxClient.%s.Update%s(nsxClient.Context, id, %s)" % (
        lower_resource, api_section, resource, lower_resource))

    f.write("\n")
    write_error_check(f, resource, "update")

    pretty_writeln(f, "return resourceNsxt%sRead(d, m)" % resource)
    unshift()
    pretty_writeln(f, "}")


def write_delete_func(f, resource, attrs, api_section):

    write_func_header(f, resource, "Delete")

    write_nsxclient(f)
    write_get_id(f)

    pretty_writeln(f, "//TerraGen TODO - select the right command for this resource, and delete this comment")
    pretty_writeln(f, "localVarOptionals := make(map[string]interface{})")
    pretty_writeln(f, "resp, err := nsxClient.%s.Delete%s(nsxClient.Context, id, localVarOptionals)" % (
        api_section, resource))
    pretty_writeln(f, "resp, err := nsxClient.%s.Delete%s(nsxClient.Context, id)" % (
        api_section, resource))

    write_error_check(f, resource, "delete")


    pretty_writeln(f, "if resp.StatusCode == http.StatusNotFound {")
    shift()
    pretty_writeln(f, "fmt.Printf(\"%s " % resource + '%s not found\", id)')
    pretty_writeln(f, 'd.SetId("")')
    unshift()
    pretty_writeln(f, "}")

    pretty_writeln(f, "return nil")
    unshift()
    pretty_writeln(f, "}")


def write_doc_header(f, resource):
    description = "Provides a resource to configure %s on NSX-T manager" % re.sub('_', ' ', resource)
    pretty_writeln(f, "---")
    pretty_writeln(f, "layout: \"nsxt\"")
    pretty_writeln(f, "page_title: \"NSXT: nsxt_%s\"" % resource)
    pretty_writeln(f, "sidebar_current: \"docs-nsxt-resource-%s\"" % re.sub('_', '-', resource))
    pretty_writeln(f, "description: |-")
    pretty_writeln(f, "  %s" % description)
    pretty_writeln(f, "---\n")
    pretty_writeln(f, "# nsxt_%s\n" % resource)
    pretty_writeln(f, "%s\n" % description)


def write_doc_example(f, resource, attrs):
    obj_name = ''.join(i[0] for i in resource.split('_')).upper()
    pretty_writeln(f, "## Example Usage\n")
    pretty_writeln(f, "```hcl")

    pretty_writeln(f, "resource \"nsxt_%s\" \"%s\" {" % (resource, obj_name))
    for attr in attrs:
        if attr['name'] == 'Revision' or attr['name'] in COMPUTED_ATTRS or attr['name'] in IGNORE_ATTRS:
            continue
        name = get_attr_fixed_name(attr)
        val = "..."
        eq = " = "
        if name == 'display_name':
            val = "\"%s\"" % obj_name
        elif name == 'description':
            val = "\"%s provisioned by Terraform\"" % obj_name
        elif name == 'tag':
            eq = ' '
            val = "{\n\tscope = \"color\"\n\ttag = \"red\"\n  }"
        pretty_writeln(f, "  %s%s%s" % (name, eq, val))

    pretty_writeln(f, "}")
    pretty_writeln(f, "```\n")


def write_arguments_reference(f, resource, attrs):
    pretty_writeln(f, "## Argument Reference\n")
    pretty_writeln(f, "The following arguments are supported:\n")
    for attr in attrs:
        if attr['name'] == 'Revision' or attr['name'] in COMPUTED_ATTRS or attr['name'] in IGNORE_ATTRS:
            continue
        name = get_attr_fixed_name(attr)
        desc = attr['comment']
        optional = 'Optional' if attr['optional'] else 'Required'
        if name == 'display_name':
            desc = "The display name of this resource. " + desc
        if name == 'tag':
            desc = "A list of scope + tag pairs to associate with this %s" % resource
        pretty_writeln(f, "* `%s` - (%s) %s." % (name, optional, desc))
    pretty_writeln(f, "\n")

def write_attributes_reference(f, resource, attrs):
    pretty_writeln(f, "## Attributes Reference\n")
    pretty_writeln(f, "In addition to arguments listed above, the following attributes are exported:\n")
    pretty_writeln(f, "* `id` - ID of the %s." % resource)
    for attr in attrs:
        if attr['name'] == 'Revision' or attr['name'] in COMPUTED_ATTRS or attr['name'] in IGNORE_ATTRS:
            name = convert_name(attr['name'])
            desc = attr['comment']
            if name == 'revision':
                desc = 'Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging'

            pretty_writeln(f, "* `%s` - %s." % (name, desc))


def main():

    if len(sys.argv) != 3:
        print("Usage: %s <sdk resource file> <api section>" % sys.argv[0])
        sys.exit()

    print("Building resource from %s" % sys.argv[1])
    api_section = sys.argv[2]

    with open(sys.argv[1], 'r') as f:
        lines = f.readlines()

    resource = None
    resource_started = False
    attr_comment = None
    attrs = []
    for line in lines:
        line = line.strip()
        match = re.match("type (.+?) struct", line)
        if match:
            resource = match.group(1)
            resource_started = True
            continue

        if not resource_started:
            continue

        if line.startswith('//'):
            attr_comment = line[3:]
            # remove dot if exists
            if attr_comment.endswith('.'):
                attr_comment = attr_comment[:-1]
            continue

        match = re.match("(.+?) (.+?) `json:\"(.+?)\"`", line)
        if match:
            attr_name = match.group(1).strip()
            attr_type = match.group(2).strip()
            attr_meta = match.group(3).strip()
            attr_optional = 'omitempty' in attr_meta
            if attr_name not in IGNORE_ATTRS:
                attrs.append({'name': attr_name,
                              'type': attr_type,
                              'meta': attr_meta,
                              'comment': attr_comment,
                              'optional': attr_optional})
            attr_comment = None


    print("Resource: %s" % resource)
    resource_lower = convert_name(resource)
    print(resource_lower)

    # write the resource file
    with open("resource_nsxt_%s.go" % resource_lower, 'w') as f:
        write_header(f)

        pretty_writeln(f, "func resourceNsxt%s() *schema.Resource {" % resource)
        shift()
        pretty_writeln(f, "return &schema.Resource{")
        shift()
        for op in ("Create", "Read", "Update", "Delete"):
            pretty_writeln(f, "%s: resourceNsxt%s%s," % (op, resource, op))

        f.write("\n")
        pretty_writeln(f, "Schema: map[string]*schema.Schema{")
        shift()

        for attr in attrs:
            write_attr(f, attr)

        unshift()
        pretty_writeln(f, "},")
        unshift()
        pretty_writeln(f, "}")
        unshift()
        pretty_writeln(f, "}")

        write_create_func(f, resource, attrs, api_section)
        write_read_func(f, resource, attrs, api_section)
        write_update_func(f, resource, attrs, api_section)
        write_delete_func(f, resource, attrs, api_section)

    # write the documentation file
    with open("%s.html.markdown" % resource_lower, 'w') as f:
        write_doc_header(f, resource_lower)
        write_doc_example(f, resource_lower, attrs)
        write_arguments_reference(f, resource_lower, attrs)
        write_attributes_reference(f, resource_lower, attrs)


main()
