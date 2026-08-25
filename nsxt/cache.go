package nsxt

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// resourceTypeCache holds every populated search bucket for one NSX resource type. The
// outer data map is keyed by the query string (see getCacheQueryKey) that produced a given
// bucket — the same resource type can have more than one bucket if it's queried under
// different tag filters or provider configs. The inner map is that bucket's contents,
// keyed by whatever CacheKeyForResourceID/data-source id was used to look objects up
// (see converListToMapByType for how an object ends up under multiple keys).
type resourceTypeCache struct {
	mu       sync.RWMutex
	data     map[string]map[string]*data.StructValue // key: query, key: resource ID
	cacheHit atomic.Int64
	cacheMis atomic.Int64
}

// typeScopedCache is the process-wide root: one resourceTypeCache per NSX resource type
// (see gcache). Its own mutex only protects the byTyp map itself (creating/looking up a
// type's cache); it is never held while a search runs or a bucket is read, so looking up
// one resource type never blocks concurrent access to another.
type typeScopedCache struct {
	mu    sync.RWMutex
	byTyp map[string]*resourceTypeCache
}

// Resource type identifiers used as cache keys and interpolated into NSX Search
// queries (resource_type:<value>). Values must match NSX Search's resource_type
// field exactly, including casing, and must not be changed here.
const (
	resourceTypeConnectivityPolicy        = "ConnectivityPolicy"
	resourceTypeDhcpV4StaticBindingConfig = "DhcpV4StaticBindingConfig"
	resourceTypeDhcpV6StaticBindingConfig = "DhcpV6StaticBindingConfig"
	// resourceTypeGatewayPolicy is shared by the Tier-0/Tier-1 gateway policy resource and the
	// VPC gateway policy resource: both cache under this one bucket. Same sharing applies to
	// resourceTypeStaticRoutes and resourceTypeDhcpV4StaticBindingConfig below. Keep these
	// path-indexed (see shouldIndexByPath) since the VPC-scoped side of each pair can collide
	// on short id once its cache bucket is project-wide-shared.
	resourceTypeGatewayPolicy = "gatewaypolicy"
	resourceTypeGroup         = "Group"
	// resourceTypeVPCGroup is a distinct NSX Search resource_type ("group", lowercase) from
	// resourceTypeGroup ("Group") above — VPC-scoped groups and Tier-0/Tier-1-scoped groups are
	// never in the same bucket, unlike the shared types noted above.
	resourceTypeVPCGroup                                       = "group"
	resourceTypeIpAddressAllocation                            = "IpAddressAllocation"
	resourceTypeIpAddressPool                                  = "IpAddressPool"
	resourceTypeIpAddressPoolBlockSubnet                       = "IpAddressPoolBlockSubnet"
	resourceTypeIpAddressPoolStaticSubnet                      = "IpAddressPoolStaticSubnet"
	resourceTypeLBPool                                         = "LBPool"
	resourceTypeLBServerSslProfile                             = "LBServerSslProfile"
	resourceTypeLBService                                      = "LBService"
	resourceTypeLBSourceIpPersistenceProfile                   = "LBSourceIpPersistenceProfile"
	resourceTypeLBTcpMonitorProfile                            = "LBTcpMonitorProfile"
	resourceTypeLBUdpMonitorProfile                            = "LBUdpMonitorProfile"
	resourceTypeLBVirtualServer                                = "LBVirtualServer"
	resourceTypePolicyFirewallFloodProtectionProfileBindingMap = "PolicyFirewallFloodProtectionProfileBindingMap"
	resourceTypePolicyNatRule                                  = "PolicyNatRule"
	resourceTypePolicyVpcNatRule                               = "PolicyVpcNatRule"
	resourceTypeProjectIpAddressAllocation                     = "ProjectIpAddressAllocation"
	resourceTypeRule                                           = "rule"
	resourceTypeSecurityPolicy                                 = "securitypolicy"
	resourceTypeSegment                                        = "Segment"
	resourceTypeSegmentPort                                    = "SegmentPort"
	resourceTypeService                                        = "Service"
	resourceTypeStaticRoutes                                   = "StaticRoutes"
	resourceTypeTier1                                          = "Tier1"
	resourceTypeTier1Interface                                 = "Tier1Interface"
	resourceTypeTransitGateway                                 = "TransitGateway"
	resourceTypeTransitGatewayAttachment                       = "TransitGatewayAttachment"
	resourceTypeVpc                                            = "Vpc"
	resourceTypeVpcAttachment                                  = "VpcAttachment"
	resourceTypeVpcConnectivityProfile                         = "VpcConnectivityProfile"
	resourceTypeVpcIpAddressAllocation                         = "VpcIpAddressAllocation"
	resourceTypeVpcServiceProfile                              = "VpcServiceProfile"
	resourceTypeVpcSubnet                                      = "VpcSubnet"
)

var gcache = &typeScopedCache{byTyp: make(map[string]*resourceTypeCache)}

// postWriteByKey tracks recently written resources to bypass cache once on the next read.
var postWriteByKey sync.Map // string -> struct{}

// errCacheUseBackendDirect indicates the cache bucket exists but the ID is missing, so use direct GET.
var errCacheUseBackendDirect = errors.New("nsxt cache: use direct API read")

// cacheLogEvent defers a log message describing bulk-populate work until after the caller
// releases tc.mu, so exclusive-lock hold time and contention on log's internal mutex don't
// grow with the number of DEBUG/WARNING lines a populate call produces.
type cacheLogEvent struct {
	level string // "DEBUG" or "WARNING"
	msg   string
}

func logCacheEvents(events []cacheLogEvent) {
	for _, e := range events {
		log.Printf("[%s] %s", e.level, e.msg) //nolint:gosec
	}
}

// compositeCacheEntry defines a parent type with an extra child search and merge step.
type compositeCacheEntry struct {
	childSearchType string
	merge           func(parents, children []*data.StructValue) ([]*data.StructValue, error)
}

// compositeCacheRegistry lists resource types whose NSX Search result doesn't already embed
// everything the Terraform schema needs: GatewayPolicy/SecurityPolicy objects come back from
// Search without their Rules, so a second search for the child rule type is required, and the
// two result sets have to be merged (by parent path) before the combined object can be cached.
var compositeCacheRegistry = map[string]compositeCacheEntry{
	resourceTypeGatewayPolicy: {
		childSearchType: resourceTypeRule,
		merge:           mergeGatewayPolicyCacheSearchResults,
	},
	resourceTypeSecurityPolicy: {
		childSearchType: resourceTypeRule,
		merge:           mergeSecurityPolicyCacheSearchResults,
	},
}

const envNSXTCacheMode = "NSXT_CACHE_MODE"

type cacheMode int

const (
	cacheDisabled cacheMode = iota
	cacheConfigScoped
	cacheGlobal
)

var invalidNSXTCacheModeWarn sync.Map // raw env value -> struct{}; log once per distinct invalid value

// sharedTypeConverter is reused across all cache conversions: bindings.TypeConverter is a
// stateless, field-less struct (all conversion state lives in per-call visitors), so a single
// instance is safe to share across goroutines and avoids an allocation on every cached
// object read/write at 100+ resource scale.
var sharedTypeConverter = bindings.NewTypeConverter()

func postWriteKey(resourceType, resourceID string) string {
	return resourceType + "::" + resourceID
}

// MarkPostWriteForResourceTypeKey marks a single resource instance as just written.
// The next CacheAwareResourceRead/TryCacheRead for that exact resourceID will bypass cache once.
func MarkPostWriteForResourceTypeKey(resourceType, resourceID string, m interface{}) {
	if !IsCacheEnabled(m) {
		return
	}
	if resourceType == "" || resourceID == "" {
		return
	}
	postWriteByKey.Store(postWriteKey(resourceType, resourceID), struct{}{})
}

// MarkPostWriteAndInvalidateCacheForResourceType marks resourceID as just written and invalidates
// the type's cache bucket. resourceID must be the same key used in the corresponding
// CacheAwareResourceRead/TryCacheRead call (typically d.Id(), or the parent's ID for child resources).
func MarkPostWriteAndInvalidateCacheForResourceType(resourceType, resourceID string, m interface{}) {
	MarkPostWriteForResourceTypeKey(resourceType, resourceID, m)
	InvalidateCacheForResourceType(resourceType, m)
}

// parseCacheModeString parses the raw cache_mode string (from
// commonProviderConfig.CacheMode, itself resolved from the cache_mode provider
// attribute / NSXT_CACHE_MODE env var via schema.EnvDefaultFunc in provider.go)
// into a cacheMode value.
func parseCacheModeString(raw string) cacheMode {
	raw = strings.TrimSpace(raw)
	lower := strings.ToLower(raw)
	switch lower {
	case "", "disabled", "off":
		return cacheDisabled
	case "config_scope":
		return cacheConfigScoped
	case "global":
		return cacheGlobal
	default:
		if raw != "" {
			if _, loaded := invalidNSXTCacheModeWarn.LoadOrStore(raw, struct{}{}); !loaded {
				// strconv.Quote neutralizes newlines/control chars in env (gosec G706).
				log.Printf("[ERROR]: Invalid %s value %s; expected disabled, off, config_scope, or global (caching disabled)", envNSXTCacheMode, strconv.Quote(raw)) //nolint:gosec
			}
		}
		return cacheDisabled
	}
}

// currentCacheMode reads the cache_mode configured on this specific provider instance
// (commonProviderConfig.CacheMode, scoped per nsxtClients) rather than any process-global
// state, so aliased provider blocks with different cache_mode values don't clobber each other.
func currentCacheMode(m interface{}) cacheMode {
	raw := m.(nsxtClients).CommonConfig.CacheMode
	return parseCacheModeString(raw)
}

func IsCacheEnabled(m interface{}) bool {
	return currentCacheMode(m) != cacheDisabled
}

func isGlobalSearchCacheMode(m interface{}) bool {
	return currentCacheMode(m) == cacheGlobal
}

// isConfigScopedCacheMode reports whether tag-filtered caching (with provider-managed tags) is active.
func isConfigScopedCacheMode(m interface{}) bool {
	return currentCacheMode(m) == cacheConfigScoped
}

func isRefreshPhase(d *schema.ResourceData) bool {
	return d.Id() != ""
}

func isCacheEnabledForRead(d *schema.ResourceData, m interface{}) bool {
	return IsCacheEnabled(m) && isRefreshPhase(d)
}

// CacheKeyForResourceID returns the key that CacheAwareResourceRead/TryCacheRead will use
// for this resource. Path-indexed types key by path (when set); all others key by id.
// Use this when calling MarkPostWriteAndInvalidateCacheForResourceType to ensure the keys match.
func CacheKeyForResourceID(resourceType string, d *schema.ResourceData) string {
	if shouldIndexByPath(resourceType) {
		if p, ok := d.GetOk("path"); ok {
			if path, ok2 := p.(string); ok2 && path != "" {
				return path
			}
		}
	}
	return d.Id()
}

// VPC-scoped resource types are path-indexed because their short NSX id (often
// user-chosen via nsx_id) is only guaranteed unique within its own VPC, not project-wide,
// while the cache populate search/bucket for utl.VPC/utl.Multitenancy contexts is now
// shared across all VPCs in a project (see projectScopedSearchContext). Path is unique
// project/system-wide, so keying by path avoids one VPC's resource shadowing another's.
// Known gap: data source cache reads (cacheAwareDataSourceReadByID) key off a short
// user-supplied id with no path available, so they remain exposed to this collision.
func shouldIndexByPath(resourceType string) bool {
	switch resourceType {
	case resourceTypePolicyNatRule, resourceTypeSegmentPort, resourceTypeIpAddressPoolStaticSubnet, resourceTypeIpAddressPoolBlockSubnet, resourceTypeService,
		resourceTypeVpc, resourceTypeVpcAttachment, resourceTypeVpcConnectivityProfile, resourceTypeVpcIpAddressAllocation, resourceTypeVpcServiceProfile,
		resourceTypeVpcSubnet, resourceTypeTransitGateway, resourceTypeTransitGatewayAttachment, resourceTypeProjectIpAddressAllocation, resourceTypePolicyVpcNatRule,
		resourceTypeVPCGroup, resourceTypeGatewayPolicy, resourceTypeStaticRoutes, resourceTypeDhcpV4StaticBindingConfig:
		return true
	default:
		return false
	}
}

// indexCacheMapKey inserts key→obj if key is not already present (first-seen wins).
// Used for both primary keys (id) and secondary keys (display_name, path).
func indexCacheMapKey(ret map[string]*data.StructValue, key string, obj *data.StructValue) {
	if _, seen := ret[key]; seen {
		return
	}
	ret[key] = obj
}

func getStructValueStringField(obj *data.StructValue, field string) (string, bool) {
	if obj == nil {
		return "", false
	}
	fields := obj.Fields()
	v, ok := fields[field]
	if !ok || v == nil {
		return "", false
	}
	if sv, ok := v.(*data.StringValue); ok {
		return sv.Value(), true
	}
	if ov, ok := v.(*data.OptionalValue); ok {
		if !ov.IsSet() {
			return "", false
		}
		inner := ov.Value()
		if sv, ok := inner.(*data.StringValue); ok {
			return sv.Value(), true
		}
	}
	return "", false
}

func converListToMapByType(list []*data.StructValue, resourceType string) map[string]*data.StructValue {
	ret := make(map[string]*data.StructValue)
	for _, obj := range list {
		dataValue, errors := sharedTypeConverter.ConvertToGolang(obj, model.PolicyConfigResourceBindingType())
		if len(errors) > 0 {
			return nil
		}
		resource := dataValue.(model.PolicyConfigResource)

		// Primary indexing via PolicyConfigResource conversion.
		idIndexed := resource.Id != nil
		if idIndexed {
			id := *resource.Id
			indexCacheMapKey(ret, id, obj)
			// Index both raw and extracted policy IDs when Search returns a full path.
			if strings.Contains(id, "/") {
				indexCacheMapKey(ret, getPolicyIDFromPath(id), obj)
			}
		}
		displayNameIndexed := resource.DisplayName != nil
		if displayNameIndexed {
			indexCacheMapKey(ret, *resource.DisplayName, obj)
		}
		indexByPath := shouldIndexByPath(resourceType)
		pathIndexed := indexByPath && resource.Path != nil
		if pathIndexed {
			indexCacheMapKey(ret, *resource.Path, obj)
		}

		// Fall back to raw StructValue fields only when the typed conversion above didn't
		// already index that field, to cover cases where SDK conversion loses fields.
		if !idIndexed {
			if id, ok := getStructValueStringField(obj, "id"); ok {
				indexCacheMapKey(ret, id, obj)
				if strings.Contains(id, "/") {
					indexCacheMapKey(ret, getPolicyIDFromPath(id), obj)
				}
			}
		}
		if !displayNameIndexed {
			if dn, ok := getStructValueStringField(obj, "display_name"); ok {
				indexCacheMapKey(ret, dn, obj)
			}
		}
		if indexByPath && !pathIndexed {
			if p, ok := getStructValueStringField(obj, "path"); ok {
				indexCacheMapKey(ret, p, obj)
			}
		}
	}
	return ret
}

func getQueryString(resourceType string, context utl.SessionContext) string {
	switch context.ClientType {
	case utl.Global:
		return fmt.Sprintf("resource_type:%s AND marked_for_delete:false AND context:Global", resourceType)
	case utl.Local:
		return fmt.Sprintf("resource_type:%s AND marked_for_delete:false AND context:Local", resourceType)
	case utl.VPC, utl.Multitenancy:
		// Scoped to the project only (VPCID intentionally omitted): NSX policy paths/IDs are
		// unique across VPCs within a project, and searchMultitenancyResources already searches
		// the whole project in one call when VPCID is empty. Narrowing the cache bucket to a
		// single VPC here caused a fresh bucket (and full search) per VPC, which regressed cache
		// mode below no-cache performance when a run touches many VPCs (e.g. 100 VpcAttachments
		// across 100 VPCs).
		return fmt.Sprintf("resource_type:%s AND marked_for_delete:false AND context:%s", resourceType, context.ProjectID)
	default:
		return fmt.Sprintf("resource_type:%s AND marked_for_delete:false", resourceType)
	}
}

// projectScopedSearchContext strips VPCID from a VPC/Multitenancy context so the cache
// populate search covers the whole project in one call instead of one call per VPC. Only
// the search/query scope is widened; the resource lookup within the populated bucket still
// keys by resourceID, which remains unique within the project.
func projectScopedSearchContext(context utl.SessionContext) utl.SessionContext {
	if context.ClientType == utl.VPC || context.ClientType == utl.Multitenancy {
		context.VPCID = ""
	}
	return context
}

func (c *resourceTypeCache) getQueryResult(query string, resourceID string) (*data.StructValue, error) {
	if inner, ok := c.data[query]; ok {
		// converListToMapByType indexes every object under its raw "id" field (in addition to
		// the typed id) at populate time, so a miss here means no O(n) scan can find it either.
		if v := inner[resourceID]; v != nil {
			return v, nil
		}
		return nil, errCacheUseBackendDirect
	}
	return nil, fmt.Errorf("element is not found")
}

func (c *resourceTypeCache) writeCache(query string, resourceType string, d *schema.ResourceData, m interface{}, connector client.Connector) ([]cacheLogEvent, error) {
	if _, ok := c.data[query]; ok {
		return []cacheLogEvent{{"DEBUG", fmt.Sprintf("Cache skip bulk refill for resourceType=%s query=%q (bucket already present)", resourceType, query)}}, nil
	}
	runID := m.(nsxtClients).CommonConfig.contextID
	events := []cacheLogEvent{{"DEBUG", fmt.Sprintf("Cache miss: populating cache for resourceType=%s query=%q", resourceType, query)}}
	childEvents, err := c.getListOfPolicyResources(query, d, m, connector, getEffectiveCacheContext(d, m), resourceType, runID)
	events = append(events, childEvents...)
	return events, err
}

// getTypeCache returns (creating if needed) the resourceTypeCache for resourceType. Held only
// long enough to read or insert the map entry, so this never serializes callers working with
// different resource types, and never blocks on a search/populate happening under the returned
// resourceTypeCache's own (separate) mutex.
func (c *typeScopedCache) getTypeCache(resourceType string) *resourceTypeCache {
	c.mu.Lock()
	defer c.mu.Unlock()
	if tc, ok := c.byTyp[resourceType]; ok {
		return tc
	}
	tc := &resourceTypeCache{data: make(map[string]map[string]*data.StructValue)}
	c.byTyp[resourceType] = tc
	return tc
}

// getEffectiveCacheContext derives the cache context, preferring parent_path when context{} is absent.
func getEffectiveCacheContext(d *schema.ResourceData, m interface{}) utl.SessionContext {
	if pp, ok := d.GetOk("parent_path"); ok {
		if parentPath := pp.(string); parentPath != "" {
			return getParentContext(d, m, parentPath)
		}
	}
	return getSessionContext(d, m)
}

// getCacheQueryKey returns the string that selects which bucket in resourceTypeCache.data a
// resource read hits. Two reads share a bucket (and its search results) only when this key is
// identical, so anything that should force a distinct search — the resource's own tag filter,
// the run's context_id, the manager host — is folded in here rather than compared separately
// after the fact.
func getCacheQueryKey(resourceType string, d *schema.ResourceData, m interface{}) string {
	clients := m.(nsxtClients)
	context := getEffectiveCacheContext(d, m)
	query := getQueryString(resourceType, context)
	runID := clients.CommonConfig.contextID
	additionalQuery := buildTagQuery(d, runID, m)
	// Prefix with manager host so aliased provider blocks pointing at different managers
	// never share a cache bucket even when context_id and resource types are identical.
	host := clients.Host
	base := fmt.Sprintf("[%s]%s", host, query)
	if additionalQuery == "" {
		return base
	}
	return fmt.Sprintf("%s AND %s ", base, additionalQuery)
}

func (c *typeScopedCache) readCache(resourceID string, resourceType string, d *schema.ResourceData, m interface{}, connector client.Connector) (interface{}, error) {
	tc := c.getTypeCache(resourceType)
	query := getCacheQueryKey(resourceType, d, m)

	// Fast path: read lock allows concurrent hits without blocking each other.
	tc.mu.RLock()
	val, qerr := tc.getQueryResult(query, resourceID)
	tc.mu.RUnlock()
	if val != nil {
		hit := tc.cacheHit.Add(1)
		log.Printf("[DEBUG] Cache hit: resourceType=%s id=%s query=%q (hit=%d miss=%d)", resourceType, resourceID, query, hit, tc.cacheMis.Load()) //nolint:gosec
		return val, nil
	}
	if errors.Is(qerr, errCacheUseBackendDirect) {
		miss := tc.cacheMis.Add(1)
		log.Printf("[DEBUG] Cache lookup miss: resourceType=%s id=%s query=%q (hit=%d miss=%d)", resourceType, resourceID, query, tc.cacheHit.Load(), miss) //nolint:gosec
		return nil, errCacheUseBackendDirect
	}

	// Slow path: take exclusive lock to populate the bucket, then double-check.
	// All logging below is accumulated into populateEvents and flushed by the deferred
	// logCacheEvents call, which is registered before (and therefore, by defer's LIFO
	// order, runs after) tc.mu.Unlock — so no log.Printf in this path ever executes while
	// tc.mu is held, regardless of which return point is taken.
	var populateEvents []cacheLogEvent
	defer func() { logCacheEvents(populateEvents) }()
	tc.mu.Lock()
	defer tc.mu.Unlock()
	val, qerr = tc.getQueryResult(query, resourceID)
	if val != nil {
		hit := tc.cacheHit.Add(1)
		populateEvents = append(populateEvents, cacheLogEvent{"DEBUG", fmt.Sprintf("Cache hit: resourceType=%s id=%s query=%q (hit=%d miss=%d)", resourceType, resourceID, query, hit, tc.cacheMis.Load())})
		return val, nil
	}
	if errors.Is(qerr, errCacheUseBackendDirect) {
		miss := tc.cacheMis.Add(1)
		populateEvents = append(populateEvents, cacheLogEvent{"DEBUG", fmt.Sprintf("Cache lookup miss: resourceType=%s id=%s query=%q (hit=%d miss=%d)", resourceType, resourceID, query, tc.cacheHit.Load(), miss)})
		return nil, errCacheUseBackendDirect
	}
	miss := tc.cacheMis.Add(1)
	var err error
	var writeCacheEvents []cacheLogEvent
	writeCacheEvents, err = tc.writeCache(query, resourceType, d, m, connector)
	populateEvents = append(populateEvents, writeCacheEvents...)
	if err != nil {
		populateEvents = append(populateEvents, cacheLogEvent{"DEBUG", fmt.Sprintf("Cache lookup miss: resourceType=%s id=%s query=%q (hit=%d miss=%d)", resourceType, resourceID, query, tc.cacheHit.Load(), miss)})
		return nil, err
	}
	val, _ = tc.getQueryResult(query, resourceID)
	populateEvents = append(populateEvents, cacheLogEvent{"DEBUG", fmt.Sprintf("Cache lookup miss: resourceType=%s id=%s query=%q (hit=%d miss=%d)", resourceType, resourceID, query, tc.cacheHit.Load(), miss)})
	if val != nil {
		return val, nil
	}
	return nil, errCacheUseBackendDirect
}

func (c *resourceTypeCache) getListOfPolicyResources(query string, d *schema.ResourceData, m interface{}, connector client.Connector, context utl.SessionContext, resourceType string, runID string) ([]cacheLogEvent, error) {
	var events []cacheLogEvent
	context = projectScopedSearchContext(context)
	additionalQuery := buildTagQuery(d, runID, m)
	events = append(events, cacheLogEvent{"DEBUG", fmt.Sprintf("Cache search query: resourceType=%s query=%q additionalQuery=%q", resourceType, query, additionalQuery)})
	resultList, err := listPolicyResources(connector, context, resourceType, &additionalQuery)
	events = append(events, cacheLogEvent{"DEBUG", fmt.Sprintf("Cache search results: resourceType=%s query=%q additionalQuery=%q results=%d", resourceType, query, additionalQuery, len(resultList))})
	if err != nil && len(resultList) == 0 {
		return events, fmt.Errorf("error listing resource %s %w", resourceType, err)
	}
	if err != nil {
		events = append(events, cacheLogEvent{"WARNING", fmt.Sprintf("Partial search results for resourceType=%s query=%q: %d parent results before error: %v", resourceType, query, len(resultList), err)})
	}

	entry, composite := compositeCacheRegistry[resourceType]
	if !composite {
		tmp := converListToMapByType(resultList, resourceType)
		if tmp == nil {
			return events, fmt.Errorf("error converting resources to cache map for resource type %s", resourceType)
		}
		c.data[query] = tmp
		return events, nil
	}

	// Child rules may be untagged, so try provider-managed tags first and then fall back to no tag filter.
	var childAdditional *string
	if !isGlobalSearchCacheMode(m) {
		if q := providerManagedTagsSearchQuery(runID); q != "" {
			childAdditional = &q
		}
	}
	childList, childErr := listPolicyResources(connector, context, entry.childSearchType, childAdditional)
	if childErr != nil && len(childList) == 0 {
		return events, fmt.Errorf("error listing composite child resource %s for parent %s: %w", entry.childSearchType, resourceType, childErr)
	}
	if childErr != nil {
		events = append(events, cacheLogEvent{"WARNING", fmt.Sprintf("Partial child search for resourceType=%s childType=%s query=%q: %d results before error: %v", resourceType, entry.childSearchType, query, len(childList), childErr)})
	}
	if childAdditional != nil && len(resultList) > 0 && len(childList) == 0 {
		childList, childErr = listPolicyResources(connector, context, entry.childSearchType, nil)
		if childErr != nil && len(childList) == 0 {
			return events, fmt.Errorf("error listing composite child resource %s for parent %s: %w", entry.childSearchType, resourceType, childErr)
		}
		if childErr != nil {
			events = append(events, cacheLogEvent{"WARNING", fmt.Sprintf("Partial child search (fallback) for resourceType=%s childType=%s query=%q: %d results before error: %v", resourceType, entry.childSearchType, query, len(childList), childErr)})
		}
	}

	mergedSVs, err := entry.merge(resultList, childList)
	if err != nil {
		return events, err
	}
	tmp := converListToMapByType(mergedSVs, resourceType)
	if tmp == nil {
		return events, fmt.Errorf("error converting merged resources to cache map for resource type %s", resourceType)
	}
	c.data[query] = tmp
	return events, nil
}

// convertCachedValue converts a raw cache value to the typed model, strips provider-managed tags,
// and returns (typedPtr, true) on success or (nil, false) on conversion/type-assertion failure.
func convertCachedValue[T any](val interface{}, resourceType, resourceID string, bindingType bindings.BindingType, m interface{}) (*T, bool) {
	sv, ok := val.(*data.StructValue)
	if !ok {
		return nil, false
	}
	goVal, convErrs := sharedTypeConverter.ConvertToGolang(sv, bindingType)
	if len(convErrs) > 0 {
		log.Printf("[WARNING] Cache: conversion failed for resourceType=%s id=%s (%v); discarding cached value", resourceType, resourceID, convErrs[0]) //nolint:gosec
		return nil, false
	}
	typedVal, ok := goVal.(T)
	if !ok {
		log.Printf("[WARNING] Cache: type assertion failed for resourceType=%s id=%s; discarding cached value", resourceType, resourceID) //nolint:gosec
		return nil, false
	}
	if !isGlobalSearchCacheMode(m) {
		stripProviderManagedTagsFromAny(&typedVal)
	}
	return &typedVal, true
}

// reflectStringField returns the value of a *string field named fieldName on obj (a
// pointer to struct), or nil if the field doesn't exist, isn't a *string, or is nil.
func reflectStringField(obj interface{}, fieldName string) *string {
	v := reflect.ValueOf(obj)
	for v.IsValid() && v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return nil
	}
	f := v.FieldByName(fieldName)
	if !f.IsValid() || f.Kind() != reflect.Pointer || f.Type().Elem().Kind() != reflect.String || f.IsNil() {
		return nil
	}
	s, ok := f.Interface().(*string)
	if !ok {
		return nil
	}
	return s
}

// cacheAwareDataSourceReadByID attempts a cache-backed lookup of a data source object by
// ID, mirroring the common id/display_name/description/path semantics of
// policyDataSourceResourceFilterAndSet. On success it sets those schema attributes and
// returns the typed object so the caller can set any resource-specific fields. Returns
// ok=false when the cache is disabled, objID is empty, the ID was just written by this
// run (see TryCacheRead), or the object can't be found/converted — callers should fall
// through to the regular (uncached) read path in all of those cases.
func cacheAwareDataSourceReadByID[T any](d *schema.ResourceData, m interface{}, connector client.Connector, objID string, resourceType string, bindingType bindings.BindingType) (*T, bool) {
	if objID == "" || !IsCacheEnabled(m) {
		return nil, false
	}
	if shouldIndexByPath(resourceType) && !strings.Contains(objID, "/") {
		// The cache bucket for this type is populated/keyed by path (see shouldIndexByPath),
		// but objID here is a short user-supplied id, which is only unique within its own
		// VPC/project, not across the shared project-wide bucket. Fall through to a live,
		// correctly-scoped read instead of risking a cross-VPC id collision.
		return nil, false
	}
	if _, ok := postWriteByKey.LoadAndDelete(postWriteKey(resourceType, objID)); ok {
		// Mirrors TryCacheRead/CacheAwareResourceRead's one-shot post-write bypass: skip
		// a cache hit that could still be serving a bucket snapshot taken before this
		// exact ID was written in this run, and let the caller fall through to a live read.
		return nil, false
	}
	val, err := gcache.readCache(objID, resourceType, d, m, connector)
	if err != nil {
		return nil, false
	}
	typedVal, ok := convertCachedValue[T](val, resourceType, objID, bindingType, m)
	if !ok {
		return nil, false
	}

	id := objID
	if idField := reflectStringField(typedVal, "Id"); idField != nil {
		id = *idField
	}
	d.SetId(id)
	d.Set("id", id)
	d.Set("display_name", reflectStringField(typedVal, "DisplayName"))
	d.Set("description", reflectStringField(typedVal, "Description"))
	d.Set("path", reflectStringField(typedVal, "Path"))
	return typedVal, true
}

// TryCacheRead reads from cache only and never falls back to a backend GET.
func TryCacheRead[T any](d *schema.ResourceData, m interface{}, connector client.Connector, resourceID string, resourceType string, bindingType bindings.BindingType) (*T, bool, bool, error) {
	cacheUsed := false
	cacheAttempted := false
	if isRefreshPhase(d) && IsCacheEnabled(m) {
		if _, ok := postWriteByKey.LoadAndDelete(postWriteKey(resourceType, resourceID)); ok {
			// Bypass cache after a write for this resource instance (one-shot).
			return nil, cacheUsed, true, nil
		}
		cacheAttempted = true
		val, err := gcache.readCache(resourceID, resourceType, d, m, connector)
		if err == nil {
			if typedVal, ok := convertCachedValue[T](val, resourceType, resourceID, bindingType, m); ok {
				return typedVal, true, cacheAttempted, nil
			}
		}
	}
	return nil, cacheUsed, cacheAttempted, nil
}

func structValuesToModels[T any](list []*data.StructValue, bt bindings.BindingType) ([]T, error) {
	out := make([]T, 0, len(list))
	for _, obj := range list {
		dataValue, errors := sharedTypeConverter.ConvertToGolang(obj, bt)
		if len(errors) > 0 {
			var zero T
			return nil, fmt.Errorf("converting %T for cache: %w", zero, errors[0])
		}
		v, ok := dataValue.(T)
		if !ok {
			return nil, fmt.Errorf("converting for cache: unexpected type %T", dataValue)
		}
		out = append(out, v)
	}
	return out, nil
}

func modelsToStructValues[T any](models []T, bt bindings.BindingType) ([]*data.StructValue, error) {
	out := make([]*data.StructValue, 0, len(models))
	for i := range models {
		dataValue, errors := sharedTypeConverter.ConvertToVapi(models[i], bt)
		if len(errors) > 0 {
			var zero T
			return nil, fmt.Errorf("converting %T to struct value: %w", zero, errors[0])
		}
		sv, ok := dataValue.(*data.StructValue)
		if !ok {
			return nil, fmt.Errorf("converting to struct value: expected *data.StructValue, got %T", dataValue)
		}
		out = append(out, sv)
	}
	return out, nil
}

func structValuesToRules(list []*data.StructValue) []model.Rule {
	out := make([]model.Rule, 0, len(list))
	for _, obj := range list {
		dataValue, errors := sharedTypeConverter.ConvertToGolang(obj, model.RuleBindingType())
		if len(errors) > 0 {
			continue
		}
		rule, ok := dataValue.(model.Rule)
		if !ok {
			continue
		}
		out = append(out, rule)
	}
	return out
}

// groupRulesByValidParentPath indexes rules by ParentPath, restricted to known parent paths.
// Rules with a missing, empty, or unrecognized ParentPath are dropped.
func groupRulesByValidParentPath(validParentPaths map[string]struct{}, rules []model.Rule) map[string][]model.Rule {
	byParent := make(map[string][]model.Rule)
	for _, r := range rules {
		if r.ParentPath == nil {
			continue
		}
		pp := strings.TrimSpace(*r.ParentPath)
		if pp == "" {
			continue
		}
		if _, ok := validParentPaths[pp]; !ok {
			continue
		}
		byParent[pp] = append(byParent[pp], r)
	}
	// NSX evaluates rules in ascending sequence_number order; the search API result
	// order is not guaranteed to match, so sort each parent's rules explicitly.
	for pp, rules := range byParent {
		sort.SliceStable(rules, func(i, j int) bool {
			var si, sj int64
			if rules[i].SequenceNumber != nil {
				si = *rules[i].SequenceNumber
			}
			if rules[j].SequenceNumber != nil {
				sj = *rules[j].SequenceNumber
			}
			return si < sj
		})
		byParent[pp] = rules
	}
	return byParent
}

// attachRulesByParentPath assigns rules to the parent whose Path matches rule.ParentPath.
// Unmatched rules are dropped; parent order is preserved.
func attachRulesByParentPath[P any](parents []P, rules []model.Rule, getPath func(P) *string, setRules func(*P, []model.Rule)) []P {
	validPaths := make(map[string]struct{})
	for _, p := range parents {
		path := getPath(p)
		if path != nil {
			k := strings.TrimSpace(*path)
			if k != "" {
				validPaths[k] = struct{}{}
			}
		}
	}
	byParent := groupRulesByValidParentPath(validPaths, rules)

	out := make([]P, len(parents))
	for i := range parents {
		p := parents[i]
		key := ""
		if path := getPath(p); path != nil {
			key = strings.TrimSpace(*path)
		}
		bucket := byParent[key]
		setRules(&p, bucket)
		out[i] = p
	}
	return out
}

func mergeGatewayPolicyCacheSearchResults(parents, children []*data.StructValue) ([]*data.StructValue, error) {
	gp, err := structValuesToModels[model.GatewayPolicy](parents, model.GatewayPolicyBindingType())
	if err != nil {
		return nil, err
	}
	rules := structValuesToRules(children)
	merged := attachRulesByParentPath(gp, rules,
		func(p model.GatewayPolicy) *string { return p.Path },
		func(p *model.GatewayPolicy, r []model.Rule) { p.Rules = r },
	)
	return modelsToStructValues(merged, model.GatewayPolicyBindingType())
}

func mergeSecurityPolicyCacheSearchResults(parents, children []*data.StructValue) ([]*data.StructValue, error) {
	sp, err := structValuesToModels[model.SecurityPolicy](parents, model.SecurityPolicyBindingType())
	if err != nil {
		return nil, err
	}
	rules := structValuesToRules(children)
	merged := attachRulesByParentPath(sp, rules,
		func(p model.SecurityPolicy) *string { return p.Path },
		func(p *model.SecurityPolicy, r []model.Rule) { p.Rules = r },
	)
	return modelsToStructValues(merged, model.SecurityPolicyBindingType())
}

// CacheAwareResourceRead is the standard read path for a policy resource's Read function: try
// the shared search-backed cache first, and fall back to backendRead (a direct GET) on a cache
// miss, a conversion failure, cache being disabled, or the resource having just been written by
// this run (postWriteBypass — avoids reading a search snapshot that predates the write). resourceID
// must be the same key used by the corresponding MarkPostWriteAndInvalidateCacheForResourceType
// call for this resource (use CacheKeyForResourceID to guarantee that for path-indexed types).
// Returns (object, cacheUsed, cacheAttempted, error): callers generally only need the object and
// error; cacheUsed/cacheAttempted exist for tests and diagnostics.
func CacheAwareResourceRead[T any](d *schema.ResourceData, m interface{}, connector client.Connector, resourceID string, resourceType string, bindingType bindings.BindingType, backendRead func() (*T, error), patchFunc func(obj *T) error) (*T, bool, bool, error) {
	cacheAttempted := false
	postWriteBypass := false

	if isRefreshPhase(d) && IsCacheEnabled(m) {
		if _, ok := postWriteByKey.LoadAndDelete(postWriteKey(resourceType, resourceID)); ok {
			// Ensure read-your-writes semantics: bypass cache once right after a write.
			// Search-backed cache may be briefly stale immediately after a write.
			postWriteBypass = true
			cacheAttempted = true
		}
		if !postWriteBypass {
			cacheAttempted = true
			val, err := gcache.readCache(resourceID, resourceType, d, m, connector)
			if err == nil {
				if typedVal, ok := convertCachedValue[T](val, resourceType, resourceID, bindingType, m); ok {
					return typedVal, true, cacheAttempted, nil
				}
			}
		}
	}

	if postWriteBypass {
		log.Printf("[DEBUG] Cache post-write bypass: direct GET for resourceType=%s id=%s", resourceType, resourceID) //nolint:gosec
	} else {
		log.Printf("[DEBUG] Cache fallback: direct GET for resourceType=%s id=%s", resourceType, resourceID) //nolint:gosec
	}
	obj, err := backendRead()
	if err != nil {
		return nil, false, cacheAttempted, err
	}

	// Only stamp provider-managed tags in config_scope mode.
	if isConfigScopedCacheMode(m) {
		_, patchErr := ensureProviderManagedTagsWithPatchFunc(obj, m, patchFunc)
		if patchErr != nil {
			log.Printf("[WARNING] Failed to patch provider-managed tags for %s %s: %v", resourceType, resourceID, patchErr) //nolint:gosec
		}
	}
	// Strip provider-managed tags from state in non-global mode.
	if !isGlobalSearchCacheMode(m) {
		stripProviderManagedTagsFromAny(obj)
	}

	return obj, false, cacheAttempted, nil
}

// providerManagedTagsSearchFragments returns NSX search fragments for provider-managed tags.
func providerManagedTagsSearchFragments(runID string) []string {
	if runID == "" {
		return nil
	}
	return []string{
		fmt.Sprintf("tags.scope:%s", escapeSpecialCharacters(managedDefaultTagScope)),
		fmt.Sprintf("tags.tag:%s", escapeSpecialCharacters(runID)),
	}
}

// providerManagedTagsSearchQuery is the AND-joined search fragment for provider-managed tags.
func providerManagedTagsSearchQuery(runID string) string {
	fr := providerManagedTagsSearchFragments(runID)
	if len(fr) == 0 {
		return ""
	}
	return strings.Join(fr, " AND ")
}

// buildTagQuery builds the NSX search tag filter, adding provider-managed tags when needed.
func buildTagQuery(d *schema.ResourceData, runID string, m interface{}) string {
	managedTagPresent := false
	if d == nil {
		return ""
	}

	shouldAddProviderTags := !isGlobalSearchCacheMode(m) && runID != ""

	tags, exists := d.GetOk("tag")
	if !exists {
		if shouldAddProviderTags {
			return providerManagedTagsSearchQuery(runID)
		}
		return ""
	}

	tagSet, ok := tags.(*schema.Set)
	if !ok {
		if shouldAddProviderTags {
			return providerManagedTagsSearchQuery(runID)
		}
		return ""
	}
	if tagSet.Len() == 0 {
		if shouldAddProviderTags {
			return providerManagedTagsSearchQuery(runID)
		}
		return ""
	}

	var tagQueries []string
	for _, tagInterface := range tagSet.List() {
		tagMap := tagInterface.(map[string]interface{})

		scope, hasScope := tagMap["scope"]
		tag, hasTag := tagMap["tag"]
		if hasScope && scope != nil && scope.(string) != "" {
			rawScope := scope.(string)
			if isManagedDefaultTagScope(&rawScope) {
				managedTagPresent = true
			}
			tagQueries = append(tagQueries, fmt.Sprintf("tags.scope:%s", escapeSpecialCharacters(rawScope)))
		}
		if hasTag && tag != nil && tag.(string) != "" {
			tagQueries = append(tagQueries, fmt.Sprintf("tags.tag:%s", escapeSpecialCharacters(tag.(string))))
		}
	}
	if !isGlobalSearchCacheMode(m) && !managedTagPresent && runID != "" {
		tagQueries = append(tagQueries, providerManagedTagsSearchFragments(runID)...)
	}

	if len(tagQueries) == 0 {
		return ""
	}

	return strings.Join(tagQueries, " AND ")
}

// InvalidateCacheForResourceType clears the cache bucket for a resource type.
func InvalidateCacheForResourceType(resourceType string, m interface{}) {
	if !IsCacheEnabled(m) {
		return
	}
	tc := gcache.getTypeCache(resourceType)
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.data = make(map[string]map[string]*data.StructValue)
	log.Printf("[DEBUG] Cache invalidated for resourceType=%s", resourceType) //nolint:gosec
}

// stripProviderManagedTagsFromAny removes nsx-tf/tf-run-id tags from an NSX model object via reflection.
func stripProviderManagedTagsFromAny(obj interface{}) {
	if obj == nil {
		return
	}
	objValue := reflect.ValueOf(obj)
	for objValue.IsValid() && objValue.Kind() == reflect.Pointer {
		if objValue.IsNil() {
			return
		}
		objValue = objValue.Elem()
	}
	if !objValue.IsValid() || objValue.Kind() != reflect.Struct {
		return
	}
	tagsField := objValue.FieldByName("Tags")
	if !tagsField.IsValid() || !tagsField.CanSet() || tagsField.Kind() != reflect.Slice {
		return
	}
	current, ok := tagsField.Interface().([]model.Tag)
	if !ok {
		return
	}
	userTags := make([]model.Tag, 0, len(current))
	for _, tag := range current {
		if !isManagedDefaultTag(tag) {
			userTags = append(userTags, tag)
		}
	}
	tagsField.Set(reflect.ValueOf(userTags))
}

// ensureProviderManagedTagsWithPatchFunc adds provider-managed tags to obj via patchFunc when missing.
func ensureProviderManagedTagsWithPatchFunc[T any](obj T, m interface{}, patchFunc func(obj T) error) (interface{}, error) {
	objValue := reflect.ValueOf(obj)
	for objValue.IsValid() && objValue.Kind() == reflect.Pointer {
		if objValue.IsNil() {
			return nil, nil
		}
		objValue = objValue.Elem()
	}
	if !objValue.IsValid() || objValue.Kind() != reflect.Struct {
		return nil, nil
	}

	tagsField := objValue.FieldByName("Tags")
	if !tagsField.IsValid() {
		return nil, nil
	}

	var currentTags []model.Tag
	if tagsField.Kind() == reflect.Slice && !tagsField.IsNil() {
		for i := 0; i < tagsField.Len(); i++ {
			tag := tagsField.Index(i).Interface().(model.Tag)
			currentTags = append(currentTags, tag)
		}
	}

	runID := m.(nsxtClients).CommonConfig.contextID
	expectedManagedTags := getProviderManagedDefaultTags(runID)

	needsPatch := false
	for _, expectedTag := range expectedManagedTags {
		found := false
		for _, currentTag := range currentTags {
			if currentTag.Scope != nil && expectedTag.Scope != nil &&
				currentTag.Tag != nil && expectedTag.Tag != nil &&
				*currentTag.Scope == *expectedTag.Scope &&
				*currentTag.Tag == *expectedTag.Tag {
				found = true
				break
			}
		}
		if !found {
			needsPatch = true
			break
		}
	}

	if !needsPatch {
		return nil, nil
	}
	log.Printf("[DEBUG] Provider-managed tag missing; patching object to add scope=%s", managedDefaultTagScope) //nolint:gosec

	userTags := make([]model.Tag, 0)
	for _, tag := range currentTags {
		if !isManagedDefaultTag(tag) {
			userTags = append(userTags, tag)
		}
	}

	mergedTags := mergeManagedDefaultAndUserTags(expectedManagedTags, userTags)

	if tagsField.CanSet() {
		tagsField.Set(reflect.ValueOf(mergedTags))

		if patchFunc != nil {
			err := patchFunc(obj)
			if err != nil {
				return nil, fmt.Errorf("failed to patch tags to NSX API: %w", err)
			}
			log.Printf("[DEBUG] Patched provider-managed tags successfully")
			// Bump Revision to match the server-side PATCH increment and avoid stale-revision errors.
			revField := objValue.FieldByName("Revision")
			if revField.IsValid() && revField.Kind() == reflect.Pointer && !revField.IsNil() {
				revField.Elem().SetInt(revField.Elem().Int() + 1)
			}
		}

		return obj, nil
	}

	return nil, fmt.Errorf("Provider managed tags field is not settable")
}
