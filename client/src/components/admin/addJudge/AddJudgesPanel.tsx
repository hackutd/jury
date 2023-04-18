import AddJudgeStatsPanel from './AddJudgeStatsPanel';
import UploadCSVForm from './UploadCSVForm';
import NewJudgeForm from './NewJudgeForm';

const AddJudgesPanel = () => {
    return (
        <div className="flex flex-col items-start justify-center w-full px-8 py-4 md:px-16 md:py-8">
            <h1 className="text-4xl font-bold">Add Judges</h1>
            <AddJudgeStatsPanel />
            <div className="mt-8 flex flex-col w-full space-y-8">
                <NewJudgeForm />
                <UploadCSVForm />
            </div>
        </div>
    );
};

export default AddJudgesPanel;
