import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import { format, formatDistanceToNow } from 'date-fns'
import { fr } from 'date-fns/locale'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDate(date: string | Date, pattern = 'dd MMM yyyy') {
  return format(new Date(date), pattern, { locale: fr })
}

export function formatDateTime(date: string | Date) {
  return format(new Date(date), 'dd MMM yyyy à HH:mm', { locale: fr })
}

export function formatRelative(date: string | Date) {
  return formatDistanceToNow(new Date(date), { addSuffix: true, locale: fr })
}

export function truncate(str: string, length: number) {
  return str.length > length ? `${str.substring(0, length)}...` : str
}

export function getInitials(firstName: string, lastName: string) {
  return `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase()
}

export function getStatusColor(status: string) {
  const colors: Record<string, string> = {
    // User/Partner status
    active: 'bg-green-100 text-green-800',
    pending: 'bg-yellow-100 text-yellow-800',
    suspended: 'bg-red-100 text-red-800',
    deleted: 'bg-gray-100 text-gray-800',
    // Verification status
    verified: 'bg-green-100 text-green-800',
    rejected: 'bg-red-100 text-red-800',
    not_submitted: 'bg-gray-100 text-gray-800',
    // Offer status
    draft: 'bg-gray-100 text-gray-800',
    paused: 'bg-yellow-100 text-yellow-800',
    expired: 'bg-orange-100 text-orange-800',
    archived: 'bg-gray-100 text-gray-800',
    // Booking status
    confirmed: 'bg-blue-100 text-blue-800',
    checked_in: 'bg-green-100 text-green-800',
    cancelled: 'bg-red-100 text-red-800',
    no_show: 'bg-orange-100 text-orange-800',
    // Review status
    approved: 'bg-green-100 text-green-800',
    reported: 'bg-red-100 text-red-800',
  }
  return colors[status.toLowerCase()] || 'bg-gray-100 text-gray-800'
}

export function getStatusLabel(status: string) {
  const labels: Record<string, string> = {
    active: 'Actif',
    pending: 'En attente',
    suspended: 'Suspendu',
    deleted: 'Supprimé',
    verified: 'Vérifié',
    rejected: 'Rejeté',
    not_submitted: 'Non soumis',
    draft: 'Brouillon',
    paused: 'En pause',
    expired: 'Expiré',
    archived: 'Archivé',
    confirmed: 'Confirmé',
    checked_in: 'Check-in',
    cancelled: 'Annulé',
    no_show: 'No-show',
    approved: 'Approuvé',
    reported: 'Signalé',
  }
  return labels[status.toLowerCase()] || status
}
