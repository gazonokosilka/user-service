CREATE TABLE IF NOT EXISTS "customer_addresses" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "address" VARCHAR(255) NOT NULL,
    "apartment" VARCHAR(50) NOT NULL,
    "floor" INTEGER NOT NULL,
    "comments" TEXT NOT NULL,
    "customer_id" UUID NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_customer_addresses_customer 
        FOREIGN KEY ("customer_id") 
        REFERENCES "customers"("id") 
        ON DELETE CASCADE
);

