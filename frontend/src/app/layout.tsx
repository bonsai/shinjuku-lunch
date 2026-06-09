import type { Metadata, Viewport } from "next"
import { Geist, Geist_Mono } from "next/font/google"
import "./globals.css"

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
})

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
})

export const metadata: Metadata = {
  title: "新宿ランチナビ",
  description: "新宿・歌舞伎町・大久保のランチ情報ナビゲーター",
  appleWebApp: { capable: true, title: "新宿ランチナビ", statusBarStyle: "black-translucent" },
}

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  maximumScale: 1,
  viewportFit: "cover",
  themeColor: "#0f1117",
}

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="ja" className={`${geistSans.variable} ${geistMono.variable} h-full`}>
      <body className="min-h-full bg-zinc-900 text-zinc-100 flex flex-col font-sans antialiased">
        <div className="flex-1 w-full sm:w-4/5 mx-auto px-4 sm:px-6">
          {children}
        </div>
      </body>
    </html>
  )
}
