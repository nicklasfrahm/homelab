locals {
  # Load all region configuration files.
  region_files = fileset("${path.cwd}/configs/regions", "*.yaml")
  region_configs = {
    for filename in local.region_files :
    replace(filename, ".yaml", "") => yamldecode(file("${path.cwd}/configs/regions/${filename}"))
  }

  machine_files = fileset("${path.cwd}/configs/machines", "*.yaml")
  machine_configs = {
    for filename in local.machine_files :
    replace(filename, ".yaml", "") => yamldecode(file("${path.cwd}/configs/machines/${filename}"))
  }

  hardware_profile_files = fileset("${path.cwd}/configs/hardwareprofiles", "*.yaml")
  hardware_profile_configs = {
    for filename in local.hardware_profile_files :
    replace(filename, ".yaml", "") => yamldecode(file("${path.cwd}/configs/hardwareprofiles/${filename}"))
  }
}

module "region" {
  source = "./modules/region"

  for_each = local.region_configs

  region = each.value
  machines = local.machine_configs
  hardware_profiles = local.hardware_profile_configs
}
