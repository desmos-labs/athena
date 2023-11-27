# Hasura
Once Athena is up and running, what you can do is use [Hasura](https://hasura.io/) to run a GraphQL sever to expose the data you want.

## Running Hasura
The easiest way to run a Hasura server is to use [Docker](https://www.docker.com/). This can be done by following [this guide](https://hasura.io/docs/2.0/graphql/core/getting-started/docker-simple.html).

### Variables
Make sure you edit the following variables accordingly: 

- `HASURA_GRAPHQL_DATABASE_URL`  
   URL used to connect to your PostgreSQL database.
- `HASURA_GRAPHQL_ENABLE_CONSOLE`  
   Tells whether you want to have the admin console enabled or not.
- `HASURA_GRAPHQL_ACCESS_KEY`  
   Represents the secret key that is going to be used to access the admin console.
- `HASURA_GRAPHQL_UNAUTHORIZED_ROLE`  
   Specifies the name of the role that should be assigned to users that are not performing authorized requests.
  
## Track the data
Once everything is setup properly, you can now access the Hasura console located at [http://localhost:8080](http://localhost:8080). From here, what you need to do is [track the data](https://hasura.io/docs/latest/graphql/core/databases/postgres/schema/using-existing-database.html).

## Start querying 
Once the data is tracked, you are now ready to [write queries](https://graphql.org/learn/queries/).