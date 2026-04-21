import { deleteSmokeNewAccounts } from './smoke-accounts'

export default async function globalSetup() {
  await deleteSmokeNewAccounts()
}
