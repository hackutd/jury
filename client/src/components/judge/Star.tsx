import star from '../../assets/star.svg';
import starGrey from '../../assets/star-grey.svg';
import { twMerge } from 'tailwind-merge';
import { postRequest } from '../../api';
import { errorAlert } from '../../util';

interface StarProps {
    active: boolean;
    clickable: boolean;
    num: number;
    project_id: string;
    update: React.Dispatch<React.SetStateAction<number>>;
}

const Star = (props: StarProps) => {
    const handleClick: React.MouseEventHandler<HTMLDivElement> = async (e) => {
        if (!props.clickable) return;

        // Update star count based on the star # clicked
        const res = await postRequest('/judge/stars', 'judge', {
            stars: props.num,
            project_id: props.project_id,
        });
        if (res.status !== 200) {
            errorAlert(res.status);
            return;
        }

        // Update star count
        props.update(props.num);
    };

    return (
        <div
            className={twMerge('w-6 h-6', props.clickable ? 'cursor-pointer' : '')}
            onClick={handleClick}
        >
            <img src={props.active ? star : starGrey} alt="star icon" />
        </div>
    );
};

export default Star;
