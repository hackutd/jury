import { useDraggable } from '@dnd-kit/core';

interface DraggableProps {
    id: string;
    children: React.ReactNode;
}

const Draggable = (props: DraggableProps) => {
    const { attributes, listeners, setNodeRef, transform } = useDraggable({
        id: props.id,
    });
    const style = transform ? {
        transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`,
    } : undefined;

    return (
        <button ref={setNodeRef} style={style} {...attributes} {...listeners}>
            {props.children}
        </button>
    )
};

export default Draggable;
