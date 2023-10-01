import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';
import ProjectDisplay from '../../components/judge/ProjectDisplay';
import Button from '../../components/Button';
import VotePopup from '../../components/judge/VotePopup';
import Back from '../../components/Back';

const JudgeLive = () => {
    const navigate = useNavigate();
    const [verified, setVerified] = useState(false);
    const [judge, setJudge] = useState<Judge | null>(null);
    const [popup, setPopup] = useState<boolean>(false);
    const [popupType, setPopupType] = useState<VotePopupState>('skip');

    useEffect(() => {
        async function fetchData() {
            // Check to see if the user is logged in
            const loggedIn = await fetch(`${process.env.REACT_APP_JURY_URL}/judge/auth`, {
                method: 'POST',
                credentials: 'include',
            });
            if (!loggedIn.ok) {
                console.error(`Judge is not logged in! ${loggedIn.status} ${loggedIn.statusText}`);
                navigate('/judge/login');
                return;
            }

            // Check for read welcome
            const readWelcomeRes = await fetch(`${process.env.REACT_APP_JURY_URL}/judge/welcome`, {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
            });
            if (!readWelcomeRes.ok) {
                alert(
                    `Unable to connect to server: ${readWelcomeRes.status} ${readWelcomeRes.statusText}. Please check your connection or reload the page.`
                );
                return;
            }
            const readWelcome = await readWelcomeRes.json();
            if (!readWelcome.ok) {
                navigate('/judge/welcome');
            }

            setVerified(true);
        }

        fetchData();
    }, []);

    // Once verification finishes, get the judge's next and prev project to judge
    async function getJudgeData() {
        const judgeRes = await fetch(`${process.env.REACT_APP_JURY_URL}/judge`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
        });
        if (!judgeRes.ok) {
            alert(
                `Unable to connect to server: ${judgeRes.status} ${judgeRes.statusText}. Please check your connection or reload the page.`
            );
            return;
        }
        const newJudge: Judge = await judgeRes.json();

        // If judge has literally no projects, query for a new project
        if (!newJudge.prev && !newJudge.next) {
            const judgeVoteRes = await fetch(`${process.env.REACT_APP_JURY_URL}/judge/vote`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({
                    judge_id: newJudge._id.$oid,
                    curr_winner: true,
                }),
            });

            const judgeVote: JudgeVoteRes = await judgeVoteRes.json();
            if (judgeVote.next_project_id) newJudge.next = { $oid: judgeVote.next_project_id };
            if (judgeVote.prev_project_id) newJudge.prev = { $oid: judgeVote.prev_project_id };
        } else if (!newJudge.next) {
            // If the judge only has a "prev" project, that means they've gone through all projects
            navigate("/judge/done");
        }

        console.log(newJudge);
        setJudge(newJudge);
    }

    useEffect(() => {
        if (!verified) return;

        getJudgeData();
    }, [verified]);

    const judgeVote = async (choice: number) => {
        const judgeId = judge?._id.$oid;
        setJudge(null);

        // Vote for the given choice
        const judgeVoteRes = await fetch(`${process.env.REACT_APP_JURY_URL}/judge/vote`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({
                judge_id: judgeId,
                curr_winner: choice === 0,
            }),
        });

        if (judgeVoteRes.status !== 200) {
            alert('Error! Unable to cast vote :(');
            return;
        }

        console.log(await judgeVoteRes.json());

        getJudgeData();
    };

    const flag = async (choice: number) => {
        alert('Flagged w choice ' + choice);
        return;
    };

    const skip = async (choice: number) => {
        alert('Skipped w choice ' + choice);
        return;
    };

    // Show loading screen
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
                {judge.next && <ProjectDisplay projectId={judge.next.$oid} />}
                {judge.prev && (
                    <>
                        <div className="my-6 h-[0.5px] shrink-0 w-full bg-light"></div>
                        <h3 className="font-bold text-light px-2 mb-1">Previous Project</h3>
                        <ProjectDisplay projectId={judge.prev.$oid} />
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
