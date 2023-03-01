DROP TABLE IF EXISTS comic;

CREATE TABLE comic (
  num INTEGER NOT NULL UNIQUE,
  comicDay NUMERIC(2) NOT NULL,
  comicMonth NUMERIC(2) NOT NULL,
  comicYear NUMERIC(4) NOT NULL,
  title BLOB NOT NULL,
  alt BLOB NOT NULL,
  transcript BLOB NOT NULL,
  img VARCHAR(200) characterSET utf8mb4 NOT NULL,
  PRIMARY KEY (`num`)
);

DROP TABLE IF EXISTS tf_idf;

CREATE TABLE tf_idf (
  term TEXT,
  freq INTEGER,
  num INTEGER,
  PRIMARY KEY (`num`)
);