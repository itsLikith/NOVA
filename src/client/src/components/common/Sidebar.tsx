import Image from 'next/image';

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarHeader,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenu,
  SidebarGroupLabel,
} from '@/components/ui/sidebar';

import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar';

import { LayoutDashboard, User2, User, LogOut, ChevronsUpDown, Signature } from 'lucide-react';

import { Separator } from '@/components/ui/separator';

const items = [
  {
    title: 'Dashboard',
    url: '/dashboard',
    icon: LayoutDashboard,
  },
  {
    title: 'Boards',
    url: '/boards',
    icon: Signature,
  },
];

const me = [
  {
    title: 'Profile',
    url: '/profile',
    icon: User,
  },
];

export function AppSidebar() {
  return (
    <Sidebar>
      <SidebarHeader>
        <div className="flex items-center ">
          <Image src="/logo.png" height={50} width={50} alt="Logo" />
          <h1 className="scroll-m-20 text-1xl font-extrabold tracking-tight lg:text-2xl">NOVA</h1>
        </div>
      </SidebarHeader>
      <Separator />
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Home</SidebarGroupLabel>
          <SidebarMenu>
            {items.map((item) => (
              <SidebarMenuItem key={item.title}>
                <SidebarMenuButton render={<a href={item.url} />}>
                  <item.icon />
                  <span>{item.title}</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            ))}
          </SidebarMenu>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupLabel>Me</SidebarGroupLabel>
          <SidebarMenu>
            {me.map((m) => (
              <SidebarMenuItem key={m.title}>
                <SidebarMenuButton render={<a href={m.url} />}>
                  <m.icon />
                  <span>{m.title}</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            ))}
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <Separator />
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem className="flex align-center space-between">
            <Avatar>
              <AvatarImage src="https://github.com/shadcn.png" />
              <AvatarFallback>CN</AvatarFallback>
            </Avatar>
            <span>Username</span>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  );
}
