---
page_title: "In-Memory Resource Caching"
description: |-
  In-Memory Resource Caching
---

# In-Memory Resource Caching

The provider can cache NSX resource lookups in memory for the duration of a single `terraform plan`/`apply`/`destroy` run, avoiding repeated NSX Search API calls for resources of the same type. This is controlled by the `cache_mode` provider argument (or the `NSXT_CACHE_MODE` environment variable) and is **disabled by default**.

```hcl
provider "nsxt" {
  host       = "192.168.110.41"
  username   = "admin"
  password   = "default"
  cache_mode = "config_scope" # or "global"
  context_id = "team-a-prod"  # required to get real isolation in config_scope mode
}
```

The cache lives only inside a single provider process, which corresponds to a single `terraform` invocation. It is never persisted to disk and is discarded as soon as the run ends, so it cannot cause cross-run staleness — only in-run staleness (see [Troubleshooting](#troubleshooting) below).

## Modes

### `disabled` (default)

No caching. Every resource/data-source read goes straight to NSX, exactly as before this feature existed. Also accepted: `off`, or simply leaving `cache_mode` unset.

### `global`

Caches resource lookups for the whole run, without any tagging or per-config scoping. The first time the provider needs to look up a resource of a given type (e.g. `Group`) within a given scope (NSX manager host + Local/Global Manager/Project/VPC context), it fetches **all** resources of that type in that scope via NSX Search and keeps them in memory; subsequent lookups of that type (by ID or display name) in the same run are served from that in-memory set instead of hitting NSX again.

`global` mode does not write anything to NSX — it is purely a read-side optimization.

### `config_scope`

Like `global` mode, but scopes the cache to only the resources managed by this specific provider configuration, using `context_id`:

- Cache population is filtered to objects tagged `nsx-tf/tf-run-id: <context_id>` (in addition to the resource type), instead of fetching every resource of that type in scope.
- **The provider will write that tracking tag to NSX objects it reads that don't already have it.** Concretely: the first time a cache-integrated resource is read (refreshed) in `config_scope` mode and its NSX object is missing the `nsx-tf/tf-run-id` tag, the provider issues a PATCH to NSX to add it — even for objects that pre-date this Terraform config, or that were created outside Terraform entirely. This is a genuine write side effect of a read operation.
- That tracking tag is stripped from Terraform state before it reaches your configuration's `tag` attribute, `terraform plan` diffs, or `terraform state show` output — you will not see it there. It **is** visible if you inspect the object directly via the NSX UI or API.
- If `context_id` is left unset, the tag is never added (there is nothing to scope by), and `config_scope` mode behaves like `global` mode functionally, just with the extra (skipped, harmless) bookkeeping overhead of checking for a tag that will never be applied. **Always set `context_id` when using `config_scope`** — otherwise you get no isolation benefit over `global` mode.

Use `config_scope` when several independent Terraform configurations (or teams/pipelines) manage different resources of the same type against the same NSX manager, and you want each configuration's cache to only "see" the resources it manages — at the cost of the tag-write side effect above. Use `global` when a single configuration owns everything of a given type in that scope, or when you cannot accept the tag being written to your objects.

## Differences at a glance

| | `disabled` | `global` | `config_scope` |
|---|---|---|---|
| Caches lookups within a run | No | Yes | Yes |
| Scope of cached data | N/A | All resources of that type in scope | Only resources tagged for this `context_id` |
| Writes tags to NSX objects | No | No | Yes (`nsx-tf/tf-run-id`), on first read if missing |
| Requires `context_id` | No | No | Effectively yes, for isolation to work |
| Visible in Terraform state | N/A | N/A | Tracking tag is always stripped from state |

Both `global` and `config_scope` share one more property: the cache is keyed by NSX manager host as well as type/scope, so multiple aliased `provider "nsxt"` blocks pointing at different NSX managers in the same `terraform` run never share cached data, even if they otherwise have identical `context_id`/resource types.

## Troubleshooting

**How do I confirm caching is active and see hits/misses?** Run with `TF_LOG=DEBUG` (or `TF_LOG_PROVIDER=DEBUG`) and look for lines like:

```text
[DEBUG] Cache hit: resourceType=Group id=... query=... (hit=3 miss=1)
[DEBUG] Cache lookup miss: resourceType=Group id=... query=...
[DEBUG] Cache post-write bypass: direct GET for resourceType=Group id=...
[DEBUG] Cache fallback: direct GET for resourceType=Group id=...
```

**I'm seeing a value that doesn't match what's actually in NSX right now.** Because the cache only lives for one `terraform` run, this can only happen *within* a single plan/apply/destroy — for example, if something else changes an object in NSX in parallel while your run is in flight. If you suspect the cache itself (rather than a real concurrent change) is the cause, re-run with `cache_mode` unset/`disabled` and compare; if the problem disappears, please file an issue with the debug log.

**I don't want any tags written to my NSX objects.** Use `global` mode instead of `config_scope`, or disable caching entirely.

**I'm using `config_scope` but don't see any caching benefit.** Check that `context_id` is actually set (via the provider argument or `NSXT_CONTEXT_ID`) — without it, no objects ever match the scoping tag on first population, so most reads still fall through to NSX. Set a stable, unique `context_id` per configuration.

**I'm running multiple aliased `nsxt` provider blocks in one run.** Each aliased provider's `cache_mode` is independent — one alias can use `config_scope` while another uses `global` or `disabled`. If two aliases point at the same NSX manager host with the same `context_id` and resource scope, they can legitimately share cached data for that overlap; this is intentional (it's the same underlying NSX state), not a bug.
