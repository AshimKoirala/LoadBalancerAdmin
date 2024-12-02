CREATE TABLE activity_logs (
    id SERIAL PRIMARY KEY,
    type VARCHAR(10) NOT NULL CHECK (type IN ('success', 'warning', 'error')),
    message TEXT NOT NULL,
    replica_id INT NULL REFERENCES replicas(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
