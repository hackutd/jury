import { useEffect, useState } from 'react';
import AdminStat from './AdminStat';

interface Stats {
    projects: number;
    seen: number;
    votes: number;
    time: number;
    avg_mu: number;
    avg_sigma: number;
    judges: number;
}

const AdminStatsPanel = () => {
    const [stats, setStats] = useState<Stats>({
        projects: 0,
        seen: 0,
        votes: 0,
        time: 0,
        avg_mu: 0,
        avg_sigma: 0,
        judges: 0,
    });
    useEffect(() => {
        (async () => {
            const fetchedStats = await fetch(`${process.env.REACT_APP_JURY_URL}/admin/stats`, {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
            }).then((data) => data.json());
            setStats(fetchedStats);
        })();
    }, []);
    return (
        <div className="flex justify-evenly w-full mt-4">
            <AdminStat name="Projects" value={stats.projects} />
            <AdminStat name="Seen" value={stats.seen} />
            <AdminStat name="Votes" value={stats.votes} />
            <AdminStat name="Judging Time" value="01:46:23" />
            <AdminStat name="Average Mu" value={stats.avg_mu} />
            <AdminStat name="Average Sigma^2" value={stats.avg_sigma} />
            <AdminStat name="Judges" value={stats.judges} />
        </div>
    );
};

export default AdminStatsPanel;
