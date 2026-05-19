// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package acctest

import (
	"regexp"
	"testing"
)

func TestStagingTestRegex_Compiles(t *testing.T) {
	_, err := regexp.Compile(StagingTestRegex)
	if err != nil {
		t.Fatalf("StagingTestRegex does not compile: %v", err)
	}
}

func TestStagingTestRegex_MatchesExpected(t *testing.T) {
	re := regexp.MustCompile(StagingTestRegex)

	shouldMatch := []string{
		"TestAccNamespaceResource_basic",
		"TestAccHealthcheckResource_basic",
		"TestAccOriginPoolResource_basic",
		"TestAccAppFirewallResource_basic",
		"TestAccVirtualSiteResource_basic",
		"TestAccAlertPolicyResource_basic",
		"TestAccAlertReceiverResource_basic",
		"TestAccServicePolicyResource_basic",
		"TestAccGlobalLogReceiverResource_basic",
	}

	shouldNotMatch := []string{
		"TestAccNamespaceResource_allAttributes",
		"TestAccNamespaceResource_update",
		"TestMockNamespaceResource_basic",
		"TestAccNamespaceDataSource_basic",
	}

	for _, name := range shouldMatch {
		if !re.MatchString(name) {
			t.Errorf("StagingTestRegex should match %q but does not", name)
		}
	}

	for _, name := range shouldNotMatch {
		if re.MatchString(name) {
			t.Errorf("StagingTestRegex should NOT match %q but does", name)
		}
	}
}
