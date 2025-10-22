resource "aws_ecr_repository" "frontend" {
  name = "${var.project_name}-frontend"
  image_tag_mutability = "MUTABLE"
  tags = var.tags
}

resource "aws_ecr_repository" "backend" {
  name = "${var.project_name}-backend"
  image_tag_mutability = "MUTABLE"
  tags = var.tags
}

