import { Handle, Position } from '@xyflow/react';

export function DatabaseNode({ data }: { data: any }) {
  return (
    <div className="rounded-md border-2 border-green-500 bg-white px-4 py-2 shadow-md">
      <div className="flex items-center">
        <div className="ml-2">
          <div className="font-bold text-sm">{data.name}</div>
          <div className="text-gray-500 text-xs">{data.engine}</div>
        </div>
      </div>
      <Handle type="target" position={Position.Top} className="w-16 bg-teal-500!" />
      <Handle type="source" position={Position.Bottom} className="w-16 bg-teal-500!" />
    </div>
  );
}
