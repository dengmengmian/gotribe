-- 添加 example 和 example_post 表（ToC API 示例业务模块）
CREATE TABLE IF NOT EXISTS public.example (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    example_id VARCHAR(10) NOT NULL UNIQUE,
    project_id BIGINT NOT NULL,
    username VARCHAR(30) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    description VARCHAR(300) NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    user_id BIGINT,
    owner_username VARCHAR(30),
    owner_nickname VARCHAR(30),
    name VARCHAR(255),
    primary_post_id VARCHAR(255)
);
CREATE INDEX IF NOT EXISTS idx_example_project_id ON public.example (project_id);
CREATE INDEX IF NOT EXISTS idx_example_user_id ON public.example (user_id);
CREATE INDEX IF NOT EXISTS idx_example_username ON public.example (username);
CREATE INDEX IF NOT EXISTS idx_example_deleted_at ON public.example (deleted_at);

CREATE TABLE IF NOT EXISTS public.example_post (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    example_record_id BIGINT,
    project_id BIGINT NOT NULL,
    user_id BIGINT,
    post_id VARCHAR(255),
    post_title VARCHAR(255),
    post_type SMALLINT,
    post_status SMALLINT,
    sort INT DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_example_post_example_record_id ON public.example_post (example_record_id);
CREATE INDEX IF NOT EXISTS idx_example_post_project_id ON public.example_post (project_id);
CREATE INDEX IF NOT EXISTS idx_example_post_user_id ON public.example_post (user_id);
CREATE INDEX IF NOT EXISTS idx_example_post_deleted_at ON public.example_post (deleted_at);
