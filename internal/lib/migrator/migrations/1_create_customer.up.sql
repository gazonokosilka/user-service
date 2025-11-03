CREATE TABLE IF NOT EXISTS "customers" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "first_name" VARCHAR(100) NOT NULL,
    "last_name" VARCHAR(100) NOT NULL,
    "gender" VARCHAR(20) NOT NULL CHECK (gender IN ('male', 'female')),
    "timezone" VARCHAR(50) NOT NULL DEFAULT 'UTC',
    "birthday" DATE NOT NULL,
    "user_id" UUID NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
