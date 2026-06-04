provider "sakura" {
  # Authentication is supplied via environment variables:
  #   SAKURACLOUD_ACCESS_TOKEN
  #   SAKURACLOUD_ACCESS_TOKEN_SECRET
  # Zone defaults to "is1a" — TiDB CR is offered in is1 only.
  zone = var.zone
}
