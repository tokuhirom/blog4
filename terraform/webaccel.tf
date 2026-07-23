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
# sacloud/sakura v3.12.0 ではバケット origin のサイトを import できず
# ("origin_parameters must be provided to keep bucket credentials")、コンパネ運用のまま
# 保留していたが、v3.12.2 の "fix(webaccel): support import for webaccel bucket origin"
# で解消したので取り込む。
#
# バケットの認証情報 (access_key_wo / secret_access_key_wo) はここに書かない。
# write-only 属性なので Terraform は既存値を読まず、書かなければ送りもしない
# = コンパネで設定済みの認証情報がそのまま維持される。
# use_document_index も書かない (書くと値の追加で in-place update が発生する)。
# 差分ゼロで import できることは plan で確認済み。
resource "sakura_webaccel" "attachments" {
  name              = "blog3-attachments"
  domain_type       = "own_domain"
  domain            = "blog-attachments.64p.org"
  request_protocol  = "https"
  normalize_ae      = "gzip"
  default_cache_ttl = 3600
  vary_support      = false

  origin_parameters = {
    type        = "bucket"
    bucket_name = "blog3-attachments"
    endpoint    = "s3.isk01.sakurastorage.jp"
    region      = "jp-north-1"
  }
}

# 既存サイトを state に取り込む。apply で state に入ったあとは削除してよい
# (残しても no-op なので、履歴として当面残す)。
import {
  to = sakura_webaccel.attachments
  id = "113602500600"
}
