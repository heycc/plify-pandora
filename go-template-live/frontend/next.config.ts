import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Acknowledge Turbopack - Monaco works fine without webpack plugin in Turbopack
  turbopack: {},
};

export default nextConfig;
