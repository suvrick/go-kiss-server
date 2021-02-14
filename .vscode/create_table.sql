
 /*
    login_time, 
    login_id,
    social,
    token, 
    token2, 
    user_id, 
    game_id, 
    name,
    profile,
    avatar,
    balance
 */

CREATE TABLE users(
    id serial,
    login_time varchar(200),
    login_id bigint,
    social integer,
    token text,
    token2 text,
    user_id varchar(50) unique,
    game_id integer,
    name text,
    profile text,
    avatar text,
    balance integer
);