import { Head, useForm } from "@inertiajs/react";
import { LoaderCircle } from "lucide-react";
import { FormEventHandler } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import InputError from "@/components/InputError";
import TextLink from "@/components/TextLink";
import AuthLayout from "@/Layouts/AuthLayout";

type RegisterForm = {
  name: string;
  email: string;
  password: string;
  password_confirmation: string;
};

export default function Register() {
  const { data, setData, post, processing, errors } = useForm<
    Required<RegisterForm>
  >({
    name: "",
    email: "",
    password: "",
    password_confirmation: "",
  });

  const submit: FormEventHandler = (e) => {
    e.preventDefault();
    post("/user/register", {
      forceFormData: true,
    });
  };

  return (
    <AuthLayout
      title="Create an account"
      description="Crafting seamless digital experiences with InertiaJS, React, and Golang"
      logo="Pagode"
    >
      <Head title="Register" />

      <form className="flex flex-col gap-6 w-full max-w-md" onSubmit={submit}>
        <div className="grid gap-6">
          <div className="grid gap-2">
            <Label htmlFor="name">Name</Label>
            <Input
              id="name"
              type="text"
              required
              autoFocus
              autoComplete="name"
              value={data.name}
              onChange={(e) => setData("name", e.target.value)}
              placeholder="Your name"
            />
            <InputError message={errors.name} />
          </div>

          <div className="grid gap-2">
            <Label htmlFor="email">Email address</Label>
            <Input
              id="email"
              type="email"
              required
              autoComplete="email"
              value={data.email}
              onChange={(e) => setData("email", e.target.value)}
              placeholder="email@example.com"
            />
            <InputError message={errors.email} />
          </div>

          <div className="grid gap-2">
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              type="password"
              required
              autoComplete="new-password"
              value={data.password}
              onChange={(e) => setData("password", e.target.value)}
              placeholder="Password"
            />
            <InputError message={errors.password} />
          </div>

          <div className="grid gap-2">
            <Label htmlFor="password-confirm">Confirm Password</Label>
            <Input
              id="password-confirm"
              type="password"
              required
              autoComplete="new-password"
              value={data.password_confirmation}
              onChange={(e) => setData("password_confirmation", e.target.value)}
              placeholder="Confirm your password"
            />
            <InputError message={errors["password_confirmation"]} />
          </div>

          <Button type="submit" className="mt-4 w-full" disabled={processing}>
            {processing && <LoaderCircle className="h-4 w-4 animate-spin" />}
            Create Account
          </Button>
        </div>

        <div className="text-center text-sm text-muted-foreground">
          Already have an account?{" "}
          <TextLink href="/user/login" tabIndex={6}>
            Log in
          </TextLink>
        </div>
      </form>
    </AuthLayout>
  );
}
