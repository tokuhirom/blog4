# 既存 AppRun 共用型アプリ (blog.64p.org の origin)。
#
# env の値は AppRun API が読み返さない (常に null/sensitive で返る) ため、
# Terraform 側からは drift 検知できない。値の正本は 1Password (blog4 vault)
# に置き、ここでは var 経由で push する (= ゼロからの再構築を可能にする DR 用途)。
#
# image は CI (.github/actions/deploy-apprun/deploy.sh) が毎デプロイで PATCH
# するため lifecycle.ignore_changes で管理外にする。traffics も CI 側。
#
# component name はイメージタグ入りの名前で、CI 側では変えていないが
# 過去のデプロイ経緯で実態とずれうる。name は forces replacement 属性なので、
# ずれたまま apply するとアプリごと作り直され public_url が変わる
# (= WebAccel の origin が切れて blog.64p.org が落ちる)。ここも管理外にする。
locals {
  # AppRun に流し込む環境変数。実値は 1Password 由来 (TF_VAR_* / op run)。
  apprun_env = {
    ADMIN_USER               = var.admin_user
    ADMIN_PW                 = var.admin_pw
    AMAZON_PAAPI5_ACCESS_KEY = var.amazon_paapi5_access_key
    AMAZON_PAAPI5_SECRET_KEY = var.amazon_paapi5_secret_key
    DATABASE_HOST            = var.database_host
    DATABASE_NAME            = var.database_name
    DATABASE_PASSWORD        = var.database_password
    DATABASE_PORT            = var.database_port
    DATABASE_USER            = var.database_user
    GIN_MODE                 = var.gin_mode
    S3_ACCESS_KEY_ID         = var.s3_access_key_id
    S3_SECRET_ACCESS_KEY     = var.s3_secret_access_key
    WEBACCEL_GUARD           = var.webaccel_guard
  }
}

resource "sakura_apprun_shared" "blog4" {
  name            = "blog4"
  port            = 8181
  timeout_seconds = 60
  min_scale       = 1
  max_scale       = 1

  components = [
    {
      # name / image は ignore_changes 対象 (実値は CI とデプロイ履歴が決める)。
      # ゼロから作り直すときの初期値として、現行の実態に合わせておく。
      name       = "blog4:3bfb434"
      max_cpu    = "0.5"
      max_memory = "1Gi"

      deploy_source = {
        container_registry = {
          # ghcr.io の public パッケージなので pull に認証は要らない。
          # server / username は optional かつ computed なので、指定しない
          # (CI も PATCH で username/password を削除している)。
          image = "ghcr.io/tokuhirom/blog4:3bfb434"
        }
      }

      env = [for k, v in local.apprun_env : { key = k, value = v }]
    },
  ]

  lifecycle {
    ignore_changes = [
      components[0].name,
      components[0].deploy_source.container_registry.image,
      traffics,
    ]
  }
}
