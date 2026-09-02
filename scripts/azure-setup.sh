#!/bin/bash
# One-time Azure resource setup for nicksite.
# Run this once locally (az login first), then add the printed values as GitHub secrets.

set -e

RG="nicksite-rg"
LOCATION="eastus"
ACR="nicksiteregistry"
APP="nicksite-app"
ENV="nicksite-env"

echo "==> Creating resource group..."
az group create --name "$RG" --location "$LOCATION"

echo "==> Creating Container Registry..."
az acr create --resource-group "$RG" --name "$ACR" --sku Basic --admin-enabled true

echo "==> Creating Container Apps environment..."
az containerapp env create --name "$ENV" --resource-group "$RG" --location "$LOCATION"

echo "==> Building and pushing initial image..."
az acr build --registry "$ACR" --image nicksite:latest .

echo "==> Creating Container App..."
az containerapp create \
  --name "$APP" \
  --resource-group "$RG" \
  --environment "$ENV" \
  --image "$ACR.azurecr.io/nicksite:latest" \
  --target-port 8080 \
  --ingress external \
  --min-replicas 0 \
  --max-replicas 2 \
  --cpu 0.25 \
  --memory 0.5Gi \
  --env-vars ENV=production

echo ""
echo "========================================"
echo "Now add these as GitHub Actions secrets:"
echo "========================================"
echo ""
echo "ACR_USERNAME:"
az acr credential show --name "$ACR" --query username -o tsv

echo ""
echo "ACR_PASSWORD:"
az acr credential show --name "$ACR" --query "passwords[0].value" -o tsv

echo ""
echo "AZURE_CREDENTIALS  (run this and paste the full JSON output):"
echo "  az ad sp create-for-rbac --name nicksite-deploy --role contributor \\"
echo "    --scopes /subscriptions/\$(az account show --query id -o tsv)/resourceGroups/$RG \\"
echo "    --sdk-auth"

echo ""
echo "ALLOWED_HOSTS:"
az containerapp show --name "$APP" --resource-group "$RG" \
  --query "properties.configuration.ingress.fqdn" -o tsv

echo ""
echo "Done. Push to main to trigger your first deploy."
