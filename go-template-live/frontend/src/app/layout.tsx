import type { Metadata } from "next";
import { Geist, Geist_Mono, Outfit } from "next/font/google";
import Metrics from '@/components/metrics/index';
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

const outfit = Outfit({
  variable: "--font-outfit",
  subsets: ["latin"],
  weight: ["400", "500", "600", "700", "800", "900"],
  display: "swap",
});

const BASE_PATH = process.env.NEXT_PUBLIC_BASE_PATH || '';
const SITE_URL = process.env.NEXT_PUBLIC_SITE_URL || '';

export const metadata: Metadata = {
  title: "Golang Template Live Preview | Plify",
  description: "Free online Go template editor with instant live preview, automatic variable extraction, and side-by-side diff comparison. Test Golang text/template syntax, debug template rendering, and visualize changes in real-time. Perfect for Go developers building dynamic templates.",
  keywords: "golang template editor, golang template live preview, go template variable extraction, golang template syntax, go template debugging, online go editor, template diff viewer",
  icons: {
    icon: `${BASE_PATH}/favicon.ico`,
  },
  openGraph: {
    title: "Golang Template Live Preview | Plify",
    description: "Free online Go template editor with instant live preview, automatic variable extraction, and side-by-side diff comparison. Test Golang text/template syntax, debug template rendering, and visualize changes in real-time. Perfect for Go developers building dynamic templates.",
    images: [
      {
        url: `${SITE_URL}${BASE_PATH}/logo_name_og.png`,
        width: 1200,
        height: 630,
        alt: "Plify Logo",
      },
    ],
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        <Metrics />
      </head>
      <body
        className={`${geistSans.variable} ${geistMono.variable} ${outfit.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}
