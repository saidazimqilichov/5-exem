CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    hotel_id int references hotels(hotel_id), 
    room_type VARCHAR(255),
    price_per_night float,
    availability boolean
)