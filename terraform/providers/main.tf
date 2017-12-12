provider "azurerm" {
  /*
  $ export ARM_SUBSCRIPTION_ID = your_azure_sub_id
  $ export ARM_CLIENT_ID = your_client_id
  $ export ARM_CLIENT_SECRET = your_client_secret
  $ export AWS_DEFAULT_REGION=us-east-1
  */
  region      = "${var.region}"
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
}

resource "aws_key_pair" "kismatic" {
  key_name   = "${var.cluster_name}"
  public_key = "${file("${var.public_ssh_key_path}")}"
}


