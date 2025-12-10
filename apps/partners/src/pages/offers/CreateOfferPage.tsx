import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { ArrowLeft, ArrowRight, Loader2, Upload, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useToast } from '@/hooks/use-toast'

const offerSchema = z.object({
  // Step 1
  title: z.string().min(3, 'Minimum 3 caractères'),
  shortDescription: z.string().max(100, 'Maximum 100 caractères').optional(),
  description: z.string().min(10, 'Minimum 10 caractères'),
  categoryId: z.string().min(1, 'Catégorie requise'),
  establishmentId: z.string().min(1, 'Établissement requis'),
  
  // Step 2
  discountType: z.enum(['percentage', 'fixed', 'formula']),
  discountValue: z.number().min(1).optional(),
  discountFormula: z.string().optional(),
  conditions: z.string().optional(),
  
  // Step 3
  startDate: z.string().min(1, 'Date de début requise'),
  endDate: z.string().min(1, 'Date de fin requise'),
  allDay: z.boolean().default(true),
  quotaTotal: z.number().optional(),
  quotaPerUser: z.number().optional(),
})

type OfferFormData = z.infer<typeof offerSchema>

const steps = [
  { id: 1, title: 'Informations' },
  { id: 2, title: 'Réduction' },
  { id: 3, title: 'Validité' },
  { id: 4, title: 'Médias' },
  { id: 5, title: 'Prévisualisation' },
]

// Mock categories and establishments
const categories = [
  { id: 'bar', name: 'Bar' },
  { id: 'restaurant', name: 'Restaurant' },
  { id: 'event', name: 'Événement' },
  { id: 'wellness', name: 'Bien-être' },
]

const establishments = [
  { id: 'e1', name: 'Le Petit Bistrot - Paris 1er' },
  { id: 'e2', name: 'Le Petit Bistrot - Paris 11e' },
]

export function CreateOfferPage() {
  const [step, setStep] = useState(1)
  const [isLoading, setIsLoading] = useState(false)
  const [images, setImages] = useState<string[]>([])
  const navigate = useNavigate()
  const { toast } = useToast()

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
    trigger,
  } = useForm<OfferFormData>({
    resolver: zodResolver(offerSchema),
    defaultValues: {
      discountType: 'percentage',
      allDay: true,
    },
  })

  const watchedValues = watch()

  const nextStep = async () => {
    let fieldsToValidate: (keyof OfferFormData)[] = []
    
    switch (step) {
      case 1:
        fieldsToValidate = ['title', 'description', 'categoryId', 'establishmentId']
        break
      case 2:
        fieldsToValidate = ['discountType']
        break
      case 3:
        fieldsToValidate = ['startDate', 'endDate']
        break
    }
    
    const isValid = await trigger(fieldsToValidate)
    if (isValid) setStep(step + 1)
  }

  const prevStep = () => setStep(step - 1)

  const onSubmit = async (_data: OfferFormData) => {
    setIsLoading(true)
    
    try {
      // TODO: Replace with actual GraphQL mutation
      await new Promise((resolve) => setTimeout(resolve, 1500))
      
      toast({
        title: 'Offre créée !',
        description: 'Votre offre a été créée avec succès.',
      })
      
      navigate('/offers')
    } catch {
      toast({
        title: 'Erreur',
        description: 'Impossible de créer l\'offre.',
        variant: 'destructive',
      })
    } finally {
      setIsLoading(false)
    }
  }

  const addImage = () => {
    // Mock adding image - in real app, this would open file picker
    const mockImages = [
      'https://images.unsplash.com/photo-1551024709-8f23befc6f87?w=400&h=300&fit=crop',
      'https://images.unsplash.com/photo-1566417713940-fe7c737a9ef2?w=400&h=300&fit=crop',
      'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=400&h=300&fit=crop',
    ]
    if (images.length < 5) {
      setImages([...images, mockImages[images.length % 3]])
    }
  }

  const removeImage = (index: number) => {
    setImages(images.filter((_, i) => i !== index))
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => navigate('/offers')}>
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div>
          <h1 className="text-2xl font-bold">Créer une offre</h1>
          <p className="text-muted-foreground">
            Étape {step} sur {steps.length} - {steps[step - 1].title}
          </p>
        </div>
      </div>

      {/* Progress */}
      <div className="flex items-center justify-between gap-2">
        {steps.map((s) => (
          <div
            key={s.id}
            className={`flex-1 h-2 rounded-full transition-colors ${
              s.id <= step ? 'bg-primary' : 'bg-muted'
            }`}
          />
        ))}
      </div>

      {/* Form */}
      <form onSubmit={handleSubmit(onSubmit)}>
        {/* Step 1: Informations */}
        {step === 1 && (
          <Card>
            <CardHeader>
              <CardTitle>Informations générales</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="title">Titre de l'offre *</Label>
                <Input
                  id="title"
                  placeholder="Ex: Happy Hour -30%"
                  {...register('title')}
                />
                {errors.title && (
                  <p className="text-sm text-destructive">{errors.title.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="shortDescription">Description courte</Label>
                <Input
                  id="shortDescription"
                  placeholder="Max 100 caractères"
                  maxLength={100}
                  {...register('shortDescription')}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="description">Description complète *</Label>
                <textarea
                  id="description"
                  className="flex min-h-[120px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  placeholder="Décrivez votre offre en détail..."
                  {...register('description')}
                />
                {errors.description && (
                  <p className="text-sm text-destructive">{errors.description.message}</p>
                )}
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="categoryId">Catégorie *</Label>
                  <select
                    id="categoryId"
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                    {...register('categoryId')}
                  >
                    <option value="">Sélectionner...</option>
                    {categories.map((cat) => (
                      <option key={cat.id} value={cat.id}>
                        {cat.name}
                      </option>
                    ))}
                  </select>
                  {errors.categoryId && (
                    <p className="text-sm text-destructive">{errors.categoryId.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="establishmentId">Établissement *</Label>
                  <select
                    id="establishmentId"
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                    {...register('establishmentId')}
                  >
                    <option value="">Sélectionner...</option>
                    {establishments.map((est) => (
                      <option key={est.id} value={est.id}>
                        {est.name}
                      </option>
                    ))}
                  </select>
                  {errors.establishmentId && (
                    <p className="text-sm text-destructive">
                      {errors.establishmentId.message}
                    </p>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Step 2: Discount */}
        {step === 2 && (
          <Card>
            <CardHeader>
              <CardTitle>Réduction</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label>Type de réduction *</Label>
                <div className="grid grid-cols-3 gap-4">
                  {[
                    { value: 'percentage', label: 'Pourcentage' },
                    { value: 'fixed', label: 'Montant fixe' },
                    { value: 'formula', label: 'Formule' },
                  ].map((type) => (
                    <label
                      key={type.value}
                      className={`flex items-center justify-center p-4 rounded-lg border cursor-pointer transition-colors ${
                        watchedValues.discountType === type.value
                          ? 'border-primary bg-primary/5'
                          : 'border-input hover:bg-muted'
                      }`}
                    >
                      <input
                        type="radio"
                        value={type.value}
                        className="sr-only"
                        {...register('discountType')}
                      />
                      <span className="font-medium">{type.label}</span>
                    </label>
                  ))}
                </div>
              </div>

              {watchedValues.discountType === 'percentage' && (
                <div className="space-y-2">
                  <Label htmlFor="discountValue">Pourcentage de réduction *</Label>
                  <div className="relative">
                    <Input
                      id="discountValue"
                      type="number"
                      min={1}
                      max={100}
                      placeholder="30"
                      {...register('discountValue', { valueAsNumber: true })}
                    />
                    <span className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                      %
                    </span>
                  </div>
                </div>
              )}

              {watchedValues.discountType === 'fixed' && (
                <div className="space-y-2">
                  <Label htmlFor="discountValue">Montant de réduction *</Label>
                  <div className="relative">
                    <Input
                      id="discountValue"
                      type="number"
                      min={1}
                      placeholder="10"
                      {...register('discountValue', { valueAsNumber: true })}
                    />
                    <span className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                      €
                    </span>
                  </div>
                </div>
              )}

              {watchedValues.discountType === 'formula' && (
                <div className="space-y-2">
                  <Label htmlFor="discountFormula">Formule *</Label>
                  <Input
                    id="discountFormula"
                    placeholder="Ex: 1 acheté = 1 offert"
                    {...register('discountFormula')}
                  />
                </div>
              )}

              <div className="space-y-2">
                <Label htmlFor="conditions">Conditions (optionnel)</Label>
                <textarea
                  id="conditions"
                  className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  placeholder="Ex: Offre non cumulable, réservée aux nouveaux clients..."
                  {...register('conditions')}
                />
              </div>
            </CardContent>
          </Card>
        )}

        {/* Step 3: Validity */}
        {step === 3 && (
          <Card>
            <CardHeader>
              <CardTitle>Validité et quotas</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="startDate">Date de début *</Label>
                  <Input
                    id="startDate"
                    type="date"
                    {...register('startDate')}
                  />
                  {errors.startDate && (
                    <p className="text-sm text-destructive">{errors.startDate.message}</p>
                  )}
                </div>
                <div className="space-y-2">
                  <Label htmlFor="endDate">Date de fin *</Label>
                  <Input
                    id="endDate"
                    type="date"
                    {...register('endDate')}
                  />
                  {errors.endDate && (
                    <p className="text-sm text-destructive">{errors.endDate.message}</p>
                  )}
                </div>
              </div>

              <div className="flex items-center gap-2">
                <input
                  id="allDay"
                  type="checkbox"
                  className="h-4 w-4 rounded border-input"
                  {...register('allDay')}
                />
                <Label htmlFor="allDay" className="font-normal">
                  Valable toute la journée
                </Label>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="quotaTotal">Quota total (optionnel)</Label>
                  <Input
                    id="quotaTotal"
                    type="number"
                    min={1}
                    placeholder="Illimité"
                    {...register('quotaTotal', { valueAsNumber: true })}
                  />
                  <p className="text-xs text-muted-foreground">
                    Nombre maximum d'utilisations
                  </p>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="quotaPerUser">Quota par utilisateur</Label>
                  <Input
                    id="quotaPerUser"
                    type="number"
                    min={1}
                    placeholder="Illimité"
                    {...register('quotaPerUser', { valueAsNumber: true })}
                  />
                  <p className="text-xs text-muted-foreground">
                    Par personne
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Step 4: Media */}
        {step === 4 && (
          <Card>
            <CardHeader>
              <CardTitle>Images</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
                {images.map((image, index) => (
                  <div key={index} className="relative group">
                    <img
                      src={image}
                      alt={`Image ${index + 1}`}
                      className="w-full h-32 object-cover rounded-lg"
                    />
                    <button
                      type="button"
                      onClick={() => removeImage(index)}
                      className="absolute top-2 right-2 p-1 rounded-full bg-black/50 text-white opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <X className="h-4 w-4" />
                    </button>
                    {index === 0 && (
                      <span className="absolute bottom-2 left-2 px-2 py-1 rounded text-xs bg-primary text-primary-foreground">
                        Principale
                      </span>
                    )}
                  </div>
                ))}
                {images.length < 5 && (
                  <button
                    type="button"
                    onClick={addImage}
                    className="h-32 border-2 border-dashed border-input rounded-lg flex flex-col items-center justify-center gap-2 text-muted-foreground hover:border-primary hover:text-primary transition-colors"
                  >
                    <Upload className="h-6 w-6" />
                    <span className="text-sm">Ajouter</span>
                  </button>
                )}
              </div>
              <p className="text-sm text-muted-foreground">
                La première image sera utilisée comme image principale.
                Maximum 5 images.
              </p>
            </CardContent>
          </Card>
        )}

        {/* Step 5: Preview */}
        {step === 5 && (
          <Card>
            <CardHeader>
              <CardTitle>Prévisualisation</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="bg-muted rounded-lg p-6">
                {images[0] && (
                  <img
                    src={images[0]}
                    alt="Preview"
                    className="w-full h-48 object-cover rounded-lg mb-4"
                  />
                )}
                <h2 className="text-xl font-bold">{watchedValues.title || 'Titre de l\'offre'}</h2>
                <p className="text-muted-foreground mt-2">
                  {watchedValues.description || 'Description...'}
                </p>
                <div className="mt-4 flex items-center gap-4 text-sm">
                  <span className="font-semibold text-primary">
                    {watchedValues.discountType === 'percentage' &&
                      `-${watchedValues.discountValue || 0}%`}
                    {watchedValues.discountType === 'fixed' &&
                      `-${watchedValues.discountValue || 0}€`}
                    {watchedValues.discountType === 'formula' &&
                      (watchedValues.discountFormula || 'Formule')}
                  </span>
                  <span className="text-muted-foreground">
                    Du {watchedValues.startDate || '...'} au {watchedValues.endDate || '...'}
                  </span>
                </div>
              </div>
              <p className="text-sm text-muted-foreground text-center">
                Vérifiez les informations avant de publier votre offre.
              </p>
            </CardContent>
          </Card>
        )}

        {/* Navigation buttons */}
        <div className="flex justify-between mt-6">
          {step > 1 ? (
            <Button type="button" variant="outline" onClick={prevStep}>
              <ArrowLeft className="mr-2 h-4 w-4" />
              Précédent
            </Button>
          ) : (
            <div />
          )}

          {step < 5 ? (
            <Button type="button" onClick={nextStep}>
              Suivant
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          ) : (
            <Button type="submit" disabled={isLoading}>
              {isLoading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Publication...
                </>
              ) : (
                'Publier l\'offre'
              )}
            </Button>
          )}
        </div>
      </form>
    </div>
  )
}
