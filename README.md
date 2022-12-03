# Gamescore API
This repository provides a backend service that keeps track of high scores records. <br />


##Notes
* The target system should have docker/docker-compose installed.
* The following ports `8080` and `3000` should be reserved or changed in the attached `docker-compose.yml` file.
* To prevent pulling the frontend image, a local version can be build by following the [front-end](https://github.com/LordRahl90/score-frontend) repository.

# Setup
* Clone the repository
* RUN `make docker-up` (Please not that this will also attempt to pull the front-end repository)
* If you desire to run this without the front-end, just run `docker-compose up backend mysql` and you can test with postman.