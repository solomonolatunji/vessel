import MonacoEditor from '@monaco-editor/react';
import { useMutation } from '@tanstack/react-query';
import { flexRender, getCoreRowModel, useReactTable } from '@tanstack/react-table';
import { Loader2, Play } from 'lucide-react';
import { useMemo, useState } from 'react';
import { toast } from 'sonner';

import { Button } from '#/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import { apiClient } from '#/lib/apiClient';

interface SqlStudioProps {
  databaseId: string;
}

export function SqlStudio({ databaseId }: SqlStudioProps) {
  const [query, setQuery] = useState('SELECT * FROM users LIMIT 10;');
  const [results, setResults] = useState<Record<string, unknown>[] | null>(null);
  const [columns, setColumns] = useState<string[]>([]);

  const executeQuery = useMutation({
    mutationFn: async (sql: string) => {
      return apiClient.post<{ rows: Record<string, unknown>[]; columns: string[] }>(
        `/databases/${databaseId}/query`,
        {
          query: sql,
        }
      );
    },
    onSuccess: (data) => {
      setResults(data.rows || []);
      setColumns(data.columns || []);
      toast.success('Query executed successfully');
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Failed to execute query');
    },
  });

  const handleRun = () => {
    if (!query.trim()) return;
    executeQuery.mutate(query);
  };

  const tableColumns = useMemo(
    () =>
      columns.map((col) => ({
        accessorKey: col,
        header: col,
        cell: (info: any) => {
          const val = info.getValue();
          return val !== null ? (
            String(val)
          ) : (
            <span className="text-muted-foreground italic">null</span>
          );
        },
      })),
    [columns]
  );

  const table = useReactTable({
    data: results || [],
    columns: tableColumns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="flex h-[calc(100vh-12rem)] flex-col gap-4">
      <Card className="flex flex-1 flex-col shadow-sm">
        <CardHeader className="flex-row items-center justify-between border-b px-4 py-3">
          <CardTitle className="font-medium text-sm">SQL Studio</CardTitle>
          <Button size="sm" onClick={handleRun} disabled={executeQuery.isPending} className="gap-2">
            {executeQuery.isPending ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Play className="h-4 w-4" />
            )}
            Run Query
          </Button>
        </CardHeader>
        <div className="min-h-50 flex-1">
          <MonacoEditor
            height="100%"
            language="sql"
            theme="vs-dark"
            value={query}
            onChange={(value) => setQuery(value || '')}
            options={{
              minimap: { enabled: false },
              padding: { top: 16 },
              fontSize: 14,
            }}
          />
        </div>
      </Card>

      <Card className="flex h-1/2 flex-col overflow-hidden shadow-sm">
        <CardHeader className="border-b bg-muted/20 px-4 py-3">
          <CardTitle className="font-medium text-sm">Results</CardTitle>
        </CardHeader>
        <CardContent className="flex-1 overflow-auto p-0">
          {results === null ? (
            <div className="flex h-full items-center justify-center text-muted-foreground text-sm">
              Enter a query and click Run to see results.
            </div>
          ) : results.length === 0 ? (
            <div className="flex h-full items-center justify-center text-muted-foreground text-sm">
              No results returned.
            </div>
          ) : (
            <Table>
              <TableHeader>
                {table.getHeaderGroups().map((headerGroup) => (
                  <TableRow key={headerGroup.id}>
                    {headerGroup.headers.map((header) => (
                      <TableHead key={header.id} className="whitespace-nowrap">
                        {header.isPlaceholder
                          ? null
                          : flexRender(header.column.columnDef.header, header.getContext())}
                      </TableHead>
                    ))}
                  </TableRow>
                ))}
              </TableHeader>
              <TableBody>
                {table.getRowModel().rows.length ? (
                  table.getRowModel().rows.map((row) => (
                    <TableRow key={row.id}>
                      {row.getVisibleCells().map((cell) => (
                        <TableCell key={cell.id} className="whitespace-nowrap">
                          {flexRender(cell.column.columnDef.cell, cell.getContext())}
                        </TableCell>
                      ))}
                    </TableRow>
                  ))
                ) : (
                  <TableRow>
                    <TableCell colSpan={columns.length} className="h-24 text-center">
                      No results.
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
