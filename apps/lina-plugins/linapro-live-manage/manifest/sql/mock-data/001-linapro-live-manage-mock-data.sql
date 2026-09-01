-- Mock data: live rooms and live content records for live streaming management demos.
-- 模拟数据：直播后台管理演示使用的直播间与直播内容记录。

-- Two demo live rooms identified by stable room codes.
-- 两个以稳定房间编码标识的演示直播间。
INSERT INTO plugin_linapro_live_manage_room ("tenant_id", "room_code", "room_name", "room_type", "status", "description", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    0,
    'ROOM-MAIN',
    '主堂直播间',
    1,
    0,
    '主堂聚会专用直播间',
    admin."id",
    admin."id",
    '2026-04-20 09:00:00',
    '2026-04-20 09:00:00'
FROM sys_user admin
WHERE admin."username" = 'admin'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_live_manage_room room
    WHERE room."tenant_id" = 0 AND room."room_code" = 'ROOM-MAIN'
);

INSERT INTO plugin_linapro_live_manage_room ("tenant_id", "room_code", "room_name", "room_type", "status", "description", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    0,
    'ROOM-CHAPEL',
    '小礼堂直播间',
    3,
    0,
    '团契与祷告会直播间',
    admin."id",
    admin."id",
    '2026-04-20 09:10:00',
    '2026-04-20 09:10:00'
FROM sys_user admin
WHERE admin."username" = 'admin'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_live_manage_room room
    WHERE room."tenant_id" = 0 AND room."room_code" = 'ROOM-CHAPEL'
);

-- One demo live content bound to the main room resolved by stable room code.
-- 一条以稳定房间编码解析主堂直播间的演示直播内容。
INSERT INTO plugin_linapro_live_manage_live ("tenant_id", "room_id", "title", "content", "live_date", "page_url", "live_url", "cover_url", "song_name", "song_list", "lead_singer", "accompaniment", "host", "sermon_title", "preacher", "preacher_identity", "scripture_ref", "scripture_content", "outline", "device_info", "reception", "state", "is_public", "start_time", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    0,
    room."id",
    '主日崇拜直播',
    '<p>本周主日崇拜线上直播，欢迎弟兄姊妹准时参加。</p>',
    '2026-04-26',
    'https://example.com/live/sunday-service',
    'https://example.com/hls/sunday-service.m3u8',
    'https://example.com/assets/sunday-cover.jpg',
    '赞美之泉',
    '[{"name":"赞美之泉","singer":"诗班","order":1},{"name":"恩典之路","singer":"诗班","order":2}]',
    '诗班',
    '钢琴伴奏',
    '张弟兄',
    '信心的操练',
    '王牧师',
    '牧师',
    '约翰福音 3:16',
    '神爱世人，甚至将他的独生子赐给他们，叫一切信他的，不至灭亡，反得永生。',
    '一、信心的来源；二、信心的实践；三、信心的果效。',
    'Canon EOS 摄像机，Shure SM7B 麦克风',
    '李姊妹 13800000000',
    0,
    1,
    NULL,
    admin."id",
    admin."id",
    '2026-04-21 10:30:00',
    '2026-04-21 10:30:00'
FROM sys_user admin
JOIN plugin_linapro_live_manage_room room
  ON room."tenant_id" = 0 AND room."room_code" = 'ROOM-MAIN'
WHERE admin."username" = 'admin'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_live_manage_live live
    WHERE live."tenant_id" = 0 AND live."title" = '主日崇拜直播'
  );
