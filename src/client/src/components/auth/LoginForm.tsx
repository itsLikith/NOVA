import Image from 'next/image';

import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
  FieldSeparator,
} from '@/components/ui/field';
import { Input } from '@/components/ui/input';

export function LoginForm({ className, ...props }: React.ComponentProps<'div'>) {
  return (
    <div className={cn('flex flex-col gap-6', className)} {...props}>
      <form>
        <FieldGroup>
          <div className="flex flex-col items-center gap-2 text-center">
            <a href="#" className="flex flex-col items-center gap-2 font-medium">
              <div className="flex items-center justify-center rounded-md">
                <Image src="/logo.png" alt="Nova Logo" width={50} height={50} />
              </div>
              <span className="sr-only">Nova.</span>
            </a>
            <h1 className="text-xl font-bold">Welcome to Nova.</h1>
          </div>
          <Field>
            <FieldLabel htmlFor="email">Email</FieldLabel>
            <Input id="email" type="email" placeholder="example@in.bosch.com" required />
          </Field>
          <Field>
            <FieldLabel htmlFor="password">Password</FieldLabel>
            <Input id="password" type="password" placeholder="•••••••" required />
          </Field>
          <Field>
            <Button type="submit">Login</Button>
          </Field>
          <FieldSeparator>Or</FieldSeparator>
          <Field className="grid">
            <Button variant="outline" className="disabled :cursor-not-allowed disabled:opacity-50">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 100 100"
                width="100"
                height="100"
              >
                <circle cx="50" cy="50" r="45" fill="none" stroke="#000" strokeWidth="5" />

                <path d="M31 41.5H69M31 58.5H69" fill="none" stroke="#000" strokeWidth="5" />

                <path
                  d="M31 27C21.5 32.5 17 40.5 17 50C17 59.5 21.5 67.5 31 73V27Z"
                  fill="#fff"
                  stroke="#000"
                  strokeWidth="5"
                  strokeLinejoin="miter"
                />

                <path
                  d="M69 27C78.5 32.5 83 40.5 83 50C83 59.5 78.5 67.5 69 73V27Z"
                  fill="#fff"
                  stroke="#000"
                  strokeWidth="5"
                  strokeLinejoin="miter"
                />
              </svg>
              Continue with Bosch
            </Button>
          </Field>
        </FieldGroup>
      </form>
      <FieldDescription className="px-6 text-center">
        By clicking continue, you agree to our <a href="#">Terms of Service</a> and{' '}
        <a href="#">Privacy Policy</a>.
      </FieldDescription>
    </div>
  );
}
