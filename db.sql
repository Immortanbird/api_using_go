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
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- (Optional but Recommended) Additional metadata for session management.
    ip_address TEXT,
    user_agent TEXT
);

CREATE TABLE `users` (
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
  	`updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE images (
    `id` BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID(), 1)),
    `user_id` BINARY(16) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    `original_filename` TEXT NOT NULL,
    `hashed_filename` TEXT NOT NULL UNIQUE,
    `image_url` TEXT NOT NULL UNIQUE,
    `download_url` TEXT NOT NULL UNIQUE,
    `size_bytes` BIGINT NOT NULL,
    `mime_type` VARCHAR(128) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts (
    -- The unique identifier for each post, using UUID.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign key linking the post to an author in the `users` table.
    -- If a user is deleted, all their posts will be deleted automatically.
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- The title of the blog post.
    title VARCHAR(255) NOT NULL,

    -- A URL-friendly version of the title (e.g., "my-first-post").
    -- This should be unique to prevent duplicate URLs.
    slug TEXT NOT NULL UNIQUE,

    -- A short summary or excerpt of the post, can be used for previews.
    excerpt TEXT,

    -- The main content of the post, stored as Markdown or HTML.
    -- TEXT type allows for very long articles.
    content TEXT NOT NULL,

    -- The URL of a cover image for the post.
    cover_image_url VARCHAR(255),

    -- The publication status of the post (e.g., 'draft', 'published', 'archived').
    -- This allows you to save posts without making them public immediately.
    status VARCHAR(50) NOT NULL DEFAULT 'draft',

    -- Timestamps for when the post was created and last updated.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes on commonly queried columns for better performance.
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_slug ON posts(slug);
CREATE INDEX idx_posts_status ON posts(status);