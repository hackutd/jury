import { useState } from 'react';
import Button from '../../Button';
import Popup from '../../Popup';
import RadioSelect from '../../RadioSelect';
import { postRequest } from '../../../api';
import { errorAlert } from '../../../util';

interface FlagPopupProps {
    /* Enabled state variable */
    enabled: boolean;

    /* True if this is for the skip popup */
    isSkip?: boolean;

    /* Function to modify the popup state variable */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;

    /* Callback function to run when flag submitted */
    onSubmit: () => void;
}

const flags = [
    {
        value: 'cannot-demo',
        title: 'Cannot Demo Project',
        subtitle: 'Team cannot prove they made the project',
    },
    {
        value: 'too-complex',
        title: 'Too Complex',
        subtitle: 'Appears too complex to be made at hackathon',
    },
    {
        value: 'offensive',
        title: 'Offensive Project',
        subtitle: 'Offensive or breaks Code of Conduct',
    },
];

const skips = [
    {
        value: 'absent',
        title: 'Absent',
        subtitle: 'Team is not present for judging at their table',
    },
    {
        value: 'busy',
        title: 'Busy',
        subtitle: 'Team is busy with another judge',
    },
];

/**
 * Popup when user clicks on "Flag".
 */
const FlagPopup = (props: FlagPopupProps) => {
    const [selected, setSelected] = useState('');

    if (!props.enabled) return null;

    // On submit click, flags project and runs callback
    const handleClick = async () => {
        if (selected === '') {
            alert('Please select an option.');
            return;
        }

        // Flag the current project
        const flagRes = await postRequest<OkResponse>('/judge/skip', 'judge', {
            reason: selected,
        });
        if (flagRes.status !== 200) {
            errorAlert(flagRes);
            return;
        }

        props.setEnabled(false);
        props.onSubmit();
    };

    // Set color based on skip or flag popup
    const color = props.isSkip ? 'gold' : 'error';

    return (
        <Popup enabled={props.enabled} setEnabled={props.setEnabled} className="items-center">
            <h1 className={`text-3xl font-bold text-${color}`}>
                {props.isSkip ? 'Skip' : 'Flag'} Project
            </h1>
            <h2 className="text-xl font-bold">Please select the reason</h2>
            <RadioSelect
                color={color}
                options={props.isSkip ? skips : flags}
                selected={selected}
                setSelected={setSelected}
            />
            <Button type="primary" onClick={handleClick} className="mt-4">
                Submit
            </Button>
            <div className="hidden text-gold bg-gold/10 border-gold text-error bg-error/10 border-error">
                Dummy div to get the tailwind colors loaded
            </div>
        </Popup>
    );
};

export default FlagPopup;
