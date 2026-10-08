variable "repositories" {
  type = list(string)
}

variable "keep_images" {
  type    = number
  default = 50
}

variable "force_delete" {
  type    = bool
  default = false
}

variable "tags" {
  type    = map(string)
  default = {}
}
