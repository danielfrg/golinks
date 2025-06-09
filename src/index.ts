import { HTTPException } from 'hono/http-exception';
import { Hono } from 'hono';

import apiApp from "./api";
import {
  getLinkByShortUrl,
  incrementClickCount,
  getAllLinks,
} from './db';

// Main Application
const app = new Hono();

// Global error handling for the main app
app.onError((err, c) => {
  if (err instanceof HTTPException) {
    return c.json({ message: err.message }, err.status);
  }
  console.error('Global error:', err);
  return c.json({ message: 'Internal Server Error' }, 500);
});


app.get('/', async (c) => {
  try {
    const links = await getAllLinks();

    if (links.length === 0) {
      return c.json({ message: 'Welcome to GoLinks! No links created yet. Use the API (POST to /api/links) to create one or manage links.' });
    }

    const linkList = links.map(link => ({
      short: `/${link.shortUrl}`,
      long: link.longUrl,
      clicks: link.clickCount
    }));
    return c.json({
      message: "Available GoLinks:",
      links: linkList,
      apiDocumentation: "/doc"
    });
  } catch (error) {
    console.error('Error fetching links for root path:', error);
    throw new HTTPException(500, { message: 'Failed to retrieve links' });
  }
});


// Redirector: short URL to long URL
app.get('/:shortUrl{[a-zA-Z0-9_-]+}', async (c) => {
  const shortUrl = c.req.param('shortUrl');

  try {
    const link = await getLinkByShortUrl(shortUrl);
    if (!link) {
      return c.text('Link not found.', 404);
    }
    incrementClickCount(shortUrl).catch(err => console.error('Failed to increment click count:', err));
    return c.redirect(link.longUrl, 301);

  } catch (error: any) {
    console.error('Error redirecting link:', error);
    return c.text('Error redirecting link.', 500);
  }
});

// Mount API app
app.route('/api', apiApp);

export default app;
