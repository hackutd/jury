import { useEffect, useState } from 'react';
import AdminStat from './AdminStat';
import PauseButton from '../../admin/PauseButton';
import useAdminStore from '../../store';

// Jank global variables for keeping track of drift in timer
let start: number;
let drift = 0;

const AdminStatsPanel = () => {
    const stats = useAdminStore((state) => state.stats);
    const fetchStats = useAdminStore((state) => state.fetchStats);
    const [time, setTime] = useState<number>(0);
    const [paused, setPaused] = useState(true);

    const handleEventSourceMessage = (e: MessageEvent) => {
        // Get the data from the message
        const data = e.data as string;

        if (data.trim() === 'clock') {
            fetchClock();
        } else {
            fetchStats();
        }
    };

    // Fetch clock info
    const fetchClock = async () => {
        const data = await fetch(`${process.env.REACT_APP_JURY_URL}/admin/clock`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
        }).then((data) => data.json());

        // If the clock is paused, set paused state to true
        setPaused(data.paused);

        // Set clock time
        if (data.start === 0 || data.paused) {
            setTime(data.prev);
        } else {
            setTime(data.prev + (Date.now() - data.start));
        }
    };

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

    // Add event source listener for when to sync stats (happens on DB change)
    // Also fetch stats on load
    useEffect(() => {
        const eventSource = new EventSource(`${process.env.REACT_APP_JURY_URL}/admin/sync`, {
            withCredentials: true,
        });
        eventSource.onmessage = handleEventSourceMessage;
        fetchStats();
        fetchClock();

        // eslint-disable-next-line
    }, []);

    // Timer for updating time each second
    useEffect(() => {
        // Start client timer
        const interval = setInterval(() => {
            // Don't update if paused
            if (paused) return;

            // Calculate the drift based on the time since the last update
            if (start) drift = Date.now() - start - 1000;

            // Update the time and set start to the current time
            start = Date.now();
            setTime(time + 1000 + drift);
        }, 1000);

        return () => clearInterval(interval);
    }, [paused, time]);

    return (
        <div className="flex flex-row mt-8 w-full">
            <PauseButton paused={paused} setPaused={setPaused} clock={time} />
            <div className="flex justify-evenly basis-2/5">
                <AdminStat name="Projects" value={stats.projects} />
                <AdminStat name="Seen" value={stats.seen} />
                <AdminStat name="Votes" value={stats.votes} />
            </div>
            <AdminStat
                name="Judging Time"
                value={msToTime(time)}
                className={'basis-1/5' + (paused ? ' text-error' : '')}
            />
            <div className="flex justify-evenly basis-2/5">
                <AdminStat name="Average Mu" value={stats.avg_mu} />
                <AdminStat name="Average Sigma^2" value={stats.avg_sigma} />
                <AdminStat name="Judges" value={stats.judges} />
            </div>
        </div>
    );
};

export default AdminStatsPanel;
