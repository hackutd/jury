import React from 'react';

/**
 * This is the main container for the site, which is optimized for mobile.
 * On desktop devices, a set width will be imposed to mimic mobile displays.
 */
const Container = (props: { children: React.ReactNode }) => {
    return (
        <div className="flex flex-col justify-center items-center grow shrink-0 my-0 mx-auto md:w-[30rem]">
            {props.children}
        </div>
    );
};

export default Container;
