import { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
  Users,
  Building2,
  Tag,
  Calendar,
  TrendingUp,
  TrendingDown,
  DollarSign,
  Eye,
  Download,
  ArrowUpRight,
  Star,
} from 'lucide-react'

interface StatCard {
  title: string
  value: string
  change: number
  changeLabel: string
  icon: typeof Users
}

const stats: StatCard[] = [
  { title: 'Utilisateurs', value: '45,892', change: 12.5, changeLabel: 'vs mois dernier', icon: Users },
  { title: 'Partenaires', value: '234', change: 8.3, changeLabel: 'vs mois dernier', icon: Building2 },
  { title: 'Offres actives', value: '1,456', change: -2.1, changeLabel: 'vs mois dernier', icon: Tag },
  { title: 'Réservations', value: '12,340', change: 18.7, changeLabel: 'vs mois dernier', icon: Calendar },
]

const revenueStats = {
  mrr: 124500,
  arr: 1494000,
  avgRevenue: 9.90,
  churnRate: 4.2,
  conversionRate: 68,
}

// Mock chart data
const monthlyData = [
  { month: 'Jan', users: 25000, bookings: 8500, revenue: 85000 },
  { month: 'Fév', users: 28000, bookings: 9200, revenue: 92000 },
  { month: 'Mar', users: 31000, bookings: 10100, revenue: 98000 },
  { month: 'Avr', users: 33500, bookings: 10800, revenue: 104000 },
  { month: 'Mai', users: 36000, bookings: 11200, revenue: 108000 },
  { month: 'Juin', users: 38500, bookings: 11500, revenue: 112000 },
  { month: 'Juil', users: 40000, bookings: 11800, revenue: 115000 },
  { month: 'Août', users: 41500, bookings: 11600, revenue: 117000 },
  { month: 'Sep', users: 42800, bookings: 11900, revenue: 119000 },
  { month: 'Oct', users: 44200, bookings: 12100, revenue: 121000 },
  { month: 'Nov', users: 45100, bookings: 12250, revenue: 123000 },
  { month: 'Déc', users: 45892, bookings: 12340, revenue: 124500 },
]

const topPartners = [
  { name: 'Le Petit Bistrot', bookings: 1245, revenue: 12450, rating: 4.8 },
  { name: 'Escape Game Paris', bookings: 987, revenue: 9870, rating: 4.6 },
  { name: 'Spa Zen', bookings: 756, revenue: 7560, rating: 4.9 },
  { name: 'Restaurant Gourmet', bookings: 654, revenue: 6540, rating: 4.7 },
  { name: 'Fitness Premium', bookings: 543, revenue: 5430, rating: 4.5 },
]

const topOffers = [
  { title: '-20% sur l\'addition', partner: 'Le Petit Bistrot', bookings: 456, views: 12340 },
  { title: 'Escape Room à 25€', partner: 'Escape Game Paris', bookings: 387, views: 9870 },
  { title: 'Spa détente -30%', partner: 'Spa Zen', bookings: 312, views: 8540 },
  { title: 'Menu découverte', partner: 'Restaurant Gourmet', bookings: 287, views: 7650 },
  { title: 'Séance découverte', partner: 'Fitness Premium', bookings: 234, views: 6780 },
]

function formatCurrency(cents: number) {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: 'EUR',
  }).format(cents / 100)
}

function SimpleBarChart({ data }: { data: typeof monthlyData }) {
  const maxValue = Math.max(...data.map(d => d.users))
  
  return (
    <div className="flex items-end gap-2 h-48">
      {data.map((d, i) => (
        <div key={i} className="flex-1 flex flex-col items-center gap-1">
          <div 
            className="w-full bg-primary/80 rounded-t hover:bg-primary transition-colors"
            style={{ height: `${(d.users / maxValue) * 100}%` }}
            title={`${d.month}: ${d.users.toLocaleString()} utilisateurs`}
          />
          <span className="text-xs text-muted-foreground">{d.month.slice(0, 3)}</span>
        </div>
      ))}
    </div>
  )
}

export function AnalyticsPage() {
  const [period, setPeriod] = useState<'7d' | '30d' | '90d' | '1y'>('30d')

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Analytics</h1>
          <p className="text-muted-foreground">Vue d'ensemble des performances de la plateforme</p>
        </div>
        <div className="flex gap-2">
          <div className="flex border rounded-lg overflow-hidden">
            {(['7d', '30d', '90d', '1y'] as const).map((p) => (
              <button
                key={p}
                className={`px-3 py-1.5 text-sm ${period === p ? 'bg-primary text-primary-foreground' : 'hover:bg-gray-100'}`}
                onClick={() => setPeriod(p)}
              >
                {p === '7d' ? '7 jours' : p === '30d' ? '30 jours' : p === '90d' ? '90 jours' : '1 an'}
              </button>
            ))}
          </div>
          <Button variant="outline" className="gap-2">
            <Download className="h-4 w-4" />
            Exporter
          </Button>
        </div>
      </div>

      {/* Main Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        {stats.map((stat, i) => {
          const Icon = stat.icon
          return (
            <Card key={i}>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                  <Icon className="h-4 w-4" />
                  {stat.title}
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stat.value}</div>
                <p className={`text-xs flex items-center gap-1 mt-1 ${stat.change >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                  {stat.change >= 0 ? <TrendingUp className="h-3 w-3" /> : <TrendingDown className="h-3 w-3" />}
                  {stat.change >= 0 ? '+' : ''}{stat.change}% {stat.changeLabel}
                </p>
              </CardContent>
            </Card>
          )
        })}
      </div>

      {/* Revenue Stats */}
      <div className="grid gap-4 md:grid-cols-5">
        <Card className="bg-gradient-to-br from-primary/10 to-primary/5">
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <DollarSign className="h-4 w-4" />
              MRR
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(revenueStats.mrr * 100)}</div>
            <p className="text-xs text-green-600 flex items-center gap-1 mt-1">
              <TrendingUp className="h-3 w-3" />
              +12% vs mois dernier
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">ARR</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(revenueStats.arr * 100)}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">ARPU</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(revenueStats.avgRevenue * 100)}</div>
            <p className="text-xs text-muted-foreground">/ utilisateur / mois</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Taux conversion</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{revenueStats.conversionRate}%</div>
            <p className="text-xs text-muted-foreground">Essai → Payant</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Churn</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-orange-600">{revenueStats.churnRate}%</div>
            <p className="text-xs text-muted-foreground">/ mois</p>
          </CardContent>
        </Card>
      </div>

      {/* Charts */}
      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Évolution des utilisateurs</CardTitle>
            <CardDescription>Nombre d'utilisateurs inscrits sur 12 mois</CardDescription>
          </CardHeader>
          <CardContent>
            <SimpleBarChart data={monthlyData} />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Répartition des abonnements</CardTitle>
            <CardDescription>Par type de plan</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div>
                <div className="flex justify-between text-sm mb-1">
                  <span>Gratuit</span>
                  <span className="font-medium">15,420 (33.6%)</span>
                </div>
                <div className="h-3 bg-gray-100 rounded-full overflow-hidden">
                  <div className="h-full bg-gray-400 rounded-full" style={{ width: '33.6%' }} />
                </div>
              </div>
              <div>
                <div className="flex justify-between text-sm mb-1">
                  <span>Mensuel</span>
                  <span className="font-medium">23,250 (50.7%)</span>
                </div>
                <div className="h-3 bg-gray-100 rounded-full overflow-hidden">
                  <div className="h-full bg-primary rounded-full" style={{ width: '50.7%' }} />
                </div>
              </div>
              <div>
                <div className="flex justify-between text-sm mb-1">
                  <span>Annuel</span>
                  <span className="font-medium">7,222 (15.7%)</span>
                </div>
                <div className="h-3 bg-gray-100 rounded-full overflow-hidden">
                  <div className="h-full bg-yellow-500 rounded-full" style={{ width: '15.7%' }} />
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Top Lists */}
      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Top partenaires</CardTitle>
            <CardDescription>Par nombre de réservations</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {topPartners.map((partner, i) => (
                <div key={i} className="flex items-center gap-4">
                  <span className="text-lg font-bold text-muted-foreground w-6">#{i + 1}</span>
                  <div className="flex-1">
                    <p className="font-medium">{partner.name}</p>
                    <p className="text-sm text-muted-foreground">
                      {partner.bookings} réservations • {formatCurrency(partner.revenue * 100)} revenus
                    </p>
                  </div>
                  <div className="flex items-center gap-1">
                    <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                    <span className="font-medium">{partner.rating}</span>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Top offres</CardTitle>
            <CardDescription>Par nombre de réservations</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {topOffers.map((offer, i) => (
                <div key={i} className="flex items-center gap-4">
                  <span className="text-lg font-bold text-muted-foreground w-6">#{i + 1}</span>
                  <div className="flex-1">
                    <p className="font-medium">{offer.title}</p>
                    <p className="text-sm text-muted-foreground">{offer.partner}</p>
                  </div>
                  <div className="text-right">
                    <p className="font-medium">{offer.bookings} résa</p>
                    <p className="text-xs text-muted-foreground flex items-center gap-1 justify-end">
                      <Eye className="h-3 w-3" />
                      {offer.views.toLocaleString()}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Quick Links */}
      <Card className="bg-gray-50">
        <CardContent className="pt-6">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="font-medium">Besoin de plus de détails ?</h3>
              <p className="text-sm text-muted-foreground">
                Consultez les rapports détaillés par section
              </p>
            </div>
            <div className="flex gap-2">
              <Button variant="outline" size="sm" className="gap-1">
                Utilisateurs <ArrowUpRight className="h-3 w-3" />
              </Button>
              <Button variant="outline" size="sm" className="gap-1">
                Partenaires <ArrowUpRight className="h-3 w-3" />
              </Button>
              <Button variant="outline" size="sm" className="gap-1">
                Offres <ArrowUpRight className="h-3 w-3" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
