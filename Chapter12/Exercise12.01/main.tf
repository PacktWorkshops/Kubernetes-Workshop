resource "aws_s3_bucket" "my_bucket" {
  bucket = "zparnold-test-bucket"
  acl    = "private"