import { useEffect, useState } from 'react';
import { getRequest, postRequest, putRequest } from '../../api';
import { errorAlert } from '../../util';
import Button from '../Button';

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
                            if (props.update) submit(newScores);
                        }}
                        className="w-full"
                    />
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
