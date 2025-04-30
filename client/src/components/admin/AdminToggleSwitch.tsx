import React from 'react';
import { twMerge } from 'tailwind-merge';

const AdminToggleSwitch = (props: {
    state: boolean;
    setState: React.Dispatch<React.SetStateAction<boolean>>;
}) => {
    const size = 'lg:w-44 lg:h-16 w-32 h-11 lg:text-2xl text-xl';
    return (
        <div className="w-full flex flex-row items-center justify-center lg:my-10 my-4">
            <div className="bg-primaryLight flex flex-row rounded-xl relative">
                <div
                    className={twMerge(
                        'absolute top-0 left-0 transition-transform duration-300 ease-in-out',
                        size,
                        !props.state ? 'transform lg:translate-x-44 translate-x-32' : ''
                    )}
                >
                    <div className="flex flex-col items-center justify-center w-full h-full lg:text-2xl text-xl rounded-xl bg-primary text-white">
                        {props.state ? 'Projects' : 'Judges'}
                    </div>
                </div>

                <button
                    className={twMerge(
                        'block rounded-xl',
                        size,
                        props.state ? '' : 'bg-transparent text-light'
                    )}
                    onClick={() => {
                        if (!props.state) props.setState(!props.state);
                    }}
                >
                    Projects
                </button>
                <button
                    className={twMerge(
                        'block rounded-xl',
                        size,
                        !props.state ? '' : 'bg-transparent text-light'
                    )}
                    onClick={() => {
                        if (props.state) props.setState(!props.state);
                    }}
                >
                    Judges
                </button>
            </div>
        </div>
    );
};

export default AdminToggleSwitch;
