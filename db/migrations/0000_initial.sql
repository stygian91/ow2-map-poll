CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  created_at TEXT,
  hash TEXT
);

CREATE INDEX users_hash_index ON users (hash);

create table maps (
  id INTEGER PRIMARY KEY,
  created_at TEXT,
  mode TEXT,
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

INSERT INTO maps (created_at, mode, name) VALUES
  ("2025-10-05 17:45:00", "flashpoint", "Aatlis"),
  ("2025-10-05 17:45:00", "control", "Antarctic Peninsula"),
  ("2025-10-05 17:45:00", "hybrid", "Blizzard World"),
  ("2025-10-05 17:45:00", "push", "Colosseo"),
  ("2025-10-05 17:45:00", "escort", "Dorado"),
  ("2025-10-05 17:45:00", "hybrid", "Eichenwalde"),
  ("2025-10-05 17:45:00", "push", "Esperanca"),
  ("2025-10-05 17:45:00", "escort", "Gibraltar"),
  ("2025-10-05 17:45:00", "clash", "Hanaoka"),
  ("2025-10-05 17:45:00", "escort", "Havana"),
  ("2025-10-05 17:45:00", "hybrid", "Hollywood"),
  ("2025-10-05 17:45:00", "control", "Illios"),
  ("2025-10-05 17:45:00", "escort", "Junkertown"),
  ("2025-10-05 17:45:00", "hybrid", "Kingsrow"),
  ("2025-10-05 17:45:00", "control", "Lijiang"),
  ("2025-10-05 17:45:00", "hybrid", "Midtown"),
  ("2025-10-05 17:45:00", "escort", "Circuit royal"),
  ("2025-10-05 17:45:00", "control", "Nepal"),
  ("2025-10-05 17:45:00", "flashpoint", "New junk city"),
  ("2025-10-05 17:45:00", "hybrid", "Numbani"),
  ("2025-10-05 17:45:00", "control", "Oasis"),
  ("2025-10-05 17:45:00", "control", "Busan"),
  ("2025-10-05 17:45:00", "hybrid", "Paraiso"),
  ("2025-10-05 17:45:00", "escort", "Rialto"),
  ("2025-10-05 17:45:00", "escort", "Route 66"),
  ("2025-10-05 17:45:00", "push", "Runasapi"),
  ("2025-10-05 17:45:00", "control", "Samoa"),
  ("2025-10-05 17:45:00", "escort", "Shambali monastery"),
  ("2025-10-05 17:45:00", "flashpoint", "Suravasa"),
  ("2025-10-05 17:45:00", "clash", "Throne of anubis"),
  ("2025-10-05 17:45:00", "push", "New queen street");
