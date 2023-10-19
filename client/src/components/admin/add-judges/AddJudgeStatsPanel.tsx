import { useEffect, useState } from 'react';
import AdminStat from '../AdminStat';
import { getRequest } from '../../../api';
import { errorAlert } from '../../../util';

interface JudgeStats {
    num: number;
    avg_votes: number;
    num_active: number;
}

const AddJudgeStatsPanel = () => {
    const [stats, setStats] = useState<JudgeStats>({ num: 0, avg_votes: 0, num_active: 0 });
    useEffect(() => {
        const fetchStats = async () => {
            const res = await getRequest('/judge/stats', 'admin');
            if (res.status === 500) {
                return;
            }
            if (res.status !== 200) {
                errorAlert(res);
                return;
            }
            setStats(res.data as JudgeStats);
        };

        fetchStats();
    }, []);

    return (
        <div className="flex flex-col justify-evenly w-full mt-8">
            <div className="flex justify-evenly basis-2/5">
                <AdminStat name="Total Judges" value={stats.num} />
                <AdminStat name="Average Votes" value={stats.avg_votes} />
                <AdminStat name="Active Judges" value={stats.num_active} />
            </div>
        </div>
    );
};

export default AddJudgeStatsPanel;
