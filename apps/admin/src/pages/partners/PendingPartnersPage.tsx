import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { getInitials, formatDate } from '@/lib/utils'
import {
  ArrowLeft,
  CheckCircle,
  XCircle,
  Eye,
  Building2,
  Mail,
  Phone,
  FileText,
  AlertCircle,
} from 'lucide-react'

interface PendingPartner {
  id: string
  companyName: string
  tradeName: string
  logo?: string
  siret: string
  category: string
  createdAt: string
  contact: {
    firstName: string
    lastName: string
    email: string
    phone: string
  }
  documents: {
    kbis: boolean
    cni: boolean
    rib: boolean
  }
}

// Mock data
const mockPendingPartners: PendingPartner[] = [
  {
    id: '1',
    companyName: 'Escape Game Paris SAS',
    tradeName: 'Escape Game Paris',
    siret: '987 654 321 00034',
    category: 'Loisirs',
    createdAt: '2024-12-08',
    contact: {
      firstName: 'Sophie',
      lastName: 'Durand',
      email: 'contact@escapegameparis.fr',
      phone: '+33 1 98 76 54 32',
    },
    documents: { kbis: true, cni: true, rib: true },
  },
  {
    id: '2',
    companyName: 'Restaurant Le Gourmet',
    tradeName: 'Le Gourmet',
    siret: '321 654 987 00090',
    category: 'Restaurant',
    createdAt: '2024-12-09',
    contact: {
      firstName: 'Antoine',
      lastName: 'Bernard',
      email: 'antoine@legourmet.fr',
      phone: '+33 1 23 45 67 89',
    },
    documents: { kbis: true, cni: false, rib: true },
  },
  {
    id: '3',
    companyName: 'Yoga Studio Zen',
    tradeName: 'Yoga Zen',
    siret: '111 222 333 00044',
    category: 'Bien-être',
    createdAt: '2024-12-09',
    contact: {
      firstName: 'Claire',
      lastName: 'Martin',
      email: 'claire@yogazen.fr',
      phone: '+33 6 12 34 56 78',
    },
    documents: { kbis: true, cni: true, rib: true },
  },
]

export function PendingPartnersPage() {
  const navigate = useNavigate()
  const [partners, setPartners] = useState(mockPendingPartners)
  const [selectedPartner, setSelectedPartner] = useState<PendingPartner | null>(null)
  const [rejectReason, setRejectReason] = useState('')

  const handleValidate = (partnerId: string) => {
    // TODO: API call
    setPartners(partners.filter((p) => p.id !== partnerId))
    setSelectedPartner(null)
  }

  const handleReject = (partnerId: string) => {
    // TODO: API call with reason
    setPartners(partners.filter((p) => p.id !== partnerId))
    setSelectedPartner(null)
    setRejectReason('')
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => navigate('/partners')}>
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div>
          <h1 className="text-3xl font-bold">Partenaires en attente</h1>
          <p className="text-muted-foreground">{partners.length} partenaire(s) à valider</p>
        </div>
      </div>

      {partners.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center">
            <CheckCircle className="h-12 w-12 text-green-500 mx-auto mb-4" />
            <h3 className="text-lg font-medium">Aucun partenaire en attente</h3>
            <p className="text-muted-foreground mt-2">Tous les partenaires ont été traités</p>
            <Link to="/partners">
              <Button className="mt-4">Retour à la liste</Button>
            </Link>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-6 lg:grid-cols-2">
          {/* Partners List */}
          <div className="space-y-4">
            {partners.map((partner) => (
              <Card 
                key={partner.id} 
                className={`cursor-pointer transition-all ${selectedPartner?.id === partner.id ? 'ring-2 ring-primary' : 'hover:shadow-md'}`}
                onClick={() => setSelectedPartner(partner)}
              >
                <CardContent className="pt-6">
                  <div className="flex items-start gap-4">
                    <Avatar className="h-12 w-12">
                      <AvatarImage src={partner.logo} />
                      <AvatarFallback className="bg-purple-100 text-purple-600">
                        {getInitials(partner.tradeName)}
                      </AvatarFallback>
                    </Avatar>
                    <div className="flex-1">
                      <h3 className="font-medium">{partner.tradeName}</h3>
                      <p className="text-sm text-muted-foreground">{partner.companyName}</p>
                      <div className="flex items-center gap-4 mt-2 text-sm text-muted-foreground">
                        <span className="flex items-center gap-1">
                          <Building2 className="h-3 w-3" />
                          {partner.category}
                        </span>
                        <span>Inscrit le {formatDate(partner.createdAt)}</span>
                      </div>
                    </div>
                    <div className="flex flex-col gap-1">
                      {Object.entries(partner.documents).map(([doc, valid]) => (
                        <span
                          key={doc}
                          className={`text-xs px-2 py-0.5 rounded ${
                            valid ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
                          }`}
                        >
                          {doc.toUpperCase()} {valid ? '✓' : '✗'}
                        </span>
                      ))}
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          {/* Detail Panel */}
          {selectedPartner ? (
            <Card className="h-fit sticky top-4">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <CardTitle>Validation</CardTitle>
                  <Link to={`/partners/${selectedPartner.id}`}>
                    <Button variant="ghost" size="sm" className="gap-2">
                      <Eye className="h-4 w-4" />
                      Voir détail
                    </Button>
                  </Link>
                </div>
                <CardDescription>{selectedPartner.tradeName}</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                {/* Company Info */}
                <div>
                  <h4 className="font-medium mb-3">Informations entreprise</h4>
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Raison sociale</span>
                      <span>{selectedPartner.companyName}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">SIRET</span>
                      <span className="font-mono">{selectedPartner.siret}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Catégorie</span>
                      <span>{selectedPartner.category}</span>
                    </div>
                  </div>
                </div>

                {/* Contact */}
                <div>
                  <h4 className="font-medium mb-3">Contact</h4>
                  <div className="space-y-2 text-sm">
                    <div className="flex items-center gap-2">
                      <Mail className="h-4 w-4 text-muted-foreground" />
                      <a href={`mailto:${selectedPartner.contact.email}`} className="text-primary hover:underline">
                        {selectedPartner.contact.email}
                      </a>
                    </div>
                    <div className="flex items-center gap-2">
                      <Phone className="h-4 w-4 text-muted-foreground" />
                      {selectedPartner.contact.phone}
                    </div>
                  </div>
                </div>

                {/* Documents */}
                <div>
                  <h4 className="font-medium mb-3">Documents</h4>
                  <div className="space-y-2">
                    {Object.entries(selectedPartner.documents).map(([doc, valid]) => (
                      <div key={doc} className="flex items-center justify-between">
                        <span className="flex items-center gap-2">
                          <FileText className="h-4 w-4 text-muted-foreground" />
                          {doc === 'kbis' ? 'Extrait Kbis' : doc === 'cni' ? 'Pièce d\'identité' : 'RIB'}
                        </span>
                        {valid ? (
                          <span className="flex items-center gap-1 text-green-600 text-sm">
                            <CheckCircle className="h-4 w-4" />
                            Fourni
                          </span>
                        ) : (
                          <span className="flex items-center gap-1 text-red-600 text-sm">
                            <AlertCircle className="h-4 w-4" />
                            Manquant
                          </span>
                        )}
                      </div>
                    ))}
                  </div>
                </div>

                {/* Actions */}
                <div className="flex gap-3 pt-4 border-t">
                  <Button 
                    className="flex-1 gap-2" 
                    onClick={() => handleValidate(selectedPartner.id)}
                    disabled={!Object.values(selectedPartner.documents).every(Boolean)}
                  >
                    <CheckCircle className="h-4 w-4" />
                    Valider
                  </Button>
                  <Button 
                    variant="destructive" 
                    className="flex-1 gap-2"
                    onClick={() => handleReject(selectedPartner.id)}
                  >
                    <XCircle className="h-4 w-4" />
                    Rejeter
                  </Button>
                </div>

                {!Object.values(selectedPartner.documents).every(Boolean) && (
                  <p className="text-sm text-yellow-600 flex items-center gap-2">
                    <AlertCircle className="h-4 w-4" />
                    Documents incomplets - validation impossible
                  </p>
                )}
              </CardContent>
            </Card>
          ) : (
            <Card className="h-fit">
              <CardContent className="py-12 text-center text-muted-foreground">
                <Building2 className="h-12 w-12 mx-auto mb-4 opacity-50" />
                <p>Sélectionnez un partenaire pour voir les détails</p>
              </CardContent>
            </Card>
          )}
        </div>
      )}
    </div>
  )
}
