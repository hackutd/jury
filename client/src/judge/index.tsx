import { useEffect, useState } from 'react';
import Container from '../components/Container';
import JuryHeader from '../components/JuryHeader';
import { useNavigate } from 'react-router-dom';
import Button from '../components/Button';
import AdminStat from '../components/admin/AdminStat';
import ProjectEntry from '../components/judge/ProjectEntry';
import Loading from '../components/Loading';

const Judge = () => {
    const navigate = useNavigate();
    const [name, setName] = useState('');
    const [projects, setProjects] = useState<Project[]>([]);
    const [loaded, setLoaded] = useState(false);

    // Verify user is logged in and read welcome before proceeding
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

            // Get the name & email of the user from the server
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
            const judge: Judge = await judgeRes.json();

            setName(judge.name);
        }

        fetchData();
    }, []);

    useEffect(() => {
        if (!name) return;

        async function getProjects() {
            const projRes = await fetch(`${process.env.REACT_APP_JURY_URL}/judge/projects`, {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
            });
            if (!projRes.ok) {
                alert(
                    `Unable to connect to server: ${projRes.status} ${projRes.statusText}. Please check your connection or reload the page.`
                );
            }
            const newProjects = await projRes.json();
            setProjects(newProjects);
            setLoaded(true);
        }

        getProjects();
    }, [name]);

    return (
        <>
            <JuryHeader withLogout />
            <Container noCenter className="px-2">
                <h1 className="text-2xl my-2">Welcome, {name}!</h1>
                <div className="w-full mb-6">
                    <Button type="primary" full square href="/judge/live">
                        Start Judging
                    </Button>
                </div>
                <div className="flex justify-evenly">
                    <AdminStat name="Seen" value={0} />
                    <AdminStat name="Projects" value={0} />
                </div>
                <h2 className="text-primary text-xl font-bold mt-4">Viewed Projects</h2>
                <div className="h-[1px] w-full bg-light my-2"></div>
                {projects.map((p) => (
                    <ProjectEntry
                        id={p._id.$oid}
                        name={p.name}
                        description={p.description}
                        stars={0}
                        key={p._id.$oid}
                    />
                ))}
                <Loading disabled={loaded} />
            </Container>
        </>
    );
};

export default Judge;
