import { CSSProperties, forwardRef } from 'react';
import ProjectEntry from './ProjectEntry';

type RankItemProps = {
    item?: SortableJudgedProject; // Change to project
    isOpacityEnabled?: boolean;
    isDragging?: boolean;
} & React.HTMLAttributes<HTMLDivElement>;

const RankItem = forwardRef<HTMLDivElement, RankItemProps>(
    ({ item, isOpacityEnabled, isDragging, style, ...props }, ref) => {
        const styles: CSSProperties = {
            opacity: isOpacityEnabled ? 0.5 : 1,
            cursor: isDragging ? 'grabbing' : 'grab',
            transform: isDragging ? 'scale(1.05)' : 'scale(1)',
            ...style,
        };
        return (
            <div {...props} ref={ref} style={styles}>
                <ProjectEntry project={item} />
            </div>
        );
    }
);

export default RankItem;
