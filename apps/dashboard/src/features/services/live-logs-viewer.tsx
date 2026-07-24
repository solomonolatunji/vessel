import { FitAddon } from '@xterm/addon-fit';
import { SearchAddon } from '@xterm/addon-search';
import { Terminal } from '@xterm/xterm';
import { env } from '#/env';
import { useAuthStore } from '#/stores/authStore';
import '@xterm/xterm/css/xterm.css';
import { useEffect, useRef, useState } from 'react';
import { AIDiagnoseDialog } from './ai-diagnose-dialog';

interface LiveLogsViewerProps {
  serviceId: string;
  deploymentId?: string;
}

export function LiveLogsViewer({ serviceId, deploymentId }: LiveLogsViewerProps) {
  const terminalRef = useRef<HTMLDivElement>(null);
  const [isConnected, setIsConnected] = useState(false);
  const searchAddonRef = useRef<SearchAddon | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [logsBuffer, setLogsBuffer] = useState('');

  useEffect(() => {
    if (!terminalRef.current) return;

    const term = new Terminal({
      cursorBlink: false,
      fontFamily: 'Menlo, Monaco, "Courier New", monospace',
      fontSize: 13,
      disableStdin: true,
      theme: {
        background: '#09090b',
        foreground: '#fafafa',
      },
    });

    const fitAddon = new FitAddon();
    const searchAddon = new SearchAddon();
    searchAddonRef.current = searchAddon;

    term.loadAddon(fitAddon);
    term.loadAddon(searchAddon);
    term.open(terminalRef.current);
    fitAddon.fit();

    const handleResize = () => fitAddon.fit();
    window.addEventListener('resize', handleResize);

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsHost = env.VITE_API_URL.replace(/^http(s?):\/\//, '');
    let wsUrl = `${protocol}//${wsHost}/api/services/${serviceId}/logs?token=${useAuthStore.getState().token || ''}`;
    if (deploymentId) {
      wsUrl += `&deploymentId=${deploymentId}`;
    }

    const socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      setIsConnected(true);
      term.writeln('\x1b[32mConnected to log stream...\x1b[0m');
    };

    socket.onmessage = (event) => {
      term.write(event.data);
      setLogsBuffer((prev) => (prev + event.data).slice(-5000)); // Keep last 5000 chars
    };

    socket.onclose = () => {
      setIsConnected(false);
      term.writeln('\n\x1b[31mLog stream disconnected.\x1b[0m');
    };

    return () => {
      window.removeEventListener('resize', handleResize);
      socket.close();
      term.dispose();
    };
  }, [serviceId, deploymentId]);

  return (
    <div className="flex h-125 w-full flex-col overflow-hidden rounded-md border border-zinc-800 bg-zinc-950">
      <div className="flex items-center justify-between border-zinc-800 border-b bg-zinc-900 px-4 py-2">
        <div className="flex items-center gap-4">
          <h3 className="font-medium text-sm text-zinc-300">Live Logs</h3>
          <div className="flex items-center gap-1 rounded-md border border-zinc-700 bg-zinc-950 px-2">
            <input
              type="text"
              placeholder="Search logs..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  searchAddonRef.current?.findNext(searchTerm);
                }
              }}
              className="w-48 bg-transparent py-1 text-xs text-zinc-300 outline-none placeholder:text-zinc-500"
            />
            <button
              type="button"
              onClick={() => searchAddonRef.current?.findPrevious(searchTerm)}
              className="px-1 text-xs text-zinc-400 hover:text-zinc-100"
            >
              ↑
            </button>
            <button
              type="button"
              onClick={() => searchAddonRef.current?.findNext(searchTerm)}
              className="px-1 text-xs text-zinc-400 hover:text-zinc-100"
            >
              ↓
            </button>
          </div>
          <AIDiagnoseDialog logs={logsBuffer} />
        </div>
        <div className="flex items-center gap-2">
          <span className={`h-2 w-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
          <span className="text-xs text-zinc-400">{isConnected ? 'Live' : 'Disconnected'}</span>
        </div>
      </div>
      <div className="flex-1 overflow-hidden p-2" ref={terminalRef} />
    </div>
  );
}
