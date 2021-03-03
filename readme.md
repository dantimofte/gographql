# Golang backend with graphql and mongodb

The project is based on the guide from [howtographql.com/graphql-go](https://www.howtographql.com/graphql-go/0-introduction/) 

Differences : 

    - It uses mongodb instead of mysql
    - Login is setting the jwt token as cookie , no extra work required in the frontend
    - Added logout mutation


## Start the graphql backend
```
go get
go run ./server.go

```

## queries and mutations

```
mutation {
  createUser(input: {username: "user1", password: "123"})
}


mutation {
  login(input: {username: "user1", password: "123"})
}

mutation {
  logout(input: {info:true})
}

mutation {
  createLink(input: {title: "real link!", address: "www.graphql.org"}){
    user{
      name
    }
  }
}

query {
  links {
    title
    address
    id
  }
}

```