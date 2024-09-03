SSO single sign on. 
it would be used for auth purposes for url shortener service. 

(I'm currently working on the deployment of the service.)

## Self-hosting instructions: 
1. git clone
2. install protoc on PC and put it to bin folder.
3. run install script from makefile. 
4. run generate script from makefile.
5. run tidy script from makefile.
6. apply migrations to the database.
   `migrations-up` script from makefile
7. to start app run `run` script from makefile.

### Testing
- `func-tests` script will run tests.

## Grpc:
It'll use declared proto file to generate the client and server code.
Grpc is used for communication between services. 
Grpc provides a lot of ready to use code for the transport layer.
It will allow to focus on the business logic of the service.

### Endpoints:
- Register user 
- Login user
- Check is user logged is admin