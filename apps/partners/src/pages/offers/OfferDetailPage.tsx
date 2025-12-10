import { useParams, Link } from 'react-router-dom'
import {
  ArrowLeft,
  Edit,
  Eye,
  Calendar,
  MapPin,
  Tag,
  Clock,
  Users,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

// Mock data
const mockOffer = {
  id: '1',
  title: 'Happy Hour -30%',
  description:
    'Profitez de 30% de réduction sur toutes les boissons pendant notre Happy Hour ! Une occasion parfaite pour se retrouver entre amis après le travail.',
  establishment: 'Le Petit Bistrot',
  address: '123 Rue de la Paix, 75001 Paris',
  category: 'Bar',
  discount: {
    type: 'percentage',
    value: 30,
  },
  schedule: {
    startDate: '2025-01-01',
    endDate: '2025-12-31',
    slots: [
      { day: 'Lundi', start: '17:00', end: '20:00' },
      { day: 'Mardi', start: '17:00', end: '20:00' },
      { day: 'Mercredi', start: '17:00', end: '20:00' },
      { day: 'Jeudi', start: '17:00', end: '20:00' },
      { day: 'Vendredi', start: '17:00', end: '21:00' },
    ],
  },
  quota: {
    total: 100,
    used: 34,
  },
  status: 'active',
  stats: {
    views: 456,
    bookings: 34,
    checkins: 28,
    conversionRate: 7.5,
  },
  images: [
    'https://images.unsplash.com/photo-1551024709-8f23befc6f87?w=600&h=400&fit=crop',
    'https://images.unsplash.com/photo-1566417713940-fe7c737a9ef2?w=600&h=400&fit=crop',
  ],
  conditions: ['Offre valable uniquement sur les boissons', 'Non cumulable avec d\'autres offres'],
  createdAt: '2025-01-01',
  updatedAt: '2025-12-08',
}

export function OfferDetailPage() {
  const { id } = useParams()

  // In real app, fetch offer by id
  const offer = mockOffer

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Link to="/offers">
          <Button variant="ghost" size="icon">
            <ArrowLeft className="h-5 w-5" />
          </Button>
        </Link>
        <div className="flex-1">
          <h1 className="text-2xl font-bold">{offer.title}</h1>
          <p className="text-muted-foreground">
            ID: {id} • Créée le {offer.createdAt}
          </p>
        </div>
        <Button variant="outline">
          <Eye className="mr-2 h-4 w-4" />
          Prévisualiser
        </Button>
        <Button>
          <Edit className="mr-2 h-4 w-4" />
          Modifier
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Vues
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{offer.stats.views}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Réservations
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{offer.stats.bookings}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Check-ins
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{offer.stats.checkins}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Conversion
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{offer.stats.conversionRate}%</div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Main content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Images */}
          <Card>
            <CardHeader>
              <CardTitle>Images</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 grid-cols-2">
                {offer.images.map((image, index) => (
                  <img
                    key={index}
                    src={image}
                    alt={`${offer.title} ${index + 1}`}
                    className="w-full h-48 object-cover rounded-lg"
                  />
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Description */}
          <Card>
            <CardHeader>
              <CardTitle>Description</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-muted-foreground">{offer.description}</p>
            </CardContent>
          </Card>

          {/* Conditions */}
          <Card>
            <CardHeader>
              <CardTitle>Conditions</CardTitle>
            </CardHeader>
            <CardContent>
              <ul className="list-disc list-inside space-y-1 text-muted-foreground">
                {offer.conditions.map((condition, index) => (
                  <li key={index}>{condition}</li>
                ))}
              </ul>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Discount */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Tag className="h-4 w-4" />
                Réduction
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-primary">
                -{offer.discount.value}%
              </div>
            </CardContent>
          </Card>

          {/* Location */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <MapPin className="h-4 w-4" />
                Établissement
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <p className="font-medium">{offer.establishment}</p>
              <p className="text-sm text-muted-foreground">{offer.address}</p>
            </CardContent>
          </Card>

          {/* Schedule */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Clock className="h-4 w-4" />
                Horaires
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              {offer.schedule.slots.map((slot, index) => (
                <div key={index} className="flex justify-between text-sm">
                  <span>{slot.day}</span>
                  <span className="text-muted-foreground">
                    {slot.start} - {slot.end}
                  </span>
                </div>
              ))}
            </CardContent>
          </Card>

          {/* Validity */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Calendar className="h-4 w-4" />
                Validité
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <div className="flex justify-between text-sm">
                <span>Début</span>
                <span className="text-muted-foreground">{offer.schedule.startDate}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span>Fin</span>
                <span className="text-muted-foreground">{offer.schedule.endDate}</span>
              </div>
            </CardContent>
          </Card>

          {/* Quota */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Users className="h-4 w-4" />
                Quota
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span>Utilisations</span>
                  <span className="font-medium">
                    {offer.quota.used} / {offer.quota.total}
                  </span>
                </div>
                <div className="h-2 bg-muted rounded-full overflow-hidden">
                  <div
                    className="h-full bg-primary transition-all"
                    style={{
                      width: `${(offer.quota.used / offer.quota.total) * 100}%`,
                    }}
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
