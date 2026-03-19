# Casdoor Integration

Pagode supports [Casdoor](https://casdoor.org) as an opt-in external authentication provider. When enabled, login and registration are handled by Casdoor via OAuth 2.0, while Pagode manages local sessions and user data.

## Quick Start

### 1. Start Casdoor

```bash
docker compose up -d casdoor
```

Casdoor will be available at http://localhost:8100.
Default credentials: `admin` / `123`

### 2. Configure a Pagode Application in Casdoor

1. Open http://localhost:8100 and log in
2. Go to **Applications** > **Add**
3. Set **Name**: `pagode`
4. Set **Organization**: `built-in`
5. Set **Redirect URL**: `http://localhost:8000/auth/casdoor/callback`
6. Copy the **Client ID** and **Client Secret**
7. Save

### 3. Get the Certificate

1. Go to **Certs** in the Casdoor admin panel
2. Find the `cert-built-in` certificate
3. Copy the certificate content (the public key)

### 4. Configure Pagode

Set the following environment variables (or update `config/config.yaml`):

```bash
export PAGODE_AUTH_PROVIDER=casdoor
export PAGODE_AUTH_CASDOOR_ENDPOINT=http://localhost:8100
export PAGODE_AUTH_CASDOOR_CLIENTID=<your-client-id>
export PAGODE_AUTH_CASDOOR_CLIENTSECRET=<your-client-secret>
export PAGODE_AUTH_CASDOOR_CERTIFICATE=<certificate-content>
export PAGODE_AUTH_CASDOOR_ORGANIZATIONNAME=built-in
export PAGODE_AUTH_CASDOOR_APPLICATIONNAME=pagode
```

### 5. Run Pagode

```bash
make run
```

Login and registration will now redirect to Casdoor.

## Production

In production, run Casdoor as a separate service (Docker, Kubernetes, or bare metal). See the [Casdoor deployment docs](https://casdoor.org/docs/basic/server-installation) for details.

Update the environment variables to point to your production Casdoor instance.

## How It Works

- When `auth.provider=casdoor`, the login/register pages redirect to Casdoor
- After authentication, Casdoor redirects back with an authorization code
- Pagode exchanges the code for a JWT token and extracts the user's email/name
- A local user is created (or linked by email) and a session is established
- All downstream features (dashboard, payments, chat, admin) work unchanged
