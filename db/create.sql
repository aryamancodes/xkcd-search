DROP TABLE IF EXISTS Comics;
CREATE TABLE Comics (
  id INTEGER NOT NULL UNIQUE,
  comicDay NUMERIC(2) NOT NULL,
  comicMonth NUMERIC(2) NOT NULL,
  comicYear NUMERIC(4) NOT NULL,
  title BLOB NOT NULL,
  alt BLOB NOT NULL,
  transcript BLOB,
  img VARCHAR(200) characterSET utf8mb4 NOT NULL,
  PRIMARY KEY (`num`)
);

DROP TABLE IF EXISTS TermFrequency;
CREATE TABLE TermFrequency (
  comic_id INTEGER, 
  term TEXT,
  frequency INTEGER,
  UNIQUE(comic_id, term)
  FOREIGN KEY (`num`) REFERENCES Comics(id) 
);