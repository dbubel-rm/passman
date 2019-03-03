# Create a cloudwatch log group
resource "aws_cloudwatch_log_group" "passman_server" {
  name = "${var.environment}-lg"

  tags {
    Environment = "${var.environment}"
    Application = "Passman API"
  }
}

# Create an ECR repository for the server container
resource "aws_ecr_repository" "passman_server_app" {
  name = "${var.environment}-ecr"
}

# Create the ECS cluster for our tasks
resource "aws_ecs_cluster" "cluster" {
  name = "${var.environment}-ecs-cluster"
}

# ECS task definition
# Task definition for the server
data "template_file" "passman_server_task" {
  template = "${file("${path.module}/tasks/passman_server_task_defintion.json")}"

  vars {
    image          = "${aws_ecr_repository.passman_server_app.repository_url}"
    mysql_endpoint = "${var.mysql_endpoint}"
    region         = "${var.region}"
    log_group      = "${aws_cloudwatch_log_group.passman_server.name}"
  }
}

resource "aws_ecs_task_definition" "passman_server" {
  family                   = "${var.environment}_passman-server"
  container_definitions    = "${data.template_file.passman_server_task.rendered}"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  depends_on               = ["aws_iam_role_policy.ecs_service_role_policy"]
  execution_role_arn       = "${aws_iam_role.ecs_execution_role.arn}"
  task_role_arn            = "${aws_iam_role.ecs_execution_role.arn}"
}

# App load balancers
resource "aws_alb_target_group" "alb_target_group" {
  name        = "passman-server-alb-tg"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = "${var.vpc_id}"
  target_type = "ip"

  lifecycle {
    create_before_destroy = true
  }

  health_check {
    path = "/v1/health"
  }
}

# Security group for ALB
resource "aws_security_group" "app_inbound_sg" {
  name        = "${var.environment}-inbound-sg"
  description = "Allow HTTP from anywhere into ALB"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8
    to_port     = 0
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "${var.environment}-app-inbound-sg"
  }
}

resource "aws_alb" "alb_passman_server" {
  name            = "passman-server-alb"
  subnets         = ["${var.public_subnet_ids}"]
  security_groups = ["${var.security_group_ids}", "${aws_security_group.app_inbound_sg.id}"]

  tags {
    Name        = "${var.environment}-alb"
    Environment = "${var.environment}"
  }
}

resource "aws_alb_listener" "passman_server_alb_listener" {
  load_balancer_arn = "${aws_alb.alb_passman_server.arn}"
  port              = "80"
  protocol          = "HTTP"
  depends_on        = ["aws_alb_target_group.alb_target_group"]

  default_action {
    target_group_arn = "${aws_alb_target_group.alb_target_group.arn}"
    type             = "forward"
  }
}

# ECS service
resource "aws_security_group" "ecs_service" {
  vpc_id      = "${var.vpc_id}"
  name        = "${var.environment}-ecs-service-sg"
  description = "Allow egress from container"

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8
    to_port     = 0
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "${var.environment}-ecs-service-sg"
    Environment = "${var.environment}"
  }
}

# Simply specify the family to find the latest ACTIVE revision in that family
data "aws_ecs_task_definition" "app" {
  task_definition = "${aws_ecs_task_definition.passman_server.family}"

  depends_on = ["aws_ecs_task_definition.passman_server"]
}

resource "aws_ecs_service" "app" {
  name            = "${var.environment}-app"
  task_definition = "${aws_ecs_task_definition.passman_server.family}:${max("${aws_ecs_task_definition.passman_server.revision}", "${data.aws_ecs_task_definition.app.revision}")}"
  desired_count   = 1
  launch_type     = "FARGATE"
  cluster         = "${aws_ecs_cluster.cluster.id}"

  network_configuration {
    security_groups = ["${var.security_group_ids}", "${aws_security_group.ecs_service.id}"]
    subnets         = ["${var.public_subnet_ids}"]
  }

  load_balancer {
    target_group_arn = "${aws_alb_target_group.alb_target_group.arn}"
    container_name   = "${var.environment}-app"
    container_port   = "80"
  }

  depends_on = ["aws_alb_target_group.alb_target_group"]
}

data "aws_iam_policy_document" "ecs_service_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ecs_role" {
  name               = "ecs_role"
  assume_role_policy = "${data.aws_iam_policy_document.ecs_service_role.json}"
}

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
resource "aws_iam_role_policy" "ecs_service_role_policy" {
  name = "ecs_service_role_policy"

  #policy = "${file("${path.module}/policies/ecs-service-role.json")}"
  policy = "${data.aws_iam_policy_document.ecs_service_policy.json}"
  role   = "${aws_iam_role.ecs_role.id}"
}

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
