import { useNavigate } from 'react-router-dom';
import StarDisplay from './StarDisplay';

interface ProjectEntryProps {
    id: string;
    name: string;
    description: string;
    stars: number;
}

const ProjectEntry = (props: ProjectEntryProps) => {
    const navigate = useNavigate();

    const openProject = () => {
        navigate(`/judge/project/${props.id}`);
    };

    return (
        <>
            <div className="cursor-pointer hover:bg-primary/10 duration-100">
                <div className="flex items-center justify-end">
                    <h3 className="text-xl grow" onClick={openProject}>
                        {props.name}
                    </h3>
                    <StarDisplay stars={props.stars} id={props.id} clickable />
                </div>
                <p className="text-light line-clamp-3" onClick={openProject}>
                    {props.description.replace('\\n', ' ')}
                </p>
            </div>
            <div className="h-[1px] w-full bg-light my-2"></div>
        </>
    );
};

export default ProjectEntry;
