import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Eye, EyeOff, Loader2, ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useToast } from '@/hooks/use-toast'

const registerSchema = z.object({
  // Company info
  companyName: z.string().min(2, 'Raison sociale requise'),
  tradeName: z.string().optional(),
  siret: z.string().length(14, 'SIRET invalide (14 chiffres)'),
  address: z.string().min(5, 'Adresse requise'),
  postalCode: z.string().regex(/^\d{5}$/, 'Code postal invalide'),
  city: z.string().min(2, 'Ville requise'),
  phone: z.string().min(10, 'Téléphone invalide'),
  
  // Admin info
  firstName: z.string().min(2, 'Prénom requis'),
  lastName: z.string().min(2, 'Nom requis'),
  email: z.string().email('Email invalide'),
  password: z.string()
    .min(8, 'Minimum 8 caractères')
    .regex(/[A-Z]/, 'Au moins une majuscule')
    .regex(/[0-9]/, 'Au moins un chiffre'),
  confirmPassword: z.string(),
  
  // Terms
  acceptTerms: z.boolean().refine(val => val === true, 'Vous devez accepter les CGU'),
}).refine(data => data.password === data.confirmPassword, {
  message: 'Les mots de passe ne correspondent pas',
  path: ['confirmPassword'],
})

type RegisterFormData = z.infer<typeof registerSchema>

export function RegisterPage() {
  const [step, setStep] = useState(1)
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const navigate = useNavigate()
  const { toast } = useToast()

  const {
    register,
    handleSubmit,
    trigger,
    formState: { errors },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    mode: 'onChange',
  })

  const nextStep = async () => {
    let fieldsToValidate: (keyof RegisterFormData)[] = []
    
    if (step === 1) {
      fieldsToValidate = ['companyName', 'siret', 'address', 'postalCode', 'city', 'phone']
    }
    
    const isValid = await trigger(fieldsToValidate)
    if (isValid) setStep(step + 1)
  }

  const prevStep = () => setStep(step - 1)

  const onSubmit = async (_data: RegisterFormData) => {
    setIsLoading(true)
    
    try {
      // TODO: Replace with actual GraphQL mutation
      await new Promise(resolve => setTimeout(resolve, 1500))
      
      toast({
        title: 'Inscription réussie !',
        description: 'Vérifiez votre email pour activer votre compte.',
      })
      
      navigate('/login')
    } catch {
      toast({
        title: 'Erreur',
        description: 'Une erreur est survenue. Veuillez réessayer.',
        variant: 'destructive',
      })
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="space-y-2 text-center">
        <h1 className="text-2xl font-bold tracking-tight">Créer un compte</h1>
        <p className="text-muted-foreground">
          Rejoignez Yousoon et développez votre activité
        </p>
      </div>

      {/* Progress */}
      <div className="flex items-center justify-center gap-2">
        {[1, 2].map((s) => (
          <div
            key={s}
            className={`h-2 w-16 rounded-full transition-colors ${
              s <= step ? 'bg-primary' : 'bg-muted'
            }`}
          />
        ))}
      </div>

      {/* Form */}
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {step === 1 && (
          <>
            <h2 className="text-lg font-semibold">Informations entreprise</h2>
            
            <div className="space-y-2">
              <Label htmlFor="companyName">Raison sociale *</Label>
              <Input
                id="companyName"
                placeholder="Ma Société SAS"
                {...register('companyName')}
              />
              {errors.companyName && (
                <p className="text-sm text-destructive">{errors.companyName.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="tradeName">Nom commercial</Label>
              <Input
                id="tradeName"
                placeholder="Le Petit Bistrot (optionnel)"
                {...register('tradeName')}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="siret">SIRET *</Label>
              <Input
                id="siret"
                placeholder="12345678901234"
                maxLength={14}
                {...register('siret')}
              />
              {errors.siret && (
                <p className="text-sm text-destructive">{errors.siret.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="address">Adresse *</Label>
              <Input
                id="address"
                placeholder="123 Rue de la Paix"
                {...register('address')}
              />
              {errors.address && (
                <p className="text-sm text-destructive">{errors.address.message}</p>
              )}
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="postalCode">Code postal *</Label>
                <Input
                  id="postalCode"
                  placeholder="75001"
                  maxLength={5}
                  {...register('postalCode')}
                />
                {errors.postalCode && (
                  <p className="text-sm text-destructive">{errors.postalCode.message}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="city">Ville *</Label>
                <Input
                  id="city"
                  placeholder="Paris"
                  {...register('city')}
                />
                {errors.city && (
                  <p className="text-sm text-destructive">{errors.city.message}</p>
                )}
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="phone">Téléphone *</Label>
              <Input
                id="phone"
                type="tel"
                placeholder="01 23 45 67 89"
                {...register('phone')}
              />
              {errors.phone && (
                <p className="text-sm text-destructive">{errors.phone.message}</p>
              )}
            </div>

            <Button type="button" className="w-full" onClick={nextStep}>
              Continuer
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </>
        )}

        {step === 2 && (
          <>
            <h2 className="text-lg font-semibold">Vos informations</h2>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="firstName">Prénom *</Label>
                <Input
                  id="firstName"
                  placeholder="Jean"
                  {...register('firstName')}
                />
                {errors.firstName && (
                  <p className="text-sm text-destructive">{errors.firstName.message}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="lastName">Nom *</Label>
                <Input
                  id="lastName"
                  placeholder="Dupont"
                  {...register('lastName')}
                />
                {errors.lastName && (
                  <p className="text-sm text-destructive">{errors.lastName.message}</p>
                )}
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="email">Email *</Label>
              <Input
                id="email"
                type="email"
                placeholder="jean.dupont@example.com"
                {...register('email')}
              />
              {errors.email && (
                <p className="text-sm text-destructive">{errors.email.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="password">Mot de passe *</Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  placeholder="••••••••"
                  {...register('password')}
                />
                <button
                  type="button"
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? (
                    <EyeOff className="h-4 w-4" />
                  ) : (
                    <Eye className="h-4 w-4" />
                  )}
                </button>
              </div>
              {errors.password && (
                <p className="text-sm text-destructive">{errors.password.message}</p>
              )}
              <p className="text-xs text-muted-foreground">
                8 caractères minimum, une majuscule, un chiffre
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="confirmPassword">Confirmer le mot de passe *</Label>
              <Input
                id="confirmPassword"
                type="password"
                placeholder="••••••••"
                {...register('confirmPassword')}
              />
              {errors.confirmPassword && (
                <p className="text-sm text-destructive">{errors.confirmPassword.message}</p>
              )}
            </div>

            <div className="flex items-start gap-2">
              <input
                id="acceptTerms"
                type="checkbox"
                className="h-4 w-4 rounded border-input mt-0.5"
                {...register('acceptTerms')}
              />
              <Label htmlFor="acceptTerms" className="text-sm font-normal leading-tight">
                J'accepte les{' '}
                <a href="#" className="text-primary hover:underline">
                  conditions générales d'utilisation
                </a>{' '}
                et la{' '}
                <a href="#" className="text-primary hover:underline">
                  politique de confidentialité
                </a>
              </Label>
            </div>
            {errors.acceptTerms && (
              <p className="text-sm text-destructive">{errors.acceptTerms.message}</p>
            )}

            <div className="flex gap-3">
              <Button type="button" variant="outline" onClick={prevStep}>
                <ArrowLeft className="mr-2 h-4 w-4" />
                Retour
              </Button>
              <Button type="submit" className="flex-1" disabled={isLoading}>
                {isLoading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Création...
                  </>
                ) : (
                  "Créer mon compte"
                )}
              </Button>
            </div>
          </>
        )}
      </form>

      {/* Login link */}
      <p className="text-center text-sm text-muted-foreground">
        Déjà un compte ?{' '}
        <Link to="/login" className="text-primary hover:underline">
          Se connecter
        </Link>
      </p>
    </div>
  )
}
