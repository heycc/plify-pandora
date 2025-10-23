// URL sharing utilities for template sharing - simplified version
// Only template text is shared since variables are automatically extracted

import pako from 'pako';

// Safe URL length limits for cross-browser compatibility
// Conservative limit to ensure compatibility with all browsers including mobile
const SAFE_URL_LENGTH_LIMIT = 2000;

/**
 * Compress and encode template text for URL sharing
 */
export function encodeTemplateText(templateText: string): string {
  try {
    // Compress using gzip
    const compressed = pako.gzip(templateText);

    // Convert to base64url (URL-safe base64)
    const base64 = btoa(String.fromCharCode(...compressed));
    const base64url = base64
      .replace(/\+/g, '-')
      .replace(/\//g, '_')
      .replace(/=/g, '');

    return base64url;
  } catch (error) {
    console.error('Failed to encode template text:', error);
    throw new Error('Failed to encode template text for sharing');
  }
}

/**
 * Decode and decompress template text from URL parameter
 */
export function decodeTemplateText(encodedData: string): string | null {
  try {
    // Convert from base64url to base64
    let base64 = encodedData
      .replace(/-/g, '+')
      .replace(/_/g, '/');

    // Add padding if needed
    while (base64.length % 4) {
      base64 += '=';
    }

    // Convert base64 to Uint8Array
    const binaryString = atob(base64);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i);
    }

    // Decompress using gzip
    const decompressed = pako.ungzip(bytes);

    // Convert back to string
    return new TextDecoder().decode(decompressed);
  } catch (error) {
    console.error('Failed to decode template text:', error);
    return null;
  }
}

/**
 * Generate shareable URL with encoded template text
 * Returns null if the URL would exceed safe length limits
 */
export function generateShareableUrl(templateText: string): string | null {
  try {
    const encodedData = encodeTemplateText(templateText);
    const currentUrl = window.location.origin + window.location.pathname;
    const shareableUrl = `${currentUrl}?template=${encodedData}`;

    // Check if URL exceeds safe length limits
    if (shareableUrl.length > SAFE_URL_LENGTH_LIMIT) {
      return null;
    }

    return shareableUrl;
  } catch (error) {
    console.error('Failed to generate shareable URL:', error);
    return null;
  }
}

/**
 * Extract template text from current URL
 */
export function getTemplateTextFromUrl(): string | null {
  const urlParams = new URLSearchParams(window.location.search);
  const encodedData = urlParams.get('template');

  if (!encodedData) {
    return null;
  }

  return decodeTemplateText(encodedData);
}

/**
 * Copy shareable URL to clipboard
 */
export async function copyShareableUrl(templateText: string): Promise<boolean> {
  try {
    const url = generateShareableUrl(templateText);

    // Check if URL generation failed due to length limits
    if (!url) {
      return false;
    }

    await navigator.clipboard.writeText(url);
    return true;
  } catch (error) {
    console.error('Failed to copy URL to clipboard:', error);

    // Fallback for older browsers
    try {
      const url = generateShareableUrl(templateText);

      if (!url) {
        return false;
      }

      const textArea = document.createElement('textarea');
      textArea.value = url;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
      return true;
    } catch (fallbackError) {
      console.error('Fallback copy also failed:', fallbackError);
      return false;
    }
  }
}

/**
 * Check if current URL contains a shared template
 */
export function hasSharedTemplateInUrl(): boolean {
  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.has('template');
}

/**
 * Check if template is too large for URL sharing
 * Only checks the final encoded string length, not the original text
 */
export function isTemplateTooLarge(templateText: string): boolean {
  // Only check the final encoded URL length
  const url = generateShareableUrl(templateText);
  return url === null;
}

/**
 * Estimate compression ratio for user feedback
 */
export function getTemplateSizeInfo(templateText: string): {
  originalSize: number;
  estimatedUrlSize: number;
  isTooLarge: boolean;
} {
  const originalSize = templateText.length;
  const url = generateShareableUrl(templateText);
  const estimatedUrlSize = url ? url.length : SAFE_URL_LENGTH_LIMIT + 1;
  const isTooLarge = estimatedUrlSize > SAFE_URL_LENGTH_LIMIT;

  return {
    originalSize,
    estimatedUrlSize,
    isTooLarge
  };
}