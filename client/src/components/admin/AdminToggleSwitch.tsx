import React from 'react';
import { twMerge } from 'tailwind-merge';

const AdminToggleSwitch = (props: {
    state: boolean;
    setState: React.Dispatch<React.SetStateAction<boolean>>;
}) => {
    return (
        <div className="w-full flex flex-row items-center justify-center my-10">
            <button
                className={twMerge("block px-4 py-2 text-3xl rounded-full", props.state ? "bg-primary text-white" : "bg-backgroundDark")}
                onClick={() => {
                    props.setState(true);
                }}
            >
                Projects
            </button>
            <button
                className={twMerge("block px-4 py-2 text-3xl rounded-full", !props.state ? "bg-primary text-white" : "bg-backgroundDark")}
                onClick={() => {
                    props.setState(false);
                }}
            >
                Judges
            </button>
        </div>
    );
};

export default AdminToggleSwitch;
