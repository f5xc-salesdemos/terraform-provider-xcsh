// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/f5xc-salesdemos/terraform-provider-xcsh/internal/acctest"
)

func TestAccNetworkPolicyRuleDataSource_basic(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test")
	resourceName := "xcsh_network_policy_rule.test"
	dataSourceName := "data.xcsh_network_policy_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkPolicyRuleDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespace", resourceName, "namespace"),
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
				),
			},
		},
	})
}

func testAccNetworkPolicyRuleDataSourceConfig_basic(name string) string {
	// Network policy rule must be in system namespace, protocol must be uppercase
	return acctest.ConfigCompose(
		acctest.ProviderConfig(),
		fmt.Sprintf(`
resource "xcsh_network_policy_rule" "test" {
  name      = %[1]q
  namespace = "system"

  action   = "ALLOW"
  protocol = "TCP"
  ports    = ["443", "8080"]

  prefix {
    prefix = ["192.168.1.0/24"]
  }
}

data "xcsh_network_policy_rule" "test" {
  name      = xcsh_network_policy_rule.test.name
  namespace = xcsh_network_policy_rule.test.namespace
}
`, name))
}
