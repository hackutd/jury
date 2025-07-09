import { useEffect, useState } from 'react';
import { putRequest } from '../../../api';
import ConfirmPopup from '../../ConfirmPopup';
import { errorAlert } from '../../../util';
import { useAdminStore, useAdminTableStore } from '../../../store';
import TextInput from '../../TextInput';

interface MovePopupProps {
    /* State variable to open popup */
    enabled: boolean;

    /* Setter for open */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;

    /* Item to move */
    item: Project;
}

const MovePopup = (props: MovePopupProps) => {
    const [newLocation, setNewLocation] = useState('');
    const fetchProjects = useAdminStore((state) => state.fetchProjects);
    const setSelected = useAdminTableStore((state) => state.setSelected);
    const selected = useAdminTableStore((state) => state.selected);

    const handleSubmit = () => {
        if (!props.item) {
            throw new Error('A project must be defined');
        }

        // Check to see if new location is valid
        const location = Number(newLocation);
        if (isNaN(location) || location < 0) {
            alert('Invalid location number');
            return;
        }

        move(location);
    };

    const move = async (location: number) => {
        const res = await putRequest<OkResponse>(`/project/move/${props.item?.id}`, 'admin', {
            location,
        });
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert(`Project moved to table ${location} successfully!`);
        fetchProjects();
        props.setEnabled(false);
    };

    useEffect(() => {
        if (!props.item) return;

        setNewLocation(String(props.item.location));
    }, [props.item]);

    return (
        <ConfirmPopup
            enabled={props.enabled}
            setEnabled={props.setEnabled}
            onSubmit={handleSubmit}
            submitText="Move"
            title="Move Project"
        >
            <div className="flex flex-col items-center">
                <p className="text-light">
                    Move the project&nbsp;
                    <span className="text-primary font-bold">{props.item?.name}</span> to a new
                    table. Enter the table number below.
                </p>
                <TextInput
                    label="New location"
                    text={newLocation}
                    setText={setNewLocation}
                    large
                    className="mt-2"
                />
            </div>
        </ConfirmPopup>
    );
};

export default MovePopup;
