import { useEffect, useState } from 'react';
import VotePopupButton from './VotePopupButton';
import Button from '../Button';
import { twMerge } from 'tailwind-merge';
import { getRequest } from '../../api';
import { errorAlert } from '../../util';

interface VotePopupProps {
    /* Type of popup to show */
    popupType: VotePopupState;

    /* Function to modify the popup state variable */
    close: React.Dispatch<React.SetStateAction<boolean>>;

    /* Judge to vote on */
    judge: Judge;

    /* Function to call to vote */
    voteFunc: (scores: { [category: string]: number }) => void;

    /* Function to call to flag */
    flagFunc: (choice: number) => void;

    /* Function to call to skip */
    skipFunc: () => void;

    /* State variable for determining if popup is open */
    open: boolean;
}

type PopupText = {
    [key in VotePopupState]: PopupTextObj[];
};

interface PopupTextObj {
    text: string;
    subtext: string;
}

const popupText: PopupText = {
    vote: [],
    flag: [
        {
            text: 'Absent',
            subtext: 'Team is not present for judging at their table',
        },
        {
            text: 'Cannot Demo Project',
            subtext: 'Team cannot prove they made the project',
        },
        {
            text: 'Too Complex',
            subtext: 'Appears too complex to be made in hackathon',
        },
        {
            text: 'Offensive Project',
            subtext: 'Offensive or breaks Code of Conduct',
        },
    ],
    busy: [],
};

const VotePopup = (props: VotePopupProps) => {
    const [categories, setCategories] = useState<string[]>([]);
    const [categoryScores, setCategoryScores] = useState<number[]>([]);
    const [selected, setSelected] = useState(-1);

    useEffect(() => {
        if (props.popupType !== 'vote') return;

        async function getCategories() {
            // Get the categories
            const res = await getRequest<string[]>(`/categories`, 'judge');
            if (res.status !== 200) {
                errorAlert(res);
                return;
            }

            const cats = res.data ?? [];
            const catScores = cats.map(() => 0);

            setCategories(cats);
            setCategoryScores(catScores);
        }

        getCategories();
    }, [props.popupType]);

    useEffect(() => {
        setSelected(-1);
    }, [props.popupType]);

    if (!props.open) return null;

    const handleClick = () => {
        if (selected === -1 && props.popupType === 'flag') {
            alert('Please select an option.');
            return;
        }

        props.close(false);

        // Create scores object if voting
        // Scores object should be { [category]: score }
        if (props.popupType === 'vote') {
            const scores = categories.map((v, i) => ({ [v]: categoryScores[i] })).reduce((a, b) => ({ ...a, ...b }), {});
            props.voteFunc(scores);
            return;
        }

        if (props.popupType === 'busy') {
            props.skipFunc();
            return;
        }

        props.flagFunc(selected);
    };

    const titleColor =
        props.popupType === 'vote'
            ? 'text-primary'
            : props.popupType === 'flag'
            ? 'text-error'
            : 'text-gold';

    const title: PopupTextObj =
        props.popupType === 'vote'
            ? { text: 'Judge Project', subtext: 'Please rate the current project' }
            : props.popupType === 'flag'
            ? { text: 'Flag Project', subtext: 'Select a flag reason' }
            : {
                  text: 'Project is Busy',
                  subtext:
                      'Confirm that the other team is being judged by another judge. If so, click below to skip this project for now.',
              };

    return (
        <>
            <div
                className="fixed left-0 top-0 z-20 w-screen h-screen bg-black/30"
                onClick={() => props.close(false)}
            ></div>
            <div className="bg-background fixed z-30 left-1/2 top-1/2 translate-x-[-50%] translate-y-[-50%] px-6 py-6 md:px-10 md:w-1/3 w-11/12 flex flex-col items-center">
                <h1 className={twMerge('text-3xl font-bold', titleColor)}>{title.text}</h1>
                <h2
                    className={twMerge(
                        'text-xl',
                        props.popupType === 'busy' ? 'text-center mt-4' : 'font-bold'
                    )}
                >
                    {title.subtext}
                </h2>
                <div className="flex flex-col items-center w-full my-4">
                    {props.popupType === 'vote' ? (
                        <div className="flex flex-col align-center">
                            {categories.map((v, i) => (
                                <div key={i}>
                                    <p className="text-center">
                                        <b>{v}</b>: {categoryScores[i]}
                                    </p>
                                    <input
                                        type="range"
                                        min="0"
                                        max="10"
                                        value={categoryScores[i]}
                                        onChange={(e) => {
                                            const newScores = [...categoryScores];
                                            newScores[i] = parseInt(e.target.value);
                                            setCategoryScores(newScores);
                                        }}
                                        className="w-full"
                                    />
                                </div>
                            ))}
                        </div>
                    ) : (
                        popupText &&
                        popupText[props.popupType].map((v, i) => (
                            <VotePopupButton
                                key={i}
                                onClick={() => {
                                    setSelected(i);
                                }}
                                text={v.text}
                                subtext={v.subtext}
                                type={props.popupType}
                                selected={selected === i}
                            />
                        ))
                    )}
                </div>
                <Button type="primary" onClick={handleClick}>
                    {props.popupType === 'busy' ? 'Skip' : 'Submit'}
                </Button>
            </div>
        </>
    );
};

export default VotePopup;
