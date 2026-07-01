# stratum-factory-gcp-dns/main.tf
#
# FACTORY MODULE — GCP private DNS managed zone for STRATUM.
#
# STATELESS and INPUT-ONLY. State owned by the calling Wrapper.
# Consumed by the Wrapper as:
#   source = "git::https://github.com/apellini/stratum-factory-gcp-dns.git?ref=v<semver>"

# ── Private DNS managed zone ──────────────────────────────────────────────────
resource "google_dns_managed_zone" "zone" {
  project     = var.project_id
  name        = "${var.name_prefix}-dns-zone"
  dns_name    = var.dns_name
  description = "STRATUM ${var.environment} private DNS zone (${var.dns_name}) — managed by OpenTofu"
  visibility  = "private"

  private_visibility_config {
    networks {
      network_url = var.network
    }
  }

  labels = var.tags
}
