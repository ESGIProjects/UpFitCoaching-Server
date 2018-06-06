START TRANSACTION;

-- Users
INSERT INTO users (type, mail, password, firstName, lastName, city, phoneNumber) VALUES (1, 'jasonpierna@icloud.com', 'motdepasse', 'Jason', 'Pierna', 'Garges-lès-Gonesse', '0123456789');
INSERT INTO users (type, mail, password, firstName, lastName, city, phoneNumber) VALUES (1, 'kevints.le@gmail.com', 'motdepasse', 'Kévin', 'Le', 'Ermont', '0123456789');
INSERT INTO users (type, mail, password, firstName, lastName, city, phoneNumber) VALUES (1, 'maeva.malih@gmail.com', 'motdepasse', 'Maëva', 'Malih', 'Saint-Brice', '0123456789');
INSERT INTO users (type, mail, password, firstName, lastName, city, phoneNumber) VALUES (2, 'coach@test.fr', 'coachtest', 'Coach', 'Super', 'Cupertino', '0123456789');

-- Coaches
INSERT INTO coaches (id, address) VALUES (4, '1, Infinite Loop');

-- Clients
INSERT INTO clients (id, birthDate) VALUES (1, '1995-08-07');
INSERT INTO clients (id, birthDate) VALUES (2, '1994-12-29');
INSERT INTO clients (id, birthDate) VALUES (3, '1994-02-04');


INSERT INTO messages (sender, receiver, date, content) VALUES (4, 1, CURRENT_TIMESTAMP, 'Salut Jason !');
INSERT INTO messages (sender, receiver, date, content) VALUES (1, 4, CURRENT_TIMESTAMP, 'Hey !');

INSERT INTO messages (sender, receiver, date, content) VALUES (4, 2, CURRENT_TIMESTAMP, 'Coucou Kévin!');
INSERT INTO messages (sender, receiver, date, content) VALUES (2, 4, CURRENT_TIMESTAMP, 'Non.');

INSERT INTO messages (sender, receiver, date, content) VALUES (4, 3, CURRENT_TIMESTAMP, 'Hello Maeva !');
INSERT INTO messages (sender, receiver, date, content) VALUES (3, 4, CURRENT_TIMESTAMP, 'Ça va ?');

COMMIT;