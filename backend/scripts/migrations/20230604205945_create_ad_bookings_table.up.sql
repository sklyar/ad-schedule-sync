CREATE TABLE ad_bookings
(
    id          SERIAL PRIMARY KEY,
    client_name TEXT        NOT NULL CHECK ( length(client_name) <= 255 ),
    booking_at  TIMESTAMPTZ NOT NULL,
    vk_post_id  TEXT        NULL,
    created_at  TIMESTAMPTZ NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL,
    deleted_at  TIMESTAMPTZ NULL
);
