import { Head, router, useForm } from "@inertiajs/react";
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

type UserFormData = {
  name: string;
  email: string;
  admin: boolean;
  emailVerified: boolean;
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

  const { data, setData, post, processing, errors, reset } =
    useForm<UserFormData>({
      name: "",
      email: "",
      admin: false,
      emailVerified: false,
    });

  const openAddModal = () => {
    setEditingUser(null);
    setData({
      name: "",
      email: "",
      admin: false,
      emailVerified: false,
    });
    setShowModal(true);
  };

  const openEditModal = (user: User) => {
    setEditingUser(user);
    setData({
      name: user.name ?? "",
      email: user.email ?? "",
      admin: !!user.admin,
      emailVerified: !!user.emailVerified,
    });
    setShowModal(true);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const url = editingUser
      ? `/admin/users/${editingUser.id}/edit`
      : "/admin/users/add";

    post(url, {
      forceFormData: true,
      preserveScroll: true,
      onSuccess: () => {
        setShowModal(false);
        reset();
      },
    });
  };

  const goToPage = (page: number) => {
    router.visit(`/admin/users?page=${page}`, {
      preserveScroll: true,
      preserveState: true,
    });
  };

  const handleDeleteUser = (userId: number) => {
    if (!confirm("Are you sure you want to delete this user?")) return;

    router.post(
      `/admin/users/${userId}/delete`,
      {},
      {
        preserveScroll: true,
        preserveState: true,
      },
    );
  };

  return (
    <AppLayout breadcrumbs={breadcrumbs}>
      <Head title="User" />

      <div className="flex flex-col w-full h-full p-6">
        <div className="w-full flex-1 min-h-0">
          <Card className="h-full flex flex-col min-h-0">
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>User</CardTitle>
              <Button onClick={openAddModal}>Add User</Button>
            </CardHeader>

            <CardContent className="flex flex-col gap-4 h-full min-h-0">
              {users.length > 0 ? (
                <>
                  <div className="flex-1 min-h-0 overflow-auto">
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
                              <Button
                                variant="destructive"
                                onClick={() => handleDeleteUser(user.id)}
                              >
                                Delete
                              </Button>
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </div>

                  <div className="flex justify-center items-center gap-4">
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
          isEdit={!!editingUser}
          handleSubmit={handleSubmit}
          data={data}
          setData={setData}
          errors={errors}
          processing={processing}
        />
      </div>
    </AppLayout>
  );
}
