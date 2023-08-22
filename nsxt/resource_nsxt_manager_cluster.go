/* Copyright © 2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

func resourceNsxtManagerCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtManagerClusterCreate,
		Read:   resourceNsxtManagerClusterRead,
		Update: resourceNsxtManagerClusterUpdate,
		Delete: resourceNsxtManagerClusterDelete,

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"node": {
				Type:        schema.TypeList,
				Description: "Nodes in the cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "ID of the cluster node",
							Computed:    true,
						},
						"fqdn": {
							Type:        schema.TypeString,
							Description: "FQDN for the cluster node",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the cluster node",
							Computed:    true,
						},
						"ip_address": {
							Type:         schema.TypeString,
							Description:  "IP address of the cluster node that will join the cluster",
							ValidateFunc: validation.IsIPAddress,
							Required:     true,
						},
						"password": {
							Type:        schema.TypeString,
							Description: "The password for login",
							Required:    true,
							Sensitive:   true,
						},
						"username": {
							Type:        schema.TypeString,
							Description: "The username for login",
							Required:    true,
						},
					},
				},
				Required: true,
			},
		},
	}
}

type NsxClusterNode struct {
	ID        string
	IPAddress string
	UserName  string
	Password  string
	Fqdn      string
	Status    string
}

func getClusterNodesFromSchema(d *schema.ResourceData) []NsxClusterNode {
	nodes := d.Get("node").([]interface{})
	var clusterNodes []NsxClusterNode
	for _, node := range nodes {
		data := node.(map[string]interface{})
		id := data["id"].(string)
		ipAddress := data["ip_address"].(string)
		userName := data["username"].(string)
		password := data["password"].(string)
		nodeObj := NsxClusterNode{
			ID:        id,
			IPAddress: ipAddress,
			UserName:  userName,
			Password:  password,
		}
		clusterNodes = append(clusterNodes, nodeObj)
	}
	return clusterNodes
}

func setNodesInSchema(d *schema.ResourceData, nodes []NsxClusterNode) error {
	// Retrieve node credential from schema and set them in the nodeList element
	// This is because the nodes are obtained using cluster client, and does not
	// contain any information about credential
	nodeIPToCredentialMap := getIPtoCredentialMap(d)

	var nodeList []map[string]interface{}
	for _, node := range nodes {
		_, nodeInSchema := nodeIPToCredentialMap[node.IPAddress]
		elem := make(map[string]interface{})
		elem["id"] = node.ID
		elem["ip_address"] = node.IPAddress
		// If node is not in schema but in cluster node list,
		// it means the node needs to be removed.
		// In this case we don't need its credential because
		// remove api does not require node credential
		if nodeInSchema {
			elem["username"] = nodeIPToCredentialMap[node.IPAddress][0]
			elem["password"] = nodeIPToCredentialMap[node.IPAddress][1]
		}
		elem["fqdn"] = node.Fqdn
		elem["status"] = node.Status
		nodeList = append(nodeList, elem)
	}
	return d.Set("node", nodeList)
}

func getIPtoCredentialMap(d *schema.ResourceData) map[string][]string {
	// returns a map, whose key is IP address of a node,
	// value is a slice of [node_login_username, node_login_password]
	nodes := d.Get("node").([]interface{})
	var idToCredentialMap = map[string][]string{}
	for _, node := range nodes {
		data := node.(map[string]interface{})
		ipAddress := data["ip_address"].(string)
		userName := data["username"].(string)
		password := data["password"].(string)
		cred := []string{userName, password}
		idToCredentialMap[ipAddress] = cred
	}
	return idToCredentialMap
}

func resourceNsxtManagerClusterCreate(d *schema.ResourceData, m interface{}) error {
	// Call Joincluster function on nodes that are not in the cluster
	nodes := getClusterNodesFromSchema(d)
	if len(nodes) == 0 {
		return fmt.Errorf("At least a manager appliance must be provided to form a cluster")
	}
	clusterID, certSha256Thumbprint, hostIP, err := getClusterInfoFromHostNode(d, m)
	if err != nil {
		return err
	}

	for _, guestNode := range nodes {
		err := joinNodeToCluster(clusterID, certSha256Thumbprint, guestNode, hostIP, d, m)
		if err != nil {
			return fmt.Errorf("Failed to join node %s: %s", guestNode.ID, err)
		}
	}
	d.SetId(clusterID)
	return resourceNsxtManagerClusterRead(d, m)
}

func getClusterInfoFromHostNode(d *schema.ResourceData, m interface{}) (string, string, string, error) {
	// function return values are:
	// clusterID, certSha256Thumbprint, hostIP, error
	connector := getPolicyConnector(m)
	client := nsx.NewClusterClient(connector)
	c := m.(nsxtClients)
	min := c.CommonConfig.MinRetryInterval
	max := c.CommonConfig.MaxRetryInterval
	maxRetries := c.CommonConfig.MaxRetries
	hostIP := ""
	for i := 0; i < maxRetries; i++ {
		clusterConfig, err := client.Get()
		if err != nil {
			return "", "", "", handleReadError(d, "Cluster Config", "", err)
		}
		if hostIP == "" {
			hostIP, err = getHostIPFromClusterConfig(m, clusterConfig)
			if err != nil {
				return "", "", "", err
			}
		}
		clusterID := *clusterConfig.ClusterId
		nodes := clusterConfig.Nodes
		node := nodes[0]
		apiListenAddr := node.ApiListenAddr
		if apiListenAddr != nil {
			certSha256Thumbprint := *apiListenAddr.CertificateSha256Thumbprint
			return clusterID, certSha256Thumbprint, hostIP, nil
		}
		interval := (rand.Intn(max-min) + min)
		time.Sleep(time.Duration(interval) * time.Millisecond)
		log.Printf("[DEBUG]: Waited %d ms before retrying getting API Listen Address, attempt %d", interval, i+1)
	}
	return "", "", "", fmt.Errorf("Failed to read ClusterConfig info from host node after %d attempts", maxRetries)
}

func getHostIPFromClusterConfig(client interface{}, clusterConfig nsxModel.ClusterConfig) (string, error) {
	c := client.(nsxtClients)
	host := c.Host
	host = strings.TrimPrefix(host, "https://")
	// Check if host is ip address or fqdn, if it's ip address then we are all set,
	// if it's fqdn, we loop through clusterconfig, find the node with same fqdn and
	// return its ip address
	ip := net.ParseIP(host)
	isIPAddress := ip != nil
	if isIPAddress {
		return host, nil
	}
	nodes := clusterConfig.Nodes
	for _, node := range nodes {
		fqdn := *node.Fqdn
		if fqdn == host {
			if node.ApiListenAddr != nil {
				// return ipv4 address if it has one, if it doesn't return ipv6 address
				addr := node.ApiListenAddr
				if *addr.IpAddress != "" {
					return *addr.IpAddress, nil
				}
				return *addr.Ipv6Address, nil
			}
			v4, v6 := getHTTPSIPFromNodeEntity(node.Entities)
			if v4 != "" {
				return v4, nil
			}
			return v6, nil
		}
	}
	return "", fmt.Errorf("Failed to get host ip from cluster config")
}

func getHTTPSIPFromNodeEntity(entities []nsxModel.NodeEntityInfo) (string, string) {
	// returns the ipv4 and ipv6 address of a node's https port 443
	for _, entity := range entities {
		if *entity.Port == 443 {
			return *entity.IpAddress, *entity.Ipv6Address
		}
	}
	return "", ""
}

func getNewNsxtClient(node NsxClusterNode, d *schema.ResourceData, clients interface{}) (interface{}, error) {
	// The API of joining a new node to cluster needs to be called with the new node's ip,
	// to do that we initialize a new nsxtClient with new node's ip as its host
	username := node.UserName
	password := node.Password
	host := node.IPAddress
	if !strings.HasPrefix(host, "https://") {
		host = fmt.Sprintf("https://%s", host)
	}
	c := clients.(nsxtClients)
	commonConfig := c.CommonConfig
	newClients := nsxtClients{
		CommonConfig: commonConfig,
	}
	err := configureNewClient(&newClients, &c, host, username, password)
	if err != nil {
		return nil, err
	}
	return newClients, nil
}

func configureNewClient(newClient *nsxtClients, oldClient *nsxtClients, host string, username string, password string) error {
	newClient.Host = host
	securityCtx, err := getConfiguredSecurityContext(newClient, "", "", "", username, password)
	if err != nil {
		return fmt.Errorf("Failed to configure new client with host %s: %s", host, err)
	}
	newClient.PolicySecurityContext = securityCtx
	newClient.PolicyHTTPClient = oldClient.PolicyHTTPClient
	err = initNSXVersion(getPolicyConnector(*newClient))
	if err != nil {
		return fmt.Errorf("Failed to configure new client with host %s: %s", host, err)
	}
	return nil
}

func joinNodeToCluster(clusterID string, certSha256Thumbprint string, guestNode NsxClusterNode, masterNodeIP string, d *schema.ResourceData, m interface{}) error {
	c, err := getNewNsxtClient(guestNode, d, m)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Cluster %s. Joining node %s", clusterID, guestNode.IPAddress)
	newNsxClients := c.(nsxtClients)
	connector := getPolicyConnector(newNsxClients)
	client := nsx.NewClusterClient(connector)
	username, password := getHostCredential(m)
	joinClusterParams := nsxModel.JoinClusterParameters{
		CertificateSha256Thumbprint: &certSha256Thumbprint,
		ClusterId:                   &clusterID,
		IpAddress:                   &masterNodeIP,
		Username:                    &username,
		Password:                    &password,
	}
	_, err = client.Joincluster(joinClusterParams)
	if err != nil {
		return fmt.Errorf("Failed to join node to cluster: %s, node ip address: %s", err, guestNode.IPAddress)
	}
	log.Printf("[INFO] Cluster %s. Completed join node %s", clusterID, guestNode.IPAddress)
	return nil
}

func getHostCredential(m interface{}) (string, string) {
	nsxtClient := m.(nsxtClients)
	username := nsxtClient.CommonConfig.Username
	password := nsxtClient.CommonConfig.Password
	return username, password
}

func resourceNsxtManagerClusterRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := nsx.NewClusterClient(connector)
	clusterConfig, err := client.Get()
	if err != nil {
		return fmt.Errorf("Failed to read Nsxt Manager Cluster: %s", err)
	}
	nodeInfo := clusterConfig.Nodes
	var nodes []NsxClusterNode
	hostIP, err := getHostIPFromClusterConfig(m, clusterConfig)
	isIPv4 := (net.ParseIP(hostIP)).To4() != nil
	if err != nil {
		return err
	}
	for _, node := range nodeInfo {
		ip := getIPFromNodeInfo(node, isIPv4)
		fqdn := *node.Fqdn
		status := *node.Status
		if ip != hostIP {
			id := *node.NodeUuid
			clusterNode := NsxClusterNode{
				ID:        id,
				IPAddress: ip,
				Fqdn:      fqdn,
				Status:    status,
			}
			nodes = append(nodes, clusterNode)
		}

	}
	d.Set("revision", clusterConfig.Revision)
	setNodesInSchema(d, nodes)
	return nil
}

func getIPFromNodeInfo(node nsxModel.ClusterNodeInfo, isIPv4 bool) string {
	// After join node into cluster, apiListen address may take some time to become available
	// if it's not available, we retrieve node ip from nodeEntityInfo
	var ip *string
	nodeEntities := node.Entities
	for _, entity := range nodeEntities {
		if *entity.Port == 443 {
			if isIPv4 {
				ip = entity.IpAddress
			} else {
				ip = entity.Ipv6Address
			}
			break
		}
	}
	return *ip
}

func resourceNsxtManagerClusterUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := nsx.NewClusterClient(connector)

	clusterID, certSha256Thumbprint, hostIP, err := getClusterInfoFromHostNode(d, m)
	if err != nil {
		return err
	}
	oldNodes, newNodes := d.GetChange("node")
	oldNodesIPs := getClusterNodesIPs(oldNodes)
	newNodesIPs := getClusterNodesIPs(newNodes)
	for _, node := range oldNodes.([]interface{}) {
		nodeMap := node.(map[string]interface{})
		ip := nodeMap["ip_address"].(string)
		if !slices.Contains(newNodesIPs, ip) {
			id := nodeMap["id"].(string)
			force := "true"
			gracefulShutdown := "true"
			ignoreRepositoryIPCheckParam := "false"
			_, err := client.Removenode(id, &force, &gracefulShutdown, &ignoreRepositoryIPCheckParam)
			if err != nil {
				return err
			}
		}
	}
	for _, node := range newNodes.([]interface{}) {
		nodeMap := node.(map[string]interface{})
		ip := nodeMap["ip_address"].(string)
		if !slices.Contains(oldNodesIPs, ip) {
			userName := nodeMap["username"].(string)
			password := nodeMap["password"].(string)
			ip := nodeMap["ip_address"].(string)
			nodeObj := NsxClusterNode{
				IPAddress: ip,
				UserName:  userName,
				Password:  password,
			}
			err = joinNodeToCluster(clusterID, certSha256Thumbprint, nodeObj, hostIP, d, m)
			if err != nil {
				return err
			}
		}
	}

	return resourceNsxtManagerClusterRead(d, m)
}

func getClusterNodesIPs(nodes interface{}) []string {
	var ips []string
	for _, node := range nodes.([]interface{}) {
		nodeMap := node.(map[string]interface{})
		ips = append(ips, nodeMap["ip_address"].(string))
	}
	return ips
}

func resourceNsxtManagerClusterDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := nsx.NewClusterClient(connector)
	nodes := getClusterNodesFromSchema(d)
	force := "true"
	gracefulShutdown := "true"
	ignoreRepositoryIPCheckParam := "false"
	for _, node := range nodes {
		guestNodeID := node.ID
		_, err := client.Removenode(guestNodeID, &force, &gracefulShutdown, &ignoreRepositoryIPCheckParam)
		if err != nil {
			return fmt.Errorf("Failed to remove node %s from cluster: %s", guestNodeID, err)
		}
	}
	return nil
}
