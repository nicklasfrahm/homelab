resource "local_file" "config" {
  content  = yamlencode(var.config)
  filename = "${path.cwd}/deploy/tofu/out/${var.config.metadata.name}.out"
}
