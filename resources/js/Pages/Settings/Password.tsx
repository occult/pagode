import { type BreadcrumbItem } from "@/types";
import { Head } from "@inertiajs/react";

import AppLayout from "@/Layouts/AppLayout";
import SettingsLayout from "@/Layouts/Settings/Layout";
import HeadingSmall from "@/components/HeadingSmall";

const breadcrumbs: BreadcrumbItem[] = [
  {
    title: "Password settings",
    href: "/settings/password",
  },
];

export default function Password() {
  return (
    <AppLayout breadcrumbs={breadcrumbs}>
      <Head title="Password settings" />
      <SettingsLayout>
        <div className="space-y-6">
          <HeadingSmall
            title="Update password"
            description="Your password is managed by your identity provider"
          />
          <p className="text-sm text-muted-foreground">
            To change your password, please visit your Casdoor account settings.
          </p>
        </div>
      </SettingsLayout>
    </AppLayout>
  );
}
