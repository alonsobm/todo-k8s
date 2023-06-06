# Your todo app, now in kubernetes

## Local development

For local development, a docker postgres container  is needed. Makefile is provided for that. Commands are:
```
make pgdev
make run
```
Migrations are under migrations/ dir, and they can be run via migrate tool, so you will need to install golang-migrate cli.

Configuration files are under config/ dir. For security reason, there are some gitignored configfiles

Kubernetes yml are under /manifests dir. For security reason, the configmap that is mounted for the deployment is being gitignored

Application is under todoapp/ dir. 

Tests: There are integration/acceptance tests under todoapp/tests/ dir, you can run them via ```make test.integration```

CI/CD: github actions yml is under .github/workflows/ dir. It is triggered on push to master


## Decisions

I have decided to deploy the app using digital ocean cloud, since classic cloud providers like aws, azure, etc. have an excessive amount of documentation.
I thought that little cloud providers like DO would have less doc to read, so that will save time and some money, so I decided to go with DO and try its 
services.

I have used GitHub Actions as CI/CD system, since I read it was very simple, and decided to give it a try and learn about it. 

For kubernetes Ingress, I decided to assign the IP of the Digital Ocean load balancer to a new subdomain that I created for a previous domain I had. Since I didn't want
to pay for a new domain, I just created a new "A" DNS record for that domain, which is managed by Cloudfare. If you visit http://kubetest.alonsobm.com you will receive a 
message from the pods under the kubernetes Service (I created a Replica Set with 2 pods). I'm not using TLS for that domain, so connection is not encrypted. This is because 
setting up https in kubernetes would have involved reading more documentation and troubleshot kubernetes networking layers. 

Unit tests are not included due to the amount of time that would require test all the different parts of the app. Integration/Acceptance test are included under 
todoapp/tests dir. I am using a TestMain function that runs before the tests, and then performing the tests against the app with a postgres container behind.

I am not using DDD or Hexagonal architecture inside the app, because this is a simple app, so I decided to follow a KISS approach. Still I have decided to use some interfaces
in order to hide the storage (postgres in this case). If in the future it is necessary to use let's say, Cassandra or Mongo, I can implement a different Service interface and
inject it to the HttpHandlers (or future CLI / GRPC  "handlers") I have.

## Improvements

Due to time, I think there are several improvements that can be done to the app, like JSON validation, which I am not doing,use https and create ssl certificates via let's
encrypt, improve the error handling via creating different errors. Other missing feature is shutdown process, which for an environment like kubernetes, where pods are going up and down,
is important to ensure the finish pending requests, We can create a channel to receive OS signals, spin up the http server in other goroutine and wait to receive a signal from that 
channel, and then shutdown the httpServer. In some cases, I am passing up the error so the client can read some SQL errors, that should not. High availability in the server is 
achieved since I have configured the Kubernetes cluster with two nodes, but not in db, that can be solved configuring the managed Postgres service that Digital Ocean
have. A swagger can also be nice to have in order to interact with the server, but again




