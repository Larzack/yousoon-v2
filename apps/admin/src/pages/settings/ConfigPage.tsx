import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'
import {
  Save,
  RefreshCw,
  AlertTriangle,
  Clock,
  MapPin,
  Shield,
  Globe,
  Bell,
  Smartphone,
} from 'lucide-react'

interface ConfigSection {
  id: string
  title: string
  description: string
  icon: typeof Clock
  settings: ConfigSetting[]
}

interface ConfigSetting {
  key: string
  label: string
  description: string
  type: 'number' | 'boolean' | 'text' | 'select'
  value: string | number | boolean
  options?: string[]
  unit?: string
}

const configSections: ConfigSection[] = [
  {
    id: 'trial',
    title: 'Période d\'essai',
    description: 'Configuration de la période d\'essai pour les nouveaux utilisateurs',
    icon: Clock,
    settings: [
      { key: 'trial_duration', label: 'Durée de l\'essai', description: 'Nombre de jours d\'essai gratuit', type: 'number', value: 30, unit: 'jours' },
      { key: 'trial_enabled', label: 'Essai activé', description: 'Activer la période d\'essai pour les nouveaux inscrits', type: 'boolean', value: true },
    ],
  },
  {
    id: 'booking',
    title: 'Réservations',
    description: 'Configuration des réservations',
    icon: Clock,
    settings: [
      { key: 'booking_expiry', label: 'Expiration QR Code', description: 'Délai avant expiration du QR code après réservation', type: 'number', value: 30, unit: 'minutes' },
      { key: 'max_active_bookings', label: 'Réservations simultanées', description: 'Nombre max de réservations actives par utilisateur', type: 'number', value: 5 },
    ],
  },
  {
    id: 'search',
    title: 'Recherche',
    description: 'Configuration de la recherche et géolocalisation',
    icon: MapPin,
    settings: [
      { key: 'default_radius', label: 'Rayon par défaut', description: 'Rayon de recherche par défaut', type: 'number', value: 10, unit: 'km' },
      { key: 'max_radius', label: 'Rayon maximum', description: 'Rayon de recherche maximum autorisé', type: 'number', value: 50, unit: 'km' },
    ],
  },
  {
    id: 'verification',
    title: 'Vérification d\'identité',
    description: 'Configuration de la vérification CNI',
    icon: Shield,
    settings: [
      { key: 'max_verification_attempts', label: 'Tentatives max', description: 'Nombre maximum de tentatives de vérification', type: 'number', value: 10 },
      { key: 'verification_required', label: 'Vérification obligatoire', description: 'Exiger la vérification CNI pour réserver', type: 'boolean', value: true },
    ],
  },
  {
    id: 'app',
    title: 'Application mobile',
    description: 'Configuration de l\'application',
    icon: Smartphone,
    settings: [
      { key: 'force_update_ios', label: 'Forcer MAJ iOS', description: 'Forcer la mise à jour sur iOS', type: 'boolean', value: false },
      { key: 'force_update_android', label: 'Forcer MAJ Android', description: 'Forcer la mise à jour sur Android', type: 'boolean', value: false },
      { key: 'min_version_ios', label: 'Version min iOS', description: 'Version minimum requise sur iOS', type: 'text', value: '1.0.0' },
      { key: 'min_version_android', label: 'Version min Android', description: 'Version minimum requise sur Android', type: 'text', value: '1.0.0' },
    ],
  },
  {
    id: 'notifications',
    title: 'Notifications',
    description: 'Configuration des notifications',
    icon: Bell,
    settings: [
      { key: 'reminder_hours', label: 'Rappel avant réservation', description: 'Heures avant la réservation pour envoyer un rappel', type: 'number', value: 2, unit: 'heures' },
      { key: 'marketing_enabled', label: 'Notifications marketing', description: 'Autoriser les notifications marketing', type: 'boolean', value: true },
    ],
  },
  {
    id: 'i18n',
    title: 'Internationalisation',
    description: 'Configuration des langues',
    icon: Globe,
    settings: [
      { key: 'default_language', label: 'Langue par défaut', description: 'Langue par défaut de l\'application', type: 'select', value: 'fr', options: ['fr', 'en'] },
      { key: 'auto_translate', label: 'Traduction auto', description: 'Activer la traduction automatique des offres', type: 'boolean', value: true },
    ],
  },
]

export function ConfigPage() {
  const [sections, setSections] = useState(configSections)
  const [hasChanges, setHasChanges] = useState(false)

  const updateSetting = (sectionId: string, key: string, value: string | number | boolean) => {
    setSections(sections.map(section => {
      if (section.id === sectionId) {
        return {
          ...section,
          settings: section.settings.map(setting => 
            setting.key === key ? { ...setting, value } : setting
          )
        }
      }
      return section
    }))
    setHasChanges(true)
  }

  const handleSave = () => {
    // In real app, call API
    alert('Configuration sauvegardée !')
    setHasChanges(false)
  }

  const handleReset = () => {
    if (confirm('Voulez-vous réinitialiser tous les paramètres ?')) {
      setSections(configSections)
      setHasChanges(false)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Configuration</h1>
          <p className="text-muted-foreground">Paramètres généraux de la plateforme</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" className="gap-2" onClick={handleReset}>
            <RefreshCw className="h-4 w-4" />
            Réinitialiser
          </Button>
          <Button className="gap-2" onClick={handleSave} disabled={!hasChanges}>
            <Save className="h-4 w-4" />
            Sauvegarder
          </Button>
        </div>
      </div>

      {hasChanges && (
        <div className="flex items-center gap-3 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
          <AlertTriangle className="h-5 w-5 text-yellow-600" />
          <p className="text-sm text-yellow-800">
            Vous avez des modifications non sauvegardées.
          </p>
        </div>
      )}

      <div className="space-y-6">
        {sections.map((section) => {
          const Icon = section.icon
          return (
            <Card key={section.id}>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Icon className="h-5 w-5" />
                  {section.title}
                </CardTitle>
                <CardDescription>{section.description}</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                {section.settings.map((setting) => (
                  <div key={setting.key} className="flex items-center justify-between py-2 border-b last:border-0">
                    <div className="flex-1">
                      <label className="text-sm font-medium">{setting.label}</label>
                      <p className="text-sm text-muted-foreground">{setting.description}</p>
                    </div>
                    <div className="ml-4">
                      {setting.type === 'boolean' && (
                        <Switch
                          checked={setting.value as boolean}
                          onCheckedChange={(checked) => updateSetting(section.id, setting.key, checked)}
                        />
                      )}
                      {setting.type === 'number' && (
                        <div className="flex items-center gap-2">
                          <Input
                            type="number"
                            value={setting.value as number}
                            onChange={(e) => updateSetting(section.id, setting.key, parseInt(e.target.value) || 0)}
                            className="w-24"
                          />
                          {setting.unit && (
                            <span className="text-sm text-muted-foreground">{setting.unit}</span>
                          )}
                        </div>
                      )}
                      {setting.type === 'text' && (
                        <Input
                          value={setting.value as string}
                          onChange={(e) => updateSetting(section.id, setting.key, e.target.value)}
                          className="w-32"
                        />
                      )}
                      {setting.type === 'select' && setting.options && (
                        <select
                          value={setting.value as string}
                          onChange={(e) => updateSetting(section.id, setting.key, e.target.value)}
                          className="px-3 py-2 border rounded-md"
                        >
                          {setting.options.map((option) => (
                            <option key={option} value={option}>
                              {option.toUpperCase()}
                            </option>
                          ))}
                        </select>
                      )}
                    </div>
                  </div>
                ))}
              </CardContent>
            </Card>
          )
        })}
      </div>

      {/* Warning */}
      <Card className="bg-red-50 border-red-200">
        <CardContent className="pt-6">
          <h3 className="font-medium text-red-800 flex items-center gap-2">
            <AlertTriangle className="h-5 w-5" />
            Attention
          </h3>
          <p className="text-sm text-red-600 mt-1">
            Certains paramètres peuvent affecter l'expérience utilisateur en temps réel.
            Testez les modifications en environnement de staging avant de les appliquer en production.
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
