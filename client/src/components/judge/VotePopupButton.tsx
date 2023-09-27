import { twMerge } from 'tailwind-merge';
import Button from '../Button';

interface VotePopupButtonProps {
    type: VotePopupState;
    text: string;
    subtext: string;
    onClick: (e: React.MouseEvent<Element>) => void;
    selected: boolean;
}

const VotePopupButton = (props: VotePopupButtonProps) => {
    const selectedStyle =
        props.type === 'vote'
            ? 'text-primary bg-primary/10 border-primary'
            : props.type === 'flag'
            ? 'text-error bg-error/10 border-error'
            : 'text-gold bg-gold/10 border-gold';
    return (
        <Button
            type="outline"
            full
            square
            className={twMerge(
                'flex flex-col items-center py-2 my-2',
                props.selected ? selectedStyle : 'text-black'
            )}
            onClick={props.onClick}
        >
            <h3 className="text-2xl">{props.text}</h3>
            <p className="text-lg text-light">{props.subtext}</p>
        </Button>
    );
};

export default VotePopupButton;
