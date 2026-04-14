import { useEffect, useState } from 'react';
import { Helmet } from 'react-helmet';
import { useNavigate } from 'react-router-dom';
import { postRequest } from '../../api';
import { errorAlert } from '../../util';
import { useAdminStore } from '../../store';
import AdminStatsPanel from '../../components/admin/AdminStatsPanel';
import AdminTable from '../../components/admin/tables/AdminTable';
import AdminToggleSwitch from '../../components/admin/AdminToggleSwitch';
import AdminToolbar from '../../components/admin/AdminToolbar';
import JuryHeader from '../../components/JuryHeader';
import Loading from '../../components/Loading';
import ToTopButton from '../../components/ToTopButton';
import AdminHeader from '../../components/admin/AdminHeader';
import AdminStatsPanelMobile from '../../components/admin/AdminStatsPanelMobile';
import AdminHeaderMobile from '../../components/admin/AdminHeaderMobile';
import ProjectsTable from '../../components/admin/tables/ProjectsTable';
import JudgesTable from '../../components/admin/tables/JudgesTable';

const Admin = () => {
    const navigate = useNavigate();
    const fetchDashboard = useAdminStore((state) => state.fetchDashboard);

    const [showProjects, setShowProjects] = useState(true);
    const [loading, setLoading] = useState(true);
    const [lastUpdate, setLastUpdate] = useState(new Date());

    useEffect(() => {
        // Check if user logged in, then load all dashboard data in one request.
        async function checkLoggedIn() {
            const loggedInRes = await postRequest<OkResponse>('/admin/auth', 'admin', null);
            if (loggedInRes.status === 401) {
                console.error(`Admin is not logged in!`);
                navigate('/admin/login');
                return;
            }
            if (loggedInRes.status === 200) {
                setLoading(false);
                fetchDashboard();
                return;
            }

            errorAlert(loggedInRes);
        }

        checkLoggedIn();
    }, []);

    useEffect(() => {
        // Poll once every 15 seconds using a single /admin/dashboard request
        // instead of the previous 6 separate API calls.
        const refresh = setInterval(async () => {
            fetchDashboard();
            setLastUpdate(new Date());
        }, 15000);

        return () => clearInterval(refresh);
    }, []);

    if (loading) {
        return <Loading disabled={!loading} />;
    }
    return (
        <>
            <Helmet>
                <title>Admin | Jury</title>
            </Helmet>
            <JuryHeader withLogout isAdmin />
            <AdminHeader />
            <AdminHeaderMobile />
            <AdminStatsPanel />
            <AdminStatsPanelMobile />
            <AdminToggleSwitch state={showProjects} setState={setShowProjects} />
            <AdminToolbar showProjects={showProjects} lastUpdate={lastUpdate} />
            {showProjects ? <ProjectsTable /> : <JudgesTable />}
            <ToTopButton />
        </>
    );
};

export default Admin;
