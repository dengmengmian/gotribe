DROP INDEX IF EXISTS idx_post_post_id;

ALTER TABLE public.post DROP COLUMN IF EXISTS post_id;
ALTER TABLE public.post DROP COLUMN IF EXISTS dynamic_type;
ALTER TABLE public.post DROP COLUMN IF EXISTS sort;
ALTER TABLE public.post DROP COLUMN IF EXISTS event_start_at;
ALTER TABLE public.post DROP COLUMN IF EXISTS event_end_at;
ALTER TABLE public.post DROP COLUMN IF EXISTS register_url;
