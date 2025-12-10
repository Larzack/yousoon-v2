import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'
import { Toaster } from '@/components/ui/toaster'

// Layouts
import { AdminLayout } from '@/components/layout/AdminLayout'

// Pages
import { LoginPage } from '@/pages/auth/LoginPage'
import { DashboardPage } from '@/pages/dashboard/DashboardPage'
import { UsersPage } from '@/pages/users/UsersPage'
import { UserDetailPage } from '@/pages/users/UserDetailPage'
import { PartnersPage, PartnerDetailPage, PartnersPendingPage } from '@/pages/partners'
import { OffersPage, OfferDetailPage, OffersPendingPage } from '@/pages/offers'
import { IdentityVerificationsPage, IdentityDetailPage } from '@/pages/identity'
import { ReviewsPage, ReviewsReportedPage } from '@/pages/reviews'
import { SubscriptionsPage, PlansPage } from '@/pages/subscriptions'
import { AnalyticsPage } from '@/pages/analytics'
import { CategoriesPage, ConfigPage, TeamPage } from '@/pages/settings'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuthStore()
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }
  
  return <>{children}</>
}

export default function App() {
  return (
    <>
      <Routes>
        {/* Auth */}
        <Route path="/login" element={<LoginPage />} />
        
        {/* Protected Routes */}
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <AdminLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<DashboardPage />} />
          
          {/* Users */}
          <Route path="users" element={<UsersPage />} />
          <Route path="users/:id" element={<UserDetailPage />} />
          
          {/* Partners */}
          <Route path="partners" element={<PartnersPage />} />
          <Route path="partners/pending" element={<PartnersPendingPage />} />
          <Route path="partners/:id" element={<PartnerDetailPage />} />
          
          {/* Offers */}
          <Route path="offers" element={<OffersPage />} />
          <Route path="offers/pending" element={<OffersPendingPage />} />
          <Route path="offers/:id" element={<OfferDetailPage />} />
          
          {/* Identity Verifications */}
          <Route path="identity" element={<IdentityVerificationsPage />} />
          <Route path="identity/:id" element={<IdentityDetailPage />} />
          
          {/* Reviews */}
          <Route path="reviews" element={<ReviewsPage />} />
          <Route path="reviews/reported" element={<ReviewsReportedPage />} />
          
          {/* Subscriptions */}
          <Route path="subscriptions" element={<SubscriptionsPage />} />
          <Route path="subscriptions/plans" element={<PlansPage />} />
          
          {/* Analytics */}
          <Route path="analytics" element={<AnalyticsPage />} />
          
          {/* Settings */}
          <Route path="settings/categories" element={<CategoriesPage />} />
          <Route path="settings/config" element={<ConfigPage />} />
          <Route path="settings/team" element={<TeamPage />} />
        </Route>
        
        {/* 404 */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
      
      <Toaster />
    </>
  )
}
