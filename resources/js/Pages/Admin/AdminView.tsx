import { Head, Link } from "@inertiajs/react";
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
};

const breadcrumbs: BreadcrumbItem[] = [
  {
    title: "Admin",
    href: "/admin",
  },
];

export default function AdminView({ users }: Props) {
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

  return (
    <AppLayout breadcrumbs={breadcrumbs}>
      <Head title="User" />

      <div className="flex flex-col items-center justify-start p-6">
        <div className="w-full max-w-6xl space-y-6 flex-1">
          <Card className="h-full flex flex-col">
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>User</CardTitle>
              <Button onClick={openAddModal}>Add User</Button>
            </CardHeader>

            <CardContent className="overflow-auto max-h-full">
              {users.length > 0 ? (
                <div>
                  <Table>
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
                          <TableCell>{user.admin ? "true" : "false"}</TableCell>
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
              ) : (
                <p className="text-sm text-muted-foreground">No users found.</p>
              )}

              <div className="pt-6">
                <Button disabled variant="ghost">
                  Previous page
                </Button>
              </div>
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
