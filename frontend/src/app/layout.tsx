import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Agendamento Automotivo",
  description: "Gestão de serviços e agenda automotiva",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return <html lang="pt-BR"><body>{children}</body></html>;
}
