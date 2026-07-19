import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import type { ListTablesResponse } from '#/interfaces/database';
import { apiClient } from '#/lib/apiClient';

interface Props {
  databaseId: string;
}

export function DataBrowser({ databaseId }: Props) {
  const [selectedTable, setSelectedTable] = useState<string>('');

  const { data: schemasData, isLoading: schemasLoading } = useQuery({
    queryKey: ['databases', databaseId, 'schemas'],
    queryFn: () => apiClient.get<ListTablesResponse>(`/databases/${databaseId}/schemas`),
  });

  const { data: tableData, isLoading: tableLoading } = useQuery({
    queryKey: ['databases', databaseId, 'data', selectedTable],
    queryFn: () =>
      apiClient.get<{ data: Record<string, unknown>[] }>(
        `/databases/${databaseId}/data/${selectedTable}`
      ),
    enabled: !!selectedTable,
  });

  const tables = schemasData?.data || [];
  const selectedSchema = tables.find((t) => t.name === selectedTable);
  const rows = tableData?.data || [];

  return (
    <Card className="flex h-full min-h-[500px] flex-col">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
        <div className="space-y-1">
          <CardTitle>Data Browser</CardTitle>
          <CardDescription>View and manage records in your database tables.</CardDescription>
        </div>
        <div className="w-[250px]">
          <Select
            value={selectedTable}
            onValueChange={setSelectedTable}
            disabled={schemasLoading || tables.length === 0}
          >
            <SelectTrigger>
              <SelectValue
                placeholder={
                  schemasLoading
                    ? 'Loading tables...'
                    : tables.length === 0
                      ? 'No tables found'
                      : 'Select a table'
                }
              />
            </SelectTrigger>
            <SelectContent>
              {tables.map((table) => (
                <SelectItem key={table.name} value={table.name}>
                  {table.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </CardHeader>
      <CardContent className="flex-1 overflow-auto">
        {!selectedTable ? (
          <div className="flex h-full items-center justify-center rounded-md border border-dashed p-8 text-muted-foreground">
            Select a table to view its data
          </div>
        ) : tableLoading ? (
          <div className="flex h-full animate-pulse items-center justify-center p-8 text-muted-foreground">
            Loading table data...
          </div>
        ) : rows.length === 0 ? (
          <div className="flex flex-col items-center justify-center rounded-md border border-dashed p-8 text-muted-foreground">
            <p>
              No records found in table <strong>{selectedTable}</strong>
            </p>
          </div>
        ) : (
          <div className="overflow-x-auto rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  {selectedSchema?.columns.map((col) => (
                    <TableHead key={col.name} className="whitespace-nowrap">
                      {col.name}
                      <span className="ml-2 font-normal text-muted-foreground text-xs">
                        {col.type}
                      </span>
                    </TableHead>
                  ))}
                  {/* Fallback if schema doesn't match data exactly */}
                  {!selectedSchema &&
                    Object.keys(rows[0] || {}).map((key) => <TableHead key={key}>{key}</TableHead>)}
                </TableRow>
              </TableHeader>
              <TableBody>
                {rows.map((row, i) => (
                  <TableRow key={i}>
                    {selectedSchema
                      ? selectedSchema.columns.map((col) => (
                          <TableCell key={col.name} className="max-w-[300px] truncate">
                            {String(row[col.name] ?? '')}
                          </TableCell>
                        ))
                      : Object.values(row).map((val, j) => (
                          <TableCell key={j} className="max-w-[300px] truncate">
                            {String(val ?? '')}
                          </TableCell>
                        ))}
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
