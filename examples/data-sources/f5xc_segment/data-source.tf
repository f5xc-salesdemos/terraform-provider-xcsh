# Segment Data Source Example
# Retrieves information about an existing Segment

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Segment by name
data "f5xc_segment" "example" {
  name      = "example-segment"
  namespace = "staging"
}

output "segment_id" {
  value = data.f5xc_segment.example.id
}
