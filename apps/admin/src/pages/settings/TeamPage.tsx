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
import { formatDate } from '@/lib/utils'
import {
  Plus,
  MoreHorizontal,
  Edit,
  Trash2,
  User,
  Shield,
  ShieldCheck,
  Mail,
  Clock,
  CheckCircle,
  XCircle,
} from 'lucide-react'

interface AdminUser {
  id: string
  email: string
  firstName: string
  lastName: string
  role: 'super_admin' | 'moderator' | 'support'
  status: 'active' | 'inactive'
  lastLoginAt: string | null
  createdAt: string
  createdBy: string
}

// Mock data
const mockAdmins: AdminUser[] = [
  {
    id: '1',
    email: 'admin@yousoon.com',
    firstName: 'Admin',
    lastName: 'Yousoon',
    role: 'super_admin',
    status: 'active',
    lastLoginAt: '2024-12-09T10:30:00',
    createdAt: '2024-01-01',
    createdBy: 'Système',
  },
  {
    id: '2',
    email: 'moderation@yousoon.com',
    firstName: 'Marie',
    lastName: 'Modératrice',
    role: 'moderator',
    status: 'active',
    lastLoginAt: '2024-12-08T16:45:00',
    createdAt: '2024-06-15',
    createdBy: 'Admin Yousoon',
  },
  {
    id: '3',
    email: 'support@yousoon.com',
    firstName: 'Pierre',
    lastName: 'Support',
    role: 'support',
    status: 'active',
    lastLoginAt: '2024-12-09T08:15:00',
    createdAt: '2024-09-01',
    createdBy: 'Admin Yousoon',
  },
  {
    id: '4',
    email: 'ancien@yousoon.com',
    firstName: 'Jean',
    lastName: 'Ancien',
    role: 'moderator',
    status: 'inactive',
    lastLoginAt: '2024-10-15T12:00:00',
    createdAt: '2024-03-15',
    createdBy: 'Admin Yousoon',
  },
]

function getRoleBadge(role: AdminUser['role']) {
  const styles = {
    super_admin: { bg: 'bg-purple-100', text: 'text-purple-700', label: 'Super Admin', icon: ShieldCheck },
    moderator: { bg: 'bg-blue-100', text: 'text-blue-700', label: 'Modérateur', icon: Shield },
    support: { bg: 'bg-green-100', text: 'text-green-700', label: 'Support', icon: User },
  }
  const style = styles[role]
  const Icon = style.icon
  return (
    <span className={`inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium ${style.bg} ${style.text}`}>
      <Icon className="h-3 w-3" />
      {style.label}
    </span>
  )
}

function getPermissions(role: AdminUser['role']) {
  const permissions = {
    super_admin: ['Tout accès', 'Gestion équipe', 'Configuration', 'Analytics'],
    moderator: ['Validation partenaires', 'Validation offres', 'Modération avis', 'Vérification CNI'],
    support: ['Lecture seule', 'Gestion utilisateurs', 'Support tickets'],
  }
  return permissions[role]
}

export function TeamPage() {
  const [admins, setAdmins] = useState(mockAdmins)
  const [isCreating, setIsCreating] = useState(false)
  const [editingAdmin, setEditingAdmin] = useState<AdminUser | null>(null)

  const toggleStatus = (id: string) => {
    setAdmins(admins.map(a => 
      a.id === id ? { ...a, status: a.status === 'active' ? 'inactive' : 'active' } : a
    ))
  }

  const deleteAdmin = (id: string) => {
    if (confirm('Êtes-vous sûr de vouloir supprimer cet administrateur ?')) {
      setAdmins(admins.filter(a => a.id !== id))
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Équipe admin</h1>
          <p className="text-muted-foreground">Gestion des comptes administrateurs</p>
        </div>
        <Button className="gap-2" onClick={() => setIsCreating(true)}>
          <Plus className="h-4 w-4" />
          Nouvel admin
        </Button>
      </div>

      {/* Roles explanation */}
      <div className="grid gap-4 md:grid-cols-3">
        {(['super_admin', 'moderator', 'support'] as const).map((role) => (
          <Card key={role}>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm">{getRoleBadge(role)}</CardTitle>
            </CardHeader>
            <CardContent>
              <ul className="text-sm text-muted-foreground space-y-1">
                {getPermissions(role).map((perm, i) => (
                  <li key={i} className="flex items-center gap-2">
                    <CheckCircle className="h-3 w-3 text-green-500" />
                    {perm}
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Admins list */}
      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b">
                <tr>
                  <th className="text-left p-4 font-medium text-muted-foreground">Administrateur</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Rôle</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Statut</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Dernière connexion</th>
                  <th className="text-left p-4 font-medium text-muted-foreground">Créé le</th>
                  <th className="text-right p-4 font-medium text-muted-foreground">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {admins.map((admin) => (
                  <tr key={admin.id} className={`hover:bg-gray-50 ${admin.status === 'inactive' ? 'opacity-60' : ''}`}>
                    <td className="p-4">
                      <div className="flex items-center gap-3">
                        <div className="h-10 w-10 rounded-full bg-gray-100 flex items-center justify-center">
                          <User className="h-5 w-5 text-gray-500" />
                        </div>
                        <div>
                          <p className="font-medium">{admin.firstName} {admin.lastName}</p>
                          <p className="text-sm text-muted-foreground flex items-center gap-1">
                            <Mail className="h-3 w-3" />
                            {admin.email}
                          </p>
                        </div>
                      </div>
                    </td>
                    <td className="p-4">{getRoleBadge(admin.role)}</td>
                    <td className="p-4">
                      {admin.status === 'active' ? (
                        <span className="inline-flex items-center gap-1 text-green-600 text-sm">
                          <CheckCircle className="h-4 w-4" />
                          Actif
                        </span>
                      ) : (
                        <span className="inline-flex items-center gap-1 text-gray-500 text-sm">
                          <XCircle className="h-4 w-4" />
                          Inactif
                        </span>
                      )}
                    </td>
                    <td className="p-4">
                      {admin.lastLoginAt ? (
                        <span className="text-sm flex items-center gap-1">
                          <Clock className="h-3 w-3 text-muted-foreground" />
                          {formatDate(admin.lastLoginAt)}
                        </span>
                      ) : (
                        <span className="text-sm text-muted-foreground">Jamais</span>
                      )}
                    </td>
                    <td className="p-4 text-sm text-muted-foreground">
                      <p>{formatDate(admin.createdAt)}</p>
                      <p className="text-xs">par {admin.createdBy}</p>
                    </td>
                    <td className="p-4">
                      <div className="flex justify-end">
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon" disabled={admin.role === 'super_admin'}>
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem onClick={() => setEditingAdmin(admin)} className="gap-2">
                              <Edit className="h-4 w-4" />
                              Modifier
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => toggleStatus(admin.id)} className="gap-2">
                              {admin.status === 'active' ? 'Désactiver' : 'Activer'}
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => deleteAdmin(admin.id)} className="gap-2 text-red-600">
                              <Trash2 className="h-4 w-4" />
                              Supprimer
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      {/* Create/Edit Modal */}
      {(isCreating || editingAdmin) && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <Card className="w-full max-w-lg">
            <CardHeader>
              <CardTitle>{isCreating ? 'Nouvel administrateur' : 'Modifier l\'administrateur'}</CardTitle>
              <CardDescription>
                {isCreating 
                  ? 'Créez un nouveau compte administrateur. Un email d\'invitation sera envoyé.'
                  : 'Modifiez les informations de l\'administrateur'}
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium">Prénom</label>
                  <Input 
                    defaultValue={editingAdmin?.firstName || ''} 
                    placeholder="Jean"
                    className="mt-1" 
                  />
                </div>
                <div>
                  <label className="text-sm font-medium">Nom</label>
                  <Input 
                    defaultValue={editingAdmin?.lastName || ''} 
                    placeholder="Dupont"
                    className="mt-1" 
                  />
                </div>
              </div>
              <div>
                <label className="text-sm font-medium">Email</label>
                <Input 
                  type="email"
                  defaultValue={editingAdmin?.email || ''} 
                  placeholder="jean.dupont@yousoon.com"
                  className="mt-1" 
                  disabled={!!editingAdmin}
                />
              </div>
              <div>
                <label className="text-sm font-medium">Rôle</label>
                <select 
                  className="w-full mt-1 px-3 py-2 border rounded-md"
                  defaultValue={editingAdmin?.role || 'support'}
                >
                  <option value="super_admin">Super Admin</option>
                  <option value="moderator">Modérateur</option>
                  <option value="support">Support</option>
                </select>
              </div>
              <div className="flex gap-2 pt-4">
                <Button className="flex-1">
                  {isCreating ? 'Envoyer l\'invitation' : 'Enregistrer'}
                </Button>
                <Button 
                  variant="outline" 
                  onClick={() => {
                    setIsCreating(false)
                    setEditingAdmin(null)
                  }}
                >
                  Annuler
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* 2FA Info */}
      <Card className="bg-blue-50 border-blue-200">
        <CardContent className="pt-6">
          <h3 className="font-medium text-blue-800 flex items-center gap-2">
            <Shield className="h-5 w-5" />
            Authentification à deux facteurs (2FA)
          </h3>
          <p className="text-sm text-blue-600 mt-1">
            La 2FA est obligatoire pour tous les comptes administrateurs. 
            Chaque admin doit configurer son application d'authentification (Google Authenticator, Authy, etc.)
            lors de sa première connexion.
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
