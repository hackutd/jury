import {
    DndContext,
    DragEndEvent,
    DragOverEvent,
    DragOverlay,
    DragStartEvent,
    KeyboardSensor,
    UniqueIdentifier,
    closestCenter,
    useSensor,
    useSensors,
} from '@dnd-kit/core';
import Droppable from './Droppable';
import RankItem from './RankItem';
import CustomPointerSensor from './CustomPointerSensor';
import { arrayMove, sortableKeyboardCoordinates } from '@dnd-kit/sortable';
import { useEffect, useState } from 'react';
import { getRequest, postRequest } from '../../../api';
import { errorAlert } from '../../../util';

interface RankingProps {
    judge: Judge | null;
}

const Ranking = ({ judge }: RankingProps) => {
    const [activeId, setActiveId] = useState<number | null>(null);
    const [activeDropzone, setActiveDropzone] = useState<string | null>(null);
    const [disabled, setDisabled] = useState(false);
    const sensors = useSensors(
        useSensor(CustomPointerSensor, {
            activationConstraint: {
                distance: 5,
            },
        }),
        useSensor(KeyboardSensor, {
            coordinateGetter: sortableKeyboardCoordinates,
        })
    );
    const [ranked, setRanked] = useState<SortableJudgedProject[]>([]);
    const [unranked, setUnranked] = useState<SortableJudgedProject[]>([]);
    const [deliberation, setDeliberation] = useState(false);
    const [loaded, setLoaded] = useState(false);

    // Load all projects when judge loads
    useEffect(() => {
        async function fetchData() {
            if (!judge) return;

            const allProjects = judge.seen_projects.map((p, i) => ({
                id: i + 1,
                ...p,
            }));

            const rankedProjects = judge.rankings.map((r) =>
                allProjects.find((p) => p.project_id === r)
            ) as SortableJudgedProject[];
            const unrankedProjects = allProjects.filter((p) =>
                judge.rankings.every((r) => r !== p.project_id)
            );
            unrankedProjects.reverse();

            // Get deliberation status
            const delibRes = await getRequest<OkResponse>('/judge/deliberation', 'judge');
            if (delibRes.status !== 200) {
                errorAlert(delibRes);
                return;
            }
            setDeliberation(delibRes.data?.ok === 1);

            setRanked(rankedProjects);
            setUnranked(unrankedProjects);
            setLoaded(true);
        }

        fetchData();
    }, [judge]);

    const handleDragStart = (event: DragStartEvent) => {
        const { active } = event;
        setActiveId(active.id as number);
    };

    const handleDragOver = (event: DragOverEvent) => {
        const { active, over } = event;
        const { id } = active;

        if (over === null) {
            setActiveId(null);
            return;
        }
        const { id: overId } = over;

        const activeRanked = isRankedObject(id);
        const overRanked = isRankedObject(overId);

        setActiveDropzone(overRanked ? 'ranked' : 'unranked');

        // If dragging from unranked to ranked and ranked has 5 items already, don't allow
        if (!activeRanked && overRanked && ranked.length >= 5) {
            setDisabled(true);
            return;
        }
        setDisabled(false);

        // If moving to new container, swap the item to the new list
        if (activeRanked !== overRanked) {
            const activeContainer = activeRanked ? ranked : unranked;
            const overContainer = overRanked ? ranked : unranked;
            const oldIndex = activeContainer.findIndex((i) => i.id === active.id);
            const newIndex = overContainer.findIndex((i) => i.id === over.id);
            const proj = activeContainer[oldIndex];
            const newActive = activeContainer.toSpliced(oldIndex, 1);
            const newOver = overContainer.toSpliced(newIndex, 0, proj);
            if (activeRanked) {
                setRanked(newActive);
                setUnranked(newOver);
            } else {
                setRanked(newOver);
                setUnranked(newActive);
            }
        }
    };

    const handleDragEnd = (event: DragEndEvent) => {
        const { active, over } = event;
        const { id } = active;

        if (over === null) {
            setActiveId(null);
            return;
        }
        const { id: overId } = over;

        const activeRanked = isRankedObject(id);
        const overRanked = isRankedObject(overId);

        if (activeRanked === overRanked) {
            const currProjs = activeRanked ? ranked : unranked;

            const oldIndex = currProjs.findIndex((i) => i.id === active.id);
            const newIndex = currProjs.findIndex((i) => i.id === over.id);
            const newProjects: SortableJudgedProject[] = arrayMove(currProjs, oldIndex, newIndex);
            activeRanked ? setRanked(newProjects) : setUnranked(newProjects);

            if (activeRanked) saveSort(newProjects);
            else saveSort(ranked);
        } else {
            saveSort(ranked);
        }

        setActiveDropzone(null);
        setActiveId(null);
        setDisabled(false);
    };

    // dnd-kit is strange. For active/over ids, it is a number most of the time,
    // representing the ID of the item that we are hovering over.
    // However, if the user is hovering NOT on an item, it will set the ID
    // to the ID of the droppable container ?!??!
    // Strange indeed.
    function isRankedObject(id: UniqueIdentifier) {
        // If drop onto the zone (id would be string)
        if (isNaN(Number(id))) {
            return id === 'ranked';
        }

        // Otherwise if dropped onto a specific object
        const ro = ranked.find((a) => a.id === id);
        return !!ro;
    }

    const saveSort = async (projects: SortableJudgedProject[]) => {
        // Save the rankings
        const saveRes = await postRequest<OkResponse>('/judge/rank', 'judge', {
            ranking: projects.map((p) => p.project_id),
        });
        if (saveRes.status !== 200) {
            // If deliberation, reload page
            if (saveRes.error.indexOf('deliberation') !== -1) {
                window.location.reload();
            }

            errorAlert(saveRes);
            return;
        }
    };

    if (!loaded) {
        return <p>Loading...</p>;
    }

    return (
        <DndContext
            sensors={sensors}
            collisionDetection={closestCenter}
            onDragStart={handleDragStart}
            onDragOver={handleDragOver}
            onDragEnd={handleDragEnd}
        >
            {deliberation && (
                <p className="p-2 text-center mt-4 m-2 border-gold border-2 bg-gold/20 rounded-md">
                    Deliberation has started. You cannot change rankings or stars.
                </p>
            )}
            <h2 className="text-primary text-xl font-bold mt-4">Ranked Projects</h2>
            <p className="text-light text-sm">
                Rank at most 5 projects. Click on titles to edit scores and see details.
            </p>
            <div className="h-[1px] w-full bg-light my-2"></div>
            <Droppable
                id="ranked"
                projects={ranked}
                active={activeDropzone}
                disabled={disabled || deliberation}
            />

            <h2 className="text-primary text-xl font-bold mt-4">Unranked Projects</h2>
            <p className="text-light text-sm">
                Projects will be sorted in reverse chronological order.
            </p>
            <div className="h-[1px] w-full bg-light my-2"></div>
            <Droppable
                id="unranked"
                projects={unranked}
                active={activeDropzone}
                disabled={deliberation}
            />

            <DragOverlay>
                {activeId ? (
                    <RankItem
                        item={
                            unranked.find((p) => p.id === activeId) ??
                            (ranked.find((p) => p.id === activeId) as SortableJudgedProject)
                        }
                        ranking={ranked.findIndex((p) => p.id === activeId) + 1}
                        disabled={deliberation}
                    />
                ) : null}
            </DragOverlay>
        </DndContext>
    );
};

export default Ranking;
