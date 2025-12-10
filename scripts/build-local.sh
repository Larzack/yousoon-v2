#!/bin/bash

# ===========================================
# Yousoon - Script de Build Local
# ===========================================

set -e

# Couleurs pour les logs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Répertoire racine
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# Docker compose avec le bon fichier
COMPOSE="docker compose -f apps/docker-compose.yml"

# Fonction d'affichage
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Afficher l'aide
show_help() {
    echo ""
    echo "🚀 Yousoon - Script de Build Local"
    echo ""
    echo "Usage: ./scripts/build-local.sh [OPTION]"
    echo ""
    echo "Options:"
    echo "  all         Build et lance tous les services"
    echo "  infra       Lance uniquement l'infrastructure (MongoDB, Redis, NATS)"
    echo "  backend     Build et lance les services backend"
    echo "  web         Build et lance les apps web (partners, admin)"
    echo "  down        Arrête tous les conteneurs"
    echo "  clean       Arrête et supprime les volumes"
    echo "  logs        Affiche les logs de tous les services"
    echo "  status      Affiche le statut des conteneurs"
    echo "  help        Affiche cette aide"
    echo ""
    echo "Exemples:"
    echo "  ./scripts/build-local.sh infra      # Lance MongoDB, Redis, NATS"
    echo "  ./scripts/build-local.sh all        # Lance tout l'environnement"
    echo "  ./scripts/build-local.sh down       # Arrête tout"
    echo ""
}

# Vérifier Docker
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker n'est pas installé"
        exit 1
    fi

    if ! docker info &> /dev/null; then
        log_error "Docker n'est pas démarré"
        exit 1
    fi

    log_success "Docker est disponible"
}

# Lancer l'infrastructure
start_infra() {
    log_info "Démarrage de l'infrastructure (MongoDB, Redis, NATS)..."
    $COMPOSE up -d mongodb redis nats
    
    log_info "Attente de la disponibilité des services..."
    sleep 5
    
    # Vérifier les services
    if $COMPOSE ps mongodb | grep -q "healthy"; then
        log_success "MongoDB est prêt (port 27017)"
    else
        log_warning "MongoDB démarre..."
    fi
    
    if $COMPOSE ps redis | grep -q "healthy"; then
        log_success "Redis est prêt (port 6379)"
    else
        log_warning "Redis démarre..."
    fi
    
    if $COMPOSE ps nats | grep -q "healthy"; then
        log_success "NATS est prêt (port 4222)"
    else
        log_warning "NATS démarre..."
    fi
}

# Lancer le backend
start_backend() {
    log_info "Build et démarrage des services backend..."
    
    $COMPOSE up -d --build \
        identity-service \
        partner-service \
        discovery-service \
        booking-service \
        engagement-service \
        notification-service \
        apollo-router
    
    log_info "Attente du démarrage des services..."
    sleep 10
    
    log_success "Services backend démarrés"
    echo ""
    echo "📡 Endpoints disponibles :"
    echo "   • API Gateway (GraphQL) : http://localhost:4000/graphql"
    echo "   • Identity Service      : http://localhost:4001/graphql"
    echo "   • Partner Service       : http://localhost:4002/graphql"
    echo "   • Discovery Service     : http://localhost:4003/graphql"
    echo "   • Booking Service       : http://localhost:4004/graphql"
    echo "   • Engagement Service    : http://localhost:4005/graphql"
    echo "   • Notification Service  : http://localhost:4006/graphql"
    echo ""
}

# Lancer les apps web
start_web() {
    log_info "Build et démarrage des applications web..."
    
    $COMPOSE up -d --build partners admin
    
    log_info "Attente du démarrage..."
    sleep 5
    
    log_success "Applications web démarrées"
    echo ""
    echo "🌐 Applications disponibles :"
    echo "   • Site Partenaires : http://localhost:3000"
    echo "   • Admin Backoffice : http://localhost:3001"
    echo ""
}

# Lancer tout
start_all() {
    log_info "🚀 Démarrage complet de la plateforme Yousoon..."
    echo ""
    
    start_infra
    echo ""
    
    start_backend
    echo ""
    
    start_web
    echo ""
    
    log_success "🎉 Plateforme Yousoon démarrée avec succès !"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📊 Résumé des services :"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "🗄️  Infrastructure :"
    echo "    MongoDB          : mongodb://localhost:27017"
    echo "    Redis            : redis://localhost:6379"
    echo "    NATS             : nats://localhost:4222"
    echo "    NATS Monitoring  : http://localhost:8222"
    echo ""
    echo "⚙️  Backend :"
    echo "    API Gateway      : http://localhost:4000/graphql"
    echo ""
    echo "🌐 Applications :"
    echo "    Site Partenaires : http://localhost:3000"
    echo "    Admin Backoffice : http://localhost:3001"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "💡 Commandes utiles :"
    echo "   ./scripts/build-local.sh logs     # Voir les logs"
    echo "   ./scripts/build-local.sh status   # Voir le statut"
    echo "   ./scripts/build-local.sh down     # Arrêter tout"
    echo ""
}

# Arrêter les services
stop_all() {
    log_info "Arrêt de tous les conteneurs..."
    $COMPOSE down
    log_success "Tous les conteneurs ont été arrêtés"
}

# Nettoyer (arrêter + supprimer volumes)
clean_all() {
    log_warning "Arrêt et suppression de tous les conteneurs et volumes..."
    $COMPOSE down -v --remove-orphans
    log_success "Nettoyage terminé"
}

# Afficher les logs
show_logs() {
    $COMPOSE logs -f
}

# Afficher le statut
show_status() {
    echo ""
    echo "📊 Statut des conteneurs Yousoon :"
    echo ""
    $COMPOSE ps
    echo ""
}

# Point d'entrée principal
main() {
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🏗️  Yousoon - Build Local"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    check_docker

    case "${1:-help}" in
        all)
            start_all
            ;;
        infra)
            start_infra
            ;;
        backend)
            start_infra
            start_backend
            ;;
        web)
            start_web
            ;;
        down)
            stop_all
            ;;
        clean)
            clean_all
            ;;
        logs)
            show_logs
            ;;
        status)
            show_status
            ;;
        help|*)
            show_help
            ;;
    esac
}

main "$@"
