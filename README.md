# Factory Module: stratum-factory-gcp-dns

## Purpose

Creates a private GCP DNS managed zone (`stratum.dev.`) for the STRATUM platform.
Enables `*.stratum.dev` resolution to the in-cluster k3s ingress during the HLD
bootstrap phase (line 451). The zone is attached to a caller-supplied VPC and is
visible only within that network.

**Factory rules applied:**
- [x] Stateless — no remote state reads, no backend configuration
- [x] Input-only — all configuration via variables, nothing hardcoded
- [x] Strict validation — every variable has a validation block with an actionable error message
- [x] No secrets — no credentials, tokens, or sensitive values in code
- [x] Documented for humans and RAG — clear inputs, outputs, and usage examples

## Usage

```hcl
module "dns" {
  source = "git::https://github.com/apellini/stratum-factory-gcp-dns.git?ref=v0.1.0"

  environment = "dev"
  project_id  = "stratum-dev-sandbox"
  name_prefix = "stratum-dev"
  network     = module.vpc.network_self_link
  dns_name    = "stratum.dev."

  tags = {
    environment = "dev"
    managed-by  = "opentofu"
    module      = "stratum-factory-gcp-dns"
  }
}
```

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
| `environment` | Deployment environment. Must be one of: `dev`, `stage`, `main`. | `string` | — | yes |
| `project_id` | GCP project ID to deploy into. Non-empty, no whitespace. | `string` | — | yes |
| `name_prefix` | Prefix for all resource names. 3-24 chars, starts with a letter, lowercase alphanumeric or hyphens. | `string` | — | yes |
| `network` | Self-link URI of the VPC network to attach the private DNS zone to. | `string` | — | yes |
| `dns_name` | Fully-qualified DNS name ending with a dot. Example: `stratum.dev.` | `string` | `"stratum.dev."` | no |
| `tags` | Map of labels to apply to all resources. Keys and values must be non-empty strings. | `map(string)` | `{}` | no |

## Outputs

| Name | Description | Type | Example |
|------|-------------|------|---------|
| `dns_zone_id` | Unique identifier of the DNS managed zone. | `string` | `"projects/stratum-dev-sandbox/managedZones/stratum-dev-dns-zone"` |
| `dns_zone_name` | Name of the DNS managed zone. Use when creating `google_dns_record_set` resources. | `string` | `"stratum-dev-dns-zone"` |
| `dns_name` | The DNS name of the managed zone (fully-qualified, with trailing dot). | `string` | `"stratum.dev."` |
| `name_servers` | List of DNS name servers for this zone. | `list(string)` | `["ns-cloud-a1.googledomains.com.", ...]` |

## Release

Consumed via git tag — never by branch:

```hcl
source = "git::https://github.com/apellini/stratum-factory-gcp-dns.git?ref=v0.1.0"
```

Tags follow semver. The Wrapper pins the tag explicitly. Do not use `?ref=main`.

## Requirements

| Tool | Version |
|------|---------|
| OpenTofu | >= 1.8.0 |
| hashicorp/google | ~> 6.0 |

Provider registry: `registry.opentofu.org`

## Testing

Tests are credential-free (no GCP API calls). Validation errors fire before provider
authentication, so both positive and negative tests run locally without credentials.

```bash
cd tests
go mod tidy
go test -v -timeout 10m ./...
```
