import { useState } from 'react';
import { postRequest } from '../../api';
import { useAdminStore, useClockStore, useOptionsStore } from '../../store';
import { errorAlert } from '../../util';
import Button from '../Button';
import ConfirmPopup from '../ConfirmPopup';
import Popup from '../Popup';
import Paragraph from '../Paragraph';
import { useNavigate } from 'react-router-dom';

interface ActionsPopupProps {
    /* State variable for open/closed */
    enabled: boolean;

    /* Function to modify the popup state variable */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;
}

const ActionsItem = ({ children }: { children: React.ReactNode }) => {
    return <div className="flex flex-col md:flex-row mb-4 md:mb-0 items-center gap-3">{children}</div>;
};

const ActionsPopup = (props: ActionsPopupProps) => {
    const navigate = useNavigate();
    const options = useOptionsStore((state) => state.options);
    const fetchOptions = useOptionsStore((state) => state.fetchOptions);
    const fetchJudges = useAdminStore((state) => state.fetchJudges);
    const clock = useClockStore((state) => state.clock);
    const [deliberationPopup, setDeliberationPopup] = useState(false);
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
        <>
            <Popup enabled={props.enabled} setEnabled={props.setEnabled}>
                <h1 className="text-5xl text-center font-bold mb-6">Actions</h1>
                <div className="flex flex-col w-full gap-3">
                    <ActionsItem>
                        <Button
                            type="primary"
                            onClick={() => navigate('/admin/log')}
                            bold
                            className="py-2 basis-1/2 grow-0 shrink-0"
                        >
                            Audit Log
                        </Button>
                        <Paragraph
                            text="View all changes that have happened in Jury."
                            className="basis-1/2 grow-0 shrink-0"
                        />
                    </ActionsItem>
                    {options && options.multi_group && options.switching_mode === 'manual' && (
                        <ActionsItem>
                            <div className="flex flex-col items-center basis-1/2 grow-0 shrink-0">
                                <Button
                                    type="gold"
                                    onClick={setSwapPopup.bind(null, true)}
                                    disabled={clock.running}
                                    tooltip="Groups can only be swapped when judging is paused"
                                    bold
                                    className="py-2"
                                >
                                    Swap Judge Groups
                                </Button>
                                <p className="text-lg text-right grow mr-4 text-lighter">
                                    <span className="font-bold">Swaps: </span>
                                    {options.manual_switches}
                                </p>
                            </div>

                            <Paragraph
                                text="Prevent judges from making changes to ranking and stars."
                                className="basis-1/2 grow-0 shrink-0"
                            />
                        </ActionsItem>
                    )}
                    <ActionsItem>
                        <Button
                            type={options.deliberation ? 'outline' : 'error'}
                            onClick={setDeliberationPopup.bind(null, true)}
                            bold
                            tooltip="Stop judges from editing their rankings during deliberation"
                            className="py-2 basis-1/2 grow-0 shrink-0"
                        >
                            {options.deliberation ? 'End Deliberation' : 'Start Deliberation'}
                        </Button>
                        <Paragraph
                            text="Prevent judges from making changes to ranking and stars."
                            className="basis-1/2 grow-0 shrink-0"
                        />
                    </ActionsItem>
                </div>
            </Popup>
            <ConfirmPopup
                enabled={deliberationPopup}
                setEnabled={setDeliberationPopup}
                onSubmit={startDeliberation}
                submitText="Confirm"
                title={options.deliberation ? 'End Deliberation' : 'Start Deliberation'}
                red
            >
                Are you sure you want to {options.deliberation ? 'end' : 'start'} deliberation? This
                will impact whether or not judges will be able to rank/star projects!
            </ConfirmPopup>
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

export default ActionsPopup;
