type ServiceIconProps = {
  icon: string;
  className?: string;
};

const iconMap: Record<string, string> = {
  nextjs: 'nextdotjs',
  vite: 'vite',
  astro: 'astro',
  nuxtjs: 'nuxtdotjs',
  nestjs: 'nestjs',
  java: 'java',
  django: 'django',
  fastapi: 'fastapi',
  python: 'python',
  go: 'go',
  nodejs: 'nodedotjs',
  express: 'express',
  react: 'react',
  vue: 'vuedotjs',
  rust: 'rust',
  php: 'php',
  ruby: 'ruby',
  git: 'git',
};

export function ServiceIcon({ icon, className = 'w-6 h-6' }: ServiceIconProps) {
  const mappedIcon = iconMap[icon] || 'git';

  return (
    <div
      className={`flex shrink-0 items-center justify-center rounded bg-white p-0.5 dark:bg-zinc-100 ${className}`}
    >
      <img
        src={`https://cdn.simpleicons.org/${mappedIcon}`}
        alt={icon}
        className="h-full w-full object-contain"
        onError={(e) => {
          e.currentTarget.src = 'https://cdn.simpleicons.org/git';
        }}
      />
    </div>
  );
}
