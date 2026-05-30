import { useState, useEffect } from 'react'
import { fetchStatus, fetchPlayers } from '../api/server'
import type { ServerStatus, Player } from '../types'

export function useServerStatus(interval = 5000) {
  const [status, setStatus] = useState<ServerStatus | null>(null)
  const [players, setPlayers] = useState<Player[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetch = async () => {
      try {
        const [s, p] = await Promise.all([fetchStatus(), fetchPlayers()])
        setStatus(s)
        setPlayers(p)
        setError(null)
      } catch (err) {
        setError('Impossible de contacter le backend')
      } finally {
        setLoading(false)
      }
    }

    fetch()
    const timer = setInterval(fetch, interval)
    return () => clearInterval(timer)
  }, [interval])

  return { status, players, loading, error }
}