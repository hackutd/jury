import star from '../../assets/star.svg';
import starGrey from '../../assets/star-grey.svg';

interface StarProps {
    active: boolean;
}

const Star = (props: StarProps) => {
    return (
        <div className="w-6 h-6">
            <img src={props.active ? star : starGrey} alt='star icon' />
        </div>
    );
};

export default Star;
