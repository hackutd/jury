import { useEffect } from 'react';
import StatBlock from '../StatBlock';
import { useAdminStore } from '../../store';
import AdminClock from './AdminClock';

const AdminStatsPanelMobile = () => {
    const stats = useAdminStore((state) => state.stats);
    const fetchStats = useAdminStore((state) => state.fetchStats);

    // Fetch stats on load
    useEffect(() => {
        fetchStats();
    }, []);

    return (
        <div className="flex md:hidden flex-col mt-4 w-full gap-2">
            <AdminClock />
            <div className="grid grid-cols-2 gap-2">
                <StatBlock name="Projects" value={stats.projects} />
                <StatBlock name="Judges" value={stats.judges} />
                <StatBlock name="Avg Project Seen" value={stats.avg_project_seen} />
                <StatBlock name="Average Judge Seen" value={stats.avg_judge_seen} />
            </div>
        </div>
    );
};

export default AdminStatsPanelMobile;
