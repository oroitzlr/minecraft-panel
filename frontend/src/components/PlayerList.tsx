import type { Player } from '../types'

interface Props {
    players: Player[]
}

export function PlayerList({ players }: Props) {
    return (
        <div className="panel">
            <div className="ph">
                <h2>Joueurs</h2>
                <span id="playerBadge" className="count-badge tab">{players.length}</span>
            </div>

            {players.length === 0 ? (
                <div id="playerEmpty">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                        <circle cx="12" cy="8" r="4" />
                        <path d="M4 21c0-4 4-6 8-6s8 2 8 6" />
                    </svg>
                    <p>Aucun joueur connecté</p>
                </div>
            ) : (
                <ul id="playerList">
                    {players.map((player) => (
                        <li key={player.name} className="pl-item">
                            <div className="pl-avatar">
                                <img
                                    src={`https://mc-heads.net/avatar/${player.name}/30`}
                                    alt={player.name}
                                    width={30}
                                    height={30}
                                />
                            </div>
                            <div className="pl-info">
                                <div className="pl-name">{player.name}</div>
                            </div>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    )
}