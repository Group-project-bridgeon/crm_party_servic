--
-- PostgreSQL database dump
--

\restrict kIHiL32QxJnkOtZuEgeMMA3h4pvfN6d0Cw6SugYUFLTnPFk1RuQZgjtQdYAtMYv

-- Dumped from database version 16.15 (Debian 16.15-1.pgdg13+2)
-- Dumped by pg_dump version 16.15 (Debian 16.15-1.pgdg13+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: customer_user_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.customer_user_accounts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    customer_id uuid NOT NULL,
    user_id uuid NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: customers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.customers (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    party_id uuid NOT NULL,
    customer_code character varying(50) NOT NULL,
    customer_group character varying(100),
    price_list_id uuid,
    credit_limit numeric(14,2) DEFAULT 0 NOT NULL,
    payment_terms_days integer DEFAULT 0 NOT NULL,
    credit_status character varying(30) DEFAULT 'ACTIVE'::character varying NOT NULL,
    sales_rep_id uuid,
    tax_exempt boolean DEFAULT false NOT NULL,
    tax_exemption_number character varying(100),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_customers_credit_limit CHECK ((credit_limit >= (0)::numeric)),
    CONSTRAINT chk_customers_credit_status CHECK (((credit_status)::text = ANY ((ARRAY['ACTIVE'::character varying, 'BLOCKED'::character varying, 'HOLD'::character varying])::text[]))),
    CONSTRAINT chk_customers_payment_terms CHECK ((payment_terms_days >= 0))
);


--
-- Name: individuals; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.individuals (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    party_id uuid NOT NULL,
    first_name character varying(100) NOT NULL,
    middle_name character varying(100),
    last_name character varying(100),
    date_of_birth date,
    gender character varying(30),
    job_title character varying(150),
    national_id character varying(100),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: organizations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.organizations (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    party_id uuid NOT NULL,
    registration_number character varying(100),
    industry character varying(150),
    organization_type character varying(50),
    founded_date date,
    employee_count integer,
    description text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_organizations_employee_count CHECK (((employee_count IS NULL) OR (employee_count >= 0)))
);


--
-- Name: parties; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.parties (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    party_type character varying(30) NOT NULL,
    display_name character varying(200) NOT NULL,
    legal_name character varying(255),
    email character varying(255),
    phone character varying(30),
    website character varying(255),
    tax_number character varying(100),
    notes text,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_parties_party_type CHECK (((party_type)::text = ANY ((ARRAY['ORGANIZATION'::character varying, 'INDIVIDUAL'::character varying])::text[])))
);


--
-- Name: party_addresses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.party_addresses (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    party_id uuid NOT NULL,
    address_type character varying(30) NOT NULL,
    address_name character varying(100),
    address_line1 character varying(255) NOT NULL,
    address_line2 character varying(255),
    landmark character varying(255),
    city character varying(100) NOT NULL,
    state character varying(100),
    postal_code character varying(20),
    country_code character varying(2) NOT NULL,
    latitude numeric(10,7),
    longitude numeric(10,7),
    is_primary boolean DEFAULT false NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_party_addresses_latitude CHECK (((latitude IS NULL) OR ((latitude >= ('-90'::integer)::numeric) AND (latitude <= (90)::numeric)))),
    CONSTRAINT chk_party_addresses_longitude CHECK (((longitude IS NULL) OR ((longitude >= ('-180'::integer)::numeric) AND (longitude <= (180)::numeric)))),
    CONSTRAINT chk_party_addresses_type CHECK (((address_type)::text = ANY ((ARRAY['BILLING'::character varying, 'SHIPPING'::character varying, 'OFFICE'::character varying, 'REGISTERED'::character varying, 'WAREHOUSE'::character varying, 'HOME'::character varying, 'OTHER'::character varying])::text[])))
);


--
-- Name: party_contact_persons; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.party_contact_persons (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    party_id uuid NOT NULL,
    first_name character varying(100) NOT NULL,
    middle_name character varying(100),
    last_name character varying(100),
    designation character varying(150),
    department character varying(150),
    email character varying(255),
    phone character varying(30),
    mobile character varying(30),
    is_primary boolean DEFAULT false NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    notes text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: party_contacts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.party_contacts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    party_id uuid NOT NULL,
    contact_type character varying(30) NOT NULL,
    contact_value character varying(255) NOT NULL,
    label character varying(100),
    is_primary boolean DEFAULT false NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_party_contacts_type CHECK (((contact_type)::text = ANY ((ARRAY['EMAIL'::character varying, 'PHONE'::character varying, 'MOBILE'::character varying, 'WHATSAPP'::character varying, 'FAX'::character varying, 'OTHER'::character varying])::text[])))
);


--
-- Name: party_notes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.party_notes (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    party_id uuid NOT NULL,
    note_type character varying(50),
    content text NOT NULL,
    created_by uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- Name: supplier_user_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.supplier_user_accounts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    supplier_id uuid NOT NULL,
    user_id uuid NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: suppliers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.suppliers (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    party_id uuid NOT NULL,
    supplier_code character varying(50) NOT NULL,
    supplier_group character varying(100),
    payment_terms_days integer DEFAULT 0 NOT NULL,
    payment_method character varying(50),
    currency_code character varying(3),
    purchase_rep_id uuid,
    tax_exempt boolean DEFAULT false NOT NULL,
    tax_exemption_number character varying(100),
    supplier_rating numeric(3,2),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_suppliers_payment_terms CHECK ((payment_terms_days >= 0)),
    CONSTRAINT chk_suppliers_rating CHECK (((supplier_rating IS NULL) OR ((supplier_rating >= (0)::numeric) AND (supplier_rating <= (5)::numeric))))
);


--
-- Name: customer_user_accounts customer_user_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer_user_accounts
    ADD CONSTRAINT customer_user_accounts_pkey PRIMARY KEY (id);


--
-- Name: customer_user_accounts customer_user_accounts_user_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer_user_accounts
    ADD CONSTRAINT customer_user_accounts_user_id_key UNIQUE (user_id);


--
-- Name: customers customers_customer_code_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_customer_code_key UNIQUE (customer_code);


--
-- Name: customers customers_party_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_party_id_key UNIQUE (party_id);


--
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);


--
-- Name: individuals individuals_party_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.individuals
    ADD CONSTRAINT individuals_party_id_key UNIQUE (party_id);


--
-- Name: individuals individuals_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.individuals
    ADD CONSTRAINT individuals_pkey PRIMARY KEY (id);


--
-- Name: organizations organizations_party_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.organizations
    ADD CONSTRAINT organizations_party_id_key UNIQUE (party_id);


--
-- Name: organizations organizations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.organizations
    ADD CONSTRAINT organizations_pkey PRIMARY KEY (id);


--
-- Name: parties parties_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.parties
    ADD CONSTRAINT parties_pkey PRIMARY KEY (id);


--
-- Name: party_addresses party_addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.party_addresses
    ADD CONSTRAINT party_addresses_pkey PRIMARY KEY (id);


--
-- Name: party_contact_persons party_contact_persons_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.party_contact_persons
    ADD CONSTRAINT party_contact_persons_pkey PRIMARY KEY (id);


--
-- Name: party_contacts party_contacts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.party_contacts
    ADD CONSTRAINT party_contacts_pkey PRIMARY KEY (id);


--
-- Name: party_notes party_notes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.party_notes
    ADD CONSTRAINT party_notes_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: supplier_user_accounts supplier_user_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_user_accounts
    ADD CONSTRAINT supplier_user_accounts_pkey PRIMARY KEY (id);


--
-- Name: supplier_user_accounts supplier_user_accounts_user_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_user_accounts
    ADD CONSTRAINT supplier_user_accounts_user_id_key UNIQUE (user_id);


--
-- Name: suppliers suppliers_party_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_party_id_key UNIQUE (party_id);


--
-- Name: suppliers suppliers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_pkey PRIMARY KEY (id);


--
-- Name: suppliers suppliers_supplier_code_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_supplier_code_key UNIQUE (supplier_code);


--
-- Name: idx_customer_user_accounts_customer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_customer_user_accounts_customer ON public.customer_user_accounts USING btree (customer_id);


--
-- Name: idx_customer_user_accounts_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_customer_user_accounts_user ON public.customer_user_accounts USING btree (user_id);


--
-- Name: idx_customers_credit_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_customers_credit_status ON public.customers USING btree (credit_status);


--
-- Name: idx_customers_customer_code; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_customers_customer_code ON public.customers USING btree (customer_code);


--
-- Name: idx_customers_group; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_customers_group ON public.customers USING btree (customer_group);


--
-- Name: idx_customers_party_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_customers_party_id ON public.customers USING btree (party_id);


--
-- Name: idx_customers_sales_rep; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_customers_sales_rep ON public.customers USING btree (sales_rep_id);


--
-- Name: idx_individuals_job_title; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_individuals_job_title ON public.individuals USING btree (job_title);


--
-- Name: idx_individuals_last_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_individuals_last_name ON public.individuals USING btree (last_name);


--
-- Name: idx_individuals_national_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_individuals_national_id ON public.individuals USING btree (national_id);


--
-- Name: idx_organizations_industry; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_organizations_industry ON public.organizations USING btree (industry);


--
-- Name: idx_organizations_registration_number; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_organizations_registration_number ON public.organizations USING btree (registration_number);


--
-- Name: idx_parties_display_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_parties_display_name ON public.parties USING btree (display_name);


--
-- Name: idx_parties_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_parties_email ON public.parties USING btree (email);


--
-- Name: idx_parties_is_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_parties_is_active ON public.parties USING btree (is_active);


--
-- Name: idx_parties_party_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_parties_party_type ON public.parties USING btree (party_type);


--
-- Name: idx_parties_phone; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_parties_phone ON public.parties USING btree (phone);


--
-- Name: idx_parties_tax_number; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_parties_tax_number ON public.parties USING btree (tax_number);


--
-- Name: idx_party_addresses_country; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_party_addresses_country ON public.party_addresses USING btree (country_code);


--
-- Name: idx_party_addresses_party_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_party_addresses_party_id ON public.party_addresses USING btree (party_id);


--
-- Name: idx_party_addresses_postal_code; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_party_addresses_postal_code ON public.party_addresses USING btree (postal_code);


--
-- Name: idx_party_addresses_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_party_addresses_type ON public.party_addresses USING btree (party_id, address_type);


--
-- Name: idx_party_contact_persons_department; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_party_contact_persons_department ON public.party_contact_persons USING btree (party_id, department);


--
-- Name: idx_party_contact_persons_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_party_contact_persons_email ON public.party_contact_persons USING btree (email);


--
-- Name: idx_party_contact_persons_party; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_party_contact_persons_party ON public.party_contact_persons USING btree (party_id);


--
-- Name: idx_party_contact_persons_phone; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_party_contact_persons_phone ON public.party_contact_persons USING btree (phone);


--
-- Name: idx_party_contacts_party_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_party_contacts_party_id ON public.party_contacts USING btree (party_id);


--
-- Name: idx_party_contacts_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_party_contacts_type ON public.party_contacts USING btree (party_id, contact_type);


--
-- Name: idx_party_contacts_value; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_party_contacts_value ON public.party_contacts USING btree (contact_value);


--
-- Name: idx_party_notes_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_party_notes_created_at ON public.party_notes USING btree (created_at);


--
-- Name: idx_party_notes_party; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_party_notes_party ON public.party_notes USING btree (party_id);


--
-- Name: idx_party_notes_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_party_notes_type ON public.party_notes USING btree (party_id, note_type);


--
-- Name: idx_supplier_user_accounts_supplier; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_supplier_user_accounts_supplier ON public.supplier_user_accounts USING btree (supplier_id);


--
-- Name: idx_supplier_user_accounts_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_supplier_user_accounts_user ON public.supplier_user_accounts USING btree (user_id);


--
-- Name: idx_suppliers_group; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_suppliers_group ON public.suppliers USING btree (supplier_group);


--
-- Name: idx_suppliers_party_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_suppliers_party_id ON public.suppliers USING btree (party_id);


--
-- Name: idx_suppliers_purchase_rep; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_suppliers_purchase_rep ON public.suppliers USING btree (purchase_rep_id);


--
-- Name: idx_suppliers_supplier_code; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_suppliers_supplier_code ON public.suppliers USING btree (supplier_code);


--
-- Name: uq_party_addresses_primary; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uq_party_addresses_primary ON public.party_addresses USING btree (party_id, address_type) WHERE (is_primary = true);


--
-- Name: uq_party_contact_persons_primary; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uq_party_contact_persons_primary ON public.party_contact_persons USING btree (party_id) WHERE (is_primary = true);


--
-- Name: customer_user_accounts fk_customer_user_accounts_customer; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customer_user_accounts
    ADD CONSTRAINT fk_customer_user_accounts_customer FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE CASCADE;


--
-- Name: customers fk_customers_party; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT fk_customers_party FOREIGN KEY (party_id) REFERENCES public.parties(id) ON DELETE CASCADE;


--
-- Name: organizations fk_organizations_party; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.organizations
    ADD CONSTRAINT fk_organizations_party FOREIGN KEY (party_id) REFERENCES public.parties(id) ON DELETE CASCADE;


--
-- Name: party_addresses fk_party_addresses_party; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.party_addresses
    ADD CONSTRAINT fk_party_addresses_party FOREIGN KEY (party_id) REFERENCES public.parties(id) ON DELETE CASCADE;


--
-- Name: party_contact_persons fk_party_contact_persons_party; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.party_contact_persons
    ADD CONSTRAINT fk_party_contact_persons_party FOREIGN KEY (party_id) REFERENCES public.parties(id) ON DELETE CASCADE;


--
-- Name: party_contacts fk_party_contacts_party; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.party_contacts
    ADD CONSTRAINT fk_party_contacts_party FOREIGN KEY (party_id) REFERENCES public.parties(id) ON DELETE CASCADE;


--
-- Name: party_notes fk_party_notes_party; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.party_notes
    ADD CONSTRAINT fk_party_notes_party FOREIGN KEY (party_id) REFERENCES public.parties(id) ON DELETE CASCADE;


--
-- Name: supplier_user_accounts fk_supplier_user_accounts_supplier; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier_user_accounts
    ADD CONSTRAINT fk_supplier_user_accounts_supplier FOREIGN KEY (supplier_id) REFERENCES public.suppliers(id) ON DELETE CASCADE;


--
-- Name: suppliers fk_suppliers_party; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT fk_suppliers_party FOREIGN KEY (party_id) REFERENCES public.parties(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict kIHiL32QxJnkOtZuEgeMMA3h4pvfN6d0Cw6SugYUFLTnPFk1RuQZgjtQdYAtMYv

