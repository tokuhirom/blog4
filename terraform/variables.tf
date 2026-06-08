variable "zone" {
  description = "Default Sakura Cloud zone. TiDB CR is in is1, so default to is1a."
  type        = string
  default     = "is1a"
}

variable "tidb_database_name" {
  description = "Database name to create on the TiDB CR instance."
  type        = string
  default     = "blog4"
}

variable "tidb_password" {
  description = "Password for the TiDB CR root user. Inject via TF_VAR_tidb_password (write-only)."
  type        = string
  sensitive   = true
}

variable "tidb_password_version" {
  description = "Bump this to rotate tidb_password without recreating the resource."
  type        = number
  default     = 1
}

# --- AppRun 環境変数 (実値は 1Password 由来。TF_VAR_* で注入) ---
# AppRun API は env 値を読み返さないため、ここに与えた値が apply 時に push される。
# import 時の挙動を no-op にしたい場合は「現在の本番値」を 1Password に入れること。

variable "admin_user" {
  description = "AppRun env ADMIN_USER (管理 UI のユーザ名)."
  type        = string
}

variable "admin_pw" {
  description = "AppRun env ADMIN_PW (管理 UI のパスワード)."
  type        = string
  sensitive   = true
}

variable "amazon_paapi5_access_key" {
  description = "AppRun env AMAZON_PAAPI5_ACCESS_KEY."
  type        = string
  sensitive   = true
}

variable "amazon_paapi5_secret_key" {
  description = "AppRun env AMAZON_PAAPI5_SECRET_KEY."
  type        = string
  sensitive   = true
}

variable "database_host" {
  description = "AppRun env DATABASE_HOST (接続先 DB ホスト)."
  type        = string
}

variable "database_name" {
  description = "AppRun env DATABASE_NAME."
  type        = string
}

variable "database_password" {
  description = "AppRun env DATABASE_PASSWORD."
  type        = string
  sensitive   = true
}

variable "database_port" {
  description = "AppRun env DATABASE_PORT (env は文字列で渡す)."
  type        = string
}

variable "database_user" {
  description = "AppRun env DATABASE_USER."
  type        = string
}

variable "gin_mode" {
  description = "AppRun env GIN_MODE (release など)."
  type        = string
}

variable "s3_access_key_id" {
  description = "AppRun env S3_ACCESS_KEY_ID (アプリの S3 アクセスキー)."
  type        = string
  sensitive   = true
}

variable "s3_secret_access_key" {
  description = "AppRun env S3_SECRET_ACCESS_KEY."
  type        = string
  sensitive   = true
}

variable "webaccel_guard" {
  description = "AppRun env WEBACCEL_GUARD (オリジン保護トークン)."
  type        = string
  sensitive   = true
}
