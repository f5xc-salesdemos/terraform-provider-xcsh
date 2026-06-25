// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/f5xc-salesdemos/terraform-provider-xcsh/internal/acctest"
)

func TestAccBotDefenseAppInfrastructureResource_basic(t *testing.T) {
	t.Skip("Skipping: requires Bot Defense license/subscription")
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-bdai")
	nsName := acctest.RandomName("tf-acc-test-ns")
	resourceName := "xcsh_bot_defense_app_infrastructure.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {Source: "hashicorp/time"},
		},
		CheckDestroy: acctest.CheckResourceDestroyed("xcsh_bot_defense_app_infrastructure"),
		Steps: []resource.TestStep{
			{
				Config: testAccBotDefenseAppInfrastructureConfig_basic(nsName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", nsName),
					resource.TestCheckResourceAttr(resourceName, "environment_type", "PRODUCTION"),
					resource.TestCheckResourceAttr(resourceName, "traffic_type", "WEB"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccBotDefenseAppInfrastructureImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccBotDefenseAppInfrastructureResource_emptyPlan(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-bot")
	nsName := acctest.RandomName("tf-acc-test-ns")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_bot_defense_app_infrastructure"),
		Steps: []resource.TestStep{
			{Config: testAccBotDefenseAppInfrastructureConfig_basic(nsName, rName)},
			{Config: testAccBotDefenseAppInfrastructureConfig_basic(nsName, rName), PlanOnly: true, ExpectNonEmptyPlan: false},
		},
	})
}

func testAccBotDefenseAppInfrastructureImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		namespace := rs.Primary.Attributes["namespace"]
		name := rs.Primary.Attributes["name"]
		return fmt.Sprintf("%s/%s", namespace, name), nil
	}
}

func testAccBotDefenseAppInfrastructureConfig_basic(nsName, name string) string {
	return acctest.ConfigCompose(
		acctest.ProviderConfig(),
		fmt.Sprintf(`
resource "xcsh_namespace" "test" {
  name = %[1]q
}

resource "time_sleep" "wait_for_namespace" {
  depends_on      = [xcsh_namespace.test]
  create_duration = "5s"
}

resource "xcsh_bot_defense_app_infrastructure" "test" {
  depends_on = [time_sleep.wait_for_namespace]
  name      = %[2]q
  namespace = xcsh_namespace.test.name

  environment_type = "PRODUCTION"
  traffic_type     = "WEB"

  cloud_hosted {
    infra_host_name = "test.example.com"
    region          = "US"
  }
}
`, nsName, name))
}
