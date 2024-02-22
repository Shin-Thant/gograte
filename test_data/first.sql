-- +gograte Up
CREATE TABLE test (
    name VARCHAR(255) NOT NULL
);

-- +gograte Down
DROP TABLE test;