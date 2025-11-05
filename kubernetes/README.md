# Kubernetes Configuration

GitOps repository for Invoice Management System running on AWS EKS.

## 📁 Directory Structure

```shell
kubernetes/
├── infrastructure/         # Cluster infrastructure (manual install)
│   └── ingress-nginx/     # NGINX Ingress Controller
│
├── argocd/                # ArgoCD installation (bootstrap)
│   ├── base/
│   └── overlays/prod/
│
├── apps/                  # Application manifests (managed by ArgoCD)
│   ├── backend/           # Go backend API
│   ├── frontend/          # Vue.js frontend
│   └── invoice-app/       # Shared ingress
│
└── argocd-apps/          # ArgoCD Application CRDs
    ├── backend-prod.yaml
    ├── frontend-prod.yaml
    └── invoice-app-prod.yaml
```

## 🚀 Installation Order

### 1. Infrastructure (One-time setup)

Infrastructure components are installed **manually** and are **prerequisite** for ArgoCD.

```shell
# Configure kubectl
aws eks update-kubeconfig --name <cluster-name> --region eu-west-1

# Install NGINX Ingress Controller
kubectl apply -f infrastructure/ingress-nginx/install.yaml

# Verify
kubectl get pods -n ingress-nginx
kubectl get svc -n ingress-nginx
```

### 2. ArgoCD (Bootstrap)

ArgoCD manages itself and all applications after initial bootstrap.

```shell
# Install ArgoCD
kubectl apply -k argocd/overlays/prod

# Wait for ready
kubectl wait --for=condition=available --timeout=600s \
  deployment/argocd-server -n argocd

# Get credentials
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath="{.data.password}" | base64 -d && echo

# Get URL
kubectl get svc argocd-server -n argocd
```

### 3. Applications (GitOps)

Applications are deployed and managed by ArgoCD.

```shell
# Deploy all applications
kubectl apply -f argocd-apps/

# Or deploy individually
kubectl apply -f argocd-apps/backend-prod.yaml
kubectl apply -f argocd-apps/frontend-prod.yaml
kubectl apply -f argocd-apps/invoice-app-prod.yaml

# Check status
kubectl get applications -n argocd
```

## 🔄 GitOps Workflow

1. **Make changes** to manifests in `apps/*/overlays/prod/`
2. **Commit & push** to Git repository
3. **ArgoCD detects** changes automatically
4. **ArgoCD syncs** changes to cluster
5. **Monitor** via ArgoCD UI or CLI

## 📊 Monitoring

```shell
# Check all resources
kubectl get all -n prod

# Check ingress
kubectl get ingress -n prod

# Check application health
kubectl get applications -n argocd

# View logs
kubectl logs -n prod -l app=backend --tail=50
kubectl logs -n prod -l app=frontend --tail=50
```

## 🔧 Troubleshooting

### Application not syncing

```shell
# Check ArgoCD application
kubectl describe application backend-prod -n argocd

# Manual sync
argocd app sync backend-prod

# View diff
argocd app diff backend-prod
```

### Pod not starting

```shell
# Describe pod
kubectl describe pod -n prod -l app=backend

# Check events
kubectl get events -n prod --sort-by='.lastTimestamp'

# View logs
kubectl logs -n prod -l app=backend --previous
```

### Ingress issues

```shell
# Check ingress controller
kubectl logs -n ingress-nginx -l app.kubernetes.io/component=controller

# Test backend connectivity
kubectl run test -n prod --rm -it --image=curlimages/curl --restart=Never -- \
  curl http://backend:8080/health

# Test frontend
kubectl run test -n prod --rm -it --image=curlimages/curl --restart=Never -- \
  curl -I http://frontend:80/
```

## 🌍 Environments

- **prod**: Production environment (namespace: `prod`)
- **dev**: Development environment (namespace: `dev`) - TBD

## 📚 Additional Resources

- [ArgoCD Documentation](https://argo-cd.readthedocs.io/)
- [Kustomize Documentation](https://kustomize.io/)
- [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
