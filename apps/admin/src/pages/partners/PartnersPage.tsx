import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu'
import { cn, formatDate, getInitials } from '@/lib/utils'
import {
  Search,
  MoreHorizontal,
  Eye,
  Ban,
  CheckCircle,
  XCircle,
  Building2,
  MapPin,
  Tag,
  Star,
  ChevronLeft,
  ChevronRight,
  AlertCircle,
} from 'lucide-react'

interface Partner {
  id: string
  companyName: string
  tradeName: string
  logo?: string
  siret: string
  category: string
  status: 'pending' | 'active' | 'suspended'
  establishmentsCount: number
  offersCount: number
  avgRating: number
  reviewsCount: number
  createdAt: string
  contact: {
    firstName: string
    lastName: string
    email: string
  }
}

// Mock data
const mockPartners: Partner[] = [
  {
    id: '1',
    companyName: 'Le Petit Bistrot SARL',
    tradeName: 'Le Petit Bistrot',
    siret: '123 456 789 00012',
    category: 'Restaurant',
    status: 'active',
    establishmentsCount: 2,
    offersCount: 5,
    avgRating: 4.5,
    reviewsCount: 128,
    createdAt: '2024-06-15',
    contact: {
      firstName: 'Pierre',
      lastName: 'Martin',
      email: 'pierre@lepetitbistrot.fr',
    },
  },
  {
    id: '2',
    companyName: 'Escape Game Paris SAS',
    tradeName: 'Escape Game Paris',
    siret: '987 654 321 00034',
    category: 'Loisirs',
    status: 'pending',
    establishmentsCount: 1,
    offersCount: 0,
    avgRating: 0,
    reviewsCount: 0,
    createdAt: '2024-12-08',
    contact: {
      firstName: 'Sophie',
      lastName: 'Durand',
      email: 'contact@escapegameparis.fr',
    },
  },
  {
    id: '3',
    companyName: 'Le Bar du Coin',
    tradeName: 'Le Bar du Coin',
    siret: '456 789 123 00056',
    category: 'Bar',
    status: 'active',
    establishmentsCount: 1,
    offersCount: 3,
    avgRating: 4.2,
    reviewsCount: 87,
    createdAt: '2024-08-20',
    contact: {
      firstName: 'Marc',
      lastName: 'Leblanc',
      email: 'marc@barducoin.fr',
    },
  },
  {
    id: '4',
    companyName: 'Club Fitness Premium',
    tradeName: 'Fitness Premium',
    siret: '789 123 456 00078',
    category: 'Sport',
    status: 'suspended',
    establishmentsCount: 3,
    offersCount: 2,
    avgRating: 3.8,
    reviewsCount: 45,
    createdAt: '2024-03-10',
    contact: {
      firstName: 'Julie',
      lastName: 'Moreau',
      email: 'julie@fitnesspremium.fr',
    },
  },
  {
    id: '5',
    companyName: 'Restaurant Le Gourmet',
    tradeName: 'Le Gourmet',
    siret: '321 654 987 00090',
    category: 'Restaurant',
    status: 'pending',
    establishmentsCount: 1,
    offersCount: 0,
    avgRating: 0,
    reviewsCount: 0,
    createdAt: '2024-12-09',
    contact: {
      firstName: 'Antoine',
      lastName: 'Bernard',
      email: 'antoine@legourmet.fr',
    },
  },
]

function getStatusBadge(status: Partner['status']) {
  switch (status) {
    case 'active':
      return (
        <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-700">
          <CheckCircle className="h-3 w-3" />
          Actif
        </span>
      )
    case 'pending':
      return (
        <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-700">
          <AlertCircle className="h-3 w-3" />
          En attente
        </span>
      )
    case 'suspended':
      return (
        <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium bg-red-100 text-red-700">
          <XCircle className="h-3 w-3" />
          Suspendu
        </span>
      )
  }
}

export function PartnersPage() {
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [categoryFilter, setCategoryFilter] = useState<string>('all')
  const [currentPage, setCurrentPage] = useState(1)

  const categories = ['all', 'Restaurant', 'Bar', 'Loisirs', 'Sport']

  const filteredPartners = mockPartners.filter((partner) => {
    const matchesSearch =
      partner.companyName.toLowerCase().includes(search.toLowerCase()) ||
      partner.tradeName.toLowerCase().includes(search.toLowerCase()) ||
      partner.siret.includes(search)
    const matchesStatus = statusFilter === 'all' || partner.status === statusFilter
    const matchesCategory = categoryFilter === 'all' || partner.category === categoryFilter
    return matchesSearch && matchesStatus && matchesCategory
  })

  const pendingCount = mockPartners.filter((p) => p.status === 'pending').length

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Partenaires</h1>
          <p className="text-muted-foreground">Gestion des partenaires de la plateforme</p>
        </div>
        {pendingCount > 0 && (
          <Link to="/partners/pending">
            <Button variant="outline" className="gap-2">
              <AlertCircle className="h-4 w-4 text-yellow-500" />
              {pendingCount} en attente
            </Button>
          </Link>
        )}
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Total</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">234</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Actifs</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">218</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">En attente</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">12</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Suspendus</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">4</div>
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
                placeholder="Rechercher par nom ou SIRET..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-9"
              />
            </div>
            <div className="flex gap-2 flex-wrap">
              {['all', 'active', 'pending', 'suspended'].map((status) => (
                <Button
                  key={status}
                  variant={statusFilter === status ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setStatusFilter(status)}
                >
                  {status === 'all' ? 'Tous' : status === 'active' ? 'Actifs' : status === 'pending' ? 'En attente' : 'Suspendus'}
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
                  <th className="text-left p-4 font-medium text-muted-foreground">Partenaire</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Catégorie</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Statut</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Établissements</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Offres</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Note</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Inscrit le</th>
                  <th className="text-right p-4 font-medium text-muted-foreground">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {filteredPartners.map((partner) => (
                  <tr key={partner.id} className="hover:bg-gray-50">
                    <td className="p-4">
                      <div className="flex items-center gap-3">
                        <Avatar>
                          <AvatarImage src={partner.logo} />
                          <AvatarFallback className="bg-purple-100 text-purple-600">
                            {getInitials(partner.tradeName)}
                          </AvatarFallback>
                        </Avatar>
                        <div>
                          <p className="font-medium">{partner.tradeName}</p>
                          <p className="text-sm text-muted-foreground">{partner.contact.email}</p>
                        </div>
                      </div>
                    </td>
                    <td className="p-4">
                      <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium bg-gray-100">
                        {partner.category}
                      </span>
                    </td>
                    <td className="p-4">{getStatusBadge(partner.status)}</td>
                    <td className="p-4">
                      <div className="flex items-center gap-1 text-muted-foreground">
                        <MapPin className="h-4 w-4" />
                        {partner.establishmentsCount}
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="flex items-center gap-1 text-muted-foreground">
                        <Tag className="h-4 w-4" />
                        {partner.offersCount}
                      </div>
                    </td>
                    <td className="p-4">
                      {partner.avgRating > 0 ? (
                        <div className="flex items-center gap-1">
                          <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                          <span className="font-medium">{partner.avgRating}</span>
                          <span className="text-sm text-muted-foreground">({partner.reviewsCount})</span>
                        </div>
                      ) : (
                        <span className="text-muted-foreground text-sm">-</span>
                      )}
                    </td>
                    <td className="p-4 text-muted-foreground">{formatDate(partner.createdAt)}</td>
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
                              <Link to={`/partners/${partner.id}`} className="flex items-center gap-2">
                                <Eye className="h-4 w-4" />
                                Voir détail
                              </Link>
                            </DropdownMenuItem>
                            {partner.status === 'pending' && (
                              <>
                                <DropdownMenuItem className="gap-2 text-green-600">
                                  <CheckCircle className="h-4 w-4" />
                                  Valider
                                </DropdownMenuItem>
                                <DropdownMenuItem className="gap-2 text-red-600">
                                  <XCircle className="h-4 w-4" />
                                  Rejeter
                                </DropdownMenuItem>
                              </>
                            )}
                            {partner.status === 'active' && (
                              <DropdownMenuItem className="gap-2 text-red-600">
                                <Ban className="h-4 w-4" />
                                Suspendre
                              </DropdownMenuItem>
                            )}
                            {partner.status === 'suspended' && (
                              <DropdownMenuItem className="gap-2 text-green-600">
                                <CheckCircle className="h-4 w-4" />
                                Réactiver
                              </DropdownMenuItem>
                            )}
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
              Affichage de {filteredPartners.length} partenaires
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
