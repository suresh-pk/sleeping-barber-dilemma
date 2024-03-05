-- Database: barbershop

-- DROP DATABASE IF EXISTS barbershop;

CREATE DATABASE barbershop
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'English_United States.1252'
    LC_CTYPE = 'English_United States.1252'
    LOCALE_PROVIDER = 'libc'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

-- Table: public.barbers

-- DROP TABLE IF EXISTS public.barbers;

CREATE TABLE IF NOT EXISTS public.barbers
(
    id integer NOT NULL DEFAULT nextval('barber_id_seq'::regclass),
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    is_available boolean NOT NULL DEFAULT true,
    CONSTRAINT barber_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.barbers
    OWNER to postgres;

INSERT INTO barber (name) VALUES
('Barber1'),
('Barber2'),
('Barber3');
