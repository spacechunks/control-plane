\restrict dbmate

-- Dumped from database version 17.2 (Debian 17.2-1.pgdg120+1)
-- Dumped by pg_dump version 18.6

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
-- Name: server_state; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.server_state AS ENUM (
    'PENDING',
    'RUNNING',
    'UPDATING',
    'DELETING',
    'ERROR'
);


--
-- Name: server_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.server_type AS ENUM (
    'EPHEMERAL',
    'PERSISTENT'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: server; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.server (
    id uuid NOT NULL,
    type public.server_type DEFAULT 'EPHEMERAL'::public.server_type NOT NULL,
    name character varying(256) DEFAULT ''::character varying NOT NULL,
    state public.server_state DEFAULT 'PENDING'::public.server_state NOT NULL,
    labels jsonb DEFAULT '{}'::jsonb NOT NULL,
    ipv4 inet,
    port integer
);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: server server_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.server
    ADD CONSTRAINT server_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

\unrestrict dbmate


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20261031125720');
