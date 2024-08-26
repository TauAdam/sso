INSERT INTO apps (id, name, secret)
VALUES (1, 'for-test', 'secret-string')
ON CONFLICT DO NOTHING;