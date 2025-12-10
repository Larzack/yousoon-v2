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
  Trash2,
  Star,
  Flag,
  User,
  ChevronLeft,
  ChevronRight,
  CheckCircle,
  AlertTriangle,
} from 'lucide-react'

interface Review {
  id: string
  user: {
    id: string
    firstName: string
    lastName: string
  }
  offer: {
    id: string
    title: string
  }
  partner: {
    id: string
    name: string
  }
  rating: number
  title: string | null
  content: string
  status: 'pending' | 'approved' | 'rejected' | 'reported'
  reportCount: number
  isVerifiedPurchase: boolean
  createdAt: string
}

// Mock data
const mockReviews: Review[] = [
  {
    id: '1',
    user: { id: 'u1', firstName: 'Jean', lastName: 'Dupont' },
    offer: { id: 'o1', title: '-20% sur l\'addition' },
    partner: { id: 'p1', name: 'Le Petit Bistrot' },
    rating: 5,
    title: 'Excellent !',
    content: 'Super expérience, le personnel était très accueillant et la réduction a bien été appliquée. Je recommande vivement !',
    status: 'approved',
    reportCount: 0,
    isVerifiedPurchase: true,
    createdAt: '2024-12-08',
  },
  {
    id: '2',
    user: { id: 'u2', firstName: 'Marie', lastName: 'Martin' },
    offer: { id: 'o2', title: 'Escape Room à 25€' },
    partner: { id: 'p2', name: 'Escape Game Paris' },
    rating: 4,
    title: 'Très bien',
    content: 'Bonne expérience, l\'escape room était bien conçue. Seul bémol, il faisait un peu chaud.',
    status: 'approved',
    reportCount: 0,
    isVerifiedPurchase: true,
    createdAt: '2024-12-07',
  },
  {
    id: '3',
    user: { id: 'u3', firstName: 'Pierre', lastName: 'Dubois' },
    offer: { id: 'o3', title: 'Happy Hour -50%' },
    partner: { id: 'p3', name: 'Le Bar du Coin' },
    rating: 2,
    title: 'Décevant',
    content: 'Service lent et le personnel n\'était pas au courant de la promotion Yousoon. J\'ai dû insister pour avoir la réduction.',
    status: 'reported',
    reportCount: 2,
    isVerifiedPurchase: true,
    createdAt: '2024-12-06',
  },
  {
    id: '4',
    user: { id: 'u4', firstName: 'Sophie', lastName: 'Leroy' },
    offer: { id: 'o1', title: '-20% sur l\'addition' },
    partner: { id: 'p1', name: 'Le Petit Bistrot' },
    rating: 1,
    title: null,
    content: 'Arnaque totale, la réduction n\'a pas été appliquée et le gérant était désagréable !!! NUL',
    status: 'pending',
    reportCount: 0,
    isVerifiedPurchase: false,
    createdAt: '2024-12-09',
  },
  {
    id: '5',
    user: { id: 'u5', firstName: 'Lucas', lastName: 'Bernard' },
    offer: { id: 'o4', title: 'Menu découverte' },
    partner: { id: 'p4', name: 'Restaurant Gourmet' },
    rating: 5,
    title: 'Parfait',
    content: 'Meilleure expérience culinaire de ma vie. Les plats étaient délicieux et le service impeccable.',
    status: 'approved',
    reportCount: 0,
    isVerifiedPurchase: true,
    createdAt: '2024-12-05',
  },
]

function getStatusBadge(status: Review['status']) {
  const styles = {
    pending: { bg: 'bg-yellow-100', text: 'text-yellow-700', label: 'En attente' },
    approved: { bg: 'bg-green-100', text: 'text-green-700', label: 'Approuvé' },
    rejected: { bg: 'bg-red-100', text: 'text-red-700', label: 'Rejeté' },
    reported: { bg: 'bg-orange-100', text: 'text-orange-700', label: 'Signalé' },
  }
  const style = styles[status]
  return (
    <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${style.bg} ${style.text}`}>
      {style.label}
    </span>
  )
}

function StarRating({ rating }: { rating: number }) {
  return (
    <div className="flex gap-0.5">
      {[1, 2, 3, 4, 5].map((star) => (
        <Star
          key={star}
          className={`h-4 w-4 ${star <= rating ? 'fill-yellow-400 text-yellow-400' : 'text-gray-300'}`}
        />
      ))}
    </div>
  )
}

export function ReviewsPage() {
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [ratingFilter, setRatingFilter] = useState<number | null>(null)
  const [currentPage, setCurrentPage] = useState(1)

  const filteredReviews = mockReviews.filter((review) => {
    const matchesSearch =
      review.content.toLowerCase().includes(search.toLowerCase()) ||
      review.user.firstName.toLowerCase().includes(search.toLowerCase()) ||
      review.user.lastName.toLowerCase().includes(search.toLowerCase()) ||
      review.partner.name.toLowerCase().includes(search.toLowerCase())
    const matchesStatus = statusFilter === 'all' || review.status === statusFilter
    const matchesRating = ratingFilter === null || review.rating === ratingFilter
    return matchesSearch && matchesStatus && matchesRating
  })

  const reportedCount = mockReviews.filter((r) => r.status === 'reported').length
  const pendingCount = mockReviews.filter((r) => r.status === 'pending').length

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Avis</h1>
          <p className="text-muted-foreground">Modération des avis utilisateurs</p>
        </div>
        <div className="flex gap-2">
          {reportedCount > 0 && (
            <Link to="/reviews/reported">
              <Button variant="outline" className="gap-2">
                <Flag className="h-4 w-4 text-orange-500" />
                {reportedCount} signalés
              </Button>
            </Link>
          )}
        </div>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-5">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Total</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">12,456</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">En attente</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">{pendingCount}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Signalés</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-orange-600">{reportedCount}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Note moyenne</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold flex items-center gap-1">
              <Star className="h-5 w-5 fill-yellow-400 text-yellow-400" />
              4.2
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Ce mois</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">+234</div>
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
                placeholder="Rechercher dans les avis..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-9"
              />
            </div>
            <div className="flex gap-2 flex-wrap">
              {['all', 'pending', 'reported', 'approved', 'rejected'].map((status) => (
                <Button
                  key={status}
                  variant={statusFilter === status ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setStatusFilter(status)}
                >
                  {status === 'all' ? 'Tous' :
                   status === 'pending' ? 'En attente' :
                   status === 'reported' ? 'Signalés' :
                   status === 'approved' ? 'Approuvés' : 'Rejetés'}
                </Button>
              ))}
            </div>
          </div>
          <div className="flex gap-2 mt-4">
            <span className="text-sm text-muted-foreground py-1">Note :</span>
            {[null, 5, 4, 3, 2, 1].map((rating) => (
              <Button
                key={rating ?? 'all'}
                variant={ratingFilter === rating ? 'secondary' : 'ghost'}
                size="sm"
                onClick={() => setRatingFilter(rating)}
                className="gap-1"
              >
                {rating === null ? 'Toutes' : (
                  <>
                    {rating}
                    <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
                  </>
                )}
              </Button>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Reviews List */}
      <div className="space-y-4">
        {filteredReviews.map((review) => (
          <Card key={review.id} className={review.status === 'reported' ? 'border-orange-200' : ''}>
            <CardContent className="pt-6">
              <div className="flex items-start gap-4">
                <div className="h-10 w-10 rounded-full bg-gray-100 flex items-center justify-center flex-shrink-0">
                  <User className="h-5 w-5 text-gray-500" />
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 flex-wrap">
                    <Link
                      to={`/users/${review.user.id}`}
                      className="font-medium hover:text-primary hover:underline"
                    >
                      {review.user.firstName} {review.user.lastName}
                    </Link>
                    <StarRating rating={review.rating} />
                    {getStatusBadge(review.status)}
                    {review.isVerifiedPurchase && (
                      <span className="inline-flex items-center gap-1 text-xs text-green-600">
                        <CheckCircle className="h-3 w-3" />
                        Achat vérifié
                      </span>
                    )}
                    {review.reportCount > 0 && (
                      <span className="inline-flex items-center gap-1 text-xs text-orange-600">
                        <AlertTriangle className="h-3 w-3" />
                        {review.reportCount} signalement{review.reportCount > 1 ? 's' : ''}
                      </span>
                    )}
                  </div>
                  <div className="flex items-center gap-2 text-sm text-muted-foreground mt-1">
                    <span>sur</span>
                    <Link
                      to={`/offers/${review.offer.id}`}
                      className="text-primary hover:underline"
                    >
                      {review.offer.title}
                    </Link>
                    <span>•</span>
                    <Link
                      to={`/partners/${review.partner.id}`}
                      className="hover:underline"
                    >
                      {review.partner.name}
                    </Link>
                    <span>•</span>
                    <span>{formatDate(review.createdAt)}</span>
                  </div>
                  {review.title && (
                    <p className="font-medium mt-3">{review.title}</p>
                  )}
                  <p className="mt-2 text-gray-700">{review.content}</p>
                </div>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon">
                      <MoreHorizontal className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem asChild>
                      <Link to={`/reviews/${review.id}`} className="flex items-center gap-2">
                        <Eye className="h-4 w-4" />
                        Voir détail
                      </Link>
                    </DropdownMenuItem>
                    {(review.status === 'pending' || review.status === 'reported') && (
                      <>
                        <DropdownMenuItem className="gap-2 text-green-600">
                          <CheckCircle className="h-4 w-4" />
                          Approuver
                        </DropdownMenuItem>
                        <DropdownMenuItem className="gap-2 text-red-600">
                          <Trash2 className="h-4 w-4" />
                          Supprimer
                        </DropdownMenuItem>
                      </>
                    )}
                    {review.status === 'approved' && (
                      <DropdownMenuItem className="gap-2 text-red-600">
                        <Trash2 className="h-4 w-4" />
                        Supprimer
                      </DropdownMenuItem>
                    )}
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Pagination */}
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          Affichage de {filteredReviews.length} avis
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
    </div>
  )
}
