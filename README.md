<p align="center">
   <img src="https://github.com/dethdkn/ldap-nel/blob/main/public/nel.png?raw=true" alt="Ldap-Nel" width="100px"/>
<h1 align="center">Ldap-Nel</h1>
</p>
<p align="center">🔐 Modern web interface for LDAP administration</p>
<p align="center">
   <a href="https://github.com/dethdkn/ldap-nel/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/dethdkn/ldap-nel?color=%233da639&logo=open%20source%20initiative" alt="License"/>
  </a>
   <a href="https://gitmoji.dev">
      <img src="https://img.shields.io/badge/gitmoji-%20😜%20😍-FFDD67" alt="Gitmoji"/>
   </a>
   <a href="https://rosa.dev.br">
      <img src="https://img.shields.io/badge/check me!-👻-F28AA9" alt="rosa.dev.br"/>
   </a>
</p>

## ✨ Reason

At [CBPF](https://cbpf.br), where I work, we use LDAP Admin to make small adjustments to our OpenLDAP server. However, LDAP Admin is a Windows-only tool, which creates issues for team members who use macOS or Linux. To solve this, I developed Ldap-Nel, a web-based application that can be accessed from any platform, anywhere—making LDAP management more accessible and convenient for the whole team.\
I also wanted to learn Go 😄

## 🚀 Setup

1. **Install Docker and Docker Compose:**\
Follow the official [Docker installation](https://docs.docker.com/get-started/get-docker/) guide for your operating system.

2. **Create a `docker-compose.yaml` file:**\
Copy the contents from the [example](https://github.com/dethdkn/ldap-nel/blob/main/docker-compose.yaml) into your own `docker-compose.yaml`.

3. **Create a `.env` file:**\
Use the [`.env.example`](https://github.com/dethdkn/ldap-nel/blob/main/.env.example) as a starting point and adjust the values to match your environment.\
🔑 JWT_SECRET: At least 32 characters\
🔐 ENCRYPTION_KEY: This must be a base64-encoded string that decodes to 32 bytes (AES-256).

✅ To generate a valid key:
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

## ⬆️ Upgrade

To upgrade `ldap-nel` to the latest version:

1. **Pull the latest image:**
```bash
docker compose pull
```

2. **Recreate and restart the containers:**
```bash
docker compose up -d --force-recreate
```

## 🎯 Todo

### Search
- [x] Add support for **JPEG image**
- [x] Enable **hover to preview JPEG image**
- [ ] Add **Search DN** option in the UI
- [x] Add **refresh** option

### Add
- [ ] Add option to create a **new DN**
- [ ] Add option to create a **new Attribute**

### Update
- [x] Implement **attribute update** functionality
- [ ] Support common **LDAP password encryption** options when updating passwords
- [ ] Add **copy** option
- [ ] Add **move** option

### Delete
- [x] Implement **attribute deletion**
- [ ] Add option to delete a **DN**
- [ ] Add **smart delete** (context-aware cleanup)

### Export
- [ ] Add **export** option

### Import
- [ ] Add **import** option

### Docs
- [ ] Add **images** to README
- [ ] Create **VitePress** documentation

## 📝 License

Copyright © 2025 [Gabriel 'DethDKN' Rosa](https://github.com/dethdkn)\
This project is under [MIT license](https://github.com/dethdkn/ldap-nel/blob/main/LICENSE)
