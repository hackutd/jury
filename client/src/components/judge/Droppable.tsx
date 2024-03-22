import { useDroppable } from '@dnd-kit/core';

interface DroppableProps {
    id: string;
    children: React.ReactNode;
}

const Droppable = (props: DroppableProps) => {
    const { isOver, setNodeRef } = useDroppable({
        id: props.id,
    });
    const style = {
        backgroundColor: isOver ? 'lightblue' : 'transparent',
    };

    return (
        <div ref={setNodeRef} style={style}>
            {props.children}
        </div>
    );
};

export default Droppable;
