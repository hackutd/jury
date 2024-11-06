import { twMerge } from 'tailwind-merge';
import DragHamburger from './dnd/DragHamburger';
import Star from './Star';
import { useEffect, useState } from 'react';
import { putRequest } from '../../api';
import { errorAlert } from '../../util';

interface ProjectEntryProps {
    project: SortableJudgedProject;
    ranking: number;
}

const ProjectEntry = ({ project, ranking }: ProjectEntryProps) => {
    const [starred, setStarred] = useState(project.starred);

    useEffect(() => {
        if (!project) return;

        setStarred(project.starred);
    }, [project]);

    if (!project) {
        return null;
    }

    let rankColor = 'text-lightest';
    switch (ranking) {
        case 1:
            rankColor = 'text-gold';
            break;
        case 2:
            rankColor = 'text-light';
            break;
        case 3:
            rankColor = 'text-bronze';
            break;
        default:
            break;
    }

    const updateStar = async () => {
        const res = await putRequest<OkResponse>(`/judge/star/${project.project_id}`, 'judge', {
            starred: !starred,
        });
        if (res.status !== 200) {
            errorAlert(res);
        }
    };

    return (
        <div className="flex items-center cursor-default">
            {ranking !== -1 && (
                <p className={twMerge('font-bold text-xl text-center w-8 shrink-0', rankColor)}>
                    {ranking}
                </p>
            )}
            <div className="m-1 pl-2 py-1 bg-background border-solid border-2 border-lightest rounded-md grow">
                <div className="flex flex-row">
                    <div className="grow">
                        <h3 className="text-lg leading-tight grow">
                            <a href={`/judge/project/${project.project_id}`}>
                                <b>Table {project.location}</b>
                                {': '}
                                {project.name}
                            </a>
                        </h3>
                        <p className="text-light text-xs line-clamp-1">{project.notes}</p>
                    </div>
                    <Star small active={starred} setActive={setStarred} onClick={updateStar} />
                    <div className="text-right flex items-center justify-end">
                        <DragHamburger />
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ProjectEntry;
