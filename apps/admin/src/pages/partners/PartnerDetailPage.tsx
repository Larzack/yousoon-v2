import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { getInitials, formatDate } from '@/lib/utils'
import {
  ArrowLeft,
  Building2,
  MapPin,
  Phone,
  Mail,
  Globe,
  Tag,
  Star,
  Users,
  Calendar,
  CheckCircle,
  XCircle,
  Ban,
  Edit,
  ExternalLink,
  AlertCircle,
} from 'lucide-react'

interface Establishment {
  id: string
  name: string
  address: string
  city: string
  isActive: boolean
}

interface Offer {
  id: string
  title: string
  discount: string
  status: string
  bookingsCount: number
}

// Mock data
const mockPartner = {
  id: '1',
  companyName: 'Le Petit Bistrot SARL',
  tradeName: 'Le Petit Bistrot',
  logo: '',
  siret: '123 456 789 00012',
  vatNumber: 'FR12345678901',
  legalForm: 'SARL',
  category: 'Restaurant',
  subcategories: ['Bistrot', 'Cuisine française'],
  status: 'active' as const,
  description: 'Un bistrot chaleureux au cœur de Paris proposant une cuisine traditionnelle française avec des produits frais du marché.',
  createdAt: '2024-06-15',
  verifiedAt: '2024-06-18',
  contact: {
    firstName: 'Pierre',
    lastName: 'Martin',
    email: 'pierre@lepetitbistrot.fr',
    phone: '+33 1 23 45 67 89',
    role: 'Gérant',
  },
  branding: {
    primaryColor: '#8B4513',
    website: 'https://lepetitbistrot.fr',
  },
  stats: {
    establishmentsCount: 2,
    offersCount: 5,
    totalBookings: 456,
    totalCheckins: 412,
    avgRating: 4.5,
    reviewsCount: 128,
  },
  establishments: [
    { id: 'e1', name: 'Le Petit Bistrot - Marais', address: '15 Rue des Rosiers', city: 'Paris 4e', isActive: true },
    { id: 'e2', name: 'Le Petit Bistrot - Bastille', address: '42 Rue de la Roquette', city: 'Paris 11e', isActive: true },
  ] as Establishment[],
  offers: [
    { id: 'o1', title: '-20% sur l\'addition', discount: '20%', status: 'active', bookingsCount: 156 },
    { id: 'o2', title: 'Menu du midi à 15€', discount: '15€', status: 'active', bookingsCount: 98 },
    { id: 'o3', title: 'Apéritif offert', discount: 'Gratuit', status: 'paused', bookingsCount: 234 },
  ] as Offer[],
  team: [
    { id: 't1', name: 'Pierre Martin', email: 'pierre@lepetitbistrot.fr', role: 'admin', status: 'active' },
    { id: 't2', name: 'Marie Dupont', email: 'marie@lepetitbistrot.fr', role: 'manager', status: 'active' },
  ],
}

function getStatusBadge(status: string) {
  switch (status) {
    case 'active':
      return (
        <span className="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-700">
          <CheckCircle className="h-4 w-4" />
          Actif
        </span>
      )
    case 'pending':
      return (
        <span className="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm font-medium bg-yellow-100 text-yellow-700">
          <AlertCircle className="h-4 w-4" />
          En attente
        </span>
      )
    case 'suspended':
      return (
        <span className="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm font-medium bg-red-100 text-red-700">
          <XCircle className="h-4 w-4" />
          Suspendu
        </span>
      )
    default:
      return null
  }
}

export function PartnerDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [activeTab, setActiveTab] = useState('overview')

  const partner = mockPartner // In real app, fetch by id

  const handleValidate = () => {
    // TODO: API call to validate partner
    console.log('Validating partner', id)
  }

  const handleSuspend = () => {
    // TODO: API call to suspend partner
    console.log('Suspending partner', id)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div className="flex-1">
          <div className="flex items-center gap-3">
            <Avatar className="h-12 w-12">
              <AvatarImage src={partner.logo} />
              <AvatarFallback className="bg-purple-100 text-purple-600 text-lg">
                {getInitials(partner.tradeName)}
              </AvatarFallback>
            </Avatar>
            <div>
              <h1 className="text-2xl font-bold">{partner.tradeName}</h1>
              <p className="text-muted-foreground">{partner.companyName}</p>
            </div>
          </div>
        </div>
        <div className="flex items-center gap-3">
          {getStatusBadge(partner.status)}
          {partner.status === 'pending' && (
            <>
              <Button onClick={handleValidate} className="gap-2">
                <CheckCircle className="h-4 w-4" />
                Valider
              </Button>
              <Button variant="destructive" className="gap-2">
                <XCircle className="h-4 w-4" />
                Rejeter
              </Button>
            </>
          )}
          {partner.status === 'active' && (
            <Button variant="destructive" onClick={handleSuspend} className="gap-2">
              <Ban className="h-4 w-4" />
              Suspendre
            </Button>
          )}
        </div>
      </div>

      {/* Tabs */}
      <div className="flex gap-2 border-b">
        {['overview', 'establishments', 'offers', 'team', 'activity'].map((tab) => (
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
            {tab === 'establishments' && `Établissements (${partner.establishments.length})`}
            {tab === 'offers' && `Offres (${partner.offers.length})`}
            {tab === 'team' && `Équipe (${partner.team.length})`}
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
                <CardTitle>Informations entreprise</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid gap-4 md:grid-cols-2">
                  <div>
                    <p className="text-sm text-muted-foreground">Raison sociale</p>
                    <p className="font-medium">{partner.companyName}</p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Forme juridique</p>
                    <p className="font-medium">{partner.legalForm}</p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">SIRET</p>
                    <p className="font-medium font-mono">{partner.siret}</p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">N° TVA</p>
                    <p className="font-medium font-mono">{partner.vatNumber}</p>
                  </div>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground mb-1">Description</p>
                  <p className="text-sm">{partner.description}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground mb-2">Catégories</p>
                  <div className="flex gap-2 flex-wrap">
                    <span className="px-2 py-1 bg-gray-100 rounded-full text-sm font-medium">
                      {partner.category}
                    </span>
                    {partner.subcategories.map((sub) => (
                      <span key={sub} className="px-2 py-1 bg-gray-50 rounded-full text-sm text-muted-foreground">
                        {sub}
                      </span>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Contact principal</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex items-center gap-3">
                  <Users className="h-4 w-4 text-muted-foreground" />
                  <span>{partner.contact.firstName} {partner.contact.lastName} ({partner.contact.role})</span>
                </div>
                <div className="flex items-center gap-3">
                  <Mail className="h-4 w-4 text-muted-foreground" />
                  <a href={`mailto:${partner.contact.email}`} className="text-primary hover:underline">
                    {partner.contact.email}
                  </a>
                </div>
                <div className="flex items-center gap-3">
                  <Phone className="h-4 w-4 text-muted-foreground" />
                  <span>{partner.contact.phone}</span>
                </div>
                {partner.branding.website && (
                  <div className="flex items-center gap-3">
                    <Globe className="h-4 w-4 text-muted-foreground" />
                    <a href={partner.branding.website} target="_blank" rel="noopener noreferrer" className="text-primary hover:underline flex items-center gap-1">
                      {partner.branding.website}
                      <ExternalLink className="h-3 w-3" />
                    </a>
                  </div>
                )}
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
                  <span className="text-muted-foreground flex items-center gap-2">
                    <MapPin className="h-4 w-4" />
                    Établissements
                  </span>
                  <span className="font-bold">{partner.stats.establishmentsCount}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground flex items-center gap-2">
                    <Tag className="h-4 w-4" />
                    Offres actives
                  </span>
                  <span className="font-bold">{partner.stats.offersCount}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground flex items-center gap-2">
                    <Calendar className="h-4 w-4" />
                    Réservations
                  </span>
                  <span className="font-bold">{partner.stats.totalBookings}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground flex items-center gap-2">
                    <CheckCircle className="h-4 w-4" />
                    Check-ins
                  </span>
                  <span className="font-bold">{partner.stats.totalCheckins}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground flex items-center gap-2">
                    <Star className="h-4 w-4" />
                    Note moyenne
                  </span>
                  <span className="font-bold flex items-center gap-1">
                    <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                    {partner.stats.avgRating} ({partner.stats.reviewsCount})
                  </span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Historique</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Inscription</span>
                  <span>{formatDate(partner.createdAt)}</span>
                </div>
                {partner.verifiedAt && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Validation</span>
                    <span>{formatDate(partner.verifiedAt)}</span>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </div>
      )}

      {activeTab === 'establishments' && (
        <div className="grid gap-4 md:grid-cols-2">
          {partner.establishments.map((est) => (
            <Card key={est.id}>
              <CardContent className="pt-6">
                <div className="flex items-start justify-between">
                  <div>
                    <h3 className="font-medium">{est.name}</h3>
                    <p className="text-sm text-muted-foreground mt-1">{est.address}</p>
                    <p className="text-sm text-muted-foreground">{est.city}</p>
                  </div>
                  <span className={`px-2 py-1 rounded-full text-xs ${est.isActive ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'}`}>
                    {est.isActive ? 'Actif' : 'Inactif'}
                  </span>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {activeTab === 'offers' && (
        <Card>
          <CardContent className="p-0">
            <table className="w-full">
              <thead className="bg-gray-50 border-b">
                <tr>
                  <th className="text-left p-4 font-medium text-muted-foreground">Offre</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Réduction</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Statut</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Réservations</th>
                  <th className="text-right p-4 font-medium text-muted-foreground">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {partner.offers.map((offer) => (
                  <tr key={offer.id} className="hover:bg-gray-50">
                    <td className="p-4 font-medium">{offer.title}</td>
                    <td className="p-4">{offer.discount}</td>
                    <td className="p-4">
                      <span className={`px-2 py-1 rounded-full text-xs ${
                        offer.status === 'active' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'
                      }`}>
                        {offer.status === 'active' ? 'Active' : 'En pause'}
                      </span>
                    </td>
                    <td className="p-4">{offer.bookingsCount}</td>
                    <td className="p-4 text-right">
                      <Link to={`/offers/${offer.id}`}>
                        <Button variant="ghost" size="sm">Voir</Button>
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      )}

      {activeTab === 'team' && (
        <Card>
          <CardContent className="p-0">
            <table className="w-full">
              <thead className="bg-gray-50 border-b">
                <tr>
                  <th className="text-left p-4 font-medium text-muted-foreground">Membre</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Rôle</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Statut</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {partner.team.map((member) => (
                  <tr key={member.id} className="hover:bg-gray-50">
                    <td className="p-4">
                      <div>
                        <p className="font-medium">{member.name}</p>
                        <p className="text-sm text-muted-foreground">{member.email}</p>
                      </div>
                    </td>
                    <td className="p-4 capitalize">{member.role}</td>
                    <td className="p-4">
                      <span className="px-2 py-1 rounded-full text-xs bg-green-100 text-green-700">
                        Actif
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      )}

      {activeTab === 'activity' && (
        <Card>
          <CardHeader>
            <CardTitle>Journal d'activité</CardTitle>
            <CardDescription>Dernières actions sur ce partenaire</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {[
                { action: 'Nouvelle offre créée', date: '10/12/2024 14:30', admin: null },
                { action: 'Partenaire validé', date: '18/06/2024 10:15', admin: 'Admin Yousoon' },
                { action: 'Inscription', date: '15/06/2024 09:00', admin: null },
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
