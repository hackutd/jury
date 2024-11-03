import PauseButton from '../../components/admin/PauseButton';
import ConfirmPopup from '../../components/ConfirmPopup';
import ToTopButton from '../../components/ToTopButton';
import AdminStatsPanel from '../../components/admin/AdminStatsPanel';
import AdminTable from '../../components/admin/tables/AdminTable';
import AdminToggleSwitch from '../../components/admin/AdminToggleSwitch';
import AdminToolbar from '../../components/admin/AdminToolbar';
import { postRequest } from '../../api';
import { errorAlert } from '../../util';
import JuryHeader from '../../components/JuryHeader';
import Button from '../../components/Button';
import { useNavigate } from 'react-router-dom';
import { useAdminStore, useClockStore, useFlagsStore, useOptionsStore } from '../../store';
import { useEffect, useState } from 'react';


    const AdminHeader = () =>{
        const navigate = useNavigate();
        const options = useOptionsStore((state) => state.options);
        const [swapPopup, setSwapPopup] = useState(false);
        const fetchJudges = useAdminStore((state) => state.fetchJudges);
        const fetchOptions = useOptionsStore((state) => state.fetchOptions);
        const clock = useClockStore((state) => state.clock);
        const [showProjects, setShowProjects] = useState(true);
        const [lastUpdate, setLastUpdate] = useState(new Date());
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
                    className="ml-4 py-2"
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
                                className="mr-36 py-2"
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
            <ToTopButton />
        </>
        
    }
    
    export default AdminHeader;