# 🚀 Plan de Génération Global - Yousoon Platform

> **Projet** : Plateforme Yousoon (App + Backend + Sites)  
> **Architecture** : Monorepo avec microservices DDD + Apollo Federation 2  
> **Dernière mise à jour** : 9 décembre 2025

---

## 📋 Vue d'Ensemble avec Parallélisation

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                    YOUSOON - PLAN DE GÉNÉRATION PARALLÉLISÉ                              │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                          │
│  STREAM A (Backend)                    STREAM B (Frontend - Parallèle)                  │
│  ══════════════════                    ════════════════════════════════                 │
│                                                                                          │
│  ┌─────────────────────────┐           ┌─────────────────────────┐                      │
│  │ Phase 1: Infrastructure │           │ Phase 7: Site Vitrine   │ ← Aucune dépendance │
│  │ (~13h)                  │           │ (~12h)                  │                      │
│  └───────────┬─────────────┘           └─────────────────────────┘                      │
│              │                                                                           │
│              ▼                                                                           │
│  ┌─────────────────────────────────────────────────────────────────┐                    │
│  │     Phase 2: Core Subgraphs (~18h) - PARALLÉLISABLE            │                    │
│  │  ┌──────────────┬──────────────┬──────────────┐                │                    │
│  │  │   Identity   │   Partner    │  Discovery   │ ← En parallèle │                    │
│  │  │   (5h)       │   (4h)       │   (5h)       │                │                    │
│  │  └──────────────┴──────────────┴──────────────┘                │                    │
│  └───────────┬─────────────────────────────────────────────────────┘                    │
│              │                                                                           │
│              ▼                                                                           │
│  ┌─────────────────────────────────────────────────────────────────┐                    │
│  │     Phase 3: Business Subgraphs (~11h) - PARALLÉLISABLE        │                    │
│  │  ┌──────────────┬──────────────┬──────────────┐                │                    │
│  │  │   Booking    │  Engagement  │ Notification │ ← En parallèle │                    │
│  │  │   (4h)       │   (4h)       │   (3h)       │                │                    │
│  │  └──────────────┴──────────────┴──────────────┘                │                    │
│  └───────────┬─────────────────────────────────────────────────────┘                    │
│              │                                                                           │
│              │ Backend prêt → Déblocage Frontend                                        │
│              ▼                                                                           │
│  ┌─────────────────────────────────────────────────────────────────────────────────┐    │
│  │              PARALLÉLISATION FRONTEND (3 streams)                               │    │
│  │  ┌────────────────────┐  ┌────────────────────┐  ┌────────────────────┐        │    │
│  │  │  Phase 4: Mobile   │  │ Phase 5: Partner   │  │ Phase 6: Admin     │        │    │
│  │  │  Flutter (~46h)    │  │  Portal (~31h)     │  │ Backoffice (~26h)  │        │    │
│  │  │                    │  │                    │  │                    │        │    │
│  │  │  4.1 Core (8h)     │  │  Setup (4h)        │  │  Setup (3h)        │        │    │
│  │  │  4.2 Auth (12h)    │  │  Auth+Layout (6h)  │  │  Auth+Layout (4h)  │        │    │
│  │  │  4.3 Core (18h)    │  │  Offers (10h)      │  │  Users/Partners    │        │    │
│  │  │  4.4 Social (8h)   │  │  Analytics (6h)    │  │   (10h)            │        │    │
│  │  │                    │  │  Settings (5h)     │  │  Modération (5h)   │        │    │
│  │  │                    │  │                    │  │  Config (4h)       │        │    │
│  │  └────────────────────┘  └────────────────────┘  └────────────────────┘        │    │
│  └─────────────────────────────────────────────────────────────────────────────────┘    │
│                                           │                                              │
│                                           ▼                                              │
│  ┌─────────────────────────────────────────────────────────────────────────────────┐    │
│  │                        Phase 8: Tests & Déploiement (~15h)                      │    │
│  └─────────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                          │
│  TEMPS SÉQUENTIEL: ~167h    │    TEMPS PARALLÉLISÉ: ~89h (~2.5 semaines)               │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 🔀 Matrice de Parallélisation

### Ce qui peut être parallélisé

| Phase | Étapes Parallélisables | Économie |
|-------|------------------------|----------|
| **Phase 1** | 1.2 MongoDB, 1.3 Redis, 1.4 NATS (après 1.1) | ~4h → ~2h |
| **Phase 2** | Identity, Partner, Discovery (après Phase 1) | ~14h → ~5h |
| **Phase 3** | Booking, Engagement, Notification | ~11h → ~4h |
| **Phase 4-6** | Mobile, Partner Portal, Admin (3 streams) | ~103h → ~46h |
| **Phase 7** | Site Vitrine (indépendant, dès le début) | Parallèle total |

### Graphe de Dépendances

```
                                    ┌──────────────────┐
                                    │ SITE VITRINE (7) │ ← Peut démarrer immédiatement
                                    │    ~12h          │
                                    └──────────────────┘
                                    
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              CHEMIN CRITIQUE                                     │
└─────────────────────────────────────────────────────────────────────────────────┘
                                           
     ┌─────────┐     ┌─────────────────┐     ┌─────────────────┐
     │ Phase 1 │────►│    Phase 2      │────►│    Phase 3      │
     │  ~13h   │     │ (// subgraphs)  │     │ (// subgraphs)  │
     │         │     │    ~5h          │     │    ~4h          │
     └─────────┘     └────────┬────────┘     └────────┬────────┘
                              │                       │
                              └───────────┬───────────┘
                                          │
                    ┌─────────────────────▼─────────────────────┐
                    │         FRONTEND PARALLÈLE                │
                    │                                           │
                    │  ┌─────────┐ ┌─────────┐ ┌─────────┐     │
                    │  │ Mobile  │ │ Partner │ │  Admin  │     │
                    │  │  ~46h   │ │  ~31h   │ │  ~26h   │     │
                    │  └─────────┘ └─────────┘ └─────────┘     │
                    │                                           │
                    │  Temps réel (max): ~46h                   │
                    └─────────────────────────────────────────────┘
                                          │
                                          ▼
                              ┌─────────────────────┐
                              │ Phase 8: Deploy     │
                              │ ~15h                │
                              └─────────────────────┘

TEMPS TOTAL OPTIMISÉ: 13h + 5h + 4h + 46h + 15h = ~83h
```

---

## 🎯 Ordre de Priorité (Optimisé)

| Stream | Priorité | Module | Durée | Dépendances | Parallèle avec |
|--------|----------|--------|-------|-------------|----------------|
| A | 🔴 **1** | Backend Phase 1 | 13h | Aucune | Stream B |
| B | 🟢 **1** | Site Vitrine | 12h | Aucune | Stream A |
| A | 🔴 **2** | Core Subgraphs (//3) | 5h | Phase 1 | - |
| A | 🔴 **3** | Business Subgraphs (//3) | 4h | Phase 2 | - |
| C | 🟠 **4** | App Mobile | 46h | Phase 3 | Streams D, E |
| D | 🟡 **4** | Site Partenaires | 31h | Phase 3 | Streams C, E |
| E | 🟡 **4** | Admin Backoffice | 26h | Phase 3 | Streams C, D |
| - | 🔵 **5** | Tests & Déploiement | 15h | Tout | - |

---

## 📦 Structure Monorepo

```
yousoon/
├── .github/
│   ├── copilot-instructions.md
│   └── workflows/
│       ├── backend.yml
│       ├── mobile.yml
│       ├── partner-portal.yml
│       ├── vitrine.yml
│       └── admin.yml
│
├── docs/
│   └── prompts/
│       ├── DATA_MODEL.md
│       ├── DESIGN_SYSTEM.md
│       ├── backend/
│       │   ├── ARCHITECTURE.md
│       │   ├── PROMPT.md
│       │   └── GENERATION_PLAN.md
│       ├── app-mobile/
│       │   ├── PROMPT.md
│       │   ├── COMPONENTS.md
│       │   └── GENERATION_PLAN.md
│       ├── site-partenaires/
│       │   ├── PROMPT.md
│       │   └── GENERATION_PLAN.md
│       ├── site-vitrine/
│       │   ├── PROMPT.md
│       │   └── GENERATION_PLAN.md
│       └── admin/
│           ├── PROMPT.md
│           └── GENERATION_PLAN.md
│
├── apps/
│   ├── mobile/                     # Flutter App
│   ├── partner-portal/             # React + Vite
│   ├── vitrine/                    # Next.js 14
│   └── admin/                      # React + Vite
│
├── services/
│   ├── shared/                     # Domain commun
│   ├── gateway/                    # GraphQL Gateway
│   ├── identity-service/
│   ├── partner-service/
│   ├── discovery-service/
│   ├── booking-service/
│   ├── engagement-service/
│   └── notification-service/
│
├── deploy/
│   ├── kubernetes/
│   │   ├── base/
│   │   ├── services/
│   │   ├── infrastructure/
│   │   └── monitoring/
│   └── docker-compose.yml          # Dev local
│
├── packages/                       # Shared packages (si besoin)
│   └── ui-kit/                     # Composants partagés web
│
└── tools/
    ├── scripts/
    └── generators/
```

---

## 🔷 Phase 1 : Backend Infrastructure (~9h)

### Objectif
Mettre en place les fondations communes à tous les microservices.

### Étapes
| # | Tâche | Durée | Fichiers clés |
|---|-------|-------|---------------|
| 1.1 | Shared Domain (Aggregates, VOs, Events) | 2h | `services/shared/domain/` |
| 1.2 | Infrastructure MongoDB | 1h | `services/shared/infrastructure/mongodb/` |
| 1.3 | Infrastructure Redis | 1h | `services/shared/infrastructure/redis/` |
| 1.4 | Infrastructure NATS | 1h | `services/shared/infrastructure/nats/` |
| 1.5 | GraphQL Gateway | 4h | `services/gateway/` |

### Livrables
- [ ] Package `shared` compilable
- [ ] Gateway GraphQL fonctionnelle
- [ ] Docker compose pour infra locale
- [ ] Health checks OK

📄 **Plan détaillé** : [backend/GENERATION_PLAN.md](./backend/GENERATION_PLAN.md)

---

## 🔷 Phase 2 : Backend Core Services (~16h)

### Objectif
Implémenter les services métier critiques.

### Étapes
| # | Service | Durée | Responsabilités |
|---|---------|-------|-----------------|
| 2.1 | Identity Service | 6h | Auth, Users, Subscriptions, CNI |
| 2.2 | Partner Service | 4h | Partners, Establishments, Teams |
| 2.3 | Discovery Service | 6h | Offers, Categories, Search |

### Livrables
- [ ] Inscription/Connexion fonctionnelle
- [ ] CRUD Partenaires
- [ ] CRUD Offres avec recherche géo
- [ ] Tests unitaires >80%

📄 **Plan détaillé** : [backend/GENERATION_PLAN.md](./backend/GENERATION_PLAN.md)

---

## 🔷 Phase 3 : Backend Business Services (~12h)

### Objectif
Compléter les services métier et support.

### Étapes
| # | Service | Durée | Responsabilités |
|---|---------|-------|-----------------|
| 3.1 | Booking Service | 4h | Outings, Check-in QR |
| 3.2 | Engagement Service | 4h | Favorites, Reviews, Messaging |
| 3.3 | Notification Service | 2h | Push, Email, SMS |
| 3.4 | Observability | 2h | Tracing, Metrics, Logging |

### Livrables
- [ ] Réservation + QR code
- [ ] Favoris et avis
- [ ] Messagerie
- [ ] Notifications push
- [ ] Tracing Jaeger

📄 **Plan détaillé** : [backend/GENERATION_PLAN.md](./backend/GENERATION_PLAN.md)

---

## 🔷 Phase 4 : Application Mobile Flutter (~46h)

### Objectif
Développer l'app mobile respectant le design Figma.

### Étapes
| # | Tâche | Durée |
|---|-------|-------|
| 4.1 | Setup & Core | 5h |
| 4.2 | Design System (composants) | 7h |
| 4.3 | Features Auth (inscription, CNI) | 9h |
| 4.4 | Features Core (offres, réservations) | 14h |
| 4.5 | Features Social (messagerie, profil) | 6h |
| 4.6 | Tests & Polish | 5h |

### Livrables
- [ ] App iOS + Android fonctionnelles
- [ ] Design Figma respecté
- [ ] Biométrie (Face ID / Touch ID)
- [ ] Mode offline
- [ ] Tests E2E inscription

📄 **Plan détaillé** : [app-mobile/GENERATION_PLAN.md](./app-mobile/GENERATION_PLAN.md)

---

## 🔷 Phase 5 : Site Partenaires (~31h)

### Objectif
Portail web pour les partenaires (business.yousoon.com).

### Étapes
| # | Tâche | Durée |
|---|-------|-------|
| 5.1 | Setup & Config | 3h |
| 5.2 | Layout & Navigation | 3h |
| 5.3 | Auth & 2FA | 4h |
| 5.4 | Dashboard & Analytics | 6h |
| 5.5 | Gestion Offres (wizard) | 6h |
| 5.6 | Établissements & Équipe | 5h |
| 5.7 | Réservations & Settings | 4h |

### Livrables
- [ ] Inscription partenaire complète
- [ ] CRUD offres avec wizard
- [ ] Analytics avec graphiques
- [ ] 2FA obligatoire

📄 **Plan détaillé** : [site-partenaires/GENERATION_PLAN.md](./site-partenaires/GENERATION_PLAN.md)

---

## 🔷 Phase 6 : Site Vitrine (~12h)

### Objectif
Site de présentation (www.yousoon.com).

### Étapes
| # | Tâche | Durée |
|---|-------|-------|
| 6.1 | Setup Next.js 14 | 1h |
| 6.2 | Homepage | 4h |
| 6.3 | Pages secondaires | 3h |
| 6.4 | SEO & Performance | 2h |
| 6.5 | Internationalisation | 2h |

### Livrables
- [ ] Site responsive
- [ ] Lighthouse > 95
- [ ] FR/EN
- [ ] SEO optimisé

📄 **Plan détaillé** : [site-vitrine/GENERATION_PLAN.md](./site-vitrine/GENERATION_PLAN.md)

---

## 🔷 Phase 7 : Admin Backoffice (~26h)

### Objectif
Interface d'administration interne.

### Étapes
| # | Tâche | Durée |
|---|-------|-------|
| 7.1 | Setup & Auth admin | 3h |
| 7.2 | Dashboard | 3h |
| 7.3 | Gestion Utilisateurs | 3h |
| 7.4 | Gestion Partenaires & Offres | 3h |
| 7.5 | Vérification CNI | 5h |
| 7.6 | Modération & Analytics | 4h |
| 7.7 | Configuration & Audit | 5h |

### Livrables
- [ ] Dashboard KPIs
- [ ] Validation CNI avec viewer
- [ ] Modération avis
- [ ] Audit logs complet

📄 **Plan détaillé** : [admin/GENERATION_PLAN.md](./admin/GENERATION_PLAN.md)

---

## 🔷 Phase 8 : Tests & Déploiement (~15h)

### Objectif
Finaliser les tests et préparer le déploiement Kubernetes.

### Étapes
| # | Tâche | Durée |
|---|-------|-------|
| 8.1 | Tests E2E cross-platform | 4h |
| 8.2 | Kubernetes manifests | 4h |
| 8.3 | CI/CD pipelines | 4h |
| 8.4 | Monitoring setup | 3h |

### Livrables
- [ ] Tests E2E passent
- [ ] Déploiement K8s fonctionnel
- [ ] CI/CD automatisé
- [ ] Dashboards Grafana

---

## 📊 Récapitulatif des Estimations

| Module | Durée | % du projet |
|--------|-------|-------------|
| Backend Infrastructure | 9h | 6% |
| Backend Core Services | 16h | 10% |
| Backend Business Services | 12h | 8% |
| App Mobile Flutter | 46h | 30% |
| Site Partenaires | 31h | 20% |
| Site Vitrine | 12h | 8% |
| Admin Backoffice | 26h | 17% |
| Tests & Déploiement | 15h | 10% |
| **TOTAL** | **~155h** | 100% |

**Équivalent** : ~4 semaines à temps plein (40h/semaine)

---

## 🚀 Commencer la Génération

### Prérequis
1. ✅ Documentation validée (ce fichier)
2. ✅ Copilot Instructions à jour
3. ✅ Plans de génération par module créés

### Commande de démarrage
```
"Génère le Backend Phase 1 en suivant le plan backend/GENERATION_PLAN.md"
```

### Workflow recommandé
1. **Générer par phase** : Une phase à la fois
2. **Valider chaque étape** : Tests avant passage à la suite
3. **Itérer** : Ajuster le plan si nécessaire

---

## 🔗 Références

- [Instructions Copilot](../.github/copilot-instructions.md)
- [Architecture Backend](./backend/ARCHITECTURE.md)
- [Data Model](./DATA_MODEL.md)
- [Design System](./DESIGN_SYSTEM.md)
- [Figma](https://www.figma.com/design/1GXJECHtsYzq46OYbSHiaj/Yousoon-Test2)
