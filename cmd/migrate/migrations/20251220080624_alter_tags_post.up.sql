CREATE TABLE IF NOT EXISTS post_tags (
  post_id bigint NOT NULL,
  tag_id bigint NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

  PRIMARY KEY (post_id, tag_id),

  CONSTRAINT fk_post_tags_post
    FOREIGN KEY (post_id)
    REFERENCES posts (id)
    ON DELETE CASCADE,

  CONSTRAINT fk_post_tags_tag
    FOREIGN KEY (tag_id)
    REFERENCES tags (id)
    ON DELETE CASCADE
);
