import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import { Terminal } from '@xterm/xterm';
import { useEffect, useRef } from 'react';
import { env } from '#/env';
import '@xterm/xterm/css/xterm.css';
import { authStore } from '#/stores/authStore';

interface WebTerminalProps {
  serviceId: string;
}

export function WebTerminal({ serviceId }: WebTerminalProps) {
  const terminalRef = useRef<HTMLDivElement>(null);

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

    term.loadAddon(fitAddon);
    term.loadAddon(webLinksAddon);
    term.open(terminalRef.current);
    fitAddon.fit();

    const handleResize = () => {
      fitAddon.fit();
    };

    window.addEventListener('resize', handleResize);

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsHost = env.VITE_API_URL.replace(/^http(s?):\/\//, '');
    const wsUrl = `${protocol}//${wsHost}/api/services/${serviceId}/terminal?token=${authStore.state.token || ''}`;

    const socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      term.writeln('\x1b[32mConnected to terminal.\x1b[0m');
    };

    socket.onmessage = (event) => {
      term.write(event.data);
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
        <h3 className="font-medium text-sm text-zinc-300">Terminal Shell</h3>
      </div>
      <div className="flex-1 overflow-hidden p-2" ref={terminalRef} />
    </div>
  );
}
