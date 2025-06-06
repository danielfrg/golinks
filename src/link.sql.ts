import { sql } from 'drizzle-orm';
import { integer, sqliteTable, text } from 'drizzle-orm/sqlite-core';

export const linksTable = sqliteTable('links', {
  id: integer('id').primaryKey({ autoIncrement: true }),
  shortUrl: text('short_url').notNull().unique(),
  longUrl: text('long_url').notNull(),
  clickCount: integer('click_count').default(0).notNull(),
  createdAt: integer('created_at', { mode: 'timestamp' })
    .default(sql`(strftime('%s', 'now'))`)
    .notNull(),
  updatedAt: integer('updated_at', { mode: 'timestamp' })
    .default(sql`(strftime('%s', 'now'))`)
    .notNull()
    .$onUpdate(() => new Date()),
});
