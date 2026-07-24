import type { Edge, Node } from '@xyflow/react';
import { useEffect, useState } from 'react';

export function useCanvasSync(envData: Record<string, unknown>) {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);

  useEffect(() => {
    if (!envData) return;

    const newNodes: Node[] = [];
    let yOffset = 0;

    (envData.apps as Record<string, unknown>[])?.forEach(
      (app: Record<string, unknown>, i: number) => {
        newNodes.push({
          id: `app-${app.id}`,
          type: 'appService',
          position: { x: 100 + i * 200, y: yOffset },
          data: { name: app.name, status: app.status },
        });
      }
    );

    yOffset += 150;

    (envData.databases as Record<string, unknown>[])?.forEach(
      (db: Record<string, unknown>, i: number) => {
        newNodes.push({
          id: `db-${db.id}`,
          type: 'database',
          position: { x: 100 + i * 200, y: yOffset },
          data: { name: db.name, engine: db.engine, status: db.status },
        });
      }
    );

    setNodes(newNodes);
  }, [envData]);

  return { nodes, edges, setNodes, setEdges };
}
