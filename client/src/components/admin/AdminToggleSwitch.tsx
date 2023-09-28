import React from 'react';

const AdminToggleSwitch = (props: {
    state: boolean;
    setState: React.Dispatch<React.SetStateAction<boolean>>;
}) => {
    return (
        <div className="w-full flex flex-row items-center justify-center my-10">
            <div className="bg-primaryLight flex flex-row rounded-full relative">
                <div
                    className={`absolute top-0 left-0 w-44 h-16 transition-transform duration-300 ease-in-out ${
                        !props.state ? 'transform translate-x-44' : ''
                    }`}
                >
                    <button
                        className="block w-full h-full text-2xl rounded-full bg-primary text-white"
                        onClick={() => {
                            props.setState(!props.state);
                        }}
                    >
                        {props.state ? 'Projects' : 'Judges'}
                    </button>
                </div>

                <button
                    className={`block w-44 h-16 text-2xl rounded-full ${
                        props.state ? '' : 'bg-transparent text-light'
                    }`}
                    onClick={() => {
                        if (!props.state) props.setState(!props.state);
                    }}
                >
                    Projects
                </button>
                <button
                    className={`block w-44 h-16 text-2xl rounded-full ${
                        !props.state ? '' : 'bg-transparent text-light'
                    }`}
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
