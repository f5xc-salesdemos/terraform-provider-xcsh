// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/f5xc-salesdemos/terraform-provider-xcsh/internal/acctest"
)

func TestAccAppTypeDataSource_basic(t *testing.T) {
	acctest.SkipIfNotAccTest(t)
	acctest.PreCheck(t)

	// app_type resources must be created in the "shared" namespace
	rName := acctest.RandomName("tf-acc-test-apptype")
	resourceName := "xcsh_app_type.test"
	dataSourceName := "data.xcsh_app_type.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppTypeDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespace", resourceName, "namespace"),
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
				),
			},
		},
	})
}

func testAccAppTypeDataSourceConfig_basic(name string) string {
	// app_type resources must be created in the "shared" namespace
	return acctest.ConfigCompose(
		acctest.ProviderConfig(),
		fmt.Sprintf(`
resource "xcsh_app_type" "test" {
  name      = %[1]q
  namespace = "shared"
}

data "xcsh_app_type" "test" {
  depends_on = [xcsh_app_type.test]
  name       = xcsh_app_type.test.name
  namespace  = xcsh_app_type.test.namespace
}
`, name))
}
