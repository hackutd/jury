import { useEffect } from 'react';
import StatBlock from '../../StatBlock';
import { useAdminStore } from '../../../store';

const AddJudgeStatsPanel = () => {
    const stats = useAdminStore((state) => state.judgeStats);
    const fetchStats = useAdminStore((state) => state.fetchJudgeStats);

    useEffect(() => {
        fetchStats();
    }, []);

    return (
        <div className="flex flex-col justify-evenly w-full mt-8">
            <div className="flex justify-evenly basis-2/5">
                <StatBlock name="Total Judges" value={stats.num} />
                <StatBlock name="Average Seen" value={stats.avg_seen} />
                <StatBlock name="Active Judges" value={stats.num_active} />
            </div>
        </div>
    );
};

export default AddJudgeStatsPanel;
