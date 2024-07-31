CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE quizzes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    is_activated BOOLEAN DEFAULT FALSE
);

CREATE TABLE scores (
    user_id INT NOT NULL,
    quiz_id INT NOT NULL,
    score INT NOT NULL,
    PRIMARY KEY (user_id, quiz_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id)
);



INSERT INTO users (name) VALUES
('Alice'), ('Bob');

INSERT INTO quizzes (name, is_activated) VALUES
('Quiz A', TRUE), ('Quiz B', FALSE);