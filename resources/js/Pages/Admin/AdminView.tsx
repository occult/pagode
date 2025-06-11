import { Head, Link, router } from "@inertiajs/react";
import AppLayout from "@/Layouts/AppLayout";
import { BreadcrumbItem, User } from "@/types";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useState } from "react";
import { UserFormModal } from "@/components/Admin/UserFormModal";

type Props = {
  users: User[];
  pagination: {
    total: number;
    page: number;
    perPage: number;
    totalPages: number;
  };
};

const breadcrumbs: BreadcrumbItem[] = [
  {
    title: "Admin",
    href: "/admin",
  },
];

export default function AdminView({ users, pagination }: Props) {
  const [showModal, setShowModal] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  const openAddModal = () => {
    setEditingUser(null);
    setShowModal(true);
  };

  const openEditModal = (user: User) => {
    setEditingUser(user);
    setShowModal(true);
  };

  const goToPage = (page: number) => {
    router.visit(`/admin/users?page=${page}`, {
      preserveScroll: true,
      preserveState: true,
    });
  };

  return (
    <AppLayout breadcrumbs={breadcrumbs}>
      <Head title="User" />

      <div className="flex flex-col w-full justify-start p-6">
        <div className="w-full space-y-6 flex-1">
          <Card className="h-full flex flex-col">
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>User</CardTitle>
              <Button onClick={openAddModal}>Add User</Button>
            </CardHeader>

            <CardContent>
              {users.length > 0 ? (
                <>
                  <div className="overflow-auto">
                    <Table className="min-w-full">
                      <TableHeader>
                        <TableRow>
                          <TableHead>ID</TableHead>
                          <TableHead>Name</TableHead>
                          <TableHead>Email</TableHead>
                          <TableHead>Verified</TableHead>
                          <TableHead>Admin</TableHead>
                          <TableHead>Created at</TableHead>
                          <TableHead></TableHead>
                          <TableHead></TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {users.map((user) => (
                          <TableRow key={user.id}>
                            <TableCell>{user.id}</TableCell>
                            <TableCell>{user.name}</TableCell>
                            <TableCell>{user.email}</TableCell>
                            <TableCell>
                              {user.verified ? "true" : "false"}
                            </TableCell>
                            <TableCell>
                              {user.admin ? "true" : "false"}
                            </TableCell>
                            <TableCell>
                              {new Date(user.created_at)
                                .toISOString()
                                .slice(0, 19)
                                .replace("T", " ")}
                            </TableCell>
                            <TableCell>
                              <Button
                                onClick={() => openEditModal(user)}
                                variant="secondary"
                              >
                                Edit
                              </Button>
                            </TableCell>
                            <TableCell>
                              <Button asChild variant="destructive">
                                <Link href={`/admin/users/${user.id}/delete`}>
                                  Delete
                                </Link>
                              </Button>
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </div>

                  <div className="flex justify-center items-center gap-4 pt-6">
                    <Button
                      variant="ghost"
                      disabled={pagination.page === 1}
                      onClick={() => goToPage(pagination.page - 1)}
                    >
                      Previous
                    </Button>
                    <span className="text-sm">
                      Page {pagination.page} of {pagination.totalPages}
                    </span>
                    <Button
                      variant="ghost"
                      disabled={pagination.page === pagination.totalPages}
                      onClick={() => goToPage(pagination.page + 1)}
                    >
                      Next
                    </Button>
                  </div>
                </>
              ) : (
                <p className="text-sm text-muted-foreground">No users found.</p>
              )}
            </CardContent>
          </Card>
        </div>

        <UserFormModal
          open={showModal}
          onClose={() => setShowModal(false)}
          user={editingUser}
        />
      </div>
    </AppLayout>
  );
}
