---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_idps_system_signatures"
description: A data source to read system IDPS signatures catalog.
---

# nsxt_policy_idps_system_signatures

This data source provides access to the comprehensive catalog of system IDPS (Intrusion Detection and Prevention System) signatures available in vDefend. It allows filtering and querying of specific system signatures for review, analysis, or cross-referencing against policies.

**Important**: The system signatures catalog contains more than 10,000 signatures. Fetching all signatures without filters can take significant time and resources. It is **strongly recommended** to always use appropriate filters (`severity`, `product_affected`, `display_name`, or `class_type`) to narrow down the results to the specific signatures you need.

* **Filter Combinations**: Combining multiple filters provides the most efficient queries. For example, filtering by both `severity` and `product_affected` will return a much smaller and more relevant result set.
* **Use Case**: This data source is primarily intended for querying specific signatures to reference in IDPS policies, not for bulk exports of the entire signature catalog.

This data source is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
# Filter by severity - recommended to always use filters
data "nsxt_policy_idps_system_signatures" "critical_sigs" {
  severity = "CRITICAL"
}

# Filter by product and severity for targeted results
data "nsxt_policy_idps_system_signatures" "web_sigs" {
  product_affected = "HTTP_SERVER"
  severity         = "HIGH"
}

# Filter by display name to find specific signature patterns
data "nsxt_policy_idps_system_signatures" "sql_sigs" {
  display_name = "SQL"
  severity     = "HIGH"
}

# Filter by class type for specific attack categories
data "nsxt_policy_idps_system_signatures" "web_attacks" {
  class_type = "web-application-attack"
  severity   = "CRITICAL"
}

# Combine multiple filters for precise results
data "nsxt_policy_idps_system_signatures" "targeted_sigs" {
  version_id       = "DEFAULT"
  severity         = "CRITICAL"
  product_affected = "HTTP_SERVER"
}

# Output signature IDs for use in policies
output "critical_signature_ids" {
  value = [for sig in data.nsxt_policy_idps_system_signatures.critical_sigs.signatures : sig.signature_id]
}
```

## Argument Reference

* `version_id` - (Optional) Signature version ID to query. If not specified, the provider uses the active signature version (or `DEFAULT`). This is also exported as an attribute so you can reference the actual version used (e.g. when omitted, `data.xxx.version_id` gives the resolved active version).
* `severity` - (Optional) Filter by severity level: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`.
* `product_affected` - (Optional) Filter by affected product (e.g., `HTTP_SERVER`, `DATABASE`).
* `display_name` - (Optional) Filter by display name (case-insensitive substring match).
* `class_type` - (Optional) Filter by class type (e.g., `web-application-attack`, `trojan-activity`).

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - The unique identifier for this data source query.
* `path` - The NSX policy path to the signatures collection.
* `signatures` - List of system signatures matching the filter criteria. Each signature contains:
    * `id` - The signature ID.
    * `signature_id` - Unique signature identifier (numeric ID).
    * `display_name` - Display name.
    * `name` - Signature name.
    * `severity` - Severity level.
    * `class_type` - Classification type.
    * `categories` - List of categories (e.g., `APPLICATION`, `NETWORK`).
    * `product_affected` - Affected product.
    * `attack_target` - Attack target.
    * `cvss` - CVSS severity rating.
    * `cvssv2` - CVSS v2 score.
    * `cvssv3` - CVSS v3 score.
    * `cves` - List of associated CVE IDs.
    * `urls` - List of reference URLs.
    * `signature_revision` - Revision number.
    * `path` - NSX policy path.
