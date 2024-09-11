import { useNavigate } from 'react-router-dom';
import DragHamburger from './dnd/DragHamburger';

interface ProjectEntryProps {
    project?: SortableJudgedProject;
}

const ProjectEntry = ({ project }: ProjectEntryProps) => {
    if (!project) {
        return null;
    }

    const navigate = useNavigate();

    const openProject = () => {
        navigate(`/judge/project/${project.project_id}`);
    };

    if (project.id === -1) {
        return (
            <div>
                <h3 className="text-xl grow mt-4">Unsorted Projects below here...</h3>
            </div>
        );
    }

    return (
        <div className="m-1 p-2 drop-shadow-md bg-background rounded-md">
            <div className="flex flex-row">
                <div>
                    <h3 className="text-xl grow">
                        <b>Table {project.location}</b>
                        {': '}
                        {project.name}
                    </h3>
                    <p className="text-light">
                        {'| '}
                        {Object.entries(project.categories).map(([name, score], i) => (
                            <span key={i}>{name.substring(0, 4) + ': ' + score + ' | '}</span>
                        ))}
                    </p>
                </div>
                <div className="grow text-right flex items-center justify-end">
                    <DragHamburger />
                    {/* <button onClick={openProject} className="text-3xl w-10 h-10 font-bold p-2 text-light duration-200 hover:text-primary leading-[0.5] rounded-full">
                        +
                    </button> */}
                </div>
            </div>
        </div>
    );
};

export default ProjectEntry;
