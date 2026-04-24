CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,
                          content TEXT NOT NULL,
                          created_at TIMESTAMP DEFAULT now(),
                          post_id INT REFERENCES posts(id) ON DELETE CASCADE,
                          author_id INT REFERENCES authors(id) ON DELETE CASCADE
);