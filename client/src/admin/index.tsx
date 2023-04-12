import { useState } from 'react';
import AdminStatsPanel from '../components/admin/AdminStatsPanel';
import AdminToggleSwitch from '../components/admin/AdminToggleSwitch';
import JuryHeader from '../components/JuryHeader';

const Admin = () => {
    const [showProjects, setShowProjects] = useState(false);
    return (
        <>
            <JuryHeader withLogout isAdmin />
            <AdminStatsPanel />
            <AdminToggleSwitch state={showProjects} setState={setShowProjects} />
        </>
    );
};

export default Admin;
