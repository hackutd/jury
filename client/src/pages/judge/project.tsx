import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';
import Paragraph from '../../components/Paragraph';
import Back from '../../components/Back';
import { getRequest, putRequest } from '../../api';
import { errorAlert } from '../../util';
import Star from '../../components/judge/Star';
import TextArea from '../../components/TextArea';

const Project = () => {
    const { id } = useParams();
    const [project, setProject] = useState<null | JudgedProjectWithUrl>(null);
    const [notes, setNotes] = useState('');
    const [starred, setStarred] = useState(false);

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
            setStarred(proj.starred);
        }

        fetchData();
    }, []);

    // Update notes with a delay for typing
    useEffect(() => {
        if (!project) return;

        async function updateNotes() {
            const url = `/judge/notes/${project?.project_id}`;
            const res = await putRequest<OkResponse>(url, 'judge', {
                notes,
            });
            if (res.status !== 200) {
                errorAlert(res);
            }
        }
        const delayDebounceFn = setTimeout(updateNotes, 1000);

        return () => {
            updateNotes();
            clearTimeout(delayDebounceFn);
        };
    }, [notes]);

    const updateStar = async () => {
        const res = await putRequest<OkResponse>('/judge/star', 'judge', {
            project: project?.project_id,
            starred: !starred,
        });
        if (res.status !== 200) {
            errorAlert(res);
        }
    };

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
                <div className="flex flex-row justify-center text-left mb-2 px-2">
                    <Star
                        active={starred}
                        setActive={setStarred}
                        className="mr-4"
                        onClick={updateStar}
                    />
                    <p className="text-light">
                        Star projects you think should win the top places in the hackathon.
                    </p>
                </div>
                <TextArea
                    label="Personal Notes"
                    value={notes}
                    setValue={setNotes}
                    className="mb-4"
                />
                <Paragraph text={project.description} className="text-black" />
            </Container>
        </>
    );
};

export default Project;
