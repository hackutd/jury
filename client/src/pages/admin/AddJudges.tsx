import JuryHeader from '../../components/JuryHeader';
import AddJudgeStatsPanel from '../../components/admin/add-judges/AddJudgeStatsPanel';
import JudgeQrCodes from '../../components/admin/add-judges/JudgeQrCodes';
import NewJudgeForm from '../../components/admin/add-judges/NewJudgeForm';
import UploadCSVForm from '../../components/admin/UploadCSVForm';

const AddJudges = () => {
    return (
        <>
            <JuryHeader withBack withLogout isAdmin />
            <div className="flex flex-col items-start justify-center w-full px-8 py-4 md:px-16 md:py-8">
                <h1 className="text-4xl font-bold">Add Judges</h1>
                <AddJudgeStatsPanel />
                <div className="mt-8 flex flex-col w-full">
                    <JudgeQrCodes />
                    <NewJudgeForm />
                    <UploadCSVForm format="judge" />
                </div>
            </div>
        </>
    );
};

export default AddJudges;
