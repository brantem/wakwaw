variable "name" {
  default = "wakwaw"
  type    = string
}

variable "do_token" {
  type      = string
  sensitive = true
}

variable "region" {
  type = string
}

variable "managers" {
  default = 1
  type    = number
}

variable "manager_size" {
  default = "s-1vcpu-2gb"
  type    = string
}

variable "workers" {
  default = 1
  type    = number
}

variable "worker_size" {
  default = "s-1vcpu-1gb"
  type    = string
}

variable "cidr_block" {
  default = "10.0.0.0/24"
  type    = string
}
