-- Cheers co-pilot i am way to lazy to write this schema by hand
USE gorm;

CREATE TABLE company (
    company_id CHAR(36) PRIMARY KEY,
    company_name VARCHAR(255)
);

CREATE TABLE jobs (
    uuid CHAR(36) PRIMARY KEY,
    job_title VARCHAR(255),
    company_id CHAR(36),
    location VARCHAR(255),
    salary DECIMAL(10, 2),
    posted_date DATE,
    FOREIGN KEY (company_id) REFERENCES company(company_id)
);

CREATE TABLE files (
    file_id CHAR(36) PRIMARY KEY,
    file_name VARCHAR(255),
    job_id CHAR(36),
    FOREIGN KEY (job_id) REFERENCES jobs(uuid)
);

-- Insert 300 entries into the company table
DELIMITER //
CREATE PROCEDURE insert_companies()
BEGIN
    DECLARE i INT DEFAULT 1;
    WHILE i <= 300 DO
        INSERT INTO company (company_id, company_name) 
        VALUES (UUID(), CONCAT('Company ', i));
        SET i = i + 1;
    END WHILE;
END //
DELIMITER ;

CALL insert_companies();

-- Insert 300 entries into the jobs table
DELIMITER //
CREATE PROCEDURE insert_jobs()
BEGIN
    DECLARE i INT DEFAULT 1;
    WHILE i <= 300 DO
        INSERT INTO jobs (uuid, job_title, company_id, location, salary, posted_date) 
        VALUES (UUID(), CONCAT('Job Title ', i), 
                (SELECT company_id FROM company ORDER BY RAND() LIMIT 1), 
                CONCAT('Location ', i), (RAND() * 100000), 
                CURDATE() - INTERVAL (i % 30) DAY);
        SET i = i + 1;
    END WHILE;
END //
DELIMITER ;

CALL insert_jobs();

-- Insert 300 entries into the files table
DELIMITER //
CREATE PROCEDURE insert_files()
BEGIN
    DECLARE i INT DEFAULT 1;
    WHILE i <= 300 DO
        INSERT INTO files (file_id, file_name, job_id) 
        VALUES (UUID(), CONCAT('File ', i), 
                (SELECT uuid FROM jobs ORDER BY RAND() LIMIT 1));
        SET i = i + 1;
    END WHILE;
END //
DELIMITER ;

CALL insert_files();
