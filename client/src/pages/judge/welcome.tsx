import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';
import Checkbox from '../../components/Checkbox';
import Button from '../../components/Button';
import { getRequest, postRequest } from '../../api';
import { errorAlert } from '../../util';
import Loading from '../../components/Loading';

const JudgeWelcome = () => {
    const navigate = useNavigate();
    const [judge, setJudge] = useState<Judge | null>(null);
    const [checkRead, setCheckRead] = useState(false);
    const [checkEmail, setCheckEmail] = useState(false);

    // Verify user is logged in and read welcome before proceeding
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
            
            // Get the name & email of the user from the server
            const judgeRes = await getRequest<Judge>('/judge', 'judge');
            if (judgeRes.status !== 200) {
                errorAlert(judgeRes.status);
                return;
            }
            setJudge(judgeRes.data as Judge);
        }

        fetchData();
    }, []);

    // Read the welcome message and mark that the user has read it
    const readWelcome = async () => {
        if (!checkRead || !checkEmail) {
            alert(
                'Please read the welcome message and confirm by checking the boxes below before proceeding.'
            );
            return;
        }

        // POST to server to mark that the user has read the welcome message
        const readWelcomeRes = await postRequest<OkResponse>('/judge/welcome', 'judge', null);
        if (readWelcomeRes.status !== 200) {
            errorAlert(readWelcomeRes.status);
            return;
        }

        navigate('/judge');
    };
    
    if (!judge) return <Loading disabled={judge !== null} />;

    return (
        <>
            <JuryHeader withLogout />
            <Container noCenter className="mx-7">
                <h1 className="text-2xl my-2">
                    Hello, {judge.name}!
                </h1>
                <h2 className="text-lg font-bold">PLEASE READ THE FOLLOWING:</h2>
                <p className="my-2">
                    Welcome to Jury, an innovative judging system that utilizes pairwise
                    comparisons for hackathon judging.
                </p>
                <p className="my-2">
                    Once you get started, you will be presented with a project and its location.
                    After you complete the first project, you will be presented with a second.
                </p>
                <p className="my-2">
                    Once completing the second project, you will be asked to compare the current
                    project and the previous project you viewed. This same process will continue for
                    all projects afterwards until all projects have been viewed.
                </p>
                <p className="my-2">
                    If a team is absent or is busy being judged, click the &apos;skip&apos; button.
                    This will NOT impact their rating!
                </p>
                <p className="my-2">
                    If you suspect a team may be cheating, please report it to the organizers with
                    the &apos;flag&apos; button.
                </p>
                <Checkbox checked={checkRead} onChange={setCheckRead}>
                    Before you continue, please acknowledge that you have read and understand the
                    above instructions.
                </Checkbox>
                <Checkbox checked={checkEmail} onChange={setCheckEmail}>
                    I certify that my email is{' '}
                    <span className="text-primary">[{judge.email}]</span>. If this is not your
                    email, contact an organizer immediately.
                </Checkbox>
                <div className="flex justify-center py-4">
                    <Button
                        type="primary"
                        disabled={!checkRead || !checkEmail}
                        onClick={readWelcome}
                        className="my-2"
                    >
                        Continue
                    </Button>
                </div>
            </Container>
        </>
    );
};

export default JudgeWelcome;
