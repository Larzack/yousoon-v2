import { test, expect } from '@playwright/test';

test.describe('Homepage', () => {
  test('should display the hero section', async ({ page }) => {
    await page.goto('/');
    
    // Check that the hero title is visible
    await expect(page.getByRole('heading', { level: 1 })).toBeVisible();
    
    // Check that CTA buttons are present
    await expect(page.getByRole('link', { name: /télécharger|download/i })).toBeVisible();
  });

  test('should navigate to features page', async ({ page }) => {
    await page.goto('/');
    
    // Click on features link
    await page.getByRole('link', { name: /fonctionnalités|features/i }).first().click();
    
    // Verify we're on the features page
    await expect(page).toHaveURL(/fonctionnalites/);
  });

  test('should display the footer', async ({ page }) => {
    await page.goto('/');
    
    // Scroll to bottom
    await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight));
    
    // Check footer is visible
    await expect(page.getByRole('contentinfo')).toBeVisible();
  });
});

test.describe('Navigation', () => {
  test('should navigate between pages', async ({ page }) => {
    // Start at homepage
    await page.goto('/');
    
    // Navigate to pricing
    await page.getByRole('link', { name: /tarifs|pricing/i }).first().click();
    await expect(page).toHaveURL(/tarifs/);
    
    // Navigate to about
    await page.getByRole('link', { name: /à propos|about/i }).first().click();
    await expect(page).toHaveURL(/a-propos/);
    
    // Navigate to contact
    await page.getByRole('link', { name: /contact/i }).first().click();
    await expect(page).toHaveURL(/contact/);
  });

  test('mobile menu should work', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/');
    
    // Open mobile menu
    const menuButton = page.getByRole('button', { name: /menu/i });
    if (await menuButton.isVisible()) {
      await menuButton.click();
      
      // Check menu items are visible
      await expect(page.getByRole('link', { name: /fonctionnalités|features/i })).toBeVisible();
    }
  });
});

test.describe('Contact Form', () => {
  test('should display contact form', async ({ page }) => {
    await page.goto('/contact');
    
    // Check form fields are present
    await expect(page.getByLabel(/nom|name/i)).toBeVisible();
    await expect(page.getByLabel(/email/i)).toBeVisible();
    await expect(page.getByLabel(/message/i)).toBeVisible();
    
    // Check submit button
    await expect(page.getByRole('button', { name: /envoyer|send/i })).toBeVisible();
  });

  test('should validate required fields', async ({ page }) => {
    await page.goto('/contact');
    
    // Try to submit empty form
    await page.getByRole('button', { name: /envoyer|send/i }).click();
    
    // Form should not submit (still on contact page)
    await expect(page).toHaveURL(/contact/);
  });
});

test.describe('Accessibility', () => {
  test('should have proper heading hierarchy', async ({ page }) => {
    await page.goto('/');
    
    // Check there's exactly one h1
    const h1Count = await page.locator('h1').count();
    expect(h1Count).toBe(1);
  });

  test('should have alt text for images', async ({ page }) => {
    await page.goto('/');
    
    // Get all images
    const images = page.locator('img');
    const count = await images.count();
    
    for (let i = 0; i < count; i++) {
      const img = images.nth(i);
      const alt = await img.getAttribute('alt');
      // Either alt should exist or the image should be decorative (aria-hidden)
      const ariaHidden = await img.getAttribute('aria-hidden');
      expect(alt !== null || ariaHidden === 'true').toBeTruthy();
    }
  });

  test('should have proper link text', async ({ page }) => {
    await page.goto('/');
    
    // Check that links have accessible text
    const links = page.locator('a:not([aria-hidden="true"])');
    const count = await links.count();
    
    for (let i = 0; i < count; i++) {
      const link = links.nth(i);
      const text = await link.textContent();
      const ariaLabel = await link.getAttribute('aria-label');
      const title = await link.getAttribute('title');
      
      // Link should have some accessible name
      expect(text?.trim() || ariaLabel || title).toBeTruthy();
    }
  });
});

test.describe('SEO', () => {
  test('should have meta title', async ({ page }) => {
    await page.goto('/');
    
    const title = await page.title();
    expect(title).toBeTruthy();
    expect(title.length).toBeGreaterThan(10);
  });

  test('should have meta description', async ({ page }) => {
    await page.goto('/');
    
    const description = await page.locator('meta[name="description"]').getAttribute('content');
    expect(description).toBeTruthy();
    expect(description?.length).toBeGreaterThan(50);
  });

  test('should have Open Graph tags', async ({ page }) => {
    await page.goto('/');
    
    const ogTitle = await page.locator('meta[property="og:title"]').getAttribute('content');
    const ogDescription = await page.locator('meta[property="og:description"]').getAttribute('content');
    
    expect(ogTitle).toBeTruthy();
    expect(ogDescription).toBeTruthy();
  });
});

test.describe('Performance', () => {
  test('should load homepage within 3 seconds', async ({ page }) => {
    const startTime = Date.now();
    
    await page.goto('/', { waitUntil: 'domcontentloaded' });
    
    const loadTime = Date.now() - startTime;
    expect(loadTime).toBeLessThan(3000);
  });
});

test.describe('i18n', () => {
  test('should support French', async ({ page }) => {
    await page.goto('/fr');
    
    // Check French content is displayed
    await expect(page.locator('html')).toHaveAttribute('lang', 'fr');
  });

  test('should support English', async ({ page }) => {
    await page.goto('/en');
    
    // Check English content is displayed
    await expect(page.locator('html')).toHaveAttribute('lang', 'en');
  });
});
