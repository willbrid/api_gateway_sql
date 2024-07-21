-- Create tables
CREATE TABLE IF NOT EXISTS school (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS teacher (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    school_id INTEGER,
    FOREIGN KEY (school_id) REFERENCES school(id)
);

CREATE TABLE IF NOT EXISTS class (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    teacher_id INTEGER,
    school_id INTEGER,
    FOREIGN KEY (teacher_id) REFERENCES teacher(id),
    FOREIGN KEY (school_id) REFERENCES school(id)
);

CREATE TABLE IF NOT EXISTS student (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    age INTEGER NOT NULL,
    class_id INTEGER,
    school_id INTEGER,
    FOREIGN KEY (class_id) REFERENCES class(id),
    FOREIGN KEY (school_id) REFERENCES school(id)
);

-- Insert data into school
INSERT INTO school (name, address) VALUES ('Greenwood High', '123 Elm Street');
INSERT INTO school (name, address) VALUES ('Sunnydale School', '456 Oak Avenue');

-- Insert data into teacher
INSERT INTO teacher (name, email, school_id) VALUES ('John Doe', 'johndoe@greenwood.com', 1);
INSERT INTO teacher (name, email, school_id) VALUES ('Jane Smith', 'janesmith@sunnydale.com', 2);

-- Insert data into class
INSERT INTO class (name, teacher_id, school_id) VALUES ('Math 101', 1, 1);
INSERT INTO class (name, teacher_id, school_id) VALUES ('Science 101', 2, 2);

-- Insert data into student
INSERT INTO student (name, age, class_id, school_id) VALUES ('Alice', 15, 1, 1);
INSERT INTO student (name, age, class_id, school_id) VALUES ('Bob', 16, 1, 1);
INSERT INTO student (name, age, class_id, school_id) VALUES ('Charlie', 17, 2, 2);
INSERT INTO student (name, age, class_id, school_id) VALUES ('Diana', 16, 2, 2);