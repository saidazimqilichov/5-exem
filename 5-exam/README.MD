User register
http://localhost:9001/users/register

{
    "age": 19,
    "email": "qilichovsaidazim006@gmail.com",
    "name": "saidazim",
    "password": "54321"
}

User verify
http://localhost:9001/users/verify

{
    "email":"qilichovsaidazim006@gmail.com",
    "code": "yuborilgan code"
}

User login
http://localhost:9001/users/login


{
    "email": "qilichovsaidazim006@gmail.com",
    "password": "54321"
}


Get user with id
http://localhost:9001/users/:id



Delete user with id
http://localhost:9001/users/:id





Create hotel
http://localhost:8080/hotels

{
    "address": "laborum",
    "location": "aliqua",
    "name": "exercitation eiusmod Lorem ad",
    "rating": 5,
    "rooms": [
        {
            "pricePerNight": 234,
            "availability": true,
            "roomType": "standart"
        },
        {
            "pricePerNight": 567,
            "availability": true,
            "roomType": "lux"
        }
    ]
}

Get  hotel with by id
http://localhost:9001/hotels/:id


Get all hotels
http://localhost:9001/hotels



available rooms
HTTP://localhost:9001/hotels/:id/check


create booking
http://localhost:9001/booking

{
    "user_id": 2,
    "hotel_id": 4,
    "roomtype": "standart",
    "check_in_date": {
        "seconds": 1640995200,
        "nanos": 0
    },
    "check_out_date": {
        "seconds": 1642211600,
        "nanos": 0
    },
    "total_amount": 250.75
}

Get booking by id
http://localhost9001/booking/:id 


Update booking by id
http://localhost9001/booking/:id
{
    "booking_id": 67,
    "check_in_date": {
        "nanos": -24758650,
        "seconds": "8427977"
    },
    "check_out_date": {
        "nanos": 1796217347,
        "seconds": "11672827"
    },
    "room_type": "mollit non"
}
