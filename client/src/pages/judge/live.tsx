import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';
import ProjectDisplay from '../../components/judge/ProjectDisplay';
import Button from '../../components/Button';
import VotePopup from '../../components/judge/VotePopup';
import Back from '../../components/Back';
import { getRequest, postRequest } from '../../api';
import { errorAlert } from '../../util';
import data from '../../data.json';
import JudgeInfoPage from '../../components/judge/info';
import alarm from '../../assets/alarm.mp3';
import InfoPopup from '../../components/InfoPopup';
import { twMerge } from 'tailwind-merge';

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
    const [popup, setPopup] = useState<boolean>(false);
    const [infoPage, setInfoPage] = useState<string>('');
    const [popupType, setPopupType] = useState<VotePopupState>('skip');
    const [started, setStarted] = useState(false);
    const [time, setTime] = useState(0);
    const [timerStart, setTimerStart] = useState(0);
    const [timerInterval, setTimerInterval] = useState<NodeJS.Timeout | null>(null);
    const [timerDisplay, setTimerDisplay] = useState('');
    const [timesUp, setTimesUp] = useState(false);
    const [stopAudio, setStopAudio] = useState(false);
    const [audioPopupOpen, setAudioPopupOpen] = useState(false);
    const [paused, setPaused] = useState(false);

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

    // Once verification finishes, get the judge's next and prev project to judge
    async function getJudgeData() {
        const judgeRes = await getRequest<Judge>('/judge', 'judge');
        if (judgeRes.status !== 200) {
            errorAlert(judgeRes);
            return;
        }
        const newJudge = judgeRes.data as Judge;

        // If the judge is disabled, redirect to the disabled page
        if (!newJudge.active) {
            setInfoPage('hidden');
            return;
        }

        // If judge has literally no projects, query for IPO (Initial Project Offering)
        if (!newJudge.prev && !newJudge.next) {
            const ipoRes = await postRequest<JudgeIpo>('/judge/ipo', 'judge', null);
            if (ipoRes.status !== 200) {
                errorAlert(ipoRes);
                return;
            }
            if (!ipoRes.data?.initial) {
                setInfoPage('no-projects');
                return;
            }

            // No project has been returned (all projects have been judged)
            if (!ipoRes.data?.project_id) {
                setInfoPage('done');
                return;
            }

            newJudge.next = ipoRes.data?.project_id as string;
        } else if (!newJudge.next) {
            // If the judge only has a "prev" project, that means they've gone through all projects
            setInfoPage('done');
        }

        setJudge(newJudge);
    }

    useEffect(() => {
        if (!verified) return;
        if (infoPage === 'paused') return;
        
        getJudgeData();
    }, [verified]);

    useEffect(() => {
        if (timerStart === 0 || time === 0) {
            return;
        }

        const interval = setInterval(() => {
            const elapsed = Date.now() - timerStart;
            const remaining = time - elapsed;

            if (remaining <= 0) {
                clearInterval(timerInterval as NodeJS.Timeout);
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

    const audioLoop = () => {
        audio.play();
        audio.addEventListener('ended', () => {
            audio.currentTime = 0;
            console.log('Audio is done!');
            if (stopAudio) return;
            audioLoop();
        });
    };

    const noAudio = () => {
        audio.pause();
        audio.currentTime = 0;
        setStopAudio(true);
    };
    
    const pauseTimer = () => {
        clearInterval(timerInterval as NodeJS.Timeout);
        
        // Calculate remaining time
        const elapsed = Date.now() - timerStart;
        const remaining = time - elapsed;
        setTime(remaining);

        setPaused(true);
    }

    const judgeVote = async (choice: number) => {
        setJudge(null);

        // Vote for the given choice
        const voteRes = await postRequest<OkResponse>('/judge/vote', 'judge', {
            curr_winner: choice === 0,
        });
        if (voteRes.status !== 200) {
            errorAlert(voteRes);
            return;
        }

        resetTimer();
        getJudgeData();
    };

    const flag = async (choice: number) => {
        setJudge(null);

        const options = ['cannot-demo', 'too-complex', 'offensive'];

        // Flag the current project
        const flagRes = await postRequest<OkResponse>('/judge/skip', 'judge', {
            reason: options[choice],
        });
        if (flagRes.status !== 200) {
            errorAlert(flagRes);
            return;
        }

        resetTimer();
        getJudgeData();
    };

    const skip = async (choice: number) => {
        setJudge(null);

        const options = ['absent', 'busy'];

        // Skip the current project
        const skipRes = await postRequest<OkResponse>('/judge/skip', 'judge', {
            reason: options[choice],
        });
        if (skipRes.status !== 200) {
            errorAlert(skipRes);
            return;
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

    const openPopup = (pop: VotePopupState) => {
        if (pop === 'vote' && !judge.prev) {
            judgeVote(0);
            return;
        }
        setPopupType(pop);
        setPopup(true);
    };
    
    const resetTimer = () => {
        clearInterval(timerInterval as NodeJS.Timeout);
        setTime(0);
        setTimerStart(0);
        setTimerInterval(null);
        setTimerDisplay('');
        setTimesUp(false);
        setStarted(false);
        setPaused(false);
    }

    // Start the judging timer
    const startJudging = async () => {
        // Get judging timer
        const timerRes = await getRequest<Timer>('/admin/timer', 'judge');
        if (timerRes.status !== 200) {
            errorAlert(timerRes);
            return;
        }

        // Start the timer client-side
        const newTime = timerRes.data?.judging_timer as number;
        const minutes = Math.floor(newTime / 60);
        const seconds = Math.floor(newTime % 60);
        setTime(newTime * 1000);
        setTimerStart(Date.now());
        setTimerDisplay(
            `${minutes < 10 ? '0' : ''}${minutes}:${seconds < 10 ? '0' : ''}${seconds}`
        );
    };

    return (
        <>
            <JuryHeader withLogout />
            <Container noCenter className="px-2 pb-4">
                <div className="p-2">
                    {!started ? (
                        <Button
                            type="primary"
                            className="py-8 text-5xl rounded-xl"
                            full
                            onClick={startJudging}
                        >
                            Start Judging
                        </Button>
                    ) : (
                        <div
                            className={twMerge(
                                'py-5 text-6xl rounded-xl w-full border-primary border-4 border-solid text-center cursor-pointer ',
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
                    <div className="flex items-center mt-4">
                        <Button
                            type="primary"
                            className="bg-error mr-2 py-1 text-xl rounded-xl basis-2/5 disabled:bg-backgroundDark"
                            disabled={judge === null || !started}
                            onClick={() => {
                                openPopup('flag');
                            }}
                        >
                            Flag
                        </Button>
                        <Button
                            type="primary"
                            className="bg-gold mx-2 py-1 text-xl rounded-xl basis-2/5 text-black disabled:bg-backgroundDark disabled:text-lighter"
                            disabled={judge === null || !started}
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
                <Back location="/judge" />
                {judge.next && <ProjectDisplay judge={judge} projectId={judge.next} />}
                {judge.prev && (
                    <>
                        <div className="my-6 h-[0.5px] shrink-0 w-full bg-light"></div>
                        <h3 className="font-bold text-light px-2 mb-1">Previous Project</h3>
                        <ProjectDisplay judge={judge} projectId={judge.prev} />
                    </>
                )}
                <VotePopup
                    popupType={popupType}
                    open={popup}
                    close={setPopup}
                    flagFunc={flag}
                    judge={judge}
                    voteFunc={judgeVote}
                    skipFunc={skip}
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
