resource "aws_ecr_repository" "this" {
  name = "${var.project_name}"
  image_tag_mutability = "MUTABLE"
  tags = var.tags
}
