import { useEffect, useState } from 'react';
import FlagPopupButton from './FlagPopupButton';
import Button from '../Button';
import { postRequest } from '../../api';
import { errorAlert } from '../../util';
import { useNavigate } from 'react-router-dom';

interface FlagPopupProps {
    /* Function to modify the popup state variable */
    close: React.Dispatch<React.SetStateAction<boolean>>;

    open: boolean;

    callback: () => void;
}

const flags = [
    {
        text: 'Absent',
        subtext: 'Team is not present for judging at their table',
    },
    {
        text: 'Cannot Demo Project',
        subtext: 'Team cannot prove they made the project',
    },
    {
        text: 'Too Complex',
        subtext: 'Appears too complex to be made in hackathon',
    },
    {
        text: 'Offensive Project',
        subtext: 'Offensive or breaks Code of Conduct',
    },
];

const FlagPopup = (props: FlagPopupProps) => {
    const navigate = useNavigate();
    const [selected, setSelected] = useState(-1);

    if (!props.open) return null;

    const handleClick = async () => {
        if (selected === -1) {
            alert('Please select an option.');
            return;
        }

        // Flag reasons
        const options = ['absent', 'cannot-demo', 'too-complex', 'offensive'];

        // Flag the current project
        const flagRes = await postRequest<OkResponse>('/judge/skip', 'judge', {
            reason: options[selected],
        });
        if (flagRes.status !== 200) {
            errorAlert(flagRes);
            return;
        }

        props.close(false);
        props.callback();
    };

    return (
        <>
            <div
                className="fixed left-0 top-0 z-20 w-screen h-screen bg-black/30"
                onClick={() => props.close(false)}
            ></div>
            <div className="bg-background fixed z-30 left-1/2 top-1/2 translate-x-[-50%] translate-y-[-50%] px-6 py-6 md:px-10 md:w-1/3 w-11/12 flex flex-col items-center">
                <h1 className="text-3xl font-bold text-error">Flag Project</h1>
                <h2 className="text-xl font-bold">Please select the reason</h2>
                {flags.map((v, i) => (
                    <FlagPopupButton
                        key={i}
                        onClick={() => {
                            setSelected(i);
                        }}
                        text={v.text}
                        subtext={v.subtext}
                        selected={selected === i}
                    />
                ))}
                <Button type="primary" onClick={handleClick} className='mt-4'>
                    Submit
                </Button>
            </div>
        </>
    );
};

export default FlagPopup;
