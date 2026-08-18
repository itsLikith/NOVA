'use client';

import Image from 'next/image';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { FormEvent, useState } from 'react';

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

type LoginResponse = {
  status: number;
  message: string;
  data?: {
    token: string;
    user: {
      userid: string;
      email: string;
      role: string;
    };
  };
};

type ResponseMessage = {
  text: string;
  type: 'error' | 'success';
};

export function LoginForm({ className, ...props }: React.ComponentProps<'div'>) {
  const router = useRouter();
  const [userID, setUserID] = useState('');
  const [password, setPassword] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [responseMessage, setResponseMessage] = useState<ResponseMessage | null>(null);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setResponseMessage(null);
    setIsSubmitting(true);

    try {
      const response = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ userid: userID, password }),
      });

      const payload = (await response.json().catch(() => null)) as LoginResponse | null;

      if (!response.ok || !payload?.data?.token || !payload.data.user) {
        setResponseMessage({
          type: 'error',
          text:
            payload?.message === 'invalid credentials'
              ? 'That user ID or password is incorrect. Please try again.'
              : payload?.message ?? 'We could not sign you in. Please try again shortly.',
        });
        return;
      }

      localStorage.setItem('nova.auth.token', payload.data.token);
      localStorage.setItem('nova.auth.user', JSON.stringify(payload.data.user));
      setResponseMessage({
        type: 'success',
        text: `Welcome back, ${payload.data.user.userid}. Redirecting to your dashboard…`,
      });
      router.push('/dashboard');
    } catch {
      setResponseMessage({
        type: 'error',
        text: 'Unable to reach Nova right now. Check your connection and try again.',
      });
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <div className={cn('flex flex-col gap-6', className)} {...props}>
      <form onSubmit={handleSubmit}>
        <FieldGroup>
          <div className="flex flex-col items-center gap-2 text-center">
            <Link href="/" className="flex flex-col items-center gap-2 font-medium">
              <div className="flex items-center justify-center rounded-md">
                <Image src="/logo.png" alt="Nova Logo" width={50} height={50} />
              </div>
              <span className="sr-only">Nova.</span>
            </Link>
            <h1 className="text-xl font-bold">Welcome to Nova.</h1>
          </div>
          <Field>
            <FieldLabel htmlFor="userid">User ID</FieldLabel>
            <Input
              id="userid"
              name="userid"
              value={userID}
              onChange={(event) => setUserID(event.target.value)}
              placeholder="Your user ID"
              autoComplete="username"
              disabled={isSubmitting}
              required
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="password">Password</FieldLabel>
            <Input
              id="password"
              name="password"
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              placeholder="••••••••"
              autoComplete="current-password"
              disabled={isSubmitting}
              required
            />
          </Field>
          {responseMessage && (
            <p
              className={cn(
                'rounded-md px-3 py-2 text-sm',
                responseMessage.type === 'success'
                  ? 'bg-emerald-500/10 text-emerald-700 dark:text-emerald-400'
                  : 'bg-destructive/10 text-destructive'
              )}
              role={responseMessage.type === 'error' ? 'alert' : 'status'}
              aria-live="polite"
            >
              {responseMessage.text}
            </p>
          )}
          <Field>
            <Button type="submit" disabled={isSubmitting} className="w-full">
              {isSubmitting ? 'Signing in…' : 'Login'}
            </Button>
          </Field>
          <FieldSeparator>Or</FieldSeparator>
          <Field className="grid">
            <Button variant="outline" disabled>
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
