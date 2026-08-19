import { SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar';
import { AppSidebar } from '@/components/common/Sidebar';
import { Separator } from '@/components/ui/separator';
import Image from 'next/image';

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <SidebarProvider>
      <AppSidebar />
      <main className="w-full">
        <div className="flex justify-between align-center p-4">
          <SidebarTrigger />
          <Image src="/banner.jpg" width={110} height={30} alt="Icon" />
        </div>
        <Separator />
        {children}
      </main>
    </SidebarProvider>
  );
}
