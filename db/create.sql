DROP TABLE IF EXISTS comic;
CREATE TABLE comic (
  id         INT AUTO_INCREMENT NOT NULL,
  num        SMALLINT NOT NULL,
  day        NUMERIC(2) NOT NULL,
  month      NUMERIC(2) NOT NULL,
  year       NUMERIC(4) NOT NULL,
  title      BLOB NOT NULL,
  alt        BLOB NOT NULL, 
  transcript BLOB NOT NULL,
  img        VARCHAR(200) CHARACTER SET utf8mb4 NOT NULL,
  PRIMARY KEY (`id`)
);
