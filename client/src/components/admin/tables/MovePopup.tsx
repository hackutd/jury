import { useEffect, useState } from 'react';
import { postRequest, putRequest } from '../../../api';
import ConfirmPopup from '../../ConfirmPopup';
import { errorAlert } from '../../../util';
import { useAdminStore, useAdminTableStore } from '../../../store';
import TextInput from '../../TextInput';

interface MovePopupProps {
    /* State variable to open popup */
    enabled: boolean;

    /* Setter for open */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;

    /* Set to true if operating on project(s) */
    isProject?: boolean;

    /* Item to move */
    item?: Judge | Project;

    /* List of item IDs (at least one of judge/judges should be defined) */
    items?: string[];
}

const MovePopup = (props: MovePopupProps) => {
    const [newGroup, setNewGroup] = useState('');
    const fetchJudges = useAdminStore((state) => state.fetchJudges);
    const fetchProjects = useAdminStore((state) => state.fetchProjects);
    const setSelected = useAdminTableStore((state) => state.setSelected);
    const selected = useAdminTableStore((state) => state.selected);

    const handleSubmit = () => {
        if (!props.item && !props.items) {
            throw new Error('At least one of judge/judges must be defined');
        }

        // Check to see if new group is valid
        const group = Number(newGroup);
        if (isNaN(group) || group < 0) {
            alert('Invalid group number');
            return;
        }

        if (props.item) {
            moveGroup(group);
        } else if (props.items) {
            moveGroupMulti(group);
        }
    };

    const moveGroup = async (group: number) => {
        const res = await putRequest<OkResponse>(
            `/${props.isProject ? 'project' : 'judge'}/move/${props.item?.id}`,
            'admin',
            {
                group,
            }
        );
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert(`${props.isProject ? 'Project' : 'Judge'} moved successfully!`);
        props.isProject ? fetchProjects() : fetchJudges();
        props.setEnabled(false);
    };

    const moveGroupMulti = async (group: number) => {
        const res = await postRequest<OkResponse>(
            `/${props.isProject ? 'project' : 'judge'}/move`,
            'admin',
            {
                items: props.items,
                group,
            }
        );
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert(`Selected ${props.isProject ? 'projects' : 'judges'} moved successfully!`);
        props.isProject ? fetchProjects() : fetchJudges();
        props.setEnabled(false);
        setSelected(new Array(selected.length).fill(false));
    };

    useEffect(() => {
        if (!props.item) return;

        setNewGroup(String(props.item.group));
    }, [props.item]);

    return (
        <ConfirmPopup
            enabled={props.enabled}
            setEnabled={props.setEnabled}
            onSubmit={handleSubmit}
            submitText="Move"
            title={`Move ${props.isProject ? 'Project' : 'Judge'}`}
        >
            <div className="flex flex-col items-center">
                {props.item && (
                    <p className="text-light">
                        Move the {props.isProject ? 'project' : 'judge'}&nbsp;
                        <span className="text-primary font-bold">{props.item?.name}</span> to
                        another group. Enter the new group number below.
                    </p>
                )}
                {props.items && (
                    <p className="text-light">
                        Move the selected {props.isProject ? 'projects' : 'judges'} to another
                        group. Enter the new group number below.
                    </p>
                )}
                <TextInput
                    label="New group"
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
