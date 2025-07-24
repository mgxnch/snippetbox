# Snippetbox

Snippetbox is an app that displays short snippets of 100 characters or less.

This is my attempt at following through with the book `Let's Go` by Alex Edwards.

## Database

### Setup

```sql
-- Create database
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create a `snippets` table
CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

-- Add an index on the created column
CREATE INDEX idx_snippets_created ON snippets(created);

-- To show indices
SHOW index from snippets;
```

### Seed data

```sql
INSERT INTO snippets (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n- Matsuo Basho',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry forest\n forest, winds howl in rage\nwith no leaves to blow.\n\n- Natsume Soseki',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n- Murakami Kijo',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);
```

### Creating a new `mysql` user

```sql
CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';

-- Randomly generated string: openssl rand -base64 18
-- In real world systems, we would never commit our password
ALTER USER 'web'@'localhost' IDENTIFIED BY '9mfOz8RWTWQSIlgt8hX9jb9V';
```

This user has restricted privileges and will be used by the web application.

## API

| Method | Pattern           | Handler           | Action                                     |
|--------|-------------------|-------------------|--------------------------------------------|
| GET    | /                 | home              | Display a home page                        |
| GET    | /snippet/view/:id | snippetView       | Display a specific snippet                 |
| GET    | /snippet/create   | snippetCreate     | Display a HTML form for creating a snippet |
| POST   | /snippet/create   | snippetCreatePost | Create a new snippet                       |
| GET    | /static/*         | http.FileServer   | Serve a specific static file               |
