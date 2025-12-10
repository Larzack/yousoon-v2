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
} from '@/components/ui/dropdown-menu'
import { cn, formatDate, getStatusColor, getStatusLabel, getInitials } from '@/lib/utils'
import {
  Search,
  Filter,
  MoreHorizontal,
  Eye,
  Ban,
  Trash2,
  Mail,
  CheckCircle,
  XCircle,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react'

interface User {
  id: string
  email: string
  firstName: string
  lastName: string
  avatar?: string
  status: 'active' | 'suspended' | 'deleted'
  identityStatus: 'verified' | 'pending' | 'rejected' | 'not_submitted'
  subscriptionPlan?: string
  createdAt: string
}

// Mock data
const mockUsers: User[] = [
  {
    id: '1',
    email: 'marie.dupont@email.com',
    firstName: 'Marie',
    lastName: 'Dupont',
    status: 'active',
    identityStatus: 'verified',
    subscriptionPlan: 'Premium',
    createdAt: '2024-11-15',
  },
  {
    id: '2',
    email: 'jean.martin@email.com',
    firstName: 'Jean',
    lastName: 'Martin',
    status: 'active',
    identityStatus: 'pending',
    subscriptionPlan: 'Gratuit',
    createdAt: '2024-11-20',
  },
  {
    id: '3',
    email: 'sophie.bernard@email.com',
    firstName: 'Sophie',
    lastName: 'Bernard',
    status: 'suspended',
    identityStatus: 'verified',
    subscriptionPlan: 'Premium',
    createdAt: '2024-10-05',
  },
  {
    id: '4',
    email: 'pierre.leroy@email.com',
    firstName: 'Pierre',
    lastName: 'Leroy',
    status: 'active',
    identityStatus: 'rejected',
    createdAt: '2024-11-25',
  },
  {
    id: '5',
    email: 'claire.moreau@email.com',
    firstName: 'Claire',
    lastName: 'Moreau',
    status: 'active',
    identityStatus: 'not_submitted',
    createdAt: '2024-12-01',
  },
]

export function UsersPage() {
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [currentPage, setCurrentPage] = useState(1)

  const filteredUsers = mockUsers.filter((user) => {
    const matchesSearch =
      user.email.toLowerCase().includes(search.toLowerCase()) ||
      user.firstName.toLowerCase().includes(search.toLowerCase()) ||
      user.lastName.toLowerCase().includes(search.toLowerCase())
    const matchesStatus = statusFilter === 'all' || user.status === statusFilter
    return matchesSearch && matchesStatus
  })

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Utilisateurs</h1>
          <p className="text-muted-foreground">Gestion des utilisateurs de la plateforme</p>
        </div>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Total</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">12,543</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Actifs</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">11,892</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Vérifiés</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">8,234</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Abonnés</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-purple-600">3,456</div>
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
              {['all', 'active', 'suspended'].map((status) => (
                <Button
                  key={status}
                  variant={statusFilter === status ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setStatusFilter(status)}
                >
                  {status === 'all' ? 'Tous' : getStatusLabel(status)}
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
                  <th className="text-left p-4 font-medium text-muted-foreground">Statut</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Identité</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Abonnement</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Inscrit le</th>
                  <th className="text-right p-4 font-medium text-muted-foreground">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {filteredUsers.map((user) => (
                  <tr key={user.id} className="hover:bg-gray-50">
                    <td className="p-4">
                      <div className="flex items-center gap-3">
                        <Avatar>
                          <AvatarImage src={user.avatar} />
                          <AvatarFallback>
                            {getInitials(user.firstName, user.lastName)}
                          </AvatarFallback>
                        </Avatar>
                        <div>
                          <p className="font-medium">{user.firstName} {user.lastName}</p>
                          <p className="text-sm text-muted-foreground">{user.email}</p>
                        </div>
                      </div>
                    </td>
                    <td className="p-4">
                      <span className={cn('px-2 py-1 rounded-full text-xs font-medium', getStatusColor(user.status))}>
                        {getStatusLabel(user.status)}
                      </span>
                    </td>
                    <td className="p-4">
                      <div className="flex items-center gap-1">
                        {user.identityStatus === 'verified' ? (
                          <CheckCircle className="h-4 w-4 text-green-500" />
                        ) : user.identityStatus === 'rejected' ? (
                          <XCircle className="h-4 w-4 text-red-500" />
                        ) : null}
                        <span className={cn('px-2 py-1 rounded-full text-xs font-medium', getStatusColor(user.identityStatus))}>
                          {getStatusLabel(user.identityStatus)}
                        </span>
                      </div>
                    </td>
                    <td className="p-4">
                      {user.subscriptionPlan ? (
                        <span className="text-sm font-medium">{user.subscriptionPlan}</span>
                      ) : (
                        <span className="text-sm text-muted-foreground">-</span>
                      )}
                    </td>
                    <td className="p-4 text-sm text-muted-foreground">
                      {formatDate(user.createdAt)}
                    </td>
                    <td className="p-4 text-right">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem asChild>
                            <Link to={`/users/${user.id}`}>
                              <Eye className="mr-2 h-4 w-4" />
                              Voir détails
                            </Link>
                          </DropdownMenuItem>
                          <DropdownMenuItem>
                            <Mail className="mr-2 h-4 w-4" />
                            Envoyer email
                          </DropdownMenuItem>
                          <DropdownMenuItem className="text-yellow-600">
                            <Ban className="mr-2 h-4 w-4" />
                            Suspendre
                          </DropdownMenuItem>
                          <DropdownMenuItem className="text-red-600">
                            <Trash2 className="mr-2 h-4 w-4" />
                            Supprimer
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          <div className="flex items-center justify-between p-4 border-t">
            <p className="text-sm text-muted-foreground">
              Affichage de 1 à {filteredUsers.length} sur {mockUsers.length} résultats
            </p>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" disabled={currentPage === 1}>
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <span className="text-sm">Page {currentPage}</span>
              <Button variant="outline" size="sm">
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
