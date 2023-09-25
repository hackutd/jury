import { useEffect, useState } from 'react';
import { twMerge } from 'tailwind-merge';
import Paragraph from '../Paragraph';

interface ProjectDisplayProps {
    /* Project ID to display */
    projectId: string;

    /* Define the className */
    className?: string;
}

const ProjectDisplay = (props: ProjectDisplayProps) => {
    const [project, setProject] = useState<Project | null>(null);

    useEffect(() => {
        console.log(props.projectId)
        async function fetchData() {
            if (!props.projectId) return;
            
            const projectRes = await fetch(
                `${process.env.REACT_APP_JURY_URL}/project/${props.projectId}`,
                {
                    method: 'GET',
                    credentials: 'include',
                }
            );
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
    }, [props.projectId]);

    if (!project) {
        return <div className={twMerge('px-2', props.className)}>Loading...</div>;
    }

    return (
        <div className={twMerge('px-2', props.className)}>
            <h1 className='text-3xl mb-1'>{project.name}</h1>
            <h2 className='text-xl mb-1'>Table {project.location}</h2>
            <Paragraph className='text-light' text={project.description} />
        </div>
    );
};

export default ProjectDisplay;
