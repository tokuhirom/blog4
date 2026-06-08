resource "sakura_ondemand_db" "blog4" {
  name          = "blog4"
  database_name = var.tidb_database_name
  database_type = "tidb"
  region        = "is1"
  description   = "blog4 primary database (TiDB CR)"
  tags          = ["blog4", "managed-by:terraform"]

  password_wo         = var.tidb_password
  password_wo_version = var.tidb_password_version
}
