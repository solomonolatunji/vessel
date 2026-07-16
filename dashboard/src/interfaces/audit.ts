export interface AuditLog {
  id: string;
  userId: string;
  action: string;
  resource: string;
  details: string;
  ipAddress: string;
  createdAt: string;
}
