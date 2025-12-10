import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { formatDate } from '@/lib/utils'
import {
  ArrowLeft,
  Star,
  Flag,
  User,
  CheckCircle,
  Trash2,
  AlertTriangle,
  MessageSquare,
  Calendar,
  Eye,
} from 'lucide-react'

interface ReportedReview {
  id: string
  user: {
    id: string
    firstName: string
    lastName: string
    email: string
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
  status: 'reported'
  reportCount: number
  reports: Array<{
    userId: string
    userName: string
    reason: string
    reportedAt: string
  }>
  isVerifiedPurchase: boolean
  createdAt: string
}

// Mock data
const mockReportedReviews: ReportedReview[] = [
  {
    id: '1',
    user: { id: 'u3', firstName: 'Pierre', lastName: 'Dubois', email: 'pierre@example.com' },
    offer: { id: 'o3', title: 'Happy Hour -50%' },
    partner: { id: 'p3', name: 'Le Bar du Coin' },
    rating: 2,
    title: 'Décevant',
    content: 'Service lent et le personnel n\'était pas au courant de la promotion Yousoon. J\'ai dû insister pour avoir la réduction.',
    status: 'reported',
    reportCount: 2,
    reports: [
      { userId: 'p3', userName: 'Le Bar du Coin (Partenaire)', reason: 'Informations inexactes, notre personnel est formé', reportedAt: '2024-12-07' },
      { userId: 'u10', userName: 'Client anonyme', reason: 'Contenu diffamatoire', reportedAt: '2024-12-08' },
    ],
    isVerifiedPurchase: true,
    createdAt: '2024-12-06',
  },
  {
    id: '2',
    user: { id: 'u6', firstName: 'Emma', lastName: 'Garcia', email: 'emma@example.com' },
    offer: { id: 'o5', title: 'Spa détente' },
    partner: { id: 'p5', name: 'Spa Zen' },
    rating: 1,
    title: 'ARNAQUE !!!',
    content: 'NE VENEZ PAS ICI !!! C\'est une arnaque totale, le personnel est impoli et l\'endroit est sale. Je vais porter plainte !!!!!',
    status: 'reported',
    reportCount: 5,
    reports: [
      { userId: 'p5', userName: 'Spa Zen (Partenaire)', reason: 'Propos diffamatoires et mensongers', reportedAt: '2024-12-05' },
      { userId: 'u11', userName: 'Alice D.', reason: 'Langage inapproprié', reportedAt: '2024-12-05' },
      { userId: 'u12', userName: 'Marc L.', reason: 'Fausses accusations', reportedAt: '2024-12-06' },
      { userId: 'u13', userName: 'Julie M.', reason: 'Spam / Abus', reportedAt: '2024-12-06' },
      { userId: 'u14', userName: 'Thomas R.', reason: 'Contenu inapproprié', reportedAt: '2024-12-07' },
    ],
    isVerifiedPurchase: false,
    createdAt: '2024-12-05',
  },
  {
    id: '3',
    user: { id: 'u7', firstName: 'Antoine', lastName: 'Petit', email: 'antoine@example.com' },
    offer: { id: 'o6', title: 'Cours de cuisine' },
    partner: { id: 'p6', name: 'Atelier Gourmand' },
    rating: 3,
    title: 'Moyen',
    content: 'Le cours était correct mais trop basique pour le prix. Le chef était sympathique mais les recettes étaient trop simples.',
    status: 'reported',
    reportCount: 1,
    reports: [
      { userId: 'p6', userName: 'Atelier Gourmand (Partenaire)', reason: 'Ce client n\'a jamais participé à notre cours', reportedAt: '2024-12-08' },
    ],
    isVerifiedPurchase: true,
    createdAt: '2024-12-07',
  },
]

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

export function ReportedReviewsPage() {
  const navigate = useNavigate()
  const [selectedReview, setSelectedReview] = useState<ReportedReview | null>(
    mockReportedReviews[0] || null
  )

  const handleApprove = (reviewId: string) => {
    // In real app, call API
    alert(`Avis ${reviewId} maintenu (signalements ignorés)`)
  }

  const handleDelete = (reviewId: string) => {
    // In real app, call API
    alert(`Avis ${reviewId} supprimé`)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => navigate('/reviews')}>
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div className="flex-1">
          <h1 className="text-2xl font-bold flex items-center gap-2">
            <Flag className="h-6 w-6 text-orange-500" />
            Avis signalés
          </h1>
          <p className="text-muted-foreground">
            {mockReportedReviews.length} avis nécessitent une modération
          </p>
        </div>
      </div>

      {/* Alert */}
      <div className="flex items-center gap-3 p-4 bg-orange-50 border border-orange-200 rounded-lg">
        <AlertTriangle className="h-5 w-5 text-orange-600 flex-shrink-0" />
        <div>
          <p className="font-medium text-orange-800">Avis signalés par des utilisateurs ou partenaires</p>
          <p className="text-sm text-orange-600">
            Examinez chaque avis et décidez s'il doit être maintenu ou supprimé.
          </p>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {/* List */}
        <div className="space-y-3">
          {mockReportedReviews.map((review) => (
            <Card
              key={review.id}
              className={`cursor-pointer transition-colors ${
                selectedReview?.id === review.id
                  ? 'border-primary ring-1 ring-primary'
                  : 'hover:border-gray-300'
              }`}
              onClick={() => setSelectedReview(review)}
            >
              <CardContent className="pt-4">
                <div className="flex items-start gap-3">
                  <div className="h-10 w-10 rounded-full bg-gray-100 flex items-center justify-center flex-shrink-0">
                    <User className="h-5 w-5 text-gray-500" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="font-medium truncate">
                        {review.user.firstName} {review.user.lastName}
                      </span>
                      <StarRating rating={review.rating} />
                    </div>
                    <p className="text-sm text-muted-foreground truncate">
                      {review.partner.name}
                    </p>
                    <p className="text-sm mt-1 line-clamp-2">{review.content}</p>
                    <div className="flex items-center gap-2 mt-2">
                      <span className="inline-flex items-center gap-1 text-xs text-orange-600 bg-orange-100 px-2 py-0.5 rounded-full">
                        <Flag className="h-3 w-3" />
                        {review.reportCount} signalement{review.reportCount > 1 ? 's' : ''}
                      </span>
                      {!review.isVerifiedPurchase && (
                        <span className="text-xs text-red-600">Non vérifié</span>
                      )}
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Detail Panel */}
        {selectedReview && (
          <div className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Détail de l'avis</CardTitle>
                <CardDescription>
                  Publié le {formatDate(selectedReview.createdAt)}
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center gap-3">
                  <div className="h-12 w-12 rounded-full bg-gray-100 flex items-center justify-center">
                    <User className="h-6 w-6 text-gray-500" />
                  </div>
                  <div>
                    <Link
                      to={`/users/${selectedReview.user.id}`}
                      className="font-medium text-primary hover:underline"
                    >
                      {selectedReview.user.firstName} {selectedReview.user.lastName}
                    </Link>
                    <p className="text-sm text-muted-foreground">{selectedReview.user.email}</p>
                  </div>
                </div>

                <div className="flex items-center gap-4">
                  <StarRating rating={selectedReview.rating} />
                  {selectedReview.isVerifiedPurchase ? (
                    <span className="inline-flex items-center gap-1 text-xs text-green-600">
                      <CheckCircle className="h-3 w-3" />
                      Achat vérifié
                    </span>
                  ) : (
                    <span className="inline-flex items-center gap-1 text-xs text-red-600">
                      <AlertTriangle className="h-3 w-3" />
                      Non vérifié
                    </span>
                  )}
                </div>

                <div>
                  <p className="text-sm text-muted-foreground">Offre concernée</p>
                  <Link
                    to={`/offers/${selectedReview.offer.id}`}
                    className="text-primary hover:underline"
                  >
                    {selectedReview.offer.title}
                  </Link>
                  <span className="text-muted-foreground"> • </span>
                  <Link
                    to={`/partners/${selectedReview.partner.id}`}
                    className="hover:underline"
                  >
                    {selectedReview.partner.name}
                  </Link>
                </div>

                {selectedReview.title && (
                  <p className="font-bold text-lg">{selectedReview.title}</p>
                )}
                <p className="bg-gray-50 p-4 rounded-lg">{selectedReview.content}</p>
              </CardContent>
            </Card>

            <Card className="border-orange-200">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-orange-700">
                  <Flag className="h-5 w-5" />
                  Signalements ({selectedReview.reportCount})
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {selectedReview.reports.map((report, i) => (
                    <div key={i} className="flex items-start gap-3 pb-3 border-b last:border-0 last:pb-0">
                      <MessageSquare className="h-4 w-4 text-muted-foreground mt-1" />
                      <div>
                        <div className="flex items-center gap-2">
                          <span className="font-medium text-sm">{report.userName}</span>
                          <span className="text-xs text-muted-foreground flex items-center gap-1">
                            <Calendar className="h-3 w-3" />
                            {formatDate(report.reportedAt)}
                          </span>
                        </div>
                        <p className="text-sm text-muted-foreground mt-1">{report.reason}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Actions</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <Button
                  variant="outline"
                  className="w-full gap-2"
                  onClick={() => handleApprove(selectedReview.id)}
                >
                  <CheckCircle className="h-4 w-4 text-green-600" />
                  Maintenir l'avis (ignorer les signalements)
                </Button>
                <Button
                  variant="destructive"
                  className="w-full gap-2"
                  onClick={() => handleDelete(selectedReview.id)}
                >
                  <Trash2 className="h-4 w-4" />
                  Supprimer l'avis
                </Button>
                <Button variant="ghost" className="w-full gap-2" asChild>
                  <Link to={`/users/${selectedReview.user.id}`}>
                    <Eye className="h-4 w-4" />
                    Voir le profil de l'auteur
                  </Link>
                </Button>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </div>
  )
}
