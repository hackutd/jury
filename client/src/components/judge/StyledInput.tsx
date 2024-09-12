import React from 'react';
import { getTrackBackground, Range } from 'react-range';

export default function StyledInput() {
    const [values, setValues] = React.useState([5]);

    return (
        <Range
            label="Select your value"
            step={1}
            min={0}
            max={10}
            values={values}
            onChange={(values) => setValues(values)}
            renderMark={({ props, index }) => (
                <div
                    {...props}
                    key={props.key}
                    style={{
                        ...props.style,
                        height: '4px',
                        width: '4px',
                        opacity: index == 0 || index == 10 ? 0 : 1,
                        backgroundColor: index * 1 < values[0] ? '#00ACE6' : '#646464',

                        borderRadius: '50%',
                        marginTop: '1px',
                    }}
                />
            )}
            renderTrack={({ props, children }) => (
                <div
                    onMouseDown={props.onMouseDown}
                    onTouchStart={props.onTouchStart}
                    style={{
                        ...props.style,
                        height: '10px',
                        display: 'flex',
                        width: '100%',
                    }}
                >
                    <div
                        ref={props.ref}
                        style={{
                            height: '6px',
                            width: '100%',
                            borderRadius: '5px',
                            background: getTrackBackground({
                                values,
                                colors: ['#00ACE6', '#ccc'],
                                min: 0,
                                max: 10,
                            }),
                            alignSelf: 'center',
                        }}
                    >
                        {children}
                    </div>
                </div>
            )}
            renderThumb={({ props }) => (
                <div
                    {...props}
                    key={props.key}
                    style={{
                        ...props.style,
                        height: '20px',
                        width: '20px',
                        backgroundColor: '#00ACE6',
                        borderRadius: '50%',
                    }}
                />
            )}
        />
    );
}
