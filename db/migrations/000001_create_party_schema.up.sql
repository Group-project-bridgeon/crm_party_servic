CREATE EXTENSION IF NOT EXISTS pgcrypto;


-- PARTIES


CREATE TABLE parties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    party_type VARCHAR(30) NOT NULL,
    display_name VARCHAR(200) NOT NULL,

    legal_name VARCHAR(255),

    email VARCHAR(255),
    phone VARCHAR(30),

    website VARCHAR(255),

    tax_number VARCHAR(100),

    notes TEXT,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_parties_party_type
        CHECK (party_type IN ('ORGANIZATION', 'INDIVIDUAL'))
);

CREATE INDEX idx_parties_party_type
    ON parties(party_type);

CREATE INDEX idx_parties_display_name
    ON parties(display_name);

CREATE INDEX idx_parties_tax_number
    ON parties(tax_number);

CREATE INDEX idx_parties_email
    ON parties(email);

CREATE INDEX idx_parties_phone
    ON parties(phone);

CREATE INDEX idx_parties_is_active
    ON parties(is_active);


-- ORGANIZATIONS


CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    party_id UUID NOT NULL UNIQUE,

    registration_number VARCHAR(100),

    industry VARCHAR(150),

    organization_type VARCHAR(50),

    founded_date DATE,

    employee_count INT,

    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_organizations_party
        FOREIGN KEY (party_id)
        REFERENCES parties(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_organizations_employee_count
        CHECK (employee_count IS NULL OR employee_count >= 0)
);

CREATE INDEX idx_organizations_registration_number
    ON organizations(registration_number);

CREATE INDEX idx_organizations_industry
    ON organizations(industry);


-- INDIVIDUALS

CREATE TABLE individuals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    party_id UUID NOT NULL UNIQUE,

    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100),

    date_of_birth DATE,

    gender VARCHAR(30),

    job_title VARCHAR(150),

    national_id VARCHAR(100),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_individuals_last_name
    ON individuals(last_name);

CREATE INDEX idx_individuals_national_id
    ON individuals(national_id);

CREATE INDEX idx_individuals_job_title
    ON individuals(job_title);



-- PARTY CONTACTS


CREATE TABLE party_contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    party_id UUID NOT NULL,

    contact_type VARCHAR(30) NOT NULL,

    contact_value VARCHAR(255) NOT NULL,

    label VARCHAR(100),

    is_primary BOOLEAN NOT NULL DEFAULT FALSE,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_party_contacts_party
        FOREIGN KEY (party_id)
        REFERENCES parties(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_party_contacts_type
        CHECK (
            contact_type IN (
                'EMAIL',
                'PHONE',
                'MOBILE',
                'WHATSAPP',
                'FAX',
                'OTHER'
            )
        )
);

CREATE INDEX idx_party_contacts_party_id
    ON party_contacts(party_id);

CREATE INDEX idx_party_contacts_type
    ON party_contacts(party_id, contact_type);

CREATE INDEX idx_party_contacts_value
    ON party_contacts(contact_value);



-- PARTY ADDRESSES


CREATE TABLE party_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    party_id UUID NOT NULL,

    address_type VARCHAR(30) NOT NULL,

    address_name VARCHAR(100),

    address_line1 VARCHAR(255) NOT NULL,
    address_line2 VARCHAR(255),

    landmark VARCHAR(255),

    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    postal_code VARCHAR(20),

    country_code VARCHAR(2) NOT NULL,

    latitude DECIMAL(10,7),
    longitude DECIMAL(10,7),

    is_primary BOOLEAN NOT NULL DEFAULT FALSE,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_party_addresses_party
        FOREIGN KEY (party_id)
        REFERENCES parties(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_party_addresses_type
        CHECK (
            address_type IN (
                'BILLING',
                'SHIPPING',
                'OFFICE',
                'REGISTERED',
                'WAREHOUSE',
                'HOME',
                'OTHER'
            )
        ),

    CONSTRAINT chk_party_addresses_latitude
        CHECK (
            latitude IS NULL
            OR latitude BETWEEN -90 AND 90
        ),

    CONSTRAINT chk_party_addresses_longitude
        CHECK (
            longitude IS NULL
            OR longitude BETWEEN -180 AND 180
        )
);

CREATE INDEX idx_party_addresses_party_id
    ON party_addresses(party_id);

CREATE INDEX idx_party_addresses_type
    ON party_addresses(party_id, address_type);

CREATE INDEX idx_party_addresses_country
    ON party_addresses(country_code);

CREATE INDEX idx_party_addresses_postal_code
    ON party_addresses(postal_code);

-- Only one primary address per party + address type.
CREATE UNIQUE INDEX uq_party_addresses_primary
    ON party_addresses(party_id, address_type)
    WHERE is_primary = TRUE;



-- PARTY CONTACT PERSONS


CREATE TABLE party_contact_persons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    party_id UUID NOT NULL,

    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100),

    designation VARCHAR(150),

    department VARCHAR(150),

    email VARCHAR(255),
    phone VARCHAR(30),
    mobile VARCHAR(30),

    is_primary BOOLEAN NOT NULL DEFAULT FALSE,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    notes TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_party_contact_persons_party
        FOREIGN KEY (party_id)
        REFERENCES parties(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_party_contact_persons_party
    ON party_contact_persons(party_id);

CREATE INDEX idx_party_contact_persons_email
    ON party_contact_persons(email);

CREATE INDEX idx_party_contact_persons_phone
    ON party_contact_persons(phone);

CREATE INDEX idx_party_contact_persons_department
    ON party_contact_persons(party_id, department);

CREATE UNIQUE INDEX uq_party_contact_persons_primary
    ON party_contact_persons(party_id)
    WHERE is_primary = TRUE;



-- CUSTOMERS


CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    party_id UUID NOT NULL UNIQUE,

    customer_code VARCHAR(50) NOT NULL UNIQUE,

    customer_group VARCHAR(100),

    price_list_id UUID,

    credit_limit DECIMAL(14,2) NOT NULL DEFAULT 0,

    payment_terms_days INT NOT NULL DEFAULT 0,

    credit_status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',

    sales_rep_id UUID,

    tax_exempt BOOLEAN NOT NULL DEFAULT FALSE,

    tax_exemption_number VARCHAR(100),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_customers_party
        FOREIGN KEY (party_id)
        REFERENCES parties(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_customers_credit_limit
        CHECK (credit_limit >= 0),

    CONSTRAINT chk_customers_payment_terms
        CHECK (payment_terms_days >= 0),

    CONSTRAINT chk_customers_credit_status
        CHECK (
            credit_status IN (
                'ACTIVE',
                'BLOCKED',
                'HOLD'
            )
        )
);

CREATE INDEX idx_customers_party_id
    ON customers(party_id);

CREATE INDEX idx_customers_customer_code
    ON customers(customer_code);

CREATE INDEX idx_customers_group
    ON customers(customer_group);

CREATE INDEX idx_customers_credit_status
    ON customers(credit_status);

CREATE INDEX idx_customers_sales_rep
    ON customers(sales_rep_id);



-- SUPPLIERS


CREATE TABLE suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    party_id UUID NOT NULL UNIQUE,

    supplier_code VARCHAR(50) NOT NULL UNIQUE,

    supplier_group VARCHAR(100),

    payment_terms_days INT NOT NULL DEFAULT 0,

    payment_method VARCHAR(50),

    currency_code VARCHAR(3),

    purchase_rep_id UUID,

    tax_exempt BOOLEAN NOT NULL DEFAULT FALSE,

    tax_exemption_number VARCHAR(100),

    supplier_rating DECIMAL(3,2),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_suppliers_party
        FOREIGN KEY (party_id)
        REFERENCES parties(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_suppliers_payment_terms
        CHECK (payment_terms_days >= 0),

    CONSTRAINT chk_suppliers_rating
        CHECK (
            supplier_rating IS NULL
            OR supplier_rating BETWEEN 0 AND 5
        )
);

CREATE INDEX idx_suppliers_party_id
    ON suppliers(party_id);

CREATE INDEX idx_suppliers_supplier_code
    ON suppliers(supplier_code);

CREATE INDEX idx_suppliers_group
    ON suppliers(supplier_group);

CREATE INDEX idx_suppliers_purchase_rep
    ON suppliers(purchase_rep_id);



-- CUSTOMER USER ACCOUNTS


CREATE TABLE customer_user_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    customer_id UUID NOT NULL,

    user_id UUID NOT NULL UNIQUE,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_customer_user_accounts_customer
        FOREIGN KEY (customer_id)
        REFERENCES customers(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_customer_user_accounts_customer
    ON customer_user_accounts(customer_id);

CREATE INDEX idx_customer_user_accounts_user
    ON customer_user_accounts(user_id);



-- SUPPLIER USER ACCOUNTS


CREATE TABLE supplier_user_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    supplier_id UUID NOT NULL,

    user_id UUID NOT NULL UNIQUE,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_supplier_user_accounts_supplier
        FOREIGN KEY (supplier_id)
        REFERENCES suppliers(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_supplier_user_accounts_supplier
    ON supplier_user_accounts(supplier_id);

CREATE INDEX idx_supplier_user_accounts_user
    ON supplier_user_accounts(user_id);



-- PARTY NOTES


CREATE TABLE party_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    party_id UUID NOT NULL,

    note_type VARCHAR(50),

    content TEXT NOT NULL,

    created_by UUID,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_party_notes_party
        FOREIGN KEY (party_id)
        REFERENCES parties(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_party_notes_party
    ON party_notes(party_id);

CREATE INDEX idx_party_notes_type
    ON party_notes(party_id, note_type);

CREATE INDEX idx_party_notes_created_at
    ON party_notes(created_at);