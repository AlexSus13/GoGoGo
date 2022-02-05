CREATE TABLE headers (filename character varying(255) not null, content_type character varying(255) not null, content_lenght character varying(255) not null, id serial not null unique);
CREATE TABLE user_data (user_name character varying(50) not null, user_password character varying(50) not null);
