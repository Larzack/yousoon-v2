import { useState } from 'react';
import {
  BarChart3,
  TrendingUp,
  Eye,
  Users,
  Calendar,
  Download,
  Filter,
  ArrowUpRight,
  ArrowDownRight,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { cn } from '@/lib/utils';

// Types
interface StatCard {
  title: string;
  value: string;
  change: number;
  changeLabel: string;
  icon: React.ReactNode;
}

interface ChartData {
  date: string;
  views: number;
  bookings: number;
  checkins: number;
}

// Mock data
const stats: StatCard[] = [
  {
    title: 'Vues totales',
    value: '24,521',
    change: 12.5,
    changeLabel: 'vs mois dernier',
    icon: <Eye className="h-5 w-5" />,
  },
  {
    title: 'Réservations',
    value: '1,234',
    change: 8.2,
    changeLabel: 'vs mois dernier',
    icon: <Calendar className="h-5 w-5" />,
  },
  {
    title: 'Check-ins',
    value: '1,089',
    change: -2.4,
    changeLabel: 'vs mois dernier',
    icon: <Users className="h-5 w-5" />,
  },
  {
    title: 'Taux de conversion',
    value: '5.03%',
    change: 0.8,
    changeLabel: 'vs mois dernier',
    icon: <TrendingUp className="h-5 w-5" />,
  },
];

const chartData: ChartData[] = [
  { date: 'Jan', views: 1200, bookings: 65, checkins: 58 },
  { date: 'Fév', views: 1800, bookings: 89, checkins: 82 },
  { date: 'Mar', views: 2400, bookings: 120, checkins: 108 },
  { date: 'Avr', views: 2100, bookings: 105, checkins: 95 },
  { date: 'Mai', views: 2800, bookings: 142, checkins: 128 },
  { date: 'Juin', views: 3200, bookings: 168, checkins: 152 },
  { date: 'Juil', views: 3800, bookings: 195, checkins: 178 },
  { date: 'Août', views: 4200, bookings: 210, checkins: 192 },
  { date: 'Sep', views: 3600, bookings: 182, checkins: 165 },
  { date: 'Oct', views: 3100, bookings: 158, checkins: 142 },
  { date: 'Nov', views: 2600, bookings: 132, checkins: 118 },
  { date: 'Déc', views: 2200, bookings: 112, checkins: 98 },
];

const topOffers = [
  { id: '1', title: 'Happy Hour -50%', views: 5420, bookings: 312, conversion: 5.76 },
  { id: '2', title: 'Brunch du dimanche', views: 4180, bookings: 245, conversion: 5.86 },
  { id: '3', title: 'Afterwork cocktails', views: 3950, bookings: 198, conversion: 5.01 },
  { id: '4', title: 'Menu découverte', views: 3200, bookings: 156, conversion: 4.88 },
  { id: '5', title: 'Soirée DJ', views: 2890, bookings: 142, conversion: 4.91 },
];

const periods = [
  { value: '7d', label: '7 derniers jours' },
  { value: '30d', label: '30 derniers jours' },
  { value: '90d', label: '3 derniers mois' },
  { value: '12m', label: '12 derniers mois' },
];

export function AnalyticsPage() {
  const [selectedPeriod, setSelectedPeriod] = useState('30d');

  const maxViews = Math.max(...chartData.map((d) => d.views));

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Analytics</h1>
          <p className="text-muted-foreground">
            Suivez les performances de vos offres et établissements
          </p>
        </div>
        <div className="flex gap-2">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline">
                <Filter className="mr-2 h-4 w-4" />
                {periods.find((p) => p.value === selectedPeriod)?.label}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              {periods.map((period) => (
                <DropdownMenuItem
                  key={period.value}
                  onClick={() => setSelectedPeriod(period.value)}
                >
                  {period.label}
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
          <Button variant="outline">
            <Download className="mr-2 h-4 w-4" />
            Exporter
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {stat.title}
              </CardTitle>
              <div className="text-muted-foreground">{stat.icon}</div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <div className="flex items-center gap-1 text-xs">
                {stat.change > 0 ? (
                  <ArrowUpRight className="h-3 w-3 text-green-500" />
                ) : (
                  <ArrowDownRight className="h-3 w-3 text-red-500" />
                )}
                <span
                  className={cn(
                    'font-medium',
                    stat.change > 0 ? 'text-green-500' : 'text-red-500'
                  )}
                >
                  {Math.abs(stat.change)}%
                </span>
                <span className="text-muted-foreground">{stat.changeLabel}</span>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Charts Grid */}
      <div className="grid gap-6 lg:grid-cols-3">
        {/* Main Chart */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <BarChart3 className="h-5 w-5" />
              Évolution des performances
            </CardTitle>
            <CardDescription>
              Vues, réservations et check-ins sur les 12 derniers mois
            </CardDescription>
          </CardHeader>
          <CardContent>
            {/* Simple bar chart visualization */}
            <div className="space-y-4">
              {/* Legend */}
              <div className="flex gap-6">
                <div className="flex items-center gap-2">
                  <div className="h-3 w-3 rounded bg-yousoon-gold" />
                  <span className="text-sm text-muted-foreground">Vues</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="h-3 w-3 rounded bg-blue-500" />
                  <span className="text-sm text-muted-foreground">Réservations</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="h-3 w-3 rounded bg-green-500" />
                  <span className="text-sm text-muted-foreground">Check-ins</span>
                </div>
              </div>

              {/* Chart */}
              <div className="h-64 flex items-end gap-2">
                {chartData.map((data) => (
                  <div key={data.date} className="flex-1 flex flex-col items-center gap-1">
                    <div className="w-full flex gap-0.5 items-end" style={{ height: '200px' }}>
                      <div
                        className="flex-1 bg-yousoon-gold/80 rounded-t transition-all hover:bg-yousoon-gold"
                        style={{ height: `${(data.views / maxViews) * 100}%` }}
                        title={`${data.views} vues`}
                      />
                      <div
                        className="flex-1 bg-blue-500/80 rounded-t transition-all hover:bg-blue-500"
                        style={{ height: `${(data.bookings / maxViews) * 100 * 20}%` }}
                        title={`${data.bookings} réservations`}
                      />
                      <div
                        className="flex-1 bg-green-500/80 rounded-t transition-all hover:bg-green-500"
                        style={{ height: `${(data.checkins / maxViews) * 100 * 20}%` }}
                        title={`${data.checkins} check-ins`}
                      />
                    </div>
                    <span className="text-xs text-muted-foreground">{data.date}</span>
                  </div>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Top Offers */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5" />
              Top 5 offres
            </CardTitle>
            <CardDescription>
              Offres les plus performantes
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {topOffers.map((offer, index) => (
                <div key={offer.id} className="flex items-center gap-3">
                  <span
                    className={cn(
                      'flex h-8 w-8 items-center justify-center rounded-full text-sm font-bold',
                      index === 0 && 'bg-yousoon-gold text-black',
                      index === 1 && 'bg-gray-300 text-gray-800',
                      index === 2 && 'bg-orange-300 text-orange-800',
                      index > 2 && 'bg-muted text-muted-foreground'
                    )}
                  >
                    {index + 1}
                  </span>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate">{offer.title}</p>
                    <p className="text-xs text-muted-foreground">
                      {offer.views.toLocaleString()} vues · {offer.bookings} réservations
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="text-sm font-medium text-green-500">
                      {offer.conversion}%
                    </p>
                    <p className="text-xs text-muted-foreground">conv.</p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Conversion Funnel */}
      <Card>
        <CardHeader>
          <CardTitle>Entonnoir de conversion</CardTitle>
          <CardDescription>
            Visualisez le parcours utilisateur de la vue à la réservation
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between gap-4">
            <FunnelStep
              label="Vues"
              value={24521}
              percentage={100}
              color="bg-yousoon-gold"
            />
            <div className="h-0.5 flex-1 bg-muted" />
            <FunnelStep
              label="Clics"
              value={8965}
              percentage={36.6}
              color="bg-blue-500"
            />
            <div className="h-0.5 flex-1 bg-muted" />
            <FunnelStep
              label="Favoris"
              value={3456}
              percentage={14.1}
              color="bg-purple-500"
            />
            <div className="h-0.5 flex-1 bg-muted" />
            <FunnelStep
              label="Réservations"
              value={1234}
              percentage={5.03}
              color="bg-green-500"
            />
            <div className="h-0.5 flex-1 bg-muted" />
            <FunnelStep
              label="Check-ins"
              value={1089}
              percentage={4.44}
              color="bg-emerald-500"
            />
          </div>
        </CardContent>
      </Card>

      {/* Heatmap Calendar */}
      <Card>
        <CardHeader>
          <CardTitle>Fréquentation par jour</CardTitle>
          <CardDescription>
            Activité des réservations sur les 30 derniers jours
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-7 gap-1">
            {/* Days of week headers */}
            {['Lun', 'Mar', 'Mer', 'Jeu', 'Ven', 'Sam', 'Dim'].map((day) => (
              <div key={day} className="text-center text-xs text-muted-foreground py-2">
                {day}
              </div>
            ))}
            {/* Calendar cells */}
            {Array.from({ length: 35 }).map((_, index) => {
              const intensity = Math.random();
              return (
                <div
                  key={index}
                  className={cn(
                    'aspect-square rounded-sm',
                    intensity < 0.2 && 'bg-muted',
                    intensity >= 0.2 && intensity < 0.4 && 'bg-yousoon-gold/20',
                    intensity >= 0.4 && intensity < 0.6 && 'bg-yousoon-gold/40',
                    intensity >= 0.6 && intensity < 0.8 && 'bg-yousoon-gold/60',
                    intensity >= 0.8 && 'bg-yousoon-gold'
                  )}
                  title={`${Math.floor(intensity * 50)} réservations`}
                />
              );
            })}
          </div>
          <div className="mt-4 flex items-center justify-end gap-2">
            <span className="text-xs text-muted-foreground">Moins</span>
            <div className="flex gap-1">
              <div className="h-3 w-3 rounded-sm bg-muted" />
              <div className="h-3 w-3 rounded-sm bg-yousoon-gold/20" />
              <div className="h-3 w-3 rounded-sm bg-yousoon-gold/40" />
              <div className="h-3 w-3 rounded-sm bg-yousoon-gold/60" />
              <div className="h-3 w-3 rounded-sm bg-yousoon-gold" />
            </div>
            <span className="text-xs text-muted-foreground">Plus</span>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

interface FunnelStepProps {
  label: string;
  value: number;
  percentage: number;
  color: string;
}

function FunnelStep({ label, value, percentage, color }: FunnelStepProps) {
  return (
    <div className="text-center">
      <div
        className={cn('mx-auto mb-2 h-16 w-16 rounded-full flex items-center justify-center', color)}
      >
        <span className="text-sm font-bold text-white">{percentage}%</span>
      </div>
      <p className="text-lg font-semibold">{value.toLocaleString()}</p>
      <p className="text-sm text-muted-foreground">{label}</p>
    </div>
  );
}

export default AnalyticsPage;
