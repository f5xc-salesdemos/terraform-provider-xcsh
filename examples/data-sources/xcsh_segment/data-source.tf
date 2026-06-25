# Segment Data Source Example
# Retrieves information about an existing Segment

terraform {
  required_version = ">= 1.0"

  required_providers {
    xcsh = {
      source  = "f5xc-salesdemos/xcsh"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing Segment by name
data "xcsh_segment" "example" {
  name      = "example-segment"
  namespace = "staging"
}

output "segment_id" {
  value = data.xcsh_segment.example.id
}
