terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "eu-west-1"
}

resource "aws_ssm_parameter" "bad" {
  name  = "/slash/at/end/"
  type  = "String"
  value = "test"
}

resource "aws_ssm_parameter" "good" {
  name  = "/no/slash/at/end"
  type  = "String"
  value = "test"
}

resource "aws_ssm_parameter" "bad2" {
  count = 3
  name  = "/slash/at/end${count.index}/"
  type  = "String"
  value = "test"
}

resource "aws_ssm_parameter" "good2" {
  count = 3
  name  = "/no/slash/at/end${count.index}"
  type  = "String"
  value = "test"
}
