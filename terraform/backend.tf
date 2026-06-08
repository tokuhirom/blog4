# tfstate を さくらのオブジェクトストレージ (S3 互換) に置く。
#
# - backend ブロックは変数を使えないため、静的値のみ。
#   認証情報 (S3 アクセスキー) は環境変数で渡す:
#     AWS_ACCESS_KEY_ID     = <Object Storage のアクセスキー>
#     AWS_SECRET_ACCESS_KEY = <同シークレット>
# - AWS ではないので各種 skip フラグと path-style が必要。
# - state ロックは S3 ネイティブのロックファイル (use_lockfile, TF 1.10+)。
#   DynamoDB 相当は不要。
# - state バケット (blog4-tfstate) は鶏卵問題のため Terraform 管理外。
#   コンパネで private バケットを作り、versioning を有効化しておく。README 参照。
terraform {
  backend "s3" {
    bucket = "blog4-tfstate"
    key    = "blog4/terraform.tfstate"
    region = "jp-north-1"

    endpoints = {
      s3 = "https://s3.isk01.sakurastorage.jp"
    }

    use_path_style = true
    use_lockfile   = true

    # 非 AWS S3 互換ストレージ向け: AWS 固有の検証/メタAPIを呼ばせない
    skip_credentials_validation = true
    skip_requesting_account_id  = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    skip_s3_checksum            = true
  }
}
