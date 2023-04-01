-- source <path/to/file> in mysql to create the db
SET FOREIGN_KEY_CHECKS=0; 
DROP TABLE IF EXISTS comics; 
SET FOREIGN_KEY_CHECKS=1;
CREATE TABLE comics (
  num INTEGER NOT NULL UNIQUE PRIMARY KEY,
  title TEXT, 
  alt_text TEXT, 
  transcript LONGTEXT, 
  explanation LONGTEXT, 
  incomplete BOOL
);

DROP TABLE IF EXISTS term_frequency;
CREATE TABLE term_frequency (
  comic_num INTEGER,
  term TEXT,
  freq INTEGER,
  FOREIGN KEY (comic_num) REFERENCES comics(num)
);

DROP TABLE IF EXISTS comic_frequency;
CREATE TABLE comic_frequency (
  id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
  term TEXT,
  freq INTEGER
);