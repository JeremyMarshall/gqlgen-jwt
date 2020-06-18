# GQLGEN-JWT


An example of using JWT tokens and RBAC to protect GQL endpoints

## Schema

There are two parts to this, handled in middleware

### JWT

JWT is processed in HTTP middleware and takes a token in the header and converts it into a set of roles.

There is a test endpoint which will generate a JWT token. There is no authentication on this endpoint so it is not suitable for production systems.

```
mutation {
  createJwt(input: { user: "user", roles: ["jwt", "rbac-rw"] })
}
```

You can interrogate the roles in a token with

```
query {
  jwt(
    token: "<token>"
  ) {
    roles
    properties {
      name
      value
    }
  }
}
```

You will need to pass a token in the header, in the playground you can use `HTTP HEADERS`

```
{
  "Authorization": "Bearer <token>"
}
```

as this end point is protected with the `jwt` role

## RBAC

Rbac middleware is gqlgen middleware and it will validate the decoded token roles to the required role for the end point

The rbac endpoints allows for querying and update of the rbac.

This too needs an auth token with roles `rbac-ro` for query and `rbac-rw` for mutate

## Schema

[schema.graphqls][1]

## Running

Download the repo and run `make server`
hit port [http://localhost:8088][2]

You can also download the docker image

`docker run -p 8088:8088 jeremymarshall/gqlgen-jwt:latest`

Or run in Kubernetes

`make deploy`


[1]: ./graph/schema.graphqls
[2]: http://localhost:8088

