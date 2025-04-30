import { useEffect, useState } from 'react';
import Button from '../Button';
import { postRequest } from '../../api';
import { errorAlert } from '../../util';
import { useClockStore } from '../../store';

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
            type="outline-primary"
            bold
            className="w-56 ml-4 py-2 shrink-0"
            onClick={() => handleClick(true)}
        >
            Pause Judging
        </Button>
    ) : (
        <Button
            type="primary"
            bold
            className="w-56 ml-4 py-2 shrink-0"
            onClick={() => handleClick(false)}
        >
            {resumeText} Judging
        </Button>
    );
};

export default PauseButton;
