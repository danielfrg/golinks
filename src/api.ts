import { OpenAPIHono, createRoute, z } from '@hono/zod-openapi'
import { swaggerUI } from '@hono/swagger-ui'
import { HTTPException } from 'hono/http-exception';

import {
  getLinkByShortUrl,
  getAllLinks,
  createLink,
  deleteLinkByShortUrl,
  type NewLink
} from './db';

import {
  LinkSchema,
  CreateLinkSchema,
  ShortUrlParamSchema,
  ErrorSchema,
  MessageSchema,
} from './schemas';

const api = new OpenAPIHono()

api.onError((err, c) => {
  if (err instanceof HTTPException) {
    return c.json({ message: err.message }, err.status);
  }

  console.error('API Error:', err);
  return c.json({ message: 'API Internal Server Error' }, 500);
});


const apiListLinksRoute = createRoute({
  method: 'get',
  path: '/links',
  summary: 'List all short links.',
  description: 'Retrieves all created short links.',
  responses: {
    200: {
      description: 'A list of links.',
      content: {
        'application/json': {
          schema: z.array(LinkSchema),
        },
      },
    },
    500: {
      description: 'Failed to retrieve links.',
      content: { 'application/json': { schema: ErrorSchema } },
    },
  },
  tags: ['API Links'],
});

api.openapi(apiListLinksRoute, async (c) => {
  try {
    const links = await getAllLinks();
    if (links.length === 0) {
      return c.json([]); // Return empty array if no links
    }
    return c.json(links.map(link => LinkSchema.parse(link)));
  } catch (error) {
    console.error('Error fetching links via API:', error);
    throw new HTTPException(500, { message: 'Failed to retrieve links' });
  }
});


const createLinkRoute = createRoute({
  method: 'post',
  path: '/links',
  summary: 'Create a new short link via API.',
  request: {
    body: {
      content: { 'application/json': { schema: CreateLinkSchema } },
    },
  },
  responses: {
    201: {
      description: 'Link created successfully.',
      content: { 'application/json': { schema: LinkSchema } },
    },
    400: {
      description: 'Bad request.',
      content: { 'application/json': { schema: ErrorSchema } },
    },
    409: {
      description: 'Conflict (e.g., custom short URL already exists).',
      content: { 'application/json': { schema: ErrorSchema } },
    },
    500: {
      description: 'Failed to create link.',
      content: { 'application/json': { schema: ErrorSchema } },
    },
  },
  tags: ['API Links'],
});

api.openapi(createLinkRoute, async (c) => {
  const { longUrl, shortUrl } = c.req.valid('json');
  try {
    let finalShortUrl = shortUrl;
    if (!finalShortUrl) {
      finalShortUrl = Math.random().toString(36).substring(2, 8);
      let attempts = 0;
      while (await getLinkByShortUrl(finalShortUrl) && attempts < 5) {
        finalShortUrl = Math.random().toString(36).substring(2, 8);
        attempts++;
      }
      if (await getLinkByShortUrl(finalShortUrl)) {
        throw new HTTPException(500, { message: 'Failed to generate a unique short URL.' });
      }
    } else {
      const existing = await getLinkByShortUrl(finalShortUrl);
      if (existing) {
        throw new HTTPException(409, { message: `Short URL '${finalShortUrl}' already exists.` });
      }
    }
    const newLinkData: NewLink = { longUrl, shortUrl: finalShortUrl };
    const [createdLink] = await createLink(newLinkData);
    return c.json(LinkSchema.parse(createdLink), 201);
  } catch (error: any) {
    console.error('Error creating link via API:', error);
    if (error instanceof HTTPException) throw error;
    throw new HTTPException(500, { message: 'Failed to create link' });
  }
});

const deleteLinkRoute = createRoute({
  method: 'delete',
  path: '/links/:shortUrl{[a-zA-Z0-9_-]+}', // Path relative to /api
  summary: 'Delete a short link via API.',
  request: {
    params: ShortUrlParamSchema,
  },
  responses: {
    200: {
      description: 'Link deleted successfully.',
      content: { 'application/json': { schema: MessageSchema } },
    },
    404: {
      description: 'Short URL not found.',
      content: { 'application/json': { schema: ErrorSchema } },
    },
    500: {
      description: 'Failed to delete link.',
      content: { 'application/json': { schema: ErrorSchema } },
    },
  },
  tags: ['API Links'],
});

api.openapi(deleteLinkRoute, async (c) => {
  const { shortUrl } = c.req.valid('param');
  try {
    const link = await getLinkByShortUrl(shortUrl);
    if (!link) {
      throw new HTTPException(404, { message: 'Short URL not found' });
    }
    await deleteLinkByShortUrl(shortUrl);
    return c.json({ message: `Link '${shortUrl}' deleted successfully` }, 200);
  } catch (error: any) {
    console.error('Error deleting link via API:', error);
    if (error instanceof HTTPException) throw error;
    throw new HTTPException(500, { message: 'Failed to delete link' });
  }
}
);

// The OpenAPI documentation will be available at /doc
api.doc('/openapi.json', {
  openapi: '3.0.0',
  info: {
    version: '1.0.0',
    title: 'golinks API',
  },
  servers: [
    {
      url: '/api',
      description: 'Default environment'
    }
  ]
})


api.get('/ui', swaggerUI({ url: '/api/openapi.json' }))

export default api
