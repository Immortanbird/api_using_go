CREATE TABLE refresh_tokens (
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID(), 1)),

    user_id BINARY(16) NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- The unique ID from the JWT's "jti" claim. This is what you'll use to look it up.
    token_jti CHAR(36) NOT NULL UNIQUE,

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
    `id` BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID(), 1)),
    `user_id` BINARY(16) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    `title` VARCHAR(255) NOT NULL,
    
    -- A URL-friendly version of the title (e.g., "my-first-post").
    -- This should be unique to prevent duplicate URLs.
    `slug` VARCHAR(255) NOT NULL UNIQUE,
    
    -- A short summary or excerpt of the post, can be used for previews.
    `excerpt` TEXT,

    `content_markdown` LONGTEXT NOT NULL,
    `content_html` LONGTEXT NOT NULL,
    
    -- Cover image
    `cover_image_url` VARCHAR(500),
    
    `status` ENUM('draft', 'published', 'archived', 'private') NOT NULL DEFAULT 'draft',
    
    -- Publishing control
    `published_at` TIMESTAMP NULL,
    `scheduled_at` TIMESTAMP NULL,
    
    -- Engagement metrics
    `view_count` BIGINT DEFAULT 0,
    `like_count` BIGINT DEFAULT 0,
    `comment_count` BIGINT DEFAULT 0,
    
    -- Reading time estimation (in minutes)
    `reading_time` INT DEFAULT 0,
    
    -- Soft delete support
    `deleted_at` TIMESTAMP NULL,
    
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create indexes on commonly queried columns for better performance.
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_slug ON posts(slug);
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_published_at ON posts(published_at);
CREATE INDEX idx_posts_status_published_at ON posts(status, published_at);
CREATE INDEX idx_posts_deleted_at ON posts(deleted_at);

-- Full-text search index for content
CREATE FULLTEXT INDEX idx_posts_content_search ON posts(title, excerpt, content_markdown);

-- Create a posts_tags junction table for many-to-many relationship
CREATE TABLE posts_tags (
    `id` BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID(), 1)),
    `post_id` BINARY(16) NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    `tag_id` BINARY(16) NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_post_tag (post_id, tag_id)
);

-- Create tags table
CREATE TABLE tags (
    `id` BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID(), 1)),
    `name` VARCHAR(100) NOT NULL UNIQUE,
    `slug` VARCHAR(100) NOT NULL UNIQUE,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create indexes for tags
CREATE INDEX idx_tags_slug ON tags(slug);
CREATE INDEX idx_posts_tags_post_id ON posts_tags(post_id);
CREATE INDEX idx_posts_tags_tag_id ON posts_tags(tag_id);

-- Create comments table
CREATE TABLE comments (
    `id` BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID(), 1)),
    `post_id` BINARY(16) NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    `user_id` BINARY(16) NULL REFERENCES users(id) ON DELETE SET NULL, -- Allow anonymous comments
    `parent_id` BINARY(16) NULL REFERENCES comments(id) ON DELETE CASCADE, -- For nested comments
    `author_name` VARCHAR(100) NOT NULL, -- For anonymous comments
    `content` TEXT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create indexes for comments
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_created_at ON comments(created_at);