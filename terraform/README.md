# Infrastructure

Infrastructure as Code for the invoice app, managed with Terraform and a simple GitOps workflow.

## What's Inside

This setup provisions AWS infrastructure for a Kubernetes-based application:

- **VPC** with public/private subnets across 2 AZs
- **EKS cluster** (v1.34) with Auto Mode compute
- **ECR repository** for container images
- **Essential EKS addons**: VPC CNI, CoreDNS, kube-proxy, EBS CSI driver

## Structure

```shell
terraform/
├── environments/prod/    # Production environment config
│   ├── main.tf           # Module composition
│   ├── providers.tf      # AWS provider & S3 backend
│   └── outputs.tf        # Output values
└── modules/              # Reusable modules
    ├── vpc/              # Network infrastructure
    ├── eks/              # Kubernetes cluster
    └── ecr/              # Container registry
```

## Technologies used

- **Terraform** 1.9.5
- **AWS Provider** ~> 6.0
- **Community modules**: `terraform-aws-modules/vpc` and `terraform-aws-modules/eks`
- **GitHub Actions** for CI/CD automation

## GitOps Workflow

### Pull Requests

When you open a PR touching Terraform files:

- Auto-formats code with `terraform fmt` and pushes it to VCS
- Runs init, validate, and plan
- Posts plan output as PR comment

### Main Branch

Merging to `main` automatically:

- Runs the same checks
- Applies changes with `terraform apply -auto-approve`

Workflows live in `.github/workflows/`:

- `terraform-validate.yaml` (PRs)
- `terraform-apply.yaml` (main branch)

## Getting Started

### Prerequisites

1. Create an S3 bucket for state storage
2. Update `backend.tf` with your bucket name
3. Set GitHub secrets:
   - `AWS_ACCESS_KEY_ID`
   - `AWS_SECRET_ACCESS_KEY`

### Local Usage

```bash
cd terraform/environments/prod

# Create new S3 bucket
# Change bucket in providers.tf

# Initialize
terraform init

# Plan changes
terraform plan -out=tfplan

# Apply changes
terraform apply tfplan
```

## Module Configuration

### VPC Module

- CIDR: `10.1.0.0/16`
- 2 private subnets, 2 public subnets
- Single NAT gateway (cost optimization)
- Proper tagging for EKS integration

### EKS Module

- Auto Mode with general-purpose node pools
- Public API endpoint enabled
- IAM roles with permission boundaries
- Pod Identity for EBS CSI driver (no IRSA)

### ECR Module

- Lifecycle policy: keeps last 5 tagged images
- Expires untagged images after 1 day
- Mutable tags enabled

## Notes

- State is stored remotely in S3 (see `providers.tf`) - change for Your bucket if needed
- IAM roles use `DefaultBoundaryPolicy` permission boundary - needed in my AWS (can be omitted in other)
- Auto-approval on main branch means changes deploy immediately—use with caution

## Outputs

After applying, you'll get:

- `vpc_id` - VPC identifier
- `cluster_endpoint` - EKS cluster ARN
- `repo_url` - ECR repository URL for pushing images
