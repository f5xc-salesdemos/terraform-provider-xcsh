// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/f5xc-salesdemos/terraform-provider-xcsh/internal/acctest"
)

// =============================================================================
// TEST: Basic http_loadbalancer creation
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccHTTPLoadBalancerResource_basic(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-test-lb")
	resourceName := "xcsh_http_loadbalancer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccHTTPLoadBalancerConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", "system"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccHTTPLoadBalancerImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"timeouts",
					"http.dns_volterra_managed",
					"l7_ddos_protection",
				},
			},
		},
	})
}

// =============================================================================
// TEST: HTTP loadbalancer with labels and description
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccHTTPLoadBalancerResource_withLabels(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-test-lb")
	resourceName := "xcsh_http_loadbalancer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLoadBalancerConfig_withLabelsSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "labels.team", "platform"),
				),
			},
		},
	})
}

// =============================================================================
// TEST: HTTP loadbalancer with multiple domains
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccHTTPLoadBalancerResource_withDomains(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-test-lb-domain")
	resourceName := "xcsh_http_loadbalancer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLoadBalancerConfig_withDomainsSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "domains.#", "2"),
				),
			},
		},
	})
}

// =============================================================================
// TEST: Update labels
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccHTTPLoadBalancerResource_updateLabels(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-test-lb")
	resourceName := "xcsh_http_loadbalancer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLoadBalancerConfig_labelsUpdateSystem(rName, "dev"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "dev"),
				),
			},
			{
				Config: testAccHTTPLoadBalancerConfig_labelsUpdateSystem(rName, "prod"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "prod"),
				),
			},
		},
	})
}

// =============================================================================
// TEST: Empty plan after apply (no drift)
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccHTTPLoadBalancerResource_emptyPlan(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-test-lb")
	resourceName := "xcsh_http_loadbalancer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLoadBalancerConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
				),
			},
			{
				Config:             testAccHTTPLoadBalancerConfig_basicSystem(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// =============================================================================
// HELPER: Import state ID function
// =============================================================================
func testAccHTTPLoadBalancerImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
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
func TestAccHTTPLoadBalancerResource_planChecks(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLoadBalancerConfig_basicSystem(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					},
				},
			},
			{
				Config: testAccHTTPLoadBalancerConfig_labelsUpdateSystem(rName, "staging"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: testAccHTTPLoadBalancerConfig_labelsUpdateSystem(rName, "staging"),
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
// TEST: Known values
// =============================================================================
func TestAccHTTPLoadBalancerResource_knownValues(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLoadBalancerConfig_basicSystem(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rName)),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("namespace"), knownvalue.StringExact("system")),
					},
				},
			},
		},
	})
}

// =============================================================================
// TEST: Invalid name
// =============================================================================
func TestAccHTTPLoadBalancerResource_invalidName(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccHTTPLoadBalancerConfig_basicSystem("Invalid-NAME"),
				ExpectError: regexp.MustCompile(`(?i)(invalid|name|must)`),
			},
		},
	})
}

// =============================================================================
// TEST: Requires replace on name change
// =============================================================================
func TestAccHTTPLoadBalancerResource_requiresReplace(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName1 := acctest.RandomName("tf-test-lb")
	rName2 := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLoadBalancerConfig_basicSystem(rName1),
				Check:  acctest.CheckResourceExists(resourceName),
			},
			{
				Config: testAccHTTPLoadBalancerConfig_basicSystem(rName2),
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
// TEST: HTTPS auto-cert TLS termination
// =============================================================================
func TestAccHTTPLoadBalancerResource_httpsAutoCert(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLBConfig_httpsAutoCertSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccHTTPLoadBalancerImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"timeouts",
					"http.dns_volterra_managed",
					"l7_ddos_protection",
					"https_auto_cert.connection_idle_timeout",
					"https_auto_cert.http_redirect",
				},
			},
		},
	})
}

// =============================================================================
// TEST: With origin pool reference
// =============================================================================
func TestAccHTTPLoadBalancerResource_withOriginPool(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			if err := acctest.CheckResourceDestroyed("xcsh_http_loadbalancer")(s); err != nil {
				return err
			}
			return acctest.CheckOriginPoolDestroyed(s)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLoadBalancerConfig_withOriginPoolSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_route_pools.#", "1"),
				),
			},
		},
	})
}

// =============================================================================
// TEST: With WAF (app_firewall reference)
// =============================================================================
func TestAccHTTPLoadBalancerResource_withWAF(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			if err := acctest.CheckResourceDestroyed("xcsh_http_loadbalancer")(s); err != nil {
				return err
			}
			return acctest.CheckAppFirewallDestroyed(s)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLoadBalancerConfig_withWAFSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
				),
			},
		},
	})
}

// =============================================================================
// TEST: Full security stack (WAF + origin pool + rate limit + threat mesh)
// =============================================================================
func TestAccHTTPLoadBalancerResource_securityStack(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			if err := acctest.CheckResourceDestroyed("xcsh_http_loadbalancer")(s); err != nil {
				return err
			}
			if err := acctest.CheckAppFirewallDestroyed(s); err != nil {
				return err
			}
			if err := acctest.CheckOriginPoolDestroyed(s); err != nil {
				return err
			}
			return acctest.CheckHealthcheckDestroyed(s)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLoadBalancerConfig_securityStackSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_route_pools.#", "1"),
				),
			},
		},
	})
}

// =============================================================================
// TEST: Full lifecycle
// =============================================================================
func TestAccHTTPLoadBalancerResource_fullLifecycle(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLoadBalancerConfig_withLabelsSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccHTTPLoadBalancerImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"timeouts",
					"http.dns_volterra_managed",
					"l7_ddos_protection",
				},
			},
			{
				Config: testAccHTTPLoadBalancerConfig_withDomainsSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "domains.#", "2"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: testAccHTTPLoadBalancerConfig_withDomainsSystem(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccHTTPLoadBalancerConfig_basicSystem(rName),
				Check:  acctest.CheckResourceExists(resourceName),
			},
		},
	})
}

// =============================================================================
// CONFIG HELPERS - Use "system" namespace
// =============================================================================

func testAccHTTPLoadBalancerConfig_basicSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name       = %[1]q
  namespace  = "system"

  labels = {
    environment = "test"
    managed_by  = "terraform"
  }

  domains = ["test.example.com"]

  http {
    port = 80
  }

  advertise_on_public_default_vip {}
}
`, name)
}

func testAccHTTPLoadBalancerConfig_withLabelsSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name       = %[1]q
  namespace  = "system"

  labels = {
    environment = "test"
    team        = "platform"
    managed_by  = "terraform"
  }

  domains = ["test.example.com"]

  http {
    port = 80
  }

  advertise_on_public_default_vip {}
}
`, name)
}

func testAccHTTPLoadBalancerConfig_withDomainsSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name       = %[1]q
  namespace  = "system"

  labels = {
    environment = "test"
  }

  domains = [
    "app.example.com",
    "api.example.com"
  ]

  http {
    port = 80
  }

  advertise_on_public_default_vip {}
}
`, name)
}

// =============================================================================
// TEST: JS challenge (OneOf: challenge variants)
// =============================================================================
func TestAccHTTPLoadBalancerResource_jsChallenge(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLBConfig_jsChallengeSystem(rName),
				Check:  acctest.CheckResourceExists(resourceName),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "http.dns_volterra_managed", "l7_ddos_protection"},
				ImportStateIdFunc:       testAccHTTPLoadBalancerImportStateIdFunc(resourceName),
			},
		},
	})
}

// =============================================================================
// TEST: LB algorithm variants
// =============================================================================
func TestAccHTTPLoadBalancerResource_leastActive(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLBConfig_leastActiveSystem(rName),
				Check:  acctest.CheckResourceExists(resourceName),
			},
		},
	})
}

func TestAccHTTPLoadBalancerResource_sourceIpStickiness(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLBConfig_sourceIpStickinessSystem(rName),
				Check:  acctest.CheckResourceExists(resourceName),
			},
		},
	})
}

// =============================================================================
// TEST: Rate limiting with rate_limit reference
// =============================================================================
func TestAccHTTPLoadBalancerResource_withRateLimit(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			if err := acctest.CheckResourceDestroyed("xcsh_http_loadbalancer")(s); err != nil {
				return err
			}
			return acctest.CheckRateLimiterDestroyed(s)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLBConfig_withRateLimitSystem(rName),
				Check:  acctest.CheckResourceExists(resourceName),
			},
		},
	})
}

// =============================================================================
// TEST: User identification reference
// =============================================================================
func TestAccHTTPLoadBalancerResource_userIdentification(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			if err := acctest.CheckResourceDestroyed("xcsh_http_loadbalancer")(s); err != nil {
				return err
			}
			return acctest.CheckUserIdentificationDestroyed(s)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLBConfig_userIdentificationSystem(rName),
				Check:  acctest.CheckResourceExists(resourceName),
			},
		},
	})
}

// =============================================================================
// TEST: Do not advertise
// =============================================================================
func TestAccHTTPLoadBalancerResource_doNotAdvertise(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLBConfig_doNotAdvertiseSystem(rName),
				Check:  acctest.CheckResourceExists(resourceName),
			},
		},
	})
}

// =============================================================================
// TEST: IP reputation enabled
// =============================================================================
func TestAccHTTPLoadBalancerResource_ipReputation(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLBConfig_ipReputationSystem(rName),
				Check:  acctest.CheckResourceExists(resourceName),
			},
		},
	})
}

// =============================================================================
// TEST: Switch protocol (http → https_auto_cert)
// =============================================================================
func TestAccHTTPLoadBalancerResource_switchProtocol(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName1 := acctest.RandomName("tf-test-lb")
	rName2 := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLoadBalancerConfig_basicSystem(rName1),
				Check:  acctest.CheckResourceExists(resourceName),
			},
			{
				Config: testAccHTTPLBConfig_httpsAutoCertSystem(rName2),
				Check:  acctest.CheckResourceExists(resourceName),
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
// TEST: Switch LB algorithm
// =============================================================================
func TestAccHTTPLoadBalancerResource_switchAlgorithm(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "xcsh_http_loadbalancer.test"
	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_http_loadbalancer"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPLoadBalancerConfig_basicSystem(rName),
				Check:  acctest.CheckResourceExists(resourceName),
			},
			{
				Config: testAccHTTPLBConfig_leastActiveSystem(rName),
				Check:  acctest.CheckResourceExists(resourceName),
			},
			{
				Config: testAccHTTPLBConfig_sourceIpStickinessSystem(rName),
				Check:  acctest.CheckResourceExists(resourceName),
			},
		},
	})
}

// =============================================================================
// TEST: NEGATIVE — conflicting OneOf blocks should fail
// =============================================================================
func TestAccHTTPLoadBalancerResource_conflictHttpAndHttps(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccHTTPLBConfig_conflictProtocolSystem(rName),
				ExpectError: regexp.MustCompile(`(?i)(conflict|mutually exclusive|only one|Client Error|BAD_REQUEST|Invalid)`),
			},
		},
	})
}

// =============================================================================
// DOMAIN-SPECIFIC CONFIG HELPERS
// =============================================================================

func testAccHTTPLBConfig_jsChallengeSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"
  domains   = ["test.example.com"]

  http {
    port = 80
  }

  js_challenge {
    js_script_delay = 5000
    cookie_expiry   = 3600
  }

  advertise_on_public_default_vip {}
}
`, name)
}

func testAccHTTPLBConfig_leastActiveSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"
  domains   = ["test.example.com"]

  http {
    port = 80
  }

  least_active {}

  advertise_on_public_default_vip {}
}
`, name)
}

func testAccHTTPLBConfig_sourceIpStickinessSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"
  domains   = ["test.example.com"]

  http {
    port = 80
  }

  source_ip_stickiness {}

  advertise_on_public_default_vip {}
}
`, name)
}

func testAccHTTPLBConfig_withRateLimitSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"
  domains   = ["test.example.com"]

  http {
    port = 80
  }

  rate_limit {
    rate_limiter {
      total_number     = 100
      unit             = "MINUTE"
      burst_multiplier = 10
    }
    no_ip_allowed_list {}
  }

  advertise_on_public_default_vip {}
}
`, name)
}

func testAccHTTPLBConfig_userIdentificationSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_user_identification" "test" {
  name      = %[1]q
  namespace = "system"

  rules {
    client_ip {}
  }
}

resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"
  domains   = ["test.example.com"]

  http {
    port = 80
  }

  user_identification {
    name      = xcsh_user_identification.test.name
    namespace = "system"
  }

  advertise_on_public_default_vip {}
}
`, name)
}

func testAccHTTPLBConfig_doNotAdvertiseSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"
  domains   = ["test.example.com"]

  http {
    port = 80
  }

  do_not_advertise {}
}
`, name)
}

func testAccHTTPLBConfig_ipReputationSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"
  domains   = ["test.example.com"]

  http {
    port = 80
  }

  enable_ip_reputation {}

  advertise_on_public_default_vip {}
}
`, name)
}

func testAccHTTPLBConfig_conflictProtocolSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"
  domains   = ["test.example.com"]

  http {
    port = 80
  }

  https_auto_cert {
    add_hsts              = false
    no_mtls               {}
    default_header        {}
    enable_path_normalize {}
    non_default_loadbalancer {}
  }

  advertise_on_public_default_vip {}
}
`, name)
}

func testAccHTTPLBConfig_httpsAutoCertSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"

  domains = ["test.example.com"]

  https_auto_cert {
    add_hsts              = false
    no_mtls               {}
    default_header        {}
    enable_path_normalize {}
    non_default_loadbalancer {}
  }

  advertise_on_public_default_vip {}
}
`, name)
}

func testAccHTTPLoadBalancerConfig_withOriginPoolSystem(name string) string {
	return fmt.Sprintf(`
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
  same_as_endpoint_port {}
}

resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"

  domains = ["test.example.com"]

  http {
    port = 80
  }

  default_route_pools {
    pool {
      name      = xcsh_origin_pool.test.name
      namespace = "system"
    }
    weight   = 1
    priority = 1
  }

  advertise_on_public_default_vip {}
}
`, name)
}

func testAccHTTPLoadBalancerConfig_withWAFSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_app_firewall" "test" {
  name      = %[1]q
  namespace = "system"

  default_detection_settings {}
  allow_all_response_codes {}
  blocking {}
  use_default_blocking_page {}
  default_bot_setting {}
  default_anonymization {}
}

resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"

  domains = ["test.example.com"]

  http {
    port = 80
  }

  app_firewall {
    name      = xcsh_app_firewall.test.name
    namespace = "system"
  }

  advertise_on_public_default_vip {}
}
`, name)
}

func testAccHTTPLoadBalancerConfig_securityStackSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 3
  unhealthy_threshold = 1
  timeout             = 3
  interval            = 15

  http_health_check {
    path        = "/health"
    host_header = "example.com"
  }
}

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

  healthcheck {
    name      = xcsh_healthcheck.test.name
    namespace = "system"
  }

  no_tls {}
  same_as_endpoint_port {}
}

resource "xcsh_app_firewall" "test" {
  name      = %[1]q
  namespace = "system"

  default_detection_settings {}
  allow_all_response_codes {}
  blocking {}
  use_default_blocking_page {}
  default_bot_setting {}
  default_anonymization {}
}

resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"

  domains = ["test.example.com"]

  http {
    port = 80
  }

  default_route_pools {
    pool {
      name      = xcsh_origin_pool.test.name
      namespace = "system"
    }
    weight   = 1
    priority = 1
  }

  app_firewall {
    name      = xcsh_app_firewall.test.name
    namespace = "system"
  }

  enable_malicious_user_detection {}
  enable_threat_mesh {}

  advertise_on_public_default_vip {}
}
`, name)
}

func testAccHTTPLoadBalancerConfig_labelsUpdateSystem(name, env string) string {
	return fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name       = %[1]q
  namespace  = "system"

  labels = {
    environment = %[2]q
    managed_by  = "terraform"
  }

  domains = ["test.example.com"]

  http {
    port = 80
  }

  advertise_on_public_default_vip {}
}
`, name, env)
}

// =============================================================================
// NEGATIVE: Conflicting advertise options (OneOf violation)
// =============================================================================
func TestAccHTTPLoadBalancerResource_conflictAdvertiseOptions(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"
  domains   = ["test.example.com"]

  http {
    port = 80
  }

  advertise_on_public_default_vip {}
  do_not_advertise {}
}
`, rName),
				ExpectError: regexp.MustCompile(`(?i)(conflict|mutually exclusive|only one|Client Error|BAD_REQUEST|Invalid|these attributes cannot)`),
			},
		},
	})
}

// =============================================================================
// NEGATIVE: Conflicting challenge options (OneOf violation)
// =============================================================================
func TestAccHTTPLoadBalancerResource_conflictChallengeOptions(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"
  domains   = ["test.example.com"]

  http {
    port = 80
  }

  js_challenge {
    js_script_delay = 5000
    cookie_expiry   = 3600
  }
  no_challenge {}

  advertise_on_public_default_vip {}
}
`, rName),
				ExpectError: regexp.MustCompile(`(?i)(conflict|mutually exclusive|only one|Client Error|BAD_REQUEST|Invalid|these attributes cannot)`),
			},
		},
	})
}

// =============================================================================
// NEGATIVE: Conflicting LB algorithm options (OneOf violation)
// =============================================================================
func TestAccHTTPLoadBalancerResource_conflictAlgorithmOptions(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-test-lb")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "xcsh_http_loadbalancer" "test" {
  name      = %[1]q
  namespace = "system"
  domains   = ["test.example.com"]

  http {
    port = 80
  }

  round_robin {}
  least_active {}

  advertise_on_public_default_vip {}
}
`, rName),
				ExpectError: regexp.MustCompile(`(?i)(conflict|mutually exclusive|only one|Client Error|BAD_REQUEST|Invalid|these attributes cannot)`),
			},
		},
	})
}
