import type { SystemStats } from '../types'

interface Props {
    stats: SystemStats | null
}

function StatBar({ label, value, max, unit, text }: {
    label: string
    value: number
    max: number
    unit: string
    text: string
}) {
    const percent = Math.min((value / max) * 100, 100)
    const level = percent > 90 ? 'crit' : percent > 70 ? 'warn' : undefined

    return (
        <div className="panel stat">
            <div className="stat-head">
                <span className="l">{label}</span>
                <span className="stat-val tab">{text}</span>
            </div>
            <div className="track">
                <div
                    className="stat-bar"
                    data-level={level}
                    style={{ width: `${percent}%` }}
                />
            </div>
        </div>
    )
}

export function StatsBar({ stats }: Props) {
    if (!stats) return null

    return (
        <section className="stats">
            <StatBar
                label="Mémoire RAM"
                value={stats.ramUsed}
                max={stats.ramTotal}
                unit="Go"
                text={`${stats.ramUsed} Mo / ${stats.ramTotal} Mo`}
            />
            <StatBar
                label="Processeur"
                value={stats.cpuPercent}
                max={100}
                unit="%"
                text={`${stats.cpuPercent.toFixed(1)}%`}
            />
            <StatBar
                label="Disque"
                value={stats.diskUsed}
                max={stats.diskTotal}
                unit="Go"
                text={`${stats.diskUsed} Go / ${stats.diskTotal} Go`}
            />
        </section>
    )
}