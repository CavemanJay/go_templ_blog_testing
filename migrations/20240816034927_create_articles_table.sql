-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists articles(
  id integer primary key,
  created_at timestamp not null,
  updated_at timestamp not null,
  title varchar(255) not null,
  slug varchar(255) not null,
  filename varchar(255) not null
);

insert into articles (created_at, updated_at, title, slug, filename) 
values 
  (datetime('now'), datetime('now'), 'Article 1', 'article-1', 'article-1'),
  (datetime('now'), datetime('now'), 'Article 2', 'article-2', 'article-2');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists articles;
-- +goose StatementEnd
