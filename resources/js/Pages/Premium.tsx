import AppLayout from "@/Layouts/AppLayout";
import { type BreadcrumbItem } from "@/types";
import { Head } from "@inertiajs/react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Crown, Sparkles, Star } from "lucide-react";

const breadcrumbs: BreadcrumbItem[] = [
  {
    title: "Dashboard",
    href: "/dashboard",
  },
  {
    title: "Premium",
    href: "/premium",
  },
];

interface PremiumProps {
  title: string;
  user: {
    id: number;
    name: string;
    email: string;
  };
}

export default function Premium({ title, user }: PremiumProps) {
  return (
    <AppLayout breadcrumbs={breadcrumbs}>
      <Head title={title} />
      <div className="flex h-full flex-1 flex-col gap-6 rounded-xl p-6">
        {/* Header */}
        <div className="text-center space-y-4">
          <div className="flex justify-center">
            <div className="relative">
              <Crown className="h-16 w-16 text-foreground" />
              <Sparkles className="h-6 w-6 text-muted-foreground absolute -top-2 -right-2 animate-pulse" />
            </div>
          </div>
          <h1 className="text-4xl font-bold text-foreground">
            Premium Page
          </h1>
          <p className="text-lg text-muted-foreground">
            Welcome to the exclusive premium content area
          </p>
        </div>

        {/* Premium Content */}
        <div className="grid gap-6 md:grid-cols-1 lg:grid-cols-2 xl:grid-cols-3">
          <Card className="border-sidebar-border/70 dark:border-sidebar-border">
            <CardHeader className="text-center">
              <div className="flex justify-center mb-2">
                <Star className="h-8 w-8 text-foreground fill-current" />
              </div>
              <CardTitle className="text-foreground">Exclusive Access</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-muted-foreground text-center">
                You have successfully unlocked premium features and content available only to paid users.
              </p>
            </CardContent>
          </Card>

          <Card className="border-sidebar-border/70 dark:border-sidebar-border">
            <CardHeader className="text-center">
              <div className="flex justify-center mb-2">
                <Sparkles className="h-8 w-8 text-foreground" />
              </div>
              <CardTitle className="text-foreground">Premium Features</CardTitle>
            </CardHeader>
            <CardContent>
              <ul className="space-y-2 text-muted-foreground">
                <li className="flex items-center gap-2">
                  <div className="h-2 w-2 bg-foreground rounded-full"></div>
                  Advanced Analytics
                </li>
                <li className="flex items-center gap-2">
                  <div className="h-2 w-2 bg-foreground rounded-full"></div>
                  Priority Support
                </li>
                <li className="flex items-center gap-2">
                  <div className="h-2 w-2 bg-foreground rounded-full"></div>
                  Custom Integrations
                </li>
              </ul>
            </CardContent>
          </Card>

          <Card className="border-sidebar-border/70 dark:border-sidebar-border">
            <CardHeader className="text-center">
              <div className="flex justify-center mb-2">
                <Crown className="h-8 w-8 text-foreground" />
              </div>
              <CardTitle className="text-foreground">Welcome, {user.name}!</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-muted-foreground text-center">
                Thank you for being a premium member. Enjoy your exclusive access to advanced features.
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Additional Premium Content */}
        <Card className="border-sidebar-border/70 dark:border-sidebar-border">
          <CardHeader>
            <CardTitle className="text-foreground flex items-center gap-2">
              <Sparkles className="h-5 w-5" />
              Premium Content Area
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <p className="text-muted-foreground">
                This section contains premium content that is only accessible to users who have:
              </p>
              <ul className="space-y-2 text-muted-foreground">
                <li className="flex items-center gap-2">
                  <div className="h-2 w-2 bg-foreground rounded-full"></div>
                  Active subscription to our premium plans
                </li>
                <li className="flex items-center gap-2">
                  <div className="h-2 w-2 bg-foreground rounded-full"></div>
                  Completed one-time product purchases
                </li>
                <li className="flex items-center gap-2">
                  <div className="h-2 w-2 bg-foreground rounded-full"></div>
                  Valid payment verification
                </li>
              </ul>
              <div className="mt-6 p-4 bg-muted rounded-lg border border-sidebar-border/70 dark:border-sidebar-border">
                <p className="text-foreground font-medium text-center">
                  🎉 Congratulations! You have successfully accessed the premium area.
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </AppLayout>
  );
}