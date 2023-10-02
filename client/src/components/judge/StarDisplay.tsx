import { useState } from 'react';
import Star from './Star';

interface StarDisplayProps {
    stars: number;
    clickable?: boolean;
    id: string;
}

const StarDisplay = (props: StarDisplayProps) => {
    const [stars, setStars] = useState(props.stars || 0);

    return (
        <>
            {Array.from(Array(stars).keys()).map((i) => (
                <Star
                    key={i}
                    active={true}
                    clickable={props.clickable || false}
                    num={i + 1}
                    project_id={props.id}
                    update={setStars}
                />
            ))}
            {Array.from(Array(5 - stars).keys()).map((i) => (
                <Star
                    key={stars + i}
                    active={false}
                    clickable={props.clickable || false}
                    num={stars + i + 1}
                    project_id={props.id}
                    update={setStars}
                />
            ))}
        </>
    );
};

export default StarDisplay;
