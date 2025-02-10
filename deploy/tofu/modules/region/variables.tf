variable "region" {
  description = "The configuration of the region."
  type = object({
    metadata = object({
      name = string
    })
    spec = object({
      provider = string
      baremetal = object({
        controlplanes = list(object({
          name = string
        }))
      })
    })
  })
}

variable "machines" {
  description = "The configuration of all machines."
  type = map(object({
    metadata = object({
      name = string
    })
    spec = object({
      hardware = object({
        vendor = string
        model = string
      })
      interfaces = list(object({
        mac = string
      }))
    })
  }))
}

variable "hardware_profiles" {
  description = "The configuration of all hardware profiles."
  type = map(object({
    spec = object({
      storage = object({
        osDisk = object({
          name = string
        })
      })
    })
  }))
}
