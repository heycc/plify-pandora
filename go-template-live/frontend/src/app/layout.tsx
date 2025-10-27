import type { Metadata } from "next";
import { Geist, Geist_Mono, Outfit } from "next/font/google";
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
  title: "Plify - Golang Template Live",
  description: "Type in Golang template, extract variables, render and view the diff in real-time",
  icons: {
    icon: `${BASE_PATH}/favicon.ico`,
  },
  openGraph: {
    title: "Plify - Golang Template Live",
    description: "Type in Golang template, extract variables, render and view the diff in real-time",
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
      <body
        className={`${geistSans.variable} ${geistMono.variable} ${outfit.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}
