import { twMerge } from 'tailwind-merge';

const StatBlock = (props: { name: string; value: string | number; className?: string }) => {
    function fixIfFloat(n: number): string {
        if (Math.round(n) === n) {
            return n.toString();
        }
        return n.toFixed(3);
    } 

    return (
        <div className={twMerge('text-center', props.className)}>
            <div className="text-5xl lg:text-6xl">
                {typeof props.value === 'number' ? fixIfFloat(props.value) : props.value}
            </div>
            <div className="text-light text-sm lg:text-md">{props.name}</div>
        </div>
    );
};

export default StatBlock;
