data "manidae_parameter" "root_volume_size_gb" {
  name         = "root_volume_size_gb"
  display_name = "Root Volume Size (GB)"
  description  = "How large should the root volume for the instance be?"
  default      = 30
  type         = "number"

  validation {
    min = 20
  }
}

data "manidae_parameter" "instance_type" {
  name         = "instance_type"
  display_name = "Instance type"
  description  = "What instance type should your workspace use?"
  default      = "SA2.MEDIUM8"

  option {
    name  = "2 vCPU, 2 GiB RAM"
    value = "SA2.MEDIUM2"
  }
  option {
    name  = "2 vCPU, 4 GiB RAM"
    value = "SA2.MEDIUM4"
  }
  option {
    name  = "2 vCPU, 8 GiB RAM"
    value = "SA2.MEDIUM8"
  }
  option {
    name  = "4 vCPU, 4 GiB RAM"
    value = "SA2.LARGE4"
  }
  option {
    name  = "4 vCPU, 8 GiB RAM"
    value = "SA2.LARGE8"
  }
  option {
    name  = "4 vCPU, 16 GiB RAM"
    value = "SA2.LARGE16"
  }
}

data "manidae_parameter" "docker_image" {
  name         = "docker_image"
  display_name = "Docker image"
  description  = "The docker image to create container"
  default      = "codercom/enterprise-base:ubuntu"
  type         = "string"
}

