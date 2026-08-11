import { NavBar } from '@/components/landing/NavBar';
import { HeroSection } from '@/components/landing/HeroSection';
import { AboutUs } from '@/components/landing/AboutUs';
import { Faq } from '@/components/landing/Faq';
import { Footer } from '@/components/landing/Footer';

export default function Home() {
  return (
    <div className="relative flex min-h-svh flex-col overflow-hidden bg-background">
      {/* Background gradient — white to light cyan at the bottom */}
      <div
        className="pointer-events-none absolute inset-0 bg-gradient-to-b from-background via-background to-cyan-100/60"
        aria-hidden="true"
      />

      <div className="relative z-10 flex flex-1 flex-col">
        <NavBar />
        <HeroSection />
        <AboutUs />
        <Faq />
        <Footer />
      </div>
    </div>
  );
}
