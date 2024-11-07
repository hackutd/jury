import { useEffect, useState } from 'react';
import ProjectEntry from '../ProjectEntry';

interface StarListProps {
    judge: Judge | null;
}

const StarList = ({ judge }: StarListProps) => {
    const [loaded, setLoaded] = useState(false);
    const [starred, setStarred] = useState<SortableJudgedProject[]>([]);
    const [unstarred, setUnstarred] = useState<SortableJudgedProject[]>([]);

    // Load all projects when judge loads
    useEffect(() => {
        if (!judge) return;

        const allProjects = judge.seen_projects
            .map((p, i) => ({
                id: i + 1,
                ...p,
            }))
            .reverse() as SortableJudgedProject[];

        const newStarred = allProjects.filter((p) => p.starred);
        const newUnstarred = allProjects.filter((p) => !p.starred);

        setStarred(newStarred);
        setUnstarred(newUnstarred);
        setLoaded(true);
    }, [judge]);

    const handleStarChange = (id: number, isStarred: boolean) => {
        // Figure out which list the project is in
        const starredIndex = starred.findIndex((p) => p.id === id);
        const unstarredIndex = unstarred.findIndex((p) => p.id === id);

        // Update the project in the correct list
        if (isStarred) {
            const project = unstarred[unstarredIndex];
            project.starred = true;
            setStarred([...starred, project].sort((a, b) => b.id - a.id));
            setUnstarred(unstarred.filter((p) => p.id !== id));
        } else {
            const project = starred[starredIndex];
            project.starred = false;
            setUnstarred([...unstarred, project].sort((a, b) => b.id - a.id));
            setStarred(starred.filter((p) => p.id !== id));
        }
    };

    if (!loaded) {
        return <p>Loading...</p>;
    }

    return (
        <div>
            <h2 className="text-primary text-xl font-bold mt-4">Starred Projects</h2>
            <p className="text-light text-sm">
                You may star as many projects as you like. Projects are sorted in reverse
                chronological order.
            </p>
            <div className="h-[1px] w-full bg-light my-2"></div>
            {starred.map((p) => (
                <ProjectEntry
                    id={p.id}
                    key={p.id}
                    project={p}
                    starCallback={handleStarChange}
                    ranking={-1}
                    noDrag
                />
            ))}
            <h2 className="text-primary text-xl font-bold mt-4">Unstarred Projects</h2>
            <p className="text-light text-sm">
                Click the star to star projects that you think will win or stand out.
            </p>
            <div className="h-[1px] w-full bg-light my-2"></div>
            {unstarred.map((p) => (
                <ProjectEntry
                    id={p.id}
                    key={p.id}
                    project={p}
                    starCallback={handleStarChange}
                    ranking={-1}
                    noDrag
                />
            ))}
        </div>
    );
};

export default StarList;
