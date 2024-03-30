import { useEffect, useState } from 'react';
import { twMerge } from 'tailwind-merge';
import Paragraph from '../Paragraph';
import { getRequest } from '../../api';
import { errorAlert } from '../../util';

interface ProjectDisplayProps {
    /* Project ID to display */
    projectId: string;

    /* Judge for the project */
    judge: Judge;

    /* Define the className */
    className?: string;
}

const ProjectDisplay = (props: ProjectDisplayProps) => {
    const [project, setProject] = useState<Project | null>(null);

    useEffect(() => {
        async function fetchData() {
            if (!props.projectId) return;

            const projRes = await getRequest<Project>(`/project/${props.projectId}`, 'judge');
            if (projRes.status !== 200) {
                errorAlert(projRes);
                return;
            }

            const newProject = projRes.data as Project;
            setProject(newProject);
        }

        fetchData();
    }, [props.projectId]);

    if (!project) {
        return <div className={twMerge('px-2', props.className)}>Loading...</div>;
    }

    return (
        <div className={twMerge('px-2', props.className)}>
            <h1 className="text-3xl mb-1">{project.name}</h1>
            <h2 className="text-xl mb-1">Table {project.location}</h2>
            <Paragraph className="text-light" text={project.description} />
        </div>
    );
};

export default ProjectDisplay;
