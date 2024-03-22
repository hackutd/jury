import { twMerge } from 'tailwind-merge';
import Button from '../Button';

interface FlagPopupButtonProps {
    text: string;
    subtext: string;
    onClick: (e: React.MouseEvent<Element>) => void;
    selected: boolean;
}

const FlagPopupButton = (props: FlagPopupButtonProps) => {
    return (
        <Button
            type="outline"
            full
            square
            className={twMerge(
                'flex flex-col items-center py-2 my-2',
                props.selected ? 'text-error bg-error/10 border-error' : 'text-black'
            )}
            onClick={props.onClick}
        >
            <h3 className="text-2xl">{props.text}</h3>
            <p className="text-lg text-light">{props.subtext}</p>
        </Button>
    );
};

export default FlagPopupButton;
