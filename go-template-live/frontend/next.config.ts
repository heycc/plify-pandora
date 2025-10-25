import type { NextConfig } from "next";

// Configure the base path for hosting under a subdirectory
// Change this to match your nginx location path (e.g., '/template-parser')
const BASE_PATH = process.env.NEXT_PUBLIC_BASE_PATH || '';

const nextConfig: NextConfig = {
  // Acknowledge Turbopack - Monaco works fine without webpack plugin in Turbopack
  turbopack: {},
  // Enable static export
  output: 'export',
  // Base path for hosting under subdirectory (e.g., /template-parser)
  basePath: BASE_PATH,
  // Asset prefix ensures all assets are loaded from the correct path
  assetPrefix: BASE_PATH,
  // Disable image optimization for static export
  images: {
    unoptimized: true,
  },
  // Ensure trailing slash behavior is consistent
  trailingSlash: false,
};

export default nextConfig;
