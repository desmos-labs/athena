# Upgrade Process

This document provides step-by-step guidance for upgrading from DJuno version to another.

## Upgrading from Desmos Core v5 to Desmos Core v6

### Database

In Desmos Core v6, we've introduced a new feature called `AdditionalFeeTokens` to the `Subspace` type. To make use of
this feature, follow these steps to modify the database:

1. Open your SQL management tool.
2. Run the following SQL query:
   ```sql
   ALTER TABLE subspace
   ADD COLUMN additional_fee_tokens COIN[];
   ```

Additionally, Desmos Core v6 has added an `Owner` field to the `Post` type. To ensure compatibility, follow these steps:

1. Open your SQL management tool.
2. Run the following SQL query:
   ```sql
   ALTER TABLE post
   ADD COLUMN owner_address TEXT;
   ```

### GraphQL

Because of these new column additions, you need to update the Hasura metadata to include the new columns.