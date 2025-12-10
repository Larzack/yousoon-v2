import { Link } from 'react-router-dom'
import {
  Tag,
  MapPin,
  CalendarCheck,
  TrendingUp,
  Eye,
  Users,
  ArrowUpRight,
  ArrowDownRight,
} from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

// Mock data - replace with actual API calls
const stats = [
  {
    title: 'Offres actives',
    value: '12',
    change: '+2',
    changeType: 'positive',
    icon: Tag,
    href: '/offers',
  },
  {
    title: 'Vues ce mois',
    value: '1,234',
    change: '+18%',
    changeType: 'positive',
    icon: Eye,
    href: '/analytics',
  },
  {
    title: 'Réservations',
    value: '89',
    change: '+12%',
    changeType: 'positive',
    icon: CalendarCheck,
    href: '/bookings',
  },
  {
    title: 'Taux de conversion',
    value: '7.2%',
    change: '-0.5%',
    changeType: 'negative',
    icon: TrendingUp,
    href: '/analytics',
  },
]

const recentBookings = [
  { id: '1', user: 'Marie D.', offer: 'Happy Hour -30%', date: '10 déc 2025', status: 'checked_in' },
  { id: '2', user: 'Pierre M.', offer: 'Menu du jour -20%', date: '10 déc 2025', status: 'confirmed' },
  { id: '3', user: 'Sophie L.', offer: 'Brunch weekend', date: '9 déc 2025', status: 'checked_in' },
  { id: '4', user: 'Jean P.', offer: 'Happy Hour -30%', date: '9 déc 2025', status: 'no_show' },
]

const topOffers = [
  { id: '1', title: 'Happy Hour -30%', views: 456, bookings: 34 },
  { id: '2', title: 'Menu du jour -20%', views: 312, bookings: 28 },
  { id: '3', title: 'Brunch weekend', views: 289, bookings: 15 },
]

const statusLabels: Record<string, { label: string; className: string }> = {
  checked_in: { label: 'Validé', className: 'bg-green-100 text-green-800' },
  confirmed: { label: 'Confirmé', className: 'bg-blue-100 text-blue-800' },
  no_show: { label: 'Absent', className: 'bg-red-100 text-red-800' },
  cancelled: { label: 'Annulé', className: 'bg-gray-100 text-gray-800' },
}

export function DashboardPage() {
  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold">Tableau de bord</h1>
          <p className="text-muted-foreground">
            Bienvenue ! Voici un aperçu de votre activité.
          </p>
        </div>
        <Link to="/offers/create">
          <Button>
            <Tag className="mr-2 h-4 w-4" />
            Nouvelle offre
          </Button>
        </Link>
      </div>

      {/* Stats cards */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <Link key={stat.title} to={stat.href}>
            <Card className="hover:shadow-md transition-shadow">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">
                  {stat.title}
                </CardTitle>
                <stat.icon className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stat.value}</div>
                <div className="flex items-center text-xs mt-1">
                  {stat.changeType === 'positive' ? (
                    <ArrowUpRight className="h-3 w-3 text-green-600 mr-1" />
                  ) : (
                    <ArrowDownRight className="h-3 w-3 text-red-600 mr-1" />
                  )}
                  <span
                    className={
                      stat.changeType === 'positive'
                        ? 'text-green-600'
                        : 'text-red-600'
                    }
                  >
                    {stat.change}
                  </span>
                  <span className="text-muted-foreground ml-1">
                    vs mois dernier
                  </span>
                </div>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>

      {/* Two columns */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Recent bookings */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle className="text-lg">Dernières réservations</CardTitle>
            <Link to="/bookings">
              <Button variant="ghost" size="sm">
                Voir tout
              </Button>
            </Link>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {recentBookings.map((booking) => (
                <div
                  key={booking.id}
                  className="flex items-center justify-between"
                >
                  <div className="flex items-center gap-3">
                    <div className="h-8 w-8 rounded-full bg-muted flex items-center justify-center">
                      <Users className="h-4 w-4 text-muted-foreground" />
                    </div>
                    <div>
                      <p className="font-medium text-sm">{booking.user}</p>
                      <p className="text-xs text-muted-foreground">
                        {booking.offer}
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <span
                      className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                        statusLabels[booking.status].className
                      }`}
                    >
                      {statusLabels[booking.status].label}
                    </span>
                    <p className="text-xs text-muted-foreground mt-1">
                      {booking.date}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Top offers */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle className="text-lg">Meilleures offres</CardTitle>
            <Link to="/offers">
              <Button variant="ghost" size="sm">
                Voir tout
              </Button>
            </Link>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {topOffers.map((offer, index) => (
                <div
                  key={offer.id}
                  className="flex items-center justify-between"
                >
                  <div className="flex items-center gap-3">
                    <div
                      className={`h-8 w-8 rounded-full flex items-center justify-center text-sm font-bold ${
                        index === 0
                          ? 'bg-yellow-100 text-yellow-800'
                          : index === 1
                          ? 'bg-gray-100 text-gray-800'
                          : 'bg-orange-100 text-orange-800'
                      }`}
                    >
                      {index + 1}
                    </div>
                    <div>
                      <p className="font-medium text-sm">{offer.title}</p>
                      <p className="text-xs text-muted-foreground">
                        {offer.views} vues
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="font-semibold text-sm">
                      {offer.bookings} réservations
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {((offer.bookings / offer.views) * 100).toFixed(1)}%
                      conversion
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Quick actions */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Actions rapides</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 sm:grid-cols-3">
            <Link to="/offers/create">
              <Button variant="outline" className="w-full h-auto py-4 flex-col">
                <Tag className="h-6 w-6 mb-2" />
                <span>Créer une offre</span>
              </Button>
            </Link>
            <Link to="/establishments">
              <Button variant="outline" className="w-full h-auto py-4 flex-col">
                <MapPin className="h-6 w-6 mb-2" />
                <span>Gérer les établissements</span>
              </Button>
            </Link>
            <Link to="/analytics">
              <Button variant="outline" className="w-full h-auto py-4 flex-col">
                <TrendingUp className="h-6 w-6 mb-2" />
                <span>Voir les statistiques</span>
              </Button>
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
