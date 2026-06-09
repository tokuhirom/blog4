# WebAccel サイト (blog.64p.org の前段 CDN / TLS 終端)。
# TLS は Let's Encrypt 自動運用のため証明書リソースは管理しない。
resource "sakura_webaccel" "blog4" {
  name              = "blog.64p.org"
  domain_type       = "own_domain"
  domain            = "blog.64p.org"
  request_protocol  = "https-redirect"
  normalize_ae      = "gzip"
  default_cache_ttl = -1
  vary_support      = false

  origin_parameters = {
    type     = "web"
    protocol = "https"
    # AppRun (sakura_apprun_shared.blog4) の ingress ホスト名。
    origin = "app-bdc9ce86-732c-4b9f-9576-a4313eed0cbd.ingress.apprun.sakura.ne.jp"
  }
}

# WebAccel サイト (blog-attachments.64p.org → Object Storage の blog3-attachments バケット)。
#
# !!! 現状 import 不可 (provider 制約) のためコンパネ運用のまま。!!!
# sacloud/sakura v3.12 の sakura_webaccel は、バケット origin のサイトを import する際の
# Read で write-only の access_key_wo を受け取れず "origin_parameters must be provided to
# keep bucket credentials" で必ず失敗する (credentials_wo_version でも回避不可)。
# provider 側が対応したら、下記の確認済み設定で取り込む (値はコンパネと一致確認済み):
#
#   resource "sakura_webaccel" "attachments" {
#     name              = "blog3-attachments"
#     domain_type       = "own_domain"
#     domain            = "blog-attachments.64p.org"
#     request_protocol  = "https"
#     normalize_ae      = "gzip"
#     default_cache_ttl = 3600
#     vary_support      = false
#     origin_parameters = {
#       type                   = "bucket"
#       bucket_name            = "blog3-attachments"
#       endpoint               = "s3.isk01.sakurastorage.jp"
#       region                 = "jp-north-1"
#       use_document_index     = false
#       access_key_wo          = var.webaccel_attachments_access_key      # 1Password
#       secret_access_key_wo   = var.webaccel_attachments_secret_access_key
#       credentials_wo_version = 1
#     }
#   }
