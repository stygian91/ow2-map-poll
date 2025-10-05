CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  created_at TEXT,
  hash TEXT
);

CREATE INDEX users_hash_index ON users (hash);

create table maps (
  id INTEGER PRIMARY KEY,
  created_at TEXT,
  name TEXT
);

CREATE TABLE polls (
  id INTEGER PRIMARY KEY,
  created_at TEXT,
  user_id INTEGER,
  map1_id INTEGER,
  map2_id INTEGER,
  map3_id INTEGER,
  vote INTEGER
);

CREATE INDEX polls_user_index ON polls (user_id);

-- DATA

INSERT INTO maps (created_at, name) VALUES
  ("2025-10-05 17:45:00", "Aatlis"),
  ("2025-10-05 17:45:00", "Antarctic Peninsula"),
  ("2025-10-05 17:45:00", "Blizzard World"),
  ("2025-10-05 17:45:00", "Colosseo"),
  ("2025-10-05 17:45:00", "Dorado"),
  ("2025-10-05 17:45:00", "Eichenwalde"),
  ("2025-10-05 17:45:00", "Esperanca"),
  ("2025-10-05 17:45:00", "Gibraltar"),
  ("2025-10-05 17:45:00", "Hanaoka"),
  ("2025-10-05 17:45:00", "Havana"),
  ("2025-10-05 17:45:00", "Hollywood"),
  ("2025-10-05 17:45:00", "Illios"),
  ("2025-10-05 17:45:00", "Junkertown"),
  ("2025-10-05 17:45:00", "Kingsrow"),
  ("2025-10-05 17:45:00", "Lijiang"),
  ("2025-10-05 17:45:00", "Midtown"),
  ("2025-10-05 17:45:00", "Circuit royal"),
  ("2025-10-05 17:45:00", "Nepal"),
  ("2025-10-05 17:45:00", "New junk city"),
  ("2025-10-05 17:45:00", "Numbani"),
  ("2025-10-05 17:45:00", "Oasis"),
  ("2025-10-05 17:45:00", "Busan"),
  ("2025-10-05 17:45:00", "Paraiso"),
  ("2025-10-05 17:45:00", "Rialto"),
  ("2025-10-05 17:45:00", "Route66"),
  ("2025-10-05 17:45:00", "Runasapi"),
  ("2025-10-05 17:45:00", "Samoa"),
  ("2025-10-05 17:45:00", "Shambali monastery"),
  ("2025-10-05 17:45:00", "Suravasa"),
  ("2025-10-05 17:45:00", "Throne of anubis"),
  ("2025-10-05 17:45:00", "New queen street");
