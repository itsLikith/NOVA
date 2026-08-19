import { Clock } from 'lucide-react';

const Recents = () => {
  return (
    <div className="flex flex-col">
      <span className="h2 flex align-center gap-2">
        <Clock className="" />
        Recents
      </span>

      <div>This is for recent boards</div>
    </div>
  );
};

export { Recents };
