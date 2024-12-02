CREATE TABLE prequal_parameters_response (
    id SERIAL PRIMARY KEY,                        
    max_life_time INT NOT NULL,                   
    pool_size INT NOT NULL,                       
    probe_factor FLOAT NOT NULL,                  
    probe_remove_factor INT NOT NULL,           
    mu INT NOT NULL,                            
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);
