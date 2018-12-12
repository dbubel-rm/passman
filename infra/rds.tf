resource "aws_db_instance" "passmandb" {
  allocated_storage    = 10
  storage_type         = "gp2"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t2.micro"
  name                 = "rename"
  username             = "passman"
  password             = "foobarbaz"
  parameter_group_name = "default.mysql5.7"
  publicly_accessible = false
  enabled_cloudwatch_logs_exports = ["general"]
  skip_final_snapshot = true
}