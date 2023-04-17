# Jury

A project designed to create a new pairwise judging system using modern technologies aimed at optimizing the user experience of judges and admin users. See the inspiration for this project: Gavel by anishathalye. Refer to this excellent article for more details on the underlying formulas of Gavel! The majority of our algorithm will be based off of research done in the field of pairwise comparison.

# Getting Started
These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

# Prerequisites
1)Docker

2)yarn

3)cargo

4)GNU Scientific Library (GSL)

# Installing
1)Clone the repository and navigate to the project directory.

2)Copy the ```.env.template file``` to ```.env``` and fill in the environmental variables. You will also need a MONGODB_URI field. 

3)Additionally, ```copy client/.env.template``` into ```client/.env``` and copy over the relevant environmental variables from ```.env```.

4)Run the following command to start the development environment with Dock.

```
docker-compose -f docker-compose.dev.yml up
```
5)Access the web application at localhost:3000.

## Alternatively, you can manually install the application by following these steps:

1)Clone the repository and navigate to the project directory.

2)Copy the ```.env.template file``` to ```.env``` and fill in the environmental variables. You will also need a ```MONGODB_URI``` field. Additionally, copy ```client/.env.template``` into ```client/.env``` and copy over the relevant environmental variables from ```.env```.

3)To start the client development server, navigate to the client directory and run the following commands:
```
yarn install
yarn start
```
4)To start the backend development server, navigate to the project directory and run the following command:
```
cargo run
```
5)Access the web application at ```localhost:3000```.

# Design
The design of the project can be found on Figma at this link: 

```https://www.figma.com/file/qwBWs4i7pJMpFbcjMffDZU/Jury-(Gavel-Plus)?node-id=8%3A100&t=egLV7iVmwvNRXef5-1```


