#!/bin/bash
set -e

# Create new features folders
mkdir -p src/features/profile

# Touch new features files
touch src/features/profile/user-profile-form.tsx
touch src/features/profile/security-2fa-setup.tsx
touch src/features/profile/access-tokens-list.tsx
touch src/features/projects/compose-deploy-form.tsx
touch src/features/projects/jobs-list.tsx
touch src/features/instance/oauth-providers-list.tsx
touch src/features/instance/git-apps-manager.tsx
touch src/features/instance/s3-destinations-list.tsx

# Touch new routes files
mkdir -p src/routes/_shell/profile
touch src/routes/_shell/profile/index.tsx
touch src/routes/_shell/settings/oauth.tsx
touch src/routes/_shell/settings/git-apps.tsx
touch src/routes/_shell/settings/backups.tsx
touch src/routes/_shell/projects/\$projectId/jobs.tsx
touch src/routes/_shell/projects/\$projectId/compose.tsx

# Initialize new route files so generate-routes works
for f in $(find src/routes -type f -name "*.tsx" -size 0); do
  echo "import { createFileRoute } from '@tanstack/react-router';" > "$f"
  echo "" >> "$f"
  echo "export const Route = createFileRoute('${f#src/routes}')({" >> "$f"
  echo "  component: () => <div>Route Component</div>," >> "$f"
  echo "});" >> "$f"
done
