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
