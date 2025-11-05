# ArgoCD Applications

ArgoCD Application manifests for deploying apps to the cluster.

## Prerequisites

1. Cluster is running
2. Infrastructure is installed (ingress-nginx)
3. ArgoCD is installed and accessible

## Installation Order

### 1. Install Infrastructure (manual)

```shell
# NGINX Ingress Controller
kubectl apply -f ../infrastructure/ingress-nginx/install.yaml

# Verify
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s
```

### 2. Install ArgoCD (manual bootstrap)
```shell
# Install ArgoCD
kubectl apply -k ../argocd/overlays/prod

# Wait for ArgoCD
kubectl wait --for=condition=available --timeout=600s \
  deployment/argocd-server -n argocd

# Get admin password
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath="{.data.password}" | base64 -d && echo

# Get ArgoCD URL
kubectl get svc argocd-server -n argocd \
  -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
```

### 3. Deploy Applications via ArgoCD

#### Option A: Via UI

1. Login to ArgoCD UI
2. Create applications manually using these manifests as reference

#### Option B: Via kubectl (GitOps way)

```shell
# Deploy all prod applications
kubectl apply -f backend-prod.yaml
kubectl apply -f frontend-prod.yaml
kubectl apply -f invoice-app-prod.yaml

# Or deploy all at once
kubectl apply -f . --recursive

# Check status
kubectl get applications -n argocd

# Watch sync progress
kubectl get applications -n argocd -w
```

## Application Structure

```shell
├── backend-prod.yaml          # Backend deployment
├── frontend-prod.yaml         # Frontend deployment
└── invoice-app-prod.yaml      # Shared ingress
```

## Sync Policy

All applications use:

- **Automated sync**: Changes in Git trigger automatic deployment
- **Self-heal**: Drift detection and auto-correction
- **Prune**: Removed resources in Git are deleted from cluster

## Troubleshooting

```shell
# Check application status
kubectl describe application backend-prod -n argocd

# Check sync status
argocd app get backend-prod

# Manual sync
argocd app sync backend-prod

# View logs
argocd app logs backend-prod --follow
```

## Development Environment

To create dev environment:

1. Copy `*-prod.yaml` to `*-dev.yaml`
2. Change `targetRevision`
3. Change `path` to `overlays/dev`
4. Add dev overlays
5. Apply: `kubectl apply -f *-dev.yaml`
