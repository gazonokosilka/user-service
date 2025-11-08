CREATE TABLE IF NOT EXISTS "seller_profiles" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid,
    "user_data_id" UUID NOT NULL,
    "inn" INTEGER NOT NULL,
    "full_name" VARCHAR(255) NOT NULL,
    "legal_type" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
) 