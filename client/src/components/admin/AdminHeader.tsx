import { useNavigate } from 'react-router-dom';
import { useAdminStore, useClockStore, useOptionsStore } from '../../store';
import Button from '../Button';
import PauseButton from './PauseButton';
import ConfirmPopup from '../ConfirmPopup';
import { useState } from 'react';
import { postRequest } from '../../api';
import { errorAlert } from '../../util';

const AdminHeader = () => {
    const navigate = useNavigate();
    const options = useOptionsStore((state) => state.options);
    const clock = useClockStore((state) => state.clock);
    const fetchJudges = useAdminStore((state) => state.fetchJudges);
    const fetchOptions = useOptionsStore((state) => state.fetchOptions);
    const [swapPopup, setSwapPopup] = useState(false);

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

    return (
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
            <Button
                type="outline"
                onClick={() => navigate('/admin/log')}
                small
                bold
                className="ml-4 py-2"
            >
                Audit Log
            </Button>
            {options && options.multi_group && options.switching_mode === 'manual' && (
                <>
                    <p className="text-lg text-right grow mr-4 text-lighter">
                        <span className="font-bold">Swaps: </span>
                        {options.manual_switches}
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
        </div>
    );
};

export default AdminHeader;
