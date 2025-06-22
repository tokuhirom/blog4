CREATE TABLE entry
(
    path         varchar(255) CHARACTER SET ascii COLLATE ascii_general_ci     NOT NULL,
    title        varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    body         text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci         NOT NULL,
    visibility   enum ('private','public')                                              DEFAULT 'private' NOT NULL,
    format       enum ('html','mkdn')                                          NOT NULL DEFAULT 'mkdn',
    published_at datetime                                                               DEFAULT NULL,
    last_edited_at datetime                                                               DEFAULT CURRENT_TIMESTAMP comment 'last manualy edited at',
    created_at   datetime                                                               DEFAULT CURRENT_TIMESTAMP,
    updated_at   datetime                                                               DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (path),
    UNIQUE title (title),
    KEY created_at (created_at),
    KEY updated_at (updated_at),
    KEY published_at (published_at),
    KEY last_edited_at (last_edited_at),
    FULLTEXT KEY idx_bigram (title, body)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

create table entry_image
(
    path       varchar(255) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
    url        varchar(5000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    created_at datetime                                                       DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (path),
    KEY created_at (created_at),
    FOREIGN KEY (path) REFERENCES entry (path) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

create table entry_link
(
    src_path  varchar(255) CHARACTER SET ascii COLLATE ascii_general_ci     NOT NULL,
    dst_title varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    PRIMARY KEY (src_path, dst_title),
    FOREIGN KEY (src_path) REFERENCES entry (path) ON DELETE CASCADE,
    INDEX (dst_title)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

create table amazon_cache
(
    asin             varchar(255) primary key,
    title            varchar(5000) default null,
    image_medium_url text          default null,
    link             varchar(5000) not null,
    created_at       datetime      DEFAULT CURRENT_TIMESTAMP,
    KEY created_at (created_at)
) engine = innodb
  default charset = utf8mb4;

CREATE TABLE admin_session
(
    session_id VARCHAR(255) CHARACTER SET ascii COLLATE ascii_general_ci PRIMARY KEY,
    username VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_accessed_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
