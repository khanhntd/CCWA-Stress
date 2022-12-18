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
  cwagent_config = "./resources/amazon-cloudwatch-agent.json"
  statsd         = "./resources/statsd.sh"
}

#####################################################################
# Generate EC2 Key Pair for log in access to EC2
#####################################################################

resource "tls_private_key" "ssh_key" {
  count     =  1
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "aws_ssh_key" {
  count      = 1 
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

}

resource "null_resource" "integration_test" {
  provisioner "file" {
    source      = local.cwagent_config
    destination = "/home/ec2-user/amazon_cloudwatch_agent.json"
    connection {
      type        = "ssh"
      user        = "ec2-user"
      private_key = local.private_key_content
      host        = aws_instance.cwagent.public_ip
    }
  }

  

  #Run sanity check and integration test
  provisioner "remote-exec" {
    inline = [
        "sudo /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl -a fetch-config -m ec2 -s -c file:/home/ec2-user/amazon_cloudwatch_agent.json",
        "export AWS_REGION=${var.region}",
        "rm -rf CCWA-Stress",
        "git clone https://github.com/khanhntd/CCWA-Stress.git",
        "cd CCWA-Stress",
        "go run  ./resources/main.go" 
        
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