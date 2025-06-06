import { drizzle } from 'drizzle-orm/bun-sqlite';
import { Database } from 'bun:sqlite';
import { sql } from 'drizzle-orm';

// Database setup
const sqlite = new Database('./sqlite.db');
export const db = drizzle(sqlite);

import { linksTable } from './link.sql';

export type Link = typeof linksTable.$inferSelect;
export type NewLink = typeof linksTable.$inferInsert;

// Database interaction functions

/**
 * Creates a new short link.
 */
export async function createLink(newLink: NewLink): Promise<Link[]> {
  return db.insert(linksTable).values(newLink).returning();
}

/**
 * Retrieves a link by its short URL.
 */
export async function getLinkByShortUrl(shortUrl: string): Promise<Link | undefined> {
  return db.select().from(linksTable).where(sql`${linksTable.shortUrl} = ${shortUrl}`).get();
}

/**
 * Increments the click count for a link.
 */
export async function incrementClickCount(shortUrl: string): Promise<void> {
  await db.update(linksTable)
    .set({
      clickCount: sql`${linksTable.clickCount} + 1`,
      updatedAt: new Date(),
    })
    .where(sql`${linksTable.shortUrl} = ${shortUrl}`);
}

/**
 * Retrieves all links.
 */
export async function getAllLinks(): Promise<Link[]> {
  return db.select().from(linksTable).all();
}

/**
 * Deletes a link by its short URL.
 */
export async function deleteLinkByShortUrl(shortUrl: string): Promise<void> {
  await db.delete(linksTable).where(sql`${linksTable.shortUrl} = ${shortUrl}`);
}

// Optional: Run migrations (if you set up drizzle-kit)
// import { migrate } from "drizzle-orm/bun-sqlite/migrator";
// migrate(db, { migrationsFolder: "./drizzle" }); // Ensure you have a migrations folder if you use this
