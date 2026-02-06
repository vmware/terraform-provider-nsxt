---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_idps_system_signatures"
description: A data source to read system IDPS signatures catalog.
---

# nsxt_policy_idps_system_signatures

This data source provides access to the comprehensive catalog of system IDPS (Intrusion Detection and Prevention System) signatures available in vDefend. It allows filtering and querying of system signatures for review, analysis, or cross-referencing against policies.

This data source is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
# Get all system signatures
data "nsxt_policy_idps_system_signatures" "all" {
}

# Filter by severity
data "nsxt_policy_idps_system_signatures" "critical_sigs" {
  severity  = "CRITICAL"
  page_size = 100
}

# Filter by product and severity
data "nsxt_policy_idps_system_signatures" "web_sigs" {
  product_affected = "HTTP_SERVER"
  severity         = "HIGH"
  page_size        = 100
}

# Filter by display name
data "nsxt_policy_idps_system_signatures" "sql_sigs" {
  display_name = "SQL"
  page_size    = 100
}

# Query specific version
data "nsxt_policy_idps_system_signatures" "default_version" {
  version_id = "DEFAULT"
}

# Output signature IDs for use in policies
output "critical_signature_ids" {
  value = [for sig in data.nsxt_policy_idps_system_signatures.critical_sigs.signatures : sig.signature_id]
}
```

## Argument Reference

* `version_id` - (Optional) Signature version ID to query. If not specified, uses the active version.
* `severity` - (Optional) Filter by severity level: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`.
* `product_affected` - (Optional) Filter by affected product (e.g., `HTTP_SERVER`, `DATABASE`).
* `display_name` - (Optional) Filter by display name (case-insensitive substring match).
* `class_type` - (Optional) Filter by class type (e.g., `web-application-attack`, `trojan-activity`).
* `page_size` - (Optional) API pagination page size. Defaults to 1000, valid range: 1-1000. See Notes for performance guidance.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - The unique identifier for this data source query.
* `path` - The NSX policy path to the signatures collection.
* `version_id` - The signature version ID being queried.
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

## Notes

* This data source provides read-only access to the vendor-provided signature database managed by NSX.
* Use signature IDs from this data source to reference signatures in IDPS profiles and policies.
* The signature catalog is automatically updated when signature updates are enabled in NSX.
* Filtering is performed client-side after fetching from the API.
* **Performance:** For large signature catalogs (10,000+ signatures), queries can take several minutes. To improve performance:
    * Use a smaller `page_size` (e.g., 100) to reduce API response time per page
    * Apply filters (`severity`, `display_name`, etc.) to reduce the result set
    * Smaller page sizes work best with filters, reducing data fetched and processed client-side
