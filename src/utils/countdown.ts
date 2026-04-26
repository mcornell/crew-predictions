export function formatCountdown(ms: number): string {
  if (ms <= 0) return 'kicks off now'
  const totalSeconds = Math.floor(ms / 1000)
  const days = Math.floor(totalSeconds / 86400)
  const hours = Math.floor((totalSeconds % 86400) / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  if (days > 0) return `locks in ${days}d ${hours}h`
  if (hours > 0) return `locks in ${hours}h ${minutes}m`
  const seconds = totalSeconds % 60
  return `locks in ${minutes}m ${seconds}s`
}
