import { twMerge } from "tailwind-merge";

const AdminStat = (props: { name: string; value: string | number; className?: string }) => {
    return (
        <div className={twMerge("text-center", props.className)}>
            <div className="text-6xl">{props.value}</div>
            <div className="text-light">{props.name}</div>
        </div>
    );
};

export default AdminStat;
