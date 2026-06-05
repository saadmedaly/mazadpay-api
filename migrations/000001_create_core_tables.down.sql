-- MazadPay — Migration 000001: Drop Core Tables (DOWN)
-- Run with: migrate -path migrations -database "$DATABASE_URL" down 1
-- Drops tables in reverse dependency order to avoid FK constraint errors.

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS contact_messages;
DROP TABLE IF EXISTS favorites;
DROP TABLE IF EXISTS bids;
DROP TABLE IF EXISTS auction_images;
DROP TABLE IF EXISTS auctions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS categories;
