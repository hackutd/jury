import AdminStatsPanel from '../components/admin/AdminStatsPanel';
import JuryHeader from '../components/JuryHeader';

const Admin = () => {
    return (
        <>
            <JuryHeader withLogout isAdmin />
            <AdminStatsPanel />
        </>
    );
};

export default Admin;
