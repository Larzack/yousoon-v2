#!/bin/bash

# ===========================================
# Yousoon - Script de Développement Local
# Lance l'infra Docker + les apps en mode dev
# ===========================================

set -e

# Couleurs
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# Docker compose avec le bon fichier
COMPOSE="docker compose -f apps/docker-compose.yml"

echo ""
echo -e "${BLUE}🚀 Yousoon - Mode Développement${NC}"
echo ""

# Démarrer l'infrastructure
echo -e "${BLUE}📦 Démarrage de l'infrastructure Docker...${NC}"
$COMPOSE up -d mongodb redis nats

echo ""
echo -e "${GREEN}✅ Infrastructure prête !${NC}"
echo ""
echo "🗄️  MongoDB : mongodb://yousoon:yousoon_dev_password@localhost:27017"
echo "📮 Redis   : redis://localhost:6379"
echo "📡 NATS    : nats://localhost:4222"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "💡 Lancez maintenant les apps en mode dev :"
echo ""
echo "   # Site Partenaires (port 5173)"
echo "   cd apps/partners && npm run dev"
echo ""
echo "   # Admin (port 5174)"
echo "   cd apps/admin && npm run dev"
echo ""
echo "   # Site Vitrine (port 3002)"
echo "   cd apps/siteweb && npm run dev"
echo ""
echo "   # App Mobile"
echo "   cd apps/mobile && flutter run"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📝 Pour lancer les services backend Go :"
echo ""
echo "   cd apps/services/identity-service"
echo "   go run ./cmd/main.go"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
