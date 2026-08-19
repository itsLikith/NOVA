import { Forward } from 'lucide-react';

const SharedByYou = () => {
  return (
    <div className="flex flex-col">
      <span className="h2 flex align-center gap-2">
        <Forward className="" />
        Shared by you
      </span>
      <div>This is for shared by you.</div>
    </div>
  );
};

export { SharedByYou };
