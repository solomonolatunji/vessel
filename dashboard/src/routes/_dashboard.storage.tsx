import { createFileRoute } from '@tanstack/react-router';
import { S3DestinationsList } from '#/features/instance/s3-destinations-list';

export const Route = createFileRoute('/_dashboard/storage')({
  component: () => <S3DestinationsList />,
});
