
INSERT INTO "user" (id, "user", pwd, role_id, exp_time, flow, in_flow, out_flow, flow_reset_time, num, created_time, updated_time, status) 
VALUES (1, 'admin_user', '3c85cdebade1c51cf64ca9f3c09d182d', 0, 2727251700000, 99999, 0, 0, 1, 99999, 1748914865000, 1754011744252, 1)
ON CONFLICT (id) DO NOTHING;

-- 同步自增序列（PG 的 BIGSERIAL 序列不随显式 id 插入而推进）
SELECT setval('user_id_seq', (SELECT MAX(id) FROM "user"));


INSERT INTO vite_config (id, name, value, time) 
VALUES (1, 'app_name', 'flux', 1755147963000)
ON CONFLICT (id) DO NOTHING;

SELECT setval('vite_config_id_seq', (SELECT MAX(id) FROM vite_config));
