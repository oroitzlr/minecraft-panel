import { useRef, useEffect, useState } from 'react'
import { useWebSocket } from '../hooks/useWebSocket'

export function Console() {
    const { messages, connected, sendMessage } = useWebSocket()
    const [input, setInput] = useState('')
    const [autoScroll, setAutoScroll] = useState(true)
    const consoleRef = useRef<HTMLDivElement>(null)

    // Auto-scroll vers le bas quand nouveau message
    useEffect(() => {
        if (autoScroll && consoleRef.current) {
            consoleRef.current.scrollTop = consoleRef.current.scrollHeight
        }
    }, [messages, autoScroll])

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        if (!input.trim()) return
        sendMessage(input.trim())
        setInput('')
    }

    const getLevel = (line: string) => {
        if (line.includes('ERROR')) return 'ERROR'
        if (line.includes('WARN')) return 'WARN'
        return 'INFO'
    }

    const clearConsole = () => {
        // On vide juste l'affichage local
        consoleRef.current!.innerHTML = ''
    }

    return (
        <div className="panel console-panel">
            <div className="console-head">
                <h2>Console</h2>
                <div className="opts">
                    <label>
                        <input
                            type="checkbox"
                            checked={autoScroll}
                            onChange={e => setAutoScroll(e.target.checked)}
                        />
                        auto-défilement
                    </label>
                    <button id="clearConsole" onClick={clearConsole}>effacer</button>
                    <span style={{ color: connected ? '#6cbf3a' : '#d24b3e', fontSize: '16px' }}>
                        {connected ? '● connecté' : '● déconnecté'}
                    </span>
                </div>
            </div>

            <div id="console" ref={consoleRef}>
                {messages.map((msg: string, i: number) => (
                    <div key={i} className="cl-line" data-level={getLevel(msg)}>
                        <span className="cl-msg">{msg}</span>
                    </div>
                ))}
            </div>

            <form id="rconForm" className="rcon-form" onSubmit={handleSubmit}>
                <span className="pr">&gt;</span>
                <input
                    id="rconInput"
                    type="text"
                    value={input}
                    onChange={e => setInput(e.target.value)}
                    disabled={!connected}
                    placeholder="list · say bonjour · time set day"
                    autoComplete="off"
                    spellCheck={false}
                />
                <button type="submit" className="send" disabled={!connected}>
                    Envoyer
                </button>
            </form>
        </div>
    )
}