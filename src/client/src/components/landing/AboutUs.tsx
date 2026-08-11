import { Badge } from '@/components/ui/badge';
import { Card, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { MousePointer2, Layout, ShieldCheck, Workflow } from 'lucide-react';

const STATS = [
  { value: '< 50ms', label: 'Real-Time Sync Latency' },
  { value: 'Infinite', label: 'Canvas Workspace Space' },
  { value: '100%', label: 'Internal Data Privacy' },
  { value: '10k+', label: 'Concurrent Canvas Elements' },
] as const;

const CAPABILITIES = [
  {
    icon: MousePointer2,
    title: 'Live Cursor Tracking',
    description:
      'See active collaborators on the board with name tags, live selection highlights, and real-time cursor positions.',
  },
  {
    icon: Layout,
    title: 'Rich Workspace Elements',
    description:
      'Create, edit, move, resize, and group sticky notes, geometric shapes, connectors, and text blocks seamlessly.',
  },
  {
    icon: Workflow,
    title: 'Architecture & Workflow Diagrams',
    description:
      'Purpose-built for system design, team brainstorming, sprint planning, and visual workflow documentation.',
  },
  {
    icon: ShieldCheck,
    title: 'Internal Office Security',
    description:
      'Customized specifically for your organization with SAML/SSO authentication, granular permissions, and data isolation.',
  },
] as const;

function AboutUs() {
  return (
    <section id="about" className="relative flex flex-col items-center px-6 py-20 md:py-28">
      <div className="mx-auto flex max-w-6xl flex-col items-center text-center">
        <Badge variant="outline" className="rounded-full px-3 py-1 text-xs font-medium">
          Platform Overview
        </Badge>

        <h2 className="mt-4 max-w-3xl text-3xl font-extrabold tracking-tight text-foreground sm:text-4xl md:text-5xl">
          Everything your organization needs for visual collaboration
        </h2>

        <p className="mt-4 max-w-2xl text-base leading-relaxed text-muted-foreground md:text-lg">
          NOVA combines the power of an infinite interactive whiteboard with enterprise security.
          Brainstorm ideas, design architecture, and map workflows together in real time.
        </p>

        {/* Stats Grid */}
        <div className="mt-12 grid w-full grid-cols-2 gap-4 sm:grid-cols-4 md:mt-16">
          {STATS.map((stat) => (
            <div
              key={stat.label}
              className="flex flex-col items-center justify-center rounded-2xl border border-border/60 bg-background/60 p-6 backdrop-blur-xs"
            >
              <span className="text-3xl font-extrabold tracking-tight text-foreground md:text-4xl">
                {stat.value}
              </span>
              <span className="mt-1 text-xs font-medium text-muted-foreground md:text-sm">
                {stat.label}
              </span>
            </div>
          ))}
        </div>

        {/* Capabilities Grid */}
        <div
          id="features"
          className="mt-12 grid w-full grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4 md:mt-16"
        >
          {CAPABILITIES.map((item) => {
            const Icon = item.icon;
            return (
              <Card
                key={item.title}
                className="border-border/60 bg-background/70 shadow-xs transition-all duration-200 hover:-translate-y-1 hover:border-border hover:shadow-md"
              >
                <CardHeader className="items-start text-left">
                  <div className="flex size-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
                    <Icon className="size-5" />
                  </div>
                  <CardTitle className="mt-3 text-lg font-bold">{item.title}</CardTitle>
                  <CardDescription className="text-sm leading-relaxed text-muted-foreground">
                    {item.description}
                  </CardDescription>
                </CardHeader>
              </Card>
            );
          })}
        </div>
      </div>
    </section>
  );
}

export { AboutUs };
