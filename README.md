# Project Jarvis

Project Jarvis is an attempt to bring consistency and order to Iron.io's APIs while simultaneously making them understandable by machines. The goal is to bring best practices to the APIs and increase the number of things that can be easily automated.

## Resources-First

Jarvis revolves around the RESTful idea of a &quot;resource&quot;. The core premise is that APIs exist to manipulate resources. Jarvis exists to let Iron.io define resources, their capabilities, and their constraints, then generate API endpoints, client libraries, documentation, and other useful tools out of those definitions. It moves resources to the centre of the API design process, a place that is currently held by representations.

The workflow is simple: resources are defined in special &quot;resource files&quot;, which enumerate the properties, actions, and constraints of the resource. The Jarvis CLI is then used to generate specific documentation, endpoints, and tooling based on the resources. Client libraries in each language exist to understand resource files, and are capable of generating client libraries for each API based on the resource files.

## Intended Tools

* `jarvis apidef`: generate a list of endpoints, the request body for that endpoint, and the response that endpoint will return.
* `jarvis docs`: generates HTML documenting the API (as defined by `apidef`) that corresponds to the resources.
* `jarvis testserver`: starts an HTTP server that will generate endpoints as in `apidef`, then keep track of requests to ensure full coverage of each client library. Errors can be coerced from the testserver using special headers.
* `jarvis clientspec`: generates a document that can be used to generate client libraries.
