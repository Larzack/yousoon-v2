import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { formatDate } from '@/lib/utils'
import {
  Search,
  MoreHorizontal,
  Eye,
  Pause,
  Play,
  Trash2,
  Calendar,
  CheckCircle,
  XCircle,
  ChevronLeft,
  ChevronRight,
  AlertCircle,
  Building2,
} from 'lucide-react'

interface Offer {
  id: string
  title: string
  partner: {
    id: string
    name: string
  }
  establishment: string
  category: string
  discount: string
  status: 'draft' | 'pending' | 'active' | 'paused' | 'expired' | 'archived'
  bookingsCount: number
  checkinsCount: number
  startDate: string
  endDate: string
  createdAt: string
}

// Mock data
const mockOffers: Offer[] = [
  {
    id: '1',
    title: '-20% sur l\'addition',
    partner: { id: 'p1', name: 'Le Petit Bistrot' },
    establishment: 'Le Petit Bistrot - Marais',
    category: 'Restaurant',
    discount: '20%',
    status: 'active',
    bookingsCount: 156,
    checkinsCount: 142,
    startDate: '2024-11-01',
    endDate: '2024-12-31',
    createdAt: '2024-10-25',
  },
  {
    id: '2',
    title: 'Escape Room à 25€',
    partner: { id: 'p2', name: 'Escape Game Paris' },
    establishment: 'Escape Game Paris',
    category: 'Loisirs',
    discount: '25€',
    status: 'pending',
    bookingsCount: 0,
    checkinsCount: 0,
    startDate: '2024-12-15',
    endDate: '2025-03-31',
    createdAt: '2024-12-08',
  },
  {
    id: '3',
    title: 'Happy Hour -50%',
    partner: { id: 'p3', name: 'Le Bar du Coin' },
    establishment: 'Le Bar du Coin',
    category: 'Bar',
    discount: '50%',
    status: 'paused',
    bookingsCount: 234,
    checkinsCount: 198,
    startDate: '2024-09-01',
    endDate: '2024-12-31',
    createdAt: '2024-08-25',
  },
  {
    id: '4',
    title: 'Séance découverte gratuite',
    partner: { id: 'p4', name: 'Fitness Premium' },
    establishment: 'Fitness Premium - Opéra',
    category: 'Sport',
    discount: 'Gratuit',
    status: 'expired',
    bookingsCount: 89,
    checkinsCount: 76,
    startDate: '2024-10-01',
    endDate: '2024-11-30',
    createdAt: '2024-09-20',
  },
  {
    id: '5',
    title: 'Menu du midi à 15€',
    partner: { id: 'p1', name: 'Le Petit Bistrot' },
    establishment: 'Le Petit Bistrot - Bastille',
    category: 'Restaurant',
    discount: '15€',
    status: 'active',
    bookingsCount: 98,
    checkinsCount: 91,
    startDate: '2024-11-15',
    endDate: '2025-02-28',
    createdAt: '2024-11-10',
  },
]

function getStatusBadge(status: Offer['status']) {
  const styles: Record<Offer['status'], { bg: string; text: string; label: string }> = {
    draft: { bg: 'bg-gray-100', text: 'text-gray-700', label: 'Brouillon' },
    pending: { bg: 'bg-yellow-100', text: 'text-yellow-700', label: 'En attente' },
    active: { bg: 'bg-green-100', text: 'text-green-700', label: 'Active' },
    paused: { bg: 'bg-orange-100', text: 'text-orange-700', label: 'En pause' },
    expired: { bg: 'bg-red-100', text: 'text-red-700', label: 'Expirée' },
    archived: { bg: 'bg-gray-100', text: 'text-gray-700', label: 'Archivée' },
  }
  const style = styles[status]
  return (
    <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${style.bg} ${style.text}`}>
      {style.label}
    </span>
  )
}

export function OffersPage() {
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [categoryFilter, setCategoryFilter] = useState<string>('all')
  const [currentPage, setCurrentPage] = useState(1)

  const categories = ['all', 'Restaurant', 'Bar', 'Loisirs', 'Sport']

  const filteredOffers = mockOffers.filter((offer) => {
    const matchesSearch =
      offer.title.toLowerCase().includes(search.toLowerCase()) ||
      offer.partner.name.toLowerCase().includes(search.toLowerCase())
    const matchesStatus = statusFilter === 'all' || offer.status === statusFilter
    const matchesCategory = categoryFilter === 'all' || offer.category === categoryFilter
    return matchesSearch && matchesStatus && matchesCategory
  })

  const pendingCount = mockOffers.filter((o) => o.status === 'pending').length

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Offres</h1>
          <p className="text-muted-foreground">Gestion des offres de la plateforme</p>
        </div>
        {pendingCount > 0 && (
          <Button variant="outline" className="gap-2">
            <AlertCircle className="h-4 w-4 text-yellow-500" />
            {pendingCount} en attente de validation
          </Button>
        )}
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-5">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Total</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">1,892</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Actives</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">1,456</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">En attente</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">23</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">En pause</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-orange-600">87</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Expirées</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">326</div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col md:flex-row gap-4">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Rechercher par titre ou partenaire..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-9"
              />
            </div>
            <div className="flex gap-2 flex-wrap">
              {['all', 'active', 'pending', 'paused', 'expired'].map((status) => (
                <Button
                  key={status}
                  variant={statusFilter === status ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setStatusFilter(status)}
                >
                  {status === 'all' ? 'Toutes' : 
                   status === 'active' ? 'Actives' : 
                   status === 'pending' ? 'En attente' : 
                   status === 'paused' ? 'En pause' : 'Expirées'}
                </Button>
              ))}
            </div>
          </div>
          <div className="flex gap-2 mt-4">
            {categories.map((cat) => (
              <Button
                key={cat}
                variant={categoryFilter === cat ? 'secondary' : 'ghost'}
                size="sm"
                onClick={() => setCategoryFilter(cat)}
              >
                {cat === 'all' ? 'Toutes catégories' : cat}
              </Button>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Table */}
      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b">
                <tr>
                  <th className="text-left p-4 font-medium text-muted-foreground">Offre</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Partenaire</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Catégorie</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Réduction</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Statut</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Validité</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Réservations</th>
                  <th className="text-right p-4 font-medium text-muted-foreground">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {filteredOffers.map((offer) => (
                  <tr key={offer.id} className="hover:bg-gray-50">
                    <td className="p-4">
                      <div>
                        <p className="font-medium">{offer.title}</p>
                        <p className="text-sm text-muted-foreground flex items-center gap-1">
                          <Building2 className="h-3 w-3" />
                          {offer.establishment}
                        </p>
                      </div>
                    </td>
                    <td className="p-4">
                      <Link to={`/partners/${offer.partner.id}`} className="text-primary hover:underline">
                        {offer.partner.name}
                      </Link>
                    </td>
                    <td className="p-4">
                      <span className="px-2 py-1 bg-gray-100 rounded-full text-xs">
                        {offer.category}
                      </span>
                    </td>
                    <td className="p-4 font-medium">{offer.discount}</td>
                    <td className="p-4">{getStatusBadge(offer.status)}</td>
                    <td className="p-4">
                      <div className="text-sm">
                        <p>{formatDate(offer.startDate)}</p>
                        <p className="text-muted-foreground">→ {formatDate(offer.endDate)}</p>
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="text-sm">
                        <p className="flex items-center gap-1">
                          <Calendar className="h-3 w-3" />
                          {offer.bookingsCount}
                        </p>
                        <p className="flex items-center gap-1 text-muted-foreground">
                          <CheckCircle className="h-3 w-3" />
                          {offer.checkinsCount} check-ins
                        </p>
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="flex justify-end">
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem asChild>
                              <Link to={`/offers/${offer.id}`} className="flex items-center gap-2">
                                <Eye className="h-4 w-4" />
                                Voir détail
                              </Link>
                            </DropdownMenuItem>
                            {offer.status === 'pending' && (
                              <>
                                <DropdownMenuItem className="gap-2 text-green-600">
                                  <CheckCircle className="h-4 w-4" />
                                  Approuver
                                </DropdownMenuItem>
                                <DropdownMenuItem className="gap-2 text-red-600">
                                  <XCircle className="h-4 w-4" />
                                  Rejeter
                                </DropdownMenuItem>
                              </>
                            )}
                            {offer.status === 'active' && (
                              <DropdownMenuItem className="gap-2 text-orange-600">
                                <Pause className="h-4 w-4" />
                                Mettre en pause
                              </DropdownMenuItem>
                            )}
                            {offer.status === 'paused' && (
                              <DropdownMenuItem className="gap-2 text-green-600">
                                <Play className="h-4 w-4" />
                                Réactiver
                              </DropdownMenuItem>
                            )}
                            <DropdownMenuItem className="gap-2 text-red-600">
                              <Trash2 className="h-4 w-4" />
                              Archiver
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          <div className="flex items-center justify-between p-4 border-t">
            <p className="text-sm text-muted-foreground">
              Affichage de {filteredOffers.length} offres
            </p>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="icon"
                disabled={currentPage === 1}
                onClick={() => setCurrentPage(currentPage - 1)}
              >
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <span className="text-sm">Page {currentPage}</span>
              <Button
                variant="outline"
                size="icon"
                onClick={() => setCurrentPage(currentPage + 1)}
              >
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
