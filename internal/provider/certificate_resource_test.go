// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/f5xc-salesdemos/terraform-provider-xcsh/internal/acctest"
)

func TestAccCertificateResource_basic(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-cert")
	nsName := acctest.RandomName("tf-acc-test-ns")
	resourceName := "xcsh_certificate.test"

	// Generate test certificates dynamically for CI/CD compatibility
	certs := acctest.MustGenerateTestCertificates()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {Source: "hashicorp/time"},
		},
		CheckDestroy: acctest.CheckResourceDestroyed("xcsh_certificate"),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateConfig_basic(nsName, rName, certs),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", nsName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "private_key"},
				ImportStateIdFunc:       testAccCertificateImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccCertificateResource_emptyPlan(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-cert")
	nsName := acctest.RandomName("tf-acc-test-ns")
	certs := acctest.MustGenerateTestCertificates()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        map[string]resource.ExternalProvider{"time": {Source: "hashicorp/time"}},
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_certificate"),
		Steps: []resource.TestStep{
			{Config: testAccCertificateConfig_basic(nsName, rName, certs)},
			{Config: testAccCertificateConfig_basic(nsName, rName, certs), PlanOnly: true, ExpectNonEmptyPlan: false},
		},
	})
}

func TestAccCertificateResource_withLabels(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-cert")
	nsName := acctest.RandomName("tf-acc-test-ns")
	resourceName := "xcsh_certificate.test"
	certs := acctest.MustGenerateTestCertificates()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders:        map[string]resource.ExternalProvider{"time": {Source: "hashicorp/time"}},
		CheckDestroy:             acctest.CheckResourceDestroyed("xcsh_certificate"),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateConfig_withLabels(nsName, rName, certs),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "labels.environment", "test"),
				),
			},
		},
	})
}

func testAccCertificateConfig_withLabels(nsName, name string, certs *acctest.TestCertificates) string {
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

resource "xcsh_certificate" "test" {
  depends_on = [time_sleep.wait_for_namespace]
  name       = %[2]q
  namespace  = xcsh_namespace.test.name

  labels = {
    environment = "test"
  }

  certificate_url = "string:///%[3]s"

  private_key {
    clear_secret_info {
      url = "string:///%[4]s"
    }
  }

  disable_ocsp_stapling {}
}
`, nsName, name, certs.ServerCertBase64, certs.ServerKeyBase64))
}

func testAccCertificateImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
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

func testAccCertificateConfig_basic(nsName, name string, certs *acctest.TestCertificates) string {
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

resource "xcsh_certificate" "test" {
  depends_on = [time_sleep.wait_for_namespace]
  name       = %[2]q
  namespace  = xcsh_namespace.test.name

  certificate_url = "string:///%[3]s"

  private_key {
    clear_secret_info {
      url = "string:///%[4]s"
    }
  }

  disable_ocsp_stapling {}
}
`, nsName, name, certs.ServerCertBase64, certs.ServerKeyBase64))
}
