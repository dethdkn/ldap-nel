## 🚀 Setup

1. **Install Docker and Docker Compose:**\
Follow the official [Docker installation](https://docs.docker.com/get-started/get-docker/) guide for your operating system.

2. **Create a `docker-compose.yaml` file:**\
Copy the contents from the [example](https://github.com/dethdkn/ldap-nel/blob/main/docker-compose.yaml) into your own `docker-compose.yaml`.\
or download it
```sh
curl -L -o docker-compose.yaml https://raw.githubusercontent.com/dethdkn/ldap-nel/main/docker-compose.yaml
# or
wget -O docker-compose.yaml https://raw.githubusercontent.com/dethdkn/ldap-nel/main/docker-compose.yaml
```

3. **Create a `.env` file:**\
Use the [`.env.example`](https://github.com/dethdkn/ldap-nel/blob/main/.env.example) as a starting point and adjust the values to match your environment.\
```sh
curl -L -o .env https://raw.githubusercontent.com/dethdkn/ldap-nel/main/.env.example
# or
wget -O .env https://raw.githubusercontent.com/dethdkn/ldap-nel/main/.env.example
```
🔑 JWT_SECRET: At least 32 characters\
🔐 ENCRYPTION_KEY: This must be a base64-encoded string that decodes to 32 bytes (AES-256).

✅ To generate a valid encryption key run:
```bash
openssl rand -base64 32
```

4. **Create the database folder:**\
Run the following command to create the required database volume directory (as defined in `docker-compose.yaml`):  
```bash
mkdir ./database
```

5. **Pull and start the container:**
```bash
docker compose pull
docker compose up -d --force-recreate
```

6.	**(Optional) Set up a reverse proxy:**\
For production use, you can put the Docker app behind an Nginx reverse proxy to serve it on ports 80 or 443.
