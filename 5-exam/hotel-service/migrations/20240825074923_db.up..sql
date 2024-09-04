CREATE TABLE hotels (
    hotel_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    rating NUMERIC(2, 1),
    address TEXT NOT NULL
);

CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    hotel_id int references hotels(hotel_id), 
    room_type VARCHAR(255),
    price_per_night float,
    availability boolean
)