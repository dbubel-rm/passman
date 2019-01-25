output "cluster_name" {
    value = "${aws_ecs_cluster.cluster.name}"
}

output "service_name" {
    value = "${aws_ecs_service.app.name}"
}

output "alb_dns_name" {
    value = "${aws_alb.alb_passman_server.dns_name}"
}

output "alb_zone_id" {
    value = "${aws_alb.alb_passman_server.zone_id}"
}