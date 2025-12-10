import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Plus,
  Edit,
  Trash2,
  MoreHorizontal,
  GripVertical,
  Utensils,
  Wine,
  Dumbbell,
  Music,
  Palette,
  Mountain,
  Sparkles,
  GamepadIcon,
  PartyPopper,
  ShoppingBag,
} from 'lucide-react'

interface Category {
  id: string
  name: {
    fr: string
    en: string
  }
  slug: string
  icon: string
  color: string
  parentId: string | null
  order: number
  isActive: boolean
  offersCount: number
}

// Mock data
const mockCategories: Category[] = [
  { id: '1', name: { fr: 'Restaurant', en: 'Restaurant' }, slug: 'restaurant', icon: 'utensils', color: '#FF6B6B', parentId: null, order: 1, isActive: true, offersCount: 456 },
  { id: '2', name: { fr: 'Bar', en: 'Bar' }, slug: 'bar', icon: 'wine', color: '#4ECDC4', parentId: null, order: 2, isActive: true, offersCount: 234 },
  { id: '3', name: { fr: 'Sport', en: 'Sport' }, slug: 'sport', icon: 'dumbbell', color: '#45B7D1', parentId: null, order: 3, isActive: true, offersCount: 123 },
  { id: '4', name: { fr: 'Loisirs', en: 'Entertainment' }, slug: 'loisirs', icon: 'gamepad', color: '#96CEB4', parentId: null, order: 4, isActive: true, offersCount: 345 },
  { id: '5', name: { fr: 'Bien-être', en: 'Wellness' }, slug: 'bien-etre', icon: 'sparkles', color: '#DDA0DD', parentId: null, order: 5, isActive: true, offersCount: 167 },
  { id: '6', name: { fr: 'Concert & Musique', en: 'Concert & Music' }, slug: 'concert-musique', icon: 'music', color: '#FFD93D', parentId: null, order: 6, isActive: true, offersCount: 89 },
  { id: '7', name: { fr: 'Arts & Culture', en: 'Arts & Culture' }, slug: 'arts-culture', icon: 'palette', color: '#FF8C42', parentId: null, order: 7, isActive: true, offersCount: 112 },
  { id: '8', name: { fr: 'Nature', en: 'Nature' }, slug: 'nature', icon: 'mountain', color: '#6BCB77', parentId: null, order: 8, isActive: true, offersCount: 78 },
  { id: '9', name: { fr: 'Événements', en: 'Events' }, slug: 'evenements', icon: 'party', color: '#FF69B4', parentId: null, order: 9, isActive: false, offersCount: 45 },
  { id: '10', name: { fr: 'Shopping', en: 'Shopping' }, slug: 'shopping', icon: 'shopping', color: '#B19CD9', parentId: null, order: 10, isActive: false, offersCount: 23 },
]

function getIconComponent(icon: string) {
  const icons: Record<string, typeof Utensils> = {
    utensils: Utensils,
    wine: Wine,
    dumbbell: Dumbbell,
    music: Music,
    palette: Palette,
    mountain: Mountain,
    sparkles: Sparkles,
    gamepad: GamepadIcon,
    party: PartyPopper,
    shopping: ShoppingBag,
  }
  const Icon = icons[icon] || Utensils
  return Icon
}

export function CategoriesPage() {
  const [categories, setCategories] = useState(mockCategories)
  const [editingCategory, setEditingCategory] = useState<Category | null>(null)
  const [isCreating, setIsCreating] = useState(false)

  const toggleActive = (id: string) => {
    setCategories(categories.map(c => 
      c.id === id ? { ...c, isActive: !c.isActive } : c
    ))
  }

  const deleteCategory = (id: string) => {
    if (confirm('Êtes-vous sûr de vouloir supprimer cette catégorie ?')) {
      setCategories(categories.filter(c => c.id !== id))
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Catégories</h1>
          <p className="text-muted-foreground">Gestion des catégories d'offres</p>
        </div>
        <Button className="gap-2" onClick={() => setIsCreating(true)}>
          <Plus className="h-4 w-4" />
          Nouvelle catégorie
        </Button>
      </div>

      {/* Categories Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {categories.map((category) => {
          const Icon = getIconComponent(category.icon)
          return (
            <Card 
              key={category.id} 
              className={`${!category.isActive ? 'opacity-60' : ''}`}
            >
              <CardContent className="pt-6">
                <div className="flex items-start gap-4">
                  <div className="flex items-center gap-2">
                    <GripVertical className="h-5 w-5 text-muted-foreground cursor-grab" />
                    <div 
                      className="h-12 w-12 rounded-lg flex items-center justify-center"
                      style={{ backgroundColor: category.color + '20' }}
                    >
                      <Icon className="h-6 w-6" style={{ color: category.color }} />
                    </div>
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <h3 className="font-medium">{category.name.fr}</h3>
                      {!category.isActive && (
                        <span className="text-xs bg-gray-100 px-2 py-0.5 rounded">Inactif</span>
                      )}
                    </div>
                    <p className="text-sm text-muted-foreground">{category.name.en}</p>
                    <p className="text-xs text-muted-foreground mt-1">
                      {category.offersCount} offres • /{category.slug}
                    </p>
                  </div>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="icon">
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem onClick={() => setEditingCategory(category)} className="gap-2">
                        <Edit className="h-4 w-4" />
                        Modifier
                      </DropdownMenuItem>
                      <DropdownMenuItem onClick={() => toggleActive(category.id)} className="gap-2">
                        {category.isActive ? 'Désactiver' : 'Activer'}
                      </DropdownMenuItem>
                      <DropdownMenuItem 
                        onClick={() => deleteCategory(category.id)} 
                        className="gap-2 text-red-600"
                        disabled={category.offersCount > 0}
                      >
                        <Trash2 className="h-4 w-4" />
                        Supprimer
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </CardContent>
            </Card>
          )
        })}
      </div>

      {/* Edit/Create Modal */}
      {(editingCategory || isCreating) && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <Card className="w-full max-w-lg">
            <CardHeader>
              <CardTitle>{isCreating ? 'Nouvelle catégorie' : 'Modifier la catégorie'}</CardTitle>
              <CardDescription>
                {isCreating ? 'Créez une nouvelle catégorie d\'offres' : 'Modifiez les informations de la catégorie'}
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium">Nom (FR)</label>
                  <Input 
                    defaultValue={editingCategory?.name.fr || ''} 
                    placeholder="Restaurant"
                    className="mt-1" 
                  />
                </div>
                <div>
                  <label className="text-sm font-medium">Nom (EN)</label>
                  <Input 
                    defaultValue={editingCategory?.name.en || ''} 
                    placeholder="Restaurant"
                    className="mt-1" 
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium">Slug</label>
                  <Input 
                    defaultValue={editingCategory?.slug || ''} 
                    placeholder="restaurant"
                    className="mt-1" 
                  />
                </div>
                <div>
                  <label className="text-sm font-medium">Couleur</label>
                  <div className="flex gap-2 mt-1">
                    <Input 
                      type="color"
                      defaultValue={editingCategory?.color || '#FF6B6B'} 
                      className="w-12 h-9 p-1"
                    />
                    <Input 
                      defaultValue={editingCategory?.color || '#FF6B6B'} 
                      placeholder="#FF6B6B"
                      className="flex-1"
                    />
                  </div>
                </div>
              </div>
              <div>
                <label className="text-sm font-medium">Icône</label>
                <div className="grid grid-cols-5 gap-2 mt-2">
                  {['utensils', 'wine', 'dumbbell', 'music', 'palette', 'mountain', 'sparkles', 'gamepad', 'party', 'shopping'].map((icon) => {
                    const IconComp = getIconComponent(icon)
                    return (
                      <button
                        key={icon}
                        className={`p-3 rounded-lg border hover:border-primary transition-colors ${
                          editingCategory?.icon === icon ? 'border-primary bg-primary/10' : ''
                        }`}
                      >
                        <IconComp className="h-5 w-5 mx-auto" />
                      </button>
                    )
                  })}
                </div>
              </div>
              <div className="flex gap-2 pt-4">
                <Button className="flex-1">
                  {isCreating ? 'Créer' : 'Enregistrer'}
                </Button>
                <Button 
                  variant="outline" 
                  onClick={() => {
                    setEditingCategory(null)
                    setIsCreating(false)
                  }}
                >
                  Annuler
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Info */}
      <Card className="bg-blue-50 border-blue-200">
        <CardContent className="pt-6">
          <h3 className="font-medium text-blue-800">Ordre d'affichage</h3>
          <p className="text-sm text-blue-600 mt-1">
            Glissez-déposez les catégories pour modifier leur ordre d'affichage dans l'application.
            Les catégories inactives ne sont pas visibles par les utilisateurs.
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
