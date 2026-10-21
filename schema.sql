-- Pembuatan manual database dan tabel (soal Database no. 1).
-- Setara dengan migrations/000001_create_posts.up.sql.
CREATE DATABASE IF NOT EXISTS article;
USE article;

CREATE TABLE IF NOT EXISTS posts (
  id INT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(200) NOT NULL,
  content TEXT NOT NULL,
  category VARCHAR(100) NOT NULL,
  created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  status VARCHAR(100) NOT NULL -- publish | draft | thrash
);
