# 既存 AppRun 共用型アプリ (blog.64p.org の origin)。
#
# env の値は AppRun API が読み返さない (常に null/sensitive で返る) ため、
# Terraform 側からは drift 検知できない。値の正本は 1Password (blog4 vault)
# に置き、ここでは var 経由で push する (= ゼロからの再構築を可能にする DR 用途)。
#
# image は CI (.github/actions/deploy-apprun/deploy.sh) が毎デプロイで PATCH
# するため lifecycle.ignore_changes で管理外にする。traffics も CI 側。
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
      # CI は image しか PATCH しないため component name は安定 (タグ込みだが固定)。
      name       = "blog4:75e6b60"
      max_cpu    = "0.5"
      max_memory = "1Gi"

      deploy_source = {
        container_registry = {
          image    = "tokuhirom-private.sakuracr.jp/blog4:75e6b60"
          server   = "tokuhirom-private.sakuracr.jp"
          username = "pull"
        }
      }

      env = [for k, v in local.apprun_env : { key = k, value = v }]
    },
  ]

  lifecycle {
    ignore_changes = [
      components[0].deploy_source.container_registry.image,
      traffics,
    ]
  }
}
