import { useEffect, useState } from 'react';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';
import { useNavigate } from 'react-router-dom';
import Button from '../../components/Button';
import StatBlock from '../../components/StatBlock';
import Loading from '../../components/Loading';
import { getRequest, postRequest } from '../../api';
import { errorAlert } from '../../util';
import Ranking from '../../components/judge/dnd/Ranking';
import StarList from '../../components/judge/dnd/StarList';
import { Helmet } from 'react-helmet';

const Judge = () => {
    const navigate = useNavigate();
    const [judge, setJudge] = useState<Judge | null>(null);
    const [loaded, setLoaded] = useState(false);
    const [projCount, setProjCount] = useState(0);
    const [deliberation, setDeliberation] = useState(false);

    // Verify user is logged in and read welcome before proceeding
    useEffect(() => {
        async function fetchData() {
            // Check to see if the user is logged in
            const loggedInRes = await postRequest<OkResponse>('/judge/auth', 'judge', null);
            if (loggedInRes.status === 401) {
                console.error(`Judge is not logged in!`);
                navigate('/judge/login');
                return;
            }
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

            // Get the name & email of the user from the server
            const judgeRes = await getRequest<Judge>('/judge', 'judge');
            if (judgeRes.status !== 200) {
                errorAlert(judgeRes);
                return;
            }
            const judge: Judge = judgeRes.data as Judge;
            setJudge(judge);

            // Get the project count
            const projCountRes = await getRequest<ProjectCount>('/project/count', 'judge');
            if (projCountRes.status !== 200) {
                errorAlert(projCountRes);
                return;
            }
            setProjCount(projCountRes.data?.count as number);

            // Get deliberation status
            const delibRes = await getRequest<OkResponse>('/judge/deliberation', 'judge');
            if (delibRes.status !== 200) {
                errorAlert(delibRes);
                return;
            }
            setDeliberation(delibRes.data?.ok === 1);

            setLoaded(true);
        }

        fetchData();
    }, []);

    if (!loaded || !judge) return <Loading disabled={!loaded || !judge} />;

    return (
        <>
            <Helmet>
                <title>Judging | Jury</title>
            </Helmet>
            <JuryHeader withLogout />
            <Container noCenter className="px-2 pb-4">
                <h1 className="text-2xl mt-2">Welcome, {judge.name}!</h1>
                {judge.track !== '' && (
                    <p className="text-lighter font-bold">Track: {judge.track}</p>
                )}
                <div className="w-full mb-6 mt-2">
                    <Button type="primary" full href="/judge/live">
                        Next Project
                    </Button>
                </div>
                <div className="flex justify-evenly">
                    <StatBlock name="Seen" value={judge.seen_projects.length as number} />
                    <StatBlock name="Total Projects" value={projCount} />
                </div>
                {judge.track === '' ? (
                    <Ranking judge={judge} deliberation={deliberation} />
                ) : (
                    <StarList judge={judge} deliberation={deliberation} />
                )}
            </Container>
        </>
    );
};

export default Judge;
