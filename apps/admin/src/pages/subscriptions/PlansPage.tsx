import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'
import {
  ArrowLeft,
  Plus,
  Edit,
  Trash2,
  Check,
  Star,
  Crown,
  Zap,
} from 'lucide-react'

interface SubscriptionPlan {
  id: string
  code: string
  name: string
  description: string
  price: number // in cents
  interval: 'month' | 'year'
  trialDays: number
  features: string[]
  isActive: boolean
  isHighlighted: boolean
  subscriberCount: number
}

// Mock data
const mockPlans: SubscriptionPlan[] = [
  {
    id: '1',
    code: 'free',
    name: 'Gratuit',
    description: 'Découvrez Yousoon gratuitement',
    price: 0,
    interval: 'month',
    trialDays: 0,
    features: [
      '5 réservations par mois',
      'Accès aux offres standards',
    ],
    isActive: true,
    isHighlighted: false,
    subscriberCount: 15420,
  },
  {
    id: '2',
    code: 'monthly',
    name: 'Mensuel',
    description: 'Abonnement mensuel sans engagement',
    price: 990,
    interval: 'month',
    trialDays: 30,
    features: [
      'Réservations illimitées',
      'Accès à toutes les offres',
      'Offres exclusives',
      'Support prioritaire',
    ],
    isActive: true,
    isHighlighted: true,
    subscriberCount: 3250,
  },
  {
    id: '3',
    code: 'yearly',
    name: 'Annuel',
    description: '2 mois offerts !',
    price: 7990,
    interval: 'year',
    trialDays: 30,
    features: [
      'Réservations illimitées',
      'Accès à toutes les offres',
      'Offres exclusives',
      'Support prioritaire',
      'Badge VIP',
      '2 mois gratuits',
    ],
    isActive: true,
    isHighlighted: false,
    subscriberCount: 890,
  },
]

function formatPrice(cents: number) {
  if (cents === 0) return 'Gratuit'
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: 'EUR',
  }).format(cents / 100)
}

function PlanIcon({ code }: { code: string }) {
  if (code === 'free') return <Star className="h-8 w-8 text-gray-400" />
  if (code === 'monthly') return <Zap className="h-8 w-8 text-primary" />
  return <Crown className="h-8 w-8 text-yellow-500" />
}

export function PlansPage() {
  const navigate = useNavigate()
  const [plans, setPlans] = useState(mockPlans)
  const [editingPlan, setEditingPlan] = useState<SubscriptionPlan | null>(null)

  const togglePlanActive = (planId: string) => {
    setPlans(plans.map(p => 
      p.id === planId ? { ...p, isActive: !p.isActive } : p
    ))
  }

  const togglePlanHighlighted = (planId: string) => {
    setPlans(plans.map(p => 
      p.id === planId ? { ...p, isHighlighted: !p.isHighlighted } : p
    ))
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => navigate('/subscriptions')}>
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div className="flex-1">
          <h1 className="text-2xl font-bold">Plans d'abonnement</h1>
          <p className="text-muted-foreground">Configuration des plans disponibles</p>
        </div>
        <Button className="gap-2">
          <Plus className="h-4 w-4" />
          Nouveau plan
        </Button>
      </div>

      {/* Plans Grid */}
      <div className="grid gap-6 md:grid-cols-3">
        {plans.map((plan) => (
          <Card 
            key={plan.id} 
            className={`relative ${plan.isHighlighted ? 'border-primary ring-1 ring-primary' : ''} ${!plan.isActive ? 'opacity-60' : ''}`}
          >
            {plan.isHighlighted && (
              <div className="absolute -top-3 left-1/2 -translate-x-1/2">
                <span className="bg-primary text-primary-foreground px-3 py-1 rounded-full text-xs font-medium">
                  Populaire
                </span>
              </div>
            )}
            <CardHeader className="text-center pt-8">
              <div className="mx-auto mb-4">
                <PlanIcon code={plan.code} />
              </div>
              <CardTitle>{plan.name}</CardTitle>
              <CardDescription>{plan.description}</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="text-center">
                <span className="text-4xl font-bold">{formatPrice(plan.price)}</span>
                {plan.price > 0 && (
                  <span className="text-muted-foreground">/{plan.interval === 'month' ? 'mois' : 'an'}</span>
                )}
              </div>

              {plan.trialDays > 0 && (
                <p className="text-center text-sm text-primary">
                  {plan.trialDays} jours d'essai gratuit
                </p>
              )}

              <ul className="space-y-2">
                {plan.features.map((feature, i) => (
                  <li key={i} className="flex items-center gap-2 text-sm">
                    <Check className="h-4 w-4 text-green-500 flex-shrink-0" />
                    {feature}
                  </li>
                ))}
              </ul>

              <div className="pt-4 border-t space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Abonnés</span>
                  <span className="font-bold">{plan.subscriberCount.toLocaleString()}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Actif</span>
                  <Switch 
                    checked={plan.isActive}
                    onCheckedChange={() => togglePlanActive(plan.id)}
                  />
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Mis en avant</span>
                  <Switch 
                    checked={plan.isHighlighted}
                    onCheckedChange={() => togglePlanHighlighted(plan.id)}
                  />
                </div>
              </div>

              <div className="flex gap-2">
                <Button variant="outline" className="flex-1 gap-2" onClick={() => setEditingPlan(plan)}>
                  <Edit className="h-4 w-4" />
                  Modifier
                </Button>
                {plan.code !== 'free' && (
                  <Button variant="ghost" size="icon" className="text-red-600">
                    <Trash2 className="h-4 w-4" />
                  </Button>
                )}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Edit Modal */}
      {editingPlan && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <Card className="w-full max-w-lg">
            <CardHeader>
              <CardTitle>Modifier le plan "{editingPlan.name}"</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <label className="text-sm font-medium">Nom</label>
                <Input defaultValue={editingPlan.name} className="mt-1" />
              </div>
              <div>
                <label className="text-sm font-medium">Description</label>
                <Input defaultValue={editingPlan.description} className="mt-1" />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium">Prix (centimes)</label>
                  <Input type="number" defaultValue={editingPlan.price} className="mt-1" />
                </div>
                <div>
                  <label className="text-sm font-medium">Jours d'essai</label>
                  <Input type="number" defaultValue={editingPlan.trialDays} className="mt-1" />
                </div>
              </div>
              <div>
                <label className="text-sm font-medium">Fonctionnalités (une par ligne)</label>
                <textarea 
                  className="w-full mt-1 p-2 border rounded-md min-h-[100px] text-sm"
                  defaultValue={editingPlan.features.join('\n')}
                />
              </div>
              <div className="flex gap-2 pt-4">
                <Button className="flex-1">Enregistrer</Button>
                <Button variant="outline" onClick={() => setEditingPlan(null)}>Annuler</Button>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Info */}
      <Card className="bg-blue-50 border-blue-200">
        <CardContent className="pt-6">
          <h3 className="font-medium text-blue-800">Note importante</h3>
          <p className="text-sm text-blue-600 mt-1">
            Les modifications de prix n'affectent que les nouveaux abonnés. Les abonnements existants 
            conservent leur tarif jusqu'à renouvellement. Les paiements sont gérés à 100% via 
            Apple Pay / Google Pay (In-App Purchase).
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
