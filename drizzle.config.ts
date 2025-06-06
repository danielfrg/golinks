import type { Config } from 'drizzle-kit';

export default {
  dialect: 'sqlite',
  schema: './src/db.ts',
  out: './migrations',
  dbCredentials: {
    url: './sqlite.db',
  },
  verbose: true,
  strict: true,
} satisfies Config;
