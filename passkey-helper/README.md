# Evolution Passkey Helper

Browser extension that completes **passkey** pairing for WhatsApp Web with Evolution GO.

Since WhatsApp began requiring passkeys when linking new devices, a QR Code alone is not sufficient for some accounts: a WebAuthn ceremony (`navigator.credentials.get`) must be executed, and the browser only permits this when the code runs on the relying party origin — `whatsapp.com`. This extension runs **only** on `web.whatsapp.com` and acts solely as a secure origin host for the biometric/PIN ceremony. It **does not** interfere with WhatsApp Web login and **does not** collect or transmit data to third parties.

## How It Works

1. Scan the QR code as usual in Evolution Manager. If the account requires a passkey, the connection modal displays the button **"Open WhatsApp Web"**.
2. This button opens `https://web.whatsapp.com/#wapk=<payload>`, where `payload` is a base64url JSON of `{ "t": "<instance token>", "b": "<API base URL>" }`.
3. The content script reads this payload, polls the ceremony status from the API, executes the WebAuthn ceremony on `whatsapp.com`, and sends the signature assertion back to Evolution.
4. If pairing requires manual confirmation, the panel shows the code and a **"Confirm code"** button. In most cases (QR scanned beforehand), Evolution confirms automatically and you only need to return to the manager.

No `host_permission` is required: requests originate from the `web.whatsapp.com` origin and the Evolution backend allows CORS for that origin.

## Endpoints Used (Evolution GO)

- `GET  {base}/passkey-ceremony/{token}` — ceremony status and challenge
- `POST {base}/passkey-ceremony/{token}/response` — sends the WebAuthn assertion
- `POST {base}/passkey-ceremony/{token}/confirm` — confirms the code (when applicable)

Where `base` is your Evolution API URL and `token` is the instance token.

## Installation in Google Chrome

1. Download / unzip this folder (containing `manifest.json`).
2. Open `chrome://extensions`.
3. Enable **Developer mode** (top right switch).
4. Click **Load unpacked**.
5. Select this folder.
6. Return to Evolution Manager and click **"Open WhatsApp Web"**.

## Installation in Microsoft Edge

1. Download / unzip this folder (containing `manifest.json`).
2. Open `edge://extensions`.
3. Enable **Developer mode** (bottom left switch).
4. Click **Load unpacked**.
5. Select this folder.
6. Return to Evolution Manager and click **"Open WhatsApp Web"**.

## Tips

- Use the browser logged into the same account (Google/iCloud) where your WhatsApp passkey is saved, or have your phone nearby for biometric confirmation.
- The ceremony challenge has a short expiration window. If it expires, generate a new QR code in Evolution.

## Whitelabel

The extension is generic: no API URL or token is hardcoded in the codebase — everything is supplied dynamically via the URL payload. To customize the name/icon, edit `manifest.json` and panel details in `content.js`.
