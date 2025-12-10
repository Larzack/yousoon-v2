import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { formatDate } from '@/lib/utils'
import {
  Search,
  Eye,
  CheckCircle,
  XCircle,
  Clock,
  AlertTriangle,
  ChevronLeft,
  ChevronRight,
  FileText,
  User,
  Calendar,
} from 'lucide-react'

interface IdentityVerification {
  id: string
  user: {
    id: string
    firstName: string
    lastName: string
    email: string
  }
  documentType: 'cni' | 'passport' | 'permit'
  status: 'pending' | 'verified' | 'rejected'
  submittedAt: string
  extractedData: {
    firstName: string
    lastName: string
    birthDate: string
    documentNumber: string
    expiryDate: string
  } | null
  attempts: number
  maxAttempts: number
}

// Mock data
const mockVerifications: IdentityVerification[] = [
  {
    id: '1',
    user: {
      id: 'u1',
      firstName: 'Jean',
      lastName: 'Dupont',
      email: 'jean.dupont@example.com',
    },
    documentType: 'cni',
    status: 'pending',
    submittedAt: '2024-12-09T10:30:00',
    extractedData: {
      firstName: 'Jean',
      lastName: 'Dupont',
      birthDate: '1985-03-15',
      documentNumber: '123456789012',
      expiryDate: '2030-03-14',
    },
    attempts: 1,
    maxAttempts: 10,
  },
  {
    id: '2',
    user: {
      id: 'u2',
      firstName: 'Marie',
      lastName: 'Martin',
      email: 'marie.martin@example.com',
    },
    documentType: 'passport',
    status: 'pending',
    submittedAt: '2024-12-09T09:15:00',
    extractedData: {
      firstName: 'Marie',
      lastName: 'Martin',
      birthDate: '1990-07-22',
      documentNumber: '15AB23456',
      expiryDate: '2028-07-21',
    },
    attempts: 2,
    maxAttempts: 10,
  },
  {
    id: '3',
    user: {
      id: 'u3',
      firstName: 'Pierre',
      lastName: 'Dubois',
      email: 'pierre.dubois@example.com',
    },
    documentType: 'cni',
    status: 'verified',
    submittedAt: '2024-12-08T16:45:00',
    extractedData: {
      firstName: 'Pierre',
      lastName: 'Dubois',
      birthDate: '1978-11-08',
      documentNumber: '987654321098',
      expiryDate: '2029-11-07',
    },
    attempts: 1,
    maxAttempts: 10,
  },
  {
    id: '4',
    user: {
      id: 'u4',
      firstName: 'Sophie',
      lastName: 'Leroy',
      email: 'sophie.leroy@example.com',
    },
    documentType: 'cni',
    status: 'rejected',
    submittedAt: '2024-12-08T14:20:00',
    extractedData: null,
    attempts: 3,
    maxAttempts: 10,
  },
  {
    id: '5',
    user: {
      id: 'u5',
      firstName: 'Lucas',
      lastName: 'Bernard',
      email: 'lucas.bernard@example.com',
    },
    documentType: 'permit',
    status: 'pending',
    submittedAt: '2024-12-09T08:00:00',
    extractedData: {
      firstName: 'Lucas',
      lastName: 'Bernard',
      birthDate: '1995-01-30',
      documentNumber: '23FR45678',
      expiryDate: '2027-01-29',
    },
    attempts: 1,
    maxAttempts: 10,
  },
]

function getStatusBadge(status: IdentityVerification['status']) {
  const styles = {
    pending: { bg: 'bg-yellow-100', text: 'text-yellow-700', label: 'En attente', icon: Clock },
    verified: { bg: 'bg-green-100', text: 'text-green-700', label: 'Vérifié', icon: CheckCircle },
    rejected: { bg: 'bg-red-100', text: 'text-red-700', label: 'Rejeté', icon: XCircle },
  }
  const style = styles[status]
  const Icon = style.icon
  return (
    <span className={`inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium ${style.bg} ${style.text}`}>
      <Icon className="h-3 w-3" />
      {style.label}
    </span>
  )
}

function getDocumentTypeLabel(type: IdentityVerification['documentType']) {
  const labels = {
    cni: 'CNI',
    passport: 'Passeport',
    permit: 'Permis de conduire',
  }
  return labels[type]
}

export function IdentityVerificationsPage() {
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('pending')
  const [currentPage, setCurrentPage] = useState(1)

  const filteredVerifications = mockVerifications.filter((v) => {
    const matchesSearch =
      v.user.firstName.toLowerCase().includes(search.toLowerCase()) ||
      v.user.lastName.toLowerCase().includes(search.toLowerCase()) ||
      v.user.email.toLowerCase().includes(search.toLowerCase())
    const matchesStatus = statusFilter === 'all' || v.status === statusFilter
    return matchesSearch && matchesStatus
  })

  const pendingCount = mockVerifications.filter((v) => v.status === 'pending').length

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Vérifications d'identité</h1>
          <p className="text-muted-foreground">Validation des documents CNI des utilisateurs</p>
        </div>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card className="border-yellow-200 bg-yellow-50">
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-yellow-700">En attente</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold text-yellow-700">{pendingCount}</div>
            <p className="text-xs text-yellow-600 mt-1">À traiter</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Vérifiés aujourd'hui</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">12</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Rejetés aujourd'hui</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">3</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Taux d'approbation</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">87%</div>
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
                placeholder="Rechercher par nom ou email..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-9"
              />
            </div>
            <div className="flex gap-2">
              {['pending', 'all', 'verified', 'rejected'].map((status) => (
                <Button
                  key={status}
                  variant={statusFilter === status ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setStatusFilter(status)}
                >
                  {status === 'pending' ? 'En attente' :
                   status === 'all' ? 'Toutes' :
                   status === 'verified' ? 'Vérifiées' : 'Rejetées'}
                </Button>
              ))}
            </div>
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
                  <th className="text-left p-4 font-medium text-muted-foreground">Utilisateur</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Document</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Données extraites</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Soumis le</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Tentatives</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Statut</th>
                  <th className="text-right p-4 font-medium text-muted-foreground">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {filteredVerifications.map((verification) => (
                  <tr key={verification.id} className="hover:bg-gray-50">
                    <td className="p-4">
                      <div className="flex items-center gap-3">
                        <div className="h-10 w-10 rounded-full bg-gray-100 flex items-center justify-center">
                          <User className="h-5 w-5 text-gray-500" />
                        </div>
                        <div>
                          <Link
                            to={`/users/${verification.user.id}`}
                            className="font-medium hover:text-primary hover:underline"
                          >
                            {verification.user.firstName} {verification.user.lastName}
                          </Link>
                          <p className="text-sm text-muted-foreground">{verification.user.email}</p>
                        </div>
                      </div>
                    </td>
                    <td className="p-4">
                      <span className="flex items-center gap-1 text-sm">
                        <FileText className="h-4 w-4 text-muted-foreground" />
                        {getDocumentTypeLabel(verification.documentType)}
                      </span>
                    </td>
                    <td className="p-4">
                      {verification.extractedData ? (
                        <div className="text-sm">
                          <p>{verification.extractedData.firstName} {verification.extractedData.lastName}</p>
                          <p className="text-muted-foreground">
                            Né le {formatDate(verification.extractedData.birthDate)}
                          </p>
                        </div>
                      ) : (
                        <span className="text-sm text-muted-foreground">Non disponible</span>
                      )}
                    </td>
                    <td className="p-4">
                      <span className="flex items-center gap-1 text-sm">
                        <Calendar className="h-4 w-4 text-muted-foreground" />
                        {formatDate(verification.submittedAt)}
                      </span>
                    </td>
                    <td className="p-4">
                      <span className={`text-sm ${verification.attempts >= 5 ? 'text-orange-600 font-medium' : ''}`}>
                        {verification.attempts} / {verification.maxAttempts}
                        {verification.attempts >= 5 && (
                          <AlertTriangle className="h-4 w-4 inline ml-1 text-orange-500" />
                        )}
                      </span>
                    </td>
                    <td className="p-4">{getStatusBadge(verification.status)}</td>
                    <td className="p-4">
                      <div className="flex justify-end gap-2">
                        <Button variant="outline" size="sm" asChild>
                          <Link to={`/identity/${verification.id}`} className="gap-1">
                            <Eye className="h-4 w-4" />
                            Voir
                          </Link>
                        </Button>
                        {verification.status === 'pending' && (
                          <>
                            <Button size="sm" className="gap-1 bg-green-600 hover:bg-green-700">
                              <CheckCircle className="h-4 w-4" />
                              Valider
                            </Button>
                            <Button variant="destructive" size="sm" className="gap-1">
                              <XCircle className="h-4 w-4" />
                              Rejeter
                            </Button>
                          </>
                        )}
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
              Affichage de {filteredVerifications.length} vérifications
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
