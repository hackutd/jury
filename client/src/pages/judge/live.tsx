import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { twMerge } from 'tailwind-merge';

import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';
import ProjectDisplay from '../../components/judge/ProjectDisplay';
import Button from '../../components/Button';
import RatePopup from '../../components/judge/popups/RatePopup';
import Back from '../../components/Back';
import JudgeInfoPage from '../../components/judge/JudgeInfoPage';
import InfoPopup from '../../components/InfoPopup';
import FlagPopup from '../../components/judge/popups/FlagPopup';

import { getRequest, postRequest } from '../../api';
import { errorAlert } from '../../util';
import alarm from '../../assets/alarm.mp3';
import data from '../../data.json';
import RawTextInput from '../../components/RawTextInput';

const infoPages = ['paused', 'hidden', 'no-projects', 'done'];
const infoData = [
    data.judgeInfo.paused,
    data.judgeInfo.hidden,
    data.judgeInfo.noProjects,
    data.judgeInfo.done,
];

const audio = new Audio(alarm);

const JudgeLive = () => {
    const navigate = useNavigate();
    const [verified, setVerified] = useState(false);
    const [judge, setJudge] = useState<Judge | null>(null);
    const [votePopup, setVotePopup] = useState<boolean>(false);
    const [flagPopup, setFlagPopup] = useState<boolean>(false);
    const [skipPopup, setSkipPopup] = useState<boolean>(false);
    const [infoPage, setInfoPage] = useState<string>('');
    const [started, setStarted] = useState(false);
    const [totalJudgingTime, setTotalJudgingTime] = useState(0);
    const [time, setTime] = useState(0);
    const [timerStart, setTimerStart] = useState(0);
    const [timerInterval, setTimerInterval] = useState<number | null>(null);
    const [timerDisplay, setTimerDisplay] = useState('');
    const [timesUp, setTimesUp] = useState(false);
    const [stopAudio, setStopAudio] = useState(false);
    const [audioPopupOpen, setAudioPopupOpen] = useState(false);
    const [paused, setPaused] = useState(false);
    const [notes, setNotes] = useState('');

    useEffect(() => {
        async function fetchData() {
            // Check to see if the user is logged in
            const loggedInRes = await postRequest<OkResponse>('/judge/auth', 'judge', null);
            if (loggedInRes.status !== 200) {
                errorAlert(loggedInRes);
                return;
            }
            if (loggedInRes.data?.ok !== 1) {
                console.error(`Judge is not logged in!`);
                navigate('/judge/login');
                return;
            }

            // Check for read welcome
            const readWelcomeRes = await getRequest<OkResponse>('/judge/welcome', 'judge');
            if (readWelcomeRes.status !== 200) {
                errorAlert(readWelcomeRes);
                return;
            }
            const readWelcome = readWelcomeRes.data?.ok === 1;
            if (!readWelcome) {
                navigate('/judge/welcome');
            }

            // Check to see if judging has started
            const startedRes = await getRequest<OkResponse>('/admin/started', '');
            if (startedRes.status !== 200) {
                errorAlert(startedRes);
                return;
            }
            if (startedRes.data?.ok !== 1) {
                setVerified(true);
                setInfoPage('paused');
                return;
            }

            setVerified(true);
        }

        fetchData();
    }, []);

    // Once verification finishes, get the judge's next project to judge, as well as the timer
    async function getJudgeData() {
        // Get judging timer
        const timerRes = await getRequest<Timer>('/admin/timer', 'judge');
        if (timerRes.status !== 200) {
            errorAlert(timerRes);
            return;
        }
        const judgingTime = timerRes.data?.judging_timer as number;
        setTotalJudgingTime(judgingTime);
        if (judgingTime === 0) setStarted(true);

        // Get the judge
        const judgeRes = await getRequest<Judge>('/judge', 'judge');
        if (judgeRes.status !== 200) {
            errorAlert(judgeRes);
            return;
        }
        const theJudge = judgeRes.data as Judge;

        // If the judge is disabled, redirect to the disabled page
        if (!theJudge.active) {
            setInfoPage('hidden');
            return;
        }

        // If the judge has a current project, use that
        if (theJudge.current) {
            setJudge(theJudge);
            return;
        }

        // Otherwise, query for the next project to judge
        const newProject = await postRequest<NextJudgeProject>('/judge/next', 'judge', null);
        if (newProject.status !== 200) {
            errorAlert(newProject);
            return;
        }

        // No project has been returned (all projects have been judged)
        if (!newProject.data?.project_id) {
            setInfoPage('done');
            return;
        }

        // Set the judge's current project
        theJudge.current = newProject.data.project_id;
        setJudge(theJudge);
    }

    // Get judge data after verification
    useEffect(() => {
        if (!verified) return;
        if (infoPage === 'paused') return;

        getJudgeData();
    }, [verified]);

    // Timer logic
    useEffect(() => {
        if (timerStart === 0 || time === 0 || totalJudgingTime === 0) {
            return;
        }

        const interval = setInterval(() => {
            const elapsed = Date.now() - timerStart;
            const remaining = time - elapsed;

            if (remaining <= 0) {
                clearInterval(timerInterval as number);
                setTimesUp(true);
                return;
            }

            const minutes = Math.floor(remaining / 60000);
            const seconds = Math.floor((remaining % 60000) / 1000);
            setTimerDisplay(
                `${minutes < 10 ? '0' : ''}${minutes}:${seconds < 10 ? '0' : ''}${seconds}`
            );
        }, 100);
        setTimerInterval(interval);
        setStarted(true);

        return () => {
            if (timerInterval) clearInterval(timerInterval);
        };
    }, [timerStart]);

    // When time is up, show popup and play sound
    useEffect(() => {
        if (!timesUp) return;

        // Show popup
        setAudioPopupOpen(true);

        // Play sound
        setStopAudio(false);
        audioLoop();
    }, [timesUp]);

    // When popup closed, stop audio
    useEffect(() => {
        if (!audioPopupOpen) noAudio();
    }, [audioPopupOpen]);

    // Make timer audio run in a loop
    const audioLoop = () => {
        audio.play();
        audio.addEventListener('ended', () => {
            audio.currentTime = 0;
            console.log('Audio is done!');
            if (stopAudio) return;
            audioLoop();
        });
    };

    // Stop audio
    const noAudio = () => {
        audio.pause();
        audio.currentTime = 0;
        setStopAudio(true);
    };

    // Pause the timer
    const pauseTimer = () => {
        clearInterval(timerInterval as number);

        // Calculate remaining time
        const elapsed = Date.now() - timerStart;
        const remaining = time - elapsed;
        setTime(remaining);

        setPaused(true);
    };

    const flagCallback = async (isVote?: boolean) => {
        setJudge(null);

        // Update notes if voting
        if (isVote) {
            const res = await postRequest<OkResponse>('/judge/notes', 'judge', {
                notes,
                project: judge?.current,
            });
            if (res.status !== 200) {
                errorAlert(res);
            }
        }

        resetTimer();
        getJudgeData();
    };

    // Display an error page if an error condition holds
    const infoIndex = infoPages.indexOf(infoPage);
    if (infoIndex !== -1) {
        return (
            <JudgeInfoPage
                title={infoData[infoIndex].title}
                description={infoData[infoIndex].description}
            />
        );
    }

    // Show loading screen if judge does not exist
    if (!judge) {
        return (
            <>
                <JuryHeader withLogout />
                <Container className="px-2">
                    <h1>Loading...</h1>
                </Container>
            </>
        );
    }

    // Open popup when clicked
    const openPopup = (pop: VotePopupState) => {
        // Pause the timer before opening popup
        pauseTimer();

        // Open specified popup
        switch (pop) {
            case 'vote':
                setVotePopup(true);
                break;
            case 'flag':
                setFlagPopup(true);
                break;
            case 'skip':
                setSkipPopup(true);
                break;
            default:
                alert('Invalid popup state!');
                break;
        }
    };

    // Reset the judging timer
    const resetTimer = () => {
        clearInterval(timerInterval as number);
        setTime(0);
        setTimerStart(0);
        setTimerInterval(null);
        setTimerDisplay('');
        setTimesUp(false);
        setStarted(false);
        setPaused(false);
    };

    // Start the judging timer
    const startJudging = async () => {
        // Start the timer client-side
        const minutes = Math.floor(totalJudgingTime / 60);
        const seconds = Math.floor(totalJudgingTime % 60);
        setTime(totalJudgingTime * 1000);
        setTimerStart(Date.now());
        setTimerDisplay(
            `${minutes < 10 ? '0' : ''}${minutes}:${seconds < 10 ? '0' : ''}${seconds}`
        );
    };

    return (
        <>
            <JuryHeader withLogout />
            <Container noCenter className="px-2 pb-4">
                <Back location="/judge" />
                <div className="p-2 pb-4">
                    {totalJudgingTime === 0 ? null : !started ? (
                        <Button
                            type="primary"
                            className="py-8 text-4xl md:text-5xl rounded-xl mb-4"
                            full
                            onClick={startJudging}
                        >
                            Start Judging
                        </Button>
                    ) : (
                        <div
                            className={twMerge(
                                'py-5 mb-4 text-6xl rounded-xl w-full border-primary border-4 border-solid text-center cursor-pointer ',
                                timesUp ? 'border-error bg-error/20' : '',
                                paused ? 'bg-lighter/20' : ''
                            )}
                            onClick={() => {
                                if (paused) {
                                    setTimerStart(Date.now());
                                    setPaused(false);
                                    return;
                                }
                                pauseTimer();
                            }}
                        >
                            {timerDisplay}
                        </div>
                    )}
                    <div className="flex items-center">
                        <Button
                            type="primary"
                            className="bg-error mr-2 py-1 text-xl rounded-xl basis-2/5 disabled:bg-backgroundDark hover:bg-errorDark"
                            disabled={judge === null}
                            onClick={() => {
                                openPopup('flag');
                            }}
                        >
                            Flag
                        </Button>
                        <Button
                            type="primary"
                            className="bg-gold mx-2 py-1 text-xl rounded-xl basis-2/5 text-black disabled:bg-backgroundDark disabled:text-lighter hover:bg-goldDark"
                            disabled={judge === null}
                            onClick={() => {
                                openPopup('skip');
                            }}
                        >
                            Skip
                        </Button>
                        <Button
                            type="primary"
                            className="ml-2 py-1 text-xl rounded-xl"
                            disabled={judge === null || !started}
                            onClick={() => {
                                openPopup('vote');
                            }}
                        >
                            Done
                        </Button>
                    </div>
                </div>
                {judge.current && <ProjectDisplay judge={judge} projectId={judge.current} />}
                {/* Dummy div for fixed text input */}
                <div className="w-full py-2 h-10"></div>
                <div className="fixed bottom-0 flex p-2 w-full left-0 bg-background">
                    <RawTextInput
                        name="notes"
                        placeholder="Personal notes..."
                        className="text-md h-10 px-2 w-auto grow"
                        text={notes}
                        setText={setNotes}
                    />
                </div>
                <RatePopup
                    enabled={votePopup}
                    setEnabled={setVotePopup}
                    judge={judge}
                    callback={flagCallback.bind(this, true)}
                />
                <FlagPopup enabled={flagPopup} setEnabled={setFlagPopup} onSubmit={flagCallback} />
                <FlagPopup
                    enabled={skipPopup}
                    setEnabled={setSkipPopup}
                    onSubmit={flagCallback}
                    isSkip
                />
                <InfoPopup
                    enabled={audioPopupOpen}
                    setEnabled={setAudioPopupOpen}
                    title="Time is Up!"
                    submitText="Close"
                >
                    Please remind the participant that they should be done with their presentation.
                </InfoPopup>
            </Container>
        </>
    );
};

export default JudgeLive;
