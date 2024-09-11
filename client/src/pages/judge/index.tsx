import { useEffect, useState } from 'react';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';
import { useNavigate } from 'react-router-dom';
import Button from '../../components/Button';
import StatBlock from '../../components/StatBlock';
import Loading from '../../components/Loading';
import { getRequest, postRequest } from '../../api';
import { errorAlert } from '../../util';
import {
    DndContext,
    DragEndEvent,
    DragOverEvent,
    DragOverlay,
    DragStartEvent,
    KeyboardSensor,
    PointerSensor,
    UniqueIdentifier,
    closestCenter,
    useSensor,
    useSensors,
} from '@dnd-kit/core';
import { arrayMove, sortableKeyboardCoordinates } from '@dnd-kit/sortable';
import Droppable from '../../components/judge/dnd/Droppable';
import RankItem from '../../components/judge/dnd/RankItem';

const Judge = () => {
    const navigate = useNavigate();
    const [judge, setJudge] = useState<Judge | null>(null);
    const [ranked, setRanked] = useState<SortableJudgedProject[]>([]);
    const [unranked, setUnranked] = useState<SortableJudgedProject[]>([]);
    const [loaded, setLoaded] = useState(false);
    const [projCount, setProjCount] = useState(0);
    const [activeId, setActiveId] = useState<number | null>(null);
    const sensors = useSensors(
        useSensor(PointerSensor, {
            activationConstraint: {
                distance: 5,
            },
        }),
        useSensor(KeyboardSensor, {
            coordinateGetter: sortableKeyboardCoordinates,
        })
    );

    // Verify user is logged in and read welcome before proceeding
    useEffect(() => {
        async function fetchData() {
            // Check to see if the user is logged in
            const loggedInRes = await postRequest<OkResponse>('/judge/auth', 'judge', null);
            if (loggedInRes.status === 401) {
                console.error(`Judge is not logged in!`);
                navigate('/judge/login');
                return;
            }
            if (loggedInRes.status !== 200) {
                errorAlert(loggedInRes);
                return;
            }
            if (loggedInRes.data?.ok !== 1) {
                console.error(`Judge is not logged in!`);
                navigate('/judge/login');
                return;
            }

            // Check for read welcome
            const readWelcomeRes = await getRequest<OkResponse>('/judge/welcome', 'judge');
            if (readWelcomeRes.status !== 200) {
                errorAlert(readWelcomeRes);
                return;
            }
            const readWelcome = readWelcomeRes.data?.ok === 1;
            if (!readWelcome) {
                navigate('/judge/welcome');
            }

            // Get the name & email of the user from the server
            const judgeRes = await getRequest<Judge>('/judge', 'judge');
            if (judgeRes.status !== 200) {
                errorAlert(judgeRes);
                return;
            }
            const judge: Judge = judgeRes.data as Judge;
            setJudge(judge);

            // Get the project count
            const projCountRes = await getRequest<ProjectCount>('/project/count', 'judge');
            if (projCountRes.status !== 200) {
                errorAlert(projCountRes);
                return;
            }
            setProjCount(projCountRes.data?.count as number);
        }

        fetchData();
    }, []);

    // Load all projects when judge loads
    useEffect(() => {
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

        setRanked(rankedProjects);
        setUnranked(unrankedProjects);
        setLoaded(true);
    }, [judge]);

    if (!loaded) return <Loading disabled={!loaded} />;

    // Lets the user take a break
    const takeBreak = async () => {
        // Check if the user is allowed to take a break
        if (judge?.current == null) {
            alert('You are already taking a break!');
            return;
        }

        const res = await postRequest<OkResponse>('/judge/break', 'judge', null);
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert('You can now take a break! Press "Next project" to continue judging.');
    };

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

        console.log(over.id, active.id);

        const activeRanked = isRankedObject(id);
        const overRanked = isRankedObject(overId);

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

        console.log(over.id, active.id);

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

        setActiveId(null);
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
            errorAlert(saveRes);
            return;
        }
    };

    return (
        <>
            <JuryHeader withLogout />
            <Container noCenter className="px-2 pb-4">
                <h1 className="text-2xl my-2">Welcome, {judge?.name}!</h1>
                <div className="w-full mb-6">
                    <Button type="primary" full square href="/judge/live">
                        Next Project
                    </Button>
                    <div className="flex align-center justify-center mt-4">
                        <Button type="outline" square onClick={takeBreak} className="text-lg p-2">
                            I want to take a break!
                        </Button>
                    </div>
                </div>
                <div className="flex justify-evenly">
                    <StatBlock name="Seen" value={judge?.seen_projects.length as number} />
                    <StatBlock name="Total Projects" value={projCount} />
                </div>
                <DndContext
                    sensors={sensors}
                    collisionDetection={closestCenter}
                    onDragStart={handleDragStart}
                    onDragOver={handleDragOver}
                    onDragEnd={handleDragEnd}
                >
                    <h2 className="text-primary text-xl font-bold mt-4">Ranked Projects</h2>
                    <div className="h-[1px] w-full bg-light my-2"></div>
                    <Droppable id="ranked" projects={ranked} />

                    <h2 className="text-primary text-xl font-bold mt-4">Unranked Projects</h2>
                    <div className="h-[1px] w-full bg-light my-2"></div>
                    <Droppable id="unranked" projects={unranked} />

                    <DragOverlay>
                        {activeId ? (
                            <RankItem
                                item={
                                    unranked.find((p) => p.id === activeId) ??
                                    ranked.find((p) => p.id === activeId)
                                }
                            />
                        ) : null}
                    </DragOverlay>
                </DndContext>
            </Container>
        </>
    );
};

export default Judge;
