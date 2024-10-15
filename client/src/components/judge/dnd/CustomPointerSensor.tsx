import { PointerSensor, PointerSensorOptions } from '@dnd-kit/core';
import { PointerEvent } from 'react';

class CustomPointerSensor extends PointerSensor {
    // I hate typescript
    static activators: {
        eventName: 'onPointerDown';
        handler: (
            { nativeEvent }: PointerEvent<Element>,
            { onActivation }: PointerSensorOptions
        ) => boolean;
    }[] = [
        {
            eventName: 'onPointerDown',
            handler: ({ nativeEvent: event }) => {
                return (
                    event.isPrimary &&
                    event.button === 0 &&
                    (event.target as HTMLElement).classList.contains('drag-handle')
                );
            },
        },
    ];
}

export default CustomPointerSensor;
