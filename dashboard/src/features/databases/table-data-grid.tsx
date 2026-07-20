import { useQuery } from '@tanstack/react-query';
import { ChevronLeft, ChevronRight, Database, Filter, Plus } from 'lucide-react';
import { useState } from 'react';

import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
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
import { RowEditorModal } from '#/features/databases/row-editor-modal';
import { apiClient } from '#/lib/apiClient';

interface TableDataGridProps {
  databaseId: string;
}

interface SchemaResponse {
  tables: string[];
}

interface DataResponse {
  rows: Record<string, unknown>[];
  columns: string[];
  totalCount: number;
}

export function TableDataGrid({ databaseId }: TableDataGridProps) {
  const [selectedTable, setSelectedTable] = useState<string>('');
  const [page, setPage] = useState(1);
  const [filter, setFilter] = useState('');

  const [isEditorOpen, setIsEditorOpen] = useState(false);
  const [editorData, setEditorData] = useState<Record<string, unknown> | null>(null);

  const { data: schema, isLoading: isLoadingSchema } = useQuery({
    queryKey: ['database-schemas', databaseId],
    queryFn: () => apiClient.get<SchemaResponse>(`/databases/${databaseId}/schemas`),
  });

  const activeTable = selectedTable || (schema?.tables?.[0] ?? '');

  const { data: tableData, isLoading: isLoadingData } = useQuery({
    queryKey: ['table-data', databaseId, activeTable, page, filter],
    queryFn: () =>
      apiClient.get<DataResponse>(
        `/databases/${databaseId}/data/${activeTable}?page=${page}&limit=50&filter=${encodeURIComponent(filter)}`
      ),
    enabled: !!activeTable,
  });

  const columns = tableData?.columns || [];
  const rows = tableData?.rows || [];

  const handleAddRow = () => {
    setEditorData(null);
    setIsEditorOpen(true);
  };

  const handleEditRow = (row: Record<string, unknown>) => {
    setEditorData(row);
    setIsEditorOpen(true);
  };

  return (
    <Card className="flex h-[calc(100vh-12rem)] flex-col">
      <CardHeader className="border-b py-4">
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <CardTitle className="flex items-center gap-2 text-lg">
              <Database className="h-5 w-5" />
              Data Browser
            </CardTitle>
            <CardDescription>View, filter, and edit relational data.</CardDescription>
          </div>
          <div className="flex items-center gap-3">
            <div className="relative w-48">
              <Filter className="absolute top-2.5 left-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Filter (e.g. id > 5)"
                value={filter}
                onChange={(e) => setFilter(e.target.value)}
                className="pl-9"
              />
            </div>

            <Select value={activeTable} onValueChange={setSelectedTable} disabled={isLoadingSchema}>
              <SelectTrigger className="w-45">
                <SelectValue placeholder="Select a table" />
              </SelectTrigger>
              <SelectContent>
                {schema?.tables?.map((table) => (
                  <SelectItem key={table} value={table}>
                    {table}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <Button onClick={handleAddRow} disabled={!activeTable} className="gap-2">
              <Plus className="h-4 w-4" />
              Add Row
            </Button>
          </div>
        </div>
      </CardHeader>

      <CardContent className="flex-1 overflow-auto p-0">
        <Table>
          <TableHeader>
            <TableRow>
              {columns.map((col) => (
                <TableHead key={col} className="whitespace-nowrap">
                  {col}
                </TableHead>
              ))}
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoadingData ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length || 1}
                  className="h-32 text-center text-muted-foreground"
                >
                  Loading data...
                </TableCell>
              </TableRow>
            ) : rows.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length || 1}
                  className="h-32 text-center text-muted-foreground"
                >
                  No data found in table {activeTable}.
                </TableCell>
              </TableRow>
            ) : (
              rows.map((row, i) => (
                <TableRow
                  key={i}
                  onDoubleClick={() => handleEditRow(row)}
                  className="cursor-pointer hover:bg-muted/50"
                >
                  {columns.map((col) => (
                    <TableCell key={col} className="max-w-50 truncate whitespace-nowrap">
                      {row[col] !== null ? (
                        String(row[col])
                      ) : (
                        <span className="text-muted-foreground italic">null</span>
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </CardContent>

      <div className="flex items-center justify-between border-t p-4">
        <div className="text-muted-foreground text-sm">Showing page {page}</div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1 || isLoadingData}
          >
            <ChevronLeft className="mr-1 h-4 w-4" />
            Previous
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => setPage((p) => p + 1)}
            disabled={rows.length < 50 || isLoadingData}
          >
            Next
            <ChevronRight className="ml-1 h-4 w-4" />
          </Button>
        </div>
      </div>

      {isEditorOpen && (
        <RowEditorModal
          isOpen={isEditorOpen}
          onClose={() => setIsEditorOpen(false)}
          databaseId={databaseId}
          tableName={activeTable}
          columns={columns}
          initialData={editorData}
        />
      )}
    </Card>
  );
}
