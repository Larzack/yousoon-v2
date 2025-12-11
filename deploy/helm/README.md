# 📦 Yousoon Helm Charts

> Charts Helm personnalisés pour déployer la plateforme Yousoon  
> **Dernière mise à jour** : 11 décembre 2025

---

## 🎯 Vue d'Ensemble

Ce dossier contient les charts Helm personnalisés pour déployer l'ensemble de la plateforme Yousoon sur Kubernetes.

### Fonctionnalité Principale : Mode Sidecar vs Classic

Le déploiement supporte **deux modes** contrôlés par le paramètre `global.infra` :

| Mode | Description | Pods créés |
|------|-------------|------------|
| **sidecar** | Multi-container pods | 4 pods (storage, services, sites, monitoring) |
| **classic** | Un container par pod | ~20 pods (un par composant) |

---

## 📁 Structure des Charts

```
deploy/helm/
├── helmfile.yaml              # Orchestration principale
├── README.md                  # Cette documentation
├── charts/
│   ├── yousoon/               # Umbrella chart (parent)
│   │   ├── Chart.yaml
│   │   ├── values.yaml        # Configuration par défaut
│   │   ├── values-sidecar.yaml
│   │   ├── values-classic.yaml
│   │   └── values-production.yaml
│   │
│   ├── yousoon-storage/       # MongoDB, Redis, NATS, Elasticsearch
│   │   ├── Chart.yaml
│   │   ├── values.yaml
│   │   └── templates/
│   │       ├── _helpers.tpl
│   │       ├── statefulset-sidecar.yaml   # Mode sidecar
│   │       ├── statefulset-classic.yaml   # Mode classic
│   │       ├── service.yaml
│   │       ├── configmap.yaml
│   │       ├── secret.yaml
│   │       └── pvc.yaml
│   │
│   ├── yousoon-services/      # 6 microservices + Router
│   │   ├── Chart.yaml
│   │   ├── values.yaml
│   │   └── templates/
│   │       ├── deployment-sidecar.yaml
│   │       ├── deployment-classic.yaml
│   │       ├── service.yaml
│   │       └── secret.yaml
│   │
│   ├── yousoon-sites/         # Admin, Partners, Siteweb
│   │   ├── Chart.yaml
│   │   ├── values.yaml
│   │   └── templates/
│   │       ├── deployment-sidecar.yaml
│   │       ├── deployment-classic.yaml
│   │       ├── service.yaml
│   │       └── ingress.yaml
│   │
│   └── yousoon-monitoring/    # Prometheus, Grafana, Loki, Jaeger
│       ├── Chart.yaml
│       ├── values.yaml
│       └── templates/
│           ├── statefulset-sidecar.yaml
│           ├── deployment-classic.yaml
│           ├── service.yaml
│           ├── configmap.yaml
│           ├── rbac.yaml
│           └── daemonset-promtail.yaml
│
└── values/                    # Values pour charts Bitnami (legacy)
    ├── mongodb.yaml
    ├── redis.yaml
    └── ...
```

---

## 🚀 Déploiement

### Prérequis

```bash
# Installer Helm 3.x
brew install helm

# Installer Helmfile
brew install helmfile

# Installer le plugin helm-diff
helm plugin install https://github.com/databus23/helm-diff
```

### Mode Sidecar (Défaut)

```bash
cd deploy/helm

# Mettre à jour les dépendances
helm dependency update charts/yousoon

# Déployer en mode sidecar
helmfile sync

# Ou explicitement
helmfile sync --state-values-set global.infra=sidecar
```

### Mode Classic

```bash
cd deploy/helm

# Déployer en mode classic
helmfile sync --state-values-set global.infra=classic
```

### Environnement Production

```bash
cd deploy/helm

# Déployer en production
helmfile sync -e production

# Ou avec fichier values spécifique
helm install yousoon ./charts/yousoon \
  -f ./charts/yousoon/values-production.yaml \
  -n yousoon-production
```

---

## ⚙️ Configuration

### Paramètres Globaux

| Paramètre | Description | Valeurs | Défaut |
|-----------|-------------|---------|--------|
| `global.infra` | Mode de déploiement | `sidecar`, `classic` | `sidecar` |
| `global.namespace` | Namespace Kubernetes | string | `yousoon-staging` |
| `global.environment` | Environnement | `staging`, `production` | `staging` |
| `global.imageRegistry` | Registry Docker | string | ECR |

### Activer/Désactiver des Composants

```yaml
# values.yaml
storage:
  enabled: true      # MongoDB, Redis, NATS, Elasticsearch

services:
  enabled: true      # 6 microservices + Router

sites:
  enabled: true      # Admin, Partners, Siteweb

monitoring:
  enabled: true      # Prometheus, Grafana, Loki, Jaeger
```

### Configuration par Composant

```yaml
# Exemple: désactiver Elasticsearch
yousoon-storage:
  elasticsearch:
    enabled: false

# Exemple: augmenter les ressources MongoDB
yousoon-storage:
  mongodb:
    resources:
      requests:
        memory: "1Gi"
        cpu: "500m"
```

---

## 📊 Mode Sidecar - Détails

En mode **sidecar**, les composants sont groupés ainsi :

### Pod Storage (StatefulSet)

```
┌─────────────────────────────────────────────────────────────┐
│                    POD: yousoon-storage                     │
├──────────────┬──────────────┬──────────────┬───────────────┤
│   MongoDB    │    Redis     │     NATS     │ Elasticsearch │
│   :27017     │    :6379     │    :4222     │    :9200      │
└──────────────┴──────────────┴──────────────┴───────────────┘
                              │
            Volumes partagés: mongo-data, redis-data, nats-data, es-data
```

### Pod Services (Deployment)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           POD: yousoon-services                                  │
├──────────┬──────────┬───────────┬──────────┬────────────┬──────────────┬───────┤
│ Identity │ Partner  │ Discovery │ Booking  │ Engagement │ Notification │ Router│
│  :4001   │  :4002   │   :4003   │  :4004   │   :4005    │    :4006     │ :4000 │
└──────────┴──────────┴───────────┴──────────┴────────────┴──────────────┴───────┘
```

### Pod Sites (Deployment)

```
┌─────────────────────────────────────────────────────┐
│                POD: yousoon-sites                    │
├─────────────────┬─────────────────┬─────────────────┤
│      Admin      │    Partners     │    Siteweb      │
│      :3001      │      :3002      │      :3000      │
└─────────────────┴─────────────────┴─────────────────┘
```

### Pod Monitoring (StatefulSet)

```
┌─────────────────────────────────────────────────────────────┐
│                  POD: yousoon-monitoring                    │
├──────────────┬──────────────┬──────────────┬───────────────┤
│  Prometheus  │   Grafana    │     Loki     │    Jaeger     │
│    :9090     │    :3000     │    :3100     │    :16686     │
└──────────────┴──────────────┴──────────────┴───────────────┘
```

---

## 📊 Mode Classic - Détails

En mode **classic**, chaque composant a son propre pod :

```
Storage:           4 pods (mongodb, redis, nats, elasticsearch)
Services:          7 pods (identity, partner, discovery, booking, engagement, notification, router)
Sites:             3 pods (admin, partners, siteweb)
Monitoring:        4 pods (prometheus, grafana, loki, jaeger) + DaemonSet promtail
─────────────────────────
Total:            ~18+ pods
```

---

## 🔧 Commandes Utiles

### Vérifier le déploiement

```bash
# Voir les pods
kubectl get pods -n yousoon-staging

# Voir les services
kubectl get svc -n yousoon-staging

# Logs d'un pod (mode sidecar)
kubectl logs -n yousoon-staging yousoon-storage -c mongodb
kubectl logs -n yousoon-staging yousoon-services -c identity

# Logs (mode classic)
kubectl logs -n yousoon-staging -l app.kubernetes.io/name=mongodb
```

### Accéder aux services

```bash
# Port-forward Grafana
kubectl port-forward -n yousoon-staging svc/yousoon-monitoring-grafana 3000:3000

# Port-forward Admin (interne uniquement)
kubectl port-forward -n yousoon-staging svc/yousoon-sites-admin 3001:3001

# Port-forward MongoDB
kubectl port-forward -n yousoon-staging svc/yousoon-storage-mongodb 27017:27017
```

### Mise à jour

```bash
# Diff avant mise à jour
helmfile diff

# Appliquer les changements
helmfile sync

# Rollback
helm rollback yousoon 1 -n yousoon-staging
```

---

## 🔐 Secrets

Les secrets doivent être créés avant le déploiement :

```bash
# Créer le secret MongoDB
kubectl create secret generic yousoon-mongodb-secret \
  --from-literal=mongodb-root-password=<PASSWORD> \
  -n yousoon-staging

# Créer le secret Redis
kubectl create secret generic yousoon-redis-secret \
  --from-literal=redis-password=<PASSWORD> \
  -n yousoon-staging

# Ou utiliser le template
cp deploy/kubernetes/secrets.template.yaml deploy/kubernetes/secrets.yaml
# Éditer secrets.yaml avec vos valeurs
kubectl apply -f deploy/kubernetes/secrets.yaml
```

---

## ⚠️ Considérations

### Mode Sidecar

**Avantages** :
- Moins de pods à gérer
- Communication inter-container rapide (localhost)
- Consommation mémoire réduite

**Inconvénients** :
- Tous les containers partagent les ressources
- Redémarrage d'un container = redémarrage du pod
- Debugging plus complexe

### Mode Classic

**Avantages** :
- Isolation complète entre composants
- Scaling indépendant
- Pattern Kubernetes standard

**Inconvénients** :
- Plus de pods = plus de ressources
- Communication réseau entre pods

### Recommandations

| Environnement | Mode recommandé | Raison |
|---------------|-----------------|--------|
| **Développement** | sidecar | Économie de ressources |
| **Staging** | sidecar | Tests rapides |
| **Production** | classic | Isolation, scaling, résilience |

---

## 🔗 Références

- [Documentation Helm](https://helm.sh/docs/)
- [Helmfile](https://github.com/helmfile/helmfile)
- [Architecture DDD](../../docs/prompts/backend/ARCHITECTURE.md)
- [copilot-instructions.md](../../.github/copilot-instructions.md)
