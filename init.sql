--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4
-- Dumped by pg_dump version 17.4

-- Started on 2025-04-22 10:22:29 MSK

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 5 (class 2615 OID 50645)
-- Name: studio; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA studio;


--
-- TOC entry 235 (class 1255 OID 50646)
-- Name: lowercase_before_insert(); Type: FUNCTION; Schema: studio; Owner: -
--

CREATE FUNCTION studio.lowercase_before_insert() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.login := LOWER(NEW.login);
    RETURN NEW;
END;
$$;


--
-- TOC entry 236 (class 1255 OID 50647)
-- Name: lowercase_before_update(); Type: FUNCTION; Schema: studio; Owner: -
--

CREATE FUNCTION studio.lowercase_before_update() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.login := LOWER(NEW.login);
    RETURN NEW;
END;
$$;


--
-- TOC entry 237 (class 1255 OID 58834)
-- Name: update_release_date(); Type: FUNCTION; Schema: studio; Owner: -
--

CREATE FUNCTION studio.update_release_date() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.status = 3 THEN
        NEW.release_date := CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 217 (class 1259 OID 50650)
-- Name: customers; Type: TABLE; Schema: studio; Owner: -
--

CREATE TABLE studio.customers (
    id integer NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    login text NOT NULL
);


--
-- TOC entry 218 (class 1259 OID 50655)
-- Name: customers_id_seq; Type: SEQUENCE; Schema: studio; Owner: -
--

CREATE SEQUENCE studio.customers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3541 (class 0 OID 0)
-- Dependencies: 218
-- Name: customers_id_seq; Type: SEQUENCE OWNED BY; Schema: studio; Owner: -
--

ALTER SEQUENCE studio.customers_id_seq OWNED BY studio.customers.id;


--
-- TOC entry 219 (class 1259 OID 50656)
-- Name: employees; Type: TABLE; Schema: studio; Owner: -
--

CREATE TABLE studio.employees (
    id integer NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    job_id integer NOT NULL,
    login text NOT NULL
);


--
-- TOC entry 220 (class 1259 OID 50661)
-- Name: employees_id_seq; Type: SEQUENCE; Schema: studio; Owner: -
--

CREATE SEQUENCE studio.employees_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3542 (class 0 OID 0)
-- Dependencies: 220
-- Name: employees_id_seq; Type: SEQUENCE OWNED BY; Schema: studio; Owner: -
--

ALTER SEQUENCE studio.employees_id_seq OWNED BY studio.employees.id;


--
-- TOC entry 221 (class 1259 OID 50662)
-- Name: job_types; Type: TABLE; Schema: studio; Owner: -
--

CREATE TABLE studio.job_types (
    id integer NOT NULL,
    title text NOT NULL
);


--
-- TOC entry 222 (class 1259 OID 50667)
-- Name: job_types_id_seq; Type: SEQUENCE; Schema: studio; Owner: -
--

CREATE SEQUENCE studio.job_types_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3543 (class 0 OID 0)
-- Dependencies: 222
-- Name: job_types_id_seq; Type: SEQUENCE OWNED BY; Schema: studio; Owner: -
--

ALTER SEQUENCE studio.job_types_id_seq OWNED BY studio.job_types.id;


--
-- TOC entry 223 (class 1259 OID 50668)
-- Name: materials; Type: TABLE; Schema: studio; Owner: -
--

CREATE TABLE studio.materials (
    id integer NOT NULL,
    title text NOT NULL,
    price real NOT NULL
);


--
-- TOC entry 224 (class 1259 OID 50673)
-- Name: materials_id_seq; Type: SEQUENCE; Schema: studio; Owner: -
--

CREATE SEQUENCE studio.materials_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3544 (class 0 OID 0)
-- Dependencies: 224
-- Name: materials_id_seq; Type: SEQUENCE OWNED BY; Schema: studio; Owner: -
--

ALTER SEQUENCE studio.materials_id_seq OWNED BY studio.materials.id;


--
-- TOC entry 225 (class 1259 OID 50674)
-- Name: model_materials; Type: TABLE; Schema: studio; Owner: -
--

CREATE TABLE studio.model_materials (
    model_id integer NOT NULL,
    material_id integer NOT NULL,
    leng real NOT NULL
);


--
-- TOC entry 226 (class 1259 OID 50677)
-- Name: models; Type: TABLE; Schema: studio; Owner: -
--

CREATE TABLE studio.models (
    id integer NOT NULL,
    title text NOT NULL,
    price real
);


--
-- TOC entry 227 (class 1259 OID 50682)
-- Name: models_id_seq; Type: SEQUENCE; Schema: studio; Owner: -
--

CREATE SEQUENCE studio.models_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3545 (class 0 OID 0)
-- Dependencies: 227
-- Name: models_id_seq; Type: SEQUENCE OWNED BY; Schema: studio; Owner: -
--

ALTER SEQUENCE studio.models_id_seq OWNED BY studio.models.id;


--
-- TOC entry 228 (class 1259 OID 50683)
-- Name: order_items; Type: TABLE; Schema: studio; Owner: -
--

CREATE TABLE studio.order_items (
    id integer NOT NULL,
    o_id integer NOT NULL,
    model integer NOT NULL,
    unit_price real NOT NULL
);


--
-- TOC entry 229 (class 1259 OID 50686)
-- Name: order_items_id_seq; Type: SEQUENCE; Schema: studio; Owner: -
--

CREATE SEQUENCE studio.order_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3546 (class 0 OID 0)
-- Dependencies: 229
-- Name: order_items_id_seq; Type: SEQUENCE OWNED BY; Schema: studio; Owner: -
--

ALTER SEQUENCE studio.order_items_id_seq OWNED BY studio.order_items.id;


--
-- TOC entry 230 (class 1259 OID 50687)
-- Name: order_status; Type: TABLE; Schema: studio; Owner: -
--

CREATE TABLE studio.order_status (
    id integer NOT NULL,
    status text NOT NULL
);


--
-- TOC entry 231 (class 1259 OID 50692)
-- Name: order_status_id_seq; Type: SEQUENCE; Schema: studio; Owner: -
--

CREATE SEQUENCE studio.order_status_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3547 (class 0 OID 0)
-- Dependencies: 231
-- Name: order_status_id_seq; Type: SEQUENCE OWNED BY; Schema: studio; Owner: -
--

ALTER SEQUENCE studio.order_status_id_seq OWNED BY studio.order_status.id;


--
-- TOC entry 232 (class 1259 OID 50693)
-- Name: orders; Type: TABLE; Schema: studio; Owner: -
--

CREATE TABLE studio.orders (
    id integer NOT NULL,
    c_id integer NOT NULL,
    e_id integer,
    status integer DEFAULT 1 NOT NULL,
    create_date timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    release_date timestamp with time zone
);


--
-- TOC entry 233 (class 1259 OID 50698)
-- Name: orders_id_seq; Type: SEQUENCE; Schema: studio; Owner: -
--

CREATE SEQUENCE studio.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3548 (class 0 OID 0)
-- Dependencies: 233
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: studio; Owner: -
--

ALTER SEQUENCE studio.orders_id_seq OWNED BY studio.orders.id;


--
-- TOC entry 234 (class 1259 OID 50699)
-- Name: users; Type: TABLE; Schema: studio; Owner: -
--

CREATE TABLE studio.users (
    login text NOT NULL,
    access_level integer NOT NULL
);


--
-- TOC entry 3322 (class 2604 OID 50794)
-- Name: customers id; Type: DEFAULT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.customers ALTER COLUMN id SET DEFAULT nextval('studio.customers_id_seq'::regclass);


--
-- TOC entry 3323 (class 2604 OID 50795)
-- Name: employees id; Type: DEFAULT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.employees ALTER COLUMN id SET DEFAULT nextval('studio.employees_id_seq'::regclass);


--
-- TOC entry 3324 (class 2604 OID 50796)
-- Name: job_types id; Type: DEFAULT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.job_types ALTER COLUMN id SET DEFAULT nextval('studio.job_types_id_seq'::regclass);


--
-- TOC entry 3325 (class 2604 OID 50797)
-- Name: materials id; Type: DEFAULT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.materials ALTER COLUMN id SET DEFAULT nextval('studio.materials_id_seq'::regclass);


--
-- TOC entry 3326 (class 2604 OID 50798)
-- Name: models id; Type: DEFAULT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.models ALTER COLUMN id SET DEFAULT nextval('studio.models_id_seq'::regclass);


--
-- TOC entry 3327 (class 2604 OID 50799)
-- Name: order_items id; Type: DEFAULT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.order_items ALTER COLUMN id SET DEFAULT nextval('studio.order_items_id_seq'::regclass);


--
-- TOC entry 3328 (class 2604 OID 50800)
-- Name: order_status id; Type: DEFAULT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.order_status ALTER COLUMN id SET DEFAULT nextval('studio.order_status_id_seq'::regclass);


--
-- TOC entry 3329 (class 2604 OID 50801)
-- Name: orders id; Type: DEFAULT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.orders ALTER COLUMN id SET DEFAULT nextval('studio.orders_id_seq'::regclass);


--
-- TOC entry 3518 (class 0 OID 50650)
-- Dependencies: 217
-- Data for Name: customers; Type: TABLE DATA; Schema: studio; Owner: -
--

INSERT INTO studio.customers VALUES (1, 'Федор', 'Миронов', 'miron');
INSERT INTO studio.customers VALUES (2, 'Михаил', 'Прохоров', 'proh');
INSERT INTO studio.customers VALUES (3, 'Михаил', 'Коновалов', 'konovalovma');
INSERT INTO studio.customers VALUES (4, 'Леонид', 'Савельев', 'savleon');
INSERT INTO studio.customers VALUES (5, 'МОеимя', 'САмсон', 'test');


--
-- TOC entry 3520 (class 0 OID 50656)
-- Dependencies: 219
-- Data for Name: employees; Type: TABLE DATA; Schema: studio; Owner: -
--

INSERT INTO studio.employees VALUES (1, 'Никита', 'Михайлов', 2, 'nicmic');


--
-- TOC entry 3522 (class 0 OID 50662)
-- Dependencies: 221
-- Data for Name: job_types; Type: TABLE DATA; Schema: studio; Owner: -
--

INSERT INTO studio.job_types VALUES (2, 'OPERATOR');
INSERT INTO studio.job_types VALUES (3, 'SYSADM');


--
-- TOC entry 3524 (class 0 OID 50668)
-- Dependencies: 223
-- Data for Name: materials; Type: TABLE DATA; Schema: studio; Owner: -
--

INSERT INTO studio.materials VALUES (1, 'Хлопок', 450);
INSERT INTO studio.materials VALUES (2, 'Шерсть', 600);
INSERT INTO studio.materials VALUES (3, 'Вискоза', 320);
INSERT INTO studio.materials VALUES (4, 'Полиэстер', 300);
INSERT INTO studio.materials VALUES (5, 'Полушерсть', 160);
INSERT INTO studio.materials VALUES (6, 'Полихлопок', 180);
INSERT INTO studio.materials VALUES (8, 'Иск. кожа', 500);
INSERT INTO studio.materials VALUES (9, 'Габардин', 1200);


--
-- TOC entry 3526 (class 0 OID 50674)
-- Dependencies: 225
-- Data for Name: model_materials; Type: TABLE DATA; Schema: studio; Owner: -
--

INSERT INTO studio.model_materials VALUES (1, 1, 3);
INSERT INTO studio.model_materials VALUES (2, 2, 2);
INSERT INTO studio.model_materials VALUES (2, 3, 2);
INSERT INTO studio.model_materials VALUES (3, 2, 3);
INSERT INTO studio.model_materials VALUES (3, 3, 3);
INSERT INTO studio.model_materials VALUES (4, 1, 3);
INSERT INTO studio.model_materials VALUES (5, 1, 3);
INSERT INTO studio.model_materials VALUES (11, 1, 5);
INSERT INTO studio.model_materials VALUES (11, 5, 5);
INSERT INTO studio.model_materials VALUES (12, 2, 3);
INSERT INTO studio.model_materials VALUES (13, 8, 1);
INSERT INTO studio.model_materials VALUES (14, 8, 4);
INSERT INTO studio.model_materials VALUES (14, 3, 4);
INSERT INTO studio.model_materials VALUES (15, 2, 3);
INSERT INTO studio.model_materials VALUES (16, 2, 3);
INSERT INTO studio.model_materials VALUES (16, 3, 3);
INSERT INTO studio.model_materials VALUES (17, 9, 5);
INSERT INTO studio.model_materials VALUES (17, 3, 4.5);


--
-- TOC entry 3527 (class 0 OID 50677)
-- Dependencies: 226
-- Data for Name: models; Type: TABLE DATA; Schema: studio; Owner: -
--

INSERT INTO studio.models VALUES (1, 'Рубашка Мод.1', 2000);
INSERT INTO studio.models VALUES (2, 'Брюки Мод.1', 4000);
INSERT INTO studio.models VALUES (3, 'Пиджак Мод.1', 6000);
INSERT INTO studio.models VALUES (4, 'Рубашка в клетку', 2000);
INSERT INTO studio.models VALUES (5, 'Рубашка в полоску', 2000);
INSERT INTO studio.models VALUES (11, 'Пальто (черное)', 44000);
INSERT INTO studio.models VALUES (12, 'Шерстяной свитер', 2500);
INSERT INTO studio.models VALUES (13, 'Кожаный ремень', 1200);
INSERT INTO studio.models VALUES (14, 'Кожаная куртка', 25000);
INSERT INTO studio.models VALUES (15, 'Джемпер', 4000);
INSERT INTO studio.models VALUES (16, 'Жилет Мод.1', 2400);
INSERT INTO studio.models VALUES (17, 'Плащ', 18000);


--
-- TOC entry 3529 (class 0 OID 50683)
-- Dependencies: 228
-- Data for Name: order_items; Type: TABLE DATA; Schema: studio; Owner: -
--

INSERT INTO studio.order_items VALUES (1, 1, 1, 2000);
INSERT INTO studio.order_items VALUES (2, 1, 2, 4000);
INSERT INTO studio.order_items VALUES (3, 2, 3, 6000);
INSERT INTO studio.order_items VALUES (4, 3, 2, 4000);
INSERT INTO studio.order_items VALUES (5, 3, 3, 6000);
INSERT INTO studio.order_items VALUES (11, 7, 4, 2000);
INSERT INTO studio.order_items VALUES (12, 7, 5, 2000);
INSERT INTO studio.order_items VALUES (13, 8, 4, 2000);
INSERT INTO studio.order_items VALUES (14, 8, 5, 2000);
INSERT INTO studio.order_items VALUES (15, 9, 1, 2000);
INSERT INTO studio.order_items VALUES (16, 9, 2, 4000);
INSERT INTO studio.order_items VALUES (17, 9, 3, 6000);
INSERT INTO studio.order_items VALUES (18, 9, 4, 2000);
INSERT INTO studio.order_items VALUES (19, 9, 5, 2000);
INSERT INTO studio.order_items VALUES (20, 10, 2, 4000);
INSERT INTO studio.order_items VALUES (21, 10, 11, 44000);
INSERT INTO studio.order_items VALUES (22, 10, 16, 2400);
INSERT INTO studio.order_items VALUES (23, 11, 17, 18000);
INSERT INTO studio.order_items VALUES (24, 12, 2, 4000);
INSERT INTO studio.order_items VALUES (25, 12, 14, 25000);
INSERT INTO studio.order_items VALUES (26, 12, 16, 2400);
INSERT INTO studio.order_items VALUES (7, 13, 2, 4000);
INSERT INTO studio.order_items VALUES (8, 13, 12, 2500);
INSERT INTO studio.order_items VALUES (9, 13, 13, 1200);
INSERT INTO studio.order_items VALUES (27, 14, 5, 2000);
INSERT INTO studio.order_items VALUES (28, 14, 12, 2500);
INSERT INTO studio.order_items VALUES (29, 15, 11, 44000);


--
-- TOC entry 3531 (class 0 OID 50687)
-- Dependencies: 230
-- Data for Name: order_status; Type: TABLE DATA; Schema: studio; Owner: -
--

INSERT INTO studio.order_status VALUES (1, 'Ожидает');
INSERT INTO studio.order_status VALUES (2, 'На исполнении');
INSERT INTO studio.order_status VALUES (3, 'Выдан');
INSERT INTO studio.order_status VALUES (4, 'Отменен');


--
-- TOC entry 3533 (class 0 OID 50693)
-- Dependencies: 232
-- Data for Name: orders; Type: TABLE DATA; Schema: studio; Owner: -
--

INSERT INTO studio.orders VALUES (2, 2, NULL, 4, '2025-03-25 22:06:55+03', NULL);
INSERT INTO studio.orders VALUES (3, 2, 1, 2, '2025-03-25 22:12:42+03', NULL);
INSERT INTO studio.orders VALUES (9, 4, NULL, 4, '2025-04-12 02:11:14+03', NULL);
INSERT INTO studio.orders VALUES (12, 2, NULL, 1, '2025-04-13 22:52:31+03', NULL);
INSERT INTO studio.orders VALUES (11, 4, 1, 2, '2025-04-12 20:28:14+03', NULL);
INSERT INTO studio.orders VALUES (1, 1, 1, 3, '2025-03-12 19:52:29+03', '2025-03-27 22:55:38+03');
INSERT INTO studio.orders VALUES (7, 3, 1, 3, '2025-04-12 00:45:06+03', '2025-04-12 01:53:56+03');
INSERT INTO studio.orders VALUES (8, 4, 1, 3, '2025-04-12 01:54:20+03', '2025-04-12 02:10:41+03');
INSERT INTO studio.orders VALUES (10, 3, 1, 3, '2025-04-12 14:20:50+03', '2025-04-12 14:23:09+03');
INSERT INTO studio.orders VALUES (14, 4, 1, 2, '2025-04-19 16:48:46.235362+03', NULL);
INSERT INTO studio.orders VALUES (15, 2, NULL, 1, '2025-04-20 00:41:47.042173+03', NULL);
INSERT INTO studio.orders VALUES (13, 1, 1, 3, '2025-04-14 23:10:59+03', '2025-04-20 00:44:16.09385+03');


--
-- TOC entry 3535 (class 0 OID 50699)
-- Dependencies: 234
-- Data for Name: users; Type: TABLE DATA; Schema: studio; Owner: -
--

INSERT INTO studio.users VALUES ('miron', 1);
INSERT INTO studio.users VALUES ('nicmic', 2);
INSERT INTO studio.users VALUES ('proh', 1);
INSERT INTO studio.users VALUES ('konovalovma', 1);
INSERT INTO studio.users VALUES ('savleon', 1);
INSERT INTO studio.users VALUES ('test', 1);
INSERT INTO studio.users VALUES ('rest', 1);


--
-- TOC entry 3549 (class 0 OID 0)
-- Dependencies: 218
-- Name: customers_id_seq; Type: SEQUENCE SET; Schema: studio; Owner: -
--

SELECT pg_catalog.setval('studio.customers_id_seq', 5, true);


--
-- TOC entry 3550 (class 0 OID 0)
-- Dependencies: 220
-- Name: employees_id_seq; Type: SEQUENCE SET; Schema: studio; Owner: -
--

SELECT pg_catalog.setval('studio.employees_id_seq', 1, true);


--
-- TOC entry 3551 (class 0 OID 0)
-- Dependencies: 222
-- Name: job_types_id_seq; Type: SEQUENCE SET; Schema: studio; Owner: -
--

SELECT pg_catalog.setval('studio.job_types_id_seq', 3, true);


--
-- TOC entry 3552 (class 0 OID 0)
-- Dependencies: 224
-- Name: materials_id_seq; Type: SEQUENCE SET; Schema: studio; Owner: -
--

SELECT pg_catalog.setval('studio.materials_id_seq', 9, true);


--
-- TOC entry 3553 (class 0 OID 0)
-- Dependencies: 227
-- Name: models_id_seq; Type: SEQUENCE SET; Schema: studio; Owner: -
--

SELECT pg_catalog.setval('studio.models_id_seq', 17, true);


--
-- TOC entry 3554 (class 0 OID 0)
-- Dependencies: 229
-- Name: order_items_id_seq; Type: SEQUENCE SET; Schema: studio; Owner: -
--

SELECT pg_catalog.setval('studio.order_items_id_seq', 29, true);


--
-- TOC entry 3555 (class 0 OID 0)
-- Dependencies: 231
-- Name: order_status_id_seq; Type: SEQUENCE SET; Schema: studio; Owner: -
--

SELECT pg_catalog.setval('studio.order_status_id_seq', 4, true);


--
-- TOC entry 3556 (class 0 OID 0)
-- Dependencies: 233
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: studio; Owner: -
--

SELECT pg_catalog.setval('studio.orders_id_seq', 15, true);


--
-- TOC entry 3333 (class 2606 OID 50713)
-- Name: customers customers_login_key; Type: CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.customers
    ADD CONSTRAINT customers_login_key UNIQUE (login);


--
-- TOC entry 3335 (class 2606 OID 50715)
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);


--
-- TOC entry 3337 (class 2606 OID 50717)
-- Name: employees employees_login_key; Type: CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.employees
    ADD CONSTRAINT employees_login_key UNIQUE (login);


--
-- TOC entry 3339 (class 2606 OID 50719)
-- Name: employees employees_pkey; Type: CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.employees
    ADD CONSTRAINT employees_pkey PRIMARY KEY (id);


--
-- TOC entry 3341 (class 2606 OID 50721)
-- Name: job_types job_types_pkey; Type: CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.job_types
    ADD CONSTRAINT job_types_pkey PRIMARY KEY (id);


--
-- TOC entry 3343 (class 2606 OID 50723)
-- Name: materials materials_pkey; Type: CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.materials
    ADD CONSTRAINT materials_pkey PRIMARY KEY (id);


--
-- TOC entry 3345 (class 2606 OID 50725)
-- Name: materials materials_title_key; Type: CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.materials
    ADD CONSTRAINT materials_title_key UNIQUE (title);


--
-- TOC entry 3347 (class 2606 OID 50727)
-- Name: model_materials model_materials_pkey; Type: CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.model_materials
    ADD CONSTRAINT model_materials_pkey PRIMARY KEY (model_id, material_id);


--
-- TOC entry 3349 (class 2606 OID 50729)
-- Name: models models_pkey; Type: CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.models
    ADD CONSTRAINT models_pkey PRIMARY KEY (id);


--
-- TOC entry 3351 (class 2606 OID 50731)
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- TOC entry 3353 (class 2606 OID 50733)
-- Name: order_status order_status_pkey; Type: CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.order_status
    ADD CONSTRAINT order_status_pkey PRIMARY KEY (id);


--
-- TOC entry 3355 (class 2606 OID 50735)
-- Name: order_status order_status_status_key; Type: CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.order_status
    ADD CONSTRAINT order_status_status_key UNIQUE (status);


--
-- TOC entry 3357 (class 2606 OID 50737)
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- TOC entry 3359 (class 2606 OID 50739)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (login);


--
-- TOC entry 3371 (class 2620 OID 50740)
-- Name: users lowercase_before_insert; Type: TRIGGER; Schema: studio; Owner: -
--

CREATE TRIGGER lowercase_before_insert BEFORE INSERT ON studio.users FOR EACH ROW EXECUTE FUNCTION studio.lowercase_before_insert();


--
-- TOC entry 3372 (class 2620 OID 50741)
-- Name: users lowercase_before_update; Type: TRIGGER; Schema: studio; Owner: -
--

CREATE TRIGGER lowercase_before_update BEFORE UPDATE ON studio.users FOR EACH ROW EXECUTE FUNCTION studio.lowercase_before_update();


--
-- TOC entry 3370 (class 2620 OID 58836)
-- Name: orders trg_update_release_date; Type: TRIGGER; Schema: studio; Owner: -
--

CREATE TRIGGER trg_update_release_date BEFORE UPDATE ON studio.orders FOR EACH ROW EXECUTE FUNCTION studio.update_release_date();


--
-- TOC entry 3360 (class 2606 OID 50807)
-- Name: customers customers_login_fkey; Type: FK CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.customers
    ADD CONSTRAINT customers_login_fkey FOREIGN KEY (login) REFERENCES studio.users(login) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3361 (class 2606 OID 50749)
-- Name: employees employees_job_id_fkey; Type: FK CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.employees
    ADD CONSTRAINT employees_job_id_fkey FOREIGN KEY (job_id) REFERENCES studio.job_types(id);


--
-- TOC entry 3362 (class 2606 OID 50754)
-- Name: employees employees_login_fkey; Type: FK CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.employees
    ADD CONSTRAINT employees_login_fkey FOREIGN KEY (login) REFERENCES studio.users(login) ON UPDATE CASCADE;


--
-- TOC entry 3363 (class 2606 OID 50759)
-- Name: model_materials model_materials_material_id_fkey; Type: FK CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.model_materials
    ADD CONSTRAINT model_materials_material_id_fkey FOREIGN KEY (material_id) REFERENCES studio.materials(id) ON DELETE CASCADE;


--
-- TOC entry 3364 (class 2606 OID 50764)
-- Name: model_materials model_materials_model_id_fkey; Type: FK CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.model_materials
    ADD CONSTRAINT model_materials_model_id_fkey FOREIGN KEY (model_id) REFERENCES studio.models(id) ON DELETE CASCADE;


--
-- TOC entry 3365 (class 2606 OID 50769)
-- Name: order_items order_items_model_fkey; Type: FK CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.order_items
    ADD CONSTRAINT order_items_model_fkey FOREIGN KEY (model) REFERENCES studio.models(id);


--
-- TOC entry 3366 (class 2606 OID 50802)
-- Name: order_items order_items_o_id_fkey; Type: FK CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.order_items
    ADD CONSTRAINT order_items_o_id_fkey FOREIGN KEY (o_id) REFERENCES studio.orders(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3367 (class 2606 OID 50779)
-- Name: orders orders_c_id_fkey; Type: FK CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.orders
    ADD CONSTRAINT orders_c_id_fkey FOREIGN KEY (c_id) REFERENCES studio.customers(id);


--
-- TOC entry 3368 (class 2606 OID 50784)
-- Name: orders orders_e_id_fkey; Type: FK CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.orders
    ADD CONSTRAINT orders_e_id_fkey FOREIGN KEY (e_id) REFERENCES studio.employees(id);


--
-- TOC entry 3369 (class 2606 OID 50789)
-- Name: orders orders_status_fkey; Type: FK CONSTRAINT; Schema: studio; Owner: -
--

ALTER TABLE ONLY studio.orders
    ADD CONSTRAINT orders_status_fkey FOREIGN KEY (status) REFERENCES studio.order_status(id);


-- Completed on 2025-04-22 10:22:29 MSK

--
-- PostgreSQL database dump complete
--

