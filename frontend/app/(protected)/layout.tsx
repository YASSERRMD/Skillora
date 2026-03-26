import { LayoutHeader } from "@/components/layout-header";

export default function ProtectedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="relative flex min-h-screen flex-col bg-background">
      <LayoutHeader />
      <main className="flex-1">{children}</main>
    </div>
  );
}
