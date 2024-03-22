import { useEffect, useState } from 'react';
import AdminStatsPanel from '../../components/admin/AdminStatsPanel';
import AdminTable from '../../components/admin/tables/AdminTable';
import AdminToggleSwitch from '../../components/admin/AdminToggleSwitch';
import AdminToolbar from '../../components/admin/AdminToolbar';
import JuryHeader from '../../components/JuryHeader';
import Loading from '../../components/Loading';
import { postRequest } from '../../api';
import { errorAlert } from '../../util';
import { useNavigate } from 'react-router-dom';
import Button from '../../components/Button';

const Admin = () => {
    const navigate = useNavigate();
    const [showProjects, setShowProjects] = useState(true);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // Check if user logged in
        async function checkLoggedIn() {
            const loggedInRes = await postRequest<OkResponse>('/admin/auth', 'admin', null);
            if (loggedInRes.status === 401) {
                console.error(`Admin is not logged in!`);
                navigate('/admin/login');
                return;
            }
            if (loggedInRes.status === 200) {
                setLoading(false);
                return;
            }

            errorAlert(loggedInRes);
        }

        checkLoggedIn();
    }, []);

    if (loading) {
        return <Loading disabled={!loading} />;
    }
    return (
        <>
            <JuryHeader withLogout isAdmin />
            <Button
                type="outline"
                onClick={() => {
                    navigate('/admin/settings');
                }}
                className="fixed top-6 left-[16rem] w-40 md:w-52 text-lg py-2 px-1 hover:scale-100 focus:scale-100 rounded-md font-bold"
            >Settings</Button>
            <AdminStatsPanel />
            <AdminToggleSwitch state={showProjects} setState={setShowProjects} />
            <AdminToolbar showProjects={showProjects} />
            <AdminTable showProjects={showProjects} />
        </>
    );
};

export default Admin;
