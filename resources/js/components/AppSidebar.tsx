import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { type NavItem } from "@/types";
import { Link, usePage } from "@inertiajs/react";
import {
  BookOpen,
  Folder,
  LayoutGrid,
  // [feature:files] start
  UploadCloud,
  // [feature:files] end
  // [feature:payment] start
  CreditCard,
  Receipt,
  ShoppingBag,
  Crown,
  // [feature:payment] end
  // [feature:chat] start
  MessageCircle,
  // [feature:chat] end
} from "lucide-react";
import { NavMain } from "./NavMain";
import { NavFooter } from "./NavFooter";
import { NavUser } from "./NavUser";
import { SharedProps } from "@/types/global";

const footerNavItems: NavItem[] = [
  {
    title: "Repository",
    href: "https://github.com/occult/pagode",
    icon: Folder,
  },
  {
    title: "Documentation",
    href: "https://github.com/occult/pagode?tab=readme-ov-file#introduction",
    icon: BookOpen,
  },
];

export function AppSidebar() {
  const { auth } = usePage<SharedProps>().props;

  const mainNavItems: NavItem[] = [
    {
      title: "Dashboard",
      href: "/dashboard",
      icon: LayoutGrid,
    },
    // [feature:chat] start
    {
      title: "Chat",
      href: "/chat",
      icon: MessageCircle,
    },
    // [feature:chat] end
    // [feature:payment] start
    {
      title: "Plans",
      href: "/plans",
      icon: CreditCard,
    },
    {
      title: "Products",
      href: "/products",
      icon: ShoppingBag,
    },
    {
      title: "Premium",
      href: "/premium",
      icon: Crown,
    },
    // [feature:payment] end
    // [feature:files] start
    {
      title: "Upload Files",
      href: "/files",
      icon: UploadCloud,
    },
    // [feature:files] end
    // [feature:payment] start
    {
      title: "Billing",
      href: "/billing",
      icon: Receipt,
    },
    // [feature:payment] end
  ];

  if (auth.user?.admin) {
    mainNavItems.push({
      title: "Admin Panel",
      href: "/admin/users",
      icon: LayoutGrid,
    });
  }

  return (
    <Sidebar collapsible="icon" variant="inset">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <Link href="/" prefetch>
                <span>Pagode</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent>
        <NavMain items={mainNavItems} />
      </SidebarContent>

      <SidebarFooter>
        <NavFooter items={footerNavItems} className="mt-auto" />
        <NavUser />
      </SidebarFooter>
    </Sidebar>
  );
}
