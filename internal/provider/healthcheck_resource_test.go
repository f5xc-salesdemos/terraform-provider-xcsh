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
// HEALTHCHECK RESOURCE ACCEPTANCE TESTS
//
// Uses "system" namespace to avoid creating test namespaces that can't be deleted
// (namespace DELETE API returns 501 Not Implemented)
//
// Run with:
//
//	TF_ACC=1 XCSH_API_URL="..." XCSH_P12_FILE="..." XCSH_P12_PASSWORD="..." \
//	go test -v ./internal/provider/ -run TestAccHealthcheckResource -timeout 30m
//
// =============================================================================

// =============================================================================
// TEST: Basic healthcheck creation with API verification
// =============================================================================
func TestAccHealthcheckResource_basic(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", "system"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Import state verification
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccHealthcheckImportStateIdFunc(resourceName),
			},
		},
	})
}

// testAccHealthcheckImportStateIdFunc returns a function that generates the import ID
func testAccHealthcheckImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		namespace := rs.Primary.Attributes["namespace"]
		name := rs.Primary.Attributes["name"]
		if namespace == "" || name == "" {
			return "", fmt.Errorf("namespace or name not set in state")
		}
		return fmt.Sprintf("%s/%s", namespace, name), nil
	}
}

// =============================================================================
// TEST: All optional attributes (labels, annotations)
// =============================================================================
func TestAccHealthcheckResource_allAttributes(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_allAttributesSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", "system"),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "labels.managed_by", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr(resourceName, "annotations.purpose", "acceptance-testing"),
					resource.TestCheckResourceAttr(resourceName, "annotations.owner", "ci-cd"),
					acctest.CheckHealthcheckAttributes(resourceName,
						map[string]string{
							"environment": "test",
							"managed_by":  "terraform-acceptance-test",
						},
						map[string]string{
							"purpose": "acceptance-testing",
							"owner":   "ci-cd",
						},
					),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "disable", "description"},
				ImportStateIdFunc:       testAccHealthcheckImportStateIdFunc(resourceName),
			},
		},
	})
}

// =============================================================================
// TEST: Update labels
// =============================================================================
func TestAccHealthcheckResource_updateLabels(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_withLabelsSystem(rName, "test", "terraform"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "labels.managed_by", "terraform"),
				),
			},
			{
				Config: testAccHealthcheckConfig_withLabelsSystem(rName, "staging", "terraform-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "staging"),
					resource.TestCheckResourceAttr(resourceName, "labels.managed_by", "terraform-updated"),
					acctest.CheckHealthcheckAttributes(resourceName,
						map[string]string{
							"environment": "staging",
							"managed_by":  "terraform-updated",
						},
						nil,
					),
				),
			},
			{
				Config: testAccHealthcheckConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckNoResourceAttr(resourceName, "labels.environment"),
					resource.TestCheckNoResourceAttr(resourceName, "labels.managed_by"),
				),
			},
		},
	})
}

// =============================================================================
// TEST: Update annotations
// =============================================================================
func TestAccHealthcheckResource_updateAnnotations(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_withAnnotationsSystem(rName, "value1", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "annotations.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "annotations.key2", "value2"),
					acctest.CheckHealthcheckAttributes(resourceName, nil,
						map[string]string{
							"key1": "value1",
							"key2": "value2",
						},
					),
				),
			},
			{
				Config: testAccHealthcheckConfig_withAnnotationsSystem(rName, "updated1", "updated2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "annotations.key1", "updated1"),
					resource.TestCheckResourceAttr(resourceName, "annotations.key2", "updated2"),
				),
			},
			{
				Config: testAccHealthcheckConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
				),
			},
		},
	})
}

// =============================================================================
// TEST: Resource disappears (deleted outside Terraform)
// =============================================================================
func TestAccHealthcheckResource_disappears(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					acctest.CheckHealthcheckDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// =============================================================================
// TEST: Empty plan after apply (no drift)
// =============================================================================
func TestAccHealthcheckResource_emptyPlan(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_allAttributesSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
				),
			},
			{
				Config: testAccHealthcheckConfig_allAttributesSystem(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// =============================================================================
// TEST: Plan checks (create, update, no-op)
// =============================================================================
func TestAccHealthcheckResource_planChecks(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_basicSystem(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
				),
			},
			{
				Config: testAccHealthcheckConfig_withLabelsSystem(rName, "test", "terraform"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: testAccHealthcheckConfig_withLabelsSystem(rName, "test", "terraform"),
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
func TestAccHealthcheckResource_knownValues(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_basicSystem(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(resourceName,
							tfjsonpath.New("name"),
							knownvalue.StringExact(rName),
						),
						plancheck.ExpectKnownValue(resourceName,
							tfjsonpath.New("namespace"),
							knownvalue.StringExact("system"),
						),
					},
				},
			},
		},
	})
}

// =============================================================================
// TEST: Invalid name error
// =============================================================================
func TestAccHealthcheckResource_invalidName(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccHealthcheckConfig_basicSystem("Invalid-NAME-Test"),
				ExpectError: regexp.MustCompile(`(?i)(invalid|name|must)`),
			},
		},
	})
}

// =============================================================================
// TEST: Name too long error
// =============================================================================
func TestAccHealthcheckResource_nameTooLong(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	longName := "tf-acc-test-this-name-is-way-too-long-and-should-fail-validation-check"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccHealthcheckConfig_basicSystem(longName),
				ExpectError: regexp.MustCompile(`(?i)(invalid|name|length|long|exceed|character)`),
			},
		},
	})
}

// =============================================================================
// TEST: Empty name error
// =============================================================================
func TestAccHealthcheckResource_emptyName(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccHealthcheckConfig_basicSystem(""),
				ExpectError: regexp.MustCompile(`(?i)(invalid|name|empty|required|blank)`),
			},
		},
	})
}

// =============================================================================
// TEST: Name change requires replacement
// =============================================================================
func TestAccHealthcheckResource_requiresReplace(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName1 := acctest.RandomName("tf-acc-test-hc")
	rName2 := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_basicSystem(rName1),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName1),
				),
			},
			{
				Config: testAccHealthcheckConfig_basicSystem(rName2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
				),
			},
		},
	})
}

// =============================================================================
// TEST: HTTP health check nested block
// =============================================================================
func TestAccHealthcheckResource_httpHealthCheck(t *testing.T) {
	// Testing actual error to debug schema type mismatch
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_httpHealthCheckSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "http_health_check.path", "/health"),
					resource.TestCheckResourceAttr(resourceName, "http_health_check.host_header", "example.com"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "http_health_check.use_http2"},
				ImportStateIdFunc:       testAccHealthcheckImportStateIdFunc(resourceName),
			},
		},
	})
}

// =============================================================================
// GROUP 1: Health Check Type Coverage
// =============================================================================

// TestAccHealthcheckResource_tcpHealthCheck_withPayload — covered in healthcheck_origin_pool_matrix_test.go
// TestAccHealthcheckResource_httpHealthCheck_allFields — covered in healthcheck_origin_pool_matrix_test.go (httpFullOptions)

func TestAccHealthcheckResource_httpHealthCheck_originServerName(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_httpOriginServerName(rName, "/health"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "http_health_check.path", "/health"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "http_health_check.use_http2"},
				ImportStateIdFunc:       testAccHealthcheckImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccHealthcheckResource_httpHealthCheck_statusCodes(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_httpStatusCodes(rName, "/health"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "http_health_check.expected_status_codes.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "http_health_check.expected_status_codes.0", "200"),
					resource.TestCheckResourceAttr(resourceName, "http_health_check.expected_status_codes.1", "201"),
					resource.TestCheckResourceAttr(resourceName, "http_health_check.expected_status_codes.2", "204"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "http_health_check.use_http2"},
				ImportStateIdFunc:       testAccHealthcheckImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccHealthcheckResource_httpHealthCheck_headersRemove(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_httpHeadersRemove(rName, "/health"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "http_health_check.request_headers_to_remove.#", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "http_health_check.use_http2"},
				ImportStateIdFunc:       testAccHealthcheckImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccHealthcheckResource_httpHealthCheck_http2(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_httpHttp2(rName, "/health"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "http_health_check.use_http2", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "http_health_check.use_http2"},
				ImportStateIdFunc:       testAccHealthcheckImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccHealthcheckResource_udpIcmpHealthCheck(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_udpIcmp(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccHealthcheckImportStateIdFunc(resourceName),
			},
		},
	})
}

// =============================================================================
// GROUP 2: Spec Attribute Coverage
// =============================================================================

func TestAccHealthcheckResource_description(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_withDescription(rName, "initial description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "initial description"),
				),
			},
			{
				Config: testAccHealthcheckConfig_withDescription(rName, "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "updated description"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "disable", "description"},
				ImportStateIdFunc:       testAccHealthcheckImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccHealthcheckResource_jitterPercent(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_withJitter(rName, 30),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "jitter_percent", "30"),
				),
			},
			{
				Config: testAccHealthcheckConfig_withJitter(rName, 50),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "jitter_percent", "50"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
				ImportStateIdFunc:       testAccHealthcheckImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccHealthcheckResource_thresholds_update(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_thresholds(rName, 1, 2, 3, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "healthy_threshold", "1"),
					resource.TestCheckResourceAttr(resourceName, "unhealthy_threshold", "2"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "3"),
					resource.TestCheckResourceAttr(resourceName, "interval", "5"),
				),
			},
			{
				Config: testAccHealthcheckConfig_thresholds(rName, 5, 3, 10, 30),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "healthy_threshold", "5"),
					resource.TestCheckResourceAttr(resourceName, "unhealthy_threshold", "3"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "interval", "30"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: testAccHealthcheckConfig_thresholds(rName, 5, 3, 10, 30),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccHealthcheckResource_thresholds_boundary(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_thresholds(rName, 1, 1, 1, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "healthy_threshold", "1"),
					resource.TestCheckResourceAttr(resourceName, "unhealthy_threshold", "1"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "1"),
					resource.TestCheckResourceAttr(resourceName, "interval", "1"),
				),
			},
			{
				Config: testAccHealthcheckConfig_thresholds(rName, 16, 16, 600, 600),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "healthy_threshold", "16"),
					resource.TestCheckResourceAttr(resourceName, "unhealthy_threshold", "16"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "600"),
					resource.TestCheckResourceAttr(resourceName, "interval", "600"),
				),
			},
		},
	})
}

// =============================================================================
// GROUP 3: Type Switching and Update Tests
// =============================================================================

func TestAccHealthcheckResource_switchType_tcpToHttp(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
				),
			},
			{
				Config: testAccHealthcheckConfig_httpHealthCheckSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "http_health_check.path", "/health"),
				),
			},
		},
	})
}

func TestAccHealthcheckResource_switchType_httpToUdp(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_httpHealthCheckSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
				),
			},
			{
				Config: testAccHealthcheckConfig_udpIcmp(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
				),
			},
		},
	})
}

func TestAccHealthcheckResource_httpHealthCheck_updatePath(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_httpWithPath(rName, "/health"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "http_health_check.path", "/health"),
				),
			},
			{
				Config: testAccHealthcheckConfig_httpWithPath(rName, "/healthz"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "http_health_check.path", "/healthz"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "http_health_check.use_http2"},
				ImportStateIdFunc:       testAccHealthcheckImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccHealthcheckResource_httpHealthCheck_switchHostToOrigin(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcheckConfig_httpHealthCheckSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "http_health_check.host_header", "example.com"),
				),
			},
			{
				Config: testAccHealthcheckConfig_httpOriginServerName(rName, "/health"),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
				),
			},
		},
	})
}

// =============================================================================
// GROUP 4: Full Lifecycle Validation
// =============================================================================

func TestAccHealthcheckResource_fullLifecycle(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")
	resourceName := "xcsh_healthcheck.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctest.CheckHealthcheckDestroyed,
		Steps: []resource.TestStep{
			// Step 1: Create with all optional fields (reuses matrix file's httpFullOptions config)
			{
				Config: testAccHealthcheckConfig_httpHealthCheckSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "http_health_check.path", "/health"),
				),
			},
			// Step 2: Import verification
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "http_health_check.use_http2"},
				ImportStateIdFunc:       testAccHealthcheckImportStateIdFunc(resourceName),
			},
			// Step 3: Update thresholds and timing
			{
				Config: testAccHealthcheckConfig_thresholds(rName, 5, 3, 10, 30),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "healthy_threshold", "5"),
					resource.TestCheckResourceAttr(resourceName, "interval", "30"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			// Step 4: Verify no drift
			{
				Config: testAccHealthcheckConfig_thresholds(rName, 5, 3, 10, 30),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 5: Strip down to minimal config
			{
				Config: testAccHealthcheckConfig_basicSystem(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckHealthcheckExists(resourceName),
				),
			},
		},
	})
}

// =============================================================================
// CONFIG HELPERS - Use "system" namespace
// =============================================================================

func testAccHealthcheckConfig_basicSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5

  tcp_health_check {}
}
`, name)
}

func testAccHealthcheckConfig_allAttributesSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5

  labels = {
    environment = "test"
    managed_by  = "terraform-acceptance-test"
  }

  annotations = {
    purpose = "acceptance-testing"
    owner   = "ci-cd"
  }

  tcp_health_check {}
}
`, name)
}

func testAccHealthcheckConfig_withLabelsSystem(name, environment, managedBy string) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5

  labels = {
    environment = %[2]q
    managed_by  = %[3]q
  }

  tcp_health_check {}
}
`, name, environment, managedBy)
}

func testAccHealthcheckConfig_withAnnotationsSystem(name, value1, value2 string) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5

  annotations = {
    key1 = %[2]q
    key2 = %[3]q
  }

  tcp_health_check {}
}
`, name, value1, value2)
}

func testAccHealthcheckConfig_httpHealthCheckSystem(name string) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5

  http_health_check {
    path        = "/health"
    host_header = "example.com"
  }
}
`, name)
}

func testAccHealthcheckConfig_httpOriginServerName(name, path string) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5

  http_health_check {
    path = %[2]q
    use_origin_server_name {}
  }
}
`, name, path)
}

func testAccHealthcheckConfig_httpStatusCodes(name, path string) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5

  http_health_check {
    path                  = %[2]q
    host_header           = "example.com"
    expected_status_codes = ["200", "201", "204"]
  }
}
`, name, path)
}

func testAccHealthcheckConfig_httpHeadersRemove(name, path string) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5

  http_health_check {
    path                       = %[2]q
    host_header                = "example.com"
    request_headers_to_remove  = ["X-Custom-Header", "X-Debug"]
  }
}
`, name, path)
}

func testAccHealthcheckConfig_httpHttp2(name, path string) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5

  http_health_check {
    path        = %[2]q
    host_header = "example.com"
    use_http2   = true
  }
}
`, name, path)
}

func testAccHealthcheckConfig_udpIcmp(name string) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5

  udp_icmp_health_check {}
}
`, name)
}

func testAccHealthcheckConfig_withDescription(name, description string) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name        = %[1]q
  namespace   = "system"
  description = %[2]q

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5

  tcp_health_check {}
}
`, name, description)
}

func testAccHealthcheckConfig_withJitter(name string, jitter int) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5
  jitter_percent      = %[2]d

  tcp_health_check {}
}
`, name, jitter)
}

func testAccHealthcheckConfig_thresholds(name string, healthy, unhealthy, timeout, interval int) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = %[2]d
  unhealthy_threshold = %[3]d
  timeout             = %[4]d
  interval            = %[5]d

  tcp_health_check {}
}
`, name, healthy, unhealthy, timeout, interval)
}

func testAccHealthcheckConfig_httpWithPath(name, path string) string {
	return fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5

  http_health_check {
    path        = %[2]q
    host_header = "example.com"
  }
}
`, name, path)
}

// =============================================================================
// NEGATIVE: Conflicting health check types (OneOf violation)
// =============================================================================
func TestAccHealthcheckResource_conflictTcpAndHttp(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-hc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "xcsh_healthcheck" "test" {
  name      = %[1]q
  namespace = "system"

  healthy_threshold   = 1
  unhealthy_threshold = 2
  timeout             = 3
  interval            = 5

  tcp_health_check {}
  http_health_check {
    path = "/health"
  }
}
`, rName),
				ExpectError: regexp.MustCompile(`(?i)(conflict|mutually exclusive|only one|Invalid|these attributes cannot)`),
			},
		},
	})
}
