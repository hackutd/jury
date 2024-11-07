import { twMerge } from 'tailwind-merge';
import DragHamburger from './dnd/DragHamburger';
import Star from './Star';
import { useEffect, useState } from 'react';
import { putRequest } from '../../api';
import { errorAlert } from '../../util';

interface ProjectEntryProps {
    project: SortableJudgedProject;
    ranking: number;
    id: number;
    starCallback?: (id: number, starred: boolean) => void;
    noDrag?: boolean;
}

const ProjectEntry = (props: ProjectEntryProps) => {
    const [starred, setStarred] = useState(props.project.starred);

    useEffect(() => {
        if (!props.project) return;

        setStarred(props.project.starred);
    }, [props.project]);

    if (!props.project) {
        return null;
    }

    let rankColor = 'text-lightest';
    switch (props.ranking) {
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
        const res = await putRequest<OkResponse>(
            `/judge/star/${props.project.project_id}`,
            'judge',
            {
                starred: !starred,
            }
        );
        if (res.status !== 200) {
            errorAlert(res);
        }

        if (props.starCallback) {
            props.starCallback(props.id, !starred);
        }
    };

    return (
        <div className="flex items-center cursor-default">
            {props.ranking !== -1 && (
                <p className={twMerge('font-bold text-xl text-center w-8 shrink-0', rankColor)}>
                    {props.ranking}
                </p>
            )}
            <div className="m-1 pl-2 py-1 bg-background border-solid border-2 border-lightest rounded-md grow">
                <div className="flex flex-row">
                    <div className="grow">
                        <h3 className="text-lg leading-tight grow">
                            <a href={`/judge/project/${props.project.project_id}`}>
                                <b>Table {props.project.location}</b>
                                {': '}
                                {props.project.name}
                            </a>
                        </h3>
                        <p className="text-light text-xs line-clamp-1">{props.project.notes}</p>
                    </div>
                    <Star small active={starred} setActive={setStarred} onClick={updateStar} />
                    {!props.noDrag ? (
                        <div className="text-right flex items-center justify-end">
                            <DragHamburger />
                        </div>
                    ) : (
                        <div className="w-2"></div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default ProjectEntry;
