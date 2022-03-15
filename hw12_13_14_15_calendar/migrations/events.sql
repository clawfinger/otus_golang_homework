-- Table: public.Events

-- DROP TABLE IF EXISTS public."Events";

CREATE TABLE IF NOT EXISTS public."Events"
(
    "ID" character varying COLLATE pg_catalog."default" NOT NULL,
    "Title" character varying COLLATE pg_catalog."default",
    "Date" date,
    "Duration" integer,
    "Description" character varying COLLATE pg_catalog."default",
    "OwnerID" character varying COLLATE pg_catalog."default",
    "NotifyTime" integer,
    CONSTRAINT "Events1_pkey" PRIMARY KEY ("ID")
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."Events"
    OWNER to admin;