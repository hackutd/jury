import JuryHeader from '../components/JuryHeader';
import AddJudgeStatsPanel from '../components/admin/add-judges/AddJudgeStatsPanel';
import NewJudgeForm from '../components/admin/add-judges/NewJudgeForm';
import UploadCSVForm from '../components/admin/add-judges/UploadCSVForm';

const AddJudges = () => {
    return (
        <>
            <JuryHeader withLogout isAdmin />
            <div className="flex flex-col items-start justify-center w-full px-8 py-4 md:px-16 md:py-8">
                <h1 className="text-4xl font-bold">Add Judges</h1>
                <AddJudgeStatsPanel />
                <div className="mt-8 flex flex-col w-full space-y-8">
                    <NewJudgeForm />
                    <UploadCSVForm format="judge" />
                </div>
            </div>
        </>
    );
};

export default AddJudges;
