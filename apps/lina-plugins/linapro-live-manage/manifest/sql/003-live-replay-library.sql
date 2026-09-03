-- 003: live replay library schema
-- 003：直播往期回放库数据结构

-- Purpose: Adds the per-live replay visibility switch consumed by the public
--          replay list and the play-info replay fallback.
-- 用途：为每场直播增加回放可见性开关，供公开回放列表与播放信息回放回退使用。

ALTER TABLE plugin_linapro_live_manage_live
    ADD COLUMN IF NOT EXISTS "replay_enabled" BOOLEAN NOT NULL DEFAULT TRUE;

COMMENT ON COLUMN plugin_linapro_live_manage_live."replay_enabled" IS
    'Whether the finished live is offered as a public replay; ongoing lives are unaffected';
