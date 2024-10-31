import { useEffect } from 'react';
import StatBlock from '../../StatBlock';
import { useAdminStore } from '../../../store';

const AddProjectsStatsPanel = () => {
    const stats = useAdminStore((state) => state.projectStats);
    const fetchStats = useAdminStore((state) => state.fetchProjectStats);

    useEffect(() => {
        fetchStats();
    }, []);

    return (
        <div className="flex flex-col justify-evenly w-full mt-8">
            <div className="flex justify-evenly basis-2/5">
                <StatBlock name="Active Projects" value={stats.num} />
                <StatBlock name="Average Seen" value={stats.avg_seen} />
            </div>
        </div>
    );
};

export default AddProjectsStatsPanel;
