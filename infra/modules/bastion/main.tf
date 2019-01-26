resource "aws_instance" "bastion" {
  ami                         = "${lookup(var.bastion_ami, var.region)}"
  instance_type               = "t2.micro"
  key_name                    = "${var.key_name}"
  monitoring                  = true
  vpc_security_group_ids      = ["${var.security_group_ids}"]
  subnet_id                   = "${element(var.public_subnet_ids, 0)}"
  associate_public_ip_address = true

  tags {
    Name        = "${var.environment}-bastion"
    Environment = "${var.environment}"
  }
}
