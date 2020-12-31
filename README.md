# OPDREWSKI ENGINE

## Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Changelog](#changelog)
- [Reference](#reference)

## About <a name = "about"></a>

The purposes of this project is learning how to creating apps with clean code and S.O.L.I.D principle, using DDD ( Domain Driven Design) approach. This repository still active maintaining during until some objectives was fulfilled. Using Go as main language and some 3rd party package make it this apps to be lighten (It's my opinion :-p )

Always check at Changelog for history of this repository. cheers B-)

## Getting Started <a name = "getting_started"></a>

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See [deployment](#deployment) for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them.

1. Go Binary >= 1.15.x
2. MongoDB >= 4.x

### Running 

Just run command "go run <root_dir_project>/app/main.go"

## Changelog <a name = "changelog"></a>

### 1.0.0
* Initiate project structure with DDD approach, guidance from Uncle Bob advice
* Embedding MongoDB Driver for persistence data
* Using echo web framework for REST API delivery layer
* Create simple create and get function

## Reference
Many thanks for article from this : 
1. https://jobel.dev/unit-testing-using-mocks-for-golangs-api-with-domain-driven-design-clean-architecture/
2. https://medium.com/hackernoon/trying-clean-architecture-on-golang-2-44d615bf8fdf

And some books :
1.  “Clean Architecture: A Craftsman’s Guide to Software Structure and Design” famous author Robert “Uncle Bob” Martin presents an architecture with some important points like testability and independence of frameworks, databases and interfaces.
