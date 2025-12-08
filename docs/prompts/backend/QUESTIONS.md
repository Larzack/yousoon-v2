# ❓ Questions - Backend Go + GraphQL

> Questions techniques et fonctionnelles à clarifier avant développement

---

## ✅ DÉCISIONS VALIDÉES

### Architecture
| Question | Réponse |
|----------|---------|
| Microservices vs Monolithe | **Microservices** (sauf si trop complexe → monolithe modulaire) |
| Communication inter-services | **gRPC** (sync) + **NATS JetStream** (async) |
| Schema Federation | **Oui** - Gateway agrège les schemas |
| GraphQL Subscriptions | **Oui** - Temps réel activé (WebSocket) |
| Persisted Queries | **Oui** - Pour performance et sécurité |

### Base de données
| Question | Réponse |
|----------|---------|
| MongoDB | **1 cluster MongoDB** avec **1 database par context** |
| MongoDB Hébergement | **Self-hosted sur EKS** |
| MongoDB HA | **Non pour commencer** (Standalone, HA plus tard) |
| Redis | **Standalone** (stockage refresh tokens) |

### JWT
| Question | Réponse |
|----------|---------|
| Génération JWT | **Identity Service** génère les tokens |
| Validation JWT | **Gateway** valide (middleware) |
| Access Token durée | **6 heures** |
| Refresh Token durée | **30 jours** |
| Stockage Refresh Token | **Redis** |

### Authentification
| Question | Réponse |
|----------|---------|
| SSO partenaires | **Non** - Pas de SSO externe |

### Vérification Identité (CNI)
| Question | Réponse |
|----------|---------|
| Méthode | OCR interne (Tesseract/OpenCV) |
| Documents acceptés | **Tous** (CNI, passeport, permis, étrangers) |
| Niveau vérification | **OCR simple** (lecture données) |
| Tentatives max | **10** avant blocage |

### Infrastructure Kubernetes
| Question | Réponse |
|----------|---------|
| Cloud | **AWS EKS** |
| Ingress Controller | **Nginx Ingress** |
| Secrets Management | **Kubernetes Secrets** |
| DNS | **Route53** |
| SSL Certificates | **Let's Encrypt** (cert-manager) |

### Observabilité
| Question | Réponse |
|----------|---------|
| Stack | OpenTelemetry + Jaeger + Prometheus + Loki + Grafana |
| Crash reporting | Sentry (self-hosted) |

### Cartographie
| Question | Réponse |
|----------|---------|
| Service cartes | **Google Maps** |

### Performance
| Question | Réponse |
|----------|---------|
| Objectif utilisateurs | **5000 min/heure** |
| Objectif requêtes | **Maximum possible** |

### SMS
| Question | Réponse |
|----------|----------|
| OTP inscription | **Oui** (SMS) |
| Rappels sortie | **Oui** (SMS) |

### Catégories d'intérêts (provisoires)
| Emoji | Catégorie |
|-------|----------|
| 🎾 | Sport |
| 🥂 | Convivialité |
| 🎨 | Arts & Culture |
| 🌎 | Voyage & Escapade |
| 🎺 | Concert & Musique |
| 🧘 | Bien-être |
| 🌱 | Nature |
| 🍴 | Gastronomie & Dégustation |
| ⛵ | Évasion |
| 💫 | Développement personnel |

> ⚠️ Liste non définitive - à affiner

---

## ❓ QUESTIONS EN SUSPENS

### Business

1. **Abonnements** :
   - Quels plans d'abonnement (noms, prix, limites) ?
   - Durée de la période d'essai gratuite ?

---

## 📝 Notes

*Espace pour noter les réponses et décisions*

| Question | Réponse | Date |
|----------|---------|------|
| | | |
