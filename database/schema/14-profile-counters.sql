/**
 * Table that contains the counters related to a profile address.
 * This is done in order to improve performance avoiding using COUNT queries.
 */
CREATE TABLE profile_counters
(
    row_id                  SERIAL NOT NULL PRIMARY KEY,

    profile_address         TEXT   NOT NULL,
    relationships_count     BIGINT NOT NULL DEFAULT 0,
    blocks_count            BIGINT NOT NULL DEFAULT 0,
    chain_links_count       BIGINT NOT NULL DEFAULT 0,
    application_links_count BIGINT NOT NULL DEFAULT 0,
    CONSTRAINT unique_profile_counters UNIQUE (profile_address)
);