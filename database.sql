CREATE TABLE `users` (
  `id`          int(11) NOT NULL AUTO_INCREMENT,
  `type`        tinyint(1) NOT NULL DEFAULT 0,
  `mail`        varchar(100) NOT NULL,
  `password`    varchar(255) NOT NULL,
  `firstName`   varchar(50) NOT NULL,
  `lastName`    varchar(50) NOT NULL,
  `city`        varchar(100),
  `phoneNumber` varchar(20),

  PRIMARY KEY (`id`)
);

CREATE TABLE `coaches` (
  `id`      int(11) NOT NULL AUTO_INCREMENT,
  `address` varchar(255) NOT NULL,

  PRIMARY KEY (`id`),
  FOREIGN KEY (`id`) REFERENCES `users`(`id`)
);

CREATE TABLE `clients` (
  `id`        int(11) NOT NULL AUTO_INCREMENT,
  `birthDate` datetime NOT NULL,
  `coachId`   int(11) DEFAULT NULL,

  PRIMARY KEY (`id`),
  FOREIGN KEY (`id`) REFERENCES `users`(`id`),
  FOREIGN KEY (`coachId`) REFERENCES `users`(`id`)
);

CREATE TABLE `messages` (
  `id`        int(11) NOT NULL AUTO_INCREMENT,
  `sender`    int(11) NOT NULL,
  `receiver`  int(11) NOT NULL,
  `date`      datetime NOT NULL,
  `content`   text NOT NULL,

  PRIMARY KEY (`id`),
  FOREIGN KEY (`sender`) REFERENCES `users`(`id`),
  FOREIGN KEY (`receiver`) REFERENCES `users`(`id`)
);

CREATE TABLE `events` (
  `id`          int(11) NOT NULL AUTO_INCREMENT,
  `name`        varchar(100) NOT NULL,
  `type`        int(1) NOT NULL DEFAULT 0,
  `status`      int(1) NOT NULL DEFAULT 0,
  `firstUser`   int(11) NOT NULL,
  `secondUser`  int(11) NOT NULL,
  `start`       datetime NOT NULL,
  `end`         datetime NOT NULL,
  `created`     datetime NOT NULL,
  `createdBy`   int(11) NOT NULL,
  `updated`     datetime NOT NULL,
  `updatedBy`   int(11) NOT NULL,

  PRIMARY KEY (`id`),
  FOREIGN KEY (`firstUser`) REFERENCES `users`(`id`),
  FOREIGN KEY (`secondUser`) REFERENCES `users`(`id`),
  FOREIGN KEY (`createdBy`) REFERENCES `users`(`id`),
  FOREIGN KEY (`updatedBy`) REFERENCES `users`(`id`)
);

CREATE TABLE `forums` (
  `id`      int(11) NOT NULL AUTO_INCREMENT,
  `title`   varchar(100) NOT NULL,

  PRIMARY KEY (`id`)
);

CREATE TABLE `threads` (
  `id`      int(11) NOT NULL AUTO_INCREMENT,
  `title`   varchar(100) NOT NULL,
  `forumId` int(11) NOT NULL,

  PRIMARY KEY (`id`),
  FOREIGN KEY (`forumId`) REFERENCES `forums`(`id`)
);

CREATE TABLE `posts` (
  `id`       int(11) NOT NULL AUTO_INCREMENT,
  `threadId` int(11) NOT NULL,
  `userId`   int(11) NOT NULL,
  `date`     datetime NOT NULL,
  `content`  text NOT NULL,

  PRIMARY KEY (`id`),
  FOREIGN KEY (`threadId`) REFERENCES `threads`(`id`),
  FOREIGN KEY (`userId`) REFERENCES `users`(`id`)
);