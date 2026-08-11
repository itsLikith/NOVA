import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ArrowUpRight, Users, Sparkles } from 'lucide-react';

function HeroSection() {
  return (
    <section className="relative flex flex-1 flex-col items-center justify-center px-6 pt-20 pb-28 text-center md:pt-28 md:pb-36">
      {/* Decorative background glow & gradient vertical lines */}
      <div className="pointer-events-none absolute inset-0 overflow-hidden" aria-hidden="true">
        <div className="absolute top-1/4 left-1/2 -z-10 h-96 w-[600px] -translate-x-1/2 rounded-full bg-cyan-200/30 blur-3xl" />
        <div className="absolute right-[15%] bottom-0 h-48 w-px bg-gradient-to-b from-transparent via-purple-300/30 to-transparent md:right-[20%] md:h-72" />
      </div>

      {/* Product Tag */}
      <Badge
        variant="outline"
        className="mb-6 gap-1.5 rounded-full px-3.5 py-1 text-xs font-medium backdrop-blur-xs"
      >
        <Sparkles className="size-3.5 text-cyan-600" />
        Internal Real-Time Collaboration Platform
      </Badge>

      <h1 className="max-w-4xl text-4xl leading-[1.1] font-extrabold tracking-tight text-foreground sm:text-5xl md:text-6xl lg:text-7xl">
        Real-Time Visual
        <br className="hidden sm:block" /> Collaboration for
        <br className="hidden sm:block" /> Enterprise Teams
      </h1>

      <p className="mt-6 max-w-2xl text-base leading-relaxed text-muted-foreground md:mt-8 md:text-lg">
        NOVA provides a shared virtual canvas where your office teams work together seamlessly.
        Create sticky notes, sketch architecture diagrams, build workflows, and move components in
        real time without refreshing.
      </p>

      <div className="mt-8 flex flex-col items-center gap-3 sm:flex-row sm:gap-4 md:mt-10">
        <Link href="/auth">
          <Button className="h-11 rounded-full px-6 text-base font-medium shadow-sm">
            Launch Workspace
            <ArrowUpRight className="ml-1 size-4" />
          </Button>
        </Link>

        <Link href="#features">
          <Button variant="outline" className="h-11 rounded-full px-6 text-base font-medium">
            Explore Features
            <ArrowUpRight className="ml-1 size-4" />
          </Button>
        </Link>
      </div>

      {/* Sub-hero info badge */}
      <div className="mt-12 flex items-center gap-3 text-xs font-medium text-muted-foreground md:mt-16">
        <div className="flex -space-x-2">
          <div className="flex size-7 items-center justify-center rounded-full border-2 border-background bg-cyan-500 text-[10px] font-bold text-white">
            JD
          </div>
          <div className="flex size-7 items-center justify-center rounded-full border-2 border-background bg-purple-500 text-[10px] font-bold text-white">
            SK
          </div>
          <div className="flex size-7 items-center justify-center rounded-full border-2 border-background bg-amber-500 text-[10px] font-bold text-white">
            AL
          </div>
        </div>
        <span className="inline-flex items-center gap-1">
          <Users className="size-3.5 text-foreground" />
          Multi-user live cursor sync active
        </span>
      </div>
    </section>
  );
}

export { HeroSection };
