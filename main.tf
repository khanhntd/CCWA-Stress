#####################################################################
# Ensure there is unique testing_id for each test
#####################################################################
resource "random_id" "testing_id" {
  byte_length = 8
}

##########################################
# Template Files
##########################################

locals {
  cwagent_config         = fileexists("./resources/amazon_cloudwatch_agent.json")
  cwagent_ecs_taskdef    = fileexists("./resources/statsd.sh")
}

#####################################################################
# Generate EC2 Key Pair for log in access to EC2
#####################################################################

resource "tls_private_key" "ssh_key" {
  count     = var.ssh_key_name == "" ? 1 : 0
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "aws_ssh_key" {
  count      = var.ssh_key_name == "" ? 1 : 0
  key_name   = "ec2-key-pair-${random_id.testing_id.hex}"
  public_key = tls_private_key.ssh_key[0].public_key_openssh
}

locals {
  ssh_key_name        = aws_key_pair.aws_ssh_key[0].key_name
  private_key_content = tls_private_key.ssh_key[0].private_key_pem
}

#####################################################################
# Generate EC2 Instance and execute test commands
#####################################################################
resource "aws_instance" "cwagent" {
  ami                         = data.aws_ami.latest.id
  instance_type               = var.ec2_instance_type
  key_name                    = local.ssh_key_name
  iam_instance_profile        = aws_iam_instance_profile.cwagent_instance_profile.name
  vpc_security_group_ids      = [aws_security_group.ec2_security_group.id]
  associate_public_ip_address = true

  tags = {
    Name = "cwagent-integ-test-ec2-${random_id.testing_id.hex}"
  }
}

resource "null_resource" "integration_test" {
  # Prepare Integration Test
  provisioner "remote-exec" {
    inline = [
        "sudo rpm -Uvh https://s3.amazonaws.com/amazoncloudwatch-agent/amazon_linux/amd64/latest/amazon-cloudwatch-agent.rpm",

    ]

    connection {
      type        = "ssh"
      user        = var.user
      private_key = local.private_key_content
      host        = aws_instance.cwagent.public_ip
    }
  }

  #Run sanity check and integration test
  provisioner "remote-exec" {
    inline = [

    ]
    connection {
      type        = "ssh"
      user        = "ec2-user"
      private_key = local.private_key_content
      host        = aws_instance.cwagent.public_ip
    }
  }

  depends_on = [aws_instance.cwagent]
}

data "aws_ami" "latest" {
  most_recent = true
  owners      = ["self", "506463145083"]

  filter {
    name   = "name"
    values = ["cloudwatch-agent-integration-test-al2*"]
  }
}
