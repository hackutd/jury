import { useNavigate, useParams } from 'react-router-dom';
import Container from '../components/Container';
import JuryHeader from '../components/JuryHeader';
import { useEffect, useState } from 'react';
import StarDisplay from '../components/judge/StarDisplay';
import Paragraph from '../components/Paragraph';
import Back from '../components/Back';

const Project = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [project, setProject] = useState<null | Project>(null);

    useEffect(() => {
        async function fetchData() {
            const projectRes = await fetch(`${process.env.REACT_APP_JURY_URL}/project/${id}`, {
                method: 'GET',
                credentials: 'include',
            });
            if (!projectRes.ok) {
                alert(
                    `Unable to get project - ${projectRes.status}: ${projectRes.statusText}. Please check your internet connection or refresh the page!`
                );
                return;
            }

            const newProject = await projectRes.json();
            setProject(newProject);
        }

        fetchData();
    }, []);

    if (!project) return <div>Loading...</div>;

    return (
        <>
            <JuryHeader withLogout />
            <Container noCenter={true} className="px-2">
                <Back location="/judge" />
                <h1 className="text-3xl mb-1">{project.name}</h1>
                <div className="flex mb-3">
                    <StarDisplay stars={project.stars} />
                </div>
                <Paragraph text={project.description} className="text-light" />
            </Container>
        </>
    );
};

export default Project;
