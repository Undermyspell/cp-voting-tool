Create TABLE Sessions (
    id VARCHAR(20) PRIMARY KEY,
    sessionSecret VARCHAR(30) NOT NULL
);

CREATE TABLE Questions (
    id VARCHAR(26) PRIMARY KEY,
    sessionId VARCHAR(20) NOT NULL,
    text VARCHAR(500) NOT NULL,
    votes INTEGER NOT NULL,
    answered BOOLEAN NOT NULL,
    anonymous BOOLEAN NOT NULL,
    creatorName VARCHAR(50),
    creatorHash VARCHAR(64) NOT NULL,
    FOREIGN KEY (sessionId) REFERENCES Sessions (id) ON DELETE CASCADE
);

CREATE INDEX idx_sessionId_name ON Questions (sessionId, id);

CREATE TABLE UserVotes (
	questionId VARCHAR(26) NOT NULL,
	userHash VARCHAR(64) NOT NULL,
	PRIMARY KEY (questionId, userHash),
	FOREIGN KEY (questionId) REFERENCES Questions (id) ON DELETE CASCADE
);