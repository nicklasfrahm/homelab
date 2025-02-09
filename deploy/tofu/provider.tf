terraform {
  backend "gcs" {
    bucket = "nicklasfrahm"
    prefix  = "tofu/state/homelab"
  }
}
