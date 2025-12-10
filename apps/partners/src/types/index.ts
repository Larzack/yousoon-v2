// Types principaux pour le Site Partenaires Yousoon

// ============================================
// Enums
// ============================================

export enum OfferStatus {
  DRAFT = 'DRAFT',
  PENDING = 'PENDING',
  ACTIVE = 'ACTIVE',
  PAUSED = 'PAUSED',
  EXPIRED = 'EXPIRED',
  ARCHIVED = 'ARCHIVED',
}

export enum BookingStatus {
  PENDING = 'PENDING',
  CONFIRMED = 'CONFIRMED',
  CHECKED_IN = 'CHECKED_IN',
  CANCELLED = 'CANCELLED',
  EXPIRED = 'EXPIRED',
  NO_SHOW = 'NO_SHOW',
}

export enum TeamRole {
  ADMIN = 'ADMIN',
  MANAGER = 'MANAGER',
  STAFF = 'STAFF',
  VIEWER = 'VIEWER',
}

export enum DiscountType {
  PERCENTAGE = 'PERCENTAGE',
  FIXED = 'FIXED',
  FORMULA = 'FORMULA',
}

export enum PartnerStatus {
  PENDING = 'PENDING',
  ACTIVE = 'ACTIVE',
  SUSPENDED = 'SUSPENDED',
}

// ============================================
// User & Auth
// ============================================

export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  avatar?: string;
  role: TeamRole;
  createdAt: string;
  updatedAt: string;
}

export interface Partner {
  id: string;
  company: Company;
  branding: Branding;
  contact: Contact;
  category: string;
  subcategories: string[];
  status: PartnerStatus;
  verifiedAt?: string;
  stats: PartnerStats;
  createdAt: string;
  updatedAt: string;
}

export interface Company {
  name: string;
  tradeName?: string;
  siret: string;
  vatNumber?: string;
  legalForm?: string;
}

export interface Branding {
  logo?: string;
  coverImage?: string;
  primaryColor?: string;
  description?: string;
}

export interface Contact {
  firstName: string;
  lastName: string;
  email: string;
  phone?: string;
  role?: string;
}

export interface PartnerStats {
  totalOffers: number;
  activeOffers: number;
  totalBookings: number;
  totalCheckins: number;
  avgRating: number;
  reviewCount: number;
  lastUpdated: string;
}

// ============================================
// Establishment
// ============================================

export interface Establishment {
  id: string;
  partnerId: string;
  name: string;
  description?: string;
  address: Address;
  location: GeoLocation;
  contact?: EstablishmentContact;
  openingHours: OpeningHour[];
  images: Image[];
  features: string[];
  type?: string;
  priceRange?: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Address {
  street: string;
  streetNumber?: string;
  complement?: string;
  postalCode: string;
  city: string;
  country: string;
  formatted?: string;
}

export interface GeoLocation {
  type: 'Point';
  coordinates: [number, number]; // [longitude, latitude]
}

export interface EstablishmentContact {
  phone?: string;
  email?: string;
  website?: string;
}

export interface OpeningHour {
  dayOfWeek: number; // 0 = Dimanche
  open: string; // "09:00"
  close: string; // "23:00"
  isClosed: boolean;
}

export interface Image {
  url: string;
  alt?: string;
  isPrimary: boolean;
  order: number;
}

// ============================================
// Offer
// ============================================

export interface Offer {
  id: string;
  partnerId: string;
  establishmentId: string;
  title: string;
  description?: string;
  shortDescription?: string;
  categoryId: string;
  tags: string[];
  discount: Discount;
  conditions: OfferCondition[];
  termsAndConditions?: string;
  validity: Validity;
  schedule: Schedule;
  quota: Quota;
  images: Image[];
  stats: OfferStats;
  status: OfferStatus;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
  publishedAt?: string;
}

export interface Discount {
  type: DiscountType;
  value: number;
  originalPrice?: number;
  formula?: string;
}

export interface OfferCondition {
  type: string;
  value: string | number | boolean;
  label: string;
}

export interface Validity {
  startDate: string;
  endDate: string;
  timezone: string;
}

export interface Schedule {
  allDay: boolean;
  slots: TimeSlot[];
}

export interface TimeSlot {
  dayOfWeek: number;
  startTime: string;
  endTime: string;
}

export interface Quota {
  total?: number;
  perUser?: number;
  perDay?: number;
  used: number;
}

export interface OfferStats {
  views: number;
  bookings: number;
  checkins: number;
  favorites: number;
}

// ============================================
// Booking (Outing)
// ============================================

export interface Booking {
  id: string;
  userId: string;
  offerId: string;
  partnerId: string;
  establishmentId: string;
  qrCode: QRCode;
  status: BookingStatus;
  timeline: TimelineEvent[];
  checkin?: CheckinInfo;
  cancellation?: CancellationInfo;
  user: BookingUser;
  offer: BookingOffer;
  establishment: BookingEstablishment;
  createdAt: string;
  updatedAt: string;
  expiresAt: string;
}

export interface QRCode {
  code: string;
  data: string;
  expiresAt: string;
}

export interface TimelineEvent {
  status: BookingStatus;
  timestamp: string;
  actor: 'user' | 'partner' | 'system';
  metadata?: Record<string, unknown>;
}

export interface CheckinInfo {
  checkedInAt: string;
  checkedInBy?: string;
  method: 'qr_scan' | 'manual';
  location?: GeoLocation;
}

export interface CancellationInfo {
  cancelledAt: string;
  cancelledBy: 'user' | 'partner' | 'system';
  reason?: string;
}

export interface BookingUser {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  avatar?: string;
}

export interface BookingOffer {
  id: string;
  title: string;
  discount: Discount;
  images: string[];
}

export interface BookingEstablishment {
  id: string;
  name: string;
  address: string;
}

// ============================================
// Team Member
// ============================================

export interface TeamMember {
  id: string;
  userId?: string;
  email: string;
  firstName?: string;
  lastName?: string;
  avatar?: string;
  role: TeamRole;
  status: 'pending' | 'active' | 'inactive';
  invitedAt: string;
  joinedAt?: string;
}

// ============================================
// Category
// ============================================

export interface Category {
  id: string;
  name: {
    fr: string;
    en: string;
  };
  slug: string;
  description?: {
    fr: string;
    en: string;
  };
  icon?: string;
  color?: string;
  image?: string;
  parentId?: string;
  order: number;
  isActive: boolean;
}

// ============================================
// Analytics
// ============================================

export interface AnalyticsSummary {
  period: {
    start: string;
    end: string;
  };
  totalViews: number;
  totalBookings: number;
  totalCheckins: number;
  conversionRate: number;
  revenue?: number;
  trends: {
    views: number; // % change
    bookings: number;
    checkins: number;
  };
}

export interface DailyStats {
  date: string;
  views: number;
  bookings: number;
  checkins: number;
  revenue?: number;
}

export interface TopOffer {
  offerId: string;
  title: string;
  views: number;
  bookings: number;
  conversionRate: number;
}

// ============================================
// Notifications & Settings
// ============================================

export interface NotificationSettings {
  email: {
    newBooking: boolean;
    bookingCancelled: boolean;
    checkin: boolean;
    newReview: boolean;
    weeklyReport: boolean;
    marketing: boolean;
  };
  push: {
    newBooking: boolean;
    bookingCancelled: boolean;
    checkin: boolean;
  };
}

export interface SecuritySettings {
  twoFactorEnabled: boolean;
  lastPasswordChange?: string;
  activeSessions: number;
}

// ============================================
// API Response Types
// ============================================

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
  hasMore: boolean;
}

export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, string>;
}

// ============================================
// Form Types
// ============================================

export interface LoginFormData {
  email: string;
  password: string;
  rememberMe?: boolean;
}

export interface RegisterFormData {
  // Step 1: Personal info
  firstName: string;
  lastName: string;
  email: string;
  phone?: string;
  // Step 2: Business info
  companyName: string;
  tradeName?: string;
  siret: string;
  category: string;
  // Step 3: Password
  password: string;
  confirmPassword: string;
  acceptTerms: boolean;
}

export interface OfferFormData {
  // Step 1: Informations
  title: string;
  description: string;
  shortDescription?: string;
  categoryId: string;
  establishmentId: string;
  tags: string[];
  // Step 2: Réduction
  discountType: DiscountType;
  discountValue: number;
  originalPrice?: number;
  formula?: string;
  conditions: OfferCondition[];
  // Step 3: Planning
  startDate: string;
  endDate: string;
  allDay: boolean;
  slots: TimeSlot[];
  quotaTotal?: number;
  quotaPerUser?: number;
  quotaPerDay?: number;
  // Step 4: Options
  images: File[];
  termsAndConditions?: string;
}

export interface EstablishmentFormData {
  name: string;
  description?: string;
  street: string;
  streetNumber?: string;
  postalCode: string;
  city: string;
  country: string;
  phone?: string;
  email?: string;
  website?: string;
  openingHours: OpeningHour[];
  features: string[];
  type?: string;
  priceRange?: number;
  images: File[];
}

export interface TeamMemberFormData {
  email: string;
  firstName?: string;
  lastName?: string;
  role: TeamRole;
}
