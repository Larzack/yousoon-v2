import { useState } from 'react';
import { Link } from 'react-router-dom';
import { 
  Plus, 
  Search, 
  MapPin, 
  Phone, 
  Clock, 
  MoreHorizontal,
  Edit,
  Trash2,
  Eye,
  Building2
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { cn } from '@/lib/utils';

// Types
interface Establishment {
  id: string;
  name: string;
  address: string;
  city: string;
  postalCode: string;
  phone: string;
  email: string;
  image: string;
  openingHours: string;
  offersCount: number;
  isActive: boolean;
}

// Mock data
const mockEstablishments: Establishment[] = [
  {
    id: '1',
    name: 'Le Comptoir Parisien',
    address: '12 Rue de la Paix',
    city: 'Paris',
    postalCode: '75001',
    phone: '+33 1 23 45 67 89',
    email: 'contact@comptoir-parisien.fr',
    image: 'https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?w=300',
    openingHours: '12h00 - 00h00',
    offersCount: 5,
    isActive: true,
  },
  {
    id: '2',
    name: 'La Terrasse du Marais',
    address: '45 Rue des Francs-Bourgeois',
    city: 'Paris',
    postalCode: '75004',
    phone: '+33 1 98 76 54 32',
    email: 'contact@terrasse-marais.fr',
    image: 'https://images.unsplash.com/photo-1552566626-52f8b828add9?w=300',
    openingHours: '10h00 - 02h00',
    offersCount: 3,
    isActive: true,
  },
  {
    id: '3',
    name: 'Le Rooftop Montmartre',
    address: '78 Rue Lepic',
    city: 'Paris',
    postalCode: '75018',
    phone: '+33 1 11 22 33 44',
    email: 'contact@rooftop-montmartre.fr',
    image: 'https://images.unsplash.com/photo-1514933651103-005eec06c04b?w=300',
    openingHours: '18h00 - 02h00',
    offersCount: 2,
    isActive: false,
  },
];

export function EstablishmentsPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [filterStatus, setFilterStatus] = useState<'all' | 'active' | 'inactive'>('all');

  const filteredEstablishments = mockEstablishments.filter((establishment) => {
    const matchesSearch = establishment.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      establishment.address.toLowerCase().includes(searchQuery.toLowerCase()) ||
      establishment.city.toLowerCase().includes(searchQuery.toLowerCase());
    
    const matchesStatus = filterStatus === 'all' ||
      (filterStatus === 'active' && establishment.isActive) ||
      (filterStatus === 'inactive' && !establishment.isActive);

    return matchesSearch && matchesStatus;
  });

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Établissements</h1>
          <p className="text-muted-foreground">
            Gérez vos établissements et leurs informations
          </p>
        </div>
        <Button asChild className="bg-yousoon-gold hover:bg-yousoon-gold/90 text-black">
          <Link to="/establishments/create">
            <Plus className="mr-2 h-4 w-4" />
            Nouvel établissement
          </Link>
        </Button>
      </div>

      {/* Filters */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder="Rechercher un établissement..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10"
          />
        </div>
        <div className="flex gap-2">
          <Button
            variant={filterStatus === 'all' ? 'default' : 'outline'}
            size="sm"
            onClick={() => setFilterStatus('all')}
            className={cn(
              filterStatus === 'all' && 'bg-yousoon-gold text-black hover:bg-yousoon-gold/90'
            )}
          >
            Tous ({mockEstablishments.length})
          </Button>
          <Button
            variant={filterStatus === 'active' ? 'default' : 'outline'}
            size="sm"
            onClick={() => setFilterStatus('active')}
            className={cn(
              filterStatus === 'active' && 'bg-yousoon-gold text-black hover:bg-yousoon-gold/90'
            )}
          >
            Actifs ({mockEstablishments.filter(e => e.isActive).length})
          </Button>
          <Button
            variant={filterStatus === 'inactive' ? 'default' : 'outline'}
            size="sm"
            onClick={() => setFilterStatus('inactive')}
            className={cn(
              filterStatus === 'inactive' && 'bg-yousoon-gold text-black hover:bg-yousoon-gold/90'
            )}
          >
            Inactifs ({mockEstablishments.filter(e => !e.isActive).length})
          </Button>
        </div>
      </div>

      {/* Establishments Grid */}
      {filteredEstablishments.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <Building2 className="h-12 w-12 text-muted-foreground mb-4" />
            <h3 className="text-lg font-semibold mb-2">Aucun établissement trouvé</h3>
            <p className="text-muted-foreground text-center mb-4">
              {searchQuery
                ? 'Modifiez votre recherche ou vos filtres'
                : 'Commencez par ajouter votre premier établissement'}
            </p>
            {!searchQuery && (
              <Button asChild className="bg-yousoon-gold hover:bg-yousoon-gold/90 text-black">
                <Link to="/establishments/create">
                  <Plus className="mr-2 h-4 w-4" />
                  Ajouter un établissement
                </Link>
              </Button>
            )}
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {filteredEstablishments.map((establishment) => (
            <EstablishmentCard key={establishment.id} establishment={establishment} />
          ))}
        </div>
      )}
    </div>
  );
}

interface EstablishmentCardProps {
  establishment: Establishment;
}

function EstablishmentCard({ establishment }: EstablishmentCardProps) {
  return (
    <Card className="overflow-hidden hover:shadow-lg transition-shadow">
      {/* Image */}
      <div className="relative h-40 overflow-hidden">
        <img
          src={establishment.image}
          alt={establishment.name}
          className="w-full h-full object-cover"
        />
        <div className="absolute top-2 right-2">
          <span
            className={cn(
              'inline-flex items-center px-2 py-1 rounded-full text-xs font-medium',
              establishment.isActive
                ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
                : 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300'
            )}
          >
            {establishment.isActive ? 'Actif' : 'Inactif'}
          </span>
        </div>
      </div>

      <CardHeader className="pb-2">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-3">
            <Avatar className="h-10 w-10">
              <AvatarImage src={establishment.image} alt={establishment.name} />
              <AvatarFallback>{establishment.name.substring(0, 2).toUpperCase()}</AvatarFallback>
            </Avatar>
            <div>
              <CardTitle className="text-lg">{establishment.name}</CardTitle>
              <p className="text-sm text-muted-foreground">
                {establishment.offersCount} offre{establishment.offersCount > 1 ? 's' : ''} active{establishment.offersCount > 1 ? 's' : ''}
              </p>
            </div>
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem asChild>
                <Link to={`/establishments/${establishment.id}`}>
                  <Eye className="mr-2 h-4 w-4" />
                  Voir
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link to={`/establishments/${establishment.id}/edit`}>
                  <Edit className="mr-2 h-4 w-4" />
                  Modifier
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem className="text-destructive focus:text-destructive">
                <Trash2 className="mr-2 h-4 w-4" />
                Supprimer
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </CardHeader>

      <CardContent className="space-y-3">
        <div className="flex items-start gap-2 text-sm">
          <MapPin className="h-4 w-4 text-muted-foreground shrink-0 mt-0.5" />
          <span className="text-muted-foreground">
            {establishment.address}, {establishment.postalCode} {establishment.city}
          </span>
        </div>
        <div className="flex items-center gap-2 text-sm">
          <Phone className="h-4 w-4 text-muted-foreground" />
          <span className="text-muted-foreground">{establishment.phone}</span>
        </div>
        <div className="flex items-center gap-2 text-sm">
          <Clock className="h-4 w-4 text-muted-foreground" />
          <span className="text-muted-foreground">{establishment.openingHours}</span>
        </div>
      </CardContent>
    </Card>
  );
}

export default EstablishmentsPage;
