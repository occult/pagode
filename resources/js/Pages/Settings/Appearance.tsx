import { Head } from "@inertiajs/react";

import { type BreadcrumbItem } from "@/types";

import AppearanceToggleTab from "@/components/AppearanceTabs";
import HeadingSmall from "@/components/HeadingSmall";
import AppLayout from "@/Layouts/AppLayout";
import SettingsLayout from "@/Layouts/Settings/Layout";

const breadcrumbs: BreadcrumbItem[] = [
  {
    title: "Appearance settings",
    href: "/settings/appearance",
  },
];

export default function Appearance() {
  return (
    <AppLayout breadcrumbs={breadcrumbs}>
      <Head title="Appearance settings" />

      <SettingsLayout>
        <div className="space-y-6">
          <HeadingSmall
            title="Appearance settings"
            description="Update your account's appearance settings"
          />
          <AppearanceToggleTab />
        </div>
      </SettingsLayout>
    </AppLayout>
  );
}
