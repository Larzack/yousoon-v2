# ❓ Questions - App Mobile Flutter

> Questions techniques et fonctionnelles à clarifier avant développement

---

## 🎨 Design & UX

### Figma
1. **Lien Figma** : Pouvez-vous partager le lien du fichier Figma ?
2. **Design System** : Y a-t-il un design system défini (couleurs, typographies, espacements) ?
3. **Animations** : Des animations spécifiques sont-elles définies dans le Figma (prototypes) ?
4. **Dark Mode** : Le design inclut-il une version dark mode ?
5. **États** : Tous les états sont-ils designés (loading, empty, error, success) ?

### Écrans à clarifier
- Quels écrans utilisent des sliders/carousels ?
- Y a-t-il des écrans avec scroll horizontal + vertical ?
- Les modales sont-elles des bottom sheets ou des dialogs ?

---

## 🔐 Authentification & Sécurité

### Inscription/Connexion
1. **Méthodes de connexion** : Email/password uniquement ou SSO (Google, Apple, Facebook) ?
2. **Double authentification** : 2FA requis ?
3. **OTP** : SMS ou Email ? Quel fournisseur (Twilio, SendGrid) ?
4. **Session** : Durée de validité du token ? Auto-refresh ?
5. **Biométrie** : Support Face ID / Touch ID pour reconnexion ?

### Validation CNI
6. **Fournisseur choisi** : Onfido, Jumio, Veriff, ID.me, autre ?
7. **Documents acceptés** : CNI française uniquement ? Passeport ? Permis ?
8. **Niveau de vérification** : Juste OCR ou vérification anti-fraude complète ?
9. **Âge minimum** : Y a-t-il une restriction d'âge (18+ pour bars) ?
10. **Fallback** : Que faire si la vérification échoue ?

---

## 📱 Fonctionnalités App

### Géolocalisation
1. **Obligatoire** : La localisation est-elle requise pour utiliser l'app ?
2. **Précision** : Besoin de géofencing pour check-in automatique ?
3. **Background** : Tracking en arrière-plan nécessaire ?

### Notifications
4. **Types** : Quels types de notifications push ?
   - Nouvelles offres à proximité ?
   - Rappel avant sortie ?
   - Offres flash/limitées ?
5. **Personnalisation** : L'utilisateur peut-il choisir ses notifications ?

### Offline
6. **Mode hors-ligne** : Quel niveau de fonctionnement offline ?
   - Consultation des offres mises en cache ?
   - Affichage du QR code de réservation ?
7. **Synchronisation** : Sync automatique au retour du réseau ?

### Paiement
8. **Paiement in-app** : Y a-t-il des achats dans l'app ?
   - Abonnement premium ?
   - Achat de crédits ?
9. **Moyens de paiement** : Apple Pay, Google Pay, CB ?

---

## 🔄 Données & Synchronisation

### Cache
1. **Données sensibles** : Quelles données peuvent être stockées localement ?
2. **Limite taille** : Limite de cache local ?
3. **Invalidation** : Comment invalider le cache (push serveur, TTL) ?

### Real-time
4. **Temps réel** : Des fonctionnalités nécessitent-elles du real-time ?
   - Disponibilité des places ?
   - Notifications instant ?
   - Chat avec partenaires ?
5. **Technologie** : WebSocket, SSE, ou polling ?

---

## 🌍 Internationalisation

1. **Langues** : Quelles langues supporter initialement ?
2. **Formats** : Dates, devises, numéros de téléphone localisés ?
3. **RTL** : Support langues RTL (arabe) prévu ?

---

## 📊 Analytics & Monitoring

1. **Analytics** : ✅ Amplitude
2. **Événements clés** : Quels événements tracker prioritairement ?
3. **Crash reporting** : ✅ Sentry (self-hosted)
4. **Performance** : Prometheus + Grafana

---

## ✅ DÉCISIONS VALIDÉES

### Performance
| Question | Réponse |
|----------|---------|
| iOS minimum | **Dernière version** (iOS 17+) |
| Android minimum | **Dernière version** (API 34+) |
| Taille app | Pas de limite stricte |

### CI/CD & Distribution
| Question | Réponse |
|----------|---------|
| CI/CD | **GitHub Actions** |
| Distribution beta iOS | **TestFlight** |
| Distribution beta Android | **Google Play Internal Testing** (défaut) |

### Stores
| Question | Réponse |
|----------|---------|
| Nom app | **Yousoon** |
| Bundle ID | **com.yousoon.yousoon** |
| Catégorie | **Lifestyle** |

### Analytics & Monitoring
| Question | Réponse |
|----------|---------|
| Analytics | **Amplitude** |
| Crash reporting | **Sentry (self-hosted)** |

---

## 📝 Notes

*Toutes les questions techniques ont été répondues.*
