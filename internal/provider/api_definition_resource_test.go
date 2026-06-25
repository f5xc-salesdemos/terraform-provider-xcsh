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

func TestAccAPIDefinitionResource_basic(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_api_definition.test"
	rName := acctest.RandomName("tf-test-def")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_api_definition"),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIDefinitionConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", "shared"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccAPIDefinitionImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccAPIDefinitionImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
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

func TestAccAPIDefinitionResource_withLabels(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_api_definition.test"
	rName := acctest.RandomName("tf-test-def")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_api_definition"),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIDefinitionConfig_withLabels(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "test"),
				),
			},
		},
	})
}

func TestAccAPIDefinitionResource_fullLifecycle(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_api_definition.test"
	rName := acctest.RandomName("tf-test-def")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_api_definition"),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIDefinitionConfig_withLabels(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "test"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccAPIDefinitionImportStateIdFunc(resourceName),
			},
			{
				Config: testAccAPIDefinitionConfig_basic(rName),
				Check:  acctest.CheckResourceExists(resourceName),
			},
			{
				Config:             testAccAPIDefinitionConfig_basic(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAPIDefinitionResource_emptyPlan(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-test-def")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_api_definition"),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIDefinitionConfig_basic(rName),
			},
			{
				Config:             testAccAPIDefinitionConfig_basic(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAPIDefinitionResource_planChecks(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_api_definition.test"
	rName := acctest.RandomName("tf-test-def")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_api_definition"),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIDefinitionConfig_basic(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					},
				},
			},
			{
				Config: testAccAPIDefinitionConfig_withLabels(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: testAccAPIDefinitionConfig_withLabels(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

func TestAccAPIDefinitionResource_invalidName(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccAPIDefinitionConfig_basic("Invalid-NAME"),
				ExpectError: regexp.MustCompile(`(?i)(invalid|name|must)`),
			},
		},
	})
}

func TestAccAPIDefinitionResource_requiresReplace(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_api_definition.test"
	rName1 := acctest.RandomName("tf-test-def")
	rName2 := acctest.RandomName("tf-test-def")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_api_definition"),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIDefinitionConfig_basic(rName1),
				Check:  acctest.CheckResourceExists(resourceName),
			},
			{
				Config: testAccAPIDefinitionConfig_basic(rName2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}

func TestAccAPIDefinitionResource_strictSchemaOrigin(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_api_definition.test"
	rName := acctest.RandomName("tf-test-def")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_api_definition"),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIDefinitionConfig_strictSchema(rName),
				Check:  acctest.CheckResourceExists(resourceName),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccAPIDefinitionImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccAPIDefinitionConfig_withLabels(name string) string {
	return fmt.Sprintf(`
resource "xcsh_api_definition" "test" {
  name      = %[1]q
  namespace = "shared"

  labels = {
    environment = "test"
  }
}
`, name)
}

func testAccAPIDefinitionConfig_strictSchema(name string) string {
	return fmt.Sprintf(`
resource "xcsh_api_definition" "test" {
  name      = %[1]q
  namespace = "shared"

  strict_schema_origin {}
}
`, name)
}

func testAccAPIDefinitionConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "xcsh_api_definition" "test" {
  name      = %[1]q
  namespace = "shared"
}
`, name)
}
