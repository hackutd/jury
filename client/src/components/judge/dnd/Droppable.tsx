import { useDroppable } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { twMerge } from 'tailwind-merge';
import SortableItem from './SortableItem';

interface DroppableProps {
    id: string;
    projects: SortableJudgedProject[];
    children?: React.ReactNode;
}

const Droppable = (props: DroppableProps) => {
    const { isOver, setNodeRef } = useDroppable({ id: props.id });
    const style = isOver ? 'bg-primary' : '';

    return (
        <SortableContext
            id={props.id}
            items={props.projects}
            strategy={verticalListSortingStrategy}
        >
            <div ref={setNodeRef} className={twMerge('bg-backgroundDark min-h-16 mb-4', style)}>
                {props.projects.map((p) => (
                    <SortableItem key={p.id} item={p} />
                ))}
            </div>
        </SortableContext>
    );
};

export default Droppable;
