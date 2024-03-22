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
                errorAlert(loggedInRes);
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
                errorAlert(judgeRes);
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
            errorAlert(readWelcomeRes);
            return;
        }

        navigate('/judge');
    };

    if (!judge) return <Loading disabled={judge !== null} />;

    return (
        <>
            <JuryHeader withLogout />
            <Container noCenter className="mx-7">
                <h1 className="text-2xl my-2">Hello, {judge.name}!</h1>
                <h2 className="text-lg font-bold">PLEASE READ THE FOLLOWING:</h2>
                <p className="my-2">
                    Welcome to Jury, an innovative judging system that uses a dynamic ranking system
                    to facilitate hackathon judging.
                </p>
                <p className="my-2">
                    Once you get started, you will be presented with a project and its location.
                    Please go to that project and listen to their presentation. Once completed,
                    please score their project on the respective categories and click "Done".
                </p>
                <p className="my-2">
                    Once you have scored a project, you will be taken to the ranking screen. Here,
                    you can rank the project against others you have seen. You can also view the
                    projects you&apos;ve seen previously and adjust their scores.
                </p>
                <p className="my-2">
                    If a team is busy being judged, click the &apos;busy&apos; button. This will NOT
                    impact their rating and you may be presented with this team again.
                </p>
                <p className="my-2">
                    If a team is absent or you suspect a team may be cheating, please report it to
                    the organizers with the &apos;flag&apos; button. We will look into the matter
                    and take the proper action.
                </p>
                <Checkbox checked={checkRead} onChange={setCheckRead}>
                    Before you continue, please acknowledge that you have read and understand the
                    above instructions.
                </Checkbox>
                <Checkbox checked={checkEmail} onChange={setCheckEmail}>
                    I certify that my email is <span className="text-primary">[{judge.email}]</span>
                    . If this is not your email, contact an organizer immediately.
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
