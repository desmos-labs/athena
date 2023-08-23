# Upgrade Process

This document provides guidance on upgrading from one DJuno version to another.

## Upgrading from Desmos Core v5 to Desmos Core v6

In Desmos Core v6, a new `owner` field has been introduced to the `Post` type. Consequently, it's necessary to modify the `post` table in the database by incorporating the new column:

```sql
ALTER TABLE post ADD COLUMN owner_address TEXT;
```