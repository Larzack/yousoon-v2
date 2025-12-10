import { useState } from 'react';
import {
  Search,
  Filter,
  Calendar,
  CheckCircle2,
  XCircle,
  Clock,
  MoreHorizontal,
  Eye,
  MapPin,
  Phone,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { cn, formatDate } from '@/lib/utils';

// Types
type BookingStatus = 'pending' | 'confirmed' | 'checked_in' | 'cancelled' | 'no_show';

interface Booking {
  id: string;
  user: {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    phone: string;
    avatar?: string;
  };
  offer: {
    id: string;
    title: string;
    establishment: string;
  };
  status: BookingStatus;
  bookedAt: string;
  checkedInAt?: string;
  qrCode: string;
}

// Mock data
const mockBookings: Booking[] = [
  {
    id: '1',
    user: {
      id: 'u1',
      firstName: 'Marie',
      lastName: 'Dupont',
      email: 'marie.dupont@email.com',
      phone: '+33 6 12 34 56 78',
      avatar: 'https://randomuser.me/api/portraits/women/1.jpg',
    },
    offer: {
      id: 'o1',
      title: 'Happy Hour -50%',
      establishment: 'Le Comptoir Parisien',
    },
    status: 'confirmed',
    bookedAt: '2024-01-15T14:30:00Z',
    qrCode: 'QR-12345',
  },
  {
    id: '2',
    user: {
      id: 'u2',
      firstName: 'Pierre',
      lastName: 'Martin',
      email: 'pierre.martin@email.com',
      phone: '+33 6 98 76 54 32',
      avatar: 'https://randomuser.me/api/portraits/men/2.jpg',
    },
    offer: {
      id: 'o2',
      title: 'Brunch du dimanche',
      establishment: 'La Terrasse du Marais',
    },
    status: 'checked_in',
    bookedAt: '2024-01-15T10:00:00Z',
    checkedInAt: '2024-01-15T11:30:00Z',
    qrCode: 'QR-12346',
  },
  {
    id: '3',
    user: {
      id: 'u3',
      firstName: 'Sophie',
      lastName: 'Bernard',
      email: 'sophie.bernard@email.com',
      phone: '+33 6 11 22 33 44',
    },
    offer: {
      id: 'o3',
      title: 'Menu découverte -30%',
      establishment: 'Le Comptoir Parisien',
    },
    status: 'pending',
    bookedAt: '2024-01-15T16:00:00Z',
    qrCode: 'QR-12347',
  },
  {
    id: '4',
    user: {
      id: 'u4',
      firstName: 'Jean',
      lastName: 'Petit',
      email: 'jean.petit@email.com',
      phone: '+33 6 55 66 77 88',
      avatar: 'https://randomuser.me/api/portraits/men/4.jpg',
    },
    offer: {
      id: 'o1',
      title: 'Happy Hour -50%',
      establishment: 'Le Comptoir Parisien',
    },
    status: 'cancelled',
    bookedAt: '2024-01-14T18:00:00Z',
    qrCode: 'QR-12348',
  },
  {
    id: '5',
    user: {
      id: 'u5',
      firstName: 'Léa',
      lastName: 'Moreau',
      email: 'lea.moreau@email.com',
      phone: '+33 6 99 88 77 66',
      avatar: 'https://randomuser.me/api/portraits/women/5.jpg',
    },
    offer: {
      id: 'o4',
      title: 'Afterwork cocktails',
      establishment: 'Le Rooftop Montmartre',
    },
    status: 'no_show',
    bookedAt: '2024-01-14T19:00:00Z',
    qrCode: 'QR-12349',
  },
];

const statusConfig: Record<BookingStatus, { label: string; icon: React.ReactNode; className: string }> = {
  pending: {
    label: 'En attente',
    icon: <Clock className="h-4 w-4" />,
    className: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300',
  },
  confirmed: {
    label: 'Confirmée',
    icon: <CheckCircle2 className="h-4 w-4" />,
    className: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300',
  },
  checked_in: {
    label: 'Check-in',
    icon: <CheckCircle2 className="h-4 w-4" />,
    className: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
  },
  cancelled: {
    label: 'Annulée',
    icon: <XCircle className="h-4 w-4" />,
    className: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300',
  },
  no_show: {
    label: 'No-show',
    icon: <XCircle className="h-4 w-4" />,
    className: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300',
  },
};

export function BookingsPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [filterStatus, setFilterStatus] = useState<BookingStatus | 'all'>('all');
  const [selectedBooking, setSelectedBooking] = useState<Booking | null>(null);

  const filteredBookings = mockBookings.filter((booking) => {
    const matchesSearch = 
      booking.user.firstName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      booking.user.lastName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      booking.user.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
      booking.offer.title.toLowerCase().includes(searchQuery.toLowerCase());

    const matchesStatus = filterStatus === 'all' || booking.status === filterStatus;

    return matchesSearch && matchesStatus;
  });

  // Stats
  const todayBookings = mockBookings.filter(b => b.status === 'pending' || b.status === 'confirmed').length;
  const todayCheckins = mockBookings.filter(b => b.status === 'checked_in').length;
  const todayNoShows = mockBookings.filter(b => b.status === 'no_show').length;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Réservations</h1>
          <p className="text-muted-foreground">
            Gérez les réservations et validez les check-ins
          </p>
        </div>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Réservations du jour</CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{todayBookings}</div>
            <p className="text-xs text-muted-foreground">en attente de check-in</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Check-ins</CardTitle>
            <CheckCircle2 className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-500">{todayCheckins}</div>
            <p className="text-xs text-muted-foreground">aujourd'hui</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">No-shows</CardTitle>
            <XCircle className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-500">{todayNoShows}</div>
            <p className="text-xs text-muted-foreground">aujourd'hui</p>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder="Rechercher par nom, email ou offre..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10"
          />
        </div>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline">
              <Filter className="mr-2 h-4 w-4" />
              {filterStatus === 'all' ? 'Tous les statuts' : statusConfig[filterStatus].label}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => setFilterStatus('all')}>
              Tous les statuts
            </DropdownMenuItem>
            {Object.entries(statusConfig).map(([status, config]) => (
              <DropdownMenuItem
                key={status}
                onClick={() => setFilterStatus(status as BookingStatus)}
              >
                <span className="flex items-center gap-2">
                  {config.icon}
                  {config.label}
                </span>
              </DropdownMenuItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {/* Bookings Table */}
      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b bg-muted/50">
                  <th className="text-left p-4 font-medium text-muted-foreground">Client</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Offre</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Date</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Statut</th>
                  <th className="text-right p-4 font-medium text-muted-foreground">Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredBookings.map((booking) => (
                  <tr key={booking.id} className="border-b hover:bg-muted/50">
                    <td className="p-4">
                      <div className="flex items-center gap-3">
                        <Avatar className="h-10 w-10">
                          <AvatarImage src={booking.user.avatar} />
                          <AvatarFallback>
                            {booking.user.firstName[0]}{booking.user.lastName[0]}
                          </AvatarFallback>
                        </Avatar>
                        <div>
                          <p className="font-medium">
                            {booking.user.firstName} {booking.user.lastName}
                          </p>
                          <p className="text-sm text-muted-foreground">
                            {booking.user.email}
                          </p>
                        </div>
                      </div>
                    </td>
                    <td className="p-4">
                      <div>
                        <p className="font-medium">{booking.offer.title}</p>
                        <p className="text-sm text-muted-foreground flex items-center gap-1">
                          <MapPin className="h-3 w-3" />
                          {booking.offer.establishment}
                        </p>
                      </div>
                    </td>
                    <td className="p-4">
                      <p className="text-sm">{formatDate(booking.bookedAt)}</p>
                      {booking.checkedInAt && (
                        <p className="text-xs text-muted-foreground">
                          Check-in: {formatDate(booking.checkedInAt)}
                        </p>
                      )}
                    </td>
                    <td className="p-4">
                      <span
                        className={cn(
                          'inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium',
                          statusConfig[booking.status].className
                        )}
                      >
                        {statusConfig[booking.status].icon}
                        {statusConfig[booking.status].label}
                      </span>
                    </td>
                    <td className="p-4 text-right">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon" className="h-8 w-8">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => setSelectedBooking(booking)}>
                            <Eye className="mr-2 h-4 w-4" />
                            Voir détails
                          </DropdownMenuItem>
                          {booking.status === 'confirmed' && (
                            <DropdownMenuItem className="text-green-600">
                              <CheckCircle2 className="mr-2 h-4 w-4" />
                              Valider check-in
                            </DropdownMenuItem>
                          )}
                          {(booking.status === 'pending' || booking.status === 'confirmed') && (
                            <DropdownMenuItem className="text-destructive">
                              <XCircle className="mr-2 h-4 w-4" />
                              Annuler
                            </DropdownMenuItem>
                          )}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          <div className="flex items-center justify-between p-4 border-t">
            <p className="text-sm text-muted-foreground">
              Affichage de 1 à {filteredBookings.length} sur {mockBookings.length} réservations
            </p>
            <div className="flex gap-1">
              <Button variant="outline" size="icon" disabled>
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <Button variant="outline" size="icon" className="bg-yousoon-gold text-black">
                1
              </Button>
              <Button variant="outline" size="icon">
                2
              </Button>
              <Button variant="outline" size="icon">
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Booking Detail Modal (simplified) */}
      {selectedBooking && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
          <Card className="w-full max-w-lg">
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>Détails de la réservation</CardTitle>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setSelectedBooking(null)}
              >
                <XCircle className="h-4 w-4" />
              </Button>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* User info */}
              <div className="flex items-center gap-4">
                <Avatar className="h-16 w-16">
                  <AvatarImage src={selectedBooking.user.avatar} />
                  <AvatarFallback className="text-lg">
                    {selectedBooking.user.firstName[0]}{selectedBooking.user.lastName[0]}
                  </AvatarFallback>
                </Avatar>
                <div>
                  <h3 className="text-lg font-semibold">
                    {selectedBooking.user.firstName} {selectedBooking.user.lastName}
                  </h3>
                  <p className="text-sm text-muted-foreground">{selectedBooking.user.email}</p>
                  <p className="text-sm text-muted-foreground flex items-center gap-1">
                    <Phone className="h-3 w-3" />
                    {selectedBooking.user.phone}
                  </p>
                </div>
              </div>

              {/* Offer info */}
              <div className="p-4 bg-muted rounded-lg">
                <h4 className="font-medium">{selectedBooking.offer.title}</h4>
                <p className="text-sm text-muted-foreground flex items-center gap-1">
                  <MapPin className="h-3 w-3" />
                  {selectedBooking.offer.establishment}
                </p>
              </div>

              {/* QR Code placeholder */}
              <div className="flex justify-center p-4 bg-white rounded-lg">
                <div className="text-center">
                  <div className="w-32 h-32 bg-gray-200 rounded-lg flex items-center justify-center mx-auto">
                    <span className="text-xs text-gray-500">QR Code</span>
                  </div>
                  <p className="text-sm font-mono mt-2">{selectedBooking.qrCode}</p>
                </div>
              </div>

              {/* Status and dates */}
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="text-muted-foreground">Statut</p>
                  <span
                    className={cn(
                      'inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium mt-1',
                      statusConfig[selectedBooking.status].className
                    )}
                  >
                    {statusConfig[selectedBooking.status].icon}
                    {statusConfig[selectedBooking.status].label}
                  </span>
                </div>
                <div>
                  <p className="text-muted-foreground">Réservé le</p>
                  <p className="font-medium">{formatDate(selectedBooking.bookedAt)}</p>
                </div>
              </div>

              {/* Actions */}
              <div className="flex gap-2 pt-4">
                {selectedBooking.status === 'confirmed' && (
                  <Button className="flex-1 bg-green-600 hover:bg-green-700">
                    <CheckCircle2 className="mr-2 h-4 w-4" />
                    Valider check-in
                  </Button>
                )}
                {(selectedBooking.status === 'pending' || selectedBooking.status === 'confirmed') && (
                  <Button variant="destructive" className="flex-1">
                    <XCircle className="mr-2 h-4 w-4" />
                    Annuler
                  </Button>
                )}
                <Button variant="outline" onClick={() => setSelectedBooking(null)}>
                  Fermer
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
}

export default BookingsPage;
