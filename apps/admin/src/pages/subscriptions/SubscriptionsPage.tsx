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
  CreditCard,
  User,
  ChevronLeft,
  ChevronRight,
  XCircle,
  RefreshCw,
  Settings,
  Calendar,
  TrendingUp,
  DollarSign,
} from 'lucide-react'

interface Subscription {
  id: string
  user: {
    id: string
    firstName: string
    lastName: string
    email: string
  }
  plan: {
    code: string
    name: string
    price: number
  }
  platform: 'apple' | 'google'
  status: 'trialing' | 'active' | 'past_due' | 'cancelled' | 'expired'
  trialEnd: string | null
  currentPeriodStart: string
  currentPeriodEnd: string
  cancelledAt: string | null
  createdAt: string
}

// Mock data
const mockSubscriptions: Subscription[] = [
  {
    id: '1',
    user: { id: 'u1', firstName: 'Jean', lastName: 'Dupont', email: 'jean.dupont@example.com' },
    plan: { code: 'monthly', name: 'Mensuel', price: 990 },
    platform: 'apple',
    status: 'active',
    trialEnd: null,
    currentPeriodStart: '2024-12-01',
    currentPeriodEnd: '2025-01-01',
    cancelledAt: null,
    createdAt: '2024-06-15',
  },
  {
    id: '2',
    user: { id: 'u2', firstName: 'Marie', lastName: 'Martin', email: 'marie.martin@example.com' },
    plan: { code: 'yearly', name: 'Annuel', price: 7990 },
    platform: 'google',
    status: 'trialing',
    trialEnd: '2025-01-08',
    currentPeriodStart: '2024-12-09',
    currentPeriodEnd: '2025-12-09',
    cancelledAt: null,
    createdAt: '2024-12-09',
  },
  {
    id: '3',
    user: { id: 'u3', firstName: 'Pierre', lastName: 'Dubois', email: 'pierre.dubois@example.com' },
    plan: { code: 'monthly', name: 'Mensuel', price: 990 },
    platform: 'apple',
    status: 'cancelled',
    trialEnd: null,
    currentPeriodStart: '2024-11-15',
    currentPeriodEnd: '2024-12-15',
    cancelledAt: '2024-12-05',
    createdAt: '2024-09-15',
  },
  {
    id: '4',
    user: { id: 'u4', firstName: 'Sophie', lastName: 'Leroy', email: 'sophie.leroy@example.com' },
    plan: { code: 'yearly', name: 'Annuel', price: 7990 },
    platform: 'google',
    status: 'past_due',
    trialEnd: null,
    currentPeriodStart: '2024-11-20',
    currentPeriodEnd: '2024-12-20',
    cancelledAt: null,
    createdAt: '2024-11-20',
  },
  {
    id: '5',
    user: { id: 'u5', firstName: 'Lucas', lastName: 'Bernard', email: 'lucas.bernard@example.com' },
    plan: { code: 'monthly', name: 'Mensuel', price: 990 },
    platform: 'apple',
    status: 'active',
    trialEnd: null,
    currentPeriodStart: '2024-12-05',
    currentPeriodEnd: '2025-01-05',
    cancelledAt: null,
    createdAt: '2024-08-05',
  },
]

function getStatusBadge(status: Subscription['status']) {
  const styles = {
    trialing: { bg: 'bg-blue-100', text: 'text-blue-700', label: 'Essai' },
    active: { bg: 'bg-green-100', text: 'text-green-700', label: 'Actif' },
    past_due: { bg: 'bg-orange-100', text: 'text-orange-700', label: 'Impayé' },
    cancelled: { bg: 'bg-red-100', text: 'text-red-700', label: 'Annulé' },
    expired: { bg: 'bg-gray-100', text: 'text-gray-700', label: 'Expiré' },
  }
  const style = styles[status]
  return (
    <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${style.bg} ${style.text}`}>
      {style.label}
    </span>
  )
}

function formatPrice(cents: number) {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: 'EUR',
  }).format(cents / 100)
}

function getPlatformIcon(platform: 'apple' | 'google') {
  if (platform === 'apple') {
    return <span className="text-xs bg-gray-100 px-2 py-0.5 rounded">Apple</span>
  }
  return <span className="text-xs bg-gray-100 px-2 py-0.5 rounded">Google</span>
}

export function SubscriptionsPage() {
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [planFilter, setPlanFilter] = useState<string>('all')
  const [currentPage, setCurrentPage] = useState(1)

  const filteredSubscriptions = mockSubscriptions.filter((sub) => {
    const matchesSearch =
      sub.user.firstName.toLowerCase().includes(search.toLowerCase()) ||
      sub.user.lastName.toLowerCase().includes(search.toLowerCase()) ||
      sub.user.email.toLowerCase().includes(search.toLowerCase())
    const matchesStatus = statusFilter === 'all' || sub.status === statusFilter
    const matchesPlan = planFilter === 'all' || sub.plan.code === planFilter
    return matchesSearch && matchesStatus && matchesPlan
  })

  // Mock stats
  const mrr = 12450 // en centimes
  const activeCount = mockSubscriptions.filter((s) => s.status === 'active').length
  const trialingCount = mockSubscriptions.filter((s) => s.status === 'trialing').length

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Abonnements</h1>
          <p className="text-muted-foreground">Gestion des abonnements utilisateurs</p>
        </div>
        <Link to="/subscriptions/plans">
          <Button variant="outline" className="gap-2">
            <Settings className="h-4 w-4" />
            Gérer les plans
          </Button>
        </Link>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <DollarSign className="h-4 w-4" />
              MRR
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatPrice(mrr * 100)}</div>
            <p className="text-xs text-green-600 flex items-center gap-1 mt-1">
              <TrendingUp className="h-3 w-3" />
              +12% vs mois dernier
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Abonnés actifs</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{activeCount}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">En essai</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">{trialingCount}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Taux de conversion</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">68%</div>
            <p className="text-xs text-muted-foreground">Essai → Payant</p>
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
            <div className="flex gap-2 flex-wrap">
              {['all', 'active', 'trialing', 'past_due', 'cancelled'].map((status) => (
                <Button
                  key={status}
                  variant={statusFilter === status ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setStatusFilter(status)}
                >
                  {status === 'all' ? 'Tous' :
                   status === 'active' ? 'Actifs' :
                   status === 'trialing' ? 'Essai' :
                   status === 'past_due' ? 'Impayés' : 'Annulés'}
                </Button>
              ))}
            </div>
          </div>
          <div className="flex gap-2 mt-4">
            <span className="text-sm text-muted-foreground py-1">Plan :</span>
            {['all', 'monthly', 'yearly'].map((plan) => (
              <Button
                key={plan}
                variant={planFilter === plan ? 'secondary' : 'ghost'}
                size="sm"
                onClick={() => setPlanFilter(plan)}
              >
                {plan === 'all' ? 'Tous' : plan === 'monthly' ? 'Mensuel' : 'Annuel'}
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
                  <th className="text-left p-4 font-medium text-muted-foreground">Utilisateur</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Plan</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Plateforme</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Statut</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Période</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Depuis</th>
                  <th className="text-right p-4 font-medium text-muted-foreground">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {filteredSubscriptions.map((subscription) => (
                  <tr key={subscription.id} className="hover:bg-gray-50">
                    <td className="p-4">
                      <div className="flex items-center gap-3">
                        <div className="h-10 w-10 rounded-full bg-gray-100 flex items-center justify-center">
                          <User className="h-5 w-5 text-gray-500" />
                        </div>
                        <div>
                          <Link
                            to={`/users/${subscription.user.id}`}
                            className="font-medium hover:text-primary hover:underline"
                          >
                            {subscription.user.firstName} {subscription.user.lastName}
                          </Link>
                          <p className="text-sm text-muted-foreground">{subscription.user.email}</p>
                        </div>
                      </div>
                    </td>
                    <td className="p-4">
                      <div>
                        <p className="font-medium">{subscription.plan.name}</p>
                        <p className="text-sm text-muted-foreground">{formatPrice(subscription.plan.price)}/mois</p>
                      </div>
                    </td>
                    <td className="p-4">
                      {getPlatformIcon(subscription.platform)}
                    </td>
                    <td className="p-4">
                      <div>
                        {getStatusBadge(subscription.status)}
                        {subscription.trialEnd && (
                          <p className="text-xs text-muted-foreground mt-1">
                            Fin essai: {formatDate(subscription.trialEnd)}
                          </p>
                        )}
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="text-sm">
                        <p className="flex items-center gap-1">
                          <Calendar className="h-3 w-3" />
                          {formatDate(subscription.currentPeriodStart)}
                        </p>
                        <p className="text-muted-foreground">→ {formatDate(subscription.currentPeriodEnd)}</p>
                      </div>
                    </td>
                    <td className="p-4 text-sm text-muted-foreground">
                      {formatDate(subscription.createdAt)}
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
                              <Link to={`/users/${subscription.user.id}`} className="flex items-center gap-2">
                                <Eye className="h-4 w-4" />
                                Voir l'utilisateur
                              </Link>
                            </DropdownMenuItem>
                            <DropdownMenuItem className="gap-2">
                              <CreditCard className="h-4 w-4" />
                              Voir les paiements
                            </DropdownMenuItem>
                            {subscription.status === 'active' && (
                              <DropdownMenuItem className="gap-2 text-red-600">
                                <XCircle className="h-4 w-4" />
                                Annuler
                              </DropdownMenuItem>
                            )}
                            {subscription.status === 'cancelled' && (
                              <DropdownMenuItem className="gap-2 text-green-600">
                                <RefreshCw className="h-4 w-4" />
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
              Affichage de {filteredSubscriptions.length} abonnements
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
