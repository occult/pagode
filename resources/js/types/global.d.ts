import { User } from ".";

export type FlashMessages = {
  success?: string[];
  info?: string[];
  warning?: string[];
  danger?: string[];
};

export type SharedProps = {
  flash: FlashMessages;
  auth: {
    user: User | null;
  };
};
