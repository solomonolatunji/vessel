import { FitAddon } from '@xterm/addon-fit';
import { SearchAddon } from '@xterm/addon-search';
import { WebLinksAddon } from '@xterm/addon-web-links';
import { Terminal } from '@xterm/xterm';
import { env } from '#/env';

import '@xterm/xterm/css/xterm.css';
import { useEffect, useRef, useState } from 'react';

interface WebTerminalProps {
  serviceId: string;
}

export function WebTerminal({ serviceId }: WebTerminalProps) {
  const terminalRef = useRef<HTMLDivElement>(null);
  const searchAddonRef = useRef<SearchAddon | null>(null);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    if (!terminalRef.current) return;

    const term = new Terminal({
      cursorBlink: true,
      fontFamily: 'Menlo, Monaco, "Courier New", monospace',
      fontSize: 14,
      theme: {
        background: '#09090b', // zinc-950
        foreground: '#fafafa', // zinc-50
      },
    });

    const fitAddon = new FitAddon();
    const webLinksAddon = new WebLinksAddon();
    const searchAddon = new SearchAddon();
    searchAddonRef.current = searchAddon;

    term.loadAddon(fitAddon);
    term.loadAddon(webLinksAddon);
    term.loadAddon(searchAddon);
    term.open(terminalRef.current);
    fitAddon.fit();

    const handleResize = () => {
      fitAddon.fit();
    };

    window.addEventListener('resize', handleResize);

    const apiUrl = new URL(env.VITE_API_URL || '', window.location.origin);
    const protocol = apiUrl.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${apiUrl.host}/api/ws/services/${serviceId}/terminal`;

    const socket = new WebSocket(wsUrl);
    socket.binaryType = 'arraybuffer';

    socket.onopen = () => {
      term.writeln('\x1b[32mConnected to terminal.\x1b[0m');
    };

    socket.onmessage = (event) => {
      if (event.data instanceof ArrayBuffer) {
        term.write(new Uint8Array(event.data));
      } else {
        term.write(event.data);
      }
    };

    socket.onclose = () => {
      term.writeln('\n\x1b[31mTerminal connection closed.\x1b[0m');
    };

    socket.onerror = () => {
      term.writeln('\n\x1b[31mWebSocket error occurred.\x1b[0m');
    };

    term.onData((data) => {
      if (socket.readyState === WebSocket.OPEN) {
        socket.send(data);
      }
    });

    return () => {
      window.removeEventListener('resize', handleResize);
      socket.close();
      term.dispose();
    };
  }, [serviceId]);

  return (
    <div className="flex h-full w-full flex-col overflow-hidden rounded-md border border-zinc-800 bg-zinc-950">
      <div className="flex items-center justify-between border-zinc-800 border-b bg-zinc-900 px-4 py-2">
        <div className="flex items-center gap-4">
          <h3 className="font-medium text-sm text-zinc-300">Terminal Shell</h3>
          <div className="flex items-center gap-1 rounded-md border border-zinc-700 bg-zinc-950 px-2">
            <input
              type="text"
              placeholder="Search output..."
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
        </div>
      </div>
      <div className="flex-1 overflow-hidden p-2" ref={terminalRef} />
    </div>
  );
}
