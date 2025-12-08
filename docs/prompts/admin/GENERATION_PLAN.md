# 🚀 Plan de Génération - Admin Backoffice

> **Module** : Administration interne (admin.yousoon.com - non public)  
> **Priorité** : 🟢 Basse (après tous les autres modules)  
> **Dépendances** : Backend complet

---

## 📋 Vue d'Ensemble

```
┌─────────────────────────────────────────────────────────────────┐
│                    ORDRE DE GÉNÉRATION                          │
├─────────────────────────────────────────────────────────────────┤
│  Phase 1: Setup & Configuration                                 │
│  Phase 2: Layout & Dashboard                                    │
│  Phase 3: Gestion Utilisateurs                                  │
│  Phase 4: Gestion Partenaires & Offres                          │
│  Phase 5: Vérification Identité (CNI)                           │
│  Phase 6: Modération & Analytics                                │
│  Phase 7: Configuration & Audit                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📁 Structure Cible

```
apps/admin/
├── src/
│   ├── main.tsx
│   ├── App.tsx
│   ├── routes/
│   │   ├── auth/
│   │   ├── dashboard/
│   │   ├── users/
│   │   ├── partners/
│   │   ├── offers/
│   │   ├── identity/
│   │   ├── reviews/
│   │   ├── subscriptions/
│   │   ├── analytics/
│   │   └── settings/
│   ├── components/
│   │   ├── ui/
│   │   ├── layout/
│   │   └── shared/
│   ├── lib/
│   ├── hooks/
│   ├── stores/
│   └── types/
├── package.json
└── vite.config.ts
```

---

## 🔷 Phase 1 : Setup & Configuration

### Étape 1.1 : Initialisation
**Même stack que Partner Portal :**
- React 19 + TypeScript
- Vite
- TailwindCSS + shadcn/ui
- urql (GraphQL)
- TanStack Query + Table
- Zustand

### Étape 1.2 : Configuration Sécurité
**Fichiers à générer :**
```
src/lib/
├── auth/
│   ├── admin-auth.ts           # Auth spécifique admin
│   └── permissions.ts          # RBAC
└── audit/
    └── logger.ts               # Log toutes les actions
```

**Rôles Admin :**
```typescript
enum AdminRole {
  SUPER_ADMIN = 'super_admin',   // Tout
  MODERATOR = 'moderator',       // Validation partenaires, offres, avis
  SUPPORT = 'support',           // Lecture + gestion utilisateurs
}
```

---

## 🔷 Phase 2 : Layout & Dashboard

### Étape 2.1 : Layout Admin
**Fichiers à générer :**
```
src/components/layout/
├── AdminLayout.tsx
├── Sidebar.tsx
├── Header.tsx
└── Breadcrumbs.tsx
```

**Menu Sidebar :**
```typescript
const adminMenu = [
  { icon: LayoutDashboard, label: 'Dashboard', href: '/' },
  { icon: Users, label: 'Utilisateurs', href: '/users' },
  { icon: Building2, label: 'Partenaires', href: '/partners' },
  { icon: Tag, label: 'Offres', href: '/offers' },
  { icon: IdCard, label: 'Vérifications CNI', href: '/identity', badge: pendingCount },
  { icon: MessageSquare, label: 'Avis', href: '/reviews' },
  { icon: CreditCard, label: 'Abonnements', href: '/subscriptions' },
  { icon: BarChart3, label: 'Analytics', href: '/analytics' },
  { icon: Settings, label: 'Configuration', href: '/settings' },
];
```

### Étape 2.2 : Dashboard Admin
**Fichiers à générer :**
```
src/routes/dashboard/
└── index.tsx

src/components/dashboard/
├── AdminStatsCards.tsx
├── PendingActions.tsx          # Actions en attente
├── RecentActivity.tsx
└── SystemHealth.tsx
```

**KPIs Dashboard :**
- Utilisateurs actifs (jour/semaine/mois)
- Nouveaux inscrits
- Partenaires en attente de validation
- Offres en attente
- Vérifications CNI en attente
- Avis signalés
- Revenus MRR

---

## 🔷 Phase 3 : Gestion Utilisateurs

### Étape 3.1 : Liste Utilisateurs
**Fichiers à générer :**
```
src/routes/users/
├── index.tsx                   # Liste avec filtres
└── [id].tsx                    # Détail utilisateur

src/components/users/
├── UserTable.tsx
├── UserFilters.tsx
├── UserStatusBadge.tsx
├── UserActions.tsx
└── UserDetailCard.tsx
```

**Colonnes Table :**
- Avatar + Nom
- Email
- Statut (actif, suspendu)
- Abonnement
- Identité vérifiée
- Date inscription
- Dernière connexion
- Actions

### Étape 3.2 : Détail Utilisateur
**Sections :**
- Informations profil
- Historique abonnements
- Historique réservations
- Favoris
- Avis postés
- Statut vérification identité
- **Actions** : Modifier, Suspendre, Réinitialiser mot de passe, Supprimer

---

## 🔷 Phase 4 : Gestion Partenaires & Offres

### Étape 4.1 : Liste Partenaires
**Fichiers à générer :**
```
src/routes/partners/
├── index.tsx
├── pending.tsx                 # En attente validation
└── [id].tsx

src/components/partners/
├── PartnerTable.tsx
├── PartnerStatusBadge.tsx
├── PartnerValidationCard.tsx
├── PartnerActions.tsx
└── QuickApproveDialog.tsx
```

**Workflow Validation :**
1. Liste des partenaires "pending"
2. Review des informations
3. Vérification SIRET (API INSEE)
4. Approuver / Rejeter avec motif

### Étape 4.2 : Gestion Offres
**Fichiers à générer :**
```
src/routes/offers/
├── index.tsx
├── pending.tsx                 # Si modération active
└── [id].tsx

src/components/offers/
├── OfferTable.tsx
├── OfferModerationCard.tsx
└── OfferActions.tsx
```

---

## 🔷 Phase 5 : Vérification Identité (CNI)

### Étape 5.1 : Liste Vérifications
**Fichiers à générer :**
```
src/routes/identity/
├── index.tsx                   # Liste pending
└── [id].tsx                    # Détail vérification

src/components/identity/
├── VerificationQueue.tsx
├── VerificationCard.tsx
├── DocumentViewer.tsx          # Affichage image CNI
├── ExtractedDataCard.tsx       # Données OCR
├── ComparisonView.tsx          # Profil vs Document
├── ValidationActions.tsx
└── RejectionDialog.tsx
```

### Étape 5.2 : Interface Validation
**Layout :**
```
┌─────────────────────────────────────────────────────────────────┐
│  Vérification #12345                                    [X]     │
├───────────────────────────────────┬─────────────────────────────┤
│                                   │  Données extraites          │
│   [Image CNI - Recto]             │  ─────────────────          │
│                                   │  Nom: DUPONT                │
│   [Image CNI - Verso]             │  Prénom: Jean               │
│                                   │  Date naissance: 15/03/1990 │
│                                   │  N° Document: 123456789     │
│                                   │                             │
├───────────────────────────────────┼─────────────────────────────┤
│  Profil utilisateur               │                             │
│  ─────────────────                │   [ ✓ VALIDER ]             │
│  Nom: Dupont                      │                             │
│  Prénom: Jean                     │   [ ✗ REJETER ]             │
│  Email: jean.dupont@mail.com      │                             │
└───────────────────────────────────┴─────────────────────────────┘
```

**Actions :**
- ✅ Valider → Utilisateur vérifié
- ❌ Rejeter avec motif → Demande nouvelle soumission
- 🔄 Demander documents supplémentaires

---

## 🔷 Phase 6 : Modération & Analytics

### Étape 6.1 : Modération Avis
**Fichiers à générer :**
```
src/routes/reviews/
├── index.tsx                   # Tous les avis
└── reported.tsx                # Avis signalés

src/components/reviews/
├── ReviewTable.tsx
├── ReviewCard.tsx
├── ReportedReviewCard.tsx
└── ModerationActions.tsx
```

**Actions Modération :**
- Conserver
- Supprimer
- Suspendre auteur

### Étape 6.2 : Analytics Global
**Fichiers à générer :**
```
src/routes/analytics/
└── index.tsx

src/components/analytics/
├── GlobalMetrics.tsx
├── UserGrowthChart.tsx
├── RevenueChart.tsx            # MRR, ARR
├── ConversionFunnel.tsx
├── TopPartnersTable.tsx
└── TopOffersTable.tsx
```

**Métriques :**
- Utilisateurs : inscrits, actifs, churn
- Partenaires : actifs, nouvelles inscriptions
- Offres : créées, réservations
- Revenus : MRR, ARR, conversion essai

---

## 🔷 Phase 7 : Configuration & Audit

### Étape 7.1 : Gestion Abonnements
**Fichiers à générer :**
```
src/routes/subscriptions/
├── index.tsx                   # Abonnements actifs
└── plans.tsx                   # Gestion des plans

src/components/subscriptions/
├── SubscriptionTable.tsx
├── PlanEditor.tsx
└── PlanCard.tsx
```

### Étape 7.2 : Configuration App
**Fichiers à générer :**
```
src/routes/settings/
├── categories.tsx              # CRUD catégories
├── config.tsx                  # Paramètres globaux
└── team.tsx                    # Équipe admin

src/components/settings/
├── CategoryManager.tsx
├── ConfigEditor.tsx
├── AdminTeamTable.tsx
└── InviteAdminDialog.tsx
```

### Étape 7.3 : Audit Logs
**Fichiers à générer :**
```
src/routes/audit/
└── index.tsx

src/components/audit/
├── AuditLogTable.tsx
├── AuditLogFilters.tsx
└── AuditLogDetail.tsx
```

**Colonnes Audit :**
- Timestamp
- Admin (email)
- Action
- Resource
- Resource ID
- IP
- Détails (before/after)

---

## ⏱️ Estimation des Temps

| Phase | Étape | Durée estimée |
|-------|-------|---------------|
| **Phase 1** | Setup | 1h |
| | Auth & Permissions | 2h |
| **Phase 2** | Layout | 1h |
| | Dashboard | 2h |
| **Phase 3** | Liste users | 2h |
| | Détail user | 1h |
| **Phase 4** | Partenaires | 2h |
| | Offres | 1h |
| **Phase 5** | Vérifications CNI | 3h |
| | Interface validation | 2h |
| **Phase 6** | Modération avis | 2h |
| | Analytics | 2h |
| **Phase 7** | Abonnements | 2h |
| | Configuration | 2h |
| | Audit logs | 1h |
| **Total** | | **~26h** |

---

## 🚀 Déploiement

### Accès Restreint
```yaml
# Pas d'Ingress public
# Accès uniquement via kubectl port-forward

# deploy/kubernetes/admin/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: admin-portal
  namespace: yousoon
spec:
  type: ClusterIP        # Pas de LoadBalancer
  ports:
    - port: 80
      targetPort: 3000
  selector:
    app: admin-portal
```

**Accès local :**
```bash
kubectl port-forward svc/admin-portal 3000:80 -n yousoon
# Puis: http://localhost:3000
```

---

## ✅ Critères de Validation

### Sécurité
- [ ] Authentification admin uniquement
- [ ] 2FA obligatoire
- [ ] Toutes actions loggées
- [ ] Pas d'accès public (ClusterIP)
- [ ] RBAC fonctionnel

### Features
- [ ] Dashboard avec métriques temps réel
- [ ] CRUD utilisateurs complet
- [ ] Validation partenaires
- [ ] Vérification CNI avec viewer images
- [ ] Modération avis
- [ ] Gestion plans abonnement
- [ ] Audit logs consultables

### UX
- [ ] Navigation rapide
- [ ] Bulk actions disponibles
- [ ] Filtres et recherche performants

---

## 🔗 Références

- [Prompt Admin](./PROMPT.md)
- [Data Model](../DATA_MODEL.md)
- [Backend API](../backend/ARCHITECTURE.md)
