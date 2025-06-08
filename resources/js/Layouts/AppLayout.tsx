import { usePage } from "@inertiajs/react";
import AppLayoutTemplate from "@/Layouts/App/AppSidebarLayout";
import { type BreadcrumbItem } from "@/types";
import { type ReactNode } from "react";
import { useFlashToasts } from "@/hooks/useFlashToast";
import { SharedProps } from "@/types/global";

interface AppLayoutProps {
  children: ReactNode;
  breadcrumbs?: BreadcrumbItem[];
}

export default function AppLayout({
  children,
  breadcrumbs,
  ...props
}: AppLayoutProps) {
  const { flash } = usePage<SharedProps>().props;

  useFlashToasts(flash);

  return (
    <AppLayoutTemplate breadcrumbs={breadcrumbs} {...props}>
      {children}
    </AppLayoutTemplate>
  );
}
