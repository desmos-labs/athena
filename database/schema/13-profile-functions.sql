/**
 * Function that allows to check if the current Hasura user is following a given profile.
 */
CREATE OR REPLACE FUNCTION is_user_following_profile(profile_row profile, hasura_session json)
    RETURNS BOOLEAN AS
$$
SELECT EXISTS (SELECT 1
               FROM user_relationship
               WHERE user_relationship.subspace_id =
                     CAST(COALESCE(hasura_session ->> 'x-hasura-selected-subspace-id', '0') AS BIGINT)
                 AND user_relationship.counterparty_address = profile_row.address
                 AND user_relationship.creator_address = hasura_session ->> 'x-hasura-user-address')
$$ LANGUAGE sql STABLE;


/**
 * Function that allows to check if the current Hasura user has been blocked by the given profile.
 */
CREATE OR REPLACE FUNCTION is_user_blocked_by_profile(profile_row profile, hasura_session json)
    RETURNS BOOLEAN AS
$$
SELECT EXISTS (SELECT 1
               FROM user_block
               WHERE user_block.subspace_id =
                     CAST(COALESCE(hasura_session ->> 'x-hasura-selected-subspace-id', '0') AS BIGINT)
                 AND user_block.blocker_address = profile_row.address
                 AND user_block.blocked_address = hasura_session ->> 'x-hasura-user-address')
$$ LANGUAGE sql STABLE;

/**
 * Function that allows to check if the current Hasura user has blocked the given profile.
 */
CREATE OR REPLACE FUNCTION has_user_blocked_profile(profile_row profile, hasura_session json)
    RETURNS BOOLEAN AS
$$
SELECT EXISTS (SELECT 1
               FROM user_block
               WHERE user_block.subspace_id =
                     CAST(COALESCE(hasura_session ->> 'x-hasura-selected-subspace-id', '0') AS BIGINT)
                 AND user_block.blocked_address = profile_row.address
                 AND user_block.blocker_address = hasura_session ->> 'x-hasura-user-address')
$$ LANGUAGE sql STABLE;