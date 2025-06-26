CREATE TABLE refresh_tokens (
    -- The unique ID for this specific token/session record.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign key to your users table. If a user is deleted, their tokens are too.
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- The unique ID from the JWT's "jti" claim. This is what you'll use to look it up.
    token_jti TEXT NOT NULL UNIQUE,

    -- A flag to manually revoke a token without deleting the record (good for auditing).
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,

    -- The token's expiration date from the JWT's "exp" claim.
    -- Storing it here allows for efficient cleanup of expired tokens.
    expires_at TIMESTAMPTZ NOT NULL,

    -- The timestamp when this record was created.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- (Optional but Recommended) Additional metadata for session management.
    ip_address TEXT,
    user_agent TEXT
);

-- Create an index on the JTI for fast lookups.
CREATE INDEX idx_refresh_tokens_jti ON refresh_tokens(token_jti);


CREATE TABLE IF NOT EXISTS `users` (
	`uid` BINARY(16) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
	`user_name` VARCHAR(20) NOT NULL UNIQUE,
	`email` VARCHAR(128) NOT NULL UNIQUE,
	`pwd_hash` VARCHAR(255) NOT NULL,
	`gender` ENUM('M', 'F', 'X') NOT NULL,
	`birth` TIMESTAMP,
	`address` VARCHAR(255),
	`tel` VARCHAR(32),
	`created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  	`updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	`deleted_at` TIMESTAMP NULL DEFAULT NULL
);