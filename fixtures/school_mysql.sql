CREATE TABLE IF NOT EXISTS school (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS teacher (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    school_id INT,
    FOREIGN KEY (school_id) REFERENCES school(id)
);

CREATE TABLE IF NOT EXISTS class (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name TEXT NOT NULL,
    teacher_id INT,
    school_id INT,
    FOREIGN KEY (teacher_id) REFERENCES teacher(id),
    FOREIGN KEY (school_id) REFERENCES school(id)
);

CREATE TABLE IF NOT EXISTS student (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name TEXT NOT NULL,
    age INT NOT NULL,
    class_id INT,
    school_id INT,
    FOREIGN KEY (class_id) REFERENCES class(id),
    FOREIGN KEY (school_id) REFERENCES school(id)
);

INSERT INTO school (id, name, address) VALUES (1, 'Greenwood High', '123 Elm Street');
INSERT INTO school (id, name, address) VALUES (2, 'Sunnydale School', '456 Oak Avenue');

INSERT INTO teacher (id, name, email, school_id) VALUES (1, 'John Doe', 'johndoe@greenwood.com', 1);
INSERT INTO teacher (id, name, email, school_id) VALUES (2, 'Jane Smith', 'janesmith@sunnydale.com', 2);

INSERT INTO class (id, name, teacher_id, school_id) VALUES (1, 'Math 101', 1, 1);
INSERT INTO class (id, name, teacher_id, school_id) VALUES (2, 'Science 101', 2, 2);

INSERT INTO student (id, name, age, class_id, school_id) VALUES (1, 'Alice', 15, 1, 1);
INSERT INTO student (id, name, age, class_id, school_id) VALUES (2, 'Bob', 16, 1, 1);
INSERT INTO student (id, name, age, class_id, school_id) VALUES (3, 'Charlie', 17, 2, 2);
INSERT INTO student (id, name, age, class_id, school_id) VALUES (4, 'Diana', 16, 2, 2);