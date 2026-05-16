-- Phase 6: Create example and example_post tables
-- These tables exist in Go models (model.Example, model.ExamplePost) but were missing from the initial schema dump

CREATE TABLE IF NOT EXISTS public.example (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    example_id character varying(10),
    project_id bigint NOT NULL,
    username character varying(30) NOT NULL,
    title character varying(255) NOT NULL,
    content text NOT NULL,
    description character varying(300) NOT NULL,
    status smallint DEFAULT 1 NOT NULL,
    user_id bigint,
    owner_username character varying(30),
    owner_nickname character varying(30),
    name character varying(255),
    primary_post_id character varying(255)
);

CREATE TABLE IF NOT EXISTS public.example_post (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    example_record_id bigint,
    project_id bigint NOT NULL,
    user_id bigint,
    post_id character varying(255),
    post_title character varying(255),
    post_type smallint,
    post_status smallint,
    sort integer DEFAULT 0
);

-- Primary keys
ALTER TABLE ONLY public.example ADD CONSTRAINT example_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.example_post ADD CONSTRAINT example_post_pkey PRIMARY KEY (id);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_example_deleted_at ON public.example USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS idx_example_project_id ON public.example USING btree (project_id);
CREATE INDEX IF NOT EXISTS idx_example_username ON public.example USING btree (username);
CREATE INDEX IF NOT EXISTS idx_example_user_id ON public.example USING btree (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_example_example_id ON public.example USING btree (example_id);

CREATE INDEX IF NOT EXISTS idx_example_post_deleted_at ON public.example_post USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS idx_example_post_project_id ON public.example_post USING btree (project_id);
CREATE INDEX IF NOT EXISTS idx_example_post_user_id ON public.example_post USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_example_post_example_record_id ON public.example_post USING btree (example_record_id);

-- Sequences
CREATE SEQUENCE IF NOT EXISTS public.example_id_seq
    START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

CREATE SEQUENCE IF NOT EXISTS public.example_post_id_seq
    START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER SEQUENCE public.example_id_seq OWNED BY public.example.id;
ALTER SEQUENCE public.example_post_id_seq OWNED BY public.example_post.id;

-- Defaults
ALTER TABLE ONLY public.example ALTER COLUMN id SET DEFAULT nextval('public.example_id_seq'::regclass);
ALTER TABLE ONLY public.example_post ALTER COLUMN id SET DEFAULT nextval('public.example_post_id_seq'::regclass);

-- Comments
COMMENT ON TABLE public.example IS '示例业务单';
COMMENT ON COLUMN public.example.example_id IS '唯一字符ID/分布式ID';
COMMENT ON COLUMN public.example.project_id IS '项目ID';
COMMENT ON COLUMN public.example.username IS '用户名';
COMMENT ON COLUMN public.example.title IS '标题';
COMMENT ON COLUMN public.example.content IS '内容';
COMMENT ON COLUMN public.example.description IS '描述';
COMMENT ON COLUMN public.example.status IS '状态，1-正常；2-禁用';
COMMENT ON COLUMN public.example.user_id IS '用户ID';
COMMENT ON COLUMN public.example.owner_username IS '拥有者用户名';
COMMENT ON COLUMN public.example.owner_nickname IS '拥有者昵称';
COMMENT ON COLUMN public.example.name IS '名称';
COMMENT ON COLUMN public.example.primary_post_id IS '主要文章ID';

COMMENT ON TABLE public.example_post IS '示例业务单与文章的关联表';
COMMENT ON COLUMN public.example_post.example_record_id IS '示例业务单ID';
COMMENT ON COLUMN public.example_post.project_id IS '项目ID';
COMMENT ON COLUMN public.example_post.user_id IS '用户ID';
COMMENT ON COLUMN public.example_post.post_id IS '文章ID';
COMMENT ON COLUMN public.example_post.post_title IS '文章标题';
COMMENT ON COLUMN public.example_post.post_type IS '文章类型';
COMMENT ON COLUMN public.example_post.post_status IS '文章状态';
COMMENT ON COLUMN public.example_post.sort IS '排序';
