import { useEffect, useState } from 'react';
import { getRequest, postRequest, putRequest } from '../../api';
import { errorAlert } from '../../util';
import Button from '../Button';
import { getTrackBackground, Range } from 'react-range';

interface RatingsProps {
    callback?: () => void;
    submitText?: string;
    prior?: { [x: string]: number }; // TODO: wtf is this type
    small?: boolean;
    update?: boolean;
    project?: JudgedProject;
}

const Ratings = (props: RatingsProps) => {
    const [categories, setCategories] = useState<string[]>([]);
    const [categoryScores, setCategoryScores] = useState<number[]>([]);

    useEffect(() => {
        if (!props.prior || categories.length === 0) return;

        const newScores = categories.map((v) => (props.prior as any)[v] ?? 0); // TODO: fix this
        setCategoryScores(newScores);
    }, [props.prior, categories]);

    // Get categories from the options
    useEffect(() => {
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
    }, []);

    // Submit the scores
    const submit = async (newScores: number[]) => {
        const scoresToUse = props.update ? newScores : categoryScores;

        // Create the scores object
        const scores = categories
            .map((v, i) => ({ [v]: scoresToUse[i] }))
            .reduce((a, b) => ({ ...a, ...b }), {});

        // Score the current project
        const scoreRes = props.update
            ? await putRequest<OkResponse>('/judge/score', 'judge', {
                  categories: scores,
                  project: props.project?.project_id,
              })
            : await postRequest<OkResponse>('/judge/score', 'judge', {
                  categories: scores,
                  initial: true,
              });
        if (scoreRes.status !== 200) {
            errorAlert(scoreRes);
            return;
        }

        if (props.callback) props.callback();
    };

    return (
        <div className="flex flex-col mx-4">
            {categories.map((v, i) => (
                <div key={i}>
                    <div className="text-light text-mono flex flex-row justify-between">
                        <p>{v}</p>
                        <p>{categoryScores[i]}</p>
                    </div>
                    <div className="pb-4 pt-1">
                        {/* React Range Slider */}
                        <Range
                            label="Select your value"
                            step={1}
                            min={0}
                            max={10}
                            values={[categoryScores[i]]}
                            onChange={(values) => {
                                const newScores = [...categoryScores];
                                newScores[i] = values[0];
                                setCategoryScores(newScores);
                                if (props.update) submit(newScores);
                            }}
                            renderMark={({ props, index }) => (
                                <div
                                    {...props}
                                    key={props.key}
                                    className={`h-0.5 w-0.5 rounded-full mt-1 ${
                                        index === 0 || index === 10 ? 'opacity-0' : 'opacity-100'
                                    } ${index < categoryScores[i] ? 'bg-primary' : 'bg-black'}`}
                                />
                            )}
                            renderTrack={({ props, children }) => {
                                return (
                                    <div
                                        onMouseDown={props.onMouseDown}
                                        onTouchStart={props.onTouchStart}
                                        className="h-6 flex w-full"
                                        style={{ ...props.style }}
                                    >
                                        <div
                                            ref={props.ref}
                                            style={{
                                                background: getTrackBackground({
                                                    values: [categoryScores[i]],
                                                    colors: ['#00ACE6', '#ccc'],
                                                    min: 0,
                                                    max: 10,
                                                }),
                                            }}
                                            className={`h-[6px] w-full rounded-lg self-center`}
                                        >
                                            {children}
                                        </div>
                                    </div>
                                );
                            }}
                            renderThumb={({ props }) => (
                                <div
                                    {...props}
                                    key={props.key}
                                    className="h-5 w-5 bg-primary rounded-full"
                                    style={{ ...props.style }}
                                />
                            )}
                        />
                    </div>
                </div>
            ))}
            {!props.update ? (
                <div className="flex justify-center mt-4">
                    <Button
                        type="primary"
                        onClick={submit.bind(this, [])}
                        className={props.small ? 'p-1 mb-4' : 'mb-4'}
                    >
                        {props.submitText ?? 'Submit'}
                    </Button>
                </div>
            ) : (
                <div className="h-4 w-full"></div>
            )}
        </div>
    );
};

export default Ratings;
