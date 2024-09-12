import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';
import Paragraph from '../../components/Paragraph';
import Back from '../../components/Back';
import { getRequest, postRequest } from '../../api';
import { errorAlert } from '../../util';
import Ratings from '../../components/judge/Ratings';
import RawTextInput from '../../components/RawTextInput';

const Project = () => {
    const { id } = useParams();
    const [project, setProject] = useState<null | JudgedProjectWithUrl>(null);
    const [notes, setNotes] = useState('');

    useEffect(() => {
        async function fetchData() {
            const projRes = await getRequest<JudgedProjectWithUrl>(`/judge/project/${id}`, 'judge');
            if (projRes.status !== 200) {
                errorAlert(projRes);
                return;
            }
            const proj = projRes.data as JudgedProjectWithUrl;
            setProject(proj);
            setNotes(proj.notes);
        }

        fetchData();
    }, []);

    useEffect(() => {
        if (!project) return;

        async function updateNotes() {
            const res = await postRequest<OkResponse>('/judge/notes', 'judge', {
                notes,
                project: project?.project_id,
            });
            if (res.status !== 200) {
                errorAlert(res);
            }
        }
        const delayDebounceFn = setTimeout(updateNotes, 1000);

        return () => {
            updateNotes();
            clearTimeout(delayDebounceFn);
        }
    }, [notes]);

    if (!project) return <div>Loading...</div>;

    return (
        <>
            <JuryHeader withLogout />
            <Container noCenter={true} className="px-2">
                <Back location="/judge" />
                <h1 className="text-3xl mb-1 font-bold">
                    <a
                        href={project.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="hover:text-primary duration-200"
                    >
                        {project.name}
                    </a>
                </h1>
                <h2 className="text-xl font-bold text-light mb-2">Table {project.location}</h2>
                <Ratings
                    prior={project.categories}
                    project={project}
                    small
                    submitText="Update"
                    update
                />
                <RawTextInput
                    name="notes"
                    placeholder="Your notes..."
                    text={notes}
                    setText={setNotes}
                    className="text-md h-10 px-2 w-auto mb-4"
                />
                <Paragraph text={project.description} className="text-light" />
            </Container>
        </>
    );
};

export default Project;
