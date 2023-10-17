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

const infoPages = ['paused', 'hidden', 'no-projects', 'done'];
const infoData = [
    data.judgeInfo.paused,
    data.judgeInfo.hidden,
    data.judgeInfo.noProjects,
    data.judgeInfo.done,
];

const JudgeLive = () => {
    const navigate = useNavigate();
    const [verified, setVerified] = useState(false);
    const [judge, setJudge] = useState<Judge | null>(null);
    const [popup, setPopup] = useState<boolean>(false);
    const [infoPage, setInfoPage] = useState<string>('');
    const [popupType, setPopupType] = useState<VotePopupState>('skip');

    useEffect(() => {
        async function fetchData() {
            // Check to see if the user is logged in
            const loggedInRes = await postRequest<OkResponse>('/judge/auth', 'judge', null);
            if (loggedInRes.status !== 200) {
                errorAlert(loggedInRes.status);
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
                errorAlert(readWelcomeRes.status);
                return;
            }
            const readWelcome = readWelcomeRes.data?.ok === 1;
            if (!readWelcome) {
                navigate('/judge/welcome');
            }

            // Check to see if judging has started
            const startedRes = await getRequest<OkResponse>('/admin/started', '');
            if (startedRes.status !== 200) {
                errorAlert(startedRes.status);
                return;
            }
            if (startedRes.data?.ok !== 1) {
                console.error(`Judging is paused!`);
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
            errorAlert(judgeRes.status);
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
                errorAlert(ipoRes.status);
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
            navigate('/judge/done');
        }

        console.log(newJudge);
        setJudge(newJudge);
    }

    useEffect(() => {
        if (!verified) return;

        getJudgeData();
    }, [verified]);

    const judgeVote = async (choice: number) => {
        setJudge(null);

        // Vote for the given choice
        const voteRes = await postRequest<OkResponse>('/judge/vote', 'judge', {
            curr_winner: choice === 0,
        });
        if (voteRes.status !== 200) {
            errorAlert(voteRes.status);
            return;
        }

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
            errorAlert(flagRes.status);
            return;
        }

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
            errorAlert(skipRes.status);
            return;
        }

        getJudgeData();
    };

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

    return (
        <>
            <JuryHeader withLogout />
            <Container noCenter className="px-2 pb-4">
                <div className="flex items-center my-2 mx-1">
                    <Button
                        type="primary"
                        className="bg-error mx-2 py-1 text-xl rounded-xl basis-2/5"
                        disabled={judge === null}
                        onClick={() => {
                            openPopup('flag');
                        }}
                    >
                        Flag
                    </Button>
                    <Button
                        type="primary"
                        className="bg-gold mx-2 py-1 text-xl rounded-xl basis-2/5 text-black"
                        disabled={judge === null}
                        onClick={() => {
                            openPopup('skip');
                        }}
                    >
                        Skip
                    </Button>
                    <Button
                        type="primary"
                        className="mx-2 py-1 text-xl rounded-xl"
                        disabled={judge === null}
                        onClick={() => {
                            openPopup('vote');
                        }}
                    >
                        Done
                    </Button>
                </div>
                <Back location="/judge" />
                {judge.next && <ProjectDisplay projectId={judge.next} />}
                {judge.prev && (
                    <>
                        <div className="my-6 h-[0.5px] shrink-0 w-full bg-light"></div>
                        <h3 className="font-bold text-light px-2 mb-1">Previous Project</h3>
                        <ProjectDisplay projectId={judge.prev} />
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
            </Container>
        </>
    );
};

export default JudgeLive;
