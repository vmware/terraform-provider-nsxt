package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"strings"
)

var defaultDomain = "default"
var defaultSite = "default"
var securityPolicyCategoryValues = []string{"Ethernet", "Emergency", "Infrastructure", "Environment", "Application"}
var securityPolicyDirectionValues = []string{model.Rule_DIRECTION_IN, model.Rule_DIRECTION_OUT, model.Rule_DIRECTION_IN_OUT}
var securityPolicyIPProtocolValues = []string{model.Rule_IP_PROTOCOL_IPV4, model.Rule_IP_PROTOCOL_IPV6, model.Rule_IP_PROTOCOL_IPV4_IPV6}
var securityPolicyActionValues = []string{model.Rule_ACTION_ALLOW, model.Rule_ACTION_DROP, model.Rule_ACTION_REJECT}
var gatewayPolicyCategoryValues = []string{"Emergency", "SystemRules", "SharedPreRules", "LocalGatewayRules", "AutoServiceRules", "Default"}
var policyFailOverModeValues = []string{model.Tier1_FAILOVER_MODE_PREEMPTIVE, model.Tier1_FAILOVER_MODE_NON_PREEMPTIVE}
var failOverModeDefaultPolicyT0Value = model.Tier0_FAILOVER_MODE_NON_PREEMPTIVE
var defaultPolicyLocaleServiceID = "default"

func getNsxIDSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "NSX ID for this resource",
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
	}
}

func getPathSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "Policy path for this resource",
		Computed:    true,
	}
}

func getDisplayNameSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "Display name for this resource",
		Required:    true,
	}
}

func getOptionalDisplayNameSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "Display name for this resource",
		Optional:    true,
	}
}

func getDescriptionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "Description for this resource",
		Optional:    true,
	}
}

func getDataSourceStringSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: description,
		Optional:    true,
		Computed:    true,
	}
}

func getDomainNameSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "The domain name to use for resources. If not specified 'default' is used",
		Optional:    true,
		Default:     defaultDomain,
	}
}

func getFailoverModeSchema(defaultValue string) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "Failover mode",
		Default:      defaultValue,
		Optional:     true,
		ValidateFunc: validation.StringInSlice(policyFailOverModeValues, false),
	}
}

func getIPv6NDRAPathSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "The path of an IPv6 NDRA profile",
		Optional:    true,
		Computed:    true,
	}
}

func getIPv6DadPathSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "The path of an IPv6 DAD profile",
		Optional:    true,
		Computed:    true,
	}
}

func getPolicyEdgeClusterPathSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "The path of the edge cluster connected to this gateway",
		Optional:     true,
		ValidateFunc: validatePolicyPath(),
	}
}

func getPolicyGatewayPathSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "The NSX-T Policy path to the Tier0 or Tier1 Gateway for this resource",
		Required:     true,
		ValidateFunc: validatePolicyPath(),
		ForceNew:     true,
	}
}

func getSecurityPolicyAndGatewayRulesSchema(scopeRequired bool) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of rules in the section",
		Optional:    true,
		MaxItems:    1000,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name": getDisplayNameSchema(),
				"description":  getDescriptionSchema(),
				"revision":     getRevisionSchema(),
				"sequence_number": {
					Type:        schema.TypeInt,
					Description: "Sequence number of the this rule",
					Computed:    true,
				},
				"destination_groups": {
					Type:        schema.TypeSet,
					Description: "List of destination groups",
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validatePolicyPath(),
					},
					Optional: true,
				},
				"destinations_excluded": {
					Type:        schema.TypeBool,
					Description: "Negation of destination groups",
					Optional:    true,
					Default:     false,
				},
				"direction": {
					Type:         schema.TypeString,
					Description:  "Traffic direction",
					Optional:     true,
					ValidateFunc: validation.StringInSlice(securityPolicyDirectionValues, false),
					Default:      model.Rule_DIRECTION_IN_OUT,
				},
				"disabled": {
					Type:        schema.TypeBool,
					Description: "Flag to disable the rule",
					Optional:    true,
					Default:     false,
				},
				"ip_version": {
					Type:         schema.TypeString,
					Description:  "IP version",
					Optional:     true,
					ValidateFunc: validation.StringInSlice(securityPolicyIPProtocolValues, false),
					Default:      model.Rule_IP_PROTOCOL_IPV4_IPV6,
				},
				"logged": {
					Type:        schema.TypeBool,
					Description: "Flag to enable packet logging",
					Optional:    true,
					Default:     false,
				},
				"notes": {
					Type:        schema.TypeString,
					Description: "Text for additional notes on changes",
					Optional:    true,
				},
				"profiles": {
					Type:        schema.TypeSet,
					Description: "List of profiles",
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validatePolicyPath(),
					},
					Optional: true,
				},
				"rule_id": {
					Type:        schema.TypeInt,
					Description: "Unique positive number that is assigned by the system and is useful for debugging",
					Computed:    true,
				},
				"scope": {
					Type:        schema.TypeSet,
					Description: "List of policy paths where the rule is applied",
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validatePolicyPath(),
					},
					Optional: !scopeRequired,
					Required: scopeRequired,
				},
				"services": {
					Type:        schema.TypeSet,
					Description: "List of services to match",
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validatePolicyPath(),
					},
					Optional: true,
				},
				"source_groups": {
					Type:        schema.TypeSet,
					Description: "List of source groups",
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validatePolicyPath(),
					},
					Optional: true,
				},
				"sources_excluded": {
					Type:        schema.TypeBool,
					Description: "Negation of source groups",
					Optional:    true,
					Default:     false,
				},
				"tag": getTagsSchema(),
				"log_label": {
					Type:        schema.TypeString,
					Description: "Additional information (string) which will be propagated to the rule syslog",
					Optional:    true,
				},
				"action": {
					Type:         schema.TypeString,
					Description:  "Action",
					Optional:     true,
					ValidateFunc: validation.StringInSlice(securityPolicyActionValues, false),
					Default:      model.Rule_ACTION_ALLOW,
				},
			},
		},
	}
}

func getPolicyGatewayPolicySchema() map[string]*schema.Schema {
	secPolicy := getPolicySecurityPolicySchema()
	// GW Policies don't support scope
	delete(secPolicy, "scope")
	secPolicy["category"].ValidateFunc = validation.StringInSlice(gatewayPolicyCategoryValues, false)
	// GW Policy rules require scope to be set
	secPolicy["rule"] = getSecurityPolicyAndGatewayRulesSchema(true)
	return secPolicy
}

func getPolicySecurityPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"nsx_id":       getNsxIDSchema(),
		"path":         getPathSchema(),
		"display_name": getDisplayNameSchema(),
		"description":  getDescriptionSchema(),
		"revision":     getRevisionSchema(),
		"tag":          getTagsSchema(),
		"domain":       getDomainNameSchema(),
		"category": {
			Type:         schema.TypeString,
			Description:  "Category",
			ValidateFunc: validation.StringInSlice(securityPolicyCategoryValues, false),
			Required:     true,
			ForceNew:     true,
		},
		"comments": {
			Type:        schema.TypeString,
			Description: "Comments for security policy lock/unlock",
			Optional:    true,
		},
		"locked": {
			Type:        schema.TypeBool,
			Description: "Indicates whether a security policy should be locked. If locked by a user, no other user would be able to modify this policy",
			Optional:    true,
			Default:     false,
		},
		"scope": {
			Type:        schema.TypeSet,
			Description: "The list of group paths where the rules in this policy will get applied",
			Optional:    true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validatePolicyPath(),
			},
		},
		// TODO - verify this is relevant not only across domains
		"sequence_number": {
			Type:        schema.TypeInt,
			Description: "This field is used to resolve conflicts between security policies across domains",
			Optional:    true,
			Default:     0,
		},
		"stateful": {
			Type:        schema.TypeBool,
			Description: "When it is stateful, the state of the network connects are tracked and a stateful packet inspection is performed",
			Optional:    true,
			Default:     true,
		},
		"tcp_strict": {
			Type:        schema.TypeBool,
			Description: "Ensures that a 3 way TCP handshake is done before the data packets are sent",
			Optional:    true,
			Computed:    true,
		},
		"rule": getSecurityPolicyAndGatewayRulesSchema(false),
	}
}

func setPolicyRulesInSchema(d *schema.ResourceData, rules []model.Rule) error {
	var rulesList []map[string]interface{}
	for _, rule := range rules {
		elem := make(map[string]interface{})
		elem["display_name"] = rule.DisplayName
		elem["description"] = rule.Description
		elem["notes"] = rule.Notes
		elem["logged"] = rule.Logged
		elem["log_label"] = rule.Tag
		elem["action"] = rule.Action
		elem["destinations_excluded"] = rule.DestinationsExcluded
		elem["sources_excluded"] = rule.SourcesExcluded
		elem["ip_version"] = rule.IpProtocol
		elem["direction"] = rule.Direction
		elem["disabled"] = rule.Disabled
		elem["revision"] = rule.Revision
		setPathListInMap(elem, "source_groups", rule.SourceGroups)
		setPathListInMap(elem, "destination_groups", rule.DestinationGroups)
		setPathListInMap(elem, "profiles", rule.Profiles)
		setPathListInMap(elem, "services", rule.Services)
		setPathListInMap(elem, "scope", rule.Scope)
		elem["sequence_number"] = rule.SequenceNumber

		var tagList []map[string]string
		for _, tag := range rule.Tags {
			tags := make(map[string]string)
			tags["scope"] = *tag.Scope
			tags["tag"] = *tag.Tag
			tagList = append(tagList, tags)
		}
		elem["tag"] = tagList

		rulesList = append(rulesList, elem)
	}

	return d.Set("rule", rulesList)
}

func getPolicyRulesFromSchema(d *schema.ResourceData) []model.Rule {
	rules := d.Get("rule").([]interface{})
	var ruleList []model.Rule
	seq := 0
	for _, rule := range rules {
		data := rule.(map[string]interface{})
		displayName := data["display_name"].(string)
		description := data["description"].(string)
		action := data["action"].(string)
		logged := data["logged"].(bool)
		tag := data["log_label"].(string)
		disabled := data["disabled"].(bool)
		sourcesExcluded := data["sources_excluded"].(bool)
		destinationsExcluded := data["destinations_excluded"].(bool)
		ipProtocol := data["ip_version"].(string)
		direction := data["direction"].(string)
		notes := data["notes"].(string)
		sequenceNumber := int64(seq)

		var tagStructs []model.Tag
		if data["tag"] != nil {
			tags := data["tag"].(*schema.Set).List()
			for _, tag := range tags {
				data := tag.(map[string]interface{})
				tagScope := data["scope"].(string)
				tagTag := data["tag"].(string)
				elem := model.Tag{
					Scope: &tagScope,
					Tag:   &tagTag}

				tagStructs = append(tagStructs, elem)
			}
		}

		// Use a different random Id each time, otherwise Update requires revision
		// to be set for existing rules, and NOT be set for new rules
		id := newUUID()

		elem := model.Rule{
			Id:                   &id,
			DisplayName:          &displayName,
			Notes:                &notes,
			Description:          &description,
			Action:               &action,
			Logged:               &logged,
			Tag:                  &tag,
			Tags:                 tagStructs,
			Disabled:             &disabled,
			SourcesExcluded:      &sourcesExcluded,
			DestinationsExcluded: &destinationsExcluded,
			IpProtocol:           &ipProtocol,
			Direction:            &direction,
			SourceGroups:         getPathListFromMap(data, "source_groups"),
			DestinationGroups:    getPathListFromMap(data, "destination_groups"),
			Services:             getPathListFromMap(data, "services"),
			Scope:                getPathListFromMap(data, "scope"),
			Profiles:             getPathListFromMap(data, "profiles"),
			SequenceNumber:       &sequenceNumber,
		}
		ruleList = append(ruleList, elem)
		seq = seq + 1
	}

	return ruleList
}

func getDataSourceDisplayNameSchema() *schema.Schema {
	return getDataSourceStringSchema("Display name of this resource")
}

func getDataSourceDescriptionSchema() *schema.Schema {
	return getDataSourceStringSchema("Description for this resource")
}

func getDataSourceIDSchema() *schema.Schema {
	return getDataSourceStringSchema("Unique ID of this resource")
}

func getIpv6ProfilePathsFromSchema(d *schema.ResourceData) []string {
	var profiles []string
	if d.Get("ipv6_ndra_profile_path") != "" {
		profiles = append(profiles, d.Get("ipv6_ndra_profile_path").(string))
	}
	if d.Get("ipv6_dad_profile_path") != "" {
		profiles = append(profiles, d.Get("ipv6_dad_profile_path").(string))
	}
	return profiles
}

func setIpv6ProfilePathsInSchema(d *schema.ResourceData, paths []string) error {
	for _, path := range paths {
		if strings.HasPrefix(path, "/infra/ipv6-ndra-profiles") {
			d.Set("ipv6_ndra_profile_path", path)
		}
		if strings.HasPrefix(path, "/infra/ipv6-dad-profiles") {
			d.Set("ipv6_dad_profile_path", path)
		}
	}
	return nil
}

func parseGatewayPolicyPath(gwPath string) (bool, string) {
	// sample path looks like "/infra/tier-0s/mytier0gw"
	isT0 := true
	segs := strings.Split(gwPath, "/")
	if segs[len(segs)-2] != "tier-0s" {
		isT0 = false
	}
	return isT0, segs[len(segs)-1]
}

func getPolicyPathSchema(isRequired bool, forceNew bool, description string) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  description,
		Optional:     !isRequired,
		Required:     isRequired,
		ForceNew:     forceNew,
		ValidateFunc: validatePolicyPath(),
	}
}

func getElemPolicyPathSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validatePolicyPath(),
	}
}

func getGatewayInterfaceSubnetsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of IP addresses and network prefixes for this interface",
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateIPCidr(),
		},
		Required: true,
	}
}

func getMtuSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Description:  "Maximum transmission unit specifies the size of the largest packet that a network protocol can transmit",
		ValidateFunc: validation.IntAtLeast(64),
	}
}

func getAllocationRangeListSchema(required bool, description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: description,
		Required:    required,
		Optional:    !required,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"start": {
					Type:         schema.TypeString,
					Description:  "The start IP Address for the range",
					Required:     true,
					ValidateFunc: validateSingleIP(),
				},
				"end": {
					Type:         schema.TypeString,
					Description:  "The end IP Address for the range",
					Required:     true,
					ValidateFunc: validateSingleIP(),
				},
			},
		},
	}
}

var gatewayInterfaceUrpfModeValues = []string{
	model.Tier0Interface_URPF_MODE_NONE,
	model.Tier0Interface_URPF_MODE_STRICT,
}

func getGatewayInterfaceUrpfModeSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "Unicast Reverse Path Forwarding mode",
		ValidateFunc: validation.StringInSlice(gatewayInterfaceUrpfModeValues, false),
		Default:      model.Tier0Interface_URPF_MODE_STRICT,
	}
}
