-- USER TABLE: Account info and physical preferences
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    preferred_gender_style VARCHAR(50),
    size_top VARCHAR(10),
    size_bottom VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CLOTHING ITEMS TABLE: The digital closet
CREATE TABLE clothing_items (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    category VARCHAR(50) NOT NULL,      
    sub_category VARCHAR(50),           
    primary_color VARCHAR(30) NOT NULL, 
    material VARCHAR(50),               
    min_temp_celsius INT DEFAULT -30,   
    max_temp_celsius INT DEFAULT 50,    
    is_waterproof BOOLEAN DEFAULT FALSE,
    suitable_seasons VARCHAR(20)[], 
    is_trending BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);