// tests/dns_test.go -- Terratest unit tests for stratum-factory-gcp-dns.
//
// Run locally:
//
//	cd tests && go mod tidy && go test -v -timeout 10m ./...
//
// Tests are credential-free:
//   - Positive tests use tofu validate (static analysis, no API calls).
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
	header := fmt.Sprintf("%-50s %-10s %s", "Test", "Result", "Detail")
	sep := strings.Repeat("-", 90)
	t.Log("\n" + sep)
	t.Log("  STRATUM-FACTORY -- GCP DNS Module Test Report")
	t.Log(sep)
	t.Log(header)
	t.Log(sep)
	for _, row := range rows {
		t.Logf("  %-50s %-10s %s", row[0], row[1], row[2])
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

func TestDnsValidate(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"name_prefix": "stratum-dev",
		"network":     "https://www.googleapis.com/compute/v1/projects/stratum-dev-sandbox/global/networks/stratum-dev-vpc",
		"dns_name":    "stratum.dev.",
	})
	terraform.Init(t, opts)
	_, err := terraform.RunTerraformCommandE(t, opts, "validate")
	result, detail := "PASS", "validate completed"
	if err != nil {
		result, detail = "FAIL", err.Error()
	}
	printReport(t, [][]string{{"DnsValidate", result, detail}})
	require.NoError(t, err, "tofu validate must pass for valid inputs")
}

func TestDnsRejectsInvalidEnvironment(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "production",
		"project_id":  "stratum-dev-sandbox",
		"name_prefix": "stratum-dev",
		"network":     "https://www.googleapis.com/compute/v1/projects/stratum-dev-sandbox/global/networks/stratum-dev-vpc",
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "PASS", "plan correctly rejected invalid environment"
	if err == nil {
		result, detail = "FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"DnsRejectsInvalidEnvironment", result, detail}})
	assert.Error(t, err, "must reject environment=production")
}

func TestDnsRejectsInvalidNamePrefix(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"name_prefix": "1bad",
		"network":     "https://www.googleapis.com/compute/v1/projects/stratum-dev-sandbox/global/networks/stratum-dev-vpc",
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "PASS", "plan correctly rejected name_prefix starting with digit"
	if err == nil {
		result, detail = "FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"DnsRejectsInvalidNamePrefix", result, detail}})
	assert.Error(t, err, "must reject name_prefix=1bad")
}

func TestDnsRejectsDnsNameWithoutTrailingDot(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"name_prefix": "stratum-dev",
		"network":     "https://www.googleapis.com/compute/v1/projects/stratum-dev-sandbox/global/networks/stratum-dev-vpc",
		"dns_name":    "stratum.dev",
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "PASS", "plan correctly rejected dns_name without trailing dot"
	if err == nil {
		result, detail = "FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"DnsRejectsDnsNameWithoutTrailingDot", result, detail}})
	assert.Error(t, err, "must reject dns_name without trailing dot")
}

func TestDnsRejectsInvalidNetwork(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"name_prefix": "stratum-dev",
		"network":     "not-a-self-link",
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "PASS", "plan correctly rejected invalid network self-link"
	if err == nil {
		result, detail = "FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"DnsRejectsInvalidNetwork", result, detail}})
	assert.Error(t, err, "must reject invalid network self-link")
}

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
	result, detail := "PASS", "no terraform binary references found"
	if len(issues) > 0 {
		result, detail = "FAIL", fmt.Sprintf("found: %v", issues)
	}
	printReport(t, [][]string{{"NoTerraformBinary", result, detail}})
	assert.Empty(t, issues, "no .tf file should reference the terraform binary")
}
