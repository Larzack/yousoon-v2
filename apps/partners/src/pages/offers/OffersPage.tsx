import { useState } from 'react'
import { Link } from 'react-router-dom'
import {
  Plus,
  Search,
  Filter,
  MoreHorizontal,
  Eye,
  Edit,
  Copy,
  Archive,
  Trash2,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

// Mock data
const mockOffers = [
  {
    id: '1',
    title: 'Happy Hour -30%',
    establishment: 'Le Petit Bistrot',
    category: 'Bar',
    discount: '-30%',
    status: 'active',
    views: 456,
    bookings: 34,
    validUntil: '31 déc 2025',
    image: 'https://images.unsplash.com/photo-1551024709-8f23befc6f87?w=100&h=100&fit=crop',
  },
  {
    id: '2',
    title: 'Menu du jour',
    establishment: 'Le Petit Bistrot',
    category: 'Restaurant',
    discount: '-20%',
    status: 'active',
    views: 312,
    bookings: 28,
    validUntil: '31 déc 2025',
    image: 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=100&h=100&fit=crop',
  },
  {
    id: '3',
    title: 'Brunch weekend',
    establishment: 'Le Petit Bistrot',
    category: 'Restaurant',
    discount: '-25%',
    status: 'paused',
    views: 289,
    bookings: 15,
    validUntil: '30 jan 2026',
    image: 'https://images.unsplash.com/photo-1504674900247-0877df9cc836?w=100&h=100&fit=crop',
  },
  {
    id: '4',
    title: 'Soirée Jazz',
    establishment: 'Le Petit Bistrot',
    category: 'Événement',
    discount: '1 = 1 offert',
    status: 'draft',
    views: 0,
    bookings: 0,
    validUntil: '15 jan 2026',
    image: 'https://images.unsplash.com/photo-1514525253161-7a46d19cd819?w=100&h=100&fit=crop',
  },
]

const statusConfig: Record<string, { label: string; className: string }> = {
  active: { label: 'Active', className: 'bg-green-100 text-green-800' },
  paused: { label: 'En pause', className: 'bg-yellow-100 text-yellow-800' },
  draft: { label: 'Brouillon', className: 'bg-gray-100 text-gray-800' },
  expired: { label: 'Expirée', className: 'bg-red-100 text-red-800' },
}

export function OffersPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [statusFilter, setStatusFilter] = useState<string | null>(null)

  const filteredOffers = mockOffers.filter((offer) => {
    const matchesSearch = offer.title
      .toLowerCase()
      .includes(searchQuery.toLowerCase())
    const matchesStatus = !statusFilter || offer.status === statusFilter
    return matchesSearch && matchesStatus
  })

  const stats = {
    total: mockOffers.length,
    active: mockOffers.filter((o) => o.status === 'active').length,
    paused: mockOffers.filter((o) => o.status === 'paused').length,
    draft: mockOffers.filter((o) => o.status === 'draft').length,
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold">Mes offres</h1>
          <p className="text-muted-foreground">
            Gérez vos offres et réductions
          </p>
        </div>
        <Link to="/offers/create">
          <Button>
            <Plus className="mr-2 h-4 w-4" />
            Nouvelle offre
          </Button>
        </Link>
      </div>

      {/* Stats */}
      <div className="grid gap-4 grid-cols-2 lg:grid-cols-4">
        <Card
          className={`cursor-pointer ${
            statusFilter === null ? 'ring-2 ring-primary' : ''
          }`}
          onClick={() => setStatusFilter(null)}
        >
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.total}</div>
          </CardContent>
        </Card>
        <Card
          className={`cursor-pointer ${
            statusFilter === 'active' ? 'ring-2 ring-primary' : ''
          }`}
          onClick={() => setStatusFilter(statusFilter === 'active' ? null : 'active')}
        >
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Actives
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{stats.active}</div>
          </CardContent>
        </Card>
        <Card
          className={`cursor-pointer ${
            statusFilter === 'paused' ? 'ring-2 ring-primary' : ''
          }`}
          onClick={() => setStatusFilter(statusFilter === 'paused' ? null : 'paused')}
        >
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              En pause
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">{stats.paused}</div>
          </CardContent>
        </Card>
        <Card
          className={`cursor-pointer ${
            statusFilter === 'draft' ? 'ring-2 ring-primary' : ''
          }`}
          onClick={() => setStatusFilter(statusFilter === 'draft' ? null : 'draft')}
        >
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Brouillons
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-gray-600">{stats.draft}</div>
          </CardContent>
        </Card>
      </div>

      {/* Search and filters */}
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Rechercher une offre..."
            className="pl-9"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
        <Button variant="outline">
          <Filter className="mr-2 h-4 w-4" />
          Filtres
        </Button>
      </div>

      {/* Offers list */}
      <div className="space-y-4">
        {filteredOffers.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">Aucune offre trouvée</p>
              <Link to="/offers/create">
                <Button className="mt-4">Créer ma première offre</Button>
              </Link>
            </CardContent>
          </Card>
        ) : (
          filteredOffers.map((offer) => (
            <Card key={offer.id} className="hover:shadow-md transition-shadow">
              <CardContent className="p-4">
                <div className="flex items-center gap-4">
                  {/* Image */}
                  <img
                    src={offer.image}
                    alt={offer.title}
                    className="w-16 h-16 rounded-lg object-cover flex-shrink-0"
                  />

                  {/* Info */}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <Link
                        to={`/offers/${offer.id}`}
                        className="font-semibold hover:underline truncate"
                      >
                        {offer.title}
                      </Link>
                      <span
                        className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${
                          statusConfig[offer.status].className
                        }`}
                      >
                        {statusConfig[offer.status].label}
                      </span>
                    </div>
                    <p className="text-sm text-muted-foreground">
                      {offer.establishment} • {offer.category}
                    </p>
                    <div className="flex items-center gap-4 mt-1 text-sm">
                      <span className="font-semibold text-primary">
                        {offer.discount}
                      </span>
                      <span className="text-muted-foreground">
                        <Eye className="inline h-3 w-3 mr-1" />
                        {offer.views} vues
                      </span>
                      <span className="text-muted-foreground">
                        {offer.bookings} réservations
                      </span>
                    </div>
                  </div>

                  {/* Valid until */}
                  <div className="hidden md:block text-right">
                    <p className="text-sm text-muted-foreground">Valide jusqu'au</p>
                    <p className="text-sm font-medium">{offer.validUntil}</p>
                  </div>

                  {/* Actions */}
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="icon">
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem asChild>
                        <Link to={`/offers/${offer.id}`}>
                          <Eye className="mr-2 h-4 w-4" />
                          Voir
                        </Link>
                      </DropdownMenuItem>
                      <DropdownMenuItem>
                        <Edit className="mr-2 h-4 w-4" />
                        Modifier
                      </DropdownMenuItem>
                      <DropdownMenuItem>
                        <Copy className="mr-2 h-4 w-4" />
                        Dupliquer
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem>
                        <Archive className="mr-2 h-4 w-4" />
                        Archiver
                      </DropdownMenuItem>
                      <DropdownMenuItem className="text-destructive">
                        <Trash2 className="mr-2 h-4 w-4" />
                        Supprimer
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </div>
    </div>
  )
}
