# Workload Identity Federation (GitHub Actions ↔ GCP)

GitHub Actions authenticates to GCP without long-lived JSON service account keys via Workload Identity Federation (WIF). When a workflow runs, GitHub mints a short-lived OIDC token; GCP's WIF provider verifies the token's claims (specifically `repository == 'mcornell/crew-predictions'`) and exchanges it for a Google access token bound to a service account.

This doc captures the configuration so it can be re-created or audited without poking through three cloud consoles.

## Service account

`github-actions@crew-predictions.iam.gserviceaccount.com`

Lives in the **prod** project, but holds IAM bindings on **both** projects (cross-project access). The CI workflow references it via the `WIF_SERVICE_ACCOUNT` GitHub repo secret.

## Workload Identity Pool & Provider

Both live in the **prod** project at the global location.

| Resource | Value |
| --- | --- |
| Pool ID | `github` |
| Pool full name | `projects/937208344837/locations/global/workloadIdentityPools/github` |
| Provider ID | `github-provider` |
| Provider full name | `projects/937208344837/locations/global/workloadIdentityPools/github/providers/github-provider` |
| Issuer URI | `https://token.actions.githubusercontent.com` |
| Attribute condition | `attribute.repository == 'mcornell/crew-predictions'` |
| Attribute mapping | `attribute.repository = assertion.repository`, `google.subject = assertion.sub` |

The attribute condition is the security boundary: only workflow runs originating from this specific GitHub repo can impersonate the service account. Forking the repo or running from a different repo will fail token exchange.

## Required IAM bindings

The `github-actions@` service account needs these roles. All bindings are unconditional (`--condition=None`).

### Both projects (`crew-predictions` and `crew-predictions-staging`)

| Role | Purpose |
| --- | --- |
| `roles/run.admin` | Deploy/update Cloud Run services |
| `roles/iam.serviceAccountUser` | Act as the Cloud Run runtime service account during deploy |
| `roles/storage.admin` | Push/pull frontend artifacts to/from `crew-predictions-frontend` GCS bucket; read/write Firebase Hosting build artifacts |
| `roles/artifactregistry.admin` | Push Docker images to `us-east5-docker.pkg.dev/crew-predictions/crew-predictions/` and manage cleanup tags |
| `roles/firebasehosting.admin` | Deploy frontend to Firebase Hosting |
| `roles/serviceusage.serviceUsageConsumer` | Required preflight for Firebase CLI deploys (`firestore`, `functions`) |
| `roles/firebaserules.admin` | Deploy `firestore.rules` |
| `roles/datastore.indexAdmin` | Deploy `firestore.indexes.json` |
| `roles/secretmanager.secretAccessor` | Read secrets during deploy when wiring `--update-secrets` |

### Prod project only

| Role | Purpose |
| --- | --- |
| `roles/cloudfunctions.developer` | Deploy `infra/billing-killswitch/` Cloud Function |

## Granting a new role

The first time you grant a role on a project that has any conditional binding, gcloud forces you to specify a condition. Always pass `--condition=None` for unconditional bindings.

```bash
SA='github-actions@crew-predictions.iam.gserviceaccount.com'
PROJECT='crew-predictions-staging'   # or crew-predictions
ROLE='roles/firebaserules.admin'

gcloud projects add-iam-policy-binding "${PROJECT}" \
  --member="serviceAccount:${SA}" \
  --role="${ROLE}" \
  --condition=None
```

## Recovery: re-create WIF from scratch

If the pool or provider is deleted, GitHub Actions will fail every deploy with a 401 from the token exchange. Recovery:

```bash
PROJECT_ID='crew-predictions'
PROJECT_NUMBER='937208344837'
REPO='mcornell/crew-predictions'

# Pool
gcloud iam workload-identity-pools create github \
  --project="${PROJECT_ID}" --location=global \
  --display-name='GitHub Actions'

# Provider
gcloud iam workload-identity-pools providers create-oidc github-provider \
  --project="${PROJECT_ID}" --location=global \
  --workload-identity-pool=github \
  --display-name='GitHub Provider' \
  --attribute-mapping='google.subject=assertion.sub,attribute.repository=assertion.repository' \
  --attribute-condition="attribute.repository=='${REPO}'" \
  --issuer-uri='https://token.actions.githubusercontent.com'

# Allow the SA to be impersonated by tokens matching the repo
gcloud iam service-accounts add-iam-policy-binding \
  "github-actions@${PROJECT_ID}.iam.gserviceaccount.com" \
  --project="${PROJECT_ID}" \
  --role='roles/iam.workloadIdentityUser' \
  --member="principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/github/attribute.repository/${REPO}"
```

Then update the GitHub repo secrets `WIF_PROVIDER` and `WIF_SERVICE_ACCOUNT` to match the new resources (the pool / provider full names from above and the SA email).

## GitHub repo secrets used by CI

| Secret | Value shape | Notes |
| --- | --- | --- |
| `WIF_PROVIDER` | `projects/937208344837/locations/global/workloadIdentityPools/github/providers/github-provider` | Pool + provider full path |
| `WIF_SERVICE_ACCOUNT` | `github-actions@crew-predictions.iam.gserviceaccount.com` | Service account email |
| `STAGING_FIREBASE_API_KEY` | `AIza...` | Public Firebase web SDK API key for staging — secret because we don't want it indexed in commit history, but it's served to clients at runtime |
| `PROD_FIREBASE_API_KEY` | `AIza...` | Same, for prod |
| `SMOKE_TEST_PASSWORD` | string | Password for the seeded staging smoke-test accounts |
