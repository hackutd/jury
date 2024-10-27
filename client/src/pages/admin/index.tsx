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
import { useAdminStore, useClockStore } from '../../store';

// TODO: Add FAB to 'return to top'
// TODO: Make pause button/settings have hover effects
const Admin = () => {
    const navigate = useNavigate();
    const fetchStats = useAdminStore((state) => state.fetchStats);
    const fetchClock = useClockStore((state) => state.fetchClock);
    const fetchProjects = useAdminStore((state) => state.fetchProjects);
    const fetchJudges = useAdminStore((state) => state.fetchJudges);

    const [showProjects, setShowProjects] = useState(true);
    const [loading, setLoading] = useState(true);
    const [lastUpdate, setLastUpdate] = useState(new Date());

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

    useEffect(() => {
        // Set 15 seconds interval to refresh admin stats and clock
        const refresh = setInterval(async () => {
            // Fetch stats and clock
            fetchStats();
            fetchClock();
            fetchProjects();
            fetchJudges();
            setLastUpdate(new Date());
        }, 15000);

        return () => clearInterval(refresh);
    }, []);

    if (loading) {
        return <Loading disabled={!loading} />;
    }

    const scrollToTop = () => {
        window.scrollTo({
            top: 0,
            behavior: 'smooth',
        });
    };

    return (
        <>
            <JuryHeader withLogout isAdmin />
            <Button
                type="outline"
                onClick={() => {
                    navigate('/admin/settings');
                }}
                className="absolute top-6 left-[16rem] w-40 md:w-52 text-lg py-2 px-1 hover:scale-100 focus:scale-100 rounded-md font-bold"
            >Settings</Button>
            <AdminStatsPanel />
            <AdminToggleSwitch state={showProjects} setState={setShowProjects} />
            <AdminToolbar showProjects={showProjects} lastUpdate={lastUpdate} />
            <AdminTable showProjects={showProjects} />
            <button
                onClick={scrollToTop}
                className="fixed bottom-8 right-8 w-12 h-12 border-2 border-primaryLight rounded-full flex items-center justify-center shadow-md bg-white"
            >
                <span className="text-light text-xl">↑</span>
            </button>
        </>
    );
};

export default Admin;
