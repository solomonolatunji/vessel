const http = require('node:http');
const fs = require('node:fs');
const path = require('node:path');

const PORT = process.env.PORT || 80;
const SCRIPTS = new Set(['/install.sh', '/upgrade.sh', '/cli']);
const POSTHOG_KEY = process.env.POSTHOG_API_KEY;
const POSTHOG_HOST = process.env.POSTHOG_HOST || 'https://us.i.posthog.com';

function trackEvent(eventName, req) {
  if (!POSTHOG_KEY) return;
  const ip = req.headers['x-forwarded-for'] || req.socket.remoteAddress;
  const userAgent = req.headers['user-agent'] || 'unknown';
  fetch(`${POSTHOG_HOST}/capture/`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      api_key: POSTHOG_KEY,
      event: eventName,
      distinct_id: ip,
      properties: {
        $ip: ip,
        $user_agent: userAgent,
        $current_url: req.url,
      }
    })
  }).catch(err => console.error('PostHog error:', err));
}

const server = http.createServer((req, res) => {
  if (req.url === '/version') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ version: '0.1.0', latest: true }));
    return;
  }

  let filename = path.basename(req.url);
  if (req.url === '/') filename = 'install.sh';
  if (req.url === '/cli') filename = 'install-cli.sh';

  const file = path.join(__dirname, filename);
  const name = path.basename(req.url === '/' ? '/install.sh' : req.url);

  if (!SCRIPTS.has(req.url === '/' ? '/install.sh' : req.url) || !fs.existsSync(file)) {
    res.writeHead(404);
    res.end('Not found');
    return;
  }

  res.writeHead(200, { 'Content-Type': 'text/x-shellscript' });
  fs.createReadStream(file).pipe(res);
  trackEvent('script_downloaded', req);
});

server.listen(PORT, () => {
  console.log(`📦 install server running on port ${PORT}`);
});
