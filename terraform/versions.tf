terraform {
  required_version = "~> 1.14"

  required_providers {
    sakura = {
      source  = "sacloud/sakura"
      version = "~> 3.12"
    }
  }
}
