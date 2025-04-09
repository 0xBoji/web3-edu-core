CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    thumbnail VARCHAR(255),
    instructor_id UUID REFERENCES users(id),
    price DECIMAL(10, 2),
    level VARCHAR(50), -- beginner, intermediate, advanced
    duration INTEGER, -- total minutes
    category VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
