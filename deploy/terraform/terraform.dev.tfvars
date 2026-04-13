environment             = "dev"
aws_region              = "us-east-1"
debug_mode              = true
deployment_package_type = "image"
image_tag               = "latest"
api_base_path           = "/hello"
api_source_label        = "terraform-dev"

# Optional Aurora PostgreSQL (Data API) backing service.
enable_postgres        = false
postgres_database_name = "app"
