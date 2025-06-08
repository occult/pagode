import { Head } from "@inertiajs/react";

interface ErrorPageProps {
  status: number;
  message?: string;
}

const defaultMessages: Record<number, string> = {
  403: "Forbidden",
  404: "Not Found",
  419: "Page Expired",
  429: "Too Many Requests",
  500: "Server Error",
};

export default function ErrorPage({ status, message }: ErrorPageProps) {
  const title = message || defaultMessages[status] || "Something went wrong.";

  return (
    <>
      <Head title={`${status} ${title}`} />
      <div className="min-h-screen flex flex-col items-center justify-center bg-background text-foreground text-center px-6">
        <h1 className="text-9xl font-bold tracking-tight text-primary">
          {status}
        </h1>
        <p className="mt-6 text-2xl font-semibold">{title}</p>
      </div>
    </>
  );
}
