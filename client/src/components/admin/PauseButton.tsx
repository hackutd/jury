import { useEffect, useState } from 'react';
import Button from '../Button';
import { postRequest } from '../../api';
import { errorAlert } from '../../util';
import { useClockStore } from '../../store';

const pauseButtonFormat =
    'absolute top-6 left-6 w-40 md:w-52 text-lg py-2 px-1 hover:scale-100 focus:scale-100 rounded-md';

const PauseButton = () => {
    const clock = useClockStore((state) => state.clock);
    const fetchClock = useClockStore((state) => state.fetchClock);
    const [resumeText, setResumeText] = useState('Resume');

    useEffect(() => {
        setResumeText(clock.time === 0 ? 'Start' : 'Resume');
    }, [clock.time]);

    const handleClick = async (pause: boolean) => {
        // Send pause/unpause request to server
        const res = await postRequest(`/admin/clock/${pause ? 'pause' : 'unpause'}`, 'admin', null);
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }
        
        // Fetch new state of clock
        await fetchClock();
    };

    return clock.running ? (
        <Button
            type="outline"
            square
            bold
            className={pauseButtonFormat}
            onClick={() => handleClick(true)}
        >
            Pause Judging
        </Button>
    ) : (
        <Button
            type="primary"
            square
            bold
            className={pauseButtonFormat + " border-[2.5px] border-transparent border-solid"}
            onClick={() => handleClick(false)}
        >
            {resumeText} Judging
        </Button>
    );
};

export default PauseButton;
