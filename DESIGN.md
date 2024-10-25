## GOSIGN

#### Clone the repository
`git clone https://github.com/GiacomoCortesi/gosign.git`

#### Build
```
cd gosign

go build -o gosign
```

#### Run
`
./gosign
`

#### Test
```
cd gosign

go test -v ./...
```

#### Documentation
You can use godoc to display package documentation

```
go get golang.org/x/tools/cmd/godoc

cd gosign

godoc -http :8888
```

Open the browser at http://localhost:8888 and scroll to "Third party" section to display gosign module documentation

### API Documentation
REST API documentation is provided as an openapi specification.

Run swagger UI:

```
cd gosign

docker run -p 80:8080 -e SWAGGER_JSON=/openapi.yaml -v $PWD/openapi.yaml:/openapi.yaml swaggerapi/swagger-ui
```

Open browser at http://localhost to display swagger UI documentation.

## Solution Design

The solution intentionally avoid the use of any external library / framework for the sake of simplicity and readability.

I only added testify mock package to simplify mocking repository and service interfaces.

Any go web framework (gin, gorilla, etc.) may be used in place of standard http package, but given the simplicity of the task I didn't see the need to add external dependencies (since features as middleware, and routing are not explicitly requested by the challenge).

The solution approach focuses on "Less is More", adding/changing only what is strictly needed to comply with the problem specification.

The only exception being the CORS middleware, added for the sake of development to allow performing swagger UI calls locally.

Go version has been updated to 1.22 to leverage enhancements in http package (such as wildcard pattern matching), more info: [go 1.22 http package](https://go.dev/blog/routing-enhancements).

Possible enhancements:
 - centralized API errors management
 - logging instrumentation
 - automated interface mocks generation through mockery

## Requirements
#### REQ - 1: The system will be used by many concurrent clients accessing the same resources.

To allow safe concurrent access of the clients to the same shared resources the inmemory repository access is controlled with mutexes. If we were choosing a persistence storage solution with a database, concurrency would have been (most likely, depending on the database and the requirements) handled internally by the database itself.

#### REQ - 2: The `signature_counter` has to be strictly monotonically increasing and ideally without any gaps.

To make the `signature_counter` strictly monotonically increasing and without any gaps, we use sync/atomic package that allows atomic thread-safe operations to increment `signature_counter

#### REQ - 3: The system currently only supports `RSA` and `ECDSA` as signature algorithms. Try to design the signing mechanism in a way that allows easy extension to other algorithms without changing the core domain logic.

To allow easy extension to other algorithms without changing the core domain logic we use a factory to instantiate the appropriate Signer in the crypto package, so that whenever we need to add a new signature algorithm it is sufficient to change the crypto package in order to:
 - include the new algorithm in the SignatureAlgorithm enum type
 - implement the Signer interface for the new signature algorithm.

#### REQ - 4: For now it is enough to store signature devices in memory. Efficiency is not a priority for this. In the future we might want to scale out. As you design your storage logic, keep in mind that we may later want to switch to a relational database.

Defining a repository interface with CRUD operations on the data allows to later switch do a different data storage solution in a simple and effective manner, leaving untouched the business logic.
Switching to a different backend storage is a matter of implementing the repository interface in the persistence package and injecting the new repository to the service.

## QA/Testing
For the sake of testing I put down some observation:
 - In general I like the golang testing table approach with subtests to handle the different test cases, here I'm mostly using reflect to check results correctness, other approaches may be used such as testify asserts.
 - meaningful testing of the Signer's interface implementations is not trivial, as we cannot follow the classic got/want approach but we need to make sure that the signature can be used to correctly verify the input data.
 - testing of the in memory repository is trivial as it is sufficient to manually fill the repo and verify that operations behave as they should. Testing a repository that involves a database would be more complex as we'd need to either use in memory database or some database driver mocking libray, such as [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
  - Testing the service layer is more complex since the service layer depend on the injected repository, being the repository an interface we can mock it, here I'm using testify mock package for simplicity. Service layer testing is limited to SignTransaction method as it is implements the most complex business logic. This is for demonstrational purposes only, real implementation would include complete set of tests.
 - Similarly testing the api handlers require mocking the injected service. Handlers test are just a basic setup for demonstrational purposes only, real implementation would include complete set of tests (and also would include additional verification on the response data)