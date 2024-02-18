## IGP Code Challenge

### Running code

You can run the code with docker compose like:

```sh
  docker compose up --build
```

It will spin up:

- PostgreSQL Database
- API Server that contains solution to the challenge
- Web Server that contains frontend to test the challenge

#### Once the server is up, the docs are available at:

```
localhost:3001/docs
```

#### Client to test the API will be served at:

```
localhost:80
```

### Code challenge task:

```
Potrebno je dizajnirati i implementirati sljedeće:

- Registraciju i prijavu korisnika
- Slanje e-maila dobrodošlice registriranom korisniku (ne mora se slati stvarni email, tj. može
  biti mock-ano)
- Mehanizam za slanje notifikacija prijavljenom korisniku
- Kroz rješenje je potrebno pokazati znanje Docker-a.

Implementaciju je potrebno odraditi koristeći Fiber framework (GO).

Osim implementacije, potrebno je razmisliti i o različitim načinima implementacije rješenja,
korištenju različitih tehnologija u rješenju, ideje koje bi bile dobre za imati, a koje nisam naveo i
slično.

Naravno, potrebna je samo jedna implementacija rješenja. Alternative, dodatne tehnologije te ideje
implementacije za pojedine dijelove se ne moraju implementirati. Dovoljno ih je dokumentirati
negdje u projektu, na bilo koji način.

Rok za rješavanje zadatka je tjedan dana, a rješenje je potrebno poslati u obliku git repozitorija
(Github, Bitbucket, Gitlab i sl.)
```
