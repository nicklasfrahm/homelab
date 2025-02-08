locals {
  region_files = fileset("${path.cwd}/configs/regions", "*.yaml")
  region_configs = { for f in local.region_files : replace(f, ".yaml", "") => yamldecode(file("${path.cwd}/configs/regions/${f}")) }
}

module "region" {
  source = "./modules/region"

  for_each = local.region_configs

  config = each.value
}
