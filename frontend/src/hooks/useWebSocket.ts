import { useState, useEffect, useRef, useCallback } from 'react'

const WS_URL = import.meta.env.VITE_WS_URL ?? 'ws://localhost:8080'

export function useWebSocket() {
  const [messages, setMessages] = useState<string[]>([])
  const [connected, setConnected] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    const token = localStorage.getItem('token')
    const ws = new WebSocket(`${WS_URL}/api/ws/console?token=${token}`)
    wsRef.current = ws

    ws.onopen = () => setConnected(true)

    ws.onmessage = (event) => {
      setMessages(prev => [...prev, event.data])
    }

    ws.onclose = () => setConnected(false)

    ws.onerror = () => setConnected(false)

    return () => ws.close()
  }, [])

  const sendMessage = useCallback((msg: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(msg)
    }
  }, [])

  return { messages, connected, sendMessage }
}