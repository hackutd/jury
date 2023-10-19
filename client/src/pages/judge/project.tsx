import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';
import StarDisplay from '../../components/judge/StarDisplay';
import Paragraph from '../../components/Paragraph';
import Back from '../../components/Back';
import { getRequest } from '../../api';
import { errorAlert } from '../../util';

const Project = () => {
    const { id } = useParams();
    const [project, setProject] = useState<null | JudgedProject>(null);

    useEffect(() => {
        async function fetchData() {
            const projRes = await getRequest<JudgedProject>(`/judge/project/${id}`, 'judge');
            if (projRes.status !== 200) {
                errorAlert(projRes);
                return;
            }
            const proj = projRes.data as JudgedProject;
            setProject(proj);
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
                    <StarDisplay stars={project.stars} clickable id={project.project_id} />
                </div>
                <Paragraph text={project.description} className="text-light" />
            </Container>
        </>
    );
};

export default Project;
