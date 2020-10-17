CREATE TABLE log (
  sha TEXT,
  date TIMESTAMP,
  author TEXT,
  added INT,
  removed INT,
  filename TEXT,

  PRIMARY KEY (sha, filename)
);

CREATE INDEX ix_date ON log (date DESC);

CREATE INDEX ix_author on log (author ASC);
