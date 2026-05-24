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

	"github.com/f5xc-salesdemos/terraform-provider-f5xc/internal/acctest"
)

// =============================================================================
// TEST: Basic app_firewall creation with default detection settings
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccAppFirewallResource_basic(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckAppFirewallExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", "system"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Import test
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccAppFirewallImportStateIdFunc(resourceName),
			},
		},
	})
}

// =============================================================================
// TEST: App firewall with labels and description
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccAppFirewallResource_withLabels(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_withLabelsSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckAppFirewallExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Test application firewall"),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "labels.team", "security"),
				),
			},
		},
	})
}

// =============================================================================
// TEST: App firewall with blocking mode
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccAppFirewallResource_blocking(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_blockingSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckAppFirewallExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

// =============================================================================
// TEST: App firewall with monitoring mode
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccAppFirewallResource_monitoring(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_monitoringSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckAppFirewallExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

// =============================================================================
// TEST: Resource disappears (deleted outside Terraform)
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// =============================================================================
func TestAccAppFirewallResource_disappears(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckAppFirewallExists(resourceName),
					acctest.CheckAppFirewallDisappears(resourceName),
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
func TestAccAppFirewallResource_emptyPlan(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckAppFirewallExists(resourceName),
				),
			},
			{
				Config:             testAccAppFirewallConfig_basicSystem(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// =============================================================================
// HELPER: Import state ID function
// =============================================================================
func testAccAppFirewallImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
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
// TEST: All attributes including labels, annotations, description
// =============================================================================
func TestAccAppFirewallResource_allAttributes(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_allAttributesSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckAppFirewallExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Full attributes test"),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "labels.managed_by", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr(resourceName, "annotations.purpose", "acceptance-testing"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "disable", "description"},
				ImportStateIdFunc:       testAccAppFirewallImportStateIdFunc(resourceName),
			},
		},
	})
}

// =============================================================================
// TEST: Update labels lifecycle
// =============================================================================
func TestAccAppFirewallResource_updateLabels(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_withLabelsSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckAppFirewallExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "test"),
				),
			},
			{
				Config: testAccAppFirewallConfig_withUpdatedLabelsSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckAppFirewallExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "staging"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: testAccAppFirewallConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckAppFirewallExists(resourceName),
				),
			},
		},
	})
}

// =============================================================================
// TEST: Plan checks (create, update, noop)
// =============================================================================
func TestAccAppFirewallResource_planChecks(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_basicSystem(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					},
				},
			},
			{
				Config: testAccAppFirewallConfig_withLabelsSystem(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: testAccAppFirewallConfig_withLabelsSystem(rName),
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
func TestAccAppFirewallResource_knownValues(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_basicSystem(rName),
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
// TEST: Invalid name validation
// =============================================================================
func TestAccAppFirewallResource_invalidName(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppFirewallConfig_basicSystem("Invalid-NAME-Test"),
				ExpectError: regexp.MustCompile(`(?i)(invalid|name|must)`),
			},
		},
	})
}

// =============================================================================
// TEST: Name too long validation
// =============================================================================
func TestAccAppFirewallResource_nameTooLong(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppFirewallConfig_basicSystem("tf-test-this-name-is-way-too-long-and-should-fail-validation-check"),
				ExpectError: regexp.MustCompile(`(?i)(invalid|name|length|long|exceed|character)`),
			},
		},
	})
}

// =============================================================================
// TEST: Empty name validation
// =============================================================================
func TestAccAppFirewallResource_emptyName(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppFirewallConfig_basicSystem(""),
				ExpectError: regexp.MustCompile(`(?i)(invalid|name|empty|required|blank)`),
			},
		},
	})
}

// =============================================================================
// TEST: Name change requires replacement
// =============================================================================
func TestAccAppFirewallResource_requiresReplace(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName1 := acctest.RandomName("tf-test-waf")
	rName2 := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_basicSystem(rName1),
				Check:  acctest.CheckAppFirewallExists(resourceName),
			},
			{
				Config: testAccAppFirewallConfig_basicSystem(rName2),
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
// TEST: Switch from blocking to monitoring mode
// =============================================================================
func TestAccAppFirewallResource_switchMode_blockingToMonitoring(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_blockingSystem(rName),
				Check:  acctest.CheckAppFirewallExists(resourceName),
			},
			{
				Config: testAccAppFirewallConfig_monitoringSystem(rName),
				Check:  acctest.CheckAppFirewallExists(resourceName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: testAccAppFirewallConfig_blockingSystem(rName),
				Check:  acctest.CheckAppFirewallExists(resourceName),
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
// TEST: Full lifecycle (create → update → import → no-drift → strip → delete)
// =============================================================================
func TestAccAppFirewallResource_fullLifecycle(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_allAttributesSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckAppFirewallExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "Full attributes test"),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "test"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "disable", "description"},
				ImportStateIdFunc:       testAccAppFirewallImportStateIdFunc(resourceName),
			},
			{
				Config: testAccAppFirewallConfig_monitoringSystem(rName),
				Check:  acctest.CheckAppFirewallExists(resourceName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: testAccAppFirewallConfig_monitoringSystem(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccAppFirewallConfig_basicSystem(rName),
				Check:  acctest.CheckAppFirewallExists(resourceName),
			},
		},
	})
}

// =============================================================================
// CONFIG HELPERS
// =============================================================================

// testAccAppFirewallConfig_basicSystem uses the "system" namespace
// to avoid creating test namespaces (namespace DELETE returns 501)
func testAccAppFirewallConfig_basicSystem(name string) string {
	return fmt.Sprintf(`
resource "f5xc_app_firewall" "test" {
  name       = %[1]q
  namespace  = "system"

  # Use default detection settings for simplicity
  default_detection_settings {}

  # Allow all response codes
  allow_all_response_codes {}

  # Blocking mode
  blocking {}

  # Use default blocking page
  use_default_blocking_page {}

  # Use default bot settings
  default_bot_setting {}

  # Use default anonymization
  default_anonymization {}
}
`, name)
}

func testAccAppFirewallConfig_withLabelsSystem(name string) string {
	return fmt.Sprintf(`
resource "f5xc_app_firewall" "test" {
  name        = %[1]q
  namespace   = "system"
  description = "Test application firewall"

  labels = {
    environment = "test"
    team        = "security"
  }

  # Use default detection settings
  default_detection_settings {}

  # Allow all response codes
  allow_all_response_codes {}

  # Blocking mode
  blocking {}

  # Use default blocking page
  use_default_blocking_page {}

  # Use default bot settings
  default_bot_setting {}

  # Use default anonymization
  default_anonymization {}
}
`, name)
}

func testAccAppFirewallConfig_blockingSystem(name string) string {
	return fmt.Sprintf(`
resource "f5xc_app_firewall" "test" {
  name       = %[1]q
  namespace  = "system"

  # Use default detection settings
  default_detection_settings {}

  # Allow all response codes
  allow_all_response_codes {}

  # Blocking mode - actively block malicious requests
  blocking {}

  # Use default blocking page
  use_default_blocking_page {}

  # Use default bot settings
  default_bot_setting {}

  # Use default anonymization
  default_anonymization {}
}
`, name)
}

func testAccAppFirewallConfig_allAttributesSystem(name string) string {
	return fmt.Sprintf(`
resource "f5xc_app_firewall" "test" {
  name        = %[1]q
  namespace   = "system"
  description = "Full attributes test"

  labels = {
    environment = "test"
    managed_by  = "terraform-acceptance-test"
  }

  annotations = {
    purpose = "acceptance-testing"
  }

  default_detection_settings {}
  allow_all_response_codes {}
  blocking {}
  use_default_blocking_page {}
  default_bot_setting {}
  default_anonymization {}
}
`, name)
}

func testAccAppFirewallConfig_withUpdatedLabelsSystem(name string) string {
	return fmt.Sprintf(`
resource "f5xc_app_firewall" "test" {
  name        = %[1]q
  namespace   = "system"
  description = "Test application firewall"

  labels = {
    environment = "staging"
    team        = "platform"
  }

  default_detection_settings {}
  allow_all_response_codes {}
  blocking {}
  use_default_blocking_page {}
  default_bot_setting {}
  default_anonymization {}
}
`, name)
}

func testAccAppFirewallConfig_monitoringSystem(name string) string {
	return fmt.Sprintf(`
resource "f5xc_app_firewall" "test" {
  name       = %[1]q
  namespace  = "system"

  default_detection_settings {}
  allow_all_response_codes {}
  monitoring {}
  use_default_blocking_page {}
  default_bot_setting {}
  default_anonymization {}
}
`, name)
}

// =============================================================================
// DOMAIN-SPECIFIC CONFIG HELPERS — OneOf variant coverage
// =============================================================================

// OneOf: blocking_page / use_default_blocking_page
func testAccAppFirewallConfig_customBlockingPage(name string) string {
	return fmt.Sprintf(`
resource "f5xc_app_firewall" "test" {
  name      = %[1]q
  namespace = "system"

  default_detection_settings {}
  allow_all_response_codes {}
  blocking {}
  default_bot_setting {}
  default_anonymization {}

  blocking_page {
    blocking_page = "https://example.com/blocked.html"
    response_code = "Forbidden"
  }
}
`, name)
}

// OneOf: bot_protection_setting / default_bot_setting
func testAccAppFirewallConfig_botProtection(name string) string {
	return fmt.Sprintf(`
resource "f5xc_app_firewall" "test" {
  name      = %[1]q
  namespace = "system"

  default_detection_settings {}
  allow_all_response_codes {}
  blocking {}
  use_default_blocking_page {}
  default_anonymization {}

  bot_protection_setting {
    good_bot_action       = "REPORT"
    malicious_bot_action  = "BLOCK"
    suspicious_bot_action = "REPORT"
  }
}
`, name)
}

// OneOf: disable_anonymization / default_anonymization / custom_anonymization
func testAccAppFirewallConfig_disableAnonymization(name string) string {
	return fmt.Sprintf(`
resource "f5xc_app_firewall" "test" {
  name      = %[1]q
  namespace = "system"

  default_detection_settings {}
  allow_all_response_codes {}
  blocking {}
  use_default_blocking_page {}
  default_bot_setting {}

  disable_anonymization {}
}
`, name)
}

// OneOf: detection_settings / default_detection_settings
func testAccAppFirewallConfig_detectionSettings(name string) string {
	return fmt.Sprintf(`
resource "f5xc_app_firewall" "test" {
  name      = %[1]q
  namespace = "system"

  allow_all_response_codes {}
  blocking {}
  use_default_blocking_page {}
  default_bot_setting {}
  default_anonymization {}

  detection_settings {
    default_violation_settings {}
    default_bot_setting {}
    enable_suppression {}
    enable_threat_campaigns {}
    signature_selection_setting {
      high_medium_accuracy_signatures {}
      default_attack_type_settings {}
    }
  }
}
`, name)
}

// OneOf: enable_ai_enhancements / disable_ai_enhancements
func testAccAppFirewallConfig_aiEnhancements(name string) string {
	return fmt.Sprintf(`
resource "f5xc_app_firewall" "test" {
  name      = %[1]q
  namespace = "system"

  default_detection_settings {}
  allow_all_response_codes {}
  blocking {}
  use_default_blocking_page {}
  default_bot_setting {}
  default_anonymization {}

  enable_ai_enhancements {
    mitigate_high_risk_action {}
  }
}
`, name)
}

// OneOf: allowed_response_codes / allow_all_response_codes
func testAccAppFirewallConfig_allowedResponseCodes(name string) string {
	return fmt.Sprintf(`
resource "f5xc_app_firewall" "test" {
  name      = %[1]q
  namespace = "system"

  default_detection_settings {}
  blocking {}
  use_default_blocking_page {}
  default_bot_setting {}
  default_anonymization {}

  allowed_response_codes {
    response_code = [200, 204, 301, 302]
  }
}
`, name)
}

// =============================================================================
// DOMAIN-SPECIFIC TESTS — OneOf variant coverage
// =============================================================================

func TestAccAppFirewallResource_customBlockingPage(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_customBlockingPage(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckAppFirewallExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "blocking_page.response_code", "Forbidden"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccAppFirewallImportStateIdFunc(resourceName),
			},
			{
				Config: testAccAppFirewallConfig_customBlockingPage(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

func TestAccAppFirewallResource_botProtection(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_botProtection(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckAppFirewallExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bot_protection_setting.good_bot_action", "REPORT"),
					resource.TestCheckResourceAttr(resourceName, "bot_protection_setting.malicious_bot_action", "BLOCK"),
					resource.TestCheckResourceAttr(resourceName, "bot_protection_setting.suspicious_bot_action", "REPORT"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "bot_protection_setting.%", "bot_protection_setting.good_bot_action", "bot_protection_setting.malicious_bot_action", "bot_protection_setting.suspicious_bot_action"},
				ImportStateIdFunc:       testAccAppFirewallImportStateIdFunc(resourceName),
			},
			// Switch to default bot settings
			{
				Config: testAccAppFirewallConfig_basicSystem(rName),
				Check:  acctest.CheckAppFirewallExists(resourceName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate)},
				},
			},
			{
				Config: testAccAppFirewallConfig_basicSystem(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

func TestAccAppFirewallResource_disableAnonymization(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_disableAnonymization(rName),
				Check:  acctest.CheckAppFirewallExists(resourceName),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccAppFirewallImportStateIdFunc(resourceName),
			},
			// Switch to default anonymization
			{
				Config: testAccAppFirewallConfig_basicSystem(rName),
				Check:  acctest.CheckAppFirewallExists(resourceName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate)},
				},
			},
		},
	})
}

func TestAccAppFirewallResource_detectionSettings(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_detectionSettings(rName),
				Check:  acctest.CheckAppFirewallExists(resourceName),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccAppFirewallImportStateIdFunc(resourceName),
			},
			{
				Config: testAccAppFirewallConfig_detectionSettings(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

func TestAccAppFirewallResource_aiEnhancements(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_aiEnhancements(rName),
				Check:  acctest.CheckAppFirewallExists(resourceName),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccAppFirewallImportStateIdFunc(resourceName),
			},
			{
				Config: testAccAppFirewallConfig_aiEnhancements(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

func TestAccAppFirewallResource_allowedResponseCodes(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resourceName := "f5xc_app_firewall.test"
	rName := acctest.RandomName("tf-test-waf")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        acctest.ExternalProviders,
		CheckDestroy:             acctest.CheckAppFirewallDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAppFirewallConfig_allowedResponseCodes(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckAppFirewallExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "allowed_response_codes.response_code.#", "4"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccAppFirewallImportStateIdFunc(resourceName),
			},
			// Switch to allow_all
			{
				Config: testAccAppFirewallConfig_basicSystem(rName),
				Check:  acctest.CheckAppFirewallExists(resourceName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate)},
				},
			},
		},
	})
}
