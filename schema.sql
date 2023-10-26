-- entities
CREATE TABLE room (
  id TEXT NOT NULL PRIMARY KEY,
  judge_capacity INTEGER NOT NULL
);

CREATE TABLE event (
  id TEXT NOT NULL PRIMARY KEY
);

CREATE TABLE student (
  email TEXT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE judge ( 
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL
);

CREATE TABLE judge_event_rel (
  judge_id INTEGER NOT NULL,
  event_id TEXT NOT NULL,

  FOREIGN KEY (judge_id) REFERENCES judge (id),
  FOREIGN KEY (event_id) REFERENCES event (id),
  PRIMARY KEY (judge_id, event_id) ON CONFLICT REPLACE
);

-- config
CREATE TABLE time_slot ( 
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  config_id INTEGER NOT NULL,
  purpose TEXT,
  duration INTEGER NOT NULL,

  FOREIGN KEY (config_id) REFERENCES config (id)
);

CREATE TABLE config (
  id INTEGER PRIMARY KEY,
  time_start TEXT NOT NULL,
  max_group_size INTEGER NOT NULL,
  exam_length INTEGER NOT NULL
);

-- event requests
CREATE TABLE student_group (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  event_id TEXT NOT NULL,

  FOREIGN KEY (event_id) REFERENCES event (id)
);

CREATE TABLE student_student_group_rel (
  student_group_id INTEGER NOT NULL,
  student_email TEXT NOT NULL,

  FOREIGN KEY (student_group_id) REFERENCES student_group (id),
  FOREIGN KEY (student_email) REFERENCES student (email),
  PRIMARY KEY (student_group_id, student_email) ON CONFLICT REPLACE
);
