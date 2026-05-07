# Secrets

Catalog of secrets the app expects in Google Secret Manager. Values are NEVER stored here — only names, purposes, and rotation procedure.

Both prod (`crew-predictions`) and staging (`crew-predictions-staging`) hold the same set of secret names. Values must differ — never reuse the same `session-secret` or `admin-key` across environments.

## Required secrets

| Name             | Purpose                                                                                                              | Consumed by                          | Rotation cadence |
| ---------------- | -------------------------------------------------------------------------------------------------------------------- | ------------------------------------ | ---------------- |
| `admin-key`      | HTTP header value required by the `AdminAuth` middleware to gate `/admin/*` endpoints (refresh-matches, poll-scores, results, seasons/close). | Cloud Run service via `--update-secrets=ADMIN_KEY=admin-key:latest` | Annually, or on suspected compromise |
| `session-secret` | HMAC key used to sign the `__session` cookie. A valid signature is the only proof a request belongs to a logged-in user.        | Cloud Run service via `--update-secrets=SESSION_SECRET=session-secret:latest` | Annually. Rotation invalidates all live sessions — schedule for low-traffic window. |

GitHub repo secrets are managed separately under repo Settings → Secrets and variables; see `infra/wif-setup.md`.

## Add a new secret

```bash
PROJECT='crew-predictions'           # repeat for crew-predictions-staging
NAME='my-new-secret'

# Create the secret container
gcloud secrets create "${NAME}" --project="${PROJECT}" --replication-policy=automatic

# Add the first version (read value from stdin to avoid shell history)
read -rs SECRET_VALUE
echo -n "${SECRET_VALUE}" | gcloud secrets versions add "${NAME}" \
  --project="${PROJECT}" --data-file=-
unset SECRET_VALUE
```

Then grant the Cloud Run runtime service account access:

```bash
# Today this is the default Compute Engine SA. After the runtime-SA
# least-privilege work in BACKLOG, replace with the dedicated runtime SA.
SA='937208344837-compute@developer.gserviceaccount.com'   # prod
# SA='99152086201-compute@developer.gserviceaccount.com'  # staging

gcloud secrets add-iam-policy-binding "${NAME}" \
  --project="${PROJECT}" \
  --member="serviceAccount:${SA}" \
  --role='roles/secretmanager.secretAccessor' \
  --condition=None
```

Reference it in the Cloud Run deploy command in `.github/workflows/ci.yml`:

```
--update-secrets=MY_NEW_SECRET=my-new-secret:latest
```

## Rotate an existing secret

```bash
PROJECT='crew-predictions'
NAME='session-secret'

# Generate a fresh value (32 random bytes, base64) and add as new version
openssl rand -base64 32 | tr -d '\n' | \
  gcloud secrets versions add "${NAME}" --project="${PROJECT}" --data-file=-

# Force Cloud Run to pick up the new version on next request:
gcloud run services update crew-predictions \
  --project="${PROJECT}" --region=us-east5 \
  --update-secrets=SESSION_SECRET=${NAME}:latest --quiet
```

Disable the old version after confirming the new one works:

```bash
gcloud secrets versions list "${NAME}" --project="${PROJECT}"
gcloud secrets versions disable <OLD_VERSION_NUMBER> --secret="${NAME}" --project="${PROJECT}"
```

## What NOT to put in Secret Manager

- Firebase web SDK config (`FIREBASE_API_KEY`, `FIREBASE_PROJECT_ID`, etc.) — these are public by design (served at `/auth/config.js`); they're stored as plain `--set-env-vars` in Cloud Run, not secrets.
- ESPN endpoints — public.
- Google Cloud project IDs — public.

If you're unsure whether a value is sensitive: treat it as sensitive until you've checked, then move it to plain env vars only after confirming it's intentionally public.
