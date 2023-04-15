SET FOREIGN_KEY_CHECKS=0; 
DROP TABLE IF EXISTS comics; 
SET FOREIGN_KEY_CHECKS=1;
CREATE TABLE comics (
  num INTEGER NOT NULL UNIQUE PRIMARY KEY,
  image_name TEXT, 
  title TEXT, 
  title_raw TEXT, 
  alt_text TEXT, 
  alt_text_raw TEXT, 
  transcript LONGTEXT, 
  transcript_raw LONGTEXT, 
  explanation LONGTEXT, 
  explanation_raw LONGTEXT, 
  incomplete BOOL
);

DROP TABLE IF EXISTS term_frequency;
CREATE TABLE term_frequency (
  comic_num INTEGER,
  term TEXT, -- stemmed term
  terms_raw LONGTEXT, -- string of raw terms that have same stem
  freq INTEGER,
  FOREIGN KEY (comic_num) REFERENCES comics(num)
);

DROP TABLE IF EXISTS comic_frequency;
CREATE TABLE comic_frequency (
  term VARCHAR(50),
  freq INTEGER
);