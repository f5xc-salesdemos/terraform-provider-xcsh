// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/f5xc-salesdemos/terraform-provider-xcsh/internal/acctest"
)

func TestAccTcpLoadbalancerDataSource_basic(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-tcp-lb")
	resourceName := "xcsh_tcp_loadbalancer.test"
	dataSourceName := "data.xcsh_tcp_loadbalancer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTcpLoadbalancerDataSourceConfig_basic("", rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespace", resourceName, "namespace"),
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
				),
			},
		},
	})
}

func testAccTcpLoadbalancerDataSourceConfig_basic(nsName, name string) string {
	return acctest.ConfigCompose(
		acctest.ProviderConfig(),
		fmt.Sprintf(`
# Origin pool is required for TCP load balancer - it needs a backend cluster
resource "xcsh_origin_pool" "test" {
  name       = "%[2]s-pool"
  namespace  = "system"

  port = 443

  origin_servers {
    labels {}
    public_name {
      dns_name = "example.com"
    }
  }

  no_tls {}
  same_as_endpoint_port {}
}

resource "xcsh_tcp_loadbalancer" "test" {
  name       = %[2]q
  namespace  = "system"

  labels = {
    environment = "test"
    managed_by  = "terraform-acceptance-test"
  }

  # Domain and SNI are required for TCP on public shared VIP
  domains = ["%[2]s.example.com"]

  listen_port = 443

  # Required: Specify protocol type (tcp, tls_tcp, or tls_tcp_auto_cert)
  tcp {}

  # Required: SNI for TCP on public shared VIP
  sni {}

  # Required: TCP LB needs origin pools for routing
  origin_pools_weights {
    pool {
      name      = xcsh_origin_pool.test.name
      namespace = "system"
    }
    weight = 1
  }

  # Required: Specify advertise configuration
  advertise_on_public_default_vip {}
}

data "xcsh_tcp_loadbalancer" "test" {
  depends_on = [xcsh_tcp_loadbalancer.test]
  name       = xcsh_tcp_loadbalancer.test.name
  namespace  = xcsh_tcp_loadbalancer.test.namespace
}
`, nsName, name))
}
