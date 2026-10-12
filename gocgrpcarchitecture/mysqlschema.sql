#############################################
###                                      ####
### Author: Austine I Owoh               ####
### Version: 1.0                         ####
### License: aIEIn                       ####
### Date: August Sunday 24 2026           ####
####                                     ####
#############################################

CREATE SCHEMA IF NOT EXISTS gocgrpcarchitecture

SET NAMES 'UTF8MB4';
SET TIME_ZONE = 'US/Pacific';
SET TIME ZONE = '3:00';

USE gocgrpcarchitecture;

DROP TABLE IF EXISTS Users;

CREATE TABLE Users
(
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    UUID       VARCHAR(50) NOT NULL,
    UniqueId    VARCHAR(50) NOT NULL,
    FirstName      VARCHAR(50) NOT NULL,
    LastName       VARCHAR(50) NOT NULL,
    Email           VARCHAR(100) NOT NULL,
    Password        VARCHAR(255) DEFAULT NULL,
    Phone           VARCHAR(25) DEFAULT NULL,
    Biograph            VARCHAR(100) DEFAULT NULL,
    EmailToken VARCHAR(255) DEFAULT NULL,
    TOTPFactorSid VARCHAR(255) DEFAULT NULL,
    QrCodeSecret VARCHAR(50) DEFAULT NULL,
    QrCodeImageUri VARCHAR(255) DEFAULT 'https://cdn-icons-png.flaticon.com/512/149/149071.png',
    QrCodeImage VARCHAR(255) DEFAULT NULL,
    FailedAttempts int DEFAULT 0,
    LastLoginDeviceName VARCHAR(50) DEFAULT NULL,
    LastLoginIpAddress VARCHAR(50) DEFAULT NULL,
    UsingMfa boolean DEFAULT FALSE,
    EmailVerified boolean DEFAULT FALSE,
    SMSVerified boolean DEFAULT FALSE,
    SMSEnabled boolean DEFAULT FALSE,
    TOTPEnabled boolean DEFAULT FALSE,
    AccountEnabled boolean DEFAULT FALSE,
    AccountDisabled boolean DEFAULT FALSE,
    AccountNonExpired boolean DEFAULT FALSE,
    AccountNonLocked boolean DEFAULT FALSE,
    LastLoginTimestamp DATETIME DEFAULT NULL,
    LastFailedAttemptTime DATETIME DEFAULT NULL,
    PasswordLastChangedAt DATETIME DEFAULT NULL,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT UQ_Users_Email UNIQUE (Email),
    CONSTRAINT UQ_Users_User_UUId UNIQUE (UUID),
    CONSTRAINT UQ_Users_User_Unique_Id UNIQUE (UniqueId)
);

DROP TABLE IF EXISTS Roles;
CREATE TABLE Roles
(
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    Name         VARCHAR(255) NOT NULL,
    Permissions  VARCHAR(255) NOT NULL,
    CreatedAt   DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt   DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT   UQ_Roles_Name UNIQUE (Name)
);

DROP TABLE IF EXISTS UserRoles;
CREATE TABLE UserRoles
(
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    UserId     BIGINT UNSIGNED NOT NULL,
    RoleId     BIGINT UNSIGNED NOT NULL,
    CONSTRAINT  UQ_UserRoles_User_Id UNIQUE (UserId),
    CONSTRAINT  FK_User_Roles_User_Id FOREIGN KEY (UserId) REFERENCES Users (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT  FK_User_Roles_Role_Id FOREIGN KEY (RoleId) REFERENCES Roles (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE RESTRICT
);



=================================================

DROP TABLE IF EXISTS EmailVerifications;
CREATE TABLE EmailVerifications
(
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id      BIGINT UNSIGNED NOT NULL,
    url          VARCHAR(255) NOT NULL,
    expiration_at  DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT UQ_AccountVerifications_User_Id UNIQUE (user_id),
    CONSTRAINT UQ_AccountVerifications_Url UNIQUE (url)
);

DROP TABLE IF EXISTS ResetPasswordVerifications;
CREATE TABLE ResetPasswordVerifications
(
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED NOT NULL,
    url             VARCHAR(255) NOT NULL,
    expiration_at   DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT UQ_ResetPasswordVerifications_User_Id UNIQUE (user_id),
    CONSTRAINT UQ_ResetPasswordVerifications_Url UNIQUE (url)
);

DROP TABLE IF EXISTS TwoFactorVerifications;
CREATE TABLE TwoFactorVerifications
(
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED NOT NULL,
    code            VARCHAR(10) NOT NULL,
    expiration_at   DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT UQ_TwoFactorVerifications_User_Id UNIQUE (user_id),
    CONSTRAINT UQ_TwoFactorVerifications_Code UNIQUE (code)
);

DROP TABLE IF EXISTS Events;
CREATE TABLE Events
(
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    type         VARCHAR(50) NOT NULL CHECK(type IN ('LOGIN_ATTEMPT', 'LOGIN_ATTEMPT_FAILURE', 'LOGIN_ATTEMPT_SUCCESS', 'PROFILE_UPDATE', 'PROFILE_PICTURE_UPDATE', 'ROLE_UPDATE', 'ACCOUNT_SETTINGS_UPDATE', 'PASSWORD_UPDATE', 'MFA_UPDATE')),
    description  VARCHAR(255) NOT NULL,
    CONSTRAINT UQ_Events_Type UNIQUE (type)
);

DROP TABLE IF EXISTS UserEvents;
CREATE TABLE UserEvents
(
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id      BIGINT UNSIGNED NOT NULL,
    event_id     BIGINT UNSIGNED NOT NULL,
    device       VARCHAR(100) DEFAULT NULL,
    ip_address   VARCHAR(100) DEFAULT NULL,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_UserEvents_User_Id FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_UserEvents_Event_Id FOREIGN KEY (event_id) REFERENCES Events (id) ON DELETE RESTRICT ON UPDATE CASCADE
);