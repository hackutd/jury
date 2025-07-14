---
sidebar_position: 3
title: Styling
description: How to style frontend components with Tailwind.
---

# Styling

All elements are to be styled with [Tailwind](https://tailwindcss.com). Tailwind lets us style components without using CSS, instead relying on classNames defined on each element.

## Configuration

Tailwind uses a `tailwind.config.js` file, located in the root directory. In this file, we define the colors, fonts, animations, and plugins that we use with Tailwind. The custom colors and fonts are standardized across Jury and should be used almost everywhere. All colors and fonts can also be viewed on the [Figma board](https://www.figma.com/file/qwBWs4i7pJMpFbcjMffDZU/Jury-(Gavel-Plus)?node-id=8%3A100&t=xYwfPwRAUeJw9jNr-1).

## Usage

We will refer to another example to demonstrate the usage of Tailwind. Here we will look at the `client/src/components/Back.tsx` component:

```jsx
 <div
    className={twMerge('cursor-pointer my-4 font-bold duration-200 hover:text-primary', props.className)}
    onClick={back}
>
    {'<'} Back
</div>
```

Notice we are using the `twMerge` function from the [tailwind merge](https://www.npmjs.com/package/tailwind-merge) library. This is used throughout the app to combine styles defined at different levels, such as `props.className` passed in through the component properties. `twMerge` will merge functions, replacing classes that **come later in the parameter list**. For example, if I have `twMerge('text-gold', 'text-primary')` as the className, the function will resolve to `'text-primary'`. The format above where we define the default styling and then merge it with `props.className` is very common among our user-defined components.

There are a ton of tailwind classes, with most pretty similar to their CSS attribute counterpart. The Tailwind website is extremely useful to finding these classes--though you can just google the CSS attribute + tailwind and the page should show up as well.
