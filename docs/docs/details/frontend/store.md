---
sidebar_position: 6
title: Zustand Store
description: How to use the global store
---

# Zustand Store

[Zustand](https://zustand.docs.pmnd.rs/getting-started/introduction) is a small and scaleable app-wide state management solution. The key benefit to using Zustand is that we don't have to fetch all the data each time we navigate between pages on the frontend. Please read the documentation to learn how to use it--but if you have experience with React Redux or other global state management solutions, it should be pretty simple.

Note that most of the stores should only be used **on the admin side** as judges generally don't need/can't access that data.

## Defining Stores

All stores are defined in `client/src/store.tsx`. To define a **store**, use the `create` Zustand function and export it from the store file. Here is an example with the flags:

```js
interface FlagsStore {
    flags: Flag[];
    fetchFlags: () => Promise<void>;
}

const useFlagsStore = create<FlagsStore>((set) => ({
    flags: [],

    fetchFlags: async () => {
        const flagsRes = await getRequest<Flag[]>('/admin/flags', 'admin');
        if (flagsRes.status !== 200) {
            errorAlert(flagsRes);
            return;
        }
        set({ flags: flagsRes.data as Flag[] });
    },
}));
```

Note that we include the state (`flags`) and a state change function (`fetchFlags`). This way, we don't have to call the `fetchFlags` function if the flags don't change between components.

## Using Stores

To use the store, we simply import it as follows (this is a simplified `client/src/components/admin/FlagsPopup.tsx`):

```jsx
const FlagsPopup = (props: FlagsPopupProps) => {
    const flags = useFlagsStore((state) => state.flags);
    const fetchFlags = useFlagsStore((state) => state.fetchFlags);
    
    // ... other code

    const clearFlag = async (id: string) => {
        // ... other code

        fetchFlags();
    }

    return (
        <>
            {/* ... other code */}

            {flags.map((flag) => (
                <div
                    key={`${flag.id}`}
                    className="flex ..."
                >
                    {/* ... other code */}
                </div>
            ))}

            {/* ... other code */}
        </>
    )
}
```

Note that we are only fetching new flags where we clear a flag. Otherwise, the flag data should already be fetched by the admin dashboard when opening this popup.

## Defined Stores

There are a couple of existing defined stores in Jury (you shouldn't need to add any more unless there's a big change):

- **Admin Store**: Used to store admin dashboard stats, list of projects, list of judges, and project/judge stats
- **Clock Store**: Stores the main judging clock state
- **Options Store**: Stores the options defined in the admin settings
- **Flags Store**: Stores all flags in Jury
- **Admin Table Store**: Stores projects/judges, as well as the current state of the table
