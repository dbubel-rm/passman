
resource "aws_instance" "bastion" {
  ami                         = "ami-0ac019f4fcb7cb7e6"
  instance_type               = "t2.micro"
  key_name                    = "prod_key"
  monitoring                  = true
  vpc_security_group_ids      = ["${aws_security_group.bastion.id}", "${aws_security_group.rds_sg.id}"]
  subnet_id                   = "${element(aws_subnet.public.*.id, count.index)}"
  associate_public_ip_address = true

  tags {
    Name        = "bastion"
  }
}



resource "aws_security_group" "bastion" {
  vpc_id      = "${aws_vpc.main.id}"
  name        = "bastion-host-sg"
  description = "Allow SSH to bastion host"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

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
    Name        = "bastion-sg"
  }
}
