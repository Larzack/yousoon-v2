# =============================================================================
# Secrets Templates - À créer manuellement ou via External Secrets
# =============================================================================
# Ces secrets doivent être créés AVANT le déploiement Helmfile
# 
# Commandes pour créer les secrets:
#
# kubectl create namespace yousoon-staging
# kubectl create namespace yousoon-production
#
# # MongoDB
# kubectl create secret generic mongodb-credentials \
#   --from-literal=mongodb-root-password=<PASSWORD> \
#   -n yousoon-<ENV>
#
# # Redis
# kubectl create secret generic redis-credentials \
#   --from-literal=password=<PASSWORD> \
#   -n yousoon-<ENV>
#
# # Grafana
# kubectl create secret generic grafana-credentials \
#   --from-literal=admin-user=admin \
#   --from-literal=admin-password=<PASSWORD> \
#   -n yousoon-<ENV>
#
# # JWT
# kubectl create secret generic jwt-secrets \
#   --from-literal=secret-key=<JWT_SECRET> \
#   --from-literal=refresh-secret=<JWT_REFRESH_SECRET> \
#   -n yousoon-<ENV>
#
# # External Services
# kubectl create secret generic external-services \
#   --from-literal=onesignal-app-id=<ONESIGNAL_APP_ID> \
#   --from-literal=onesignal-api-key=<ONESIGNAL_API_KEY> \
#   --from-literal=aws-access-key-id=<AWS_ACCESS_KEY_ID> \
#   --from-literal=aws-secret-access-key=<AWS_SECRET_ACCESS_KEY> \
#   -n yousoon-<ENV>
# =============================================================================
