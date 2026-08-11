import Link from 'next/link';
import Image from 'next/image';
import { Separator } from '@/components/ui/separator';
import { ArrowUpRight } from 'lucide-react';

const FOOTER_NAV = {
  product: [
    { label: 'Features', href: '#features' },
    { label: 'Components', href: '/docs' },
    { label: 'Design System', href: '/docs' },
    { label: 'Templates', href: '/templates' },
    { label: 'Releases', href: '/releases' },
  ],
  company: [
    { label: 'About Us', href: '#about' },
    { label: 'Careers', href: '/careers' },
    { label: 'Blog', href: '/blog' },
    { label: 'Press Kit', href: '/press' },
    { label: 'Contact', href: '#contact' },
  ],
  resources: [
    { label: 'Documentation', href: '/docs' },
    { label: 'Help Center', href: '/help' },
    { label: 'Community', href: '/community' },
    { label: 'Status Page', href: '/status' },
    { label: 'GitHub Repository', href: 'https://github.com' },
  ],
  legal: [
    { label: 'Privacy Policy', href: '/privacy' },
    { label: 'Terms of Service', href: '/terms' },
    { label: 'Security', href: '/security' },
    { label: 'Cookie Settings', href: '/cookies' },
  ],
};

function Footer() {
  return (
    <footer
      id="contact"
      className="relative border-t border-border/60 bg-background/80 pt-16 pb-12 backdrop-blur-md"
    >
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="grid grid-cols-1 gap-10 lg:grid-cols-12 lg:gap-8">
          {/* Corporate Branding & Logo Area */}
          <div className="flex w-full flex-col items-start space-y-4 lg:col-span-4">
            <Link href="/" className="flex items-center gap-3">
              <div className="w-full max-w-[180px] sm:max-w-[220px]">
                <Image
                  src="/banner.jpg"
                  alt="Bosch Logo"
                  width={220}
                  height={60}
                  className="h-auto w-full object-contain"
                />
              </div>
            </Link>

            <p className="max-w-sm text-sm leading-relaxed text-muted-foreground">
              Invented for life
            </p>
          </div>

          {/* Navigation Columns */}
          <div className="grid grid-cols-2 gap-6 sm:gap-8 lg:grid-cols-4 lg:col-span-8">
            <div>
              <h3 className="text-xs font-semibold uppercase tracking-wider text-foreground">
                Product
              </h3>
              <ul className="mt-4 space-y-2.5">
                {FOOTER_NAV.product.map((item) => (
                  <li key={item.label}>
                    <Link
                      href={item.href}
                      className="text-sm text-muted-foreground transition-colors hover:text-foreground"
                    >
                      {item.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>

            <div>
              <h3 className="text-xs font-semibold uppercase tracking-wider text-foreground">
                Company
              </h3>
              <ul className="mt-4 space-y-2.5">
                {FOOTER_NAV.company.map((item) => (
                  <li key={item.label}>
                    <Link
                      href={item.href}
                      className="text-sm text-muted-foreground transition-colors hover:text-foreground"
                    >
                      {item.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>

            <div>
              <h3 className="text-xs font-semibold uppercase tracking-wider text-foreground">
                Resources
              </h3>
              <ul className="mt-4 space-y-2.5">
                {FOOTER_NAV.resources.map((item) => (
                  <li key={item.label}>
                    <Link
                      href={item.href}
                      className="inline-flex items-center text-sm text-muted-foreground transition-colors hover:text-foreground"
                    >
                      {item.label}
                      {item.href.startsWith('http') && <ArrowUpRight className="ml-0.5 size-3" />}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>

            <div>
              <h3 className="text-xs font-semibold uppercase tracking-wider text-foreground">
                Legal
              </h3>
              <ul className="mt-4 space-y-2.5">
                {FOOTER_NAV.legal.map((item) => (
                  <li key={item.label}>
                    <Link
                      href={item.href}
                      className="text-sm text-muted-foreground transition-colors hover:text-foreground"
                    >
                      {item.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>

        <Separator className="my-10" />

        {/* Bottom Copyright Bar */}
        <div className="flex flex-col items-center justify-between gap-4 text-center sm:flex-row sm:text-left">
          <p className="text-xs text-muted-foreground">
            &copy; {new Date().getFullYear()} Bosch Global Software Technologies. All rights
            reserved.
          </p>
          <div className="flex gap-6 text-xs text-muted-foreground">
            <Link href="/privacy" className="hover:text-foreground">
              Privacy Policy
            </Link>
            <Link href="/terms" className="hover:text-foreground">
              Terms of Service
            </Link>
            <Link href="/security" className="hover:text-foreground">
              Security
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}

export { Footer };
