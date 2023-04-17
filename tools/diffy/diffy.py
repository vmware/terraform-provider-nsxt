#!/usr/bin/python
# -*- coding: latin-1 -*-

#  Copyright (c) 2023 VMware, Inc. All Rights Reserved.
#  SPDX-License-Identifier: MPL-2.0

# This script detects difference in model between different versions of NSX spec. List of relevant objects is for now hard-coded.

import sys
import re
import os
import json

OBJECTS = ["Segment", "StaticRoutes", "QoSProfile", "L2VPNService", "IPSecVpnTunnelProfile", "IPSecVpnService", "IPSecVpnDpdProfile", "Service", "Group", "EvpnConfig", "EvpnTenantConfig", "Tier0", "Tier0Interface", "Tier1", "Tier1Interface", "GatewayPolicy", "SecurityPolicy", "LBPool", "LBService", "ContextProfile", "IpAddressPool", "PrefixList", "PolicyDnsForwarder", "PolicyNatRule", "OspfConfig", "MacDiscoveryProfile", "IpAddressBlock", "DhcpRelayConfig", "DhcpServerConfig", "IdsProfile", "IdsPolicy", "Domain", "L2VPNSession"]


def load_spec(path):

    with open(path, 'r') as f:
        spec = json.load(f)

        obj_map = {}
        for key in spec["definitions"]:
            if "allOf" not in spec["definitions"][key]:
                if "properties" in spec["definitions"][key]:
                    # no inheritance
                    obj_map[key] = spec["definitions"][key]["properties"]
                continue
            if len(spec["definitions"][key]["allOf"]) == 2:
                # object inheritance
                if "properties" not in spec["definitions"][key]["allOf"][1]:
                    continue
                obj_map[key] = spec["definitions"][key]["allOf"][1]["properties"]

        return obj_map

def ref_to_def(ref):
    return ref[len("#definitions/ "):]

def print_ident(text, level):
    ident = "  "*level
    print("%s%s" % (ident, text))

class color:
    PURPLE = '\033[95m'
    BLUE = '\033[94m'
    END = '\033[0m'

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
        print_ident("analyzing %s.." % obj, level)
        if obj not in baseline_map:
            return
        for attr in target_map[obj]:
            if attr not in baseline_map[obj]:
                print_ident(color.BLUE + "new attribute %s" % attr + color.END, level + 1)

            if "$ref" in target_map[obj][attr]:
                analyze_obj(ref_to_def(target_map[obj][attr]["$ref"]), level + 1)
            if "items" in target_map[obj][attr] and "$ref" in target_map[obj][attr]["items"]:
                analyze_obj(ref_to_def(target_map[obj][attr]["items"]["$ref"]), level + 1)

        for attr in baseline_map[obj]:
            if obj not in target_map:
                continue
            if attr not in target_map[obj]:
                print_ident(color.PURPLE + "deleted attribute %s" % attr + color.END, level + 1)

    for obj in OBJECTS:
        analyze_obj(obj, level + 1)


    print("Done.")

main()
