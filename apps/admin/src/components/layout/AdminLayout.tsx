import { useState } from 'react'
import { Outlet, Link, useLocation, useNavigate } from 'react-router-dom'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/stores/authStore'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  LayoutDashboard,
  Users,
  Building2,
  Tag,
  ShieldCheck,
  MessageSquare,
  CreditCard,
  BarChart3,
  Settings,
  ChevronLeft,
  ChevronRight,
  LogOut,
  Bell,
  FolderTree,
  Cog,
  UserCog,
  Clock,
  AlertTriangle,
} from 'lucide-react'

interface NavItem {
  label: string
  icon: React.ElementType
  href: string
  badge?: number
  children?: { label: string; href: string; badge?: number }[]
}

const navigation: NavItem[] = [
  { label: 'Dashboard', icon: LayoutDashboard, href: '/' },
  { label: 'Utilisateurs', icon: Users, href: '/users' },
  {
    label: 'Partenaires',
    icon: Building2,
    href: '/partners',
    children: [
      { label: 'Tous', href: '/partners' },
      { label: 'En attente', href: '/partners/pending', badge: 5 },
    ],
  },
  { label: 'Offres', icon: Tag, href: '/offers' },
  { label: 'Vérifications CNI', icon: ShieldCheck, href: '/identity', badge: 12 },
  {
    label: 'Avis',
    icon: MessageSquare,
    href: '/reviews',
    children: [
      { label: 'Tous', href: '/reviews' },
      { label: 'Signalés', href: '/reviews/reported', badge: 3 },
    ],
  },
  {
    label: 'Abonnements',
    icon: CreditCard,
    href: '/subscriptions',
    children: [
      { label: 'Actifs', href: '/subscriptions' },
      { label: 'Plans', href: '/subscriptions/plans' },
    ],
  },
  { label: 'Analytics', icon: BarChart3, href: '/analytics' },
]

const settingsNavigation: NavItem[] = [
  { label: 'Catégories', icon: FolderTree, href: '/settings/categories' },
  { label: 'Configuration', icon: Cog, href: '/settings/config' },
  { label: 'Équipe Admin', icon: UserCog, href: '/settings/team' },
]

export function AdminLayout() {
  const [isCollapsed, setIsCollapsed] = useState(false)
  const location = useLocation()
  const navigate = useNavigate()
  const { user, logout } = useAuthStore()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const isActive = (href: string) => {
    if (href === '/') return location.pathname === '/'
    return location.pathname.startsWith(href)
  }

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Sidebar */}
      <aside
        className={cn(
          'flex flex-col bg-white border-r transition-all duration-300',
          isCollapsed ? 'w-16' : 'w-64'
        )}
      >
        {/* Logo */}
        <div className="flex items-center h-16 px-4 border-b">
          {!isCollapsed && (
            <span className="text-xl font-bold text-indigo-600">
              Yousoon Admin
            </span>
          )}
          {isCollapsed && (
            <span className="text-xl font-bold text-indigo-600 mx-auto">Y</span>
          )}
        </div>

        {/* Navigation */}
        <nav className="flex-1 overflow-y-auto py-4">
          <div className="px-3 space-y-1">
            {navigation.map((item) => (
              <div key={item.href}>
                <Link
                  to={item.children ? item.children[0].href : item.href}
                  className={cn(
                    'flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors',
                    isActive(item.href)
                      ? 'bg-indigo-50 text-indigo-600'
                      : 'text-gray-700 hover:bg-gray-100'
                  )}
                >
                  <item.icon className="h-5 w-5 flex-shrink-0" />
                  {!isCollapsed && (
                    <>
                      <span className="flex-1">{item.label}</span>
                      {item.badge && (
                        <span className="bg-red-100 text-red-600 text-xs font-semibold px-2 py-0.5 rounded-full">
                          {item.badge}
                        </span>
                      )}
                    </>
                  )}
                </Link>
                {!isCollapsed && item.children && isActive(item.href) && (
                  <div className="ml-8 mt-1 space-y-1">
                    {item.children.map((child) => (
                      <Link
                        key={child.href}
                        to={child.href}
                        className={cn(
                          'flex items-center gap-2 px-3 py-1.5 rounded text-sm transition-colors',
                          location.pathname === child.href
                            ? 'text-indigo-600 font-medium'
                            : 'text-gray-600 hover:text-gray-900'
                        )}
                      >
                        <span className="flex-1">{child.label}</span>
                        {child.badge && (
                          <span className="bg-red-100 text-red-600 text-xs font-semibold px-2 py-0.5 rounded-full">
                            {child.badge}
                          </span>
                        )}
                      </Link>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>

          {/* Settings */}
          <div className="px-3 mt-6">
            <div className={cn('text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2', isCollapsed && 'hidden')}>
              Paramètres
            </div>
            <div className="space-y-1">
              {settingsNavigation.map((item) => (
                <Link
                  key={item.href}
                  to={item.href}
                  className={cn(
                    'flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors',
                    isActive(item.href)
                      ? 'bg-indigo-50 text-indigo-600'
                      : 'text-gray-700 hover:bg-gray-100'
                  )}
                >
                  <item.icon className="h-5 w-5 flex-shrink-0" />
                  {!isCollapsed && <span>{item.label}</span>}
                </Link>
              ))}
            </div>
          </div>
        </nav>

        {/* Collapse button */}
        <div className="p-3 border-t">
          <Button
            variant="ghost"
            size="sm"
            className="w-full justify-center"
            onClick={() => setIsCollapsed(!isCollapsed)}
          >
            {isCollapsed ? (
              <ChevronRight className="h-4 w-4" />
            ) : (
              <ChevronLeft className="h-4 w-4" />
            )}
          </Button>
        </div>
      </aside>

      {/* Main content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <header className="flex items-center justify-between h-16 px-6 bg-white border-b">
          <div>
            {/* Breadcrumb or page title could go here */}
          </div>

          <div className="flex items-center gap-4">
            {/* Pending actions */}
            <div className="flex items-center gap-2 text-sm">
              <div className="flex items-center gap-1 text-yellow-600 bg-yellow-50 px-3 py-1 rounded-full">
                <Clock className="h-4 w-4" />
                <span>12 CNI</span>
              </div>
              <div className="flex items-center gap-1 text-red-600 bg-red-50 px-3 py-1 rounded-full">
                <AlertTriangle className="h-4 w-4" />
                <span>3 signalements</span>
              </div>
            </div>

            {/* Notifications */}
            <Button variant="ghost" size="icon" className="relative">
              <Bell className="h-5 w-5" />
              <span className="absolute -top-1 -right-1 h-4 w-4 bg-red-500 text-white text-xs rounded-full flex items-center justify-center">
                5
              </span>
            </Button>

            {/* User menu */}
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="flex items-center gap-2">
                  <Avatar className="h-8 w-8">
                    <AvatarImage src={user?.avatar} />
                    <AvatarFallback>
                      {user?.firstName?.charAt(0)}{user?.lastName?.charAt(0)}
                    </AvatarFallback>
                  </Avatar>
                  <span className="hidden md:inline-block text-sm font-medium">
                    {user?.firstName} {user?.lastName}
                  </span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-56">
                <DropdownMenuLabel>
                  <div className="flex flex-col">
                    <span>{user?.firstName} {user?.lastName}</span>
                    <span className="text-xs text-muted-foreground">{user?.email}</span>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem>
                  <Settings className="mr-2 h-4 w-4" />
                  Paramètres
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={handleLogout} className="text-red-600">
                  <LogOut className="mr-2 h-4 w-4" />
                  Déconnexion
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </header>

        {/* Page content */}
        <main className="flex-1 overflow-y-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
