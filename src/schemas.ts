import { z } from '@hono/zod-openapi';

export const LinkSchema = z.object({
  id: z.number().int().openapi({
    description: 'The unique identifier for the link.',
    example: 1,
  }),
  shortUrl: z.string().min(1).max(50).regex(/^[a-zA-Z0-9_-]+$/).openapi({
    description: 'The short URL alias.',
    example: 'mylink',
  }),
  longUrl: z.string().url().openapi({
    description: 'The original long URL.',
    example: 'https://example.com/very/long/url',
  }),
  clickCount: z.number().int().min(0).openapi({
    description: 'Number of times the link has been clicked.',
    example: 10,
  }),
  createdAt: z.number().int().openapi({
    description: 'Timestamp of when the link was created (Unix epoch seconds).',
    example: 1678886400,
  }),
  updatedAt: z.number().int().openapi({
    description: 'Timestamp of when the link was last updated (Unix epoch seconds).',
    example: 1678889900,
  }),
}).openapi('Link');

export const CreateLinkSchema = z.object({
  longUrl: z.string().url({ message: "Invalid URL format for longUrl" }).openapi({
    description: 'The original long URL to shorten.',
    example: 'https://example.com/my-long-article-url',
  }),
  shortUrl: z.string().min(1).max(50).regex(/^[a-zA-Z0-9_-]+$/, { message: "Short URL can only contain alphanumeric characters, underscores, and hyphens." }).optional().openapi({
    description: 'Optional custom short URL alias. If not provided, one will be generated.',
    example: 'custom-link',
  }),
}).openapi('CreateLinkPayload');

export const ShortUrlParamSchema = z.object({
  shortUrl: z.string().min(1).max(50)
    .regex(/^[a-zA-Z0-9_-]+$/, { message: "Short URL parameter is invalid. It can only contain alphanumeric characters, underscores, and hyphens." })
    .openapi({
      param: {
        name: 'shortUrl',
        in: 'path',
      },
      description: 'The short URL alias.',
      example: 'mylink',
    }),
});

export const ErrorSchema = z.object({
  message: z.string().openapi({
    description: 'Error message.',
    example: 'Resource not found',
  }),
}).openapi('Error');

export const MessageSchema = z.object({
  message: z.string().openapi({
    description: 'A descriptive message.',
    example: 'Operation successful.'
  })
}).openapi('Message');
