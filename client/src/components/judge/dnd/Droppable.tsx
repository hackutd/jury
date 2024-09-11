import { useDroppable } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { twMerge } from 'tailwind-merge';
import SortableItem from './SortableItem';

interface DroppableProps {
    id: string;
    projects: SortableJudgedProject[];
    active: string | null;
    children?: React.ReactNode;
}

const Droppable = (props: DroppableProps) => {
    const { isOver, setNodeRef } = useDroppable({ id: props.id });
    const style =
        isOver || props.active === props.id || props.projects.length === 0
            ? 'border-backgroundDark'
            : 'border-transparent';

    return (
        <SortableContext
            id={props.id}
            items={props.projects}
            strategy={verticalListSortingStrategy}
        >
            <div
                ref={setNodeRef}
                className={twMerge('border-dashed border-2 duration-100 min-h-16', style)}
            >
                {props.projects.map((p, i) => (
                    <SortableItem
                        key={p.id}
                        item={p}
                        ranking={props.id === 'ranked' ? i + 1 : -1}
                    />
                ))}
            </div>
        </SortableContext>
    );
};

export default Droppable;
