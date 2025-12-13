import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import {
  Building2,
  Bell,
  Shield,
  CreditCard,
  Save,
  Upload,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useToast } from '@/hooks/use-toast';
import { cn } from '@/lib/utils';

// Schema de validation
const companySchema = z.object({
  companyName: z.string().min(2, 'Raison sociale requise'),
  tradeName: z.string().optional(),
  siret: z.string().regex(/^\d{14}$/, 'SIRET invalide (14 chiffres)'),
  vatNumber: z.string().optional(),
  address: z.string().min(5, 'Adresse requise'),
  postalCode: z.string().regex(/^\d{5}$/, 'Code postal invalide'),
  city: z.string().min(2, 'Ville requise'),
  phone: z.string().min(10, 'Téléphone invalide'),
  email: z.string().email('Email invalide'),
  website: z.string().url('URL invalide').optional().or(z.literal('')),
  description: z.string().optional(),
});

type CompanyFormData = z.infer<typeof companySchema>;

// Tabs
const tabs = [
  { id: 'company', label: 'Entreprise', icon: Building2 },
  { id: 'notifications', label: 'Notifications', icon: Bell },
  { id: 'security', label: 'Sécurité', icon: Shield },
  { id: 'billing', label: 'Facturation', icon: CreditCard },
];

// Mock data
const mockCompanyData = {
  companyName: 'SARL Le Comptoir Parisien',
  tradeName: 'Le Comptoir Parisien',
  siret: '12345678901234',
  vatNumber: 'FR12345678901',
  address: '12 Rue de la Paix',
  postalCode: '75001',
  city: 'Paris',
  phone: '+33 1 23 45 67 89',
  email: 'contact@comptoir-parisien.fr',
  website: 'https://www.comptoir-parisien.fr',
  description: 'Bar-restaurant au cœur de Paris proposant une cuisine française traditionnelle et des cocktails.',
  logo: 'https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?w=200',
};

const notificationSettings = [
  { id: 'booking_new', label: 'Nouvelle réservation', description: 'Recevoir une notification à chaque nouvelle réservation', enabled: true },
  { id: 'booking_cancelled', label: 'Réservation annulée', description: 'Être notifié des annulations', enabled: true },
  { id: 'checkin', label: 'Check-in effectué', description: 'Notification lors d\'un check-in client', enabled: false },
  { id: 'review_new', label: 'Nouvel avis', description: 'Recevoir une notification pour chaque nouvel avis', enabled: true },
  { id: 'marketing', label: 'Actualités Yousoon', description: 'Nouveautés et conseils de la plateforme', enabled: false },
];

export function SettingsPage() {
  const { toast } = useToast();
  const [activeTab, setActiveTab] = useState('company');
  const [isLoading, setIsLoading] = useState(false);
  const [notifications, setNotifications] = useState(notificationSettings);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CompanyFormData>({
    resolver: zodResolver(companySchema),
    defaultValues: mockCompanyData,
  });

  const onSubmitCompany = async (data: CompanyFormData) => {
    setIsLoading(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      console.log('Company data:', data);
      toast({
        title: 'Modifications enregistrées',
        description: 'Les informations de votre entreprise ont été mises à jour.',
      });
    } catch (_error) {
      toast({
        title: 'Erreur',
        description: 'Une erreur est survenue. Veuillez réessayer.',
        variant: 'destructive',
      });
    } finally {
      setIsLoading(false);
    }
  };

  const toggleNotification = (id: string) => {
    setNotifications(notifications.map(n => 
      n.id === id ? { ...n, enabled: !n.enabled } : n
    ));
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-foreground">Paramètres</h1>
        <p className="text-muted-foreground">
          Gérez les paramètres de votre compte partenaire
        </p>
      </div>

      {/* Tabs */}
      <div className="flex flex-wrap gap-2 border-b">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={cn(
              'flex items-center gap-2 px-4 py-2 text-sm font-medium border-b-2 transition-colors',
              activeTab === tab.id
                ? 'border-yousoon-gold text-yousoon-gold'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            )}
          >
            <tab.icon className="h-4 w-4" />
            {tab.label}
          </button>
        ))}
      </div>

      {/* Content */}
      {activeTab === 'company' && (
        <form onSubmit={handleSubmit(onSubmitCompany)} className="space-y-6">
          <div className="grid gap-6 lg:grid-cols-3">
            {/* Main form */}
            <div className="lg:col-span-2 space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Informations légales</CardTitle>
                  <CardDescription>
                    Les informations officielles de votre entreprise
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid gap-4 sm:grid-cols-2">
                    <div className="space-y-2">
                      <Label htmlFor="companyName">Raison sociale *</Label>
                      <Input id="companyName" {...register('companyName')} />
                      {errors.companyName && (
                        <p className="text-sm text-destructive">{errors.companyName.message}</p>
                      )}
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="tradeName">Nom commercial</Label>
                      <Input id="tradeName" {...register('tradeName')} />
                    </div>
                  </div>
                  <div className="grid gap-4 sm:grid-cols-2">
                    <div className="space-y-2">
                      <Label htmlFor="siret">SIRET *</Label>
                      <Input id="siret" {...register('siret')} placeholder="14 chiffres" />
                      {errors.siret && (
                        <p className="text-sm text-destructive">{errors.siret.message}</p>
                      )}
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="vatNumber">N° TVA intracommunautaire</Label>
                      <Input id="vatNumber" {...register('vatNumber')} />
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Adresse</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="address">Adresse *</Label>
                    <Input id="address" {...register('address')} />
                    {errors.address && (
                      <p className="text-sm text-destructive">{errors.address.message}</p>
                    )}
                  </div>
                  <div className="grid gap-4 sm:grid-cols-2">
                    <div className="space-y-2">
                      <Label htmlFor="postalCode">Code postal *</Label>
                      <Input id="postalCode" {...register('postalCode')} />
                      {errors.postalCode && (
                        <p className="text-sm text-destructive">{errors.postalCode.message}</p>
                      )}
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="city">Ville *</Label>
                      <Input id="city" {...register('city')} />
                      {errors.city && (
                        <p className="text-sm text-destructive">{errors.city.message}</p>
                      )}
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Contact</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid gap-4 sm:grid-cols-2">
                    <div className="space-y-2">
                      <Label htmlFor="phone">Téléphone *</Label>
                      <Input id="phone" {...register('phone')} />
                      {errors.phone && (
                        <p className="text-sm text-destructive">{errors.phone.message}</p>
                      )}
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="email">Email *</Label>
                      <Input id="email" type="email" {...register('email')} />
                      {errors.email && (
                        <p className="text-sm text-destructive">{errors.email.message}</p>
                      )}
                    </div>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="website">Site web</Label>
                    <Input id="website" {...register('website')} />
                    {errors.website && (
                      <p className="text-sm text-destructive">{errors.website.message}</p>
                    )}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="description">Description</Label>
                    <textarea
                      id="description"
                      {...register('description')}
                      rows={4}
                      className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                      placeholder="Présentez votre entreprise..."
                    />
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Sidebar */}
            <div className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Logo</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex flex-col items-center gap-4">
                    <div className="w-32 h-32 rounded-lg overflow-hidden bg-muted">
                      <img
                        src={mockCompanyData.logo}
                        alt="Logo"
                        className="w-full h-full object-cover"
                      />
                    </div>
                    <Button variant="outline" size="sm">
                      <Upload className="mr-2 h-4 w-4" />
                      Changer le logo
                    </Button>
                  </div>
                </CardContent>
              </Card>

              <Button
                type="submit"
                disabled={isLoading}
                className="w-full bg-yousoon-gold hover:bg-yousoon-gold/90 text-black"
              >
                <Save className="mr-2 h-4 w-4" />
                {isLoading ? 'Enregistrement...' : 'Enregistrer'}
              </Button>
            </div>
          </div>
        </form>
      )}

      {activeTab === 'notifications' && (
        <Card>
          <CardHeader>
            <CardTitle>Préférences de notification</CardTitle>
            <CardDescription>
              Choisissez quelles notifications vous souhaitez recevoir
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {notifications.map((notification) => (
              <div
                key={notification.id}
                className="flex items-center justify-between p-4 border rounded-lg"
              >
                <div>
                  <p className="font-medium">{notification.label}</p>
                  <p className="text-sm text-muted-foreground">{notification.description}</p>
                </div>
                <button
                  onClick={() => toggleNotification(notification.id)}
                  className={cn(
                    'relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
                    notification.enabled ? 'bg-yousoon-gold' : 'bg-muted'
                  )}
                >
                  <span
                    className={cn(
                      'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                      notification.enabled ? 'translate-x-6' : 'translate-x-1'
                    )}
                  />
                </button>
              </div>
            ))}
          </CardContent>
        </Card>
      )}

      {activeTab === 'security' && (
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Changer le mot de passe</CardTitle>
              <CardDescription>
                Mettez à jour votre mot de passe régulièrement pour plus de sécurité
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="currentPassword">Mot de passe actuel</Label>
                <Input id="currentPassword" type="password" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="newPassword">Nouveau mot de passe</Label>
                <Input id="newPassword" type="password" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="confirmPassword">Confirmer le nouveau mot de passe</Label>
                <Input id="confirmPassword" type="password" />
              </div>
              <Button className="bg-yousoon-gold hover:bg-yousoon-gold/90 text-black">
                Modifier le mot de passe
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Authentification à deux facteurs</CardTitle>
              <CardDescription>
                Ajoutez une couche de sécurité supplémentaire à votre compte
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">2FA activé</p>
                  <p className="text-sm text-muted-foreground">
                    Votre compte est protégé par l'authentification à deux facteurs
                  </p>
                </div>
                <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                  Actif
                </span>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {activeTab === 'billing' && (
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Plan actuel</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-between p-4 border rounded-lg bg-muted/50">
                <div>
                  <p className="text-lg font-semibold">Plan Gratuit</p>
                  <p className="text-sm text-muted-foreground">
                    Accès complet à toutes les fonctionnalités
                  </p>
                </div>
                <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-yousoon-gold text-black">
                  Actif
                </span>
              </div>
              <p className="text-sm text-muted-foreground mt-4">
                Yousoon est actuellement gratuit pour tous les partenaires. 
                Aucun frais ni commission sur les réservations.
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Historique de facturation</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground text-center py-8">
                Aucune facture pour le moment
              </p>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
}

export default SettingsPage;
