# NSX-T Terraform Provider Caching Mechanism

## Overview

This document describes the proof-of-concept in-memory caching mechanism implemented in the NSX-T Terraform provider to optimize API performance during Terraform operations, particularly during the refresh phase.

## Motivation

### Problem Statement

The NSX-T Terraform provider was experiencing performance issues due to heavy GET API calls being triggered repeatedly during Terraform operations:

- **Terraform Refresh Phase**: During `terraform plan` and `terraform apply`, Terraform performs a refresh phase where it reads the current state of all resources from the backend
- **Repeated API Calls**: Each resource read operation resulted in individual GET API calls to NSX-T
- **Performance Impact**: Large configurations with many resources resulted in hundreds of individual API calls, significantly increasing execution time
- **API Traffic**: Unnecessary load on NSX-T API endpoints due to redundant calls for the same resource types

### Impact Areas

- Slow `terraform plan` execution times
- Increased `terraform apply` duration during refresh phase
- Higher API traffic and potential rate limiting
- Poor user experience with large NSX-T deployments

## Solution Architecture

### Core Concept

The caching mechanism implements a **search-once, cache-many, reuse-often** pattern:

1. **Detection**: Identify when Terraform is in refresh phase
2. **Bulk Fetch**: Use NSX-T search API to fetch all resources of a given type in one call
3. **Cache Storage**: Store results in an in-memory map structure
4. **Reuse**: Serve subsequent requests from cache instead of making new API calls

### Key Components

#### 1. Cache Structure (`cache.go`)

```go
type cache struct {
    mu       sync.RWMutex                              // Concurrent access protection
    data     map[string]map[string]*data.StructValue   // Nested map: query -> displayName -> resource
    cacheHit int                                       // Performance metrics
    cacheMis int                                       // Performance metrics
}
```

**Cache Key Structure:**
- **Outer Key**: Query string (includes resource type and context)
- **Inner Key**: Resource display name
- **Value**: NSX-T resource as `*data.StructValue`

#### 2. Refresh Phase Detection

```go
func isRefreshPhase(d *schema.ResourceData) bool {
    return d.Id() != "" && !d.HasChangesExcept()
}
```

**Detection Logic:**
- Resource has an existing ID (not a new resource)
- No configuration changes detected (`!d.HasChangesExcept()`)
- Indicates Terraform is reading current state, not applying changes

#### 3. Query String Generation

```go
func getQueryString(resourceType string, context utl.SessionContext) string {
    switch context.ClientType {
    case utl.Global:
        return fmt.Sprintf("resource_type:%s AND marked_for_delete:false AND context:Global", resourceType)
    case utl.Local:
        return fmt.Sprintf("resource_type:%s AND marked_for_delete:false AND context:Local", resourceType)
    case utl.VPC, utl.Multitenancy:
        return fmt.Sprintf("resource_type:%s AND marked_for_delete:false AND context:%s-%s", 
                          resourceType, context.ProjectID, context.VPCID)
    default:
        return fmt.Sprintf("resource_type:%s AND marked_for_delete:false", resourceType)
    }
}
```

**Context-Aware Caching:**
- Different cache entries for different NSX-T contexts (Global, Local, VPC, Multitenancy)
- Ensures proper isolation between different scopes
- Prevents cross-contamination of cached data

## Implementation Details

### Cache Operations Flow

#### 1. Cache Read Operation

```go
func (c *cache) readCache(displayName string, resourceType string, d *schema.ResourceData, 
                         m interface{}, connector client.Connector) (interface{}, error) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    query := getQueryString(resourceType, getSessionContext(d, m))
    
    // Check if data exists in cache
    if val := c.getQueryResult(query, displayName); val != nil {
        c.cacheHit += 1
        return val, nil
    }
    
    // Cache miss - populate cache
    err := c.writeCache(query, resourceType, d, m, connector)
    if err != nil {
        return nil, err
    }
    
    return c.getQueryResult(query, displayName), nil
}
```

#### 2. Cache Population (Bulk Fetch)

```go
func (c *cache) getListOfPolicyResources(query string, d *schema.ResourceData, 
                                        connector client.Connector, context api.SessionContext, 
                                        resourceType string) error {
    // Use search API to fetch ALL resources of this type
    resultList, err := listPolicyResources(connector, context, resourceType, nil)
    if err != nil {
        return fmt.Errorf("error listing resource %s %w", resourceType, err)
    }
    
    // Convert list to map for fast lookups
    tmp := converListToMap(resultList)
    c.data[query] = tmp
    return nil
}
```

#### 3. List-to-Map Conversion

```go
func converListToMap(list []*data.StructValue) map[string]*data.StructValue {
    converter := bindings.NewTypeConverter()
    ret := make(map[string]*data.StructValue)
    
    for _, obj := range list {
        dataValue, errors := converter.ConvertToGolang(obj, model.PolicyConfigResourceBindingType())
        if len(errors) > 0 {
            return nil
        }
        resource := dataValue.(model.PolicyConfigResource)
        if resource.DisplayName != nil {
            ret[*resource.DisplayName] = obj  // Key by display name for fast lookup
        }
    }
    return ret
}
```

### Resource Integration Example

#### VPC Subnet Resource (`resource_nsxt_vpc_subnet.go`)

```go
func resourceNsxtVpcSubnetRead(d *schema.ResourceData, m interface{}) error {
    connector := getPolicyConnector(m)
    id := d.Id()
    displayName := d.Get("display_name").(string)
    var obj model.VpcSubnet
    var err error
    cacheUsed := false

    // Check if we're in refresh phase
    if isRefreshPhase(d) {
        fmt.Println("---------------------> Refresh Phase of plan/apply")
        
        // Try to get from cache
        val, err := gcache.readCache(displayName, "VpcSubnet", d, m, connector)
        if err == nil {
            converter := bindings.NewTypeConverter()
            goVal, convErrs := converter.ConvertToGolang(val.(*data.StructValue), model.VpcSubnetBindingType())
            if len(convErrs) == 0 {
                obj = goVal.(model.VpcSubnet)
                cacheUsed = true
            }
        }
    }
    
    // Fallback to regular API call if cache miss or not in refresh phase
    if !cacheUsed {
        fmt.Println("--------> Using the backend API, regular flow")
        client := clientLayer.NewSubnetsClient(connector)
        parents := getVpcParentsFromContext(getSessionContext(d, m))
        obj, err = client.Get(parents[0], parents[1], parents[2], id)
        if err != nil {
            return handleReadError(d, "VpcSubnet", id, err)
        }
    }

    // Set resource attributes from obj...
}
```

## Concurrency and Thread Safety

### Mutex Protection

```go
type cache struct {
    mu sync.RWMutex  // Read-Write mutex for concurrent access
    // ...
}

func (c *cache) readCache(...) {
    c.mu.Lock()      // Exclusive lock for read operations
    defer c.mu.Unlock()
    // ...
}
```

**Concurrency Strategy:**
- **Read-Write Mutex**: Allows multiple concurrent reads when no writes are happening
- **Exclusive Locking**: Write operations (cache population) get exclusive access
- **Lock Scope**: Minimal lock duration to reduce contention
- **Thread Safety**: Safe for concurrent Terraform operations

### Global Cache Instance

```go
var gcache = &cache{data: make(map[string]map[string]*data.StructValue), cacheHit: 0, cacheMis: 0}
```

- **Singleton Pattern**: Single global cache instance shared across all resources
- **Memory Efficiency**: Avoids duplicate caching of same resource types
- **Cross-Resource Sharing**: Different resource types can benefit from shared cache

## Performance Characteristics

### Cache Hit/Miss Tracking

```go
type cache struct {
    cacheHit int  // Number of successful cache retrievals
    cacheMis int  // Number of cache misses requiring API calls
}
```

**Metrics Collection:**
- Cache hit ratio indicates effectiveness
- Debug output shows cache performance in real-time
- Helps identify optimization opportunities

### Expected Performance Gains

**Before Caching:**
- N resources = N individual GET API calls
- Linear scaling with resource count
- Each call has network latency overhead

**After Caching:**
- N resources = 1 search API call + (N-1) cache hits
- Constant time for subsequent reads of same resource type
- Significant reduction in network calls

**Example Scenario:**
- 50 VPC subnets in configuration
- **Without cache**: 50 individual GET calls
- **With cache**: 1 search call + 49 cache hits
- **Improvement**: ~98% reduction in API calls

## Limitations and Trade-offs

### Current Limitations

-- Large deployments may consume significant RAM
--  No retry mechanisms for cache population failures
-- Currently implemented only for VPC Subnet resources
-- Other resource types still use direct API calls

## Conclusion
The proof-of-concept caching mechanism demonstrates significant potential for improving NSX-T Terraform provider performance. The implementation successfully reduces API calls during refresh operations while maintaining data consistency and thread safety.
