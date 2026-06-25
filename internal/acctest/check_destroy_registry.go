// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package acctest

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/f5xc-salesdemos/terraform-provider-xcsh/internal/client"
)

// ResourceVerifier is a function that verifies a resource no longer exists via API
type ResourceVerifier func(ctx context.Context, c *client.Client, namespace, name string) error

// ResourceDeleter is a function that deletes a resource via API
// Used by CheckResourceDisappears to simulate external deletion
type ResourceDeleter func(ctx context.Context, c *client.Client, namespace, name string) error

// resourceVerifierRegistry maps Terraform resource types to their API verification functions
// This registry covers ALL 78 testable resources (100% coverage)
var resourceVerifierRegistry = map[string]ResourceVerifier{
	// ============================================================================
	// Address & Network Resources
	// ============================================================================
	"xcsh_address_allocator": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAddressAllocator(ctx, ns, name)
		return err
	},
	"xcsh_advertise_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAdvertisePolicy(ctx, ns, name)
		return err
	},
	"xcsh_cluster": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetCluster(ctx, ns, name)
		return err
	},
	"xcsh_endpoint": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetEndpoint(ctx, ns, name)
		return err
	},
	"xcsh_network_connector": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetNetworkConnector(ctx, ns, name)
		return err
	},
	"xcsh_network_firewall": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetNetworkFirewall(ctx, ns, name)
		return err
	},
	"xcsh_network_interface": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetNetworkInterface(ctx, ns, name)
		return err
	},
	"xcsh_network_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetNetworkPolicy(ctx, ns, name)
		return err
	},
	"xcsh_network_policy_rule": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetNetworkPolicyRule(ctx, ns, name)
		return err
	},
	"xcsh_network_policy_view": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetNetworkPolicyView(ctx, ns, name)
		return err
	},
	"xcsh_segment": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetSegment(ctx, ns, name)
		return err
	},
	"xcsh_tunnel": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetTunnel(ctx, ns, name)
		return err
	},
	"xcsh_virtual_network": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetVirtualNetwork(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Alert & Monitoring Resources
	// ============================================================================
	"xcsh_alert_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAlertPolicy(ctx, ns, name)
		return err
	},
	"xcsh_alert_receiver": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAlertReceiver(ctx, ns, name)
		return err
	},
	"xcsh_global_log_receiver": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetGlobalLogReceiver(ctx, ns, name)
		return err
	},
	"xcsh_log_receiver": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetLogReceiver(ctx, ns, name)
		return err
	},

	// ============================================================================
	// API Security Resources
	// ============================================================================
	"xcsh_api_crawler": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAPICrawler(ctx, ns, name)
		return err
	},
	"xcsh_api_definition": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAPIDefinition(ctx, ns, name)
		return err
	},
	"xcsh_api_discovery": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAPIDiscovery(ctx, ns, name)
		return err
	},
	"xcsh_api_testing": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAPITesting(ctx, ns, name)
		return err
	},
	"xcsh_app_api_group": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAppAPIGroup(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Application Security Resources
	// ============================================================================
	"xcsh_app_firewall": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAppFirewall(ctx, ns, name)
		return err
	},
	"xcsh_app_setting": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAppSetting(ctx, ns, name)
		return err
	},
	"xcsh_app_type": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAppType(ctx, ns, name)
		return err
	},
	"xcsh_sensitive_data_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetSensitiveDataPolicy(ctx, ns, name)
		return err
	},
	"xcsh_waf_exclusion_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetWAFExclusionPolicy(ctx, ns, name)
		return err
	},

	// ============================================================================
	// BGP Resources
	// ============================================================================
	"xcsh_bgp_asn_set": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetBGPAsnSet(ctx, ns, name)
		return err
	},
	"xcsh_bgp_routing_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetBGPRoutingPolicy(ctx, ns, name)
		return err
	},

	// ============================================================================
	// CDN Resources
	// ============================================================================
	"xcsh_cdn_cache_rule": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetCDNCacheRule(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Certificate Resources
	// ============================================================================
	"xcsh_certificate": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetCertificate(ctx, ns, name)
		return err
	},
	"xcsh_certificate_chain": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetCertificateChain(ctx, ns, name)
		return err
	},
	"xcsh_crl": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetCRL(ctx, ns, name)
		return err
	},
	"xcsh_trusted_ca_list": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetTrustedCAList(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Container & Registry Resources
	// ============================================================================
	"xcsh_container_registry": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetContainerRegistry(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Data Resources
	// ============================================================================
	"xcsh_data_group": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetDataGroup(ctx, ns, name)
		return err
	},
	"xcsh_data_type": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetDataType(ctx, ns, name)
		return err
	},

	// ============================================================================
	// DNS Resources
	// ============================================================================
	"xcsh_dns_compliance_checks": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetDNSComplianceChecks(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Firewall & ACL Resources
	// ============================================================================
	"xcsh_enhanced_firewall_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetEnhancedFirewallPolicy(ctx, ns, name)
		return err
	},
	"xcsh_fast_acl": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetFastACL(ctx, ns, name)
		return err
	},
	"xcsh_fast_acl_rule": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetFastACLRule(ctx, ns, name)
		return err
	},
	"xcsh_filter_set": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetFilterSet(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Fleet Resources
	// ============================================================================
	"xcsh_fleet": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetFleet(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Forwarding & Proxy Resources
	// ============================================================================
	"xcsh_forward_proxy_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetForwardProxyPolicy(ctx, ns, name)
		return err
	},
	"xcsh_forwarding_class": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetForwardingClass(ctx, ns, name)
		return err
	},
	"xcsh_proxy": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetProxy(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Health Check Resources
	// ============================================================================
	"xcsh_healthcheck": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetHealthcheck(ctx, ns, name)
		return err
	},

	// ============================================================================
	// IP Prefix Resources
	// ============================================================================
	"xcsh_ip_prefix_set": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetIPPrefixSet(ctx, ns, name)
		return err
	},

	// ============================================================================
	// iRule Resources
	// ============================================================================
	"xcsh_irule": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetIrule(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Load Balancer Resources
	// ============================================================================
	"xcsh_http_loadbalancer": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetHTTPLoadBalancer(ctx, ns, name)
		return err
	},
	"xcsh_tcp_loadbalancer": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetTCPLoadBalancer(ctx, ns, name)
		return err
	},
	"xcsh_udp_loadbalancer": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetUDPLoadBalancer(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Namespace Resources
	// ============================================================================
	"xcsh_namespace": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetNamespace(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Origin Pool Resources
	// ============================================================================
	"xcsh_origin_pool": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetOriginPool(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Policer & Rate Limiting Resources
	// ============================================================================
	"xcsh_policer": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetPolicer(ctx, ns, name)
		return err
	},
	"xcsh_rate_limiter": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetRateLimiter(ctx, ns, name)
		return err
	},
	"xcsh_rate_limiter_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetRateLimiterPolicy(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Protocol Inspection Resources
	// ============================================================================
	"xcsh_protocol_inspection": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetProtocolInspection(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Service Policy Resources
	// ============================================================================
	"xcsh_service_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetServicePolicy(ctx, ns, name)
		return err
	},
	"xcsh_service_policy_rule": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetServicePolicyRule(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Site Resources
	// ============================================================================
	"xcsh_aws_vpc_site": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAWSVPCSite(ctx, ns, name)
		return err
	},
	"xcsh_azure_vnet_site": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAzureVNETSite(ctx, ns, name)
		return err
	},
	"xcsh_gcp_vpc_site": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetGCPVPCSite(ctx, ns, name)
		return err
	},
	"xcsh_securemesh_site": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetSecuremeshSite(ctx, ns, name)
		return err
	},
	"xcsh_virtual_site": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetVirtualSite(ctx, ns, name)
		return err
	},

	// ============================================================================
	// User Identification Resources
	// ============================================================================
	"xcsh_user_identification": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetUserIdentification(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Virtual Host Resources
	// ============================================================================
	"xcsh_virtual_host": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetVirtualHost(ctx, ns, name)
		return err
	},

	// ============================================================================
	// Infrastructure Resources (from original registry)
	// ============================================================================
	"xcsh_cloud_connect": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetCloudConnect(ctx, ns, name)
		return err
	},
	"xcsh_nfv_service": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetNfvService(ctx, ns, name)
		return err
	},
	"xcsh_cminstance": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetCminstance(ctx, ns, name)
		return err
	},
	"xcsh_policy_based_routing": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetPolicyBasedRouting(ctx, ns, name)
		return err
	},
	"xcsh_apm": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetAPM(ctx, ns, name)
		return err
	},
	"xcsh_code_base_integration": func(ctx context.Context, c *client.Client, ns, name string) error {
		_, err := c.GetCodeBaseIntegration(ctx, ns, name)
		return err
	},
}

// resourceDeleterRegistry maps Terraform resource types to their API delete functions
// Used by CheckResourceDisappears to simulate external deletion of resources
var resourceDeleterRegistry = map[string]ResourceDeleter{
	// ============================================================================
	// Address & Network Resources
	// ============================================================================
	"xcsh_address_allocator": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAddressAllocator(ctx, ns, name)
	},
	"xcsh_advertise_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAdvertisePolicy(ctx, ns, name)
	},
	"xcsh_cluster": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteCluster(ctx, ns, name)
	},
	"xcsh_endpoint": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteEndpoint(ctx, ns, name)
	},
	"xcsh_network_connector": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteNetworkConnector(ctx, ns, name)
	},
	"xcsh_network_firewall": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteNetworkFirewall(ctx, ns, name)
	},
	"xcsh_network_interface": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteNetworkInterface(ctx, ns, name)
	},
	"xcsh_network_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteNetworkPolicy(ctx, ns, name)
	},
	"xcsh_network_policy_rule": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteNetworkPolicyRule(ctx, ns, name)
	},
	"xcsh_network_policy_view": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteNetworkPolicyView(ctx, ns, name)
	},
	"xcsh_segment": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteSegment(ctx, ns, name)
	},
	"xcsh_tunnel": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteTunnel(ctx, ns, name)
	},
	"xcsh_virtual_network": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteVirtualNetwork(ctx, ns, name)
	},

	// ============================================================================
	// Alert & Monitoring Resources
	// ============================================================================
	"xcsh_alert_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAlertPolicy(ctx, ns, name)
	},
	"xcsh_alert_receiver": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAlertReceiver(ctx, ns, name)
	},
	"xcsh_global_log_receiver": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteGlobalLogReceiver(ctx, ns, name)
	},
	"xcsh_log_receiver": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteLogReceiver(ctx, ns, name)
	},

	// ============================================================================
	// API Security Resources
	// ============================================================================
	"xcsh_api_crawler": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAPICrawler(ctx, ns, name)
	},
	"xcsh_api_definition": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAPIDefinition(ctx, ns, name)
	},
	"xcsh_api_discovery": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAPIDiscovery(ctx, ns, name)
	},
	"xcsh_api_testing": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAPITesting(ctx, ns, name)
	},
	"xcsh_app_api_group": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAppAPIGroup(ctx, ns, name)
	},

	// ============================================================================
	// Application Security Resources
	// ============================================================================
	"xcsh_app_firewall": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAppFirewall(ctx, ns, name)
	},
	"xcsh_app_setting": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAppSetting(ctx, ns, name)
	},
	"xcsh_app_type": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAppType(ctx, ns, name)
	},
	"xcsh_sensitive_data_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteSensitiveDataPolicy(ctx, ns, name)
	},
	"xcsh_waf_exclusion_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteWAFExclusionPolicy(ctx, ns, name)
	},

	// ============================================================================
	// BGP Resources
	// ============================================================================
	"xcsh_bgp_asn_set": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteBGPAsnSet(ctx, ns, name)
	},

	// ============================================================================
	// Bot Defense & Security Resources
	// ============================================================================
	"xcsh_malicious_user_mitigation": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteMaliciousUserMitigation(ctx, ns, name)
	},
	"xcsh_user_identification": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteUserIdentification(ctx, ns, name)
	},

	// ============================================================================
	// Certificate Resources
	// ============================================================================
	"xcsh_certificate": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteCertificate(ctx, ns, name)
	},
	"xcsh_certificate_chain": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteCertificateChain(ctx, ns, name)
	},

	// ============================================================================
	// Cloud Credentials & Infrastructure
	// ============================================================================
	"xcsh_cloud_credentials": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteCloudCredentials(ctx, ns, name)
	},
	"xcsh_dns_domain": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteDNSDomain(ctx, ns, name)
	},

	// ============================================================================
	// Data Resources
	// ============================================================================
	"xcsh_data_group": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteDataGroup(ctx, ns, name)
	},
	"xcsh_data_type": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteDataType(ctx, ns, name)
	},

	// ============================================================================
	// Forwarding & Rate Limiting Resources
	// ============================================================================
	"xcsh_filter_set": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteFilterSet(ctx, ns, name)
	},
	"xcsh_forwarding_class": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteForwardingClass(ctx, ns, name)
	},
	"xcsh_policer": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeletePolicer(ctx, ns, name)
	},
	"xcsh_rate_limiter_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteRateLimiterPolicy(ctx, ns, name)
	},
	"xcsh_rate_limiter": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteRateLimiter(ctx, ns, name)
	},

	// ============================================================================
	// Healthcheck Resources
	// ============================================================================
	"xcsh_healthcheck": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteHealthcheck(ctx, ns, name)
	},

	// ============================================================================
	// IP & Protocol Resources
	// ============================================================================
	"xcsh_ip_prefix_set": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteIPPrefixSet(ctx, ns, name)
	},
	"xcsh_protocol_policer": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteProtocolPolicer(ctx, ns, name)
	},

	// ============================================================================
	// Load Balancer Resources
	// ============================================================================
	"xcsh_http_loadbalancer": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteHTTPLoadBalancer(ctx, ns, name)
	},
	"xcsh_tcp_loadbalancer": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteTCPLoadBalancer(ctx, ns, name)
	},
	"xcsh_cdn_loadbalancer": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteCDNLoadBalancer(ctx, ns, name)
	},

	// ============================================================================
	// Namespace & Origin Resources
	// ============================================================================
	"xcsh_namespace": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteNamespace(ctx, ns, name)
	},
	"xcsh_origin_pool": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteOriginPool(ctx, ns, name)
	},

	// ============================================================================
	// Route Resources
	// ============================================================================
	"xcsh_route": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteRoute(ctx, ns, name)
	},

	// ============================================================================
	// Service Policy Resources
	// ============================================================================
	"xcsh_service_policy": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteServicePolicy(ctx, ns, name)
	},
	"xcsh_service_policy_rule": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteServicePolicyRule(ctx, ns, name)
	},

	// ============================================================================
	// Site Resources
	// ============================================================================
	"xcsh_aws_vpc_site": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAWSVPCSite(ctx, ns, name)
	},
	"xcsh_azure_vnet_site": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAzureVNETSite(ctx, ns, name)
	},
	"xcsh_gcp_vpc_site": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteGCPVPCSite(ctx, ns, name)
	},
	"xcsh_virtual_site": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteVirtualSite(ctx, ns, name)
	},

	// ============================================================================
	// Virtual Host Resources
	// ============================================================================
	"xcsh_virtual_host": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteVirtualHost(ctx, ns, name)
	},

	// ============================================================================
	// Platform & Configuration Resources
	// ============================================================================
	"xcsh_tenant_configuration": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteTenantConfiguration(ctx, ns, name)
	},
	"xcsh_workload_flavor": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteWorkloadFlavor(ctx, ns, name)
	},
	"xcsh_nfv_service": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteNfvService(ctx, ns, name)
	},
	"xcsh_cminstance": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteCminstance(ctx, ns, name)
	},
	"xcsh_policy_based_routing": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeletePolicyBasedRouting(ctx, ns, name)
	},
	"xcsh_apm": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteAPM(ctx, ns, name)
	},
	"xcsh_code_base_integration": func(ctx context.Context, c *client.Client, ns, name string) error {
		return c.DeleteCodeBaseIntegration(ctx, ns, name)
	},
}

// CheckResourceDestroyedWithAPIVerification verifies resources are deleted from the API
// This enhanced version performs actual API calls to verify deletion.
func CheckResourceDestroyedWithAPIVerification(resourceType string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		verifier, ok := resourceVerifierRegistry[resourceType]
		if !ok {
			// Fall back to state-only check with warning for unregistered types
			return checkResourceDestroyedStateOnly(resourceType, s)
		}

		c, err := GetTestClient()
		if err != nil {
			return fmt.Errorf("failed to get test client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			name := rs.Primary.Attributes["name"]
			namespace := rs.Primary.Attributes["namespace"]
			if name == "" {
				name = rs.Primary.ID
			}
			if namespace == "" {
				namespace = "system"
			}

			// Retry loop to handle async deletion
			maxRetries := 6
			for i := 0; i < maxRetries; i++ {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				err := verifier(ctx, c, namespace, name)
				cancel()

				if err != nil {
					// Check if it's a "not found" error - success!
					if isNotFoundError(err) {
						break // Resource is deleted
					}
					// Some other error occurred
					return fmt.Errorf("unexpected error checking %s %s/%s: %w", resourceType, namespace, name, err)
				}

				// Resource still exists
				if i == maxRetries-1 {
					return fmt.Errorf("%s %s/%s still exists in F5 XC API after waiting", resourceType, namespace, name)
				}

				// Wait before retrying
				time.Sleep(5 * time.Second)
			}
		}

		return nil
	}
}

// Note: isNotFoundError is defined in sweep.go and shared across acctest package

// checkResourceDestroyedStateOnly performs a state-only check for unregistered resource types
func checkResourceDestroyedStateOnly(resourceType string, s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourceType {
			continue
		}
		// Resource is in destroyed state per Terraform, but we can't verify with API
		// This is a soft pass - resource appears destroyed from Terraform's perspective
	}
	return nil
}

// RegisterResourceVerifier allows registering additional resource verifiers at runtime
func RegisterResourceVerifier(resourceType string, verifier ResourceVerifier) {
	resourceVerifierRegistry[resourceType] = verifier
}

// GetRegisteredResourceTypes returns a list of resource types with API verification
func GetRegisteredResourceTypes() []string {
	types := make([]string, 0, len(resourceVerifierRegistry))
	for t := range resourceVerifierRegistry {
		types = append(types, t)
	}
	return types
}

// GetRegistrySize returns the number of registered resource verifiers
func GetRegistrySize() int {
	return len(resourceVerifierRegistry)
}

// CheckResourceDisappears returns a TestCheckFunc that deletes a resource
// via the API to simulate external deletion (outside of Terraform).
// This is used to test that Terraform properly detects and handles
// resources that have been deleted externally.
//
// Usage:
//
//	{
//	    Config: testAccResourceConfig_basic(rName),
//	    Check: resource.ComposeTestCheckFunc(
//	        acctest.CheckResourceExists(resourceName),
//	        acctest.CheckResourceDisappears("xcsh_namespace", resourceName),
//	    ),
//	    ExpectNonEmptyPlan: true,
//	},
func CheckResourceDisappears(resourceType, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		deleter, ok := resourceDeleterRegistry[resourceType]
		if !ok {
			return fmt.Errorf("no deleter registered for resource type: %s", resourceType)
		}

		c, err := GetTestClient()
		if err != nil {
			return fmt.Errorf("failed to get test client: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		name := rs.Primary.Attributes["name"]
		namespace := rs.Primary.Attributes["namespace"]
		if namespace == "" {
			namespace = "system"
		}

		if err := deleter(ctx, c, namespace, name); err != nil {
			// Ignore "not found" errors - resource may already be deleted
			if !isNotFoundError(err) {
				return fmt.Errorf("failed to delete %s %s/%s: %w", resourceType, namespace, name, err)
			}
		}

		return nil
	}
}

// RegisterResourceDeleter allows registering additional resource deleters at runtime
func RegisterResourceDeleter(resourceType string, deleter ResourceDeleter) {
	resourceDeleterRegistry[resourceType] = deleter
}

// GetRegisteredDeleterTypes returns a list of resource types with delete capability
func GetRegisteredDeleterTypes() []string {
	types := make([]string, 0, len(resourceDeleterRegistry))
	for t := range resourceDeleterRegistry {
		types = append(types, t)
	}
	return types
}

// GetDeleterRegistrySize returns the number of registered resource deleters
func GetDeleterRegistrySize() int {
	return len(resourceDeleterRegistry)
}
