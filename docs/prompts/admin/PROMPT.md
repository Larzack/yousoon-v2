# 🔐 Admin Backoffice - Prompt Détaillé

> **Module** : Administration interne Yousoon  
> **URL** : admin.yousoon.com (accès interne via port-forward)  
> **Technologie** : React + TypeScript  
> **Accès** : Équipe Yousoon uniquement (non public)

---

## 🎯 Objectifs

Le backoffice admin permet à l'équipe Yousoon de :
- Valider/bloquer les comptes partenaires
- Valider/bloquer les offres
- Valider les vérifications d'identité (CNI)
- Gérer les utilisateurs et abonnements
- Modérer les avis
- Consulter les analytics globaux
- Gérer la configuration (plans, catégories, etc.)

---

## 🛠️ Stack Technique

### Core
| Technologie | Version | Justification |
|-------------|---------|---------------|
| React | 19.x | Dernière version stable, cohérence avec site partenaires |
| TypeScript | 5.x | Type safety, maintenabilité |
| Vite | 5.x | Build rapide, HMR instantané |

### State Management
| Technologie | Usage |
|-------------|-------|
| **Zustand** | État global léger |
| **TanStack Query** | Cache serveur, sync API |

### Communication API
| Technologie | Usage |
|-------------|-------|
| **urql** | Client GraphQL léger |
| **graphql-codegen** | Types auto-générés |

### UI/Styling
| Technologie | Usage |
|-------------|-------|
| **TailwindCSS** | Utility-first CSS |
| **shadcn/ui** | Composants accessibles |
| **Radix UI** | Primitives headless |
| **Lucide Icons** | Iconographie |

### Formulaires
| Technologie | Usage |
|-------------|-------|
| **React Hook Form** | Gestion formulaires |
| **Zod** | Validation schémas |

### Tableaux & Graphiques
| Technologie | Usage |
|-------------|-------|
| **TanStack Table** | Tableaux de données |
| **Recharts** | Graphiques |

### Tests
| Type | Technologie |
|------|-------------|
| Unit | Vitest |
| Component | Testing Library |
| E2E | Playwright |

### Autres
| Technologie | Usage |
|-------------|-------|
| **React Router** | Routing |
| **date-fns** | Manipulation dates |

---

## 🏗️ Architecture

### Structure du Projet

```
apps/web-admin/
├── src/
│   ├── main.tsx
│   ├── App.tsx
│   ├── routes/
│   │   ├── index.tsx                # Dashboard
│   │   ├── auth/
│   │   │   └── login.tsx
│   │   ├── users/
│   │   │   ├── index.tsx            # Liste utilisateurs
│   │   │   └── [id].tsx             # Détail utilisateur
│   │   ├── partners/
│   │   │   ├── index.tsx            # Liste partenaires
│   │   │   ├── [id].tsx             # Détail partenaire
│   │   │   └── pending.tsx          # En attente validation
│   │   ├── offers/
│   │   │   ├── index.tsx            # Liste offres
│   │   │   ├── [id].tsx             # Détail offre
│   │   │   └── pending.tsx          # En attente validation
│   │   ├── identity/
│   │   │   ├── index.tsx            # Vérifications en attente
│   │   │   └── [id].tsx             # Détail vérification
│   │   ├── reviews/
│   │   │   ├── index.tsx            # Liste avis
│   │   │   └── reported.tsx         # Avis signalés
│   │   ├── subscriptions/
│   │   │   ├── index.tsx            # Abonnements actifs
│   │   │   └── plans.tsx            # Gestion des plans
│   │   ├── analytics/
│   │   │   └── index.tsx            # Stats globales
│   │   └── settings/
│   │       ├── categories.tsx       # Gestion catégories
│   │       ├── config.tsx           # Configuration app
│   │       └── team.tsx             # Équipe admin
│   ├── components/
│   │   ├── layout/
│   │   │   ├── Sidebar.tsx
│   │   │   └── AdminLayout.tsx
│   │   └── shared/
│   │       ├── DataTable.tsx
│   │       ├── StatusBadge.tsx
│   │       ├── ActionMenu.tsx
│   │       └── ConfirmDialog.tsx
│   └── lib/
│       └── graphql/
├── package.json
└── vite.config.ts
```

---

## 📱 Fonctionnalités par Section

### 1. Dashboard

Vue d'ensemble avec KPIs :
- Utilisateurs actifs (jour/semaine/mois)
- Nouveaux inscrits
- Partenaires en attente de validation
- Offres en attente de validation
- Vérifications CNI en attente
- Avis signalés
- Revenus (abonnements)

### 2. Gestion Utilisateurs

**Liste utilisateurs** :
- Recherche par email, nom
- Filtres : statut, abonnement, identité vérifiée
- Actions : voir, suspendre, supprimer

**Détail utilisateur** :
- Informations profil
- Historique abonnements
- Historique réservations
- Favoris
- Avis postés
- Statut vérification identité
- Actions : modifier, suspendre, réinitialiser mot de passe

### 3. Gestion Partenaires

**Liste partenaires** :
- Recherche par nom, SIRET
- Filtres : statut (pending, active, suspended)
- Actions : voir, valider, bloquer

**Détail partenaire** :
- Informations entreprise
- Établissements
- Offres publiées
- Statistiques
- Actions : valider, suspendre, supprimer

**En attente de validation** :
- Liste des nouveaux partenaires
- Quick actions : valider / rejeter

### 4. Gestion Offres

**Liste offres** :
- Recherche par titre, partenaire
- Filtres : statut, catégorie, partenaire
- Actions : voir, valider, bloquer

**En attente de validation** (optionnel) :
- Si modération active
- Quick actions : approuver / rejeter

### 5. Vérification Identité (CNI)

**Liste vérifications** :
- Vérifications en attente
- Données extraites par OCR
- Photo du document
- Photo selfie (si applicable)

**Actions** :
- Valider le compte
- Rejeter avec motif
- Demander nouvelle soumission

**Interface validation** :
- Affichage image CNI
- Données extraites (nom, prénom, date naissance, numéro)
- Comparaison avec profil utilisateur
- Boutons Valider / Rejeter

### 6. Modération Avis

**Liste avis** :
- Tous les avis
- Filtres : note, signalés, récents

**Avis signalés** :
- Avis reportés par utilisateurs/partenaires
- Actions : conserver, supprimer, suspendre auteur

### 7. Gestion Abonnements

**Abonnements actifs** :
- Liste avec filtres (plan, statut, période)
- Actions : annuler, rembourser

**Gestion des plans** :
- Modifier prix, durée essai
- Activer/désactiver un plan
- Créer nouveau plan

### 8. Analytics

**Métriques globales** :
- Utilisateurs : inscrits, actifs, churn
- Partenaires : actifs, nouvelles inscriptions
- Offres : créées, réservations
- Revenus : MRR, ARR, conversion essai

**Graphiques** :
- Évolution temporelle
- Répartition par plan
- Top partenaires/offres

### 9. Configuration

**Catégories** :
- CRUD catégories d'offres
- Ordre d'affichage
- Icônes et couleurs

**Configuration app** :
- Durée période d'essai (défaut 30j)
- Paramètres divers

**Équipe admin** :
- Gestion des comptes admin
- Rôles et permissions

---

## 🔐 Sécurité

### Accès
- **Pas d'accès public** : uniquement via `kubectl port-forward`
- **Authentification** : JWT avec rôle admin
- **IP Whitelist** : optionnel si exposé

### Rôles Admin
| Rôle | Permissions |
|------|-------------|
| **super_admin** | Tout |
| **moderator** | Validation partenaires, offres, avis |
| **support** | Lecture + gestion utilisateurs |

### Audit Log
- Toutes les actions admin sont loggées
- Qui, quoi, quand

---

## 🚀 Déploiement

### Accès via Port-Forward

```bash
# Accès local à l'admin
kubectl port-forward svc/web-admin 3000:80 -n yousoon

# Puis accéder à http://localhost:3000
```

### Pas d'Ingress public
- Service de type ClusterIP
- Aucune exposition externe
- Accès uniquement via kubectl

---

## 📋 Checklist

- [ ] Authentification admin sécurisée
- [ ] Toutes les actions loggées
- [ ] Validation CNI avec affichage images
- [ ] Modération avis
- [ ] Gestion plans abonnement
- [ ] Analytics de base
- [ ] Responsive (desktop principalement)

---

## 🔗 Références

- [MASTER_PROMPT.md](../MASTER_PROMPT.md)
- [Backend API](../backend/PROMPT.md)
- [Modèle de données](../DATA_MODEL.md)
