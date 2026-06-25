// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/f5xc-salesdemos/terraform-provider-xcsh/internal/acctest"
)

func TestAccCrlDataSource_basic(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	rName := acctest.RandomName("tf-acc-test-crl")
	nsName := acctest.RandomName("tf-acc-test-ns")
	resourceName := "xcsh_crl.test"
	dataSourceName := "data.xcsh_crl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {Source: "hashicorp/time"},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCrlDataSourceConfig_basic(nsName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespace", resourceName, "namespace"),
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
				),
			},
		},
	})
}

func testAccCrlDataSourceConfig_basic(nsName, name string) string {
	// CRL resource uses server_address, server_port, refresh_interval, timeout, and http_access
	// NOT crl_url which is not a valid attribute
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

resource "xcsh_crl" "test" {
  depends_on = [time_sleep.wait_for_namespace]
  name       = %[2]q
  namespace  = xcsh_namespace.test.name

  server_address   = "crl.example.com"
  server_port      = 80
  refresh_interval = 24
  timeout          = 30

  http_access {
    path = "/crl/test.crl"
  }
}

data "xcsh_crl" "test" {
  depends_on = [xcsh_crl.test]
  name       = xcsh_crl.test.name
  namespace  = xcsh_crl.test.namespace
}
`, nsName, name))
}
