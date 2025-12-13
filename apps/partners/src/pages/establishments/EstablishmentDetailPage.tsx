import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import {
  ArrowLeft,
  Save,
  MapPin,
  Phone,
  Mail,
  Globe,
  Clock,
  Upload,
  Trash2,
  Plus,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useToast } from '@/hooks/use-toast';

// Schema de validation
const establishmentSchema = z.object({
  name: z.string().min(2, 'Le nom doit contenir au moins 2 caractères'),
  description: z.string().optional(),
  address: z.string().min(5, 'Adresse requise'),
  postalCode: z.string().regex(/^\d{5}$/, 'Code postal invalide'),
  city: z.string().min(2, 'Ville requise'),
  phone: z.string().min(10, 'Numéro de téléphone invalide'),
  email: z.string().email('Email invalide'),
  website: z.string().url('URL invalide').optional().or(z.literal('')),
});

type EstablishmentFormData = z.infer<typeof establishmentSchema>;

// Mock data pour un établissement
const mockEstablishment = {
  id: '1',
  name: 'Le Comptoir Parisien',
  description: 'Un bar-restaurant au cœur de Paris avec une ambiance chaleureuse et conviviale.',
  address: '12 Rue de la Paix',
  city: 'Paris',
  postalCode: '75001',
  phone: '+33 1 23 45 67 89',
  email: 'contact@comptoir-parisien.fr',
  website: 'https://www.comptoir-parisien.fr',
  image: 'https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?w=600',
  openingHours: [
    { day: 'Lundi', open: '12:00', close: '00:00', closed: false },
    { day: 'Mardi', open: '12:00', close: '00:00', closed: false },
    { day: 'Mercredi', open: '12:00', close: '00:00', closed: false },
    { day: 'Jeudi', open: '12:00', close: '00:00', closed: false },
    { day: 'Vendredi', open: '12:00', close: '02:00', closed: false },
    { day: 'Samedi', open: '18:00', close: '02:00', closed: false },
    { day: 'Dimanche', open: '', close: '', closed: true },
  ],
  features: ['Terrasse', 'WiFi', 'Parking', 'Accessibilité PMR'],
  isActive: true,
};

export function EstablishmentDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { toast } = useToast();
  const [isLoading, setIsLoading] = useState(false);
  const [openingHours, setOpeningHours] = useState(mockEstablishment.openingHours);
  const [features, setFeatures] = useState(mockEstablishment.features);
  const [newFeature, setNewFeature] = useState('');

  const isNewEstablishment = id === 'create';

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<EstablishmentFormData>({
    resolver: zodResolver(establishmentSchema),
    defaultValues: isNewEstablishment
      ? {}
      : {
          name: mockEstablishment.name,
          description: mockEstablishment.description,
          address: mockEstablishment.address,
          postalCode: mockEstablishment.postalCode,
          city: mockEstablishment.city,
          phone: mockEstablishment.phone,
          email: mockEstablishment.email,
          website: mockEstablishment.website,
        },
  });

  const onSubmit = async (data: EstablishmentFormData) => {
    setIsLoading(true);
    try {
      // Simulated API call
      await new Promise((resolve) => setTimeout(resolve, 1000));
      console.log('Form data:', { ...data, openingHours, features });
      
      toast({
        title: isNewEstablishment ? 'Établissement créé' : 'Modifications enregistrées',
        description: isNewEstablishment
          ? 'Votre établissement a été créé avec succès.'
          : 'Les modifications ont été enregistrées.',
      });
      
      navigate('/establishments');
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

  const updateOpeningHours = (index: number, field: 'open' | 'close' | 'closed', value: string | boolean) => {
    const updated = [...openingHours];
    updated[index] = { ...updated[index], [field]: value };
    setOpeningHours(updated);
  };

  const addFeature = () => {
    if (newFeature.trim() && !features.includes(newFeature.trim())) {
      setFeatures([...features, newFeature.trim()]);
      setNewFeature('');
    }
  };

  const removeFeature = (feature: string) => {
    setFeatures(features.filter((f) => f !== feature));
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" onClick={() => navigate('/establishments')}>
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-foreground">
              {isNewEstablishment ? 'Nouvel établissement' : mockEstablishment.name}
            </h1>
            <p className="text-muted-foreground">
              {isNewEstablishment
                ? 'Ajoutez les informations de votre établissement'
                : 'Modifiez les informations de votre établissement'}
            </p>
          </div>
        </div>
        <Button
          onClick={handleSubmit(onSubmit)}
          disabled={isLoading}
          className="bg-yousoon-gold hover:bg-yousoon-gold/90 text-black"
        >
          <Save className="mr-2 h-4 w-4" />
          {isLoading ? 'Enregistrement...' : 'Enregistrer'}
        </Button>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <div className="grid gap-6 lg:grid-cols-3">
          {/* Main content */}
          <div className="lg:col-span-2 space-y-6">
            {/* Informations générales */}
            <Card>
              <CardHeader>
                <CardTitle>Informations générales</CardTitle>
                <CardDescription>
                  Les informations principales de votre établissement
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Nom de l'établissement *</Label>
                  <Input
                    id="name"
                    {...register('name')}
                    placeholder="Ex: Le Comptoir Parisien"
                  />
                  {errors.name && (
                    <p className="text-sm text-destructive">{errors.name.message}</p>
                  )}
                </div>
                <div className="space-y-2">
                  <Label htmlFor="description">Description</Label>
                  <textarea
                    id="description"
                    {...register('description')}
                    rows={4}
                    className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    placeholder="Décrivez votre établissement..."
                  />
                </div>
              </CardContent>
            </Card>

            {/* Adresse */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <MapPin className="h-5 w-5" />
                  Adresse
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="address">Adresse *</Label>
                  <Input
                    id="address"
                    {...register('address')}
                    placeholder="12 Rue de la Paix"
                  />
                  {errors.address && (
                    <p className="text-sm text-destructive">{errors.address.message}</p>
                  )}
                </div>
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="postalCode">Code postal *</Label>
                    <Input
                      id="postalCode"
                      {...register('postalCode')}
                      placeholder="75001"
                    />
                    {errors.postalCode && (
                      <p className="text-sm text-destructive">{errors.postalCode.message}</p>
                    )}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="city">Ville *</Label>
                    <Input
                      id="city"
                      {...register('city')}
                      placeholder="Paris"
                    />
                    {errors.city && (
                      <p className="text-sm text-destructive">{errors.city.message}</p>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Contact */}
            <Card>
              <CardHeader>
                <CardTitle>Contact</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="phone" className="flex items-center gap-2">
                      <Phone className="h-4 w-4" />
                      Téléphone *
                    </Label>
                    <Input
                      id="phone"
                      {...register('phone')}
                      placeholder="+33 1 23 45 67 89"
                    />
                    {errors.phone && (
                      <p className="text-sm text-destructive">{errors.phone.message}</p>
                    )}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="email" className="flex items-center gap-2">
                      <Mail className="h-4 w-4" />
                      Email *
                    </Label>
                    <Input
                      id="email"
                      type="email"
                      {...register('email')}
                      placeholder="contact@example.fr"
                    />
                    {errors.email && (
                      <p className="text-sm text-destructive">{errors.email.message}</p>
                    )}
                  </div>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="website" className="flex items-center gap-2">
                    <Globe className="h-4 w-4" />
                    Site web
                  </Label>
                  <Input
                    id="website"
                    {...register('website')}
                    placeholder="https://www.example.fr"
                  />
                  {errors.website && (
                    <p className="text-sm text-destructive">{errors.website.message}</p>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Horaires */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Clock className="h-5 w-5" />
                  Horaires d'ouverture
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {openingHours.map((hours, index) => (
                    <div key={hours.day} className="flex items-center gap-4">
                      <span className="w-24 text-sm font-medium">{hours.day}</span>
                      <label className="flex items-center gap-2">
                        <input
                          type="checkbox"
                          checked={hours.closed}
                          onChange={(e) => updateOpeningHours(index, 'closed', e.target.checked)}
                          className="h-4 w-4 rounded border-gray-300 text-yousoon-gold focus:ring-yousoon-gold"
                        />
                        <span className="text-sm text-muted-foreground">Fermé</span>
                      </label>
                      {!hours.closed && (
                        <>
                          <Input
                            type="time"
                            value={hours.open}
                            onChange={(e) => updateOpeningHours(index, 'open', e.target.value)}
                            className="w-32"
                          />
                          <span className="text-muted-foreground">à</span>
                          <Input
                            type="time"
                            value={hours.close}
                            onChange={(e) => updateOpeningHours(index, 'close', e.target.value)}
                            className="w-32"
                          />
                        </>
                      )}
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Image */}
            <Card>
              <CardHeader>
                <CardTitle>Photo principale</CardTitle>
              </CardHeader>
              <CardContent>
                {!isNewEstablishment && mockEstablishment.image ? (
                  <div className="relative aspect-video rounded-lg overflow-hidden">
                    <img
                      src={mockEstablishment.image}
                      alt={mockEstablishment.name}
                      className="w-full h-full object-cover"
                    />
                    <Button
                      variant="secondary"
                      size="sm"
                      className="absolute bottom-2 right-2"
                    >
                      <Upload className="mr-2 h-4 w-4" />
                      Changer
                    </Button>
                  </div>
                ) : (
                  <div className="flex aspect-video items-center justify-center rounded-lg border border-dashed border-muted-foreground/25 bg-muted/10">
                    <div className="text-center">
                      <Upload className="mx-auto h-8 w-8 text-muted-foreground" />
                      <p className="mt-2 text-sm text-muted-foreground">
                        Glissez une image ou
                      </p>
                      <Button variant="link" size="sm">
                        parcourir
                      </Button>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Équipements */}
            <Card>
              <CardHeader>
                <CardTitle>Équipements</CardTitle>
                <CardDescription>
                  Indiquez les équipements disponibles
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex flex-wrap gap-2">
                  {features.map((feature) => (
                    <span
                      key={feature}
                      className="inline-flex items-center gap-1 rounded-full bg-muted px-3 py-1 text-sm"
                    >
                      {feature}
                      <button
                        type="button"
                        onClick={() => removeFeature(feature)}
                        className="text-muted-foreground hover:text-foreground"
                      >
                        <Trash2 className="h-3 w-3" />
                      </button>
                    </span>
                  ))}
                </div>
                <div className="flex gap-2">
                  <Input
                    placeholder="Ajouter un équipement"
                    value={newFeature}
                    onChange={(e) => setNewFeature(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addFeature())}
                  />
                  <Button type="button" variant="outline" size="icon" onClick={addFeature}>
                    <Plus className="h-4 w-4" />
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </form>
    </div>
  );
}

export default EstablishmentDetailPage;
