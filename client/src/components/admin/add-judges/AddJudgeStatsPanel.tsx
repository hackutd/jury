import { useEffect, useState } from 'react';
import AddJudgeStat from './AddJudgeStat';

interface JudgeStats {
    num: number;
    avg_votes: number;
    num_active: number;
}

const AddJudgeStatsPanel = () => {
    const [stats, setStats] = useState<JudgeStats>({ num: 0, avg_votes: 0, num_active: 0 });
    useEffect(() => {
        const fetchStats = async () => {
            const fetchedStats = await fetch(`${process.env.REACT_APP_JURY_URL}/judge/stats`, {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
            }).then((data) => data.json());
            // TODO: Add error checking
            setStats(fetchedStats);
        };

        fetchStats();
    }, []);

    return (
        <div className="flex flex-col justify-evenly w-full mt-8">
            <div className="flex justify-evenly basis-2/5">
                <AddJudgeStat name="Total Judges" value={stats.num} />
                <AddJudgeStat name="Average Votes" value={stats.avg_votes} />
                <AddJudgeStat name="Active Judges" value={stats.num_active} />
            </div>
        </div>
    );
};

export default AddJudgeStatsPanel;
