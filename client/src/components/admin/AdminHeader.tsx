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

    const startDeliberation = async () => {
        const res = await postRequest<OkResponse>('/admin/deliberation', 'admin', {
            start: !options.deliberation,
        });
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        fetchOptions();
        if (options.deliberation) {
            alert('Deliberations has been stopped. Judges can now edit their rankings.');
        } else {
            alert(
                'Deliberation has started! Judges will no longer be able to edit their rankings.'
            );
        }
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
            <div className="grow mr-36 flex flex-row justify-end items-center">
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
                            className="py-2"
                        >
                            Swap Judge Groups
                        </Button>
                    </>
                )}
                <Button
                    type={options.deliberation ? 'outline' : 'error'}
                    onClick={startDeliberation}
                    small
                    bold
                    tooltip="Stop judges from editing their rankings during deliberation"
                    className="ml-4 py-2"
                >
                    {options.deliberation ? 'End Deliberation' : 'Start Deliberation'}
                </Button>
            </div>
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
