import JuryHeader from '../components/JuryHeader';
import AddJudgesPanel from '../components/admin/add-judges/AddJudgesPanel';

const AddJudges = () => {
    return (
        <>
            <JuryHeader withLogout isAdmin />
            <AddJudgesPanel />
        </>
    );
};

export default AddJudges;
