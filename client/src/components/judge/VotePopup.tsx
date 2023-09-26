import { useEffect, useState } from 'react';
import VotePopupButton from './VotePopupButton';
import Button from '../Button';
import { twMerge } from 'tailwind-merge';

interface VotePopupProps {
    /* Type of popup to show */
    popupType: VotePopupState;

    /* Function to modify the popup state variable */
    close: React.Dispatch<React.SetStateAction<boolean>>;

    /* Judge to vote on */
    judge: Judge;

    /* Function to call to vote */
    voteFunc: (choice: number) => void;

    /* Function to call to flag */
    flagFunc: (choice: number) => void;

    /* Function to call to skip */
    skipFunc: (choice: number) => void;

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

const VotePopup = (props: VotePopupProps) => {
    const [projectInfo, setProjectInfo] = useState<VotingProjectInfo | null>(null);
    const [popupText, setPopupText] = useState<PopupText | null>(null);
    const [selected, setSelected] = useState(-1);

    useEffect(() => {
        if (props.popupType !== 'vote') return;

        async function getProjectInfo() {
            // Gets the project info
            const piRes = await fetch(`${process.env.REACT_APP_JURY_URL}/judge/vote/info`, {
                method: 'GET',
                credentials: 'include',
            });

            if (piRes.status !== 200) {
                alert('Unable to fetch project info :(');
                return;
            }

            // Set project info
            setProjectInfo(await piRes.json());
        }

        getProjectInfo();
    }, [props.popupType]);

    useEffect(() => {
        if (!projectInfo) return;

        const pt: PopupText = {
            vote: [
                {
                    text: projectInfo.curr_name,
                    subtext: 'Current: Table ' + projectInfo.curr_location,
                },
                {
                    text: projectInfo.prev_name,
                    subtext: 'Previous: Table ' + projectInfo.prev_location,
                },
            ],
            flag: [
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
            skip: [
                {
                    text: 'Not Present',
                    subtext: 'Group is not present for judging',
                },
                {
                    text: 'Busy (Being Judged)',
                    subtext: 'Group is busy with another judge',
                },
            ],
        };

        setPopupText(pt);
    }, [projectInfo]);

    useEffect(() => {
        setSelected(-1);
    }, [props.popupType]);

    if (!props.open) return null;

    const handleClick = () => {
        props.close(false);
        props.popupType === 'vote'
            ? props.voteFunc(selected)
            : props.popupType === 'flag'
            ? props.flagFunc(selected)
            : props.skipFunc(selected);
    };

    const titleColor =
        props.popupType === 'vote'
            ? 'text-primary'
            : props.popupType === 'flag'
            ? 'text-error'
            : 'text-gold';

    const title: PopupTextObj =
        props.popupType === 'vote'
            ? { text: 'Judge Project', subtext: 'Select the better project' }
            : props.popupType === 'flag'
            ? { text: 'Flag Project', subtext: 'Select a flag reason' }
            : { text: 'Skip Project', subtext: 'Select a reason for skipping' };

    return (
        <>
            <div
                className="fixed left-0 top-0 z-20 w-screen h-screen bg-black/30"
                onClick={() => props.close(false)}
            ></div>
            <div className="bg-background fixed z-30 left-1/2 top-1/2 translate-x-[-50%] translate-y-[-50%] px-6 py-6 md:px-10 md:w-1/3 w-11/12 flex flex-col items-center">
                <h1 className={twMerge('text-3xl font-bold', titleColor)}>{title.text}</h1>
                <h2 className="text-xl font-bold">{title.subtext}</h2>
                <div className="flex flex-col items-center w-full my-4">
                    {popupText &&
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
                        ))}
                </div>
                <Button type="primary" onClick={handleClick}>
                    Submit
                </Button>
            </div>
        </>
    );
};

export default VotePopup;
