import { useEffect } from 'react';
import StatBlock from '../StatBlock';
import { useAdminStore } from '../../store';
import AdminClock from './AdminClock';

// Displays the admin stats, will be hidden on mobile
const AdminStatsPanel = () => {
    const stats = useAdminStore((state) => state.stats);
    const fetchStats = useAdminStore((state) => state.fetchStats);

    // Fetch stats on load
    useEffect(() => {
        fetchStats();
    }, []);

    return (
        <div className="hidden md:flex flex-row mt-8 w-full">
            <div className="flex justify-evenly basis-2/5">
                <StatBlock name="Projects" value={stats.projects} />
                <StatBlock name="Avg Project Seen" value={stats.avg_project_seen} />
            </div>
            <AdminClock />
            <div className="flex justify-evenly basis-2/5">
                <StatBlock name="Average Judge Seen" value={stats.avg_judge_seen} />
                <StatBlock name="Judges" value={stats.judges} />
            </div>
        </div>
    );
};

export default AdminStatsPanel;
