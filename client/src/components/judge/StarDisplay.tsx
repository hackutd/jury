import Star from "./Star";

const StarDisplay = (props: { stars: number }) => {
    const stars = props.stars ? props.stars : 0;
    return (
        <>
            {Array.from(Array(stars).keys()).map((i) => (
                <Star key={i} active={true} />
            ))}
            {Array.from(Array(5 - stars).keys()).map((i) => (
                <Star key={i} active={false} />
            ))}
        </>
    );
};

export default StarDisplay;
