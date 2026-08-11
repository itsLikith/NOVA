import { Badge } from '@/components/ui/badge';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { Card } from '@/components/ui/card';

const FAQS = [
  {
    id: 'item-1',
    question: 'How does real-time synchronization work on NOVA?',
    answer:
      'NOVA utilizes high-performance WebSocket connections and conflict-free state synchronization. Any action—such as creating sticky notes, moving shapes, or drawing connectors—is instantly broadcast and reflected on all collaborators’ screens in under 50 milliseconds.',
  },
  {
    id: 'item-2',
    question: 'What visual components are supported on the canvas?',
    answer:
      'You can create, edit, resize, and group geometric shapes, sticky notes, rich text blocks, smart connectors/arrows, freehand drawings, workflow cards, and system architecture components.',
  },
  {
    id: 'item-3',
    question: 'How many team members can work on a board at the same time?',
    answer:
      'NOVA is engineered for internal office scale, supporting dozens to hundreds of concurrent users on a single board with live cursor indicators and active element selection overlays.',
  },
  {
    id: 'item-4',
    question: 'Is NOVA secure for internal enterprise discussions?',
    answer:
      'Yes. NOVA is designed for internal organizational deployment. It integrates directly with your corporate SAML/Single Sign-On (SSO), supports role-based access control (RBAC), and keeps all board data hosted within your private network.',
  },
  {
    id: 'item-5',
    question: 'Can we export boards or save reusable templates?',
    answer:
      'Abolutely. Workspaces can be exported to high-resolution PNG, SVG, PDF, or JSON board backups. You can also save custom board configurations as reusable templates for future team sprints or architecture reviews.',
  },
  {
    id: 'item-6',
    question: 'Can NOVA be customized for our organization’s workflows?',
    answer:
      'Because NOVA is built for internal office deployment, its component library, board permissions, workflow templates, and user access rules can be tailored to align with your organization’s governance.',
  },
] as const;

function Faq() {
  return (
    <section id="faq" className="relative flex flex-col items-center px-6 py-20 md:py-28">
      <div className="mx-auto flex w-full max-w-4xl flex-col items-center text-center">
        <Badge variant="outline" className="rounded-full px-3 py-1 text-xs font-medium">
          FAQ
        </Badge>

        <h2 className="mt-4 max-w-2xl text-3xl font-extrabold tracking-tight text-foreground sm:text-4xl md:text-5xl">
          Frequently Asked Questions
        </h2>

        <p className="mt-4 max-w-xl text-base leading-relaxed text-muted-foreground md:text-lg">
          Got questions about NOVA? Everything you need to know about our internal real-time
          collaborative whiteboard.
        </p>

        {/* Accordion List Container */}
        <Card className="mt-12 w-full border-border/60 bg-background/80 p-4 shadow-sm backdrop-blur-md sm:p-6 md:mt-16">
          <Accordion defaultValue={['item-1']} className="w-full">
            {FAQS.map((faq) => (
              <AccordionItem key={faq.id} value={faq.id} className="border-border/60 px-2 py-1">
                <AccordionTrigger className="text-left text-base font-semibold text-foreground hover:no-underline hover:text-foreground/80 md:text-lg">
                  {faq.question}
                </AccordionTrigger>
                <AccordionContent className="text-left text-sm leading-relaxed text-muted-foreground md:text-base">
                  {faq.answer}
                </AccordionContent>
              </AccordionItem>
            ))}
          </Accordion>
        </Card>
      </div>
    </section>
  );
}

export { Faq };
