output "alb_hostname" {
  value = "${aws_alb.main.dns_name}"
}

output "rds_address" {
  value = "${aws_db_instance.rds.address}"
}