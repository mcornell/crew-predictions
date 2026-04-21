import { deleteSmokeNewAccounts } from './smoke-accounts'

export default async function globalTeardown() {
  await deleteSmokeNewAccounts()
}
