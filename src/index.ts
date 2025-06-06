import { Hono } from 'hono';
import { HTTPException } from 'hono/http-exception';
import { RegExpRouter } from 'hono/router/reg-exp-router';
import {
  createLink,
  getLinkByShortUrl,
  incrementClickCount,
  getAllLinks,
  deleteLinkByShortUrl,
  NewLink,
} from './db';

const app = new Hono({
  router: new RegExpRouter(), // Using RegExpRouter for more complex route matching if needed
});

// List all links or show a welcome message
app.get('/', async (c) => {
  try {
    const links = await getAllLinks();
    if (links.length === 0) {
      return c.json({ message: 'Welcome to GoLinks! No links created yet. POST to /api/links to create one.' });
    }
    return c.json(links);
  } catch (error) {
    console.error('Error fetching links:', error);
    throw new HTTPException(500, { message: 'Failed to retrieve links' });
  }
});

// Create a new short link
app.post('/api/links', async (c) => {
  try {
    const { longUrl, shortUrl } = await c.req.json<{ longUrl: string; shortUrl?: string }>();

    if (!longUrl) {
      throw new HTTPException(400, { message: 'longUrl is required' });
    }

    let finalShortUrl = shortUrl;
    if (!finalShortUrl) {
      // Generate a simple random short URL if not provided (e.g., 6 random alphanumeric chars)
      // For a real application, ensure uniqueness and a better generation strategy
      finalShortUrl = Math.random().toString(36).substring(2, 8);
    }

    // Check if custom short URL already exists
    if (shortUrl) {
        const existing = await getLinkByShortUrl(shortUrl);
        if (existing) {
            throw new HTTPException(409, { message: `Short URL '${shortUrl}' already exists.` });
        }
    }

    const newLinkData: NewLink = { longUrl, shortUrl: finalShortUrl };
    const [created] = await createLink(newLinkData);
    return c.json(created, 201);
  } catch (error: any) {
    console.error('Error creating link:', error);
    if (error instanceof HTTPException) throw error;
    throw new HTTPException(500, { message: 'Failed to create link' });
  }
});

// Delete a short link
app.delete('/api/links/:shortUrl{[a-zA-Z0-9_-]+}', async (c) => {
  try {
    const shortUrl = c.req.param('shortUrl');
    const link = await getLinkByShortUrl(shortUrl);

    if (!link) {
      throw new HTTPException(404, { message: 'Short URL not found' });
    }

    await deleteLinkByShortUrl(shortUrl);
    return c.json({ message: `Link '${shortUrl}' deleted successfully` }, 200);
  } catch (error: any) {
    console.error('Error deleting link:', error);
    if (error instanceof HTTPException) throw error;
    throw new HTTPException(500, { message: 'Failed to delete link' });
  }
});


// Redirect short URL to long URL
app.get('/:shortUrl{[a-zA-Z0-9_-]+}', async (c) => {
  try {
    const shortUrl = c.req.param('shortUrl');
    const link = await getLinkByShortUrl(shortUrl);

    if (!link) {
      throw new HTTPException(404, { message: 'Short URL not found' });
    }

    // Increment click count asynchronously (don't wait for it to redirect faster)
    incrementClickCount(shortUrl).catch(err => console.error('Failed to increment click count:', err));

    return c.redirect(link.longUrl, 301); // 301 for permanent redirect, 302 for temporary
  } catch (error: any) {
    console.error('Error redirecting link:', error);
    if (error instanceof HTTPException) throw error;
    // Avoid redirecting to an error page for the actual shortlink itself.
    // Instead, show a clear message. The root '/' can list available links.
    if (c.req.path !== '/') { // Check to ensure this isn't the root path trying to redirect itself on initial load
        return c.text('Link not found or an error occurred.', 404);
    }
    // For the root path or other general errors, rethrow
    throw new HTTPException(500, { message: 'Failed to process redirect' });
  }
});

export default app;
