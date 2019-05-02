-- rambler up

CREATE TABLE keys (
	id UUID PRIMARY KEY,
	alg TEXT,
	owner TEXT,
	priv BYTEA
);
CREATE INDEX keys_alg_idx ON keys (alg);

-- rambler down

DROP INDEX keys_alg_idx;
DROP TABLE keys;
