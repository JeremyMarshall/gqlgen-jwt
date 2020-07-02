# GQLGEN-JWT

[![Build Status](https://travis-ci.org/JeremyMarshall/gqlgen-jwt.svg?branch=master)](https://travis-ci.org/JeremyMarshall/gqlgen-jwt)
[![codecov](https://codecov.io/gh/JeremyMarshall/gqlgen-jwt/branch/master/graph/badge.svg)](https://codecov.io/gh/JeremyMarshall/gqlgen-jwt)
[![Go Report Card](https://goreportcard.com/badge/github.com/JeremyMarshall/gqlgen-jwt)](https://goreportcard.com/report/github.com/JeremyMarshall/gqlgen-jwt)


An example of using JWT tokens and RBAC to protect GQL endpoints

## Schema

There are two parts to this, handled in middleware

### JWT

JWT is processed in HTTP middleware and takes a token in the header and converts it into a set of roles.

There is a test endpoint which will generate a JWT token. There is no authentication on this endpoint so it is not suitable for production systems.

```gql
mutation {
  createJwt(input: { user: "user", roles: ["jwt", "rbac-rw"] })
}
```

You can interrogate the roles in a token with

```gql
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

```http
{
  "Authorization": "Bearer <token>"
}
```

as this end point is protected with the `jwt` role

## RBAC

Rbac middleware is gqlgen middleware and it will validate the decoded token roles to the required role for the end point

The rbac endpoints allows for querying and update of the rbac.

This too needs an auth token with roles `rbac-ro` for query and `rbac-rw` for mutate


## Payload

The payload is some endpints which are protected by RBAC. There are two types

1. RBAC only
2. RBAC with a domain

### RBAC Only

This works in the same way as the RBAC above and allows users with the correct role (which has the correct permission) to access the endpoint.

### RBAC with domain

This is as above but will also check a defined field in the args for access

## Schema

[schema.graphqls][1]

## Running

Download the repo and run `make server`
hit port [http://localhost:8088][2]

You can also download the docker image

`docker run -p 8088:8088 jeremymarshall/gqlgen-jwt:latest`

Or run in Kubernetes

`make deploy`


## GQLGEN middleware

The middleware is defined by directives in the schema

```gql
enum RBAC {
    JWT_QUERY
    JWT_MUTATE

    RBAC_QUERY
    RBAC_MUTATE

    MOD_NEWSPAPER
    MOD_STAFF
    MOD_STORY
    MOD_PHOTO
    DEL_MEDIA
}

enum DOMAIN {
  newspaper
}

directive @HasRbac(rbac: RBAC!) on ARGUMENT_DEFINITION | INPUT_FIELD_DEFINITION | FIELD_DEFINITION
directive @HasRbacDomain(rbac: RBAC!, domainField: DOMAIN! ) on ARGUMENT_DEFINITION | INPUT_FIELD_DEFINITION | FIELD_DEFINITION
```

This is then implemented (here in `main.go`) 

```go
type RbacMiddlewareFunc func(ctx context.Context, obj interface{}, next graphql.Resolver, rbac model.Rbac) (res interface{}, err error)
type RbacDomainMiddlewareFunc func(ctx context.Context, obj interface{}, next graphql.Resolver, rbac model.Rbac, domainFiled model.Domain) (res interface{}, err error)

func RbacMiddleware(rbacChecker types.Rbac) RbacMiddlewareFunc {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, rbac model.Rbac) (res interface{}, err error) {
		if !rbacChecker.Check(GetCurrentUser(ctx).Roles, rbac.String()) {
			// block calling the next resolver
			return nil, fmt.Errorf("Access denied")
		}

		// or let it pass through
		return next(ctx)
	}
}

func RbacDomainMiddleware(rbacChecker types.Rbac) RbacDomainMiddlewareFunc {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, rbac model.Rbac, domainString model.Domain) (res interface{}, err error) {

		if args, ok := obj.(map[string]interface{}); ok {
			if domain, ok := args[domainString.String()].(string); ok {
				if rbacChecker.CheckDomain(GetCurrentUser(ctx).Roles, &domain, rbac.String()) {
					return next(ctx)
				}
			}
		}
		return nil, fmt.Errorf("Access denied")
	}
}
```

Then tied together in the config for GQLGEN
```go
	c := generated.Config{
		Resolvers: resolver,
		Directives: generated.DirectiveRoot{
			HasRbac:       RbacMiddleware(rbac),
			HasRbacDomain: RbacDomainMiddleware(rbac),
		},
	}
```

This middleware is called before the main schema functions and can be used to validate the request

[1]: ./graph/schema.graphqls
[2]: http://localhost:8088

