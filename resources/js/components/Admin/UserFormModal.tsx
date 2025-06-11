import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useForm } from "@inertiajs/react";
import { User } from "@/types";

type UserFormData = {
  name: string;
  email: string;
  admin: boolean;
  emailVerified: boolean;
};

type Props = {
  open: boolean;
  onClose: () => void;
  user: User | null;
};

export function UserFormModal({ open, onClose, user }: Props) {
  const isEdit = !!user;

  const { data, setData, post, processing, errors } = useForm<UserFormData>({
    name: user?.name ?? "",
    email: user?.email ?? "",
    admin: !!user?.admin,
    emailVerified: !!user?.emailVerified,
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (isEdit) {
      post(`/admin/users/${user!.id}/edit`, {
        forceFormData: true,
        preserveScroll: true,
        onSuccess: onClose,
      });
    } else {
      post("/admin/users/add", {
        forceFormData: true,
        preserveScroll: true,
        onSuccess: onClose,
      });
    }
  };

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEdit ? "Edit User" : "Add User"}</DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid gap-2">
            <Label htmlFor="name">Name</Label>
            <Input
              id="name"
              value={data.name}
              onChange={(e) => setData("name", e.target.value)}
            />
            {errors.name && (
              <p className="text-sm text-destructive">{errors.name}</p>
            )}
          </div>

          <div className="grid gap-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              value={data.email}
              onChange={(e) => setData("email", e.target.value)}
            />
            {errors.email && (
              <p className="text-sm text-destructive">{errors.email}</p>
            )}
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="grid gap-2">
              <Label>Admin</Label>
              <Button
                type="button"
                variant={data.admin ? "default" : "outline"}
                onClick={() => setData("admin", !data.admin)}
              >
                {data.admin ? "Yes (Click to unset)" : "No (Click to set)"}
              </Button>
            </div>

            <div className="grid gap-2">
              <Label>Verified</Label>
              <Button
                type="button"
                variant={data.emailVerified ? "default" : "outline"}
                onClick={() => setData("emailVerified", !data.emailVerified)}
              >
                {data.emailVerified
                  ? "Yes (Click to unset)"
                  : "No (Click to set)"}
              </Button>
            </div>
          </div>

          <div className="flex justify-end pt-4">
            <Button type="submit" disabled={processing}>
              {isEdit ? "Update" : "Create"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
