import { describe, it, expect } from 'vitest'
import {
  cn,
  formatDate,
  formatDateTime,
  formatRelative,
  truncate,
  getInitials,
  getStatusColor,
  getStatusLabel,
} from '../utils'

describe('utils', () => {
  describe('cn (className merge)', () => {
    it('should merge class names', () => {
      const result = cn('px-4', 'py-2')
      expect(result).toBe('px-4 py-2')
    })

    it('should handle conditional classes', () => {
      const isActive = true
      const result = cn('base', isActive && 'active')
      expect(result).toBe('base active')
    })

    it('should handle falsy conditions', () => {
      const isActive = false
      const result = cn('base', isActive && 'active')
      expect(result).toBe('base')
    })

    it('should merge tailwind classes correctly', () => {
      const result = cn('p-4', 'p-2')
      expect(result).toBe('p-2')
    })

    it('should handle arrays', () => {
      const result = cn(['px-4', 'py-2'])
      expect(result).toBe('px-4 py-2')
    })

    it('should handle objects', () => {
      const result = cn({ 'px-4': true, 'py-2': false })
      expect(result).toBe('px-4')
    })
  })

  describe('formatDate', () => {
    it('should format date with default pattern', () => {
      const result = formatDate('2024-12-10')
      expect(result).toMatch(/10 déc\. 2024/i)
    })

    it('should format date with custom pattern', () => {
      const result = formatDate('2024-12-10', 'yyyy-MM-dd')
      expect(result).toBe('2024-12-10')
    })

    it('should handle Date objects', () => {
      const date = new Date(2024, 11, 10) // Month is 0-indexed
      const result = formatDate(date)
      expect(result).toMatch(/10 déc\. 2024/i)
    })
  })

  describe('formatDateTime', () => {
    it('should format date with time', () => {
      const result = formatDateTime('2024-12-10T14:30:00')
      expect(result).toMatch(/10 déc\. 2024/i)
      expect(result).toMatch(/14:30/)
    })

    it('should include "à" separator', () => {
      const result = formatDateTime('2024-12-10T14:30:00')
      expect(result).toContain('à')
    })
  })

  describe('formatRelative', () => {
    it('should return relative time', () => {
      const now = new Date()
      const result = formatRelative(now)
      expect(result).toMatch(/il y a|moins|maintenant/i)
    })

    it('should include suffix', () => {
      const yesterday = new Date(Date.now() - 24 * 60 * 60 * 1000)
      const result = formatRelative(yesterday)
      expect(result).toMatch(/il y a/i)
    })
  })

  describe('truncate', () => {
    it('should truncate long strings', () => {
      const result = truncate('This is a very long string', 10)
      expect(result).toBe('This is a ...')
    })

    it('should not truncate short strings', () => {
      const result = truncate('Short', 10)
      expect(result).toBe('Short')
    })

    it('should handle exact length', () => {
      const result = truncate('Exactly', 7)
      expect(result).toBe('Exactly')
    })

    it('should handle empty string', () => {
      const result = truncate('', 10)
      expect(result).toBe('')
    })
  })

  describe('getInitials', () => {
    it('should return initials in uppercase', () => {
      const result = getInitials('John', 'Doe')
      expect(result).toBe('JD')
    })

    it('should handle lowercase names', () => {
      const result = getInitials('jean', 'dupont')
      expect(result).toBe('JD')
    })

    it('should handle single character names', () => {
      const result = getInitials('A', 'B')
      expect(result).toBe('AB')
    })
  })

  describe('getStatusColor', () => {
    it('should return green for active status', () => {
      const result = getStatusColor('active')
      expect(result).toContain('green')
    })

    it('should return green for verified status', () => {
      const result = getStatusColor('verified')
      expect(result).toContain('green')
    })

    it('should return yellow for pending status', () => {
      const result = getStatusColor('pending')
      expect(result).toContain('yellow')
    })

    it('should return red for suspended status', () => {
      const result = getStatusColor('suspended')
      expect(result).toContain('red')
    })

    it('should return red for rejected status', () => {
      const result = getStatusColor('rejected')
      expect(result).toContain('red')
    })

    it('should return gray for draft status', () => {
      const result = getStatusColor('draft')
      expect(result).toContain('gray')
    })

    it('should return gray for unknown status', () => {
      const result = getStatusColor('unknown_status')
      expect(result).toContain('gray')
    })

    it('should be case insensitive', () => {
      const result = getStatusColor('ACTIVE')
      expect(result).toContain('green')
    })

    it('should handle booking statuses', () => {
      expect(getStatusColor('confirmed')).toContain('blue')
      expect(getStatusColor('checked_in')).toContain('green')
      expect(getStatusColor('cancelled')).toContain('red')
      expect(getStatusColor('no_show')).toContain('orange')
    })
  })

  describe('getStatusLabel', () => {
    it('should return French label for active', () => {
      expect(getStatusLabel('active')).toBe('Actif')
    })

    it('should return French label for pending', () => {
      expect(getStatusLabel('pending')).toBe('En attente')
    })

    it('should return French label for verified', () => {
      expect(getStatusLabel('verified')).toBe('Vérifié')
    })

    it('should return French label for rejected', () => {
      expect(getStatusLabel('rejected')).toBe('Rejeté')
    })

    it('should return original status for unknown', () => {
      expect(getStatusLabel('custom_status')).toBe('custom_status')
    })

    it('should be case insensitive', () => {
      expect(getStatusLabel('ACTIVE')).toBe('Actif')
    })

    it('should handle all user statuses', () => {
      expect(getStatusLabel('suspended')).toBe('Suspendu')
      expect(getStatusLabel('deleted')).toBe('Supprimé')
    })

    it('should handle all offer statuses', () => {
      expect(getStatusLabel('draft')).toBe('Brouillon')
      expect(getStatusLabel('paused')).toBe('En pause')
      expect(getStatusLabel('expired')).toBe('Expiré')
      expect(getStatusLabel('archived')).toBe('Archivé')
    })

    it('should handle all booking statuses', () => {
      expect(getStatusLabel('confirmed')).toBe('Confirmé')
      expect(getStatusLabel('checked_in')).toBe('Check-in')
      expect(getStatusLabel('cancelled')).toBe('Annulé')
      expect(getStatusLabel('no_show')).toBe('No-show')
    })

    it('should handle review statuses', () => {
      expect(getStatusLabel('approved')).toBe('Approuvé')
      expect(getStatusLabel('reported')).toBe('Signalé')
    })
  })
})
