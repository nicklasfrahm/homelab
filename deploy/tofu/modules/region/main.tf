locals {
  name = var.region.metadata.name
}

# Get the latest Talos version from image factory.
data "talos_image_factory_versions" "this" {
  filters = {
    stable_versions_only = true
  }
}

# Extract the latest Talos version.
locals {
  talos_version = element(data.talos_image_factory_versions.this.talos_versions, length(data.talos_image_factory_versions.this.talos_versions) - 1)
}

# Create the Talos secret bundle.
resource "talos_machine_secrets" "this" {
  talos_version = local.talos_version
}

# Create a file with the Talos version.
resource "local_file" "version" {
  filename = "${path.cwd}/deploy/tofu/out/${local.name}"
  content = local.talos_version
}
