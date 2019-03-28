resource "aws_instance" "bastion" {
  ami                         = "ami-0ac019f4fcb7cb7e6"
  instance_type               = "t2.micro"
  key_name                    = "bastion_key"
  monitoring                  = true
  vpc_security_group_ids      = ["${var.security_group_ids}"]
  subnet_id                   = "${element(var.public_subnet_ids, 0)}"
  associate_public_ip_address = true

  tags {
    Name        = "bastion"
  }
}
