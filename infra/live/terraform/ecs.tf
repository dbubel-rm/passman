resource "aws_ecs_cluster" "main" {
  name = "cb-cluster"
}

data "template_file" "cb_app" {
  template = "${file("templates/ecs/cb_app.json.tpl")}"

  vars {
    app_image      = "${var.app_image}"
    fargate_cpu    = "${var.fargate_cpu}"
    fargate_memory = "${var.fargate_memory}"
    aws_region     = "${var.aws_region}"
    app_port       = "${var.app_port}"
    db_url         = "${aws_db_instance.rds.address}"
  }
}

resource "aws_ecs_task_definition" "app" {
  family                   = "cb-app-task"
  execution_role_arn       = "${aws_iam_role.ecs_execution_role.arn}"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.fargate_cpu}"
  memory                   = "${var.fargate_memory}"
  container_definitions    = "${data.template_file.cb_app.rendered}"
}

resource "aws_ecs_service" "main" {
  name            = "cb-service"
  cluster         = "${aws_ecs_cluster.main.id}"
  task_definition = "${aws_ecs_task_definition.app.arn}"
  desired_count   = "${var.app_count}"
  launch_type     = "FARGATE"

  network_configuration {
    security_groups  = ["${aws_security_group.ecs_tasks.id}"]
    subnets          = ["${aws_subnet.private.*.id}"]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = "${aws_alb_target_group.app.id}"
    container_name   = "cb-app"
    container_port   = "${var.app_port}"
  }

  depends_on = [
    "aws_alb_listener.front_end",
  ]
}




# data "aws_iam_policy_document" "ecs_service_role" {
#   statement {
#     effect  = "Allow"
#     actions = ["sts:AssumeRole"]

#     principals {
#       type        = "Service"
#       identifiers = ["ecs.amazonaws.com"]
#     }
#   }
# }

# resource "aws_iam_role" "ecs_role" {
#   name               = "ecs_role"
#   assume_role_policy = "${data.aws_iam_policy_document.ecs_service_role.json}"
# }

data "aws_iam_policy_document" "ecs_service_policy" {
  statement {
    effect    = "Allow"
    resources = ["*"]

    actions = [
      "elasticloadbalancing:Describe*",
      "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
      "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
      "ec2:Describe*",
      "ec2:AuthorizeSecurityGroupIngress",
    ]
  }
}

/* ecs service scheduler role */
# resource "aws_iam_role_policy" "ecs_service_role_policy" {
#   name = "ecs_service_role_policy"

#   #policy = "${file("${path.module}/policies/ecs-service-role.json")}"
#   policy = "${data.aws_iam_policy_document.ecs_service_policy.json}"
#   role   = "${aws_iam_role.ecs_role.id}"
# }

/* role that the Amazon ECS container agent and the Docker daemon can assume */
resource "aws_iam_role" "ecs_execution_role" {
  name               = "ecs_task_execution_role"
  assume_role_policy = "${file("${path.module}/policies/ecs-task-execution-role.json")}"
}

resource "aws_iam_role_policy" "ecs_execution_role_policy" {
  name   = "ecs_execution_role_policy"
  policy = "${file("${path.module}/policies/ecs-execution-role-policy.json")}"
  role   = "${aws_iam_role.ecs_execution_role.id}"
}
