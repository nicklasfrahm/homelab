locals {
  name = var.config.metadata.name
}

resource "local_file" "config" {
  filename = "${path.cwd}/deploy/tofu/out/${local.name}.yaml"
  content  = yamlencode(var.config)
}
