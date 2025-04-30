import { useEffect, useRef } from 'react';
import { useClockStore } from '../../store';
import StatBlock from '../StatBlock';

const AdminClock = () => {
    const clock = useClockStore((state) => state.clock);
    const fetchClock = useClockStore((state) => state.fetchClock);
    const setTime = useClockStore((state) => state.setTime);
    const start = useRef<number>(0);

    useEffect(() => {
        fetchClock();
    }, []);

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

    // Timer for updating time each second
    useEffect(() => {
        // Start client timer
        const interval = setInterval(() => {
            // Don't update if paused
            if (!clock.running) return;

            // Calculate the drift based on the time since the last update
            const drift = start.current === 0 ? 0 : Date.now() - start.current - 1000;

            // Update the time and set start to the current time
            start.current = Date.now();
            setTime(clock.time + 1000 + drift);
        }, 1000);

        return () => clearInterval(interval);
    }, [clock]);

    return (
        <StatBlock
            name="Judging Time"
            value={msToTime(clock.time)}
            className={'basis-1/5 ' + (!clock.running && 'text-error')}
        />
    );
};

export default AdminClock;
