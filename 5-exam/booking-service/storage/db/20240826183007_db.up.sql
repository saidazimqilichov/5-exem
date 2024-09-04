create table hotel_book(
    booking_id serial primary key,
    user_id int,
    hotel_id int,
    room_type text,
    check_in_date date,
    check_out_date date,
    total_amount float,
    status varchar(64)
);