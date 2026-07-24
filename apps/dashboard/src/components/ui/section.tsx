import { Label } from '#/components/ui/label';

export type SectionProps = {
  icon: React.ReactNode;
  title: string;
  action?: React.ReactNode;
  children: React.ReactNode;
};

export const Section = ({ icon, title, action, children }: SectionProps) => (
  <div className="rounded-xl border border-border/50 bg-card/40 p-6">
    <div className="mb-4 flex items-center justify-between">
      <div className="flex items-center gap-3">
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
          {icon}
        </div>
        <span className="font-semibold text-sm">{title}</span>
      </div>
      {action && <div className="flex shrink-0">{action}</div>}
    </div>
    <div className="divide-y divide-border/50">{children}</div>
  </div>
);

export type RowProps = { label: string; description?: string; children: React.ReactNode };
export const Row = ({ label, description, children }: RowProps) => (
  <div className="flex flex-col gap-4 py-4 md:flex-row md:items-center md:justify-between">
    <div className="flex-1 pr-4">
      <Label className="font-medium text-sm">{label}</Label>
      {description && <p className="mt-1 text-muted-foreground text-sm">{description}</p>}
    </div>
    <div className="flex w-full shrink-0 md:w-1/2 md:justify-end">{children}</div>
  </div>
);
