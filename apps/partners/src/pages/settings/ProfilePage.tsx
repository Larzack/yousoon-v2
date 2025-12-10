import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import {
  User,
  Mail,
  Phone,
  Camera,
  Save,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { useToast } from '@/hooks/use-toast';
import { useAuthStore } from '@/stores/authStore';

// Schema de validation
const profileSchema = z.object({
  firstName: z.string().min(2, 'Prénom requis'),
  lastName: z.string().min(2, 'Nom requis'),
  email: z.string().email('Email invalide'),
  phone: z.string().min(10, 'Numéro de téléphone invalide').optional().or(z.literal('')),
  jobTitle: z.string().optional(),
});

type ProfileFormData = z.infer<typeof profileSchema>;

// Mock data
const mockProfileData = {
  firstName: 'Jean',
  lastName: 'Dupont',
  email: 'jean.dupont@comptoir-parisien.fr',
  phone: '+33 6 12 34 56 78',
  jobTitle: 'Gérant',
  avatar: 'https://randomuser.me/api/portraits/men/1.jpg',
};

export function ProfilePage() {
  const { toast } = useToast();
  const [isLoading, setIsLoading] = useState(false);
  const user = useAuthStore((state) => state.user);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ProfileFormData>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      firstName: user?.firstName || mockProfileData.firstName,
      lastName: user?.lastName || mockProfileData.lastName,
      email: user?.email || mockProfileData.email,
      phone: mockProfileData.phone,
      jobTitle: mockProfileData.jobTitle,
    },
  });

  const onSubmit = async (data: ProfileFormData) => {
    setIsLoading(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      console.log('Profile data:', data);
      toast({
        title: 'Profil mis à jour',
        description: 'Vos informations personnelles ont été enregistrées.',
      });
    } catch (error) {
      toast({
        title: 'Erreur',
        description: 'Une erreur est survenue. Veuillez réessayer.',
        variant: 'destructive',
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-foreground">Mon profil</h1>
        <p className="text-muted-foreground">
          Gérez vos informations personnelles
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Avatar section */}
        <Card>
          <CardHeader>
            <CardTitle>Photo de profil</CardTitle>
          </CardHeader>
          <CardContent className="flex flex-col items-center gap-4">
            <div className="relative">
              <Avatar className="h-32 w-32">
                <AvatarImage src={mockProfileData.avatar} />
                <AvatarFallback className="text-3xl">
                  {mockProfileData.firstName[0]}{mockProfileData.lastName[0]}
                </AvatarFallback>
              </Avatar>
              <Button
                size="icon"
                className="absolute bottom-0 right-0 h-8 w-8 rounded-full bg-yousoon-gold hover:bg-yousoon-gold/90"
              >
                <Camera className="h-4 w-4 text-black" />
              </Button>
            </div>
            <div className="text-center">
              <p className="font-medium">
                {mockProfileData.firstName} {mockProfileData.lastName}
              </p>
              <p className="text-sm text-muted-foreground">{mockProfileData.jobTitle}</p>
            </div>
          </CardContent>
        </Card>

        {/* Form */}
        <div className="lg:col-span-2">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <User className="h-5 w-5" />
                  Informations personnelles
                </CardTitle>
                <CardDescription>
                  Mettez à jour vos informations de profil
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="firstName">Prénom *</Label>
                    <Input
                      id="firstName"
                      {...register('firstName')}
                      placeholder="Jean"
                    />
                    {errors.firstName && (
                      <p className="text-sm text-destructive">{errors.firstName.message}</p>
                    )}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="lastName">Nom *</Label>
                    <Input
                      id="lastName"
                      {...register('lastName')}
                      placeholder="Dupont"
                    />
                    {errors.lastName && (
                      <p className="text-sm text-destructive">{errors.lastName.message}</p>
                    )}
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="jobTitle">Fonction</Label>
                  <Input
                    id="jobTitle"
                    {...register('jobTitle')}
                    placeholder="Gérant, Manager, etc."
                  />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Coordonnées</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="email" className="flex items-center gap-2">
                    <Mail className="h-4 w-4" />
                    Email *
                  </Label>
                  <Input
                    id="email"
                    type="email"
                    {...register('email')}
                    placeholder="jean.dupont@email.com"
                  />
                  {errors.email && (
                    <p className="text-sm text-destructive">{errors.email.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="phone" className="flex items-center gap-2">
                    <Phone className="h-4 w-4" />
                    Téléphone
                  </Label>
                  <Input
                    id="phone"
                    {...register('phone')}
                    placeholder="+33 6 12 34 56 78"
                  />
                  {errors.phone && (
                    <p className="text-sm text-destructive">{errors.phone.message}</p>
                  )}
                </div>
              </CardContent>
            </Card>

            <div className="flex justify-end">
              <Button
                type="submit"
                disabled={isLoading}
                className="bg-yousoon-gold hover:bg-yousoon-gold/90 text-black"
              >
                <Save className="mr-2 h-4 w-4" />
                {isLoading ? 'Enregistrement...' : 'Enregistrer les modifications'}
              </Button>
            </div>
          </form>
        </div>
      </div>

      {/* Account info */}
      <Card>
        <CardHeader>
          <CardTitle>Informations du compte</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 sm:grid-cols-3">
            <div>
              <p className="text-sm text-muted-foreground">Rôle</p>
              <p className="font-medium">Administrateur</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Membre depuis</p>
              <p className="font-medium">Janvier 2023</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Dernière connexion</p>
              <p className="font-medium">Aujourd'hui, 14:32</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Danger zone */}
      <Card className="border-destructive/50">
        <CardHeader>
          <CardTitle className="text-destructive">Zone de danger</CardTitle>
          <CardDescription>
            Actions irréversibles sur votre compte
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium">Supprimer mon compte</p>
              <p className="text-sm text-muted-foreground">
                Cette action est irréversible et supprimera toutes vos données
              </p>
            </div>
            <Button variant="destructive">
              Supprimer le compte
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

export default ProfilePage;
