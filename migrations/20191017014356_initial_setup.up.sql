CREATE TABLE IF NOT EXISTS chats (
  id VARCHAR(255) PRIMARY KEY,
  default_bucket VARCHAR(255)
);
CREATE TABLE IF NOT EXISTS buckets (
  id VARCHAR(255),
  chat_id VARCHAR(255) REFERENCES chats(id),
  PRIMARY KEY (chat_id, id)
);
CREATE TABLE IF NOT EXISTS options (
  id SERIAL PRIMARY KEY,
  chat_id VARCHAR(255) REFERENCES chats(id),
  bucket_id VARCHAR(255),
  content VARCHAR(255),
  FOREIGN KEY (chat_id, bucket_id) REFERENCES buckets(chat_id, id)
);
