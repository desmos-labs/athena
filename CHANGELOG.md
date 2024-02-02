## v2.0.0

With this version of Athena, we are introducing a new PostgreSQL table named `profile_counters` that will be used to
store the number of relationships, blocks, chain links and application links associated to a user. This is done to
improve the performance of queries that need to count the number of these objects without having to scan the entire
database.

Please note that this change requires a new database table to be created. To do so, you can use the following SQL
statement:

```sql
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
```

If you also want to track such table inside Hasura, you can run the following commands to import the updated Hasura
metadata into your endpoint:

```shell
$ cd hasura
$ hasura metadata apply --endpoint <your-hasura-endpoint> --admin-secret <your-admin-secret>
```

## Athena v1.0.0

We are excited to announce the release of Athena v1.0.0, marking a significant milestone in the evolution of our
project. Formerly known as DJuno, we have rebranded our Golang package to Athena for branding reasons.

### Why the Change?

The decision to rebrand from DJuno to Athena is rooted in our commitment to clarity, consistency, and a more cohesive
identity. Athena, named after the Greek goddess of wisdom, symbolizes the intelligence and strength that our project
embodies. This change aligns with our vision for the future and reflects the evolution of our goals and values.

### What's New?

- **Name Change:** The Golang package has transitioned from `github.com/desmos-labs/djuno`
  to `github.com/desmos-labs/athena`. Update your dependencies accordingly.

- **Brand Refresh:** Alongside the name change, we've refreshed the branding to better represent the essence of our
  project.

### Important Note:

With the release of Athena v1.0.0, DJuno is officially discontinued and will no longer receive support or updates. All
future developments, enhancements, and bug fixes will be focused on Athena. We encourage all users to migrate to Athena
to take advantage of the latest features and improvements.

### How to Migrate?

Updating your project to use Athena is a straightforward process. Simply replace all references
to `github.com/desmos-labs/djuno/v2` with `github.com/desmos-labs/athena` in your Go module files.

Thank you for your continued support, and we look forward to building a smarter, more robust future with Athena!

For any questions or assistance during the migration, feel free to reach out to us on
our [official support channels](https://desmos.discord.network).

Happy coding!

[Desmos Labs](https://github.com/desmos-labs)