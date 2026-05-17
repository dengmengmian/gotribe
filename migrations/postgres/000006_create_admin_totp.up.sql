-- 新增 admin_totp 表：管理员 TOTP 二次校验绑定记录
CREATE TABLE IF NOT EXISTS public.admin_totp (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    admin_id BIGINT NOT NULL,
    secret_cipher TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    recovery_codes TEXT,
    last_used_at BIGINT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_admin_totp_admin_id ON public.admin_totp (admin_id);
CREATE INDEX IF NOT EXISTS idx_admin_totp_deleted_at ON public.admin_totp (deleted_at);
