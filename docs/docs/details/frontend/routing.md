---
sidebar_position: 1
title: Routing
description: How routing works on the frontend application.
---

# Routing (Specifially React Routing)

Jury uses [React Router](https://reactrouter.com/en/main) to do routing. All routes are defined in the `client/src/index.tsx` page. Every component used on this page for the main pages should be placed in the `client/src/pages` directory. While the placement of the page components don't influence their actual location, it's good practice to place the pages in subdirectories corresponding to their path.

## Page Format

Generally, a page should look like the following. The `Container` component wraps the entire page to limit its width, used for all of the judging pages. The `Helmet` component sets the title, and the `JuryHeader` component displays the title of the app and some other things, based on props set (such as a logout button and a back button).

```jsx
import Container from '../components/Container';

const PageName = () => {
    return (
        <>
            <Helmet>
                <title>Title of Page | Jury</title>
            </Helmet>
            <JuryHeader withLogout />
            <Container>
                {/* Page content */}
            </Container>
        </>
    );
};

export default PageName;
```
