import { useEffect, useState } from 'react';
import Loading from '../../components/Loading';
import { postRequest } from '../../api';
import { errorAlert } from '../../util';
import { useNavigate } from 'react-router-dom';
import AdminHeader  from '../admin/AdminHeader';

import { useAdminStore, useClockStore, useFlagsStore, useOptionsStore } from '../../store';


const Admin = () => {
    const navigate = useNavigate();
    const fetchStats = useAdminStore((state) => state.fetchStats);
    const clock = useClockStore((state) => state.clock);
    const fetchClock = useClockStore((state) => state.fetchClock);
    const fetchProjects = useAdminStore((state) => state.fetchProjects);
    const fetchJudges = useAdminStore((state) => state.fetchJudges);
    const options = useOptionsStore((state) => state.options);
    const fetchOptions = useOptionsStore((state) => state.fetchOptions);
    const fetchFlags = useFlagsStore((state) => state.fetchFlags);

    const [showProjects, setShowProjects] = useState(true);
    const [loading, setLoading] = useState(true);
    const [lastUpdate, setLastUpdate] = useState(new Date());
    const [swapPopup, setSwapPopup] = useState(false);

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
        fetchOptions();
        fetchFlags();
    }, []);

    useEffect(() => {
        // Set 15 seconds interval to refresh admin stats and clock
        const refresh = setInterval(async () => {
            // Fetch stats and clock
            fetchStats();
            fetchClock();
            fetchProjects();
            fetchJudges();
            fetchOptions();
            fetchFlags();
            setLastUpdate(new Date());
        }, 15000);

        return () => clearInterval(refresh);
    }, []);

    const swapJudgeGroups = async () => {
        const res = await postRequest<OkResponse>('/admin/groups/swap', 'admin', null);
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }
        setSwapPopup(false);
        fetchJudges();
        fetchOptions();
        alert('Judge groups swapped successfully');
    };

    if (loading) {
        return <Loading disabled={!loading} />;
    }
    return AdminHeader;
};

export default Admin;
