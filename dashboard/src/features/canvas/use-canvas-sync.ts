import type { Edge, Node } from '@xyflow/react';
import { useEffect, useState } from 'react';

export function useCanvasSync(envData: any) {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);

  useEffect(() => {
    if (!envData) return;

    const newNodes: Node[] = [];
    let yOffset = 0;

    envData.apps?.forEach((app: any, i: number) => {
      newNodes.push({
        id: `app-${app.id}`,
        type: 'appService',
        position: { x: 100 + i * 200, y: yOffset },
        data: { name: app.name, status: app.status },
      });
    });

    yOffset += 150;

    envData.databases?.forEach((db: any, i: number) => {
      newNodes.push({
        id: `db-${db.id}`,
        type: 'database',
        position: { x: 100 + i * 200, y: yOffset },
        data: { name: db.name, engine: db.engine, status: db.status },
      });
    });

    setNodes(newNodes);
  }, [envData]);

  return { nodes, edges, setNodes, setEdges };
}
