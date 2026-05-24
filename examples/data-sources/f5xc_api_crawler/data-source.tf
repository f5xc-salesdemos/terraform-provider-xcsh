# API Crawler Data Source Example
# Retrieves information about an existing API Crawler

terraform {
  required_version = ">= 1.0"

  required_providers {
    f5xc = {
      source  = "f5xc-salesdemos/f5xc"
      version = ">= 0.1.0"
    }
  }
}

# Look up an existing API Crawler by name
data "f5xc_api_crawler" "example" {
  name      = "example-api-crawler"
  namespace = "staging"
}

output "api_crawler_id" {
  value = data.f5xc_api_crawler.example.id
}
