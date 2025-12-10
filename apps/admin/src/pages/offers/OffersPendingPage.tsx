import { useState } from 'react'
import { Check, X, Eye, Clock, Search, Filter } from 'lucide-react'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'

// Mock data - À remplacer par GraphQL
const mockPendingOffers = [
  {
    id: '1',
    title: '-30% sur les cocktails',
    partner: { name: 'Le Petit Bar', logo: '' },
    establishment: 'Le Petit Bar - Bastille',
    category: 'Bar',
    discount: { type: 'percentage', value: 30 },
    submittedAt: '2024-12-09T10:00:00Z',
    status: 'pending',
  },
  {
    id: '2',
    title: 'Menu découverte à 25€',
    partner: { name: 'Chez Marcel', logo: '' },
    establishment: 'Chez Marcel - Marais',
    category: 'Restaurant',
    discount: { type: 'fixed', value: 25 },
    submittedAt: '2024-12-08T15:30:00Z',
    status: 'pending',
  },
  {
    id: '3',
    title: 'Happy Hour 2 pour 1',
    partner: { name: 'The Irish Pub', logo: '' },
    establishment: 'The Irish Pub - Opéra',
    category: 'Bar',
    discount: { type: 'formula', value: 0 },
    submittedAt: '2024-12-08T09:00:00Z',
    status: 'pending',
  },
]

export function OffersPendingPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedOffer, setSelectedOffer] = useState<typeof mockPendingOffers[0] | null>(null)
  const [actionType, setActionType] = useState<'approve' | 'reject' | null>(null)
  const [rejectReason, setRejectReason] = useState('')
  const [isProcessing, setIsProcessing] = useState(false)

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('fr-FR', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const formatDiscount = (discount: { type: string; value: number }) => {
    switch (discount.type) {
      case 'percentage':
        return `-${discount.value}%`
      case 'fixed':
        return `${discount.value}€`
      case 'formula':
        return 'Formule'
      default:
        return '-'
    }
  }

  const handleApprove = async () => {
    if (!selectedOffer) return
    setIsProcessing(true)
    
    // TODO: Appel GraphQL mutation
    console.log('Approving offer:', selectedOffer.id)
    
    setTimeout(() => {
      setIsProcessing(false)
      setSelectedOffer(null)
      setActionType(null)
    }, 1000)
  }

  const handleReject = async () => {
    if (!selectedOffer || !rejectReason.trim()) return
    setIsProcessing(true)
    
    // TODO: Appel GraphQL mutation
    console.log('Rejecting offer:', selectedOffer.id, 'Reason:', rejectReason)
    
    setTimeout(() => {
      setIsProcessing(false)
      setSelectedOffer(null)
      setActionType(null)
      setRejectReason('')
    }, 1000)
  }

  const openActionDialog = (offer: typeof mockPendingOffers[0], action: 'approve' | 'reject') => {
    setSelectedOffer(offer)
    setActionType(action)
  }

  const filteredOffers = mockPendingOffers.filter(offer =>
    offer.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    offer.partner.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    offer.establishment.toLowerCase().includes(searchQuery.toLowerCase())
  )

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Offres en attente</h1>
        <p className="text-muted-foreground">
          Validez ou rejetez les nouvelles offres soumises par les partenaires.
        </p>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">En attente</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{mockPendingOffers.length}</div>
            <p className="text-xs text-muted-foreground">
              offres à examiner
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Approuvées aujourd'hui</CardTitle>
            <Check className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">12</div>
            <p className="text-xs text-muted-foreground">
              offres validées
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Rejetées aujourd'hui</CardTitle>
            <X className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">3</div>
            <p className="text-xs text-muted-foreground">
              offres refusées
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Search */}
      <Card>
        <CardHeader>
          <CardTitle>File d'attente</CardTitle>
          <CardDescription>
            Examinez les offres dans l'ordre de soumission
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-4 mb-4">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Rechercher par titre, partenaire..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10"
              />
            </div>
            <Button variant="outline">
              <Filter className="mr-2 h-4 w-4" />
              Filtres
            </Button>
          </div>

          {/* Table */}
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Offre</TableHead>
                  <TableHead>Partenaire</TableHead>
                  <TableHead>Établissement</TableHead>
                  <TableHead>Réduction</TableHead>
                  <TableHead>Soumise le</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredOffers.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="h-24 text-center">
                      <div className="flex flex-col items-center gap-2">
                        <Clock className="h-8 w-8 text-muted-foreground" />
                        <p className="text-muted-foreground">Aucune offre en attente</p>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredOffers.map((offer) => (
                    <TableRow key={offer.id}>
                      <TableCell>
                        <div className="font-medium">{offer.title}</div>
                        <Badge variant="outline" className="mt-1">
                          {offer.category}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <div className="h-8 w-8 rounded-full bg-muted flex items-center justify-center">
                            {offer.partner.name.charAt(0)}
                          </div>
                          <span>{offer.partner.name}</span>
                        </div>
                      </TableCell>
                      <TableCell>{offer.establishment}</TableCell>
                      <TableCell>
                        <Badge variant="secondary">
                          {formatDiscount(offer.discount)}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {formatDate(offer.submittedAt)}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center justify-end gap-2">
                          <Button variant="ghost" size="icon" asChild>
                            <Link to={`/offers/${offer.id}`}>
                              <Eye className="h-4 w-4" />
                            </Link>
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            className="text-green-500 hover:text-green-600 hover:bg-green-50"
                            onClick={() => openActionDialog(offer, 'approve')}
                          >
                            <Check className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            className="text-red-500 hover:text-red-600 hover:bg-red-50"
                            onClick={() => openActionDialog(offer, 'reject')}
                          >
                            <X className="h-4 w-4" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      {/* Action Dialog */}
      <Dialog open={!!actionType} onOpenChange={() => {
        setActionType(null)
        setSelectedOffer(null)
        setRejectReason('')
      }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {actionType === 'approve' ? 'Approuver l\'offre' : 'Rejeter l\'offre'}
            </DialogTitle>
            <DialogDescription>
              {actionType === 'approve'
                ? 'Cette offre sera publiée et visible par les utilisateurs.'
                : 'Veuillez indiquer la raison du rejet.'}
            </DialogDescription>
          </DialogHeader>

          {selectedOffer && (
            <div className="space-y-4">
              <div className="rounded-lg border p-4 bg-muted/50">
                <p className="font-medium">{selectedOffer.title}</p>
                <p className="text-sm text-muted-foreground">
                  {selectedOffer.partner.name} - {selectedOffer.establishment}
                </p>
              </div>

              {actionType === 'reject' && (
                <div className="space-y-2">
                  <label className="text-sm font-medium">Raison du rejet</label>
                  <Textarea
                    placeholder="Expliquez pourquoi cette offre est rejetée..."
                    value={rejectReason}
                    onChange={(e) => setRejectReason(e.target.value)}
                    rows={4}
                  />
                </div>
              )}
            </div>
          )}

          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => {
                setActionType(null)
                setSelectedOffer(null)
                setRejectReason('')
              }}
            >
              Annuler
            </Button>
            {actionType === 'approve' ? (
              <Button
                onClick={handleApprove}
                disabled={isProcessing}
                className="bg-green-600 hover:bg-green-700"
              >
                {isProcessing ? 'Validation...' : 'Approuver'}
              </Button>
            ) : (
              <Button
                variant="destructive"
                onClick={handleReject}
                disabled={isProcessing || !rejectReason.trim()}
              >
                {isProcessing ? 'Traitement...' : 'Rejeter'}
              </Button>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
