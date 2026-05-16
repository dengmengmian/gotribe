-- Phase 6: Add ToC-specific fields to post table
-- These fields exist in the Go model (model.Post) but were missing from the initial schema dump

ALTER TABLE public.post ADD COLUMN IF NOT EXISTS post_id character varying(255);
ALTER TABLE public.post ADD COLUMN IF NOT EXISTS dynamic_type character varying(50);
ALTER TABLE public.post ADD COLUMN IF NOT EXISTS sort bigint DEFAULT 0;
ALTER TABLE public.post ADD COLUMN IF NOT EXISTS event_start_at timestamp;
ALTER TABLE public.post ADD COLUMN IF NOT EXISTS event_end_at timestamp;
ALTER TABLE public.post ADD COLUMN IF NOT EXISTS register_url character varying(255);

-- Add unique index on post_id (matches Go model gorm tag: uniqueIndex)
CREATE UNIQUE INDEX IF NOT EXISTS idx_post_post_id ON public.post USING btree (post_id);
