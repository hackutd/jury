import { useEffect, useState } from 'react';
import StatBlock from '../StatBlock';
import PauseButton from './PauseButton';
import { useAdminStore, useClockStore } from '../../store';

// Jank global variables for keeping track of drift in timer
let start: number;
let drift = 0;

const AdminStatsPanel = () => {
    const stats = useAdminStore((state) => state.stats);
    const fetchStats = useAdminStore((state) => state.fetchStats);
    const clock = useClockStore((state) => state.clock);
    const fetchClock = useClockStore((state) => state.fetchClock);
    const setTime = useClockStore((state) => state.setTime);

    // Convert ms time to time string
    const msToTime = (ms: number) => {
        const seconds = Math.floor((ms / 1000) % 60);
        const minutes = Math.floor((ms / (1000 * 60)) % 60);
        const hours = Math.floor((ms / (1000 * 60 * 60)) % 24);

        return `${pad(hours)}:${pad(minutes)}:${pad(seconds)}`;
    };

    // Pad numbers to 2 digits
    const pad = (num: number) => {
        return num.toString().padStart(2, '0');
    };

    // Fetch clock and stats on load
    useEffect(() => {
        fetchStats();
        fetchClock();
    }, []);

    // Timer for updating time each second
    useEffect(() => {
        // Start client timer
        const interval = setInterval(() => {
            // Don't update if paused
            if (!clock.running) return;

            // Calculate the drift based on the time since the last update
            if (start) drift = Date.now() - start - 1000;

            // Update the time and set start to the current time
            start = Date.now();
            setTime(clock.time + 1000 + drift);
        }, 1000);

        return () => clearInterval(interval);
    }, [clock]);

    return (
        <div className="flex flex-row mt-8 w-full">
            <div className="flex justify-evenly basis-2/5">
                <StatBlock name="Projects" value={stats.projects} />
                <StatBlock name="Avg Project Seen" value={stats.avg_project_seen} />
            </div>
            <StatBlock
                name="Judging Time"
                value={msToTime(clock.time)}
                className={'basis-1/5 ' + (!clock.running && 'text-error')}
            />
            <div className="flex justify-evenly basis-2/5">
                <StatBlock name="Average Judge Seen" value={stats.avg_judge_seen} />
                <StatBlock name="Judges" value={stats.judges} />
            </div>
        </div>
    );
};

export default AdminStatsPanel;
