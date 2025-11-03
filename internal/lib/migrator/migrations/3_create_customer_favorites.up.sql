CREATE TABLE IF NOT EXISTS "favorites" (
    "product_id" UUID NOT NULL,
    "customer_id" UUID NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY ("product_id", "customer_id"),
    CONSTRAINT fk_favorites_customer 
        FOREIGN KEY ("customer_id") 
        REFERENCES "customers"("id") 
        ON DELETE CASCADE
);
