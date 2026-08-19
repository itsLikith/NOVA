import { Recents } from '@/components/dashboard/Recents';
import { SharedByYou } from '@/components/dashboard/SharedByYou';
import { SharedWithYou } from '@/components/dashboard/SharedWithYou';
import { Separator } from '@/components/ui/separator';

export default function BoardsPage() {
  return (
    <section className="p-4 flex flex-col gap-5">
      <Recents />
      <Separator />
      <SharedByYou />
      <Separator />
      <SharedWithYou />
    </section>
  );
}
