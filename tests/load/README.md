# k6 Load Tests

Ce dossier contient les tests de charge pour la plateforme Yousoon.

## Installation

```bash
# macOS
brew install k6

# Linux (Debian/Ubuntu)
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6

# Docker
docker pull grafana/k6
```

## Exécution

### Tests API publiques

```bash
# Local
k6 run tests/load/backend.js -e BASE_URL=http://localhost:8080

# Staging
k6 run tests/load/backend.js -e BASE_URL=https://api.staging.yousoon.com

# Avec output InfluxDB (pour Grafana)
k6 run tests/load/backend.js -e BASE_URL=http://localhost:8080 --out influxdb=http://localhost:8086/k6
```

### Tests utilisateurs authentifiés

```bash
k6 run tests/load/authenticated.js -e BASE_URL=http://localhost:8080
```

### Tests avec paramètres personnalisés

```bash
# Plus d'utilisateurs virtuels
k6 run tests/load/backend.js -e BASE_URL=http://localhost:8080 --vus 100 --duration 10m

# Mode cloud (k6 Cloud)
k6 cloud tests/load/backend.js
```

## Configuration

### Seuils de performance

| Métrique | Seuil | Description |
|----------|-------|-------------|
| `http_req_duration` | p(95) < 500ms | 95% des requêtes sous 500ms |
| `http_req_failed` | < 1% | Moins de 1% d'erreurs |
| `booking_success` | > 95% | Plus de 95% des réservations réussies |

### Scénarios de charge

#### backend.js - API Publiques
- **Ramp up**: 10 → 50 → 100 utilisateurs
- **Durée**: 15 minutes
- **Endpoints testés**: offers, categories, search

#### authenticated.js - Utilisateurs authentifiés
- **Average load**: 30 utilisateurs pendant 5 minutes
- **Peak hour**: Jusqu'à 100 req/s
- **Endpoints testés**: login, profile, bookings, favorites

## Interprétation des résultats

### Métriques clés

- **http_req_duration**: Temps de réponse
- **http_req_failed**: Taux d'erreur
- **vus**: Utilisateurs virtuels actifs
- **iterations**: Nombre total d'itérations

### Exemple de sortie

```
     ✓ status is 200
     ✓ no errors in response

     checks.........................: 99.87% ✓ 15432  ✗ 20
     data_received..................: 12 MB  80 kB/s
     data_sent......................: 2.3 MB 15 kB/s
     http_req_duration..............: avg=45.2ms min=12ms med=38ms max=892ms p(90)=85ms p(95)=120ms
     http_req_failed................: 0.12%  ✓ 20     ✗ 15432
     http_reqs......................: 15452  103/s
     iterations.....................: 5150   34.3/s
     vus............................: 50     min=0    max=100
```

## Intégration CI/CD

```yaml
# .github/workflows/load-test.yml
name: Load Test

on:
  schedule:
    - cron: '0 2 * * *'  # Chaque nuit à 2h
  workflow_dispatch:

jobs:
  load-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install k6
        run: |
          sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
          echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
          sudo apt-get update
          sudo apt-get install k6
      
      - name: Run load test
        run: k6 run tests/load/backend.js -e BASE_URL=${{ secrets.API_URL }}
        
      - name: Upload results
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: k6-results
          path: summary.json
```

## Monitoring avec Grafana

Pour visualiser les résultats en temps réel :

1. Démarrer InfluxDB et Grafana
2. Importer le dashboard k6 (ID: 2587)
3. Exécuter les tests avec `--out influxdb`

```bash
# Docker Compose pour le monitoring
docker-compose -f docker-compose.monitoring.yml up -d

# Exécuter les tests
k6 run tests/load/backend.js --out influxdb=http://localhost:8086/k6
```
