import { useEffect, useState } from 'react';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';
import { useNavigate } from 'react-router-dom';
import Button from '../../components/Button';
import StatBlock from '../../components/StatBlock';
import ProjectEntry from '../../components/judge/ProjectEntry';
import Loading from '../../components/Loading';
import { getRequest, postRequest } from '../../api';
import { errorAlert } from '../../util';

const Judge = () => {
    const navigate = useNavigate();
    const [judge, setJudge] = useState<Judge | null>(null);
    const [projects, setProjects] = useState<JudgedProject[]>([]);
    const [loaded, setLoaded] = useState(false);
    const [projCount, setProjCount] = useState(0);

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
        }

        fetchData();
    }, []);

    useEffect(() => {
        if (!judge) return;

        async function getProjects() {
            const projRes = await getRequest<JudgedProject[]>('/judge/projects', 'judge');
            if (projRes.status !== 200) {
                errorAlert(projRes);
                return;
            }
            const newProjects = projRes.data as JudgedProject[];
            setProjects(newProjects);
            setLoaded(true);
        }

        getProjects();
    }, [judge]);

    if (!loaded) return <Loading disabled={!loaded} />;

    const takeBreak = async () => {
        const res = await postRequest<OkResponse>('/judge/break', 'judge', null);
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert('You can now take a break! Press "Next project" to continue judging.');
    };

    return (
        <>
            <JuryHeader withLogout />
            <Container noCenter className="px-2">
                <h1 className="text-2xl my-2">Welcome, {judge?.name}!</h1>
                <div className="w-full mb-6">
                    <Button type="primary" full square href="/judge/live">
                        Next Project
                    </Button>
                    <div className="flex align-center justify-center mt-4">
                        <Button type="outline" square onClick={takeBreak} className="text-lg p-2">
                            I want to take a break!
                        </Button>
                    </div>
                </div>
                <div className="flex justify-evenly">
                    <StatBlock name="Seen" value={judge?.seen_projects.length as number} />
                    <StatBlock name="Total Projects" value={projCount} />
                </div>
                <h2 className="text-primary text-xl font-bold mt-4">Viewed Projects</h2>
                <div className="h-[1px] w-full bg-light my-2"></div>
                {projects.map((p) => (
                    <ProjectEntry
                        id={p.project_id}
                        name={p.name}
                        description={p.description}
                        stars={5}
                        key={p.project_id}
                    />
                ))}
            </Container>
        </>
    );
};

export default Judge;
