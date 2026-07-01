# stratum-factory-gcp-dns/outputs.tf
#
# FACTORY MODULE outputs — documented for humans and RAG.

output "dns_zone_id" {
  description = <<-EOT
    Unique identifier of the DNS managed zone.
    Type: string. Example: "projects/stratum-dev-sandbox/managedZones/stratum-dev-dns-zone"
  EOT
  value       = google_dns_managed_zone.zone.id
}

output "dns_zone_name" {
  description = <<-EOT
    Name of the DNS managed zone. Use this when creating google_dns_record_set resources.
    Type: string. Example: "stratum-dev-dns-zone"
  EOT
  value       = google_dns_managed_zone.zone.name
}

output "dns_name" {
  description = <<-EOT
    The DNS name of the managed zone (fully-qualified, with trailing dot).
    Type: string. Example: "stratum.dev."
  EOT
  value       = google_dns_managed_zone.zone.dns_name
}

output "name_servers" {
  description = <<-EOT
    List of DNS name servers for this zone.
    Type: list(string).
  EOT
  value       = google_dns_managed_zone.zone.name_servers
}

output "record_names" {
  description = <<-EOT
    Map of "<name>/<type>" keys to the DNS record set names created by this module.
    Type: map(string).
    Example: { "bastion.stratum.dev./A" = "bastion.stratum.dev." }
    Empty map when records = [].
  EOT
  value = { for k, r in google_dns_record_set.records : k => r.name }
}
