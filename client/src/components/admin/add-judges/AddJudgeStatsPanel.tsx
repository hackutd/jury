import AddJudgeStat from './AddJudgeStat';

const AddJudgeStatsPanel = () => {
    return (
        <div className="flex flex-col justify-evenly w-full mt-8">
            <div className="flex justify-evenly basis-2/5">
                <AddJudgeStat name="Active Judges" value={100} />
                <AddJudgeStat name="Average Votes" value={21.2} />
                <AddJudgeStat name="Average Seen" value={21.2} />
            </div>
        </div>
    );
};

export default AddJudgeStatsPanel;
