// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package provider_test

import (
	"fmt"
	"testing"

	"regexp"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/f5xc-salesdemos/terraform-provider-xcsh/internal/acctest"
)

// =============================================================================
// TEST: Basic origin_pool creation with public_name origin
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccOriginPoolResource_basic(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_origin_pool.test"
	rName := acctest.RandomName("tf-test-op")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckOriginPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOriginPoolConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckOriginPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", "system"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "port", "443"),
				),
			},
			// Import test
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccOriginPoolImportStateIdFunc(resourceName),
			},
		},
	})
}

// =============================================================================
// TEST: Origin pool with labels and description
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccOriginPoolResource_withLabels(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_origin_pool.test"
	rName := acctest.RandomName("tf-test-op")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckOriginPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOriginPoolConfig_withLabelsSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckOriginPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Test origin pool"),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "labels.team", "platform"),
				),
			},
		},
	})
}

// =============================================================================
// TEST: Update labels
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccOriginPoolResource_updateLabels(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_origin_pool.test"
	rName := acctest.RandomName("tf-test-op")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckOriginPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOriginPoolConfig_labelsUpdateSystem(rName, "dev"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckOriginPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "dev"),
				),
			},
			{
				Config: testAccOriginPoolConfig_labelsUpdateSystem(rName, "prod"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckOriginPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "prod"),
				),
			},
		},
	})
}

// =============================================================================
// TEST: Resource disappears (deleted outside Terraform)
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccOriginPoolResource_disappears(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_origin_pool.test"
	rName := acctest.RandomName("tf-test-op")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckOriginPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOriginPoolConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckOriginPoolExists(resourceName),
					acctest.CheckOriginPoolDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// =============================================================================
// TEST: Empty plan after apply (no drift)
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccOriginPoolResource_emptyPlan(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_origin_pool.test"
	rName := acctest.RandomName("tf-test-op")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckOriginPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOriginPoolConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckOriginPoolExists(resourceName),
				),
			},
			{
				Config:             testAccOriginPoolConfig_basicSystem(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// =============================================================================
// HELPER: Import state ID function
// =============================================================================
func testAccOriginPoolImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
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

// =============================================================================
// TEST: Plan checks (create, update, noop)
// =============================================================================
func TestAccOriginPoolResource_planChecks(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_origin_pool.test"
	rName := acctest.RandomName("tf-test-op")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckOriginPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOriginPoolConfig_basicSystem(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					},
				},
			},
			{
				Config: testAccOriginPoolConfig_labelsUpdateSystem(rName, "staging"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: testAccOriginPoolConfig_labelsUpdateSystem(rName, "staging"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// =============================================================================
// TEST: Known values plan check
// =============================================================================
func TestAccOriginPoolResource_knownValues(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_origin_pool.test"
	rName := acctest.RandomName("tf-test-op")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckOriginPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOriginPoolConfig_basicSystem(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rName)),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("namespace"), knownvalue.StringExact("system")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("port"), knownvalue.Int64Exact(443)),
					},
				},
			},
		},
	})
}

// =============================================================================
// TEST: Name validation
// =============================================================================
func TestAccOriginPoolResource_invalidName(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccOriginPoolConfig_basicSystem("Invalid-NAME"),
				ExpectError: regexp.MustCompile(`(?i)(invalid|name|must)`),
			},
		},
	})
}

// =============================================================================
// TEST: Name change requires replacement
// =============================================================================
func TestAccOriginPoolResource_requiresReplace(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_origin_pool.test"
	rName1 := acctest.RandomName("tf-test-op")
	rName2 := acctest.RandomName("tf-test-op")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckOriginPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOriginPoolConfig_basicSystem(rName1),
				Check:  acctest.CheckOriginPoolExists(resourceName),
			},
			{
				Config: testAccOriginPoolConfig_basicSystem(rName2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}

// =============================================================================
// TEST: Origin with public_ip instead of public_name
// =============================================================================
func TestAccOriginPoolResource_publicIp(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_origin_pool.test"
	rName := acctest.RandomName("tf-test-op")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckOriginPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOriginPoolConfig_publicIpSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckOriginPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "port", "8080"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccOriginPoolImportStateIdFunc(resourceName),
			},
		},
	})
}

// =============================================================================
// TEST: Multiple origin servers
// =============================================================================
func TestAccOriginPoolResource_multipleOrigins(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_origin_pool.test"
	rName := acctest.RandomName("tf-test-op")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckOriginPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOriginPoolConfig_multipleOriginsSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckOriginPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "origin_servers.#", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccOriginPoolImportStateIdFunc(resourceName),
			},
		},
	})
}

// =============================================================================
// TEST: Origin pool with healthcheck reference
// =============================================================================
func TestAccOriginPoolResource_withHealthcheckRef(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_origin_pool.test"
	rName := acctest.RandomName("tf-test-op")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			if err := acctest.CheckOriginPoolDestroyed(s); err != nil {
				return err
			}
			return acctest.CheckHealthcheckDestroyed(s)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccOriginPoolConfig_withHealthcheckRefSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckOriginPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "healthcheck.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccOriginPoolImportStateIdFunc(resourceName),
			},
		},
	})
}

// =============================================================================
// TEST: Port update
// =============================================================================
func TestAccOriginPoolResource_updatePort(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_origin_pool.test"
	rName := acctest.RandomName("tf-test-op")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckOriginPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOriginPoolConfig_portSystem(rName, 443),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckOriginPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "port", "443"),
				),
			},
			{
				Config: testAccOriginPoolConfig_portSystem(rName, 8443),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckOriginPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "port", "8443"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})
}

// =============================================================================
// TEST: Full lifecycle
// =============================================================================
func TestAccOriginPoolResource_fullLifecycle(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_origin_pool.test"
	rName := acctest.RandomName("tf-test-op")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckOriginPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOriginPoolConfig_withLabelsSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckOriginPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "Test origin pool"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "disable", "description"},
				ImportStateIdFunc:       testAccOriginPoolImportStateIdFunc(resourceName),
			},
			{
				Config: testAccOriginPoolConfig_portSystem(rName, 8080),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckOriginPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "port", "8080"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: testAccOriginPoolConfig_portSystem(rName, 8080),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccOriginPoolConfig_basicSystem(rName),
				Check:  acctest.CheckOriginPoolExists(resourceName),
			},
		},
	})
}

// =============================================================================
// CONFIG HELPERS - Use "system" namespace
// =============================================================================

func testAccOriginPoolConfig_basicSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_origin_pool" "test" {
  name       = %[1]q
  namespace  = "system"

  port = 443

  origin_servers {
    labels {}  # API returns this even if not set
    public_name {
      dns_name = "example.com"
    }
  }

  no_tls {}
  same_as_endpoint_port {}
}
`, name)
}

func testAccOriginPoolConfig_withLabelsSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_origin_pool" "test" {
  name        = %[1]q
  namespace   = "system"
  description = "Test origin pool"

  port = 443

  labels = {
    environment = "test"
    team        = "platform"
  }

  origin_servers {
    labels {}  # API returns this even if not set
    public_name {
      dns_name = "example.com"
    }
  }

  no_tls {}
  same_as_endpoint_port {}
}
`, name)
}

func testAccOriginPoolConfig_publicIpSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_origin_pool" "test" {
  name      = %[1]q
  namespace = "system"

  port = 8080

  origin_servers {
    labels {}
    public_ip {
      ip = "93.184.216.34"
    }
  }

  no_tls {}
  same_as_endpoint_port {}
}
`, name)
}

func testAccOriginPoolConfig_multipleOriginsSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_origin_pool" "test" {
  name      = %[1]q
  namespace = "system"

  port = 443

  origin_servers {
    labels {}
    public_name {
      dns_name = "backend1.example.com"
    }
  }

  origin_servers {
    labels {}
    public_name {
      dns_name = "backend2.example.com"
    }
  }

  no_tls {}
  same_as_endpoint_port {}
}
`, name)
}

func testAccOriginPoolConfig_withHealthcheckRefSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 3
  unhealthy_threshold = 1
  timeout             = 3
  interval            = 15

  tcp_health_check {}
}

resource "xcsh_origin_pool" "test" {
  name      = %[1]q
  namespace = "system"

  port = 443

  origin_servers {
    labels {}
    public_name {
      dns_name = "example.com"
    }
  }

  healthcheck {
    name      = xcsh_healthcheck.test.name
    namespace = xcsh_healthcheck.test.namespace
  }

  no_tls {}
  same_as_endpoint_port {}
}
`, name)
}

func testAccOriginPoolConfig_portSystem(name string, port int) string {
	return fmt.Sprintf(`
resource "xcsh_origin_pool" "test" {
  name      = %[1]q
  namespace = "system"

  port = %[2]d

  origin_servers {
    labels {}
    public_name {
      dns_name = "example.com"
    }
  }

  no_tls {}
  same_as_endpoint_port {}
}
`, name, port)
}

func testAccOriginPoolConfig_labelsUpdateSystem(name, env string) string {
	return fmt.Sprintf(`
resource "xcsh_origin_pool" "test" {
  name       = %[1]q
  namespace  = "system"

  port = 443

  labels = {
    environment = %[2]q
  }

  origin_servers {
    labels {}  # API returns this even if not set
    public_name {
      dns_name = "example.com"
    }
  }

  no_tls {}
  same_as_endpoint_port {}
}
`, name, env)
}

// =============================================================================
// NEGATIVE: Conflicting TLS options (OneOf violation)
// =============================================================================
func TestAccOriginPoolResource_conflictTlsOptions(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-test-op")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "xcsh_origin_pool" "test" {
  name      = %[1]q
  namespace = "system"
  port      = 443

  origin_servers {
    labels {}
    public_name {
      dns_name = "example.com"
    }
  }

  no_tls {}
  use_tls {
    tls_config {
      default_security {}
    }
    no_mtls {}
    volterra_trusted_ca {}
  }
  same_as_endpoint_port {}
}
`, rName),
				ExpectError: regexp.MustCompile(`(?i)(conflict|mutually exclusive|only one|Invalid|these attributes cannot)`),
			},
		},
	})
}
