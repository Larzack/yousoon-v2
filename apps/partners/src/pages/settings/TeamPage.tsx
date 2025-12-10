import { useState } from 'react';
import {
  Plus,
  Search,
  MoreHorizontal,
  Mail,
  Shield,
  UserX,
  Trash2,
  CheckCircle,
  Clock,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useToast } from '@/hooks/use-toast';
import { cn } from '@/lib/utils';

// Types
type TeamRole = 'admin' | 'manager' | 'staff' | 'viewer';
type TeamStatus = 'active' | 'pending';

interface TeamMember {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  role: TeamRole;
  status: TeamStatus;
  avatar?: string;
  invitedAt?: string;
  joinedAt?: string;
}

// Mock data
const mockTeamMembers: TeamMember[] = [
  {
    id: '1',
    firstName: 'Jean',
    lastName: 'Dupont',
    email: 'jean.dupont@comptoir-parisien.fr',
    role: 'admin',
    status: 'active',
    avatar: 'https://randomuser.me/api/portraits/men/1.jpg',
    joinedAt: '2023-01-15',
  },
  {
    id: '2',
    firstName: 'Marie',
    lastName: 'Martin',
    email: 'marie.martin@comptoir-parisien.fr',
    role: 'manager',
    status: 'active',
    avatar: 'https://randomuser.me/api/portraits/women/2.jpg',
    joinedAt: '2023-03-20',
  },
  {
    id: '3',
    firstName: 'Pierre',
    lastName: 'Bernard',
    email: 'pierre.bernard@comptoir-parisien.fr',
    role: 'staff',
    status: 'active',
    joinedAt: '2023-06-01',
  },
  {
    id: '4',
    firstName: '',
    lastName: '',
    email: 'nouveau.membre@email.com',
    role: 'viewer',
    status: 'pending',
    invitedAt: '2024-01-10',
  },
];

const roleConfig: Record<TeamRole, { label: string; description: string; color: string }> = {
  admin: {
    label: 'Administrateur',
    description: 'Accès complet à toutes les fonctionnalités',
    color: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300',
  },
  manager: {
    label: 'Manager',
    description: 'Gestion des offres, établissements et statistiques',
    color: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300',
  },
  staff: {
    label: 'Staff',
    description: 'Gestion des réservations et check-ins',
    color: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
  },
  viewer: {
    label: 'Lecteur',
    description: 'Consultation uniquement',
    color: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300',
  },
};

export function TeamPage() {
  const { toast } = useToast();
  const [searchQuery, setSearchQuery] = useState('');
  const [showInviteModal, setShowInviteModal] = useState(false);
  const [inviteEmail, setInviteEmail] = useState('');
  const [inviteRole, setInviteRole] = useState<TeamRole>('viewer');

  const filteredMembers = mockTeamMembers.filter(
    (member) =>
      member.firstName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      member.lastName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      member.email.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const activeMembers = mockTeamMembers.filter((m) => m.status === 'active').length;
  const pendingInvites = mockTeamMembers.filter((m) => m.status === 'pending').length;

  const handleInvite = () => {
    if (!inviteEmail) return;
    
    toast({
      title: 'Invitation envoyée',
      description: `Une invitation a été envoyée à ${inviteEmail}`,
    });
    setShowInviteModal(false);
    setInviteEmail('');
    setInviteRole('viewer');
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Équipe</h1>
          <p className="text-muted-foreground">
            Gérez les membres de votre équipe et leurs permissions
          </p>
        </div>
        <Button 
          onClick={() => setShowInviteModal(true)}
          className="bg-yousoon-gold hover:bg-yousoon-gold/90 text-black"
        >
          <Plus className="mr-2 h-4 w-4" />
          Inviter un membre
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Membres actifs</CardTitle>
            <CheckCircle className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{activeMembers}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Invitations en attente</CardTitle>
            <Clock className="h-4 w-4 text-yellow-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{pendingInvites}</div>
          </CardContent>
        </Card>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          placeholder="Rechercher un membre..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-10"
        />
      </div>

      {/* Team list */}
      <Card>
        <CardContent className="p-0">
          <div className="divide-y">
            {filteredMembers.map((member) => (
              <div
                key={member.id}
                className="flex items-center justify-between p-4 hover:bg-muted/50"
              >
                <div className="flex items-center gap-4">
                  <Avatar className="h-10 w-10">
                    <AvatarImage src={member.avatar} />
                    <AvatarFallback>
                      {member.status === 'pending' 
                        ? '?' 
                        : `${member.firstName[0]}${member.lastName[0]}`}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <div className="flex items-center gap-2">
                      <p className="font-medium">
                        {member.status === 'pending'
                          ? member.email
                          : `${member.firstName} ${member.lastName}`}
                      </p>
                      {member.status === 'pending' && (
                        <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300">
                          <Clock className="mr-1 h-3 w-3" />
                          En attente
                        </span>
                      )}
                    </div>
                    {member.status === 'active' && (
                      <p className="text-sm text-muted-foreground">{member.email}</p>
                    )}
                  </div>
                </div>
                <div className="flex items-center gap-4">
                  <span
                    className={cn(
                      'inline-flex items-center px-2 py-1 rounded-full text-xs font-medium',
                      roleConfig[member.role].color
                    )}
                  >
                    <Shield className="mr-1 h-3 w-3" />
                    {roleConfig[member.role].label}
                  </span>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="icon" className="h-8 w-8">
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      {member.status === 'pending' ? (
                        <>
                          <DropdownMenuItem>
                            <Mail className="mr-2 h-4 w-4" />
                            Renvoyer l'invitation
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem className="text-destructive">
                            <Trash2 className="mr-2 h-4 w-4" />
                            Annuler l'invitation
                          </DropdownMenuItem>
                        </>
                      ) : (
                        <>
                          <DropdownMenuItem>
                            <Shield className="mr-2 h-4 w-4" />
                            Modifier le rôle
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem className="text-destructive">
                            <UserX className="mr-2 h-4 w-4" />
                            Retirer de l'équipe
                          </DropdownMenuItem>
                        </>
                      )}
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Roles explanation */}
      <Card>
        <CardHeader>
          <CardTitle>Rôles et permissions</CardTitle>
          <CardDescription>
            Comprendre les différents niveaux d'accès
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 sm:grid-cols-2">
            {Object.entries(roleConfig).map(([role, config]) => (
              <div key={role} className="p-4 border rounded-lg">
                <div className="flex items-center gap-2 mb-2">
                  <span className={cn('px-2 py-1 rounded-full text-xs font-medium', config.color)}>
                    {config.label}
                  </span>
                </div>
                <p className="text-sm text-muted-foreground">{config.description}</p>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Invite Modal */}
      {showInviteModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
          <Card className="w-full max-w-md">
            <CardHeader>
              <CardTitle>Inviter un membre</CardTitle>
              <CardDescription>
                Envoyez une invitation par email pour rejoindre votre équipe
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="inviteEmail">Adresse email</Label>
                <Input
                  id="inviteEmail"
                  type="email"
                  placeholder="email@example.com"
                  value={inviteEmail}
                  onChange={(e) => setInviteEmail(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label>Rôle</Label>
                <div className="grid grid-cols-2 gap-2">
                  {(Object.keys(roleConfig) as TeamRole[]).map((role) => (
                    <button
                      key={role}
                      type="button"
                      onClick={() => setInviteRole(role)}
                      className={cn(
                        'p-3 border rounded-lg text-left transition-colors',
                        inviteRole === role
                          ? 'border-yousoon-gold bg-yousoon-gold/10'
                          : 'hover:border-muted-foreground'
                      )}
                    >
                      <p className="font-medium text-sm">{roleConfig[role].label}</p>
                      <p className="text-xs text-muted-foreground">{roleConfig[role].description}</p>
                    </button>
                  ))}
                </div>
              </div>
              <div className="flex gap-2 pt-4">
                <Button
                  onClick={handleInvite}
                  className="flex-1 bg-yousoon-gold hover:bg-yousoon-gold/90 text-black"
                >
                  <Mail className="mr-2 h-4 w-4" />
                  Envoyer l'invitation
                </Button>
                <Button variant="outline" onClick={() => setShowInviteModal(false)}>
                  Annuler
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
}

export default TeamPage;
