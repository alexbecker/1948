CREATE TABLE users (
    rowid INTEGER NOT NULL PRIMARY KEY,
    username varchar(100) UNIQUE NOT NULL,
    hash varchar(100) NOT NULL
);

CREATE TABLE user_inheritance (
    parentid INTEGER NOT NULL REFERENCES users(rowid),
    childid INTEGER NOT NULL REFERENCES users(rowid)
);

CREATE TABLE user_roles (
    userid INTEGER NOT NULL REFERENCES users(rowid),
    role varchar(100) NOT NULL
);
