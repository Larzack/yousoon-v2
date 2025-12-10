import { Routes, Route, Navigate } from 'react-router-dom'
import { Toaster } from '@/components/ui/toaster'
import { useAuthStore } from '@/stores/authStore'

// Layouts
import { AuthLayout } from '@/components/layout/AuthLayout'
import { DashboardLayout } from '@/components/layout/DashboardLayout'

// Auth pages
import { LoginPage } from '@/pages/auth/LoginPage'
import { RegisterPage } from '@/pages/auth/RegisterPage'
import { ForgotPasswordPage } from '@/pages/auth/ForgotPasswordPage'

// Dashboard pages
import { DashboardPage } from '@/pages/dashboard/DashboardPage'
import { OffersPage } from '@/pages/offers/OffersPage'
import { OfferDetailPage } from '@/pages/offers/OfferDetailPage'
import { CreateOfferPage } from '@/pages/offers/CreateOfferPage'
import { EstablishmentsPage } from '@/pages/establishments/EstablishmentsPage'
import { EstablishmentDetailPage } from '@/pages/establishments/EstablishmentDetailPage'
import { AnalyticsPage } from '@/pages/analytics/AnalyticsPage'
import { BookingsPage } from '@/pages/bookings/BookingsPage'
import { SettingsPage } from '@/pages/settings/SettingsPage'
import { TeamPage } from '@/pages/settings/TeamPage'
import { ProfilePage } from '@/pages/settings/ProfilePage'

// Protected Route wrapper
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }
  
  return <>{children}</>
}

// Public Route wrapper (redirect if already authenticated)
function PublicRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  
  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />
  }
  
  return <>{children}</>
}

export function App() {
  return (
    <>
      <Routes>
        {/* Public routes */}
        <Route element={<AuthLayout />}>
          <Route
            path="/login"
            element={
              <PublicRoute>
                <LoginPage />
              </PublicRoute>
            }
          />
          <Route
            path="/register"
            element={
              <PublicRoute>
                <RegisterPage />
              </PublicRoute>
            }
          />
          <Route
            path="/forgot-password"
            element={
              <PublicRoute>
                <ForgotPasswordPage />
              </PublicRoute>
            }
          />
        </Route>

        {/* Protected routes */}
        <Route
          element={
            <ProtectedRoute>
              <DashboardLayout />
            </ProtectedRoute>
          }
        >
          <Route path="/dashboard" element={<DashboardPage />} />
          
          {/* Offers */}
          <Route path="/offers" element={<OffersPage />} />
          <Route path="/offers/create" element={<CreateOfferPage />} />
          <Route path="/offers/:id" element={<OfferDetailPage />} />
          
          {/* Establishments */}
          <Route path="/establishments" element={<EstablishmentsPage />} />
          <Route path="/establishments/:id" element={<EstablishmentDetailPage />} />
          
          {/* Analytics */}
          <Route path="/analytics" element={<AnalyticsPage />} />
          
          {/* Bookings */}
          <Route path="/bookings" element={<BookingsPage />} />
          
          {/* Settings */}
          <Route path="/settings" element={<SettingsPage />} />
          <Route path="/settings/team" element={<TeamPage />} />
          <Route path="/settings/profile" element={<ProfilePage />} />
        </Route>

        {/* Default redirect */}
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
      
      <Toaster />
    </>
  )
}
