import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Textarea } from '@/components/ui/textarea'
import { formatDate } from '@/lib/utils'
import {
  ArrowLeft,
  CheckCircle,
  XCircle,
  Clock,
  User,
  FileText,
  AlertTriangle,
  Shield,
  Eye,
  ZoomIn,
  AlertCircle,
} from 'lucide-react'

// Mock data
const mockVerification = {
  id: '1',
  user: {
    id: 'u1',
    firstName: 'Jean',
    lastName: 'Dupont',
    email: 'jean.dupont@example.com',
    phone: '+33 6 12 34 56 78',
    createdAt: '2024-11-15',
    profile: {
      firstName: 'Jean',
      lastName: 'Dupont',
      birthDate: '1985-03-15',
    },
  },
  documentType: 'cni',
  status: 'pending' as const,
  submittedAt: '2024-12-09T10:30:00',
  documents: {
    frontImage: '/placeholder-cni-front.jpg',
    backImage: '/placeholder-cni-back.jpg',
    selfie: '/placeholder-selfie.jpg',
  },
  extractedData: {
    firstName: 'Jean',
    lastName: 'Dupont',
    birthDate: '1985-03-15',
    birthPlace: 'Paris',
    documentNumber: '123456789012',
    issueDate: '2020-03-14',
    expiryDate: '2030-03-14',
    nationality: 'Française',
  },
  ocrConfidence: 92,
  attempts: 1,
  maxAttempts: 10,
  history: [
    { action: 'Document soumis', date: '2024-12-09T10:30:00', actor: null },
  ],
}

function getStatusBadge(status: string) {
  const styles: Record<string, { bg: string; text: string; label: string; icon: typeof CheckCircle }> = {
    pending: { bg: 'bg-yellow-100', text: 'text-yellow-700', label: 'En attente', icon: Clock },
    verified: { bg: 'bg-green-100', text: 'text-green-700', label: 'Vérifié', icon: CheckCircle },
    rejected: { bg: 'bg-red-100', text: 'text-red-700', label: 'Rejeté', icon: XCircle },
  }
  const style = styles[status] || styles.pending
  const Icon = style.icon
  return (
    <span className={`inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm font-medium ${style.bg} ${style.text}`}>
      <Icon className="h-4 w-4" />
      {style.label}
    </span>
  )
}

function DataComparison({ label, extracted, profile, match }: { label: string; extracted: string; profile: string; match: boolean }) {
  return (
    <div className="grid grid-cols-3 gap-4 py-2 border-b last:border-0">
      <div className="text-sm text-muted-foreground">{label}</div>
      <div className="text-sm font-medium">{extracted}</div>
      <div className={`text-sm flex items-center gap-1 ${match ? 'text-green-600' : 'text-red-600'}`}>
        {profile}
        {match ? (
          <CheckCircle className="h-4 w-4" />
        ) : (
          <AlertCircle className="h-4 w-4" />
        )}
      </div>
    </div>
  )
}

export function IdentityDetailPage() {
  const { id: _id } = useParams()
  const navigate = useNavigate()
  const [rejectReason, setRejectReason] = useState('')
  const [showRejectModal, setShowRejectModal] = useState(false)
  const [selectedImage, setSelectedImage] = useState<string | null>(null)

  const verification = mockVerification // In real app, fetch by id

  const handleValidate = () => {
    // In real app, call API
    alert('Vérification validée !')
    navigate('/identity')
  }

  const handleReject = () => {
    if (!rejectReason.trim()) {
      alert('Veuillez indiquer un motif de rejet')
      return
    }
    // In real app, call API with reason
    alert(`Vérification rejetée: ${rejectReason}`)
    navigate('/identity')
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div className="flex-1">
          <h1 className="text-2xl font-bold">Vérification d'identité</h1>
          <p className="text-muted-foreground">
            {verification.user.firstName} {verification.user.lastName} • {verification.user.email}
          </p>
        </div>
        <div className="flex items-center gap-3">
          {getStatusBadge(verification.status)}
        </div>
      </div>

      {/* Alert for multiple attempts */}
      {verification.attempts >= 3 && (
        <div className="flex items-center gap-3 p-4 bg-orange-50 border border-orange-200 rounded-lg">
          <AlertTriangle className="h-5 w-5 text-orange-600" />
          <div>
            <p className="font-medium text-orange-800">Attention: Tentatives multiples</p>
            <p className="text-sm text-orange-600">
              L'utilisateur a soumis {verification.attempts} tentatives de vérification.
            </p>
          </div>
        </div>
      )}

      <div className="grid gap-6 md:grid-cols-3">
        {/* Document Images */}
        <div className="md:col-span-2 space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <FileText className="h-5 w-5" />
                Documents soumis
              </CardTitle>
              <CardDescription>
                Cliquez sur une image pour l'agrandir
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 md:grid-cols-2">
                {/* Front */}
                <div>
                  <p className="text-sm font-medium mb-2">Recto du document</p>
                  <div
                    className="relative aspect-[1.6] bg-gray-100 rounded-lg border-2 border-dashed border-gray-300 flex items-center justify-center cursor-pointer hover:border-primary transition-colors group"
                    onClick={() => setSelectedImage('front')}
                  >
                    <div className="text-center">
                      <Eye className="h-8 w-8 mx-auto text-gray-400 group-hover:text-primary" />
                      <p className="text-sm text-muted-foreground mt-2">Cliquer pour voir</p>
                    </div>
                    <div className="absolute inset-0 bg-black/0 group-hover:bg-black/10 transition-colors rounded-lg flex items-center justify-center">
                      <ZoomIn className="h-6 w-6 text-white opacity-0 group-hover:opacity-100 transition-opacity" />
                    </div>
                  </div>
                </div>

                {/* Back */}
                <div>
                  <p className="text-sm font-medium mb-2">Verso du document</p>
                  <div
                    className="relative aspect-[1.6] bg-gray-100 rounded-lg border-2 border-dashed border-gray-300 flex items-center justify-center cursor-pointer hover:border-primary transition-colors group"
                    onClick={() => setSelectedImage('back')}
                  >
                    <div className="text-center">
                      <Eye className="h-8 w-8 mx-auto text-gray-400 group-hover:text-primary" />
                      <p className="text-sm text-muted-foreground mt-2">Cliquer pour voir</p>
                    </div>
                    <div className="absolute inset-0 bg-black/0 group-hover:bg-black/10 transition-colors rounded-lg flex items-center justify-center">
                      <ZoomIn className="h-6 w-6 text-white opacity-0 group-hover:opacity-100 transition-opacity" />
                    </div>
                  </div>
                </div>

                {/* Selfie (optional) */}
                {verification.documents.selfie && (
                  <div className="md:col-span-2">
                    <p className="text-sm font-medium mb-2">Selfie de vérification</p>
                    <div
                      className="relative aspect-video max-w-sm bg-gray-100 rounded-lg border-2 border-dashed border-gray-300 flex items-center justify-center cursor-pointer hover:border-primary transition-colors group"
                      onClick={() => setSelectedImage('selfie')}
                    >
                      <div className="text-center">
                        <User className="h-8 w-8 mx-auto text-gray-400 group-hover:text-primary" />
                        <p className="text-sm text-muted-foreground mt-2">Cliquer pour voir</p>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>

          {/* Data Comparison */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Shield className="h-5 w-5" />
                Comparaison des données
              </CardTitle>
              <CardDescription>
                Données extraites par OCR vs profil utilisateur
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="mb-4">
                <div className="flex items-center gap-2 mb-2">
                  <span className="text-sm text-muted-foreground">Confiance OCR:</span>
                  <span className={`text-sm font-bold ${verification.ocrConfidence >= 80 ? 'text-green-600' : 'text-orange-600'}`}>
                    {verification.ocrConfidence}%
                  </span>
                </div>
                <div className="h-2 bg-gray-100 rounded-full overflow-hidden">
                  <div 
                    className={`h-full rounded-full ${verification.ocrConfidence >= 80 ? 'bg-green-500' : 'bg-orange-500'}`}
                    style={{ width: `${verification.ocrConfidence}%` }}
                  />
                </div>
              </div>

              <div className="grid grid-cols-3 gap-4 py-2 border-b bg-gray-50 -mx-6 px-6 mb-2">
                <div className="text-sm font-medium text-muted-foreground">Champ</div>
                <div className="text-sm font-medium text-muted-foreground">Document (OCR)</div>
                <div className="text-sm font-medium text-muted-foreground">Profil</div>
              </div>

              <DataComparison
                label="Prénom"
                extracted={verification.extractedData.firstName}
                profile={verification.user.profile.firstName}
                match={verification.extractedData.firstName.toLowerCase() === verification.user.profile.firstName.toLowerCase()}
              />
              <DataComparison
                label="Nom"
                extracted={verification.extractedData.lastName}
                profile={verification.user.profile.lastName}
                match={verification.extractedData.lastName.toLowerCase() === verification.user.profile.lastName.toLowerCase()}
              />
              <DataComparison
                label="Date de naissance"
                extracted={formatDate(verification.extractedData.birthDate)}
                profile={formatDate(verification.user.profile.birthDate)}
                match={verification.extractedData.birthDate === verification.user.profile.birthDate}
              />

              <div className="grid grid-cols-2 gap-4 mt-4 pt-4 border-t">
                <div>
                  <p className="text-sm text-muted-foreground">N° Document</p>
                  <p className="font-medium">{verification.extractedData.documentNumber}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Date d'expiration</p>
                  <p className="font-medium">{formatDate(verification.extractedData.expiryDate)}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Lieu de naissance</p>
                  <p className="font-medium">{verification.extractedData.birthPlace}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Nationalité</p>
                  <p className="font-medium">{verification.extractedData.nationality}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Rejection Reason (only show if pending) */}
          {verification.status === 'pending' && showRejectModal && (
            <Card className="border-red-200">
              <CardHeader>
                <CardTitle className="text-red-600">Motif de rejet</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <Textarea
                  placeholder="Indiquez le motif du rejet (ex: Document illisible, Photo floue, Données ne correspondent pas...)"
                  value={rejectReason}
                  onChange={(e) => setRejectReason(e.target.value)}
                  rows={3}
                />
                <div className="flex gap-2">
                  <Button variant="destructive" onClick={handleReject}>
                    Confirmer le rejet
                  </Button>
                  <Button variant="outline" onClick={() => setShowRejectModal(false)}>
                    Annuler
                  </Button>
                </div>
              </CardContent>
            </Card>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Actions */}
          {verification.status === 'pending' && (
            <Card>
              <CardHeader>
                <CardTitle>Actions</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <Button className="w-full gap-2 bg-green-600 hover:bg-green-700" onClick={handleValidate}>
                  <CheckCircle className="h-4 w-4" />
                  Valider l'identité
                </Button>
                <Button
                  variant="destructive"
                  className="w-full gap-2"
                  onClick={() => setShowRejectModal(true)}
                >
                  <XCircle className="h-4 w-4" />
                  Rejeter
                </Button>
              </CardContent>
            </Card>
          )}

          {/* User Info */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <User className="h-5 w-5" />
                Utilisateur
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div>
                <p className="text-sm text-muted-foreground">Nom complet</p>
                <Link to={`/users/${verification.user.id}`} className="font-medium text-primary hover:underline">
                  {verification.user.firstName} {verification.user.lastName}
                </Link>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Email</p>
                <p className="font-medium">{verification.user.email}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Téléphone</p>
                <p className="font-medium">{verification.user.phone}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Inscrit le</p>
                <p className="font-medium">{formatDate(verification.user.createdAt)}</p>
              </div>
            </CardContent>
          </Card>

          {/* Verification Info */}
          <Card>
            <CardHeader>
              <CardTitle>Informations</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Type de document</span>
                <span className="font-medium">
                  {verification.documentType === 'cni' ? 'CNI' :
                   verification.documentType === 'passport' ? 'Passeport' : 'Permis'}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Soumis le</span>
                <span className="font-medium">{formatDate(verification.submittedAt)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Tentatives</span>
                <span className={`font-medium ${verification.attempts >= 5 ? 'text-orange-600' : ''}`}>
                  {verification.attempts} / {verification.maxAttempts}
                </span>
              </div>
            </CardContent>
          </Card>

          {/* History */}
          <Card>
            <CardHeader>
              <CardTitle>Historique</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {verification.history.map((log, i) => (
                  <div key={i} className="flex items-start gap-3 text-sm">
                    <div className="h-2 w-2 rounded-full bg-gray-300 mt-2" />
                    <div>
                      <p className="font-medium">{log.action}</p>
                      <p className="text-muted-foreground">{formatDate(log.date)}</p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Image Modal */}
      {selectedImage && (
        <div
          className="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-4"
          onClick={() => setSelectedImage(null)}
        >
          <div className="bg-white rounded-lg p-4 max-w-4xl max-h-[90vh] overflow-auto">
            <div className="aspect-[1.6] bg-gray-200 rounded-lg flex items-center justify-center min-h-[400px]">
              <div className="text-center">
                <FileText className="h-16 w-16 mx-auto text-gray-400" />
                <p className="text-muted-foreground mt-4">
                  Image du document ({selectedImage === 'front' ? 'Recto' : selectedImage === 'back' ? 'Verso' : 'Selfie'})
                </p>
                <p className="text-sm text-muted-foreground mt-2">
                  En production, l'image réelle serait affichée ici
                </p>
              </div>
            </div>
            <Button className="mt-4" onClick={() => setSelectedImage(null)}>
              Fermer
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}
