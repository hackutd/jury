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
import { useAdminStore, useClockStore, useFlagsStore, useOptionsStore } from '../../store';
import PauseButton from '../../components/admin/PauseButton';
import ConfirmPopup from '../../components/ConfirmPopup';

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
    return (
        <>
            <JuryHeader withLogout isAdmin />
            <div className="absolute top-4 w-full flex items-center z-0">
                <PauseButton />
                <Button
                    type="outline"
                    onClick={() => {
                        navigate('/admin/settings');
                    }}
                    small
                    bold
                    className="ml-4"
                >
                    Settings
                </Button>
                {options &&
                    options.multi_group &&
                    options.main_group.switching_mode === 'manual' && (
                        <>
                            <p className="text-lg text-right grow mr-4 text-lighter">
                                <span className="font-bold">Swaps: </span>
                                {options.main_group.manual_switches}
                            </p>
                            <Button
                                type="gold"
                                onClick={setSwapPopup.bind(null, true)}
                                disabled={clock.running}
                                tooltip="Groups can only be swapped when judging is paused"
                                small
                                bold
                                className="mr-36"
                            >
                                Swap Judge Groups
                            </Button>
                        </>
                    )}
            </div>
            <AdminStatsPanel />
            <AdminToggleSwitch state={showProjects} setState={setShowProjects} />
            <AdminToolbar showProjects={showProjects} lastUpdate={lastUpdate} />
            <AdminTable showProjects={showProjects} />
            <ConfirmPopup
                enabled={swapPopup}
                setEnabled={setSwapPopup}
                onSubmit={swapJudgeGroups}
                submitText="Swap"
                title="Swap Judge Groups"
                red
            >
                Are you sure you want to swap judge groups? This will increment the group number of
                every judge.
            </ConfirmPopup>
        </>
    );
};

export default Admin;
