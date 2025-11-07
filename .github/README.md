# GitHub Workflows Overview

This folder contains CI/CD workflows used to build, test, scan and deploy the project via GitHub Actions.

The code CD part is provided via deploied ArgoCD.

## What the workflows do (short)

- ci.yaml — detects changes, runs linters, build backend/frontend, versions code, push tags and create build artifacts. Also supports caching for faster pipelines.
- security-scans.yaml — run container, code security and terraform scans (Trivy).
- terraform-validate.yaml / terraform-apply.yaml — validate and apply infrastructure as code.
- Supporting composite actions: aws-setup, backend-lint, docker-build-and-push, terraform-init-validate-plan, versioning.

## Technologies used

- GitHub Actions (workflows & composite actions)
- Go + Gin (backend build)
- Node / Vite / Vue 3 + Vuetify (frontend build)
- Docker / Docker Buildx (image build & push)
- Trivy (container scanning)
- Terraform (infrastructure provisioning)
- Kubernetes / ArgoCD (optional deployment targets)
- Caching (go_sum, buildx cache, package-json.lock)

## Standards enforced

- Semantic versioning for releases (versioning action)
- Linting and unit tests on PRs
- Automated container and IaC security scanning
- Terraform validation and plan before apply
- Central `.versions.json` for tool versions management

## What is required for full operation

- GitHub repository settings: Workflows enabled, secrets allowed for actions.
- Required secrets:
  - AWS_ACCESS_KEY_ID
  - AWS_SECRET_ACCESS_KEY
  - ECR_REPO
  - AWS_REGION
- Ensure branch naming matches workflow triggers (main / release branches).
- Local environment variables used in Docker builds: GO_VERSION, NODE_VERSION (can be set in workflow via `.versions.json`).
