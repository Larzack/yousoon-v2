import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Users,
  Building2,
  Tag,
  ShieldCheck,
  TrendingUp,
  TrendingDown,
  Activity,
  CreditCard,
} from 'lucide-react'

interface StatCardProps {
  title: string
  value: string | number
  change?: number
  icon: React.ElementType
  color: string
}

function StatCard({ title, value, change, icon: Icon, color }: StatCardProps) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          {title}
        </CardTitle>
        <div className={`h-8 w-8 rounded-lg ${color} flex items-center justify-center`}>
          <Icon className="h-4 w-4 text-white" />
        </div>
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        {change !== undefined && (
          <p className={`text-xs flex items-center gap-1 mt-1 ${change >= 0 ? 'text-green-600' : 'text-red-600'}`}>
            {change >= 0 ? <TrendingUp className="h-3 w-3" /> : <TrendingDown className="h-3 w-3" />}
            {Math.abs(change)}% vs mois dernier
          </p>
        )}
      </CardContent>
    </Card>
  )
}

interface PendingItem {
  id: string
  type: string
  title: string
  createdAt: string
}

const pendingItems: PendingItem[] = [
  { id: '1', type: 'partner', title: 'Le Petit Bistrot - En attente de validation', createdAt: 'Il y a 2h' },
  { id: '2', type: 'identity', title: 'Marie Dupont - Vérification CNI', createdAt: 'Il y a 3h' },
  { id: '3', type: 'review', title: 'Avis signalé sur Restaurant Le Gourmet', createdAt: 'Il y a 5h' },
  { id: '4', type: 'partner', title: 'Escape Game Paris - En attente de validation', createdAt: 'Il y a 1j' },
  { id: '5', type: 'identity', title: 'Jean Martin - Vérification CNI', createdAt: 'Il y a 1j' },
]

export function DashboardPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground">Vue d'ensemble de la plateforme Yousoon</p>
      </div>

      {/* Stats Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Utilisateurs actifs"
          value="12,543"
          change={12.5}
          icon={Users}
          color="bg-blue-500"
        />
        <StatCard
          title="Partenaires"
          value="234"
          change={8.2}
          icon={Building2}
          color="bg-purple-500"
        />
        <StatCard
          title="Offres actives"
          value="1,892"
          change={-2.4}
          icon={Tag}
          color="bg-orange-500"
        />
        <StatCard
          title="Revenus MRR"
          value="€45,230"
          change={15.3}
          icon={CreditCard}
          color="bg-green-500"
        />
      </div>

      {/* Secondary Stats */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              CNI en attente
            </CardTitle>
            <ShieldCheck className="h-4 w-4 text-yellow-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">12</div>
            <p className="text-xs text-muted-foreground mt-1">À traiter aujourd'hui</p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Partenaires en attente
            </CardTitle>
            <Building2 className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-orange-600">5</div>
            <p className="text-xs text-muted-foreground mt-1">À valider</p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Avis signalés
            </CardTitle>
            <Activity className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">3</div>
            <p className="text-xs text-muted-foreground mt-1">À modérer</p>
          </CardContent>
        </Card>
      </div>

      {/* Pending Actions */}
      <Card>
        <CardHeader>
          <CardTitle>Actions en attente</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {pendingItems.map((item) => (
              <div
                key={item.id}
                className="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
              >
                <div className="flex items-center gap-3">
                  <div
                    className={`h-2 w-2 rounded-full ${
                      item.type === 'partner'
                        ? 'bg-purple-500'
                        : item.type === 'identity'
                        ? 'bg-yellow-500'
                        : 'bg-red-500'
                    }`}
                  />
                  <div>
                    <p className="font-medium">{item.title}</p>
                    <p className="text-sm text-muted-foreground">{item.createdAt}</p>
                  </div>
                </div>
                <a
                  href={`/${item.type === 'partner' ? 'partners' : item.type === 'identity' ? 'identity' : 'reviews'}`}
                  className="text-indigo-600 text-sm font-medium hover:underline"
                >
                  Voir
                </a>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
