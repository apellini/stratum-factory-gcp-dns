# stratum-factory-gcp-dns/variables.tf
#
# FACTORY pattern: all configuration received as input. NOTHING is hardcoded.
# Every variable carries a strict validation block.

variable "environment" {
  description = "Deployment environment. Must be one of the approved environment names."
  type        = string

  validation {
    condition     = contains(["dev", "stage", "main"], var.environment)
    error_message = "environment must be one of: dev, stage, main."
  }
}

variable "project_id" {
  description = "GCP project ID to deploy into. Must be non-empty with no whitespace."
  type        = string

  validation {
    condition     = length(var.project_id) > 0 && !can(regex("\\s", var.project_id))
    error_message = "project_id must be a non-empty string with no whitespace."
  }
}

variable "name_prefix" {
  description = "Prefix applied to all resource names. 3-24 lowercase alphanumeric or hyphens, starts with a letter."
  type        = string

  validation {
    condition     = can(regex("^[a-z][a-z0-9-]{2,23}$", var.name_prefix))
    error_message = "name_prefix must be 3-24 chars, start with a letter, lowercase alphanumeric or hyphens only."
  }
}

variable "network" {
  description = "Self-link URI of the VPC network to attach the private DNS zone to. Passed from the VPC module via the Wrapper."
  type        = string

  validation {
    condition     = can(regex("^https://www\\.googleapis\\.com/compute/v1/projects/.+/global/networks/.+$", var.network))
    error_message = "network must be a valid GCP compute network self-link URI (https://www.googleapis.com/compute/v1/projects/...)."
  }
}

variable "dns_name" {
  description = "DNS name for the managed zone. Must be a fully-qualified domain name ending with a dot. Example: stratum.dev."
  type        = string
  default     = "stratum.dev."

  validation {
    condition     = can(regex("^[a-z0-9][a-z0-9.-]+\\.$", var.dns_name))
    error_message = "dns_name must be a fully-qualified domain name ending with a dot (e.g. stratum.dev.)."
  }
}

variable "tags" {
  description = "Map of labels to apply to all resources. Keys and values must be non-empty strings."
  type        = map(string)
  default     = {}

  validation {
    condition     = alltrue([for k, v in var.tags : length(k) > 0 && length(v) > 0])
    error_message = "All tag keys and values must be non-empty strings."
  }
}

variable "records" {
  description = <<-EOT
    List of DNS record sets to create in the managed zone. Each entry maps to one
    google_dns_record_set resource.

    Attributes:
      name    - (required) Fully-qualified DNS name ending with a dot (e.g. "bastion.stratum.dev.").
      type    - (required) Record type. One of: A, AAAA, CNAME, MX, NS, PTR, SOA, SRV, TXT, CAA.
      ttl     - (optional, default 300) Time-to-live in seconds (must be > 0).
      rrdatas - (required) Non-empty list of record data strings.

    Default [] — no record sets created when empty.
  EOT
  type = list(object({
    name    = string
    type    = string
    ttl     = optional(number, 300)
    rrdatas = list(string)
  }))
  default = []

  validation {
    condition = alltrue([
      for r in var.records :
      contains(["A", "AAAA", "CNAME", "MX", "NS", "PTR", "SOA", "SRV", "TXT", "CAA"], r.type)
    ])
    error_message = "Each record type must be one of: A, AAAA, CNAME, MX, NS, PTR, SOA, SRV, TXT, CAA."
  }

  validation {
    condition     = alltrue([for r in var.records : r.ttl > 0])
    error_message = "Each record ttl must be greater than 0."
  }

  validation {
    condition     = alltrue([for r in var.records : length(r.rrdatas) > 0])
    error_message = "Each record rrdatas must be a non-empty list."
  }

  validation {
    condition     = alltrue([for r in var.records : can(regex("^[a-zA-Z0-9._-]+\\.$", r.name))])
    error_message = "Each record name must be a fully-qualified DNS name ending with a dot (e.g. \"bastion.stratum.dev.\")."
  }
}
