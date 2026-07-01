// tests/dns_test.go — Terratest unit tests for stratum-factory-gcp-dns.
//
// Run locally:
//
//	cd tests && go mod tidy && go test -v -timeout 10m ./...
//
// Tests are credential-free:
//   - Positive tests use `tofu validate` (static analysis, no API calls).
//   - Negative tests use InitAndPlanE with invalid var values (validation fires before provider auth).
package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validNetwork = "https://www.googleapis.com/compute/v1/projects/stratum-dev-sandbox/global/networks/stratum-dev-vpc"

func moduleDir(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("..")
	require.NoError(t, err)
	return root
}

func tofuOptions(t *testing.T, vars map[string]interface{}) *terraform.Options {
	t.Helper()
	return &terraform.Options{
		TerraformDir:    moduleDir(t),
		TerraformBinary: "tofu",
		Vars:            vars,
		NoColor:         true,
	}
}

func printReport(t *testing.T, rows [][]string) {
	t.Helper()
	header := fmt.Sprintf("%-60s %-10s %s", "Test", "Result", "Detail")
	sep := strings.Repeat("─", 100)
	t.Log("\n" + sep)
	t.Log("  STRATUM-FACTORY — GCP DNS Module Test Report")
	t.Log(sep)
	t.Log(header)
	t.Log(sep)
	for _, row := range rows {
		t.Logf("  %-60s %-10s %s", row[0], row[1], row[2])
	}
	t.Log(sep)
	if summaryFile := os.Getenv("GITHUB_STEP_SUMMARY"); summaryFile != "" {
		f, err := os.OpenFile(summaryFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		defer f.Close()
		fmt.Fprintln(f, "## GCP DNS Module Test Report")
		fmt.Fprintln(f, "| Test | Result | Detail |")
		fmt.Fprintln(f, "|------|--------|--------|")
		for _, row := range rows {
			fmt.Fprintf(f, "| %s | %s | %s |\n", row[0], row[1], row[2])
		}
	}
}

// TestDnsValidate verifies that tofu validate succeeds with minimal valid inputs.
func TestDnsValidate(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"name_prefix": "stratum-dev",
		"network":     validNetwork,
		"dns_name":    "stratum.dev.",
	})
	terraform.Init(t, opts)
	_, err := terraform.RunTerraformCommandE(t, opts, "validate")
	result, detail := "✅ PASS", "validate completed"
	if err != nil {
		result, detail = "❌ FAIL", err.Error()
	}
	printReport(t, [][]string{{"DnsValidate", result, detail}})
	require.NoError(t, err, "tofu validate must pass for valid inputs")
}

// TestDnsValidateWithRecord verifies that validate succeeds when a DNS A record is provided.
func TestDnsValidateWithRecord(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"name_prefix": "stratum-dev",
		"network":     validNetwork,
		"dns_name":    "stratum.dev.",
		"records": []interface{}{
			map[string]interface{}{
				"name":    "bastion.stratum.dev.",
				"type":    "A",
				"ttl":     float64(300),
				"rrdatas": []interface{}{"203.0.113.10"},
			},
		},
	})
	terraform.Init(t, opts)
	_, err := terraform.RunTerraformCommandE(t, opts, "validate")
	result, detail := "✅ PASS", "validate completed with DNS A record"
	if err != nil {
		result, detail = "❌ FAIL", err.Error()
	}
	printReport(t, [][]string{{"DnsValidateWithRecord", result, detail}})
	require.NoError(t, err, "tofu validate must pass with a valid DNS record")
}

// TestDnsRejectsInvalidEnvironment verifies that an invalid environment is rejected.
func TestDnsRejectsInvalidEnvironment(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "production",
		"project_id":  "stratum-dev-sandbox",
		"name_prefix": "stratum-dev",
		"network":     validNetwork,
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "✅ PASS", "plan correctly rejected invalid environment"
	if err == nil {
		result, detail = "❌ FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"DnsRejectsInvalidEnvironment", result, detail}})
	assert.Error(t, err, "must reject environment='production'")
}

// TestDnsRejectsInvalidNamePrefix verifies that a name_prefix starting with a digit is rejected.
func TestDnsRejectsInvalidNamePrefix(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"name_prefix": "1bad",
		"network":     validNetwork,
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "✅ PASS", "plan correctly rejected name_prefix starting with digit"
	if err == nil {
		result, detail = "❌ FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"DnsRejectsInvalidNamePrefix", result, detail}})
	assert.Error(t, err, "must reject name_prefix='1bad'")
}

// TestDnsRejectsDnsNameWithoutTrailingDot verifies that a dns_name without trailing dot is rejected.
func TestDnsRejectsDnsNameWithoutTrailingDot(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"name_prefix": "stratum-dev",
		"network":     validNetwork,
		"dns_name":    "stratum.dev",
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "✅ PASS", "plan correctly rejected dns_name without trailing dot"
	if err == nil {
		result, detail = "❌ FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"DnsRejectsDnsNameWithoutTrailingDot", result, detail}})
	assert.Error(t, err, "must reject dns_name without trailing dot")
}

// TestDnsRejectsInvalidNetwork verifies that an invalid network self-link is rejected.
func TestDnsRejectsInvalidNetwork(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"name_prefix": "stratum-dev",
		"network":     "not-a-self-link",
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "✅ PASS", "plan correctly rejected invalid network self-link"
	if err == nil {
		result, detail = "❌ FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"DnsRejectsInvalidNetwork", result, detail}})
	assert.Error(t, err, "must reject invalid network self-link")
}

// TestDnsRejectsInvalidRecordType verifies that an unsupported record type is rejected.
func TestDnsRejectsInvalidRecordType(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"name_prefix": "stratum-dev",
		"network":     validNetwork,
		"records": []interface{}{
			map[string]interface{}{
				"name":    "test.stratum.dev.",
				"type":    "UNKNOWN",
				"rrdatas": []interface{}{"1.2.3.4"},
			},
		},
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "✅ PASS", "plan correctly rejected unknown record type"
	if err == nil {
		result, detail = "❌ FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"DnsRejectsInvalidRecordType", result, detail}})
	assert.Error(t, err, "must reject type='UNKNOWN'")
}

// TestDnsRejectsZeroTtl verifies that a record with ttl = 0 is rejected.
func TestDnsRejectsZeroTtl(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"name_prefix": "stratum-dev",
		"network":     validNetwork,
		"records": []interface{}{
			map[string]interface{}{
				"name":    "test.stratum.dev.",
				"type":    "A",
				"ttl":     float64(0),
				"rrdatas": []interface{}{"10.0.0.1"},
			},
		},
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "✅ PASS", "plan correctly rejected record with ttl=0"
	if err == nil {
		result, detail = "❌ FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"DnsRejectsZeroTtl", result, detail}})
	assert.Error(t, err, "must reject record with ttl=0")
}

// TestDnsRejectsRecordNameWithoutTrailingDot verifies that a record name without a trailing dot is rejected.
func TestDnsRejectsRecordNameWithoutTrailingDot(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"name_prefix": "stratum-dev",
		"network":     validNetwork,
		"records": []interface{}{
			map[string]interface{}{
				"name":    "bastion.stratum.dev", // missing trailing dot
				"type":    "A",
				"rrdatas": []interface{}{"10.0.0.1"},
			},
		},
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "✅ PASS", "plan correctly rejected record name without trailing dot"
	if err == nil {
		result, detail = "❌ FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"DnsRejectsRecordNameWithoutTrailingDot", result, detail}})
	assert.Error(t, err, "must reject record name without trailing dot")
}

// TestNoTerraformBinary verifies that no .tf file references the `terraform` binary.
func TestNoTerraformBinary(t *testing.T) {
	t.Parallel()
	tfFiles, _ := filepath.Glob("../*.tf")
	issues := []string{}
	for _, f := range tfFiles {
		content, err := os.ReadFile(f)
		require.NoError(t, err)
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.Contains(trimmed, "terraform") &&
				!strings.HasPrefix(trimmed, "#") &&
				!strings.HasPrefix(trimmed, "//") &&
				trimmed != "terraform {" &&
				!strings.HasPrefix(trimmed, "required_version") &&
				!strings.HasPrefix(trimmed, "required_providers") &&
				!strings.HasPrefix(trimmed, "backend") &&
				!strings.Contains(trimmed, "TerraformBinary") {
				if strings.Contains(trimmed, "\"terraform\"") || strings.Contains(trimmed, "`terraform`") {
					issues = append(issues, fmt.Sprintf("%s:%d: %s", f, i+1, trimmed))
				}
			}
		}
	}
	result, detail := "✅ PASS", "no terraform binary references found"
	if len(issues) > 0 {
		result, detail = "❌ FAIL", fmt.Sprintf("found: %v", issues)
	}
	printReport(t, [][]string{{"NoTerraformBinary", result, detail}})
	assert.Empty(t, issues, "no .tf file should reference the 'terraform' binary")
}
