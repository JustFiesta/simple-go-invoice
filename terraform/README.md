# Invoice App infrastructure

Infra is managed with Terraform with very simple GitOps approach - not ideal but sufficient for simple project.

## Resources

User needs to create S3 Bucket for backend first and paste it in `backend.tf`.

Then VPC, ECR, EKS are created.
