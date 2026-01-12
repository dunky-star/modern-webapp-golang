CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name  VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    access_level INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rooms (
    id BIGSERIAL PRIMARY KEY,
    room_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE reservations (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name  VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    phone      VARCHAR(255),
    start_date DATE NOT NULL,
    end_date   DATE NOT NULL,
    room_id    BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_reservations_room
        FOREIGN KEY (room_id)
        REFERENCES rooms(id)
        ON DELETE CASCADE
);


CREATE TABLE restrictions (
    id BIGSERIAL PRIMARY KEY,
    restriction_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE room_restrictions (
    id BIGSERIAL PRIMARY KEY,
    start_date DATE NOT NULL,
    end_date   DATE NOT NULL,
    room_id BIGINT NOT NULL,
    reservation_id BIGINT,
    restriction_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_rr_room
        FOREIGN KEY (room_id)
        REFERENCES rooms(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_rr_reservation
        FOREIGN KEY (reservation_id)
        REFERENCES reservations(id)
        ON DELETE SET NULL,

    CONSTRAINT fk_rr_restriction
        FOREIGN KEY (restriction_id)
        REFERENCES restrictions(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_reservations_room_id ON reservations(room_id);
CREATE INDEX idx_reservations_email ON reservations(email);
CREATE INDEX idx_reservations_last_name ON reservations(last_name);
CREATE INDEX idx_room_restrictions_room_id ON room_restrictions(room_id);
CREATE INDEX idx_room_restrictions_reservation_id ON room_restrictions(reservation_id);
CREATE INDEX idx_room_restrictions_restriction_id ON room_restrictions(restriction_id);
CREATE INDEX idx_room_restrictions_start_date ON room_restrictions(start_date);
CREATE INDEX idx_room_restrictions_end_date ON room_restrictions(end_date);

