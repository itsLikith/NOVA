import { CornerUpLeft } from 'lucide-react';

const SharedWithYou = () => {
  return (
    <div className="flex flex-col">
      <span className="h2 flex align-center gap-2">
        <CornerUpLeft className="" />
        Shared by you
      </span>
      <div>This is for shared with you.</div>
    </div>
  );
};

export { SharedWithYou };
