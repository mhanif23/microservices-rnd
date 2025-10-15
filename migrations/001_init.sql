-- users
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  membership_status TEXT NOT NULL CHECK (membership_status IN ('ACTIVE','SUSPENDED')),
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

-- inventory
CREATE TABLE IF NOT EXISTS inventory (
  book_id UUID PRIMARY KEY,
  total INT NOT NULL CHECK (total >= 0),
  available INT NOT NULL CHECK (available >= 0),
  updated_at TIMESTAMPTZ DEFAULT now()
);

-- borrows
CREATE TABLE IF NOT EXISTS borrows (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  book_id UUID NOT NULL,
  borrowed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  due_at TIMESTAMPTZ NOT NULL,
  returned_at TIMESTAMPTZ NULL,
  status TEXT NOT NULL CHECK (status IN ('BORROWED','RETURNED','LATE')),
  fine_amount NUMERIC(12,2) NOT NULL DEFAULT 0
);

-- outbox (future)
CREATE TABLE IF NOT EXISTS borrow_events (
  id UUID PRIMARY KEY,
  borrow_id UUID NOT NULL,
  type TEXT NOT NULL,
  payload JSONB NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS borrows_user_status_idx ON borrows(user_id, status);
CREATE INDEX IF NOT EXISTS borrows_book_status_idx ON borrows(book_id, status);
