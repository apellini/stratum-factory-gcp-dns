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

# ── DNS record sets ───────────────────────────────────────────────────────────
# Creates one google_dns_record_set per entry in var.records.
# When records is empty (the default), no record set resources are created.
resource "google_dns_record_set" "records" {
  for_each = { for r in var.records : "${r.name}/${r.type}" => r }

  project      = var.project_id
  name         = each.value.name
  type         = each.value.type
  ttl          = each.value.ttl
  managed_zone = google_dns_managed_zone.zone.name
  rrdatas      = each.value.rrdatas
}
