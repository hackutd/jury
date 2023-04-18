import { twMerge } from 'tailwind-merge';

const AddJudgeStat = (props: { name: string; value: string | number; className?: string }) => {
    return (
        <div className={twMerge('text-center', props.className)}>
            <div className="text-4xl md:text-6xl">{props.value}</div>
            <div className="text-light text-lg md:text-2xl">{props.name}</div>
        </div>
    );
};

export default AddJudgeStat;
