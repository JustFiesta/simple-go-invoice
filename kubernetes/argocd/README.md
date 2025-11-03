# ArgoCD Configuration manifests

These are deployed manually due to lacking need of frequent updates.

It will create AWS LB so terraform might be stuck sometimes when recreating VPC parts.

## To update Argo Version

```shell
# Fetch new base manifest
curl -o kubernetes/argocd/base/install.yaml \
  https://raw.githubusercontent.com/argoproj/argo-cd/v2.13.0/manifests/install.yaml

# Review changes
git diff kubernetes/argocd/base/install.yaml

# Commit if good
git add kubernetes/argocd/
git commit -m "Some commit message"

# Apply update
kubectl apply -k kubernetes/argocd/overlays/prod

# Check rollout
kubectl rollout status deployment/argocd-server -n argocd
```

## First time install

```shell
# Configure kubectl
aws eks update-kubeconfig --name <cluster-name> --region <region>

# Apply kustomization
kubectl apply -k kubernetes/argocd/overlays/prod

# Wait for argo to come up
kubectl wait --for=condition=available --timeout=600s \
  deployment/argocd-server -n argocd

# Wait for AWS LB to come up online

# Fetch admin password
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath="{.data.password}" | base64 -d && echo

# Fetch LB URL
kubectl get svc argocd-server -n argocd \
  -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
```
