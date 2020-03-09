DROP TABLE IF EXISTS news;

CREATE TABLE news (
    id serial primary key,
    title varchar(128),
    url varchar(256) unique,
    date_added timestamp default now()
);

-- SELECT * FROM news LIMIT 10 OFFSET 19;
-- 
-- SELECT count(*) FROM news;
