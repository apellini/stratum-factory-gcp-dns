# Factory Module: `stratum-factory-gcp-dns`

Provisions a GCP private DNS managed zone for the STRATUM platform, with optional custom record sets.

## Purpose

Creates a `google_dns_managed_zone` (private visibility, attached to the VPC) and optional
`google_dns_record_set` resources. When `records` is empty (the default), no record set
resources are created.

## Usage

```hcl
module "dev_dns" {
  source = "git::https://github.com/apellini/stratum-factory-gcp-dns.git?ref=v0.2.0"

  environment = "dev"
  project_id  = "stratum-dev-sandbox"
  name_prefix = "stratum-dev"
  network     = module.dev_vpc.vpc_self_link
  dns_name    = "stratum.dev."
  tags        = { environment = "dev", managed_by = "opentofu" }

  # Optional: custom DNS records
  records = [
    {
      name    = "bastion.stratum.dev."
      type    = "A"
      ttl     = 300
      rrdatas = ["203.0.113.10"]
    },
  ]
}
```

## Inputs

| Name | Type | Required | Validation | Description |
|------|------|----------|------------|-------------|
| `environment` | `string` | yes | one of `dev`, `stage`, `main` | Deployment environment |
| `project_id` | `string` | yes | non-empty, no whitespace | GCP project ID |
| `name_prefix` | `string` | yes | 3–24 chars, lowercase alphanumeric/hyphens, starts with letter | Resource name prefix |
| `network` | `string` | yes | valid GCP network self-link URI | VPC self-link (from vpc module output) |
| `dns_name` | `string` | optional (default `"stratum.dev."`) | FQDN ending with `.` | DNS name for the managed zone |
| `records` | `list(object)` | optional (default `[]`) | see below | DNS record sets to create |
| `tags` | `map(string)` | optional (default `{}`) | all keys/values non-empty | Labels applied to all resources |

### `records` object attributes

| Attribute | Type | Default | Validation | Description |
|-----------|------|---------|------------|-------------|
| `name` | `string` | — | FQDN ending with `.` | Fully-qualified DNS name for the record |
| `type` | `string` | — | A, AAAA, CNAME, MX, NS, PTR, SOA, SRV, TXT, or CAA | DNS record type |
| `ttl` | `number` | `300` | > 0 | Time-to-live in seconds |
| `rrdatas` | `list(string)` | — | non-empty | Record data strings |

## Outputs

| Name | Type | Description |
|------|------|-------------|
| `dns_zone_id` | `string` | Unique identifier of the DNS managed zone |
| `dns_zone_name` | `string` | Name of the DNS managed zone (use when creating record sets) |
| `dns_name` | `string` | DNS name of the zone (FQDN with trailing dot) |
| `name_servers` | `list(string)` | DNS name servers for this zone |
| `record_names` | `map(string)` | Map of `"<name>/<type>"` keys → DNS record set names |

## Factory rules applied

- **Stateless** — no local state, no remote state reads
- **Input-only** — all configuration via `variables.tf`; nothing hardcoded
- **Strict validation** — every variable has a `validation` block
- **No secrets** — no sensitive values in code or outputs
- **Documented for humans and RAG** — this README + inline comments

## Release

```hcl
source = "git::https://github.com/apellini/stratum-factory-gcp-dns.git?ref=v0.2.0"
```
