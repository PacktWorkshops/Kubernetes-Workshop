resource "aws_s3_bucket" "my_bucket" {
  bucket = "<<NAME>>-test-bucket"
  acl    = "private"
}
