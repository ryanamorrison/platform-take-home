BEGIN TRANSACTION;
CREATE TABLE items (id SERIAL PRIMARY KEY,created_at date,updated_at date,deleted_at date,name text,description text);
INSERT INTO items VALUES('1','2024-11-19 15:37:14.109777+01:00','2024-11-19 15:37:14.109777+01:00',NULL,'test','ttt');
CREATE INDEX idx_items_deleted_at ON items(deleted_at);
COMMIT;
