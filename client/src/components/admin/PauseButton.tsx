import { useEffect, useState } from 'react';
import Button from '../Button';
import { postRequest } from '../../api';
import { errorAlert } from '../../util';

interface PauseButtonProps {
    paused: boolean;
    setPaused: React.Dispatch<React.SetStateAction<boolean>>;
    clock: number;
}

const pauseButtonFormat =
    'absolute top-6 left-6 w-40 md:w-52 text-lg py-2 px-1 hover:scale-100 focus:scale-100 rounded-md';

const PauseButton = ({ paused, setPaused, clock }: PauseButtonProps) => {
    const [resumeText, setResumeText] = useState('Resume');

    useEffect(() => {
        setResumeText(clock === 0 ? 'Start' : 'Resume');
    }, [clock]);

    const handleClick = async (pause: boolean) => {
        // Send pause/unpause request to server
        const res = await postRequest(`/admin/clock/${pause ? 'pause' : 'unpause'}`, 'admin', null);
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }
        
        setPaused(pause);
    };

    return paused ? (
        <Button
            type="primary"
            square
            bold
            className={pauseButtonFormat}
            onClick={() => handleClick(false)}
        >
            {resumeText} Judging
        </Button>
    ) : (
        <Button
            type="outline"
            square
            bold
            className={pauseButtonFormat}
            onClick={() => handleClick(true)}
        >
            Pause Judging
        </Button>
    );
};

export default PauseButton;
