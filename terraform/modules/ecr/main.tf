resource "aws_ecr_repository" "this" {
  name = "${var.project_name}"
  image_tag_mutability = "MUTABLE"
  tags = var.tags
}

resource "aws_ecr_lifecycle_policy" "backend_policy" {
  repository = aws_ecr_repository.this.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last 4 tagged images only (multi-arch safe)"
        selection = {
          tagStatus   = "tagged"
          countType   = "imageCountMoreThan"
          countNumber = 5
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 2
        description  = "Expire untagged images (architektur pomocnicze)"
        selection = {
          tagStatus = "untagged"
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}
