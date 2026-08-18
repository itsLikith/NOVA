import { NavBar } from '@/components/landing/NavBar';
import { LoginForm } from '@/components/auth/LoginForm';

export default function LoginPage() {
  return (
    <div className="relative flex min-h-svh flex-col overflow-hidden bg-background">
      <div
        className="pointer-events-none absolute inset-0 bg-gradient-to-b from-background via-background to-cyan-100/60"
        aria-hidden="true"
      />

      <div className="relative z-10 flex flex-1 flex-col">
        <NavBar />

        <main className="flex flex-1 items-center justify-center px-4 py-10 sm:px-6 lg:px-8">
          <div className="w-full max-w-md rounded-[28px] border border-border/60 bg-background/80 p-5 shadow-sm backdrop-blur-md sm:p-7">
            <LoginForm />
          </div>
        </main>
      </div>
    </div>
  );
}
