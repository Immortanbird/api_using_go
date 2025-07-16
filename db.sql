CREATE TABLE refresh_tokens (
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID(), 1)),

    user_id BINARY(16) NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- The unique ID from the JWT's "jti" claim. This is what you'll use to look it up.
    token_jti TEXT NOT NULL UNIQUE,

    -- A flag to manually revoke a token without deleting the record (good for auditing).
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,

    -- The token's expiration date from the JWT's "exp" claim.
    -- Storing it here allows for efficient cleanup of expired tokens.
    expires_at TIMESTAMP NOT NULL,

    -- The timestamp when this record was created.
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- (Optional but Recommended) Additional metadata for session management.
    ip_address TEXT,
    user_agent TEXT
);

CREATE TABLE IF NOT EXISTS `users` (
	`id` BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID(), 1)),
	`user_name` VARCHAR(20) NOT NULL UNIQUE,
	`email` VARCHAR(128) NOT NULL UNIQUE,
	`pwd_hash` VARCHAR(255) NOT NULL,
	`gender` ENUM('M', 'F', 'X') NOT NULL,
	`birth` TIMESTAMP,
	`address` VARCHAR(255),
	`tel` VARCHAR(32),
    `verified` BOOLEAN DEFAULT FALSE,
	`created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  	`updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);