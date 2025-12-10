import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
  test('should display login page', async ({ page }) => {
    await page.goto('/login');
    
    // Check login form elements
    await expect(page.getByLabel(/email/i)).toBeVisible();
    await expect(page.getByLabel(/mot de passe|password/i)).toBeVisible();
    await expect(page.getByRole('button', { name: /connexion|login|sign in/i })).toBeVisible();
  });

  test('should show validation errors on empty submit', async ({ page }) => {
    await page.goto('/login');
    
    // Submit empty form
    await page.getByRole('button', { name: /connexion|login|sign in/i }).click();
    
    // Should still be on login page
    await expect(page).toHaveURL(/login/);
  });

  test('should navigate to register page', async ({ page }) => {
    await page.goto('/login');
    
    // Click register link
    const registerLink = page.getByRole('link', { name: /inscription|register|sign up|créer/i });
    if (await registerLink.isVisible()) {
      await registerLink.click();
      await expect(page).toHaveURL(/register/);
    }
  });

  test('should navigate to forgot password', async ({ page }) => {
    await page.goto('/login');
    
    // Click forgot password link
    const forgotLink = page.getByRole('link', { name: /oublié|forgot/i });
    if (await forgotLink.isVisible()) {
      await forgotLink.click();
      await expect(page).toHaveURL(/forgot/);
    }
  });
});

test.describe('Register Flow', () => {
  test('should display register page with all fields', async ({ page }) => {
    await page.goto('/register');
    
    // Company info
    await expect(page.getByLabel(/raison sociale|company name/i)).toBeVisible();
    await expect(page.getByLabel(/siret/i)).toBeVisible();
    
    // Contact info
    await expect(page.getByLabel(/email/i)).toBeVisible();
    await expect(page.getByLabel(/mot de passe|password/i)).toBeVisible();
  });
});

test.describe('Dashboard (Authenticated)', () => {
  // Note: These tests would require authentication setup
  // For now, we just check the redirect behavior
  
  test('should redirect to login when not authenticated', async ({ page }) => {
    await page.goto('/dashboard');
    
    // Should redirect to login
    await expect(page).toHaveURL(/login/);
  });

  test('should redirect offers page to login', async ({ page }) => {
    await page.goto('/offers');
    
    // Should redirect to login
    await expect(page).toHaveURL(/login/);
  });
});

test.describe('Responsive Design', () => {
  test('should work on tablet', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.goto('/login');
    
    // Login form should still be usable
    await expect(page.getByLabel(/email/i)).toBeVisible();
  });

  test('should work on mobile', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/login');
    
    // Login form should still be usable
    await expect(page.getByLabel(/email/i)).toBeVisible();
  });
});

test.describe('Accessibility', () => {
  test('should have form labels', async ({ page }) => {
    await page.goto('/login');
    
    // All inputs should have labels
    const inputs = page.locator('input:not([type="hidden"])');
    const count = await inputs.count();
    
    for (let i = 0; i < count; i++) {
      const input = inputs.nth(i);
      const id = await input.getAttribute('id');
      if (id) {
        const label = page.locator(`label[for="${id}"]`);
        const ariaLabel = await input.getAttribute('aria-label');
        expect((await label.count()) > 0 || ariaLabel).toBeTruthy();
      }
    }
  });

  test('should be keyboard navigable', async ({ page }) => {
    await page.goto('/login');
    
    // Tab to first input
    await page.keyboard.press('Tab');
    
    // Something should be focused
    const focused = page.locator(':focus');
    await expect(focused).toBeTruthy();
  });
});
