CREATE TABLE school (
    id INT IDENTITY(1,1) PRIMARY KEY,
    name NVARCHAR(100) NOT NULL,
    address NVARCHAR(255) NOT NULL
);

CREATE TABLE teacher (
    id INT IDENTITY(1,1) PRIMARY KEY,
    name NVARCHAR(100) NOT NULL,
    email NVARCHAR(255) NOT NULL,
    school_id INT,
    FOREIGN KEY (school_id) REFERENCES school(id)
);

CREATE TABLE class (
    id INT IDENTITY(1,1) PRIMARY KEY,
    name NVARCHAR(100) NOT NULL,
    teacher_id INT,
    school_id INT,
    FOREIGN KEY (teacher_id) REFERENCES teacher(id),
    FOREIGN KEY (school_id) REFERENCES school(id)
);

CREATE TABLE student (
    id INT IDENTITY(1,1) PRIMARY KEY,
    name NVARCHAR(100) NOT NULL,
    age INT NOT NULL,
    class_id INT,
    school_id INT,
    FOREIGN KEY (class_id) REFERENCES class(id),
    FOREIGN KEY (school_id) REFERENCES school(id)
);

SET IDENTITY_INSERT school ON;
INSERT INTO school (id, name, address) VALUES (1, 'Greenwood High', '123 Elm Street');
INSERT INTO school (id, name, address) VALUES (2, 'Sunnydale School', '456 Oak Avenue');
SET IDENTITY_INSERT school OFF;

SET IDENTITY_INSERT teacher ON;
INSERT INTO teacher (id, name, email, school_id) VALUES (1, 'John Doe', 'johndoe@greenwood.com', 1);
INSERT INTO teacher (id, name, email, school_id) VALUES (2, 'Jane Smith', 'janesmith@sunnydale.com', 2);
SET IDENTITY_INSERT teacher OFF;

SET IDENTITY_INSERT class ON;
INSERT INTO class (id, name, teacher_id, school_id) VALUES (1, 'Math 101', 1, 1);
INSERT INTO class (id, name, teacher_id, school_id) VALUES (2, 'Science 101', 2, 2);
SET IDENTITY_INSERT class OFF;

SET IDENTITY_INSERT student ON;
INSERT INTO student (id, name, age, class_id, school_id) VALUES (1, 'Alice', 15, 1, 1);
INSERT INTO student (id, name, age, class_id, school_id) VALUES (2, 'Bob', 16, 1, 1);
INSERT INTO student (id, name, age, class_id, school_id) VALUES (3, 'Charlie', 17, 2, 2);
INSERT INTO student (id, name, age, class_id, school_id) VALUES (4, 'Diana', 16, 2, 2);
SET IDENTITY_INSERT student OFF;