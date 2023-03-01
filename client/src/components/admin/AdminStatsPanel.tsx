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
        <div className="flex flex-row mt-8 w-full">
            <div className="flex justify-evenly basis-2/5">
                <AdminStat name="Projects" value={stats.projects} />
                <AdminStat name="Seen" value={stats.seen} />
                <AdminStat name="Votes" value={stats.votes} />
            </div>
            <AdminStat name="Judging Time" value="01:46:23" className="basis-1/5" />
            <div className="flex justify-evenly basis-2/5">
                <AdminStat name="Average Mu" value={stats.avg_mu} />
                <AdminStat name="Average Sigma^2" value={stats.avg_sigma} />
                <AdminStat name="Judges" value={stats.judges} />
            </div>
        </div>
    );
};

export default AdminStatsPanel;
