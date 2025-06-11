import { Head, useForm } from "@inertiajs/react";
import { FormEventHandler, useRef } from "react";
import { LoaderCircle, UploadCloud } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import InputError from "@/components/InputError";
import AppLayout from "@/Layouts/AppLayout";
import { BreadcrumbItem } from "@/types";

type UploadedFile = {
  id: string;
  name: string;
  size: string;
  modified: string;
};

type Props = {
  files: UploadedFile[];
};

const breadcrumbs: BreadcrumbItem[] = [
  {
    title: "Upload Files",
    href: "/files",
  },
];

export default function UploadFile({ files }: Props) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const { setData, post, processing, errors, reset } = useForm<{
    file: File | null;
  }>({
    file: null,
  });

  const submit: FormEventHandler = (e) => {
    e.preventDefault();
    post("/files", {
      forceFormData: true,
      onSuccess: () => {
        reset();
        if (fileInputRef.current) {
          fileInputRef.current.value = "";
        }
      },
    });
  };

  return (
    <AppLayout breadcrumbs={breadcrumbs}>
      <Head title="Upload File" />

      <div className="flex flex-1 flex-col items-center justify-center p-6">
        <div className="w-full max-w-xl space-y-10">
          <form onSubmit={submit} className="grid gap-6">
            <div className="grid gap-2">
              <Label htmlFor="file">Choose File</Label>
              <Input
                ref={fileInputRef}
                id="file"
                type="file"
                accept="*"
                onChange={(e) => setData("file", e.target.files?.[0] || null)}
              />
              <InputError message={errors.file} />
            </div>

            <Button
              type="submit"
              disabled={processing}
              className="flex items-center gap-2"
            >
              {processing && <LoaderCircle className="h-4 w-4 animate-spin" />}
              <UploadCloud className="h-4 w-4" />
              Upload
            </Button>
          </form>

          <div>
            <h3 className="text-lg font-semibold mb-2">Uploaded Files</h3>
            {files.length > 0 ? (
              <div className="border rounded-lg p-4">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b text-left">
                      <th className="py-2 pr-4">Filename</th>
                      <th className="py-2 pr-4">Size</th>
                      <th className="py-2">Modified</th>
                    </tr>
                  </thead>
                  <tbody>
                    {files.map((file) => (
                      <tr key={file.id} className="border-b">
                        <td className="py-2 pr-4">{file.name}</td>
                        <td className="py-2 pr-4">{file.size}</td>
                        <td className="py-2">{file.modified}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">
                No files uploaded yet.
              </p>
            )}
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
