locals {
  name = var.region.metadata.name
  talos_version = yamldecode(file("${path.cwd}/config.yaml")).talos.version
}

# Create the Talos secret bundle.
resource "talos_machine_secrets" "this" {
  talos_version = local.talos_version
}
