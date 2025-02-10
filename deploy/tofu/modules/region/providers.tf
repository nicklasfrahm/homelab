terraform {
  required_version = ">= 1.9.0"

  required_providers {
    kind = {
      source = "tehcyx/kind"
      version = "0.7.0"
    }
    local = {
      source = "hashicorp/local"
      version = "2.5.2"
    }
  }
}
