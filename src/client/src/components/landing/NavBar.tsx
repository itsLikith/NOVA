'use client';

import Link from 'next/link';
import Image from 'next/image';

import { Button } from '@/components/ui/button';
import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
} from '@/components/ui/navigation-menu';
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from '@/components/ui/sheet';
import { Separator } from '@/components/ui/separator';
import { ChevronRight, Menu } from 'lucide-react';

const NAV_LINKS = [
  { label: 'About', href: '#about' },
  { label: 'FAQs', href: '#faq' },
  { label: 'Docs', href: '/docs' },
] as const;

const FEATURE_ITEMS = [
  {
    title: 'Real-Time Sync',
    description: 'Sub-millisecond state updates & live cursor tracking for seamless collaboration.',
    href: '#features',
  },
  {
    title: 'Infinite Workspace',
    description: 'Shapes, sticky notes, connectors, and diagrams tailored for team planning.',
    href: '#features',
  },
  {
    title: 'Internal & Secure',
    description: 'Enterprise SSO, role-based access control, and end-to-end office data privacy.',
    href: '#features',
  },
] as const;

function NavBar() {
  return (
    <header className="sticky top-0 z-50 w-full px-4 pt-4 md:px-6 md:pt-6">
      {/* Desktop Navigation */}
      <nav className="mx-auto hidden max-w-4xl items-center justify-between rounded-full border border-border/60 bg-background/80 px-4 py-2 shadow-sm backdrop-blur-md md:flex">
        <Link href="/" className="flex shrink-0 items-center gap-2">
          <Image src="/logo.png" alt="NOVA Logo" width={32} height={32} priority />
          <span className="text-base font-bold tracking-tight text-foreground">NOVA</span>
        </Link>

        <NavigationMenu className="flex-1 justify-center">
          <NavigationMenuList>
            <NavigationMenuItem>
              <NavigationMenuTrigger className="bg-transparent text-sm font-normal text-foreground/80 hover:bg-transparent hover:text-foreground data-popup-open:bg-transparent">
                Features
              </NavigationMenuTrigger>
              <NavigationMenuContent className="w-[320px] p-3">
                <ul className="flex flex-col gap-1">
                  {FEATURE_ITEMS.map((item) => (
                    <li key={item.title}>
                      <NavigationMenuLink
                        href={item.href}
                        className="flex flex-col gap-0.5 rounded-md p-2.5 transition-colors hover:bg-muted"
                      >
                        <span className="text-sm font-medium text-foreground">{item.title}</span>
                        <span className="text-xs text-muted-foreground">{item.description}</span>
                      </NavigationMenuLink>
                    </li>
                  ))}
                </ul>
              </NavigationMenuContent>
            </NavigationMenuItem>

            {NAV_LINKS.map((link) => (
              <NavigationMenuItem key={link.label}>
                <NavigationMenuLink
                  href={link.href}
                  className="bg-transparent px-3 py-2 text-sm font-normal text-foreground/80 hover:bg-transparent hover:text-foreground"
                >
                  {link.label}
                </NavigationMenuLink>
              </NavigationMenuItem>
            ))}
          </NavigationMenuList>
        </NavigationMenu>

        <Link href="/auth">
          <Button variant="outline" className="rounded-full px-5 text-sm font-medium">
            Join Board
          </Button>
        </Link>
      </nav>

      {/* Mobile Navigation */}
      <nav className="mx-auto flex max-w-lg items-center justify-between rounded-full border border-border/60 bg-background/80 px-4 py-2.5 shadow-sm backdrop-blur-md md:hidden">
        <Link href="/" className="flex shrink-0 items-center gap-2">
          <Image src="/logo.png" alt="NOVA Logo" width={28} height={28} priority />
          <span className="text-base font-bold tracking-tight text-foreground">NOVA</span>
        </Link>

        <Sheet>
          <SheetTrigger
            render={<Button variant="ghost" size="icon" aria-label="Open navigation menu" />}
          >
            <Menu className="size-5" />
          </SheetTrigger>

          <SheetContent
            side="top"
            className="rounded-b-2xl border-b border-border/60 p-0"
            showCloseButton
          >
            <SheetHeader className="flex-row items-center justify-between border-b-0 px-6 pt-5 pb-3">
              <SheetTitle className="sr-only">Navigation Menu</SheetTitle>
              <Link href="/" className="flex items-center gap-2">
                <Image src="/logo.png" alt="NOVA Logo" width={28} height={28} />
                <span className="text-base font-bold tracking-tight text-foreground">NOVA</span>
              </Link>
            </SheetHeader>

            <div className="mx-4 mb-6 flex flex-col rounded-xl border border-border/60 bg-background">
              <Link
                href="#features"
                className="flex items-center justify-between px-5 py-4 text-base font-medium text-foreground transition-colors hover:bg-muted/50"
              >
                Features
                <ChevronRight className="size-4 text-muted-foreground" />
              </Link>
              <Separator />

              {NAV_LINKS.map((link, index) => (
                <div key={link.label}>
                  <Link
                    href={link.href}
                    className="flex items-center px-5 py-4 text-base font-medium text-foreground transition-colors hover:bg-muted/50"
                  >
                    {link.label}
                  </Link>
                  {index < NAV_LINKS.length - 1 && <Separator />}
                </div>
              ))}
            </div>
          </SheetContent>
        </Sheet>
      </nav>
    </header>
  );
}

export { NavBar };
