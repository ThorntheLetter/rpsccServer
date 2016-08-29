provider "aws" {}

resource "aws_instance" "rpscc-host" {
    ami                         = "ami-dbc380bb"
    availability_zone           = "us-west-1c"
    instance_type               = "t2.small"
    associate_public_ip_address = true
    disable_api_termination     = false
    ebs_optimized               = false
    monitoring                  = false
    key_name                    = "seavey_generic"
    subnet_id                   = "subnet-275e707e"
    vpc_security_group_ids      = ["sg-a53473c1"]
    source_dest_check           = true
    iam_instance_profile        = "jenkins"
    user_data                   = "${file("./user-data.sh")}"
    root_block_device {
        volume_type           = "gp2"
        volume_size           = 75
        delete_on_termination = true
    }
    tags {
        "Name" = "RPSCC-EC2"
    }
}
