# iitk-coin

## Introduction
IITk-coin is a vision of a pseudo-currency for use in the IITK Campus. 

---
## Endpoints
These are the following functional endpoints in the current status of project
- ```/``` -> A general endpoint that just returns a message
- ```/signup``` -> For the registration of new users to the portal
- ```/login``` -> For the login of existing users to the portal
- ```/secretpage``` -> A special endpoint that can be used only if you are a verified user.
- ```/refresh``` -> A special endpoint that refreshes the expiration time of your JWT token.

---
## Usage

To use iitcoin , first clone the repo-
```
git clone <repo name>
```

Build the executable 

```
go build -o out *go
```

Run the executable

```
./out
```
