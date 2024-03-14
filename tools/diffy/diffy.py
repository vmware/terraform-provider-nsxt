#!/usr/bin/python
# -*- coding: latin-1 -*-

#  Copyright (c) 2023 VMware, Inc. All Rights Reserved.
#  SPDX-License-Identifier: MPL-2.0

# This script detects difference in model between different versions of NSX spec. List of relevant objects is for now hard-coded.

import sys
import re
import os
import json

def load_types(path):
    with open(path, 'r') as f:
        return f.read().splitlines()

def load_spec(path):

    with open(path, 'r') as f:
        spec = json.load(f)

        obj_map = {}
        defs = spec["definitions"]
        for key in defs:
            if "allOf" not in defs[key]:
                if "properties" in defs[key]:
                    # no inheritance
                    obj_map[key] = defs[key]["properties"]
                continue
            if len(defs[key]["allOf"]) == 2:
                # object inheritance
                if "properties" not in defs[key]["allOf"][1]:
                    continue
                obj_map[key] = defs[key]["allOf"][1]["properties"]

        return obj_map

def ref_to_def(ref):
    return ref[len("#definitions/ "):]

def print_indent(text, level):
    ident = "  "*level
    print("%s%s" % (ident, text))

class color:
    PURPLE = '\033[95m'
    BLUE = '\033[94m'
    GREEN = '\033[92m'
    PURPLE = '\033[45m'
    END = '\033[0m'

# This script greps out list of 'model' types that are used in provider codebase, removes enum definitions, binding types and Child objects. Those will be used as baseline for further analysis.
object_scanner_script = "ls ../../nsxt/resource_nsxt_policy_*.go | grep -v test | xargs awk -F \"model.\" 'NF>1{ sub(/ .*/,\"\",$NF); print $NF }' | grep -v '_' | grep -v 'BindingType'| grep -v ListResult | grep -v Child | sed -E 's/[^[:alnum:][:space:]]+//g' | sort | uniq"

def main():

    if len(sys.argv) < 3:
        print("Usage: %s <baseline.json> <target.json>" % sys.argv[0])
        sys.exit()

    baseline_file_path = sys.argv[1]
    target_file_path = sys.argv[2]
    baseline_map = load_spec(baseline_file_path)
    target_map = load_spec(target_file_path)
    level = 0

    def analyze_obj(obj, level):
        print_indent("analyzing %s.." % obj, level)
        if obj not in baseline_map:
            return
        if obj not in target_map:
            print_indent(color.PURPLE + "skipping type %s" % obj + color.END, level)
            return
        for attr in target_map[obj]:
            if attr not in baseline_map[obj]:
                print_indent(color.BLUE + "new attribute %s" % attr + color.END, level + 1)

            if "$ref" in target_map[obj][attr]:
                analyze_obj(ref_to_def(target_map[obj][attr]["$ref"]), level + 1)
            if "items" in target_map[obj][attr] and "$ref" in target_map[obj][attr]["items"]:
                analyze_obj(ref_to_def(target_map[obj][attr]["items"]["$ref"]), level + 1)
            if "enum" in target_map[obj][attr] and attr in baseline_map[obj]:
                target_set = set(target_map[obj][attr]["enum"])
                baseline_set = set(baseline_map[obj][attr]["enum"])
                diff = target_set - baseline_set
                if diff:
                    print_indent(color.GREEN + "new enum values for attribute %s: %s" % (attr, diff) + color.END, level + 1)

        for attr in baseline_map[obj]:
            if obj not in target_map:
                continue
            if attr not in target_map[obj]:
                print_indent(color.PURPLE + "deleted attribute %s" % attr + color.END, level + 1)

    types_file = "object_types.tmp"
    os.system("%s > %s" % (object_scanner_script, types_file))
    objects = load_types(types_file)
    for obj in objects:
        analyze_obj(obj, level + 1)


    print("Done.")

main()
