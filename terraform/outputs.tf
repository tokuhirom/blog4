output "tidb_hostname" {
  description = "Hostname to connect to the TiDB CR instance (port 4000)."
  value       = sakura_ondemand_db.blog4.hostname
}

output "tidb_database_name" {
  description = "Database name created on the TiDB CR instance."
  value       = sakura_ondemand_db.blog4.database_name
}

output "tidb_max_connections" {
  description = "Max connections allowed by the TiDB CR instance."
  value       = sakura_ondemand_db.blog4.max_connections
}
