import { useNavigate } from 'react-router-dom';
import { twMerge } from 'tailwind-merge';

const Back = (props: { location: string, className?: string }) => {
    const navigate = useNavigate();
    const back = () => {
        navigate(props.location);
    };
    return (
        <div className={twMerge("cursor-pointer my-4 font-bold", props.className)} onClick={back}>
            {'<'} Back
        </div>
    );
};

export default Back;
