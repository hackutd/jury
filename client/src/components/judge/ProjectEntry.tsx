import { useNavigate } from 'react-router-dom';

interface ProjectEntryProps {
    project?: SortableJudgedProject;
}

const ProjectEntry = ({ project }: ProjectEntryProps) => {
    // const navigate = useNavigate();

    // const openProject = () => {
    //     navigate(`/judge/project/${project.project_id}`);
    // };

    if (!project) {
        return null;
    }

    if (project.id === -1) {
        return <div>
            <h3 className="text-xl grow mt-4">Unsorted Projects below here...</h3>
        </div>
    }

    return (
            <div className="m-1 p-2 drop-shadow-md bg-background rounded-md">
                <h3 className="text-xl grow">
                    <b>Table {project.location}</b>{': '}
                    {project.name}
                </h3>
                <p className="text-light">
                    {'| '}{Object.entries(project.categories).map(([name, score], i) => (
                        <span key={i}>{name.substring(0, 4) + ': ' + score + ' | '}</span>
                    ))}
                </p>
            </div>
    );
};

export default ProjectEntry;
