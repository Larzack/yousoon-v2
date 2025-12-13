import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { formatDate } from '@/lib/utils'
import {
  ArrowLeft,
  Calendar,
  MapPin,
  CheckCircle,
  XCircle,
  Pause,
  Play,
  Trash2,
  Building2,
  Clock,
  Percent,
  AlertCircle,
  Star,
} from 'lucide-react'

// Mock data
const mockOffer = {
  id: '1',
  title: '-20% sur l\'addition',
  description: 'Profitez de 20% de réduction sur l\'ensemble de votre addition au déjeuner comme au dîner. Offre valable du lundi au jeudi.',
  shortDescription: 'Réduction sur l\'addition',
  partner: { id: 'p1', name: 'Le Petit Bistrot', logo: '' },
  establishment: {
    id: 'e1',
    name: 'Le Petit Bistrot - Marais',
    address: '15 Rue des Rosiers',
    city: 'Paris 4e',
  },
  category: 'Restaurant',
  discount: {
    type: 'percentage',
    value: 20,
    originalPrice: null,
    formula: null,
  },
  conditions: [
    { type: 'min_purchase', value: 25, label: 'Minimum 25€ d\'achat' },
    { type: 'days', value: ['lundi', 'mardi', 'mercredi', 'jeudi'], label: 'Du lundi au jeudi' },
  ],
  validity: {
    startDate: '2024-11-01',
    endDate: '2024-12-31',
    timezone: 'Europe/Paris',
  },
  schedule: {
    allDay: false,
    slots: [
      { day: 'Lundi', start: '12:00', end: '14:30' },
      { day: 'Lundi', start: '19:00', end: '22:30' },
      { day: 'Mardi', start: '12:00', end: '14:30' },
      { day: 'Mardi', start: '19:00', end: '22:30' },
    ],
  },
  quota: {
    total: 500,
    perUser: 3,
    perDay: 20,
    used: 298,
  },
  images: ['/placeholder1.jpg', '/placeholder2.jpg'],
  status: 'active' as const,
  stats: {
    views: 2345,
    bookings: 156,
    checkins: 142,
    conversionRate: 6.6,
    avgRating: 4.5,
    reviewsCount: 45,
  },
  moderation: {
    status: 'approved',
    reviewedBy: 'Admin Yousoon',
    reviewedAt: '2024-10-28',
    comment: null,
  },
  createdAt: '2024-10-25',
  publishedAt: '2024-10-28',
}

function getStatusBadge(status: string) {
  const styles: Record<string, { bg: string; text: string; label: string; icon: typeof CheckCircle }> = {
    active: { bg: 'bg-green-100', text: 'text-green-700', label: 'Active', icon: CheckCircle },
    pending: { bg: 'bg-yellow-100', text: 'text-yellow-700', label: 'En attente', icon: AlertCircle },
    paused: { bg: 'bg-orange-100', text: 'text-orange-700', label: 'En pause', icon: Pause },
    expired: { bg: 'bg-red-100', text: 'text-red-700', label: 'Expirée', icon: XCircle },
  }
  const style = styles[status] || styles.active
  const Icon = style.icon
  return (
    <span className={`inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm font-medium ${style.bg} ${style.text}`}>
      <Icon className="h-4 w-4" />
      {style.label}
    </span>
  )
}

export function OfferDetailPage() {
  const { id: _id } = useParams()
  const navigate = useNavigate()
  const [activeTab, setActiveTab] = useState('overview')

  const offer = mockOffer // In real app, fetch by id

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div className="flex-1">
          <h1 className="text-2xl font-bold">{offer.title}</h1>
          <div className="flex items-center gap-2 text-muted-foreground mt-1">
            <Link to={`/partners/${offer.partner.id}`} className="text-primary hover:underline">
              {offer.partner.name}
            </Link>
            <span>•</span>
            <span className="flex items-center gap-1">
              <Building2 className="h-4 w-4" />
              {offer.establishment.name}
            </span>
          </div>
        </div>
        <div className="flex items-center gap-3">
          {getStatusBadge(offer.status)}
          {offer.status === 'active' && (
            <Button variant="outline" className="gap-2">
              <Pause className="h-4 w-4" />
              Mettre en pause
            </Button>
          )}
          {offer.status === 'paused' && (
            <Button className="gap-2">
              <Play className="h-4 w-4" />
              Réactiver
            </Button>
          )}
          <Button variant="destructive" className="gap-2">
            <Trash2 className="h-4 w-4" />
            Archiver
          </Button>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex gap-2 border-b">
        {['overview', 'bookings', 'reviews', 'activity'].map((tab) => (
          <button
            key={tab}
            className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
              activeTab === tab
                ? 'border-primary text-primary'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
            onClick={() => setActiveTab(tab)}
          >
            {tab === 'overview' && 'Vue d\'ensemble'}
            {tab === 'bookings' && `Réservations (${offer.stats.bookings})`}
            {tab === 'reviews' && `Avis (${offer.stats.reviewsCount})`}
            {tab === 'activity' && 'Activité'}
          </button>
        ))}
      </div>

      {/* Content */}
      {activeTab === 'overview' && (
        <div className="grid gap-6 md:grid-cols-3">
          {/* Main Info */}
          <div className="md:col-span-2 space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Description</CardTitle>
              </CardHeader>
              <CardContent>
                <p>{offer.description}</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Réduction</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center gap-4">
                  <div className="h-16 w-16 rounded-lg bg-primary/10 flex items-center justify-center">
                    <Percent className="h-8 w-8 text-primary" />
                  </div>
                  <div>
                    <p className="text-3xl font-bold">{offer.discount.value}%</p>
                    <p className="text-muted-foreground">de réduction sur l'addition</p>
                  </div>
                </div>
                <div>
                  <p className="text-sm font-medium mb-2">Conditions</p>
                  <ul className="space-y-1">
                    {offer.conditions.map((cond, i) => (
                      <li key={i} className="text-sm text-muted-foreground flex items-center gap-2">
                        <CheckCircle className="h-4 w-4 text-green-500" />
                        {cond.label}
                      </li>
                    ))}
                  </ul>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Validité & Horaires</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center gap-4">
                  <div className="flex items-center gap-2 text-muted-foreground">
                    <Calendar className="h-4 w-4" />
                    Du {formatDate(offer.validity.startDate)} au {formatDate(offer.validity.endDate)}
                  </div>
                </div>
                <div>
                  <p className="text-sm font-medium mb-2">Créneaux disponibles</p>
                  <div className="grid gap-2 md:grid-cols-2">
                    {offer.schedule.slots.map((slot, i) => (
                      <div key={i} className="flex items-center gap-2 text-sm">
                        <Clock className="h-4 w-4 text-muted-foreground" />
                        <span className="font-medium">{slot.day}</span>
                        <span className="text-muted-foreground">{slot.start} - {slot.end}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Quotas</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div>
                    <div className="flex justify-between text-sm mb-1">
                      <span>Utilisation globale</span>
                      <span>{offer.quota.used} / {offer.quota.total}</span>
                    </div>
                    <div className="h-2 bg-gray-100 rounded-full overflow-hidden">
                      <div 
                        className="h-full bg-primary rounded-full" 
                        style={{ width: `${(offer.quota.used / offer.quota.total) * 100}%` }}
                      />
                    </div>
                  </div>
                  <div className="grid gap-4 md:grid-cols-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Par utilisateur</span>
                      <span className="font-medium">{offer.quota.perUser} max</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Par jour</span>
                      <span className="font-medium">{offer.quota.perDay} max</span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Statistiques</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Vues</span>
                  <span className="font-bold">{offer.stats.views.toLocaleString()}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Réservations</span>
                  <span className="font-bold">{offer.stats.bookings}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Check-ins</span>
                  <span className="font-bold">{offer.stats.checkins}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Taux de conversion</span>
                  <span className="font-bold">{offer.stats.conversionRate}%</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground flex items-center gap-1">
                    <Star className="h-4 w-4" />
                    Note moyenne
                  </span>
                  <span className="font-bold flex items-center gap-1">
                    <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                    {offer.stats.avgRating} ({offer.stats.reviewsCount})
                  </span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Établissement</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                <p className="font-medium">{offer.establishment.name}</p>
                <p className="text-sm text-muted-foreground flex items-center gap-1">
                  <MapPin className="h-4 w-4" />
                  {offer.establishment.address}, {offer.establishment.city}
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Modération</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Statut</span>
                  <span className="text-green-600 font-medium">Approuvé</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Par</span>
                  <span>{offer.moderation.reviewedBy}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Date</span>
                  <span>{formatDate(offer.moderation.reviewedAt!)}</span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Historique</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Créée le</span>
                  <span>{formatDate(offer.createdAt)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Publiée le</span>
                  <span>{formatDate(offer.publishedAt!)}</span>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      )}

      {activeTab === 'bookings' && (
        <Card>
          <CardHeader>
            <CardTitle>Réservations récentes</CardTitle>
            <CardDescription>Liste des dernières réservations pour cette offre</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground text-center py-8">
              Fonctionnalité à venir - Affichage des réservations
            </p>
          </CardContent>
        </Card>
      )}

      {activeTab === 'reviews' && (
        <Card>
          <CardHeader>
            <CardTitle>Avis</CardTitle>
            <CardDescription>Avis des utilisateurs sur cette offre</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground text-center py-8">
              Fonctionnalité à venir - Affichage des avis
            </p>
          </CardContent>
        </Card>
      )}

      {activeTab === 'activity' && (
        <Card>
          <CardHeader>
            <CardTitle>Journal d'activité</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {[
                { action: 'Statut modifié en "Active"', date: '28/10/2024 10:15', admin: 'Admin Yousoon' },
                { action: 'Offre approuvée', date: '28/10/2024 10:15', admin: 'Admin Yousoon' },
                { action: 'Offre créée', date: '25/10/2024 14:30', admin: null },
              ].map((log, i) => (
                <div key={i} className="flex items-center justify-between py-2 border-b last:border-0">
                  <div>
                    <p className="font-medium">{log.action}</p>
                    {log.admin && <p className="text-sm text-muted-foreground">Par {log.admin}</p>}
                  </div>
                  <span className="text-sm text-muted-foreground">{log.date}</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
