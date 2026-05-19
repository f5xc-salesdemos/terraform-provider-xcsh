// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package resource

// AICategories maps resource names to their functional categories for AI discovery.
// This is used by getResourceAIMetadata to generate structured metadata prefixes
// that help AI tools understand resource classification.
var AICategories = map[string]string{
	// Load Balancing
	"http_loadbalancer": "Load Balancing",
	"tcp_loadbalancer":  "Load Balancing",
	"udp_loadbalancer":  "Load Balancing",
	"dns_load_balancer": "Load Balancing",
	"cdn_loadbalancer":  "Load Balancing",
	"origin_pool":       "Load Balancing",
	"healthcheck":       "Load Balancing",
	"route":             "Load Balancing",

	// Security
	"app_firewall":                   "Security",
	"service_policy":                 "Security",
	"network_firewall":               "Security",
	"rate_limiter":                   "Security",
	"bot_defense_app_infrastructure": "Security",
	"malicious_user_mitigation":      "Security",
	"waf_exclusion_policy":           "Security",
	"enhanced_firewall_policy":       "Security",
	"forward_proxy_policy":           "Security",

	// Networking
	"network_connector": "Networking",
	"virtual_network":   "Networking",
	"cloud_connect":     "Networking",
	"cloud_link":        "Networking",
	"bgp":               "Networking",
	"ip_prefix_set":     "Networking",
	"network_interface": "Networking",
	"virtual_site":      "Networking",

	// Sites & Infrastructure
	"securemesh_site":    "Sites",
	"securemesh_site_v2": "Sites",
	"aws_vpc_site":       "Sites",
	"azure_vnet_site":    "Sites",
	"gcp_vpc_site":       "Sites",
	"aws_tgw_site":       "Sites",
	"voltstack_site":     "Sites",

	// DNS
	"dns_zone":              "DNS",
	"dns_domain":            "DNS",
	"dns_lb_pool":           "DNS",
	"dns_lb_health_check":   "DNS",
	"dns_compliance_checks": "DNS",

	// Kubernetes
	"k8s_cluster":              "Kubernetes",
	"virtual_k8s":              "Kubernetes",
	"k8s_cluster_role":         "Kubernetes",
	"k8s_cluster_role_binding": "Kubernetes",
	"k8s_pod_security_policy":  "Kubernetes",
	"container_registry":       "Kubernetes",

	// Authentication & Credentials
	"authentication":    "Authentication",
	"cloud_credentials": "Authentication",
	"api_credential":    "Authentication",
	"token":             "Authentication",
	"secret_policy":     "Authentication",

	// Certificates
	"certificate":       "Certificates",
	"certificate_chain": "Certificates",
	"trusted_ca_list":   "Certificates",

	// Monitoring
	"log_receiver":        "Monitoring",
	"global_log_receiver": "Monitoring",
	"alert_policy":        "Monitoring",
	"alert_receiver":      "Monitoring",

	// API Security
	"api_definition": "API Security",
	"api_discovery":  "API Security",
	"api_testing":    "API Security",
	"api_crawler":    "API Security",

	// Organization
	"namespace":      "Organization",
	"tenant":         "Organization",
	"role":           "Organization",
	"allowed_tenant": "Organization",
}

// Dependencies maps resource names to their common dependencies.
// This helps AI tools understand creation order.
var Dependencies = map[string][]string{
	"http_loadbalancer": {"namespace", "origin_pool"},
	"tcp_loadbalancer":  {"namespace", "origin_pool"},
	"udp_loadbalancer":  {"namespace", "origin_pool"},
	"origin_pool":       {"namespace", "healthcheck"},
	"healthcheck":       {"namespace"},
	"route":             {"namespace", "http_loadbalancer"},
	"app_firewall":      {"namespace"},
	"service_policy":    {"namespace"},
	"rate_limiter":      {"namespace"},
	"certificate":       {"namespace"},
	"api_definition":    {"namespace"},
	"dns_zone":          {},
	"virtual_site":      {},
	"namespace":         {},
}
