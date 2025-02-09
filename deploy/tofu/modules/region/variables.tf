variable "config" {
  type = object({
    metadata = object({
      name = string
    })
    spec = object({
      provider = string
    })
  })
  description = "The configuration of the region."
}
