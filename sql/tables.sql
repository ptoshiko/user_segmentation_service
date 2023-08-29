-- Active: 1693057408979@@127.0.0.1@5432@postgres

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users(
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50)
);

-- Таблица сегментов
CREATE TABLE IF NOT EXISTS segments(
    seg_id SERIAL PRIMARY KEY,
    seg_name VARCHAR(50)
);

-- Таблица для связи пользователей и сегментов (многие ко многим)
CREATE TABLE IF NOT EXISTS user_segment (
    user_id INT,
    segment_id INT,
    PRIMARY KEY (user_id, segment_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (segment_id) REFERENCES segments(seg_id)
);


