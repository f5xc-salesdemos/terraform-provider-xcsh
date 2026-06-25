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

func TestAccAdvertisePolicyResource_basic(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_advertise_policy.test"
	rName := acctest.RandomName("tf-test-adv")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_advertise_policy"),
		Steps: []resource.TestStep{
			{
				Config: testAccAdvertisePolicyConfig_basic(rName),
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
				ImportStateIdFunc:       testAccAdvertisePolicyImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccAdvertisePolicyImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
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

func TestAccAdvertisePolicyResource_withLabels(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_advertise_policy.test"
	rName := acctest.RandomName("tf-test-adv")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_advertise_policy"),
		Steps: []resource.TestStep{
			{
				Config: testAccAdvertisePolicyConfig_withLabels(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "test"),
				),
			},
		},
	})
}

func TestAccAdvertisePolicyResource_emptyPlan(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-test-adv")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_advertise_policy"),
		Steps: []resource.TestStep{
			{Config: testAccAdvertisePolicyConfig_basic(rName)},
			{Config: testAccAdvertisePolicyConfig_basic(rName), PlanOnly: true, ExpectNonEmptyPlan: false},
		},
	})
}

func TestAccAdvertisePolicyResource_planChecks(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_advertise_policy.test"
	rName := acctest.RandomName("tf-test-adv")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_advertise_policy"),
		Steps: []resource.TestStep{
			{
				Config: testAccAdvertisePolicyConfig_basic(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate)},
				},
			},
			{
				Config: testAccAdvertisePolicyConfig_withLabels(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate)},
				},
			},
			{
				Config: testAccAdvertisePolicyConfig_withLabels(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop)},
				},
			},
		},
	})
}

func TestAccAdvertisePolicyResource_invalidName(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: testAccAdvertisePolicyConfig_basic("Invalid-NAME"), ExpectError: regexp.MustCompile(`(?i)(invalid|name|must)`)},
		},
	})
}

func TestAccAdvertisePolicyResource_requiresReplace(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_advertise_policy.test"
	rName1 := acctest.RandomName("tf-test-adv")
	rName2 := acctest.RandomName("tf-test-adv")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_advertise_policy"),
		Steps: []resource.TestStep{
			{Config: testAccAdvertisePolicyConfig_basic(rName1), Check: acctest.CheckResourceExists(resourceName)},
			{
				Config: testAccAdvertisePolicyConfig_basic(rName2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate)},
				},
			},
		},
	})
}

func TestAccAdvertisePolicyResource_updatePort(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_advertise_policy.test"
	rName := acctest.RandomName("tf-test-adv")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_advertise_policy"),
		Steps: []resource.TestStep{
			{
				Config: testAccAdvertisePolicyConfig_port(rName, 80),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
				),
			},
			{
				Config: testAccAdvertisePolicyConfig_port(rName, 443),
				Check:  resource.TestCheckResourceAttr(resourceName, "port", "443"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate)},
				},
			},
		},
	})
}

func testAccAdvertisePolicyConfig_withLabels(name string) string {
	return fmt.Sprintf(`
resource "xcsh_advertise_policy" "test" {
  name            = %[1]q
  namespace       = "system"
  address         = "0.0.0.0"
  port            = 80
  protocol        = "TCP"
  skip_xff_append = false
  labels = {
    environment = "test"
  }
}
`, name)
}

func testAccAdvertisePolicyConfig_port(name string, port int) string {
	return fmt.Sprintf(`
resource "xcsh_advertise_policy" "test" {
  name            = %[1]q
  namespace       = "system"
  address         = "0.0.0.0"
  port            = %[2]d
  protocol        = "TCP"
  skip_xff_append = false
}
`, name, port)
}

func testAccAdvertisePolicyConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "xcsh_advertise_policy" "test" {
  name            = %[1]q
  namespace       = "system"
  address         = "0.0.0.0"
  port            = 80
  protocol        = "TCP"
  skip_xff_append = false
}
`, name)
}
