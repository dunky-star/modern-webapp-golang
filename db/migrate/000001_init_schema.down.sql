DROP INDEX IF EXISTS idx_room_restrictions_restriction_id;
DROP INDEX IF EXISTS idx_room_restrictions_reservation_id;
DROP INDEX IF EXISTS idx_room_restrictions_room_id;
DROP INDEX IF EXISTS idx_reservations_room_id;

DROP TABLE IF EXISTS room_restrictions;
DROP TABLE IF EXISTS restrictions;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS users;
