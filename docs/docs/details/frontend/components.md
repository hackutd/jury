---
sidebar_position: 2
title: Components
description: The component system and how to organize the frontend.
---

# Components

All components should be placed in the `client/src/components` directory. In that directory, the components not in a subdirectory are common shared components, such as `Button` and `Card`. These are used across the application.

Components specific to a page are placed in subdirectories, such as `client/src/components/admin/AdminToolbar.tsx`.

## Making a component

To create a component, refer to the common components. For example, here is the `Card` component:

```jsx
import React from 'react';
import { twMerge } from 'tailwind-merge';

interface ContainerProps {
    /* Element children */
    children?: React.ReactNode;

    /* Tailwind classes, will override defaults */
    className?: string;

    /* If true, don't justify center the entire page */
    noCenter?: boolean;
}

/**
 * This is the main container for the site, which is optimized for mobile.
 * On desktop devices, a set width will be imposed to mimic mobile displays.
 */
const Container = (props: ContainerProps) => {
    return (
        <div
            className={twMerge(
                'flex flex-col grow shrink-0 my-0 md:mx-auto md:w-[30rem] min-h-full',
                props.noCenter? '' : 'justify-center items-center',
                props.className
            )}
        >
            {props.children}
        </div>
    );
};

export default Container;
```

Props should be defined in a separate interface, with doc comments for each prop. If you are passing in children to the component, use `React.ReactNode` as the type. If you are passing in a Tailwind `className`, make sure to make it an optional string. Note that all props that are NOT optional will be required to be used otherwise an error will be thrown.

## Using Components

Refer to the official React docs on [how to use a component](https://react.dev/learn/your-first-component#using-a-component) if you don't know how to.
