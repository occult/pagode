import { AppContent } from "@/components/AppContent";
import { AppShell } from "@/components/AppShell";
import { AppSidebar } from "@/components/AppSidebar";
import { AppSidebarHeader } from "@/components/AppSidebarHeader";
import { type BreadcrumbItem } from "@/types";
import { type PropsWithChildren } from "react";

export default function AppSidebarLayout({
  children,
  breadcrumbs = [],
}: PropsWithChildren<{ breadcrumbs?: BreadcrumbItem[] }>) {
  return (
    <AppShell variant="sidebar">
      <AppSidebar />
      <AppContent variant="sidebar" className="overflow-x-hidden flex flex-col">
        <AppSidebarHeader breadcrumbs={breadcrumbs} />
        <div className="flex-1 overflow-y-auto">{children}</div>
      </AppContent>
    </AppShell>
  );
}
