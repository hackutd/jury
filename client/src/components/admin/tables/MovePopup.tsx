import { useEffect, useState } from 'react';
import { postRequest, putRequest } from '../../../api';
import ConfirmPopup from '../../ConfirmPopup';
import { errorAlert } from '../../../util';
import { useAdminStore, useAdminTableStore } from '../../../store';
import RawTextInput from '../../RawTextInput';

interface MovePopupProps {
    /* State variable to open popup */
    open: boolean;

    /* Setter for open */
    setOpen: React.Dispatch<React.SetStateAction<boolean>>;

    /* Judge to move */
    judge?: Judge;

    /* List of judge IDs (at least one of judge/judges should be defined) */
    judges?: string[];
}

const MovePopup = (props: MovePopupProps) => {
    const [newGroup, setNewGroup] = useState('');
    const fetchJudges = useAdminStore((state) => state.fetchJudges);
    const setSelected = useAdminTableStore((state) => state.setSelected);
    const selected = useAdminTableStore((state) => state.selected);

    const handleSubmit = () => {
        if (!props.judge && !props.judges) {
            throw new Error('At least one of judge/judges must be defined');
        }

        // Check to see if new group is valid
        const group = Number(newGroup);
        if (isNaN(group) || group < 0) {
            alert('Invalid group number');
            return;
        }

        if (props.judge) {
            moveJudgeGroup(group);
        } else if (props.judges) {
            moveJudgesGroup(group);
        }
    };

    const moveJudgeGroup = async (group: number) => {
        const res = await putRequest<OkResponse>(`/judge/move/${props.judge?.id}`, 'admin', {
            group,
        });
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert(`Judge moved successfully!`);
        fetchJudges();
        props.setOpen(false);
    };

    const moveJudgesGroup = async (group: number) => {
        const res = await postRequest<OkResponse>('/judge/move', 'admin', {
            judges: props.judges,
            group,
        });
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert(`Selected judges moved successfully!`);
        fetchJudges();
        props.setOpen(false);
        setSelected(new Array(selected.length).fill(false));
    };

    useEffect(() => {
        if (!props.judge) return;

        setNewGroup(String(props.judge.group));
    }, [props.judge]);

    return (
        <ConfirmPopup
            enabled={props.open}
            setEnabled={props.setOpen}
            onSubmit={handleSubmit}
            submitText="Move"
            title="Move Judge"
        >
            <div className="flex flex-col items-center">
                {props.judge && (
                    <p className="text-light">
                        Move the judge{' '}
                        <span className="text-primary font-bold">{props.judge?.name}</span> to
                        another group. Enter the new group number below.
                    </p>
                )}
                {props.judges && (
                    <p className="text-light">
                        Move the selected judges to another group. Enter the new group number below.
                    </p>
                )}
                <RawTextInput
                    placeholder="New group"
                    text={newGroup}
                    setText={setNewGroup}
                    large
                    className="mt-2"
                />
            </div>
        </ConfirmPopup>
    );
};

export default MovePopup;
