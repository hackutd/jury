import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import RankItem from './RankItem';

interface SortableItemProps {
    item: SortableJudgedProject;
    ranking: number;
    children?: React.ReactNode;
}

const SortableItem = (props: SortableItemProps) => {
    if (!props.item) {
        return null;
    }
    const { attributes, isDragging, listeners, setNodeRef, transform, transition } = useSortable({
        id: props.item.id,
    });

    const style = {
        transform: CSS.Transform.toString(transform),
        transition,
        touchAction: 'initial',
    };

    return (
        <RankItem
            ref={setNodeRef}
            style={style}
            isOpacityEnabled={isDragging}
            {...props}
            {...attributes}
            {...listeners}
        >
            {props.children}
        </RankItem>
    );
};

export default SortableItem;
