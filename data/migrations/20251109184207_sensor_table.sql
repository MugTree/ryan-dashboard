-- +goose Up
-- +goose StatementBegin
CREATE TABLE sensor_data (
    id INTEGER PRIMARY KEY,
    depth INTEGER DEFAULT 0 NOT NULL,
    record_created TEXT
);

INSERT INTO sensor_data (depth, record_created) VALUES (7, '2025-11-09T12:00:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (2, '2025-11-09T12:01:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (4, '2025-11-09T12:02:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (8, '2025-11-09T12:03:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (3, '2025-11-09T12:04:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (10, '2025-11-09T12:05:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (6, '2025-11-09T12:06:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (1, '2025-11-09T12:07:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (5, '2025-11-09T12:08:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (9, '2025-11-09T12:09:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (8, '2025-11-09T12:10:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (3, '2025-11-09T12:11:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (10, '2025-11-09T12:12:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (5, '2025-11-09T12:13:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (1, '2025-11-09T12:14:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (7, '2025-11-09T12:15:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (4, '2025-11-09T12:16:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (6, '2025-11-09T12:17:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (9, '2025-11-09T12:18:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (2, '2025-11-09T12:19:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (10, '2025-11-09T12:20:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (5, '2025-11-09T12:21:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (8, '2025-11-09T12:22:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (4, '2025-11-09T12:23:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (7, '2025-11-09T12:24:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (3, '2025-11-09T12:25:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (9, '2025-11-09T12:26:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (1, '2025-11-09T12:27:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (6, '2025-11-09T12:28:00Z');
INSERT INTO sensor_data (depth, record_created) VALUES (2, '2025-11-09T12:29:00Z');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sensor_data;
-- +goose StatementEnd
