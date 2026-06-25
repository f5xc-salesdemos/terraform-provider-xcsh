// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/f5xc-salesdemos/terraform-provider-xcsh/internal/acctest"
)

func TestAccAPIDiscoveryResource_basic(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_api_discovery.test"
	rName := acctest.RandomName("tf-test-disc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_api_discovery"),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIDiscoveryConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", "system"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccAPIDiscoveryImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccAPIDiscoveryImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
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

func TestAccAPIDiscoveryResource_withLabels(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_api_discovery.test"
	rName := acctest.RandomName("tf-test-disc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_api_discovery"),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIDiscoveryConfig_withLabels(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "test"),
				),
			},
		},
	})
}

func TestAccAPIDiscoveryResource_emptyPlan(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-test-disc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_api_discovery"),
		Steps: []resource.TestStep{
			{Config: testAccAPIDiscoveryConfig_basic(rName)},
			{Config: testAccAPIDiscoveryConfig_basic(rName), PlanOnly: true, ExpectNonEmptyPlan: false},
		},
	})
}

func TestAccAPIDiscoveryResource_planChecks(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_api_discovery.test"
	rName := acctest.RandomName("tf-test-disc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_api_discovery"),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIDiscoveryConfig_basic(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate)},
				},
			},
			{
				Config: testAccAPIDiscoveryConfig_withLabels(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate)},
				},
			},
			{
				Config: testAccAPIDiscoveryConfig_withLabels(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop)},
				},
			},
		},
	})
}

func TestAccAPIDiscoveryResource_invalidName(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: testAccAPIDiscoveryConfig_basic("Invalid-NAME"), ExpectError: regexp.MustCompile(`(?i)(invalid|name|must)`)},
		},
	})
}

func TestAccAPIDiscoveryResource_requiresReplace(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_api_discovery.test"
	rName1 := acctest.RandomName("tf-test-disc")
	rName2 := acctest.RandomName("tf-test-disc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_api_discovery"),
		Steps: []resource.TestStep{
			{Config: testAccAPIDiscoveryConfig_basic(rName1), Check: acctest.CheckResourceExists(resourceName)},
			{
				Config: testAccAPIDiscoveryConfig_basic(rName2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate)},
				},
			},
		},
	})
}

func testAccAPIDiscoveryConfig_withLabels(name string) string {
	return fmt.Sprintf(`
resource "xcsh_api_discovery" "test" {
  name      = %[1]q
  namespace = "system"
  labels = {
    environment = "test"
  }
}
`, name)
}

func testAccAPIDiscoveryConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "xcsh_api_discovery" "test" {
  name      = %[1]q
  namespace = "system"
}
`, name)
}
