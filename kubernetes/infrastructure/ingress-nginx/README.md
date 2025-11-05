# NGINX Ingress Controller

## Installation

This is cluster infrastructure and should be installed **before** ArgoCD applications.

### Version

- **Current**: v1.14.0
- **Installed**: 2025-11-05
- [**Source**](https://github.com/kubernetes/ingress-nginx)

### Install/Update

```bash
curl -o install.yaml \
  https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml

# Review changes
git diff install.yaml

# Apply
kubectl apply -f install.yaml

# Verify
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s

# Get Load Balancer URL
kubectl get svc -n ingress-nginx ingress-nginx-controller \
  -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
```

## Configuration

Default configuration is used. Custom configuration can be added via ConfigMap:

```bash
kubectl edit configmap ingress-nginx-controller -n ingress-nginx
```

## Monitoring

```bash
# Check controller logs
kubectl logs -n ingress-nginx -l app.kubernetes.io/component=controller

# Check ingresses
kubectl get ingress --all-namespaces
```
