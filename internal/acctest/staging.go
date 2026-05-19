// Copyright (c) 2026 Robin Mordasiewicz. MIT License.

package acctest

import (
	"os"
	"testing"
)

// StagingTestRegex is the Go test regex matching the curated staging test subset.
// These are representative basic lifecycle tests across all domains.
const StagingTestRegex = "^TestAcc(Namespace|Healthcheck|OriginPool|AppFirewall|VirtualSite|AlertPolicy|AlertReceiver|ServicePolicy|GlobalLogReceiver)Resource_basic$"

// StagingPreCheck verifies staging-specific prerequisites.
// Requires F5XC_API_URL and F5XC_API_TOKEN (token auth for staging).
func StagingPreCheck(t *testing.T) {
	t.Helper()

	if os.Getenv(EnvF5XCURL) == "" {
		t.Skip("F5XC_API_URL must be set for staging tests")
	}
	if os.Getenv(EnvF5XCToken) == "" {
		t.Skip("F5XC_API_TOKEN must be set for staging tests")
	}
}
